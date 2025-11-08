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

### UI Design Guidelines (Portal Aesthetic)
- **Color Palette**: Primarily white backgrounds (#FFFFFF), black text/borders (#000000), with subtle blue accents (#00BFFF or similar) for interactive elements. Avoid bright colors, gradients, or complex color schemes.
- **Borders**: Thin black lines (1px solid #000000) for all borders, buttons, inputs, and containers. Use rounded corners sparingly (border-radius: 4px max).
- **Typography**: Monospace font family (e.g., 'Courier New', 'Monaco', or 'Consolas') for all text. Use consistent font sizes: headings 24px+, body 16px, small text 12px.
- **Layout**: Clean, grid-based layouts with ample white space. Use flexbox/grid for alignment. Avoid cluttered designs; prioritize simplicity and readability.
- **Components**: Minimalist design - buttons as simple rectangles with black borders, inputs with thin borders, cards with subtle shadows (box-shadow: 0 2px 4px rgba(0,0,0,0.1)).
- **Interactive Elements**: Hover states with light blue background (#E0F7FF), active states with darker blue. No animations unless essential.
- **Icons**: Simple geometric shapes or text-based icons. Avoid complex illustrations.
- **Overall Feel**: Techy, scientific, clean interface reminiscent of Aperture Science labs - functional, unadorned, with a retro-tech vibe. 
### General
- **Formatting**: Use `gofmt` for Go, Prettier for frontend (if configured)
- **Comments**: Document exported functions and complex logic
- **Security**: Never log sensitive data, validate all inputs</content>
<parameter name="filePath">/home/mckusa/Code/aperture/AGENTS.md
