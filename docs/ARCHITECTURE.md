# Arsitektur LaporKampus

## Gaya

Berlapis (Layered) + Repository + Service Layer, mengikuti pola modul praktikum:

```
HTTP request
   │
   ▼
route      → memetakan URL ke method service
   │
   ▼
service    → aturan bisnis + sekaligus handler Fiber
   │
   ▼
repository → akses PostgreSQL / MongoDB
   │
   ▼
database   → pool pgx, client Mongo
```

Tidak ada paket `controller`, `handler`, atau `dto`. Service-lah yang menerima
`*fiber.Ctx`, sama seperti contoh pada modul.

## Aturan dependensi

```
main → config → route → middleware → helper
                  │
                  └→ app/service → app/repository → app/model
```

- `app/model` tidak mengimpor paket internal lain.
- `app/repository` hanya mengimpor `model` dan driver database.
- `app/service` mengimpor `model`, `repository`, `helper`.
- `route` mengimpor `service`, `middleware`, `helper`.
- `main` adalah composition root: satu-satunya tempat merakit semua bagian.

## Peran biner

Satu biner, tiga peran (`-role` atau `APP_ROLE`):

| Peran | Isi |
|---|---|
| `api` | server Fiber |
| `worker` | server Asynq (notifikasi, eskalasi SLA) |
| `all` | keduanya dalam satu proses (khusus pengembangan) |

Di production, `api` dan `worker` dijalankan sebagai container terpisah agar
dapat diskalakan sendiri-sendiri.

## Alur error

Service mengembalikan `*helper.AppError`. `config.newErrorHandler` adalah
satu-satunya tempat error diterjemahkan menjadi respons HTTP, sehingga bentuk
JSON konsisten dan `request_id` selalu ikut.

## Keputusan yang dapat ditinjau

Lihat `docs/adr/` untuk catatan keputusan arsitektur.
