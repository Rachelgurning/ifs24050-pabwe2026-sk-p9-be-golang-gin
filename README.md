# ifs24050-pabwe2026-sk-p9-be-golang-gin — Delcom Blog API

## Jalankan lokal (port 8081)
```bash
psql -U postgres -c "CREATE DATABASE pabwe_p9_blog;"
cp .env.example .env        # sesuaikan DB_PASSWORD, OPENAI_KEY, OPENAI_MODEL
go mod tidy
go run cmd/server/main.go
```
- API: http://localhost:8081 · Swagger: http://localhost:8081/swagger/index.html
- Akun admin otomatis: `admin` / `admin` (ADMIN_USERNAME/ADMIN_PASSWORD) + 4 artikel contoh terbit.

## Build binary
```bash
go build -o server.exe cmd/server/main.go   # Windows -> .\server.exe
```

## Deploy
Set environment variable di hosting: `APP_ENV=production`, `APP_KEY`, `ADMIN_USERNAME`, `ADMIN_PASSWORD`,
`DATABASE_URL` (PostgreSQL hosting) **atau** `DB_HOST/DB_PORT/DB_DATABASE/DB_USERNAME/DB_PASSWORD/DB_SSLMODE`,
`SEARCH_ENGINE_ENDPOINT`, `OPENAI_ENDPOINT`, `OPENAI_KEY`, `OPENAI_MODEL`. `PORT` dari hosting otomatis dipakai.
Tanpa Postgres: `DB_CONNECTION=sqlite` + `DB_DATABASE=/data/blog.db`. Tersedia `Dockerfile`.
URL hasil deploy wajib diawali `https://ifs24050-pabwe2026-sk-p9-be-golang-gin`.
