# Go API Project Setup & Migration Guide

## Part 1: Setting Up New Go Project

### Step 1: Project Structure

```bash
# Create project directory
mkdir my-new-api
cd my-new-api

# Initialize Go module
go mod init github.com/yourusername/my-new-api

# Create directory structure
mkdir -p cmd/{api,worker,migrate,scheduler}
mkdir -p internal/{handlers,services,repositories,models,middleware,database,utils}
mkdir -p {migrations,tests,config,docs,scripts}
mkdir -p tests/{integration,unit,fixtures}
```

### Step 2: Install Dependencies (First Time Only)

**Core Dependencies:**

```bash
# Web framework
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors

# Database & ORM
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Authentication & Security
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt

# Configuration
go get github.com/joho/godotenv

# Validation
go get github.com/go-playground/validator/v10

# Migration dependencies (for cmd/migrate)
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
go get github.com/lib/pq
```

**Development Tools (Install Globally):**

```bash
# Migration CLI tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Optional: Live reload for development
go install github.com/cosmtrek/air@latest

# Optional: Linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Step 3: For Subsequent Projects (After First Setup)

**Copy Template Structure:**

```bash
# Clone your template repository
git clone https://github.com/yourusername/go-api-template my-new-project
cd my-new-project

# Update module name
go mod edit -module github.com/yourusername/my-new-project

# Download dependencies (automatic)
go mod download
go mod tidy

# Update import paths in all files
find . -name "*.go" -type f -exec sed -i 's/go-api-template/my-new-project/g' {} \;

# Ready to go!
```

### Step 4: Essential Project Files

#### `.env.example`

```env
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapi_development
DB_PORT=5432
DB_SSLMODE=disable

# Server
PORT=8080
GIN_MODE=debug

# JWT Secret - Change this in production!
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

#### `Makefile`

```makefile
.PHONY: setup run migrate-up migrate-down migrate-create test build clean

# Setup project
setup:
	cp .env.example .env
	go mod download
	go mod tidy
	@echo "Setup complete! Update .env with your settings."

# Development
run:
	go run cmd/api/main.go

dev:
	air -c .air.toml

# Migrations
migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-status:
	go run cmd/migrate/main.go version

migrate-create:
	@echo "Usage: make migrate-create name=create_users_table"
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)

# Testing
test:
	go test ./...

test-coverage:
	go test -cover ./...

# Building
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/migrate cmd/migrate/main.go

clean:
	rm -rf bin/

# Database
db-create:
	createdb $(shell grep DB_NAME .env | cut -d '=' -f2)

db-drop:
	dropdb $(shell grep DB_NAME .env | cut -d '=' -f2)

db-reset: db-drop db-create migrate-up
```

#### `scripts/setup.sh`

```bash
#!/bin/bash

echo "🚀 Setting up Go API project..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Check if PostgreSQL is available
if ! command -v createdb &> /dev/null; then
    echo "❌ PostgreSQL is not available. Please install PostgreSQL."
    exit 1
fi

# Install migration tool if not present
if ! command -v migrate &> /dev/null; then
    echo "📦 Installing migration tool..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Setup environment
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ Please update .env with your database credentials"
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Create database if it doesn't exist
DB_NAME=$(grep DB_NAME .env | cut -d '=' -f2)
if [ ! -z "$DB_NAME" ]; then
    echo "🗄️  Creating database: $DB_NAME"
    createdb $DB_NAME 2>/dev/null || echo "Database might already exist"
fi

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env with your database credentials"
echo "2. Run migrations: make migrate-up"
echo "3. Start development: make run"
```

#### `.gitignore`

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/

# Test binary
*.test
*.out

# Go workspace
go.work
go.work.sum

# Environment variables
.env
.env.local
.env.production

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Log files
*.log

# Air (live reload) temp files
tmp/

# Database
*.db
*.sqlite
```

#### `README.md`

````markdown
# My Go API

A clean Go API with authentication and CRUD operations.

## Quick Start

```bash
# Clone template
git clone https://github.com/yourusername/go-api-template my-new-project
cd my-new-project

# Setup project
make setup

# Update .env with your database credentials

# Run migrations
make migrate-up

