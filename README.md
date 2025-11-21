# Users Service

Layanan microservice sederhana berbasis Gin untuk mengelola data pengguna dan terhubung ke PostgreSQL.

## Setup
- Salin `.env.example` ke `.env` dan sesuaikan nilai koneksi database.
- Pastikan database sudah dibuat dan migrasi `db/migrations/0001_create_users.sql` telah dijalankan.
- Jalankan server dengan `go run cmd/server/main.go`.

## Endpoint Ringkas
- `GET /health` – pengecekan layanan.
- `POST /api/v1/users` – buat pengguna (`{ "name": "...", "email": "..." }`).
- `GET /api/v1/users` – daftar pengguna.
- `GET /api/v1/users/:id` – detail pengguna.
