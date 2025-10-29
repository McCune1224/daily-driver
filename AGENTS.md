# AGENTS.md

## Build/Run Commands
- **Dev**: `make run` (runs air with live reload)
- **Build**: `make build` (generates CSS, Templ, and SQLC)
- **Test All**: `go test ./...`
- **Test Single**: `go test -v -run TestName ./path/to/package`
- **CSS Watch**: `make css`
- **Templ Watch**: `make templ`

## Tech Stack
Go 1.24, Echo v4, Templ, TailwindCSS, HTMX, PostgreSQL (pgx/v5), SQLC, Zap (logging)

## Code Style

**Imports**: Group stdlib, external, then local packages with blank lines between groups. Use named imports for clarity when needed (e.g., `handler "daily-driver/internal/handler"`).

**Error Handling**: Return errors immediately; log with zap.Logger using structured fields (e.g., `logger.Info("msg", zap.String("key", val))`). Use `logger.Panic()` for fatal startup errors.

**Naming**: Use camelCase for private, PascalCase for exported. Prefer descriptive names over abbreviations. Constants in PascalCase (e.g., `DevPort`).

**Types**: Define struct types for domain models. Use pointers for handlers/services. Nil-check pointers before dereferencing (see service/garmin.go).

**Functions**: Keep functions focused and small. Extract helper functions for reusability. Document exported functions with comments.

**Database**: Use SQLC for type-safe queries. Pool connections via pgxpool. Define queries in `internal/db/query.sql`, schema in `internal/db/schema.sql`.

**Templates**: Use Templ for type-safe HTML. Custom `Render()` function wraps templ.Component (see handler/handler.go:38).

**Structure**: Follow existing folder layout: `cmd/` for entry, `internal/` for app code, `web/static/` for assets, `test/` for tests.