# Start development server
make run
```
````

## Development Commands

```bash
# Development
make run                    # Start API server
make dev                    # Start with live reload (requires air)

# Migrations
make migrate-up             # Run all migrations
make migrate-down           # Rollback last migration
make migrate-create name=   # Create new migration
make migrate-status         # Check migration status

# Database
make db-create              # Create database
make db-drop                # Drop database
make db-reset               # Drop, create, and migrate

# Testing
make test                   # Run tests
make test-coverage          # Run tests with coverage

# Building
make build                  # Build binaries
make clean                  # Clean build artifacts
```

---

## Part 2: Migration Management

### Migration Commands Overview

```bash
# Create new migration
make migrate-create name=create_users_table

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check status
make migrate-status

# Force version (if dirty)
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" force 1
```

### Step-by-Step Migration Workflow

#### 1. Create Migration

```bash
# Create migration files
make migrate-create name=create_users_table

# This creates:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql
```

#### 2. Write Migration SQL

**`migrations/000001_create_users_table.up.sql`**

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

**`migrations/000001_create_users_table.down.sql`**

```sql
DROP TABLE IF EXISTS users;
```

#### 3. Run Migration

```bash
# Apply migration
make migrate-up

# Check status
make migrate-status
# Output: Version: 1, Dirty: false
```

### Common Migration Patterns

#### Adding New Table

```sql
-- 000002_create_posts_table.up.sql
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
```

#### Adding New Column

```sql
-- 000003_add_email_verified_to_users.up.sql
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

-- 000003_add_email_verified_to_users.down.sql
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;
```

#### Adding Constraints

```sql
-- 000004_add_title_length_constraint.up.sql
ALTER TABLE posts ADD CONSTRAINT posts_title_length CHECK (length(title) >= 3);

-- 000004_add_title_length_constraint.down.sql
ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_title_length;
```

#### Adding Indexes

```sql
-- 000005_add_posts_title_index.up.sql
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_title ON posts(title);

-- 000005_add_posts_title_index.down.sql
DROP INDEX IF EXISTS idx_posts_title;
```

### Migration Best Practices

#### 1. Always Test Down Migrations

```bash
# Test up
make migrate-up

# Test down
make migrate-down

# Test up again
make migrate-up
```

#### 2. Use Transactions for Complex Migrations

```sql
BEGIN;

-- Multiple related changes
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';
UPDATE users SET status = 'inactive' WHERE last_login < '2023-01-01';
ALTER TABLE users ALTER COLUMN status SET NOT NULL;

COMMIT;
```

#### 3. Handle Data Carefully

```sql
-- Safe: Add column with default
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user';

-- Unsafe: Don't do this in one migration
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL;  -- ❌ Will fail if data exists
```

#### 4. Use Descriptive Names

```bash
# ✅ Good names
make migrate-create name=create_users_table
make migrate-create name=add_email_verified_to_users
make migrate-create name=create_posts_author_index

# ❌ Bad names
make migrate-create name=update_db
make migrate-create name=fix_stuff
```

### Troubleshooting Migrations

#### Dirty Database State

```bash
# Check status
make migrate-status
# Output: Version: 2, Dirty: true

# Force clean (be careful!)
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" force 2

# Try again
make migrate-up
```

#### Migration Failed

```bash
# Check what happened
psql -d myapi_development -c "\dt"  # List tables
psql -d myapi_development -c "SELECT * FROM schema_migrations;"

# Fix the migration file and force to previous version
migrate -path migrations -database "postgres://..." force 1
make migrate-up
```

#### Development Reset

```bash
# Nuclear option: reset everything
make db-reset

# This runs:
# 1. dropdb myapi_development
# 2. createdb myapi_development
# 3. make migrate-up
```

### Production Migration Strategy

#### 1. Backup First

```bash
pg_dump myapi_production > backup.sql
```

#### 2. Test on Staging

```bash
# Copy production to staging
pg_dump myapi_production | psql myapi_staging

# Test migration
make migrate-up
```

#### 3. Run with Monitoring

```bash
# Run migration with timing
time make migrate-up

# Monitor locks
# In another terminal:
psql -d myapi_production -c "SELECT * FROM pg_locks WHERE granted = false;"
```

This guide provides everything you need to set up new Go projects quickly and manage database migrations professionally!
