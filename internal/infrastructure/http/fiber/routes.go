package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/handler"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/middleware"
)

type Routes struct {
	handlers *handler.Handlers
}

func InitRoutes(fiberApp *fiber.App, handlers *handler.Handlers, services *service.Services) *fiber.App {
	routes := &Routes{
		handlers: handlers,
	}

	v1 := fiberApp.Group("/api").Group("/v1")
	routes.authRoutes(v1)
	routes.userRoutes(v1, services)

	return fiberApp
}

func (r *Routes) authRoutes(fiberRouter fiber.Router) {
	groupAuth := fiberRouter.Group("/auth")
	groupAuth.Post("/register", r.handlers.AuthHandler.Register)
	groupAuth.Post("/login", r.handlers.AuthHandler.Login)
}

func (r *Routes) userRoutes(fiberRouter fiber.Router, services *service.Services) {
	groupUser := fiberRouter.Group("/user", middleware.AuthMiddleware(services.JWT))
	groupUser.Get("/me", r.handlers.UserHandler.Me)
	r.friendRoutes(groupUser, services)
}

func (r *Routes) friendRoutes(fiberRouter fiber.Router, services *service.Services) {
	groupFriend := fiberRouter.Group("/friend")
	groupFriend.Post("/add", r.handlers.FriendHandler.AddFriend)
	groupFriend.Get("/pending", r.handlers.FriendHandler.GetPending)
	groupFriend.Post("/accept", r.handlers.FriendHandler.Accept)
	groupFriend.Post("/decline", r.handlers.FriendHandler.Decline)
	groupFriend.Delete("/delete", r.handlers.FriendHandler.Delete)
	groupFriend.Get("/list", r.handlers.FriendHandler.List)
}
