# CV Platform

A minimal, production-ready monorepo for a CV upload platform featuring:

- Go backend with clean architecture, Google Cloud Storage (GCS) and Firestore adapters
- Zap structured logging and Viper-based configuration with .env support
- Next.js (App Router) frontend with a modern, accessible upload UI

## Overview

Users upload CV files from the web UI. The backend issues a signed URL (PUT to object storage). After the client uploads directly to storage, the backend finalizes the record (size, mime-type) in the database.

## Architecture

- Domain-driven layering with ports/adapters
- Direct-to-storage uploads with server-issued signed URLs
- Minimal external dependencies

## Repository Structure

```
cv-platform/
  cmd/api/                    # Backend entrypoint
  internal/
    adapter/
      gcp/                    # GCP adapters (GCS, Firestore)
      http/                   # HTTP router & handlers
    config/                   # Viper-based configuration loader
    domain/                   # Core domain models
    log/                      # Zap logger wrapper
    port/                     # Interfaces (ports)
    usecase/                  # Application use cases
  pkg/                        # Reusable packages (errors, validation)
  web/                        # Next.js frontend (App Router)
  .env                        # Backend environment (optional, local dev)
  Dockerfile
  Makefile
  README.md
```

## Tech Stack

- Backend: Go, Viper, Zap, Firestore, GCS
- Frontend: Next.js 14, React 18, Tailwind CSS
- Tooling: golangci-lint, npm

## Prerequisites

- Go (>= 1.21)
- Node.js (LTS) and npm
- GCP project (for GCS/Firestore), or stub adapters while developing

## Backend

### Configuration (env)

Place a `.env` at repo root (optional in development). Viper reads OS env; `godotenv` loads `.env` if present.

Required/optional variables:

- PORT: Backend HTTP port (default: 8080)
- LOG_LEVEL: debug | info | warn | error (default: info)
- LOG_FORMAT: json | text (default: json)
- GCP_PROJECT_ID: GCP project id (required with GCP adapters)
- GCS_BUCKET_NAME: GCS bucket for CV uploads (required with GCP adapters)
- GOOGLE_APPLICATION_CREDENTIALS: Path to service account JSON
- GOOGLE_APPLICATION_CREDENTIALS_JSON: Inline JSON credentials (alternative)

Example `.env`:

```
PORT=8080
LOG_LEVEL=debug
LOG_FORMAT=text
GCP_PROJECT_ID=your-project
GCS_BUCKET_NAME=your-bucket
# GOOGLE_APPLICATION_CREDENTIALS=/absolute/path/to/key.json
# GOOGLE_APPLICATION_CREDENTIALS_JSON={"type":"service_account",...}
```

### Run

From repo root:

```
go mod download
go run ./cmd/api
```

Zap logs are emitted to stdout. Adjust with `LOG_LEVEL` and `LOG_FORMAT`.

## Frontend (web)

### Environment

Create `web/.env.local`:

```
NEXT_PUBLIC_API_BASE=http://localhost:8080
```

### Run

```
cd web
npm install
npm run dev
# open http://localhost:3000
```

## API

Endpoints used by the UI (subject to change as handlers are implemented):

- POST `${API_BASE}/api/v1/cvs/uploads`
  - Body: `{ "file_name": string, "mime_type": string }`
  - Response: `{ "id": string, "object_key": string, "signed_url": string, "expires_at": RFC3339 }`
- POST `${API_BASE}/api/v1/cvs/{id}/complete`
  - Finalizes the upload by reading object head (size, content-type) and updating metadata

Curl examples:

```
# 1) Request signed URL
curl -s -X POST \
  -H 'Content-Type: application/json' \
  -d '{"file_name":"cv.pdf","mime_type":"application/pdf"}' \
  ${API_BASE:-http://localhost:8080}/api/v1/cvs/uploads

# 2) Upload file to signed_url (PUT)
# curl -X PUT -H 'Content-Type: application/pdf' --data-binary @cv.pdf "<signed_url>"

# 3) Finalize
curl -s -X POST ${API_BASE:-http://localhost:8080}/api/v1/cvs/<id>/complete
```

## Development Notes

- `internal/domain`: Domain entities such as `CV` and status
- `internal/port`: Interfaces for `BlobStorage` and `CVRepository`
- `internal/adapter/gcp/gcs_storage.go`: GCS implementation (signed URLs, head)
- `internal/adapter/gcp/firestore_repo.go`: Firestore `CVRepository`
- `internal/usecase/cv_upload.go`: StartUpload/CompleteUpload use cases
- `internal/adapter/http`: HTTP transport (router, handlers)
- `internal/config/config.go`: Viper config loader with .env support
- `internal/log/logger.go`: Zap logger initialization and helpers
- `web/app`: Next.js App Router pages, including `/upload`

## Linting & Testing

- Go lint: `golangci-lint run`
- Go tests: `go test ./...`
- Web lint/format: `cd web && npm run lint` (add as needed)

## Docker

Build images (examples):

```
# backend
docker build -t cv-backend:local .
# frontend
cd web && docker build -t cv-web:local .
```

## Troubleshooting

- 404 at `/`: ensure `web/app/page.tsx` exists or use redirect in `next.config.js`
- `next: command not found`: install Node.js and `npm install` inside `web`
- Missing Go modules: `go mod tidy` or `go mod download`
- GCP imports missing: `go get cloud.google.com/go/firestore google.golang.org/api`
- `@/` alias unresolved in web: ensure `web/tsconfig.json` has `baseUrl` and `paths`, then restart dev server

## License

MIT (or update as appropriate)
