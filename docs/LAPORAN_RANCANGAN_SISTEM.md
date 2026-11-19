<!--
Catatan konversi PDF (pandoc + xelatex):
  ./docs/build-laporan.sh
Hasil: docs/LAPORAN_RANCANGAN_SISTEM.pdf (A4, halaman sampul, daftar isi).
Komentar HTML ini tidak akan tampil pada PDF.
-->

# RANCANGAN SISTEM HELPDESK TIKET LAYANAN AKADEMIK "LAPORKAMPUS" BERBASIS GO DAN FIBER DENGAN ARSITEKTUR BERLAPIS

**Laporan Rancangan Sistem**

| Keterangan | Isi |
|---|---|
| Nama | Vyerru |
| NIM | 434241021 |
| Program Studi | D4 Teknik Informatika |
| Fakultas | Vokasi |
| Universitas | Universitas Airlangga |
| Mata Kuliah | Praktikum Pemrograman Backend Lanjut |
| Semester | 5 |
| Tahun | 2026 |

---

## ABSTRAK

Layanan akademik di lingkungan kampus masih banyak dikelola melalui saluran yang tersebar, seperti pesan instan, surel, dan kedatangan langsung. Kondisi tersebut menyulitkan pelacakan status permintaan, pengukuran waktu penyelesaian, dan pembagian beban kerja petugas. Laporan ini merancang **LaporKampus**, yaitu sistem helpdesk tiket layanan akademik berbasis *Application Programming Interface* (API) yang menerapkan arsitektur berlapis dengan pola *repository*, *service layer*, dan *role-based access control* (RBAC). Sistem dibangun menggunakan bahasa Go dengan kerangka kerja Fiber, PostgreSQL sebagai basis data transaksional, MongoDB sebagai basis data dokumen untuk percakapan tiket dan survei kepuasan, serta Redis dan Asynq untuk pemrosesan tugas latar. Untuk menjamin konsistensi antara perubahan data dan pengiriman notifikasi, rancangan menerapkan pola *transactional outbox* dengan proses *relay* dan *worker* terpisah. Rancangan mencakup empat aktor (mahasiswa, agen, supervisor, dan administrator), mesin status tiket, kebijakan *Service Level Agreement* (SLA) beserta eskalasi otomatis, serta kuota beban kerja agen yang dijaga secara transaksional. Hasil rancangan diharapkan menjadi acuan implementasi yang terukur, dapat diuji, dan mudah dipelihara.

**Kata kunci:** helpdesk, arsitektur berlapis, Go, Fiber, RBAC, transactional outbox, SLA

---

## DAFTAR ISI

1. BAB I Pendahuluan
2. BAB II Landasan Teori
3. BAB III Analisis Kebutuhan
4. BAB IV Rancangan Sistem
5. BAB V Rencana Implementasi dan Pengujian
6. BAB VI Penutup
7. Daftar Pustaka
8. Lampiran

---

# BAB I PENDAHULUAN

## 1.1 Latar Belakang

Permintaan layanan akademik di perguruan tinggi mencakup legalisir dokumen, cuti akademik, perbaikan Kartu Rencana Studi (KRS), penggantian kata sandi, hingga penanganan gangguan sistem informasi. Pada praktiknya, permintaan tersebut sering disampaikan melalui pesan instan pribadi, surel yang berbeda-beda, atau kunjungan langsung ke unit layanan. Pola ini menimbulkan tiga masalah utama: permintaan mudah terlewat, status penanganan tidak dapat dilacak oleh pemohon, dan beban kerja petugas tidak terukur.

Sistem helpdesk merupakan pendekatan yang lazim digunakan untuk mengatasi persoalan tersebut. Permintaan diubah menjadi **tiket** yang memiliki identitas, status, penanggung jawab, dan batas waktu penyelesaian. Setiap perubahan tercatat sehingga dapat diaudit, sedangkan *Service Level Agreement* (SLA) memberikan dasar objektif untuk menilai mutu layanan.

Perancangan sistem helpdesk menuntut lebih dari sekadar operasi *Create, Read, Update, Delete* (CRUD). Terdapat empat kebutuhan yang saling terkait. Pertama, penugasan tiket harus menghormati kapasitas agen sehingga diperlukan transaksi dan penguncian baris. Kedua, notifikasi tidak boleh memperlambat atau menggagalkan pembuatan tiket sehingga diperlukan pemrosesan asinkron. Ketiga, perubahan data dan pengiriman peristiwa harus konsisten sehingga diperlukan pola *transactional outbox* [3]. Keempat, hak akses berbeda antarperan sehingga diperlukan RBAC yang diperiksa pada lapisan *middleware*.

Arsitektur berlapis dengan pemisahan *model*, *repository*, dan *service* dipilih karena memisahkan aturan bisnis dari detail teknis. Pemisahan ini memudahkan pengujian unit, penggantian teknologi penyimpanan, dan penelusuran kesalahan [1]. Bahasa Go dipilih karena mendukung konkurensi yang baik, menghasilkan biner tunggal, dan memiliki ekosistem yang memadai untuk kebutuhan API modern. Fiber digunakan sebagai kerangka kerja HTTP karena ringkas dan berkinerja tinggi di atas `fasthttp` [6].

## 1.2 Rumusan Masalah

Berdasarkan latar belakang tersebut, rumusan masalah laporan ini adalah sebagai berikut.

1. Bagaimana merancang sistem helpdesk layanan akademik yang mampu melacak status tiket, menegakkan SLA, dan membagi beban kerja agen secara adil?
2. Bagaimana menerapkan arsitektur berlapis, RBAC, dan *transactional outbox* pada sistem tersebut agar konsisten, aman, dan tetap tersedia ketika layanan pendukung gagal?
3. Bagaimana merencanakan pengujian dan implementasi sistem dalam batas waktu satu bulan?

