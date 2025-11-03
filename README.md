# Aperture Science Data Portal

A daily driver web app for tracking fitness activities, tournament placements, and rotating artwork.

**Tech Stack:**
- Backend: Go + Echo Framework + SQLC + PostgreSQL
- Frontend: SvelteKit + TypeScript
- APIs: Start.GG, Art Institute of Chicago

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Make

### Installation

1. **Clone and setup**
   ```bash
   git clone <your-repo>
   cd aperture-science-portal
   make setup
   ```

2. **Configure environment**
   ```bash
   cp backend/.env.example backend/.env
   # Edit backend/.env with your database credentials and API tokens
   ```

3. **Setup database**
   ```bash
   # Create database
   make db-create

   # Copy schema to migrations
   cp backend/db/schema.sql backend/db/migrations/000001_init_schema.up.sql
   echo "DROP TABLE IF EXISTS activities, tournaments, art_pieces CASCADE;" > backend/db/migrations/000001_init_schema.down.sql

   # Run migrations
   make db-migrate-up

   # Generate SQLC code
   make sqlc-generate
   ```

4. **Run development servers**
   ```bash
   # Terminal 1: Backend
   make backend-run

   # Terminal 2: Frontend
   make frontend-dev
   ```

5. **Open** http://localhost:5173

## Project Structure

```
aperture-science-portal/
├── backend/
│   ├── cmd/api/              # Application entry point
│   ├── internal/
│   │   ├── handlers/         # HTTP handlers
│   │   ├── middleware/       # Custom middleware
│   │   ├── services/         # Business logic & API clients
│   │   └── repository/       # SQLC generated code
│   ├── db/
│   │   ├── migrations/       # Database migrations
│   │   ├── queries/          # SQLC queries
│   │   └── schema.sql        # Database schema
│   ├── pkg/                  # Public libraries
│   ├── .env                  # Environment variables (gitignored)
│   ├── go.mod                # Go dependencies
│   └── sqlc.yaml             # SQLC configuration
│
├── frontend/
│   ├── src/
│   │   ├── routes/           # SvelteKit routes
│   │   ├── lib/
│   │   │   ├── api/          # API clients
│   │   │   └── components/   # Svelte components
│   │   └── app.html          # HTML template
│   ├── static/               # Static assets
│   ├── package.json
│   ├── svelte.config.js
│   └── vite.config.ts
│
└── Makefile                  # Development commands
```

## Available Commands

See all commands:
```bash
make help
```

### Database
- `make db-create` - Create database
- `make db-migrate-up` - Run migrations
- `make db-migrate-down` - Rollback migrations
- `make db-migrate-create name=<name>` - Create new migration

### Backend
- `make backend-run` - Run server
- `make backend-build` - Build binary
- `make sqlc-generate` - Generate repository code

### Frontend
- `make frontend-dev` - Dev server (hot reload)
- `make frontend-build` - Production build

## API Endpoints

### Health
- `GET /api/health` - Health check

### Activities
- `GET /api/activities` - List activities
- `POST /api/activities` - Create activity

### Tournaments
- `GET /api/tournaments` - List tournaments
- `POST /api/tournaments` - Create tournament

### Art
- `GET /api/art/random` - Get random artwork

## Environment Variables

```env
DATABASE_URL=postgres://user:pass@localhost:5432/aperture_db?sslmode=disable
SERVER_PORT=8080
STARTGG_API_TOKEN=your_token_here
```

## Development Workflow

1. Make schema changes in `backend/db/schema.sql`
2. Create migration: `make db-migrate-create name=your_change`
3. Edit migration files in `backend/db/migrations/`
4. Run migration: `make db-migrate-up`
5. Update queries in `backend/db/queries/`
6. Regenerate code: `make sqlc-generate`

## Production Build

```bash
# Build frontend
make frontend-build

# Build backend
make backend-build

# Run
./backend/bin/server
```

## API Integration

### Start.GG
Get your API token from: https://start.gg/admin/profile/developer

### Art Institute of Chicago
No API key required. Docs: https://api.artic.edu/docs/

## License

MIT
