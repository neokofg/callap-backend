package service

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
)

type Message struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
	Data   any    `json:"data"`
}

type Hub struct {
	conns  map[string]*websocket.Conn
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		conns:  make(map[string]*websocket.Conn),
		ctx:    ctx,
		cancel: cancel,
	}
}

type WebsocketService struct {
	cTimeout time.Duration
	hub      *Hub
	logger   *zap.Logger
}

func NewWebsocketService(cTimeout time.Duration, logger *zap.Logger) *WebsocketService {
	return &WebsocketService{
		cTimeout: cTimeout,
		hub:      NewHub(),
		logger:   logger,
	}
}

func (ws *WebsocketService) Join(userId string, conn *websocket.Conn) {
	ws.hub.mu.Lock()
	if ws.hub.conns[userId] != nil {
		delete(ws.hub.conns, userId)
	}
	ws.hub.conns[userId] = conn
	ws.hub.mu.Unlock()
}

func (ws *WebsocketService) Leave(userId string) {
	ws.hub.mu.Lock()
	delete(ws.hub.conns, userId)
	ws.hub.mu.Unlock()
}

func (ws *WebsocketService) SendToUser(userId string, msg Message) {
	ws.hub.mu.RLock()
	conn, ok := ws.hub.conns[userId]
	ws.hub.mu.RUnlock()
	if !ok {
		ws.logger.Warn("No connection for user", zap.String("userId", userId))
		return
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		ws.logger.Warn("Error writing to websocket", zap.Error(err), zap.String("userId", userId))
		ws.hub.mu.Lock()
		delete(ws.hub.conns, userId)
		ws.hub.mu.Unlock()
	} else {
		ws.logger.Debug("Message sent to user", zap.String("userId", userId))
	}
}
