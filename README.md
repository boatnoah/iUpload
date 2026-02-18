# iUpload (Backend)

Backend-only Go service for uploading, storing, retrieving, and transforming images. Provides REST endpoints for auth, image upload, on-the-fly transforms, and deletion. No frontend is shipped in this repository.

## Stack

- Go with chi router
- PostgreSQL for metadata (users, sessions, image records)
- S3-compatible object storage for image blobs

## Setup

Prerequisites: Go, PostgreSQL, and an S3-compatible bucket (AWS S3, MinIO, or LocalStack).

Run locally:

- `make build` to produce `bin/iupload`
- `make run` or `go run ./cmd/api` to start the HTTP API on `:3000`
- `make test` to run unit tests

## API (v1)

All responses are JSON unless downloading an image. Authenticated image routes expect `Authorization: Bearer <session_token>` header.

- `POST /v1/register` — create user
  - Body: `{ "first_name": "...", "last_name": "...", "user_name": "...", "password": "..." }`
  - Returns session token

- `POST /v1/login` — login existing user
  - Body: `{ "user_name": "...", "password": "..." }`
  - Returns session token

- `POST /v1/images` — upload JPG/PNG
  - Multipart field `image` (file)
  - Returns image metadata (id, object key, content type, created_at)

- `GET /v1/images/{id}` — fetch original image by UUID

- `DELETE /v1/images/{id}` — delete image

- `POST /v1/images/{id}/transform` — apply transforms and return bytes
  - JSON body: `{ "transformation": { "resize": {"width": 640, "height": 480}, "crop": {"width": 200, "height": 200, "x": 0, "y": 0}, "rotate": 90, "content_type": "image/png" } }`
  - Supports `resize`, `crop`, `rotate` (degrees), and optional `content_type` (`image/jpeg` or `image/png`)
