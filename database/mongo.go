package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/vyerru/lapor-kampus/config"
)

// NewMongoClient menyiapkan client MongoDB dan memastikan server benar-benar
// dapat dihubungi. Pada driver v2, mongo.Connect TIDAK menghubungi server;
// Ping wajib agar aplikasi menolak menyala bila Mongo tidak siap.
func NewMongoClient(ctx context.Context) (*mongo.Client, error) {
	maxPool := config.GetEnvInt("MONGO_MAX_POOL", 20)
	if maxPool < 1 {
		// Nilai negatif tidak bermakna untuk ukuran pool dan akan menjadi
		// angka sangat besar bila dikonversi ke tipe tanpa tanda.
		maxPool = 20
	}

	opts := options.Client().
		ApplyURI(config.GetEnv("MONGO_URI", "mongodb://localhost:27017")).
		SetServerSelectionTimeout(5 * time.Second).
		// maxPool dijamin positif di atas.
		SetMaxPoolSize(uint64(maxPool)) //nolint:gosec // dibatasi sebelum konversi

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("menyiapkan client MongoDB: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	return client, nil
}
