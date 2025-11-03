# Aperture Science Portal - Agent Guidelines

## Build/Lint/Test Commands

### Frontend (SvelteKit/TypeScript)
- **Build**: `npm run build` (frontend/)
- **Dev server**: `npm run dev` (frontend/)
- **Lint/Type check**: `npm run check` (includes svelte-check)
- **Install deps**: `npm install` (frontend/)

### Backend (Go)
- **Build**: `make backend-build` or `go build -o bin/server cmd/api/main.go` (backend/)
- **Run**: `make backend-run` or `go run cmd/api/main.go` (backend/)
- **Install deps**: `make backend-deps` or `go mod download` (backend/)
- **Tests**: No test files found, but testify framework available for future tests

### Database
- **Migrations up**: `make db-migrate-up`
- **Generate SQLC**: `make sqlc-generate`

## Code Style Guidelines

### Go Backend
- **Imports**: Standard library first, then third-party, then local packages
- **Error handling**: Use `fmt.Errorf` with `%w` verb for error wrapping
- **Logging**: Use zap logger with structured logging
- **Naming**: PascalCase for exported, camelCase for unexported
- **Types**: Use struct tags for JSON marshaling (`json:"field_name"`)

### TypeScript/Svelte Frontend
- **Strict mode**: Enabled in tsconfig.json
- **Imports**: ES6 imports with relative paths (`$lib/` for lib imports)
- **Types**: Avoid `any` types - use proper TypeScript interfaces
- **Components**: PascalCase for component names, kebab-case for CSS classes
- **Error handling**: Throw errors in API functions, catch in components
- UI: Should be using TailwindCSS to design a UI aesthetic like that of the Portal Game Series (white background, thin black lines for borders, mono-spaced / techy font and design) 
### General
- **Formatting**: Use `gofmt` for Go, Prettier for frontend (if configured)
- **Comments**: Document exported functions and complex logic
- **Security**: Never log sensitive data, validate all inputs</content>
<parameter name="filePath">/home/mckusa/Code/aperture/AGENTS.md
