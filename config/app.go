package config

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// NewApp merakit Fiber dengan konfigurasi server saja. Middleware, route, dan
// handler HTTP adalah tanggung jawab pemilik module masing-masing (helper,
// middleware, route) yang belum digabung ke repository ini.
func NewApp(logger *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      GetEnv("APP_NAME", "LaporKampus"),
		ErrorHandler: newErrorHandler(logger),
		BodyLimit:    1 * 1024 * 1024,
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"code":    "NOT_FOUND",
			"message": "endpoint tidak ditemukan",
		})
	})

	return app
}

// newErrorHandler adalah satu-satunya tempat error diterjemahkan menjadi HTTP.
func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		requestID := c.Get("X-Request-ID")

		status := fiber.StatusInternalServerError
		code := "INTERNAL"
		message := "kesalahan internal"

		var fiberErr *fiber.Error
		switch {
		case errors.Is(err, fiber.ErrRequestEntityTooLarge):
			status, code, message = fiber.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE",
				"ukuran body melebihi batas yang diizinkan"
		case errors.As(err, &fiberErr):
			status, code, message = fiberErr.Code, "HTTP_ERROR", fiberErr.Message
		}

		if status >= fiber.StatusInternalServerError {
			logger.Error("request_failed",
				slog.String("request_id", requestID),
				slog.String("path", c.Path()),
				slog.String("code", code),
				slog.Int("status", status),
				slog.String("error", err.Error()))
		} else {
			logger.Warn("request_rejected",
				slog.String("request_id", requestID),
				slog.String("path", c.Path()),
				slog.String("code", code),
				slog.Int("status", status))
		}

		return c.Status(status).JSON(fiber.Map{
			"success":    false,
			"code":       code,
			"message":    message,
			"request_id": requestID,
		})
	}
}
