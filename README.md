# response-std

Starter **framework API Golang** berbasis **Gin + GORM (MySQL)** dengan standar respons seragam, autentikasi token, RBAC (roles/permissions), routing versi, middleware umum, upload gambar, serta tooling migrasi & seeder.

> Fokus README ini: **jalan dulu dengan kelebihan yang sudah ada**. Tanpa roadmap atau rencana perbaikan.

---

## Fitur Utama
- **Standar respons konsisten** (`app/pkg/response` + `libs/responses`).
- **Auth**: login/register, token personal access ala `id|token` (hash SHA-256), `expires_at`.
- **RBAC**: roles & permissions mirip Spatie (middleware `role`, `permission`, `any-permission`).
- **Routing versi** via `libs/router` (`RouteRegistry`) dan `routes/api.go`.
- **Middleware**: CORS, recovery/panic handler, rate limiter global, Auth Bearer.
- **Upload gambar** via `goupload` (ext & size filter) + static serving `/storage`.
- **Migrations (SQL)** & **Seeder (Go)** siap jalan.
- **Logging** terstruktur + opsi kirim ke Discord Webhook.

## Arsitektur Singkat
```
app/
  http/controllers, middleware, requests
  models/entities
  pkg/response, pkg/permissions
config/ (ENV, DB, external API)
database/migrations/*.sql (up/down)
database/seeds/*.go
libs/external (api client, logger, hooks)
libs/router (RouteRegistry)
routes/ (api & web, versi v1)
main.go
```

---

## Requirements
- **Go** ≥ 1.22
- **MySQL** ≥ 8.0 (atau kompatibel, mis. MariaDB 10.6+)
- **golang-migrate** CLI (untuk mengeksekusi file SQL `up/down`)
- (Opsional) **Discord Webhook** untuk notifikasi log

### Instalasi golang-migrate (contoh)
- macOS (Homebrew): `brew install golang-migrate`
- Linux: lihat rilis resmi `migrate` sesuai OS/arch di [GitHub](https://github.com/golang-migrate/migrate/releases)
- Windows: disarakan ke [GitHub](https://github.com/golang-migrate/migrate/releases)
---

## Konfigurasi
Gunakan file **`.env`** di root proyek (dibaca oleh Viper). Contoh minimal:

```env
# App
APP_NAME=response-std
APP_PORT=5220
ENVIRONMENT=development
API_VERSION=v1,web
BASE_URL=http://localhost:5220
API_BASE_URL=http://localhost:5220/api/v1

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=golang_api

# Auth & Logging
JWT_SECRET=supersecretkey
LOG_LEVEL=info           # debug|info|warn|error
ENABLE_LOGGING=true
LOG_CHANNEL=file         # file|console
LOG_TO_FILE=true
LOG_DIR=logs

# Discord (opsional)
DISCORD_WEBHOOK_URL=
DISCORD_MIN_LOG_LEVEL=error

# External API (opsional)
EXTERNAL_API_BASE_URL=http://localhost:8080
EXTERNAL_API_ENDPOINT=/api/v1/
```

> Nilai default ada di `config/app.go`. DB diinisialisasi oleh `config.LoadDBMysql()`.

---

## Setup Database
Struktur tabel disediakan di `database/migrations/*.sql` (format timestamped `*.up.sql` / `*.down.sql`).

### 1) Jalankan migrasi dengan **wrapper** bawaan
Wrapper memanggil **golang-migrate**.

```bash
# UP semua migrasi
go run app/console/cmd/migrate/migrate.go up

# DOWN satu langkah
go run app/console/cmd/migrate/migrate.go down 1

# (opsional) FORCE ke versi tertentu
go run app/console/cmd/migrate/migrate.go force 20250614000102
```

> Pastikan `golang-migrate` sudah terpasang di PATH.

### 2) Seed data awal (roles & users)
```bash
go run app/console/cmd/scripts/seed/run.go
```

---

## Menjalankan Aplikasi
```bash
# 1) pastikan .env terisi
# 2) jalankan migrasi & seeder (opsional)
# 3) run server

go run .
# atau build
# go build -o bin/app && ./bin/app
```
Server akan berjalan di `:APP_PORT` (default `:5220`).

Health check sederhana:
- `GET /` (root)
- `GET /health`

---

## Endpoint Inti (v1)
Semua endpoint di-*mount* pada **`/api/v1`**. Cuplikan definisi rute (lihat `routes/api/v1.go`):

### Public
- `GET /` – root ping
- `GET /hello` – contoh teks
- `GET /health` – health check
- `POST /login` – login (Bearer token `id|raw-token`)
- `POST /register` – registrasi user

### Auth
- `POST /auth/logout` – logout (revoke token aktif)
- `POST /auth/refresh` – refresh token
- `GET /auth/me` – profil user saat ini

### Users (protected, contoh)
- `GET /users` – list users
- `POST /users` – create
- `GET /users/:id` – detail
- `PUT /users/:id` – update
- `DELETE /users/:id` – delete

> Akses ke endpoint tertentu dapat menggunakan middleware **role/permission** (`app/http/middleware`).

### Upload Gambar
- `POST /upload` – unggah gambar (allowed: `jpg,jpeg,png,webp`, max 2MB)
- Static files: `GET /storage/...` (otomatis dilayani)

---

## Format Respons Standar
Semua respons melalui wrapper `response.Success/Error` agar seragam. Contoh (sukses):
```json
{
  "success": true,
  "message": "OK",
  "data": { "example": 123 },
  "code": 200,
}
```
Contoh (error):
```json
{
  "success": false,
  "message": "Unauthorized",
  "errors": {"detail": "invalid token"},
  "code": 401,
}
```

---

## Otentikasi
- **Skema token**: `id|raw-token` (disimpan hash SHA-256 di tabel `personal_access_tokens`).
- **Header**: `Authorization: Bearer <id|raw-token>`
- **Kadaluarsa**: dikelola pada saat pembuatan token (lihat controller `auth_controller.go`).

---

## Middleware Utama
- **CORS**: diaktifkan via `gin-contrib/cors`.
- **Rate Limit (global)**: limiter proses `10r/s` burst `20`.
- **Recovery**: menangani panic → respons 500 standar.
- **AuthMiddleware**: validasi Bearer token.

---

## Logging
- Logger terintegrasi (`libs/external/services/logger.go`), level via `LOG_LEVEL`.
- Output ke file (default) di `storage/logs/` atau console.
- (Opsional) Kirim ke **Discord**: set `DISCORD_WEBHOOK_URL` & `DISCORD_MIN_LOG_LEVEL`.

---

## Contoh cURL
```bash
# register
curl -X POST "$API_BASE_URL/register" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret"}'

# login
curl -X POST "$API_BASE_URL/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"secret"}'
# → simpan token "id|raw-token"

# me
curl -H "Authorization: Bearer id|raw-token" "$API_BASE_URL/auth/me"

# upload file
curl -X POST "$BASE_URL/upload" \
  -H "Authorization: Bearer id|raw-token" \
  -F "file=@/path/to/image.jpg"
```

---

## Catatan
- Default CORS longgar untuk memudahkan dev; sesuaikan di `app/http/middleware` saat produksi.
- Pastikan direktori `storage/app/public/uploads/images` dapat ditulis oleh proses aplikasi.

---

## Lisensi
Terserah pemilik proyek.