## 1.3 Tujuan

Tujuan dari laporan rancangan ini adalah:

1. Menghasilkan rancangan fungsional dan non-fungsional sistem helpdesk tiket layanan akademik.
2. Menghasilkan rancangan arsitektur, basis data, antarmuka API, mesin status, dan alur asinkron yang siap diimplementasikan.
3. Menyusun rencana implementasi dan pengujian yang terukur beserta kriteria keberhasilannya.

## 1.4 Batasan Masalah

Agar pembahasan terfokus, rancangan dibatasi pada hal-hal berikut.

1. Sistem berupa **backend API**; antarmuka pengguna berbasis web atau seluler tidak termasuk dalam lingkup.
2. Lingkup layanan adalah layanan akademik dalam satu kampus, tanpa integrasi ke sistem akademik eksternal.
3. Pembayaran, kecerdasan buatan, obrolan waktu nyata, dan aplikasi seluler tidak dirancang.
4. SLA dihitung 24 jam sehari, tujuh hari seminggu; kalender jam kerja merupakan pengembangan lanjutan.
5. Sistem dijalankan sebagai satu biner dengan dua peran proses, yaitu `api` dan `worker`.

## 1.5 Manfaat

**Manfaat teoretis.** Laporan ini memperlihatkan penerapan arsitektur berlapis, pemisahan *port* dan adaptor secara pragmatis, RBAC, serta pola *transactional outbox* pada satu studi kasus utuh. Rancangan dapat dijadikan bahan pembelajaran mengenai keterkaitan antara aturan bisnis, transaksi basis data, dan pemrosesan asinkron.

**Manfaat praktis.** Bagi unit layanan akademik, rancangan menyediakan dasar sistem yang mampu melacak tiket, mengukur SLA, dan memantau beban agen. Bagi mahasiswa, rancangan menjadi acuan implementasi yang jelas sehingga pekerjaan tidak melebar dari lingkup.

## 1.6 Sistematika Penulisan

Bab I memaparkan latar belakang, rumusan masalah, tujuan, batasan, dan manfaat. Bab II membahas landasan teori yang digunakan. Bab III menyajikan analisis kebutuhan fungsional dan non-fungsional. Bab IV menguraikan rancangan arsitektur, data, API, mesin status, SLA, keamanan, dan penyebaran. Bab V memuat rencana implementasi serta strategi pengujian. Bab VI menutup laporan dengan kesimpulan dan saran.

---

# BAB II LANDASAN TEORI

## 2.1 Bahasa Go dan Kerangka Kerja Fiber

Go adalah bahasa pemrograman terkompilasi yang menekankan kesederhanaan, keamanan tipe, dan konkurensi melalui *goroutine* serta *channel*. Go menghasilkan biner tunggal tanpa ketergantungan *runtime* eksternal sehingga memudahkan penyebaran dalam kontainer [5].

Fiber adalah kerangka kerja web untuk Go yang dibangun di atas `fasthttp`. Fiber menyediakan abstraksi *routing*, *middleware*, dan penanganan konteks permintaan. Pada rancangan ini, Fiber berperan sebagai lapisan transportasi HTTP yang meneruskan permintaan ke lapisan *service* [6].

## 2.2 Arsitektur Berlapis, Service Layer, dan Repository

Arsitektur berlapis membagi sistem menjadi beberapa lapisan dengan tanggung jawab yang berbeda. Fowler menjelaskan bahwa *service layer* menetapkan batas aplikasi dan mengoordinasikan operasi bisnis, sedangkan *repository* memediasi antara domain dan lapisan pemetaan data [1].

Rancangan ini menggunakan empat lapisan:

1. **Transportasi** — `route` dan `middleware` menangani pemetaan URL, autentikasi, otorisasi, serta pencatatan permintaan.
2. **Aplikasi** — `app/service` memuat use case dan sekaligus menjadi penangan HTTP, mengikuti pola modul praktikum. Tidak terdapat paket `controller` atau `handler` terpisah.
3. **Domain dan akses data** — `app/model` memuat entitas serta objek transfer data, sedangkan `app/repository` memuat kontrak dan implementasi akses PostgreSQL maupun MongoDB.
4. **Infrastruktur** — `config`, `database`, `helper`, dan `worker` menyediakan konfigurasi, koneksi, utilitas, serta pemrosesan latar.

Aturan dependensi diterapkan searah: lapisan atas boleh bergantung pada lapisan bawah, tetapi tidak sebaliknya. `app/model` tidak mengimpor paket internal lain, dan `app/repository` tidak mengetahui keberadaan Fiber.

## 2.3 Autentikasi JWT dan Role-Based Access Control

*JSON Web Token* (JWT) adalah token bertanda tangan yang memuat klaim identitas. Sifatnya *stateless* sehingga server tidak perlu menyimpan sesi pada setiap permintaan. Rancangan menggunakan dua token: *access token* berdurasi pendek dan *refresh token* berdurasi panjang yang disimpan dalam basis data dalam bentuk *hash*.

RBAC mengatur hak akses berdasarkan peran. Pada rancangan ini, peran mahasiswa, agen, supervisor, dan administrator memiliki himpunan izin yang berbeda. Pemeriksaan izin dilakukan pada *middleware* sebelum permintaan mencapai *service*, sehingga aturan otorisasi tidak tersebar di banyak tempat [10].

## 2.4 PostgreSQL, MongoDB, dan Persistensi Poliglot

PostgreSQL adalah basis data relasional yang mendukung transaksi ACID, kunci baris, dan tipe data JSON. Karakteristik tersebut sesuai untuk data transaksional seperti tiket, penugasan, dan kuota agen [7].

MongoDB adalah basis data dokumen yang menyimpan data dalam bentuk BSON. Skema yang fleksibel sesuai untuk data yang tumbuh dan beragam, seperti utas percakapan tiket, lampiran, dan ulasan kepuasan [8].

