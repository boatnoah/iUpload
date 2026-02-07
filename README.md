# iUpload

A image processing service built with Go.

# User Authentication

    - Sign-Up: Allow users to create an account.

    - Log-In: Allow users to log into their account.

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

- id UUID PRIMARY KEY
- user_id UUID NOT NULL (owner)
- object_key TEXT NOT NULL (full key like users/<uid>/<image-id>/original.jpg)
- content_type TEXT NOT NULL (e.g. image/jpeg)
- size_bytes BIGINT NOT NULL
- created_at TIMESTAMPTZ NOT NULL DEFAULT now()
