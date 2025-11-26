package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neokofg/callap-backend/internal/domain/entity"
	"github.com/oklog/ulid/v2"
)

type FriendRepository struct {
	pool          *pgxpool.Pool
	tableName     string
	userTableName string
}

func NewFriendRepository(pool *pgxpool.Pool) *FriendRepository {
	return &FriendRepository{
		pool:          pool,
		tableName:     friendsTableName,
		userTableName: userTableName,
	}
}

func (fr *FriendRepository) List(c context.Context, userId string, limit int, offset int) ([]entity.FriendUser, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(
		"SELECT u.id, u.name, u.tag FROM %s f JOIN %s u ON f.friend_id = u.id WHERE f.user_id = $1 AND f.status = 'accepted' ORDER BY f.created_at DESC LIMIT $2 OFFSET $3",
		fr.tableName,
		fr.userTableName,
	)
	rows, err := fr.pool.Query(c, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []entity.FriendUser
	for rows.Next() {
		var fu entity.FriendUser
		var idStr string
		err := rows.Scan(&idStr, &fu.Name, &fu.Tag)
		if err != nil {
			return nil, err
		}
		fu.Id = ulid.MustParse(idStr)
		friends = append(friends, fu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}

func (fr *FriendRepository) Delete(c context.Context, userId string, friendId string) error {
	tx, err := fr.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	var status string
	query := fmt.Sprintf(
		"SELECT status FROM %s WHERE user_id = $1 AND friend_id = $2",
		fr.tableName,
	)
	err = tx.QueryRow(c, query, userId, friendId).Scan(&status)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("no connection to delete: user=%s, friend=%s", userId, friendId)
	}
	if err != nil {
		return err
	}

	if status == "accepted" {
		query = fmt.Sprintf(
			"DELETE FROM %s WHERE user_id = $1 AND friend_id = $2",
			fr.tableName,
		)
		_, err = tx.Exec(c, query, userId, friendId)
		if err != nil {
			return err
		}
		query = fmt.Sprintf(
			"DELETE FROM %s WHERE user_id = $2 AND friend_id = $1",
			fr.tableName,
		)
		_, err = tx.Exec(c, query, friendId, userId)
		if err != nil {
			return err
		}
	} else if status == "pending" {
		query = fmt.Sprintf(
			"DELETE FROM %s WHERE user_id = $1 AND friend_id = $2",
			fr.tableName,
		)
		_, err = tx.Exec(c, query, userId, friendId)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("cannot delete: invalid status %s", status)
	}

	return tx.Commit(c)
}

func (fr *FriendRepository) Decline(c context.Context, userId string, id string) error {
	tx, err := fr.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	query := fmt.Sprintf(
		"UPDATE %s SET status = 'rejected', updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE id = $1 AND friend_id = $2",
		fr.tableName,
	)
	result, err := tx.Exec(c, query, id, userId)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("friend request not found or no permission: id=%s, user=%s", id, userId)
	}

	return tx.Commit(c)
}

func (fr *FriendRepository) Accept(c context.Context, userId string, id string) error {
	tx, err := fr.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	query := fmt.Sprintf(
		"UPDATE %s SET status = 'accepted', updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND friend_id = $2",
		fr.tableName,
	)
	result, err := tx.Exec(c, query, id, userId)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("friend request not found or no permission: id=%s, user=%s", id, userId)
	}
	var senderID string
	querySelect := fmt.Sprintf(
		"SELECT user_id FROM %s WHERE id = $1",
		fr.tableName,
	)
	err = tx.QueryRow(c, querySelect, id).Scan(&senderID)
	if err != nil {
		return fmt.Errorf("failed to get sender: %w", err)
	}
	queryInsert := fmt.Sprintf(
		"INSERT INTO %s (id, user_id, friend_id, status, created_at, updated_at) VALUES ($1, $2, $3, 'accepted', CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC') ON CONFLICT (user_id, friend_id) DO UPDATE SET status = 'accepted', updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE %s.user_id = $2 AND %s.friend_id = $3",
		fr.tableName, fr.tableName, fr.tableName,
	)
	_, err = tx.Exec(c, queryInsert, ulid.Make().String(), userId, senderID)
	if err != nil {
		return fmt.Errorf("failed to insert reverse friend: %w", err)
	}
	return tx.Commit(c)
}

func (fr *FriendRepository) GetPending(c context.Context, userId string) ([]*entity.PendingFriend, error) {
	query := fmt.Sprintf(
		"SELECT f.id, u.id, u.name, u.tag FROM %s f JOIN %s u ON f.user_id = u.id WHERE f.friend_id = $1 AND f.status = 'pending' ORDER BY f.created_at DESC",
		fr.tableName,
		fr.userTableName,
	)
	rows, err := fr.pool.Query(c, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pending []*entity.PendingFriend
	for rows.Next() {
		var pf entity.PendingFriend
		err = rows.Scan(&pf.ID, &pf.SenderID, &pf.Name, &pf.Tag)
		if err != nil {
			return nil, err
		}
		pending = append(pending, &pf)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pending, nil
}

func (fr *FriendRepository) AddFriend(c context.Context, userId string, friendId string) error {
	tx, err := fr.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	parsedUserId := ulid.MustParse(userId)
	parsedFriendId := ulid.MustParse(friendId)

	var existingStatus string
	query := fmt.Sprintf(
		"SELECT status FROM %s WHERE user_id = $1 AND friend_id = $2",
		fr.tableName,
	)
	err = tx.QueryRow(c, query, parsedUserId.String(), parsedFriendId.String()).Scan(&existingStatus)
	if err == nil {
		if existingStatus == "pending" || existingStatus == "accepted" {
			return fmt.Errorf("request already exists: status %s", existingStatus)
		} else if existingStatus == "rejected" {
			query = fmt.Sprintf(""+
				"UPDATE %s SET status = 'pending', updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC' WHERE user_id = $1 AND friend_id = $2",
				fr.tableName,
			)
			_, err = tx.Exec(c, query, parsedUserId, parsedFriendId)
			if err != nil {
				return err
			}
			return tx.Commit(c)
		}
	} else if err != pgx.ErrNoRows {
		return err
	}

	newFriend := entity.NewFriend(entity.Friend{
		UserId:   parsedUserId,
		FriendId: parsedFriendId,
		Status:   "pending",
	})

	query = fmt.Sprintf(
		"INSERT INTO %s (id, user_id, friend_id, status) VALUES ($1, $2, $3, $4)",
		fr.tableName,
	)
	_, err = tx.Exec(
		c, query,
		newFriend.Id.String(),
		newFriend.UserId.String(),
		newFriend.FriendId.String(),
		newFriend.Status,
	)
	if err != nil {
		return err
	}
	return tx.Commit(c)
}