Penggunaan dua jenis basis data dalam satu sistem disebut *polyglot persistence*. Rancangan ini menempatkan data yang membutuhkan konsistensi kuat di PostgreSQL dan data dokumen di MongoDB.

## 2.5 Antrean Tugas dan Asynq

Asynq adalah pustaka antrean tugas untuk Go yang menggunakan Redis sebagai *broker*. Produsen mengirim tugas ke antrean, sedangkan *worker* mengambil dan mengerjakannya. Asynq mendukung percobaan ulang, penundaan, serta pengarsipan tugas yang gagal [9].

Pemisahan API dan *worker* memberikan dua manfaat. Pertama, permintaan pengguna tidak tertahan oleh proses yang lambat seperti pengiriman surel. Kedua, pekerja dapat diperbanyak secara terpisah sesuai beban.

## 2.6 Transactional Outbox

Pola *transactional outbox* menyelesaikan masalah *dual write*, yaitu penulisan ke basis data dan ke *broker* pesan yang tidak berada dalam satu transaksi [3]. Pada pola ini, peristiwa disimpan ke tabel `outbox` dalam transaksi yang sama dengan perubahan data. Proses *relay* kemudian membaca tabel tersebut dan mengirimkannya ke antrean. Dengan demikian, kegagalan *broker* tidak menyebabkan kehilangan peristiwa; pengiriman hanya tertunda.

Konsekuensi pola ini adalah jaminan pengiriman *at-least-once*. Oleh karena itu, penangan tugas harus idempoten, misalnya dengan memeriksa kunci idempotensi sebelum menulis efek samping.

## 2.7 Pengujian Perangkat Lunak

Pengujian unit memeriksa fungsi murni secara terisolasi. Pengujian integrasi memeriksa interaksi antarkomponen dengan basis data nyata. *Race detector* digunakan untuk mendeteksi akses data bersamaan yang tidak aman. Pengujian beban mengukur perilaku sistem pada jumlah permintaan tinggi. Rancangan ini menerapkan keempat jenis pengujian tersebut secara bertahap.

## 2.8 Docker dan Penyebaran

Docker mengemas aplikasi beserta dependensinya ke dalam image. *Multi-stage build* memisahkan tahap kompilasi dari tahap akhir sehingga image lebih kecil dan permukaan serangan berkurang. Docker Compose digunakan untuk menjalankan PostgreSQL, MongoDB, Redis, API, dan *worker* dalam satu lingkungan pengembangan yang konsisten.

---

# BAB III ANALISIS KEBUTUHAN

## 3.1 Gambaran Umum Sistem

LaporKampus menerima permintaan layanan akademik dalam bentuk tiket. Mahasiswa membuat tiket dengan memilih kategori layanan dan prioritas. Sistem menghitung batas waktu respons dan penyelesaian berdasarkan kebijakan SLA kategori tersebut. Agen mengambil atau ditugaskan untuk menangani tiket, berkomunikasi melalui utas pesan, dan mengubah status sesuai alur yang sah. Supervisor memantau beban kerja, menerima eskalasi ketika SLA terlewati, dan membaca laporan. Administrator mengelola pengguna, peran, kategori, dan kebijakan SLA.

## 3.2 Aktor dan Hak Akses

**Tabel 3.1** Aktor sistem

| Aktor | Peran | Tanggung Jawab utama |
|---|---|---|
| Mahasiswa | `mahasiswa` | Membuat tiket, membalas utas, menutup tiket, mengisi CSAT |
| Agen | `agen` | Mengambil dan menangani tiket, membalas, menyelesaikan tiket |
| Supervisor | `supervisor` | Menugaskan tiket, memantau SLA dan beban, menerima eskalasi |
| Administrator | `admin` | Mengelola pengguna, peran, kategori, dan kebijakan SLA |

**Tabel 3.2** Matriks hak akses utama

| Aksi | mahasiswa | agen | supervisor | admin |
|---|:---:|:---:|:---:|:---:|
| Membuat tiket | ✓ | ✓ | ✓ | ✓ |
| Melihat tiket sendiri | ✓ | — | — | — |
| Melihat tiket pada layanan terkait | — | ✓ | ✓ | ✓ |
| Mengambil atau menugaskan tiket | — | ✓ | ✓ | ✓ |
| Mengubah status tiket | terbatas | ✓ | ✓ | ✓ |
| Membaca dan menulis utas | ✓ | ✓ | ✓ | ✓ |
| Mengisi CSAT | ✓ | — | — | — |
| Melihat dasbor dan laporan | — | — | ✓ | ✓ |
| Mengelola kategori dan SLA | — | — | — | ✓ |
| Mengelola pengguna dan peran | — | — | — | ✓ |

## 3.3 Kebutuhan Fungsional

**Tabel 3.3** Kebutuhan fungsional

| Kode | Kebutuhan | Aktor |
|---|---|---|
| FR-01 | Sistem menyediakan registrasi pengguna dengan validasi surel dan kekuatan kata sandi | Publik |
| FR-02 | Sistem menyediakan login, pembaruan *access token*, dan logout | Semua |
| FR-03 | Sistem memungkinkan administrator mengelola pengguna dan menetapkan peran | Administrator |
| FR-04 | Sistem memungkinkan administrator mengelola kategori layanan dan kebijakan SLA | Administrator |
| FR-05 | Sistem memungkinkan pembuatan tiket dengan kategori, subjek, deskripsi, dan prioritas | Mahasiswa |
| FR-06 | Sistem menyediakan daftar dan detail tiket dengan penyaringan serta paginasi kursor | Semua sesuai hak |
| FR-07 | Sistem memungkinkan penugasan tiket kepada agen dengan batas beban kerja | Agen, Supervisor |
| FR-08 | Sistem mengubah status tiket hanya melalui transisi yang sah | Agen, Supervisor |
| FR-09 | Sistem menyediakan utas pesan dan metadata lampiran pada setiap tiket | Semua sesuai hak |
| FR-10 | Sistem mengirim notifikasi perubahan status secara asinkron | Sistem |
| FR-11 | Sistem memantau batas SLA dan mengeskalasi tiket yang terlewati | Sistem |
| FR-12 | Sistem menerima penilaian kepuasan setelah tiket berstatus selesai | Mahasiswa |
| FR-13 | Sistem menyajikan ringkasan dasbor dan laporan kepatuhan SLA | Supervisor, Administrator |
| FR-14 | Sistem mencatat riwayat perubahan tiket pada tabel audit | Sistem |
| FR-15 | Sistem menyediakan *health check* dan *readiness check* | Operator |

