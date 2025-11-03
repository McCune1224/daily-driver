.PHONY: help db-migrate-up db-migrate-down db-migrate-create sqlc-generate backend-run frontend-dev frontend-build clean

help:
	@echo "Aperture Science Data Portal - Available Commands"
	@echo "=================================================="
	@echo ""
	@echo "Database:"
	@echo "  db-migrate-up      Run database migrations"
	@echo "  db-migrate-down    Rollback database migrations"
	@echo "  db-migrate-create  Create new migration (use: make db-migrate-create name=migration_name)"
	@echo ""
	@echo "Backend:"
	@echo "  backend-deps       Install Go dependencies"
	@echo "  backend-run        Run Go backend server"
	@echo "  backend-build      Build Go binary"
	@echo "  sqlc-generate      Generate SQLC code"
	@echo ""
	@echo "Frontend:"
	@echo "  frontend-deps      Install npm dependencies"
	@echo "  frontend-dev       Run SvelteKit dev server"
	@echo "  frontend-build     Build SvelteKit for production"
	@echo ""
	@echo "Development:"
	@echo "  setup              Full project setup"
	@echo "  dev                Instructions for running dev servers"
	@echo "  clean              Clean build artifacts"

# Database operations

db-migrate-up:
	cd backend && migrate -database "$(shell grep -E '^DATABASE_URL=' backend/.env | cut -d '=' -f2)" -path db/migrations up

db-migrate-down:
	cd backend && migrate -database "$(shell grep -E '^DATABASE_URL=' backend/.env | cut -d '=' -f2)" -path db/migrations down

db-migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name: make db-migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	cd backend && migrate create -ext sql -dir db/migrations -seq $(name)

# SQLC
sqlc-generate:
	cd backend && sqlc generate

# Backend
backend-deps:
	cd backend && go mod download

backend-run:
	cd backend && go run cmd/api/main.go

backend-build:
	cd backend && go build -o bin/server cmd/api/main.go

# Frontend
frontend-deps:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

# Setup
setup: backend-deps frontend-deps
	@echo ""
	@echo "Setup complete! Next steps:"
	@echo "2. Create initial migration: make db-migrate-create name=init_schema"
	@echo "3. Run migrations: make db-migrate-up"
	@echo "4. Generate SQLC code: make sqlc-generate"

# Development
dev:
	@echo "Run these commands in separate terminals:"
	@echo ""
	@echo "Terminal 1 (Backend):"
	@echo "  make backend-run"
	@echo ""
	@echo "Terminal 2 (Frontend):"
	@echo "  make frontend-dev"
	@echo ""
	@echo "Then open http://localhost:5173"

# Cleanup
clean:
	rm -rf backend/bin
	rm -rf frontend/build
	rm -rf frontend/.svelte-kit
	rm -rf frontend/node_modules
