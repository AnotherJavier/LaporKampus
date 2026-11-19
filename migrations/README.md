# Migrasi Database

Skema dikelola dengan [golang-migrate](https://github.com/golang-migrate/migrate).
Setiap perubahan skema adalah satu pasang berkas:

```
migrations/NNNNNN_nama.up.sql
migrations/NNNNNN_nama.down.sql
```

## Menjalankan

```bash
make migrate-up          # terapkan semua migrasi
make migrate-down        # batalkan satu migrasi terakhir
make migrate-version     # lihat versi aktif
make migrate-create NAME=tambah_tabel_csat
```

Makefile memakai CLI `migrate` bila tersedia. Jika tidak, image Docker
`migrate/migrate` dijalankan otomatis. Alamat database diambil dari variabel
`DB_*` di `.env`.

## Aturan

1. Jangan mengubah migrasi yang sudah pernah dijalankan — buat migrasi baru.
2. `down.sql` harus benar-benar mengembalikan keadaan sebelumnya.
3. Satu migrasi = satu tujuan yang jelas (jangan campur domain berbeda).
4. Nama berkas memakai `snake_case`.
