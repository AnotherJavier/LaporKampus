package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vyerru/lapor-kampus/config"
)

// NewPool membuat connection pool ke PostgreSQL dan memastikan koneksi
// benar-benar hidup sebelum aplikasi dinyatakan siap.
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := config.GetEnv("DATABASE_URL", "")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.GetEnv("DB_USER", "lapor"),
			config.GetEnv("DB_PASSWORD", "lapor"),
			config.GetEnv("DB_HOST", "localhost"),
			config.GetEnv("DB_PORT", "5432"),
			config.GetEnv("DB_NAME", "lapor_kampus"),
			config.GetEnv("DB_SSLMODE", "disable"),
		)
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("konfigurasi database tidak valid: %w", err)
	}

	maxConns := config.GetEnvInt("DB_MAX_CONNS", 10)
	if maxConns < 1 {
		maxConns = 1
	}
	if maxConns > 1000 {
		maxConns = 1000
	}
	// Nilai sudah dibatasi 1..1000 sehingga aman dikonversi ke int32.
	cfg.MaxConns = int32(maxConns) //nolint:gosec // dibatasi sebelum konversi
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	return pool, nil
}
