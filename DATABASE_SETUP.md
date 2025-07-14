# PostgreSQL Setup for GoFr Blog Service

## Prerequisites

1. **Install PostgreSQL**
   - Download from: https://www.postgresql.org/download/
   - Or use Docker: `docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres`

2. **Create Database**
   ```sql
   CREATE DATABASE gofr_blog;
   ```

## Configuration

1. **Environment Variables** (`.env` file):
   ```
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=gofr_blog
   DB_PORT=5432
   DB_DIALECT=postgres
   ```

2. **Update credentials** in `.env` file to match your PostgreSQL setup

## Running the Application

1. **Start PostgreSQL** (if not running)
2. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```
3. **Database tables** will be created automatically via GoFr migrations

## Testing Database Connection

Test endpoints:
- `GET /health` - Check if service is running
- `POST /api/v1/posts` - Create a test post
- `GET /api/v1/posts` - List posts

## Database Schema

The `posts` table will be created with:
- `id` (Primary Key)
- `title` (VARCHAR 200)
- `content` (TEXT)
- `slug` (VARCHAR 200, Unique)
- `author_id` (INTEGER)
- `status` (ENUM: draft, published, archived)
- `created_at`, `updated_at` (TIMESTAMPS)
