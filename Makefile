# ── Variabel ────────────────────────────────────────────────────────────────
# .env dimuat lebih dulu agar DB_DSN di bawah memakai nilai yang benar.
-include .env
export

APP        := lapor-kampus
BIN_DIR    := bin
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w \
	-X github.com/vyerru/lapor-kampus/config.Version=$(VERSION) \
	-X github.com/vyerru/lapor-kampus/config.Commit=$(COMMIT) \
	-X github.com/vyerru/lapor-kampus/config.BuildDate=$(BUILD_DATE)

# Migrasi memakai CLI lokal bila tersedia; jika tidak, image resmi
# golang-migrate dijalankan lewat Docker.
MIGRATE_IMAGE := migrate/migrate:v4.17.0

# Nilai bawaan harus sama dengan docker-compose agar `make migrate-*` tetap
# berjalan sebelum .env dibuat.
DB_HOST    ?= localhost
DB_PORT    ?= 5433
DB_USER    ?= lapor
DB_PASSWORD ?= lapor
DB_NAME    ?= lapor_kampus
DB_SSLMODE ?= disable
DB_DSN     ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

ifeq ($(shell command -v migrate 2>/dev/null),)
MIGRATE_RUN  := docker run --rm --network host -v "$(PWD)/migrations:/migrations" $(MIGRATE_IMAGE)
MIGRATE_PATH := -path=/migrations
MIGRATE_DIR  := /migrations
else
MIGRATE_RUN  := migrate
MIGRATE_PATH := -path=$(PWD)/migrations
MIGRATE_DIR  := $(PWD)/migrations
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Tampilkan daftar target
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ── Infrastruktur ───────────────────────────────────────────────────────────
.PHONY: up down logs ps
up: ## Nyalakan Postgres, Mongo, Redis
	docker compose up -d

down: ## Matikan seluruh container (data tetap tersimpan)
	docker compose down

logs: ## Ikuti log container
	docker compose logs -f --tail=100

ps: ## Status container
	docker compose ps

# ── Aplikasi ────────────────────────────────────────────────────────────────
.PHONY: run run-api run-worker build
run: ## Jalankan api + worker (role=all)
	go run .

run-api: ## Jalankan hanya HTTP API
	APP_ROLE=api go run .

run-worker: ## Jalankan hanya worker
	APP_ROLE=worker go run .

build: ## Build binary ke ./bin
	@mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP) .

# ── Kualitas ────────────────────────────────────────────────────────────────
.PHONY: test test-race cover lint vet fmt vuln tidy
test: ## Unit test
	go test ./... -count=1

test-race: ## Unit test dengan race detector
	go test ./... -count=1 -race -shuffle=on

cover: ## Coverage -> coverage.html
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint: ## golangci-lint (wajib terpasang)
	golangci-lint run ./...

vet: ## go vet
	go vet ./...

fmt: ## Rapikan format & import
	gofmt -l -w .
	@command -v goimports >/dev/null && goimports -w . || true

vuln: ## Pemindaian kerentanan dependensi
	@command -v govulncheck >/dev/null && govulncheck ./... || \
		go run golang.org/x/vuln/cmd/govulncheck@latest ./...

tidy: ## Rapikan go.mod/go.sum
	go mod tidy

# ── Migrasi ─────────────────────────────────────────────────────────────────
.PHONY: migrate-up migrate-down migrate-version migrate-create
migrate-up: ## Terapkan semua migrasi
	$(MIGRATE_RUN) "$(MIGRATE_PATH)" -database "$(DB_DSN)" up

migrate-down: ## Batalkan satu migrasi terakhir
	$(MIGRATE_RUN) "$(MIGRATE_PATH)" -database "$(DB_DSN)" down 1

migrate-version: ## Lihat versi migrasi aktif
	$(MIGRATE_RUN) "$(MIGRATE_PATH)" -database "$(DB_DSN)" version

migrate-create: ## Buat pasangan migrasi: make migrate-create NAME=nama
	@test -n "$(NAME)" || (echo "NAME wajib diisi"; exit 1)
	$(MIGRATE_RUN) create -ext sql -dir "$(MIGRATE_DIR)" -seq $(NAME)

# ── Docker ──────────────────────────────────────────────────────────────────
.PHONY: docker-build app-up app-down
docker-build: ## Build image aplikasi
	docker build -t $(APP):$(VERSION) .

app-up: ## Nyalakan infra + api + worker dalam Docker
	docker compose --profile app up -d --build

app-down: ## Matikan termasuk api + worker
	docker compose --profile app down

# ── Kebersihan ──────────────────────────────────────────────────────────────
.PHONY: clean
clean: ## Hapus artefak build & coverage
	rm -rf $(BIN_DIR) coverage.out coverage.html