## 3.4 Kebutuhan Non-Fungsional

**Tabel 3.4** Kebutuhan non-fungsional

| Kode | Kebutuhan | Tolok ukur |
|---|---|---|
| NFR-01 | Keamanan | Kata sandi di-*hash* bcrypt; JWT bertanda tangan; RBAC pada setiap endpoint terlindungi; pembatasan laju login |
| NFR-02 | Kinerja | Waktu respons baca pada persentil ke-95 di bawah 300 ms untuk data uji 10.000 tiket |
| NFR-03 | Ketahanan | Pembuatan tiket tetap berhasil ketika Redis atau *worker* tidak tersedia |
| NFR-04 | Konsistensi | Tidak terjadi pelampauan kuota agen pada pengujian konkurensi |
| NFR-05 | Skalabilitas | Proses `api` dan `worker` dapat diperbanyak secara terpisah |
| NFR-06 | Keteramatan | Log terstruktur, *request identifier*, dan pemeriksaan kesiapan dependensi |
| NFR-07 | Portabilitas | Seluruh layanan dapat dijalankan melalui Docker Compose |
| NFR-08 | Pemeliharaan | Pemisahan lapisan, gaya kode konsisten, serta *lint* dan pengujian otomatis |
| NFR-09 | Auditabilitas | Setiap perubahan status tiket tercatat beserta pelaku dan waktu |
| NFR-10 | Dokumentasi | Kontrak API terdokumentasi dalam OpenAPI |

## 3.5 Use Case Utama

Use case utama sistem adalah: (UC-01) registrasi dan login; (UC-02) pembuatan tiket; (UC-03) penugasan tiket dengan pemeriksaan kuota; (UC-04) penanganan tiket melalui utas pesan; (UC-05) perubahan status tiket; (UC-06) pemantauan SLA dan eskalasi; (UC-07) pengisian CSAT; serta (UC-08) pembacaan dasbor dan laporan.

## 3.6 Asumsi dan Ketergantungan

1. Tersedia layanan PostgreSQL, MongoDB, dan Redis yang dapat dihubungi dari proses aplikasi.
2. Jam sistem seluruh layanan disinkronkan dan perhitungan SLA menggunakan waktu UTC.
3. Pengiriman surel pada tahap awal disimulasikan melalui *logger* agar tidak bergantung pada penyedia surel eksternal.
4. Data pengguna awal dan kategori layanan diisi melalui migrasi dan *seed*.

---

# BAB IV RANCANGAN SISTEM

## 4.1 Arsitektur Sistem

Rancangan menerapkan arsitektur berlapis dengan satu biner yang dapat berjalan sebagai proses `api`, `worker`, atau keduanya. Lapisan transportasi tidak mengetahui detail basis data, sedangkan lapisan akses data tidak mengetahui keberadaan Fiber.

```
Klien HTTP
    │
    ▼
┌──────────────────────────────────────────────────────────┐
│ route + middleware                                       │
│ request id · recover · CORS · helmet · auth · authz      │
└───────────────────────────┬──────────────────────────────┘
                            ▼
┌──────────────────────────────────────────────────────────┐
│ app/service  (use case + handler HTTP)                   │
│ aturan status · kuota agen · SLA · validasi              │
└───────────────┬───────────────────────────┬──────────────┘
                ▼                           ▼
┌───────────────────────────┐   ┌──────────────────────────┐
│ app/repository            │   │ tabel outbox             │
│ PostgreSQL (pgx)          │   │ ditulis dalam transaksi   │
│ MongoDB (driver v2)       │   └────────────┬─────────────┘
└───────────────────────────┘                │
                                             ▼
                                  relay ──► Redis (Asynq) ──► worker
                                                              │
                                                              ▼
                                                     notifikasi / eskalasi
```

*Gambar 4.1* Arsitektur umum LaporKampus

Aturan dependensi yang berlaku:

```
main ─► config ─► route ─► middleware ─► helper
                    │
                    └─► app/service ─► app/repository ─► app/model
```

## 4.2 Struktur Direktori

```
lapor-kampus/
├── app/
│   ├── model/          entitas, DTO, tipe bersama
│   ├── repository/     kontrak dan implementasi akses data, TxManager
│   └── service/        use case dan handler Fiber
├── config/             env, logger, perakitan Fiber, opsi antrean
├── database/           pool PostgreSQL, client MongoDB
├── docs/               arsitektur, ADR, laporan
├── helper/             respons, error, konteks, JWT, keamanan, validasi, kursor
├── middleware/         request id, recover, CORS, auth, authz
├── migrations/         berkas SQL golang-migrate
├── route/              pemetaan URL
├── tests/integration/  pengujian integrasi
├── tools/              loadtest, redeliver
├── main.go             composition root dan graceful shutdown
└── worker.go           server Asynq
```

## 4.3 Rancangan Basis Data

### 4.3.1 Skema PostgreSQL

**Tabel 4.1** Tabel autentikasi dan otorisasi

