# golang-s3-microservices

A small example microservices project in Go demonstrating generating S3 presigned upload URLs, verifying uploads, and creating products that reference uploaded images. The project runs two services:

- `app-service` - API for requesting upload URLs, creating products and (for testing) generating JWT tokens.
- `storage-service` - Service responsible for verifying upload metadata and generating S3 presigned PUT URLs.

## Features

- Request presigned S3 upload URLs (server-side verification of metadata signature)
- Upload files directly to S3 using presigned PUT URLs
- Verify uploaded images and mark them as available
- Create product resources referencing uploaded images
- Dockerized services for easy local development and testing
- Small test harness to exercise the full flow with a `test.png`

## Tech stack

- Go (services are written in Go modules)
- Gin HTTP framework
- AWS SDK for Go v2 (S3 presigned URLs)
- Docker & Docker Compose for running services locally
- jq and curl used in test scripts

## Prerequisites

- Docker & Docker Compose installed on your machine
- Go (if you want to run services locally without Docker)
- An AWS account or a local S3-compatible endpoint and valid credentials in environment variables
- `jq` for parsing JSON in test scripts (optional but recommended)

## Environment / Configuration

Copy or edit the `.env` file in the repository root (a sample `.env` already exists). Key variables:

- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret
- `AWS_DEFAULT_REGION` - S3 region
- `AWS_BUCKET` - S3 bucket name used for demo uploads
- `AWS_USE_PATH_STYLE_ENDPOINT` - true/false for path-style S3 endpoints
- `JWT_SECRET` - secret used to sign and verify JWTs and metadata HMACs
- `JWT_EXPIRATION` - token expiration (e.g. `24h`)
- `APP_PORT` - app-service port (default `8080`)
- `STORAGE_SERVICE_URL` - URL for storage-service (used by app-service), usually `http://storage-service:8081` when running with Docker Compose

Important: keep `JWT_SECRET` safe. The repository includes a testing endpoint that generates tokens — remove or protect this for production.

## Run with Docker Compose (recommended)

From the repository root:

```bash
# build and start both services
docker-compose up --build -d

# follow logs
docker-compose logs -f

# stop and remove
docker-compose down
```

If port 8080 is already used on your machine, stop the process using it or update `APP_PORT` and the `ports` mapping in `docker-compose.yml`.

## Run services locally (without Docker)

Each service has its own `go.mod` and can be run individually. Example for `app-service`:

```bash
cd app-service
go mod tidy
go run .
```

Do the same in `storage-service`.

## APIs

The following endpoints are available (paths are for the local Docker setup):

### 1) Generate test JWT (development only)

POST /auth/token

Response:

```json
{ "token": "<signed-jwt>" }
```

Use this token for authenticated endpoints during local testing.

### 2) Request presigned upload URL

POST /upload-url

Headers:
- Authorization: Bearer <token>
- Content-Type: application/json

Body:

```json
{ "filename": "test.png", "size": 12345, "content_type": "image/png" }
```

Response (200):

```json
{ "upload_url": "https://...", "image_id": "123" }
```

Use the `upload_url` to perform a PUT with the file contents (Content-Type should match).

Example (uploaded file `test.png`):

```bash
curl -X PUT "<upload_url>" -H "Content-Type: image/png" --upload-file test.png
```

### 3) Verify upload

GET /verify/:id

Response example:

```json
{ "valid": true }
```

This endpoint confirms whether the upload metadata and presence are valid.

### 4) Create product

POST /products

Headers: 
- Content-Type: application/json
- Authorization: Bearer <token>

Body example:

```json
{
  "name": "Test Product",
  "description": "A test product",
  "price": 99.99,
  "image_id": "123"
}
```

Response example:

```json
{ "id": "123", "name": "Test Product", "image_id": "123", "price": 99.99 }
```

## Testing the full flow

A test harness `test-api.sh` has been added to the project root. It:

1. Requests a test JWT from `/auth/token`
2. Requests a presigned upload URL from `/upload-url`
3. Uploads `test.png` to the returned URL
4. Calls `/verify/:id` to confirm the upload
5. Calls `/products` to create a product referencing the upload

Run it after starting the services (Docker Compose):

```bash
./test-api.sh
```

There is also `test-presigned.sh` for quick presigned URL testing.

## Notes and security

- `/auth/token` is intentionally provided for local testing. Do not expose this in production.
- `TEST_MODE` environment variable is used to bypass JWT validation in local testing flows. Remove or set to `false` in production.
- The metadata HMAC and JWT use the same `JWT_SECRET` in this demo; in production you might separate secrets and use stronger guardrails.

## Prepare repository for GitHub

1. Initialize a git repo (if not already):

```bash
git init
git add .
git commit -m "Initial commit: golang s3 microservices demo"
```

2. Create a new repository on GitHub and push:

```bash
# replace origin URL with your repo
git remote add origin git@github.com:youruser/golang-s3-microservices.git
git branch -M main
git push -u origin main
```

Add a `.gitignore` if you'd like to exclude binaries and `.env` (recommended) — keep secrets out of VCS.

## Next steps / improvements

- Add automated integration tests in CI that run the `test-api.sh` flow in a controlled environment.
- Add proper storage (DB) for products instead of in-memory maps.
- Harden security: remove test endpoints, rotate secrets, verify CORS, don't trust proxies blindly.

## License

MIT
