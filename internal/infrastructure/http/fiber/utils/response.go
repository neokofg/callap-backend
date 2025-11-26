package utils

import "github.com/gofiber/fiber/v2"

func MakeSuccessResponseWithData(data any) map[string]interface{} {
	return fiber.Map{
		"success": true,
		"data":    data,
	}
}

func MakeSuccessResponse() map[string]interface{} {
	return fiber.Map{
		"success": true,
	}
}