| Tabel | Kolom utama | Keterangan |
|---|---|---|
| `users` | `id`, `username`, `email`, `password`, `role`, `is_active`, `created_at` | Akun pengguna; kata sandi disimpan sebagai *hash* |
| `roles` | `name`, `description` | Daftar peran |
| `permissions` | `name`, `description` | Daftar izin |
| `role_permissions` | `role`, `permission` | Pemetaan peran ke izin |
| `refresh_tokens` | `id`, `user_id`, `token_hash`, `expires_at`, `revoked_at` | Token penyegar yang dapat dicabut |

**Tabel 4.2** Tabel domain tiket

| Tabel | Kolom utama | Keterangan |
|---|---|---|
| `categories` | `id`, `name`, `description`, `sla_policy_id`, `is_active` | Kategori layanan |
| `sla_policies` | `id`, `priority`, `first_response_minutes`, `resolution_minutes` | Kebijakan SLA per prioritas |
| `tickets` | `id`, `code`, `category_id`, `requester_id`, `assignee_id`, `priority`, `status`, `subject`, `sla_due_at`, `first_response_at`, `resolved_at` | Entitas tiket |
| `ticket_events` | `id`, `ticket_id`, `actor_id`, `type`, `payload`, `created_at` | Jejak audit |
| `outbox` | `id`, `event_type`, `payload`, `published_at`, `created_at` | Peristiwa yang menunggu dikirim |
| `notifications` | `id`, `user_id`, `title`, `body`, `read_at`, `created_at` | Notifikasi dalam aplikasi |

Indeks yang direncanakan mencakup `tickets(status, assignee_id)` untuk perhitungan beban, `tickets(requester_id, created_at)` untuk paginasi kursor, `outbox(published_at)` untuk proses *relay*, dan `refresh_tokens(user_id)` untuk pencabutan token.

### 4.3.2 Koleksi MongoDB

| Koleksi | Struktur | Keterangan |
|---|---|---|
| `ticket_messages` | `ticket_id`, `author_id`, `body`, `attachments[]`, `created_at` | Utas percakapan; jumlah dan bentuk lampiran dapat berubah |
| `csat_reviews` | `ticket_id`, `score`, `comment`, `created_at` | Penilaian kepuasan setelah tiket selesai |

Indeks unik pada `csat_reviews.ticket_id` mencegah penilaian ganda, sedangkan indeks pada `ticket_messages.ticket_id` dan `created_at` mempercepat pembacaan utas.

### 4.3.3 Relasi Antar-Entitas

```
users 1 ──── N tickets (sebagai requester)
users 1 ──── N tickets (sebagai assignee)
categories 1 ──── N tickets
sla_policies 1 ──── N categories
tickets 1 ──── N ticket_events
users 1 ──── N refresh_tokens
users 1 ──── N notifications

MongoDB:
tickets 1 ──── N ticket_messages
tickets 1 ──── 1 csat_reviews
```

*Gambar 4.2* Relasi antar-entitas

## 4.4 Mesin Status Tiket

Status tiket dibatasi pada himpunan berikut: `OPEN`, `ASSIGNED`, `IN_PROGRESS`, `WAITING_USER`, `RESOLVED`, dan `CLOSED`. Transisi yang sah ditetapkan sebagai fungsi murni sehingga dapat diuji tanpa basis data.

```
  OPEN
   | assign
   v
  ASSIGNED
   | mulai
   v
  IN_PROGRESS --tunggu pengguna--> WAITING_USER
   |                                    |
   |<-------- pengguna membalas --------+
   | selesai
   v
  RESOLVED --tutup--> CLOSED
   |
   +-- buka kembali --> OPEN

  Pembatalan oleh supervisor: OPEN/ASSIGNED --> CLOSED
```

*Gambar 4.3* Mesin status tiket

**Tabel 4.3** Matriks transisi status

| Dari | Ke | Pemicu | Pelaku |
|---|---|---|---|
| `OPEN` | `ASSIGNED` | Penugasan agen | Agen, Supervisor |
| `ASSIGNED` | `IN_PROGRESS` | Agen mulai menangani | Agen |
| `IN_PROGRESS` | `WAITING_USER` | Agen menunggu tanggapan | Agen |
| `WAITING_USER` | `IN_PROGRESS` | Pengguna membalas | Mahasiswa |
| `IN_PROGRESS` | `RESOLVED` | Agen menyatakan selesai | Agen |
| `RESOLVED` | `CLOSED` | Konfirmasi penutupan | Mahasiswa, Supervisor |
| `RESOLVED` | `OPEN` | Tiket dibuka kembali | Mahasiswa |
| `OPEN`/`ASSIGNED` | `CLOSED` | Pembatalan | Supervisor |

Transisi di luar tabel tersebut ditolak dengan kode kesalahan `INVALID_TRANSITION`.

## 4.5 Rancangan SLA dan Eskalasi

Setiap kategori merujuk pada satu kebijakan SLA. Kebijakan menetapkan batas waktu respons pertama dan batas waktu penyelesaian untuk setiap tingkat prioritas. Batas waktu dihitung saat tiket dibuat.

```
batas_respons   = waktu_dibuat + first_response_minutes(prioritas)
batas_selesai   = waktu_dibuat + resolution_minutes(prioritas)
terlampaui      = sekarang > batas_selesai
                  DAN status ∉ {RESOLVED, CLOSED}
```

Proses `worker` memeriksa tiket yang belum selesai secara berkala. Tiket yang melewati batas ditandai, dinaikkan prioritasnya bila perlu, dan peristiwanya dikirim kepada supervisor melalui outbox. Pemeriksaan dijalankan dengan kueri berindeks dan dibatasi jumlah barisnya agar tidak membebani basis data.

## 4.6 Rancangan Antarmuka API

Seluruh endpoint berada di bawah prefiks `/api/v1`. Respons sukses menggunakan amplop `success`, `message`, `data`, dan `meta`; respons gagal menggunakan `success`, `code`, `message`, `fields`, dan `request_id`.

