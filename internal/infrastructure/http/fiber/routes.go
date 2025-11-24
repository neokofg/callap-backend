package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/handler"
)

type Routes struct {
	handlers *handler.Handlers
}

func InitRoutes(fiberApp *fiber.App, handlers *handler.Handlers) *fiber.App {
	routes := &Routes{
		handlers: handlers,
	}

	v1 := fiberApp.Group("/api").Group("/v1")
	routes.authRoutes(v1)
	routes.userRoutes(v1)

	return fiberApp
}

func (r *Routes) authRoutes(fiberRouter fiber.Router) {
	groupAuth := fiberRouter.Group("/auth")
	groupAuth.Post("/register", r.handlers.AuthHandler.Register)
	groupAuth.Post("/login")
}

func (r *Routes) userRoutes(fiberRouter fiber.Router) {
	groupUser := fiberRouter.Group("/user")
	groupUser.Get("/me")
}
