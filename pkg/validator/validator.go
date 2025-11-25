package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type XValidator struct {
	validator *validator.Validate
}

func Validate(c *fiber.Ctx, logger *zap.Logger, data interface{}) error {
	validationErrors := []ErrorResponse{}
	errs := validator.New().Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse
			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true
			validationErrors = append(validationErrors, elem)
		}
	}

	if len(validationErrors) > 0 {
		errMsgs := make([]string, 0, len(validationErrors))
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf("[%s]: %s", err.FailedField, err.Tag))
		}
		logger.Warn("Validation failed", zap.Strings("errors", errMsgs))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  strings.Join(errMsgs, ", "),
		})
	}

	return nil
}