**Tabel 4.4** Rancangan endpoint

| Metode | Path | Keterangan | Hak akses |
|:------|:--------------------------------|:--------------------------|:-------------|
| GET | `/health` | Pemeriksaan hidup proses | Publik |
| GET | `/readyz` | Pemeriksaan kesiapan dependensi | Publik |
| POST | `/auth/register` | Registrasi pengguna | Publik |
| POST | `/auth/login` | Login dan penerbitan token | Publik |
| POST | `/auth/refresh` | Pembaruan access token | Publik |
| POST | `/auth/logout` | Pencabutan refresh token | Terautentikasi |
| GET | `/auth/me` | Profil pengguna aktif | Terautentikasi |
| GET | `/categories` | Daftar kategori | Terautentikasi |
| POST | `/categories` | Membuat kategori | `category:manage` |
| PUT | `/categories/:id` | Mengganti kategori | `category:manage` |
| PATCH | `/categories/:id` | Mengubah sebagian kategori | `category:manage` |
| DELETE | `/categories/:id` | Menonaktifkan kategori | `category:manage` |
| GET | `/tickets` | Daftar tiket dengan paginasi kursor | Terautentikasi |
| POST | `/tickets` | Membuat tiket | `ticket:create` |
| GET | `/tickets/:id` | Detail tiket | Pemilik atau agen |
| POST | `/tickets/:id/assign` | Menugaskan tiket | `ticket:assign` |
| PATCH | `/tickets/:id/status` | Mengubah status | `ticket:update` |
| GET | `/tickets/:id/messages` | Membaca utas | Pemilik atau agen |
| POST | `/tickets/:id/messages` | Menulis utas | Pemilik atau agen |
| POST | `/tickets/:id/csat` | Mengisi penilaian | Pemilik tiket |
| GET | `/dashboard/summary` | Ringkasan tiket dan SLA | `dashboard:read` |
| GET | `/dashboard/agents` | Beban kerja agen | `dashboard:read` |

## 4.7 Rancangan Autentikasi dan Otorisasi

Alur autentikasi dimulai ketika pengguna mengirim kredensial ke `/auth/login`. Sistem memverifikasi kata sandi dengan bcrypt, menerbitkan *access token* berdurasi 15 menit dan *refresh token* berdurasi 7 hari. *Refresh token* disimpan sebagai *hash* SHA-256 sehingga kebocoran basis data tidak langsung memberikan token yang dapat dipakai.

Setiap permintaan terlindungi melewati dua tahap. Pertama, `RequireAuth` mengambil token dari header `Authorization: Bearer`, memverifikasi tanda tangan dan masa berlaku, lalu menempelkan identitas ke konteks permintaan. Kedua, `RequirePermission` memeriksa apakah peran pengguna memiliki izin yang dibutuhkan. Kegagalan pada tahap pertama menghasilkan `401`, sedangkan kegagalan pada tahap kedua menghasilkan `403`.

**Tabel 4.5** Daftar izin yang direncanakan

| Izin | mahasiswa | agen | supervisor | admin |
|---|:---:|:---:|:---:|:---:|
| `ticket:create` | ✓ | ✓ | ✓ | ✓ |
| `ticket:read:own` | ✓ | — | — | — |
| `ticket:read:all` | — | ✓ | ✓ | ✓ |
| `ticket:assign` | — | ✓ | ✓ | ✓ |
| `ticket:update` | — | ✓ | ✓ | ✓ |
| `ticket:close` | ✓ | ✓ | ✓ | ✓ |
| `ticket:manage` | — | — | ✓ | ✓ |
| `category:manage` | — | — | — | ✓ |
| `user:manage` | — | — | — | ✓ |
| `dashboard:read` | — | — | ✓ | ✓ |

## 4.8 Rancangan Pemrosesan Asinkron

Pembuatan tiket dan penugasan menulis dua hal dalam satu transaksi: perubahan pada tabel `tickets` dan satu baris pada tabel `outbox`. API tidak berkomunikasi langsung dengan Redis. Proses `relay` membaca baris outbox yang belum terkirim, mengirimkannya ke antrean, lalu menandainya terkirim. *Worker* mengambil tugas dan menjalankan efek samping, seperti menyimpan notifikasi atau mengirim surel.

```
API            PostgreSQL          Relay            Redis           Worker
 │  BEGIN            │                │                │               │
 │  INSERT ticket    │                │                │               │
 │  INSERT outbox    │                │                │               │
 │  COMMIT ─────────►│                │                │               │
 │  201 Created      │                │                │               │
 │                   │◄── poll ───────│                │               │
 │                   │─── event ─────►│─── enqueue ───►│               │
 │                   │◄────────────── tandai terkirim │               │
 │                   │                │                │─── task ─────►│
 │                   │◄────────────── tulis notifikasi ◄───────────────│
```

*Gambar 4.4* Alur transactional outbox

Konsekuensi rancangan ini adalah pengiriman bersifat *at-least-once*. Oleh karena itu, penangan tugas memeriksa kunci idempotensi sebelum menulis notifikasi. Jika Redis tidak tersedia, API tetap dapat membuat tiket; peristiwa tertahan di outbox dan dikirim ketika Redis pulih.

## 4.9 Rancangan Penanganan Kesalahan dan Respons

Kesalahan bisnis direpresentasikan sebagai `AppError` yang memuat status HTTP, kode, pesan, dan detail validasi. Satu penangan kesalahan terpusat di `config` menerjemahkan kesalahan tersebut menjadi respons JSON. Kesalahan tak terduga dipetakan menjadi `500` tanpa membocorkan detail internal, tetapi tetap dicatat pada log beserta `request_id`.

