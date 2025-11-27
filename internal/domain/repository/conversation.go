package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neokofg/callap-backend/internal/domain/entity"
	"github.com/neokofg/callap-backend/internal/domain/repository/utils"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
)

type ConversationRepository struct {
	pool                  *pgxpool.Pool
	rdb                   *redis.Client
	tableName             string
	participantsTableName string
	userTableName         string
}

func NewConversationRepository(pool *pgxpool.Pool, rdb *redis.Client) *ConversationRepository {
	return &ConversationRepository{
		pool:                  pool,
		rdb:                   rdb,
		tableName:             conversationTableName,
		participantsTableName: conversationParticipantsTableName,
		userTableName:         userTableName,
	}
}

func (cr *ConversationRepository) ListMessages(c context.Context, userId string, id string, limit int, offset int) ([]entity.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := cr.pool.Query(c, `
	SELECT 
            m.id,
            m.sender_id,
            u.name AS sender_name,
            m.content,
            m.created_at,
            m.is_read
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.conversation_id = $1
          AND EXISTS (
              SELECT 1 FROM conversation_participants cp 
              WHERE cp.conversation_id = $1 AND cp.user_id = $2 AND cp.left_at IS NULL
          )
        ORDER BY m.created_at DESC
        LIMIT $3 OFFSET $4
	`, id, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.SenderName, &msg.Content, &msg.CreatedAt, &msg.IsRead)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (cr *ConversationRepository) DeleteMessage(c context.Context, userId string, id string) error {
	tx, err := cr.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
		DELETE FROM messages WHERE id = $1 AND sender_id = $2;
	`, id, userId)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func (cr *ConversationRepository) NewMessage(c context.Context, userId string, id string, content string) error {
	tx, err := cr.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
        INSERT INTO messages (id, conversation_id, sender_id, content) 
        VALUES ($1, $2, $3, $4)
    `, ulid.Make().String(), id, userId, content)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, `
        UPDATE conversations SET updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE id = $1
    `, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, `
        UPDATE conversation_participants 
        SET left_at = NULL 
        WHERE conversation_id = $1 AND left_at IS NOT NULL
		`, id)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func (cr *ConversationRepository) Hide(c context.Context, userId string, id string) error {
	tx, err := cr.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	_, err = tx.Exec(c, `
		UPDATE 
		    conversation_participants cp
		SET left_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
		WHERE cp.conversation_id = $1 AND cp.user_id = $2
	`, id, userId)
	if err != nil {
		return err
	}

	return tx.Commit(c)
}

func (cr *ConversationRepository) GetConversationByID(c context.Context, userId string, convId string) (*entity.ConversationDetails, error) {
	row := cr.pool.QueryRow(c, `
        SELECT 
            c.id,
            c.type,
            c.created_at,
            c.updated_at,
            array_agg(cp.user_id) OVER () AS participants
        FROM conversations c
        JOIN conversation_participants cp ON c.id = cp.conversation_id
        WHERE c.id = $1
          AND EXISTS (SELECT 1 FROM conversation_participants cp2 WHERE cp2.conversation_id = c.id AND cp2.user_id = $2)
        GROUP BY c.id, c.type, c.created_at, c.updated_at
    `, convId, userId)
	var cd entity.ConversationDetails
	var participants []string
	err := row.Scan(&cd.ID, &cd.Type, &cd.CreatedAt, &cd.UpdatedAt, &participants)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, err
	}

	cd.Participants = participants

	return &cd, nil
}

