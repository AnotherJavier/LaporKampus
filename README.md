# LaporKampus

Sistem helpdesk tiket layanan akademik kampus: mahasiswa membuka tiket,
agen menangani sesuai SLA, supervisor memantau beban kerja dan eskalasi.

Dibangun dengan Go, Fiber, PostgreSQL, MongoDB, dan Redis (Asynq), mengikuti
arsitektur berlapis yang dipakai pada modul praktikum.

## Fitur yang direncanakan

- Autentikasi JWT (access + refresh) dengan RBAC: mahasiswa, agen, supervisor, admin.
- Tiket: buat, tugaskan, ubah status, komentar (thread di MongoDB).
- SLA per kategori dengan penghitung jatuh tempo dan eskalasi otomatis.
- Kuota/beban kerja agen dan penugasan berbasis kapasitas.
- Notifikasi asinkron lewat outbox + worker.
- Survei kepuasan (CSAT) setelah tiket selesai.
- Health, readiness, migrasi, CI, dan Docker.

## Status

Boilerplate selesai: kerangka HTTP, data layer, auth scaffolding, worker,
migrasi awal, Docker, CI, dan dokumentasi. Fitur bisnis menyusul per hari.

## Prasyarat

- Go 1.26+
- Docker + Docker Compose
- (opsional) `golangci-lint`, `govulncheck`, `golang-migrate`

> Port host sengaja digeser agar tidak bentrok dengan instalasi lokal:
> Postgres `5433`, MongoDB `27018`, Redis `6380`. Di dalam jaringan Docker
> aplikasi tetap memakai port bawaannya (`5432`, `27017`, `6379`).

## Mulai cepat

```bash
git clone https://github.com/vyerru/lapor-kampus.git
cd lapor-kampus

cp .env.example .env
# isi JWT_SECRET minimal 32 karakter:
#   openssl rand -hex 32

make up            # Postgres, Mongo, Redis
make migrate-up    # skema awal
make run           # api + worker
```

Uji:

```bash
curl -s localhost:3000/api/v1/health | jq
curl -s localhost:3000/api/v1/readyz | jq
```

Menjalankan terpisah:

```bash
make run-api       # APP_ROLE=api
make run-worker    # APP_ROLE=worker
```

Seluruh stack di dalam Docker:

```bash
make app-up        # profil "app": infra + api + worker
make app-down
```

## Endpoint

| Method | Path | Keterangan |
|---|---|---|
| GET | `/api/v1/health` | liveness: proses hidup |
| GET | `/api/v1/readyz` | readiness: Postgres + Mongo dapat dihubungi |

Respons selalu berbentuk:

```json
{ "success": true, "message": "...", "data": {} }
```

Error:

```json
{
  "success": false,
  "code": "NOT_FOUND",
  "message": "endpoint tidak ditemukan",
  "request_id": "..."
}
```

## Perintah

```bash
make help          # daftar target
make test          # unit test
make test-race     # unit test + race detector
make cover         # coverage.html
make lint          # golangci-lint
make vuln          # govulncheck
make build         # bin/lapor-kampus (version, commit, date via -ldflags)
make migrate-up    # terapkan migrasi
make migrate-down  # batalkan satu migrasi
make migrate-create NAME=nama_migrasi
```

## Struktur

```
app/
  model/         entitas + DTO + tipe bersama
  repository/    kontrak & implementasi akses data, TxManager
  service/       use case + handler Fiber (tanpa controller terpisah)
config/          env, logger, perakitan Fiber, queue
database/        pool PostgreSQL, client MongoDB
docs/            arsitektur + ADR
helper/          response, error, context, JWT, keamanan, validasi, cursor
middleware/      request id, recover, CORS, auth, authz
migrations/      SQL golang-migrate
route/           pemetaan URL
tests/           uji integrasi (build tag `integration`)
tools/           loadtest, redeliver
main.go          composition root + graceful shutdown
worker.go        server Asynq
```

Detail: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Health & graceful shutdown

`/api/v1/health` tidak menyentuh dependency dan dipakai sebagai liveness probe.
`/api/v1/readyz` memeriksa Postgres dan Mongo dan dipakai sebagai readiness probe.
Saat `SIGINT`/`SIGTERM`, HTTP ditutup lebih dulu, kemudian worker, dengan batas
10 detik.

## Lisensi

MIT — lihat [LICENSE](LICENSE).