Kode kesalahan yang direncanakan mencakup `VALIDATION_ERROR`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `INVALID_TRANSITION`, `AGENT_CAPACITY_EXCEEDED`, `TOO_MANY_REQUESTS`, dan `INTERNAL_ERROR`.

## 4.10 Rancangan Keamanan

1. Kata sandi di-*hash* dengan bcrypt pada biaya 12 dan tidak pernah dikembalikan melalui API.
2. *Access token* berumur pendek; *refresh token* dapat dicabut dan disimpan dalam bentuk *hash*.
3. Login dibatasi lima percobaan per menit per alamat IP.
4. Header keamanan dipasang melalui Helmet; CORS dibatasi pada origin yang dikonfigurasi.
5. Ukuran badan permintaan dibatasi 1 MB.
6. Validasi masukan dilakukan sebelum menyentuh basis data.
7. Rahasia disimpan dalam variabel lingkungan dan tidak dimasukkan ke repositori.
8. Swagger dinonaktifkan secara bawaan pada lingkungan produksi.

## 4.11 Rancangan Penyebaran

Lingkungan pengembangan dijalankan dengan Docker Compose yang memuat PostgreSQL 16, MongoDB 7, dan Redis 7. API dan *worker* dibangun dari satu Dockerfile *multi-stage* dan dijalankan sebagai dua kontainer terpisah pada profil `app`. Pemeriksaan kesehatan kontainer memakai endpoint `/health`, sedangkan kesiapan aplikasi memakai `/readyz`.

---

# BAB V RENCANA IMPLEMENTASI DAN PENGUJIAN

## 5.1 Metode dan Lingkungan Implementasi

Implementasi menggunakan metode iteratif-incremental dengan siklus mingguan. Setiap iterasi menghasilkan irisan vertikal yang dapat didemontrasikan, mulai dari fondasi, domain tiket, pemrosesan asinkron, hingga pengerasan dan dokumentasi.

**Tabel 5.1** Lingkungan implementasi

| Komponen | Keterangan |
|---|---|
| Bahasa | Go 1.25 atau lebih baru |
| Kerangka kerja | Fiber v2 |
| Basis data relasional | PostgreSQL 16 |
| Basis data dokumen | MongoDB 7 |
| Antrean | Redis 7 dan Asynq |
| Autentikasi | JWT dan bcrypt |
| Validasi | go-playground/validator |
| Migrasi | golang-migrate |
| Kontainerisasi | Docker dan Docker Compose |
| Integrasi berkelanjutan | GitHub Actions |

## 5.2 Jadwal Implementasi

**Tabel 5.2** Jadwal implementasi empat minggu

| Minggu | Fokus | Keluaran |
|---|---|---|
| 1 | Fondasi, autentikasi, RBAC | API menyala, health/readyz, login empat peran, migrasi awal |
| 2 | Domain tiket | Kategori, SLA, tiket, mesin status, penugasan dengan kuota |
| 3 | MongoDB dan asinkron | Utas pesan, CSAT, outbox, relay, worker, notifikasi |
| 4 | SLA, laporan, pengerasan | Eskalasi otomatis, dasbor, Swagger, pengujian integrasi, Docker, dokumentasi |

## 5.3 Strategi Pengujian

**Tabel 5.3** Jenis dan cakupan pengujian

| Jenis | Sasaran | Perkakas | Kriteria |
|---|---|---|---|
| Unit | Fungsi aturan: transisi status, SLA, kuota, validasi | `go test`, tabel kasus uji | Cakupan aturan ≥ 90% |
| Integrasi | Repository, endpoint, transaksi | `go test -tags=integration`, basis data uji | Skenario utama lulus |
| Konkurensi | Penugasan tiket bersamaan | `go test -race`, goroutine | Tidak ada pelampauan kuota |
| Beban | Daftar tiket dan pembuatan tiket | `hey` atau `k6` | p95 < 300 ms pada 10.000 tiket |
| Statis | Kualitas kode dan kerentanan | `go vet`, `golangci-lint`, `govulncheck` | Tanpa temuan kritis |

## 5.4 Kriteria Keberhasilan

1. Seluruh kebutuhan fungsional prioritas tinggi terimplementasi dan dapat didemontrasikan.
2. Pengujian unit aturan bisnis lulus dengan cakupan minimal 90%.
3. Pengujian konkurensi membuktikan tidak terjadi pelampauan kuota agen.
4. Notifikasi tetap terkirim setelah Redis pulih, tanpa kehilangan peristiwa.
5. API tetap melayani pembuatan tiket ketika *worker* dimatikan.
6. `docker compose up` menjalankan seluruh sistem dan pemeriksaan kesiapan berhasil.
7. *Pipeline* integrasi berkelanjutan berstatus hijau.

## 5.5 Manajemen Risiko

**Tabel 5.4** Risiko dan mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Lingkup melebar | Pekerjaan tidak selesai | Menetapkan daftar fitur di luar lingkup sejak awal |
| Perhitungan SLA keliru akibat zona waktu | Eskalasi salah waktu | Menyimpan dan menghitung waktu dalam UTC |
| Tugas ganda akibat percobaan ulang | Notifikasi berulang | Kunci idempotensi pada penangan tugas |
| Kuota agen terlampaui | Beban tidak adil | Transaksi dan `SELECT ... FOR UPDATE` |
| Rahasia masuk repositori | Kebocoran kredensial | `.gitignore`, `.env.example`, pemindaian rahasia |
| Perbedaan lingkungan | Aplikasi gagal di mesin lain | Docker Compose sejak awal |

---

# BAB VI PENUTUP

## 6.1 Kesimpulan

