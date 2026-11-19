# ADR 0001 — Arsitektur berlapis mengikuti modul praktikum

- Status: diterima
- Tanggal: 2026-11-18

## Konteks

Proyek akhir semester harus dapat dipertanggungjawabkan sebagai pekerjaan
mandiri, tetapi juga harus konsisten dengan pola yang diajarkan pada modul
(MODUL_2–MODUL_12). Dosen pembimbing menilai kesesuaian pola, bukan hanya
fitur yang berjalan.

## Keputusan

1. Memakai struktur `app/model`, `app/repository`, `app/service`, `config`,
   `database`, `helper`, `middleware`, `route`, `migrations`.
2. Tidak membuat paket `controller`, `handler`, atau `dto`. Service menerima
   `*fiber.Ctx` dan mengembalikan `error`.
3. Satu biner dengan flag peran `api | worker | all`.
4. Aturan bisnis yang berdiri sendiri diletakkan di `<fitur>_rules.go` dan
   diuji unit, terpisah dari kode akses data.

## Alasan

- Konsisten dengan contoh yang menjadi acuan penilaian.
- Lebih sedikit lapisan berarti lebih sedikit kode perpindahan data tanpa
  kehilangan pemisahan tanggung jawab yang penting.
- Aturan bisnis dapat diuji tanpa database maupun HTTP.

## Konsekuensi

- Service menjadi lebih besar; diimbangi dengan memecah per fitur
  (`auth_service.go`, `ticket_service.go`, `ticket_rules.go`).
- Perubahan pada kontrak HTTP mudah merembet ke service. Diterima karena
  respons HTTP justru bagian dari use case pada aplikasi sebesar ini.
