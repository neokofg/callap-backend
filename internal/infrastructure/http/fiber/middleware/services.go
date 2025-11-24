package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
)

func ServicesMiddleware(services *service.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("services", services)
		return c.Next()
	}
}