Laporan ini merancang LaporKampus, sistem helpdesk tiket layanan akademik berbasis Go dan Fiber dengan arsitektur berlapis. Rancangan mencakup empat aktor dengan hak akses berbeda, mesin status tiket yang eksplisit, kebijakan SLA beserta eskalasi otomatis, penugasan berbasis kuota yang dijaga secara transaksional, serta pemrosesan notifikasi asinkron melalui pola *transactional outbox*. Pemisahan lapisan *model*, *repository*, dan *service* menjaga aturan bisnis tetap teruji dan tidak bergantung pada kerangka kerja HTTP. Penggunaan PostgreSQL untuk data transaksional dan MongoDB untuk data dokumen disesuaikan dengan karakteristik masing-masing data.

## 6.2 Saran

Pengembangan lanjutan dapat diarahkan pada: (1) penghitungan SLA berbasis kalender jam kerja dan hari libur; (2) antarmuka pengguna web yang memakai API ini; (3) penambahan basis pengetahuan dan deteksi tiket serupa; (4) pemantauan metrik dan penelusuran terdistribusi; serta (5) mekanisme *dead-letter* dan pengiriman ulang otomatis untuk tugas yang gagal permanen.

---

# DAFTAR PUSTAKA

[1] Fowler, M. (2002). *Patterns of Enterprise Application Architecture*. Boston: Addison-Wesley.

[2] Evans, E. (2003). *Domain-Driven Design: Tackling Complexity in the Heart of Software*. Boston: Addison-Wesley.

[3] Richardson, C. (2019). *Microservices Patterns: With Examples in Java*. Shelter Island: Manning.

[4] Kleppmann, M. (2017). *Designing Data-Intensive Applications*. Sebastopol: O'Reilly Media.

[5] The Go Authors. (2026). *The Go Programming Language Documentation*. https://go.dev/doc/

[6] Fiber Contributors. (2026). *Fiber Documentation*. https://docs.gofiber.io/

[7] PostgreSQL Global Development Group. (2026). *PostgreSQL 16 Documentation*. https://www.postgresql.org/docs/16/

[8] MongoDB, Inc. (2026). *MongoDB Manual*. https://www.mongodb.com/docs/manual/

[9] Hibiken. (2026). *Asynq: Simple, Reliable, and Efficient Distributed Task Queue in Go*. https://github.com/hibiken/asynq

[10] OWASP Foundation. (2021). *OWASP Top 10:2021*. https://owasp.org/Top10/

[11] Newman, S. (2021). *Building Microservices: Designing Fine-Grained Systems* (2nd ed.). Sebastopol: O'Reilly Media.

[12] Nygard, M. T. (2018). *Release It! Design and Deploy Production-Ready Software* (2nd ed.). Raleigh: Pragmatic Bookshelf.

---

# LAMPIRAN

## Lampiran A — Contoh Payload

**A.1** Permintaan pembuatan tiket

```json
{
  "category_id": 3,
  "subject": "Perbaikan KRS semester gasal",
  "description": "Mata kuliah yang sudah diambil tidak muncul pada KRS.",
  "priority": "HIGH"
}
```

**A.2** Respons pembuatan tiket

```json
{
  "success": true,
  "message": "tiket berhasil dibuat",
  "data": {
    "id": 101,
    "code": "TK-2026-000101",
    "status": "OPEN",
    "priority": "HIGH",
    "sla_due_at": "2026-11-19T09:00:00Z"
  }
}
```

**A.3** Respons kesalahan validasi

```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "message": "validasi gagal",
  "fields": {
    "subject": "minimal 5 karakter"
  },
  "request_id": "8f1c2a9e-4b7d-4f0e-9a51-2c6d7e8f9a0b"
}
```

## Lampiran B — Variabel Lingkungan

**Tabel B.1** Variabel lingkungan utama

| Variabel | Contoh | Keterangan |
|---|---|---|
| `APP_NAME` | `LaporKampus` | Nama aplikasi |
| `APP_PORT` | `3000` | Port HTTP |
| `APP_ENV` | `development` | Lingkungan |
| `APP_ROLE` | `all` | Peran proses: `api`, `worker`, `all` |
| `DB_HOST` | `localhost` | Hos PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL |
| `DB_USER` | `lapor` | Pengguna basis data |
| `DB_PASSWORD` | *(rahasia)* | Kata sandi basis data |
| `DB_NAME` | `lapor_kampus` | Nama basis data |
| `MONGO_URI` | `mongodb://localhost:27017` | URI MongoDB |
| `MONGO_DATABASE` | `lapor_kampus` | Nama basis data MongoDB |
| `REDIS_ADDR` | `localhost:6379` | Alamat Redis |
| `JWT_SECRET` | *(minimal 32 karakter)* | Kunci tanda tangan token |
| `JWT_ISSUER` | `lapor-kampus` | Penerbit token |
| `JWT_ACCESS_TTL_MINUTES` | `15` | Masa berlaku access token |
| `JWT_REFRESH_TTL_DAYS` | `7` | Masa berlaku refresh token |
| `SWAGGER_ENABLED` | `false` | Mengaktifkan Swagger UI |

## Lampiran C — Ringkasan Keputusan Arsitektur

**Tabel C.1** Keputusan utama dan alasannya

| Keputusan | Alternatif | Alasan pemilihan |
|---|---|---|
| Arsitektur berlapis dengan *service layer* | Monolit terstruktur per fitur | Sesuai pola modul, mudah dipahami dan diuji |
| Satu biner dengan dua peran proses | Dua biner terpisah | Menyederhanakan pembangunan dan penyebaran awal |
| PostgreSQL untuk data transaksional | Seluruhnya MongoDB | Membutuhkan transaksi dan kunci baris |
| MongoDB untuk utas dan CSAT | Seluruhnya PostgreSQL | Data dokumen yang tumbuh dan beragam |
| *Transactional outbox* | Penulisan ganda langsung | Menghindari kehilangan peristiwa saat Redis gagal |
| RBAC berbasis izin | Pemeriksaan peran langsung | Izin dapat berubah tanpa mengubah kode endpoint |
