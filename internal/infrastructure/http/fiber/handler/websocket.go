package handler

import (
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
	"go.uber.org/zap"
)

type WebsocketHandler struct {
	logger           *zap.Logger
	websocketService *service.WebsocketService
}

func NewWebsocketHandler(websocketService *service.WebsocketService, logger *zap.Logger) *WebsocketHandler {
	return &WebsocketHandler{
		logger:           logger,
		websocketService: websocketService,
	}
}

func (wh *WebsocketHandler) Connect() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		userId, exists := c.Locals("userId").(string)
		if !exists {
			wh.logger.Warn("User ID required")
			c.Conn.Close()
			return
		}
		wh.logger.Info("connected", zap.String("userId", userId))
		defer wh.logger.Info("disconnected", zap.String("userId", userId))
		defer c.Conn.Close()

		wh.websocketService.Join(userId, c)
		defer wh.websocketService.Leave(userId)

		closeCh := make(chan struct{})
		go func() {
			pingTicker := time.NewTicker(30 * time.Second)
			defer pingTicker.Stop()
			for {
				select {
				case <-pingTicker.C:
					if c.Conn == nil {
						wh.logger.Warn("Ping attempted on nil Conn", zap.String("userId", userId))
						return
					}
					err := c.WriteMessage(websocket.PingMessage, nil)
					if err != nil {
						wh.logger.Warn("Failed to send ping", zap.Error(err), zap.String("userId", userId))
						return
					}
					wh.logger.Debug("ping sent", zap.String("userId", userId))
				case <-closeCh:
					wh.logger.Debug("Ping goroutine stopped", zap.String("userId", userId))
					return
				}
			}
		}()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				wh.logger.Warn("Error reading message", zap.Error(err), zap.String("userId", userId))
				close(closeCh)
				break
			}
			if mt == websocket.CloseMessage {
				wh.logger.Debug("Received close message", zap.String("userId", userId))
				close(closeCh)
				break
			}
			if mt == websocket.BinaryMessage {
				wh.logger.Warn("Binary message skipped", zap.String("userId", userId))
				continue
			}

			var msg map[string]interface{}
			if err = json.Unmarshal(message, &msg); err != nil {
				wh.logger.Warn("Failed to unmarshal JSON", zap.Error(err), zap.ByteString("raw", message), zap.String("userId", userId))
				continue
			}
			wh.logger.Debug("Received message", zap.String("userId", userId), zap.Any("msg", msg))

			if action, ok := msg["action"].(string); ok {
				switch action {
				default:
					wh.logger.Debug("Unhandled action", zap.String("action", action), zap.String("userId", userId))
				}
			}
		}
	}, websocket.Config{
		EnableCompression: false,
	})
}
