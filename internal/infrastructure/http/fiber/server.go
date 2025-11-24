package fiber

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/handler"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/middleware"
	"go.uber.org/zap"
)

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func InitFiber(cfg *config.Config, logger *zap.Logger, services *service.Services) {
	fiberApp := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: false,
		StrictRouting: false,
		AppName:       cfg.Host,
		Concurrency:   256 * 1024,
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	fiberApp.Use(helmet.New())

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	fiberApp.Use(middleware.ConfigMiddleware(cfg))

	handlers := handler.NewHandlers(services, logger)

	fiberApp = InitRoutes(fiberApp, handlers)

	defer func(fiberApp *fiber.App) {
		err := fiberApp.Shutdown()
		if err != nil {
			logger.Fatal("failed to shutdown http server", zap.Error(err))
		}
	}(fiberApp)

	go func() {
		logger.Info("Starting server")
		if err := fiberApp.Listen(":8000"); err != nil {
			logger.Error("Error starting server", zap.Error(err))
		}
	}()
}
