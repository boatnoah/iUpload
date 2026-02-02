# iUpload

A image processing service built with Go.

# Setup

## Environment

- `DATABASE_URL` - Postgres connection string (Supabase Postgres).
- `SUPABASE_S3_ENDPOINT` - S3-compatible endpoint for Supabase Storage.
- `SUPABASE_S3_BUCKET` - Storage bucket name.
- `SUPABASE_S3_ACCESS_KEY` - S3 access key for Supabase Storage.
- `SUPABASE_S3_SECRET_KEY` - S3 secret key for Supabase Storage.
- `JWT_SECRET` - Secret for signing JWTs.

## Local setup

1. Create a Supabase project.
2. Get Postgres connection info and Storage S3 credentials.
3. Create a Storage bucket for images.
4. Set env vars in `.env`.
5. Run the API:

```bash
make run
```

# User Authentication

    - Sign-Up: Allow users to create an account.

    - Log-In: Allow users to log into their account.

    - JWT Authentication: Secure endpoints using JWTs for authenticated access.

# Image Management

    - Upload Image: Allow users to upload images.

    - Transform Image: Allow users to perform various transformations (resize, crop, rotate, watermark etc.).

    - Retrieve Image: Allow users to retrieve a saved image in different formats.

    - List Images: List all uploaded images by the user with metadata.

# Image Transformation

    - Resize

    - Crop

    - Rotate

    - Watermark

    - Flip

    - Mirror

    - Compress

    - Change format (JPEG, PNG, etc.)

    - Apply filters (grayscale, sepia, etc

# Work plan

## Data model + SQL

- `users`: id, email, password_hash, created_at
- `sessions`: id, user_id, token, expires_at
- `uploads`: id, user_id, original_key, mime, size, width, height, created_at
- `transforms`: id, upload_id, params_json, result_key, created_at

## API endpoints

- `POST /v2/register` - create user
- `POST /v2/login` - create session and JWT
- `POST /v2/images` - upload image (multipart or signed URL flow)
- `GET /v2/images/{id}` - image metadata + access URL
- `GET /v2/images/{id}/file` - stream original image
- `POST /v2/images/{id}/transform` - create transformed image
- `GET /v2/transforms/{id}` - transform metadata + access URL

## Storage adapter (Supabase Storage)

- Use AWS S3 Go SDK v2 with `SUPABASE_S3_ENDPOINT`.
- Implement upload, download, delete, and signed URL helpers.

## Image processing

- Use `bimg` (libvips) or `imaging` for transforms.
- Validate params (resize, crop, rotate, format, quality).
- Store transform params in `params_json`.

## Security + validation

- Enforce max upload size and allowed MIME types.
- Verify image headers match MIME.
- Require JWT on image routes.
