# Deploying daily-driver to Railway

This project ships with a production-ready multi-stage Dockerfile. Railway will build from the Dockerfile and run the produced container.

## Prerequisites
- Railway account and CLI installed: https://docs.railway.app/quick-start
- A PostgreSQL database (Railway plugin recommended)

## How the container runs
- The server binds to the PORT environment variable, defaulting to 8080 when unset.
- The Dockerfile sets EXPOSE 8080. Railway automatically injects PORT in the container; no extra config needed.
- The app requires DATABASE_URL to connect to Postgres.
- Static assets are served from /static and are copied into the image from web/static at build time.

## One-time set up (CLI)
```bash
# Login, initialize and link the project
railway login
railway init
railway link # if this repo already has a Railway project

# Add a PostgreSQL database (creates DATABASE_URL automatically)
railway add --plugin postgresql

# Deploy using the Dockerfile in the repo root
railway up
```

Notes:
- When using a Dockerfile, Railway will use the image entrypoint defined there (no "Start Command" necessary).
- DATABASE_URL will be injected automatically by the Postgres plugin. You can inspect with:
  railway variables

## Dashboard flow
1) Create a new project or open your existing project
2) Add the Postgres plugin (Resources -> New -> PostgreSQL)
3) Create a new Service from your GitHub repo (set Root Directory to repository root where Dockerfile lives)
4) No Start Command needed since the Dockerfile has an ENTRYPOINT
5) Deploy

## Local test (optional)
Build and run locally to verify:
```bash
# Build image
docker build -t daily-driver .

# Run with env vars
# Replace the example DATABASE_URL with your own connection string
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=disable" \
  daily-driver
```

Visit http://localhost:8080

## Database schema
This repo contains a schema at internal/db/schema.sql. Apply it to your database once (if not already applied):
```bash
# Using psql
psql "$DATABASE_URL" -f internal/db/schema.sql
```

## Tailwind and templ
- Templ templates are generated during docker build: `templ generate`
- Tailwind CSS output (web/static/css/output.css) is committed to the repo. If you want Docker to build Tailwind too, you can add a step to download the Tailwind standalone binary in the builder stage and run:
  tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --minify

The provided Dockerfile is sufficient as long as output.css is committed or otherwise present at build time.

