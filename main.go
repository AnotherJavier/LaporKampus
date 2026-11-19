package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/vyerru/lapor-kampus/config"
	"github.com/vyerru/lapor-kampus/database"
)

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	// Satu program, dua peran. Di production, api dan worker dijalankan
	// sebagai proses terpisah agar masing-masing dapat diperbanyak.
	role := flag.String("role", config.GetEnv("APP_ROLE", "all"), "api | worker | all")
	flag.Parse()
	runAPI := *role == "api" || *role == "all"
	runWorker := *role == "worker" || *role == "all"
	if !runAPI && !runWorker {
		logger.Error("role tidak dikenal", slog.String("role", *role))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var app *fiber.App
	if runAPI {
		app = startAPI(logger)
		if app == nil {
			os.Exit(1)
		}
	}

	var stopWorker func()
	if runWorker {
		stopWorker = startWorker(logger)
		if stopWorker == nil {
			os.Exit(1)
		}
	}

	logger.Info("aplikasi berjalan",
		slog.String("role", *role),
		slog.String("version", config.Version),
		slog.String("commit", config.Commit))
	<-ctx.Done()

	// Urutan berhenti: pintu masuk (HTTP) ditutup lebih dulu, baru worker.
	if app != nil {
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			logger.Error("shutdown HTTP tidak rapi", slog.String("error", err.Error()))
		}
	}
	if stopWorker != nil {
		stopWorker()
	}
	logger.Info("aplikasi berhenti dengan rapi")
}

// startAPI menyiapkan connection pool (PostgreSQL + MongoDB) lalu menyalakan
// server Fiber. Route & middleware adalah tanggung jawab modul lain yang belum
// digabung ke repository ini.
func startAPI(logger *slog.Logger) *fiber.App {
	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke PostgreSQL", slog.String("error", err.Error()))
		return nil
	}
	defer pool.Close()

	mongoClient, err := database.NewMongoClient(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke MongoDB", slog.String("error", err.Error()))
		return nil
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mongoClient.Disconnect(ctx)
	}()

	app := config.NewApp(logger)

	port := config.GetEnv("APP_PORT", "3000")
	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("HTTP berjalan", slog.String("port", port))

	// TODO: sambungkan pool & mongoClient ke handler setelah route/middleware tersedia
	_ = pool
	_ = mongoClient

	return app
}
