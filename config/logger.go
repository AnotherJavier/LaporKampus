package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger membuat logger JSON yang menulis ke stdout SEKALIGUS ke berkas
// yang dirotasi. Satu format log untuk seluruh aplikasi.
func NewLogger() *slog.Logger {
	if err := os.MkdirAll("logs", 0o750); err != nil {
		panic("gagal membuat folder logs: " + err.Error())
	}
	rotator := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "app.log"),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   true,
	}
	writer := io.MultiWriter(os.Stdout, rotator)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: parseLevel(GetEnv("LOG_LEVEL", "info")),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func parseLevel(value string) slog.Level {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
