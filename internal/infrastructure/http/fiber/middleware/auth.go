package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/pkg/jwt"
)

func AuthMiddleware(jwtService *jwt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var accessToken string

		authHeader := c.Request().Header.Peek("Authorization")
		if authHeader != nil && strings.HasPrefix(string(authHeader), "Bearer ") {
			accessToken = strings.TrimPrefix(string(authHeader), "Bearer ")
		}
		if accessToken != "" {
			body, err := jwtService.ValidateToken(accessToken)
			if err != nil {
				return err
			}
			c.Locals("userId", body.UserId.String())
		}

		return c.Next()
	}
}