func (cr *ConversationRepository) List(c context.Context, userId string, limit int, offset int) ([]entity.ConversationSummary, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := cr.pool.Query(c, `
        SELECT 
            c.id AS conversation_id,
            mp.user_id AS other_user_id,
            u.name AS other_user_name,
            u.tag AS other_user_tag,
            m.content AS last_message,
            m.created_at AS last_message_at,
            COALESCE(unread.unread_count, 0) AS unread_count
        FROM conversations c
        JOIN conversation_participants cp ON c.id = cp.conversation_id
        JOIN conversation_participants mp ON c.id = mp.conversation_id AND mp.user_id != $1
        JOIN users u ON mp.user_id = u.id
        LEFT JOIN (
            SELECT 
                conversation_id,
                content,
                created_at,
                ROW_NUMBER() OVER (PARTITION BY conversation_id ORDER BY created_at DESC) as rn
            FROM messages
        ) m ON c.id = m.conversation_id AND m.rn = 1
        LEFT JOIN (
            SELECT 
                m.conversation_id,
                COUNT(*) AS unread_count
            FROM messages m
            JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id
            WHERE m.sender_id != $1 AND m.is_read = FALSE AND cp.user_id = $1
            GROUP BY m.conversation_id
        ) unread ON c.id = unread.conversation_id
        WHERE cp.user_id = $1 AND cp.left_at IS NULL
        ORDER BY c.updated_at DESC
        LIMIT $2 OFFSET $3
    `, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []entity.ConversationSummary
	for rows.Next() {
		var cs entity.ConversationSummary
		var lastAt *time.Time
		var lastMessage *string
		err = rows.Scan(&cs.ID, &cs.OtherUserID, &cs.OtherUserName, &cs.OtherUserTag, &lastMessage, &lastAt, &cs.UnreadCount)
		if err != nil {
			return nil, err
		}
		cs.LastMessageAt = lastAt
		cs.LastMessage = lastMessage
		chats = append(chats, cs)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return chats, nil
}

func (cr *ConversationRepository) GetOrCreate(c context.Context, userId string, targetId string) (string, error) {
	if userId == targetId {
		return "", fmt.Errorf("cannot chat with self")
	}

	minId, maxId := utils.GetOrderedIds(userId, targetId)
	cacheKey := fmt.Sprintf("private_chat:%s:%s", minId, maxId)

	var convId string
	var err error

	convId, err = cr.rdb.Get(c, cacheKey).Result()
	if err == nil {
		return convId, nil
	}
	if err != redis.Nil {
		return "", err
	}

	tx, err := cr.pool.Begin(c)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(c)

	err = tx.QueryRow(c, `
        SELECT c.id FROM conversations c
        JOIN conversation_participants cp1 ON c.id = cp1.conversation_id
        JOIN conversation_participants cp2 ON c.id = cp2.conversation_id AND cp2.user_id != cp1.user_id
        WHERE c.type = 'private'
          AND ((cp1.user_id = $1 AND cp2.user_id = $2) OR (cp1.user_id = $2 AND cp2.user_id = $1))
          AND NOT EXISTS (SELECT 1 FROM conversation_participants cp3 WHERE cp3.conversation_id = c.id AND cp3.user_id NOT IN ($1, $2))
    `, userId, targetId).Scan(&convId)
	if err == pgx.ErrNoRows {
		convId = ulid.Make().String()
		_, err = tx.Exec(c, `
            INSERT INTO conversations (id, type) VALUES ($1, 'private')
        `, convId)
		if err != nil {
			return "", err
		}

		_, err = tx.Exec(c, `
            INSERT INTO conversation_participants (id, conversation_id, user_id) 
            VALUES ($1, $3, $4), ($2, $3, $5) ON CONFLICT DO NOTHING
        `, ulid.Make().String(), ulid.Make().String(), convId, userId, targetId)
		if err != nil {
			return "", err
		}

		go cr.asyncCacheSet(c, cacheKey, convId)
	} else if err != nil {
		return "", err
	}

	if err = tx.Commit(c); err != nil {
		return "", err
	}

	return convId, nil
}

func (cr *ConversationRepository) asyncCacheSet(c context.Context, key, convId string) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	cr.rdb.Set(ctx, key, convId, 24*time.Hour)
}
