# Perplexity Lab Instructions - Daily Driver Web App

## Project Overview

**Project Name:** Daily-Driver Web App (Feel free to come up with something creative. I like naming my projects after Valve games, mainly Portal and Half-Life.)

**Description:** _This web app will be a collection of my personal interests and tracking data / stats for them. This includes things like running data / analysis progress (Strava Clone), Tournament History / Placements on Start.GG as well as upcoming tounraments of interest. Artsy / rotating art via things like the Chicago Institute of Art's API._

**Target Users:** _Just myself. Might be something cool to show off to friends._

**Primary Use Cases:** _Keep my self engaged with daily hobbies and that I can show off / share with friends_

---

## Tech Stack

### Backend
- **Language:** Golang
- **Web Framework:** Echo (https://echo.labstack.com/)
- **Database ORM/Query Builder:** SQLC (https://sqlc.dev/)
- **Database Driver:** SQLX (https://github.com/jmoiron/sqlx)
- **Database:** PostgreSQL

### Frontend
- **Framework:** SvelteKit (https://kit.svelte.dev/)
- **Adapter:** @sveltejs/adapter-static (for static site generation)

### Architecture Notes
- **Dual-Server Setup:** The project uses two servers:
  1. **Go Backend Server (Echo):** Acts as a REST API that handles business logic, database operations, and serves the static SvelteKit build
  2. **SvelteKit Development Server:** Used during development for hot-reloading and development features

---

## API Integrations

### 1. Chicago Art Institute API
- **Purpose:** _Display artwork collections, with details of the artwork and the artist. Some sort of rotating panel]
- **API Documentation:** https://api.artic.edu/docs/
- **Features to Implement:** _See Purpose._

### 2. Start.GG API
- **Purpose:** _Display tournament brackets, player rankings, match results_
- **API Type:** GraphQL
- **API Documentation:** https://developer.start.gg/docs/intro/
- **Authentication:** Bearer token required
- **Features to Implement:** _Tournament tracking, match history, etc._

### 3. Garmin FIT Data (muktihari/fit package)
- **Purpose:** _Strava-like running/fitness tracking and analysis, weekly running stats, weekly streak, calendar view of days ran, best/fastest runs so far (best 5k, 10k, etc)._
- **Package:** https://github.com/muktihari/fit
- **Data Sources:** _Data will come from me manually uploading .FIT files from Garmin Watch activites, which will then be ported into a SQL Database table. Can also use Strava API if one is available_
- **Features to Implement:** _Activity tracking, statistics, charts, best runs, etc._
- **Analysis Requirements:** _Pace, distance, heart rate, etc._

---

## Project Structure

### Backend (Go + Echo + SQLC)

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # Echo HTTP handlers
│   ├── middleware/              # Custom middleware (auth, logging, etc.)
│   ├── services/                # Business logic layer
│   └── repository/              # Data access layer (SQLC generated code)
├── db/
│   ├── migrations/              # Database migration files
│   ├── queries/                 # SQL queries for SQLC
│   │   ├── activities.sql
│   │   └── ...
│   └── schema.sql               # Database schema
├── pkg/                         # Shared packages
├── sqlc.yaml                    # SQLC configuration
└── go.mod
```

### Frontend (SvelteKit)

```
frontend/
├── src/
│   ├── routes/                  # SvelteKit file-based routing
│   │   ├── +page.svelte        # Homepage
│   │   ├── art/                # Art Institute features
│   │   ├── tournaments/        # Start.GG features
│   │   ├── fitness/            # Garmin/running data
│   │   └── calendar/           # Google Calendar UI
│   ├── lib/                    # Shared components and utilities
│   │   ├── components/
│   │   ├── stores/             # Svelte stores for state management
│   │   └── api/                # API client functions
│   └── app.html
├── static/                     # Static assets
├── svelte.config.js            # SvelteKit configuration with adapter-static
└── package.json
```

---

## Development Workflow

### Plase make all dev workflow uses in a Makefile

### Database Setup
1. Install PostgreSQL
2. Create database: _Whatever name you think best fits_
3. Run migrations: _golang-migrate/migrate CLI tool will be used_
4. Generate SQLC code: `sqlc generate + whatever CLI / usage should be for migrate CLI`

### Backend Development
1. Install dependencies: `go mod download`
2. Configure environment variables via .env: _Start GG API Token, Strava API if possible_
3. Run server: `go run cmd/api/main.go`
4. Server runs on: _[Fill in: port, e.g., localhost:8080]_

### Frontend Development
1. Install dependencies: `npm install`
2. Run dev server: `npm run dev`
3. Dev server runs on: _[Fill in: port, e.g., localhost:5173]_
4. Build for production: `npm run build`
5. Output directory: `build/` (static files)

### Integration
- Go backend serves the SvelteKit static build from the `/build` directory
- SvelteKit makes API calls to Go backend at: _[Fill in: API base URL, e.g., http://localhost:8080/api]_

---

## SQLC Configuration

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries/"
    schema: "db/schema.sql"
    gen:
      go:
        package: "repository"
        out: "internal/repository"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
```

_[Adjust paths and settings as needed]_

---

## SvelteKit Configuration

```javascript
import adapter from '@sveltejs/adapter-static';

const config = {
  kit: {
    adapter: adapter({
      pages: 'build',
      assets: 'build',
      fallback: 'index.html',  // For SPA mode
      precompress: false,
      strict: true
    })
  }
};

export default config;
```

---

## Authentication & Security

**Auth Method:** _None. This is being used by a single user, just don't leak API / env keys and variables_

**API Security:** _Rate limiting, CORS settings, etc._

---

## Data Models to Consider

### Activity (Fitness Data)
_Fields for running/fitness activities_

### Tournament Data
_Fields for tournament/gaming data like placements, best stats, upsets done, etc]_

### _Add more models as needed_

---

## Key Features to Implement

1. **Rotating Art Panel**
   - _Rotating Panel of Art provided by Art API, with details such as artist, history, context, etc. Only ever get collection items that contain an image to display_

2. **Tournament Tracker**
   - _Placements at tournament (showing expected seed and actual placement with an indicator if palced higher or lower), upsets, etc._

3. **Fitness Dashboard**
   - _Pretty much just copy Strava and how that app works. Raw data can be provided via .FIT files I will have uploaded via the SQL database. If API key is needed for Strava I will provide it_

---

## Custom Instructions for Perplexity

When working on this project, please:

1. **Keep answers minimal but informed** - Provide concise, technical responses without unnecessary elaboration
2. **Do not use emojis or glyphs** in responses
3. **Focus on Go and SvelteKit best practices** for the tech stack mentioned
4. **Assume familiarity with the architecture** - Don't over-explain basic concepts unless asked
5. **Provide code examples** when relevant, using proper syntax highlighting
6. **Consider the dual-server architecture** when discussing routing, API calls, or deployment
7. **Reference official documentation** for Echo, SQLC, SvelteKit, and the APIs when providing guidance
8. **Prioritize type safety** in both Go (SQLC) and TypeScript (if using in SvelteKit)

---

## Additional Resources

- **Echo Framework Docs:** https://echo.labstack.com/docs
- **SQLC Documentation:** https://docs.sqlc.dev/
- **SvelteKit Documentation:** https://kit.svelte.dev/docs
- **PostgreSQL Documentation:** https://www.postgresql.org/docs/
- **Go Project Layout:** https://github.com/golang-standards/project-layout

---

## Questions to Answer

_Use these to further refine the project scope:_

1. What is the primary goal of the application? _Show off my cool stats in a flashy website_
2. Will this be a single-user or multi-user application? _single user_
3. What level of real-time data updates is needed? _as real time as possible, maybe updates every minute or so to avoid API rate limiting myself. So long as a user is on the page it should self update every so set interval_
4. Are there any specific UI/UX requirements? _Retro tech / Valve Game UI aesthetic. Think the white panels used in the Portal Game Series._
5. What are the deployment targets? (Cloud provider, self-hosted, etc.) _Likely will be hosted via railway.app, prob easiest to use a Dockerfile to ship there._
6. What is the expected scale? (users, data volume) _1 whole user (should not be heavy usage besides )_

---

## Notes
