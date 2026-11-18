package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/vyerru/lapor-kampus/config"
	"github.com/vyerru/lapor-kampus/helper"
)

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	app := fiber.New(fiber.Config{
		AppName:   config.GetEnv("APP_NAME", "LaporKampus"),
		BodyLimit: 1 * 1024 * 1024,
	})

	app.Get("/health", healthCheck)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := config.GetEnv("APP_PORT", "3000")
	go func() {
		logger.Info("HTTP berjalan", slog.String("port", port))
		if err := app.Listen(":" + port); err != nil && !errors.Is(err, fiber.ErrServiceUnavailable) {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("sinyal berhenti diterima")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Error("shutdown tidak rapi", slog.String("error", err.Error()))
	}
	logger.Info("aplikasi berhenti dengan rapi")
}

func healthCheck(c *fiber.Ctx) error {
	return helper.Success(c, fiber.StatusOK, "aplikasi sehat", fiber.Map{
		"service": config.GetEnv("APP_NAME", "LaporKampus"),
		"env":     config.GetEnv("APP_ENV", "development"),
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}
