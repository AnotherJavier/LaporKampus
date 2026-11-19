package main

import (
	"log/slog"

	"github.com/hibiken/asynq"

	"github.com/vyerru/lapor-kampus/config"
)

// startWorker menyalakan server Asynq. Saat ini belum ada task yang
// terdaftar; kerangka ini memastikan koneksi Redis, retry, dan shutdown
// sudah benar sebelum fitur notifikasi masuk.
func startWorker(logger *slog.Logger) func() {
	client := asynq.NewClient(config.RedisOpt())

	// ServeMux memetakan JENIS task ke handler-nya, seperti route memetakan
	// URL ke handler pada Fiber.
	mux := asynq.NewServeMux()
	// TODO(outbox): mux.HandleFunc(model.TaskSendNotification, handler.HandleSend)

	srv := config.NewWorkerServer(logger)
	if err := srv.Start(mux); err != nil {
		logger.Error("gagal menyalakan worker", slog.String("error", err.Error()))
		_ = client.Close()
		return nil
	}
	logger.Info("worker berjalan")

	return func() {
		srv.Shutdown()
		_ = client.Close()
	}
}
