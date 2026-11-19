package config

// Nilai di bawah diisi saat build melalui -ldflags, misalnya:
//
//	go build -ldflags "-X github.com/vyerru/lapor-kampus/config.Version=1.0.0"
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)
