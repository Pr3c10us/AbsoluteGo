# AbsoluteGo

A Go backend + Next.js UI for processing comic pages/panels and generating AI narration, audio, and video using Google's Gemini models.

## Prerequisites

- **Go** 1.24+
- **Node.js** 20+ (for the `ui/`)
- **Docker + Docker Compose** (for MinIO and RabbitMQ)
- **FFmpeg** on your `PATH`
- **Tesseract + Leptonica** via MSYS2 (Windows) — required for CGO OCR bindings
- **golang-migrate** CLI: `go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- **air** (hot reload, optional): `go install github.com/air-verse/air@latest`

### Windows-specific (MSYS2)

The `Makefile` expects Tesseract/Leptonica installed at `C:\msys64\mingw64`:

```bash
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-tesseract-ocr mingw-w64-x86_64-leptonica
```

## Getting a Gemini API key

The backend requires a Google Gemini API key. To get one:

1. Go to **https://aistudio.google.com/apikey**.
2. Sign in with your Google account.
3. Click **Create API key**.
4. Select an existing Google Cloud project, or let AI Studio create one for you.
5. Copy the generated key — you will paste it into `.env` as `GEMINI_API_KEY` (see below).

Notes:
- Preview models (e.g. `gemini-3.1-pro-preview`, `gemini-3.1-flash-live-preview`) may require a **paid tier** project. If you hit quota errors on the free tier, enable billing on the Cloud project linked to the key at https://console.cloud.google.com/billing.
- Keep the key out of version control — `.env` is git-ignored.

## Installation

### 1. Clone and install dependencies

```bash
git clone https://github.com/Pr3c10us/AbsoluteGo
cd AbsoluteGo
go mod download
cd ui && npm install && cd ..
```

### 2. Start infra (MinIO + RabbitMQ)

```bash
docker compose up -d
```

- MinIO console: http://localhost:9001 (user: `minioadmin`, pass: `minioadmin`)
- RabbitMQ console: http://localhost:15672 (user: `guest`, pass: `guest`)

The `minio-init` container auto-creates the required buckets (`pages`, `panels`, `comics`, `audios`, `videos`, `vabs`).

### 3. Create `.env` at the repo root

```env
PORT=:5000
ALLOWED_ORIGINS=http://localhost:3000

DATABASE_PATH="database.db"

S3_Endpoint=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_ACCESS_KEY=minioadmin

COMIC_BUCKET=comics
PAGE_BUCKET=pages
PANEL_BUCKET=panels
AUDIOS_BUCKET=audios
VIDEOS_BUCKET=videos
VABS_BUCKET=vabs

GEMINI_API_KEY=paste-your-key-here
GEMINI_MODEL=gemini-3.1-pro-preview
GEMINI_FAST_MODEL=gemini-3-pro-preview
GEMINI_LIVE_MODEL=gemini-3.1-flash-live-preview

AMQ_CONNECTION_STRING="amqp://guest:guest@localhost:5672/"

HARDWARE_ACCELERATOR=none
```

Set `HARDWARE_ACCELERATOR` to `nvidia` or `apple` if you have GPU-accelerated FFmpeg available.

### 4. Run database migrations

```bash
make db_up
```

### 5. Run the backend

```bash
make dev     # hot-reload with air
# or
make run     # build + run once
```

Backend listens on `http://localhost:5000`.

### 6. Run the UI

```bash
cd ui
npm run dev
```

UI available at `http://localhost:3000`.

## Useful Make targets

| Target | Description |
| --- | --- |
| `make dev` | Run backend with hot reload (`air`) |
| `make build` | Build backend binary to `bin/main` |
| `make run` | Build and run |
| `make db_up` | Apply all migrations |
| `make db_down` | Roll back all migrations |
| `make db_create migration=name` | Create a new migration pair |
| `make db_force version=N` | Force migration version |
