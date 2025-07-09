# GoFr Blog Service

A lightweight Content Management System (CMS) built with GoFr framework offering developer-focused blog functionality.

## Features

- **Markdown Content APIs**: Full CRUD operations for blog posts written in Markdown
- **Media/Image Upload Support**: File upload and management system
- **Content Management**: Posts, tags, authors, and comments management
- **SEO Metadata Support**: Built-in SEO optimization features
- **Search Functionality**: Advanced search capabilities across blog content

## Project Structure

```
├── main.go                  # Application entry point
├── handlers/                # HTTP handlers
│   └── posts.go
├── models/                  # Data models
│   └── post.go
├── services/                # Business logic
│   ├── post_service.go
│   └── post_service_test.go
├── database/                # Database layer
│   └── migrations/
│       └── 001_create_posts.sql
├── integration_test.go      # Integration tests
├── swagger.yaml             # API documentation
├── .golangci.yml           # Linter configuration
├── .env.example
├── go.mod
├── go.sum
└── bin/                    # Built binaries
    └── blog-service
```

## Development Roadmap

### PR #1: Project Foundation & Basic Post CRUD (Week 1)
- Project setup with GoFr
- Basic post model and CRUD operations
- Database setup and migrations
- Health check endpoint
- Unit tests for post operations

### PR #2: Author Management & Authentication (Week 1-2)
- Author model and CRUD operations
- Basic authentication system
- Author-post relationships
- Unit tests for author operations

### PR #3: Tags & Categorization System (Week 2)
- Tag model and CRUD operations
- Post-tag many-to-many relationships
- Tag-based filtering and search
- Unit tests for tag operations

### PR #4: Media Upload & File Management (Week 2-3)
- File upload endpoints
- Image processing and optimization
- Media model and database storage
- File validation and security
- Unit tests for media operations

### PR #5: Comments System (Week 3)
- Comment model and CRUD operations
- Comment-post relationships
- Comment moderation features
- Unit tests for comment operations

### PR #6: SEO & Search Features (Week 3-4)
- SEO metadata management
- Full-text search implementation
- Search indexing and optimization
- Advanced filtering capabilities
- Unit tests for search functionality

### PR #7: API Documentation & Production Ready (Week 4)
- Complete API documentation
- Docker containerization
- Environment configuration
- Production optimizations
- Integration tests

## Getting Started

### Prerequisites
- Go 1.21 or higher
- PostgreSQL (or your preferred database)
- Git

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd gofr-blog-service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. Run database migrations:
```bash
go run cmd/migrate/main.go
```

5. Start the server:
```bash
go run cmd/server/main.go
```

## API Endpoints

### Health Check
- `GET /health` - Service health check

### Posts (Currently Implemented)
- `GET /posts` - List all posts with pagination
- `GET /posts/{id}` - Get specific post
- `POST /posts` - Create new post
- `PUT /posts/{id}` - Update post
- `DELETE /posts/{id}` - Delete post

### Future Endpoints (Planned)
- `GET /authors` - List all authors
- `GET /authors/{id}` - Get specific author
- `POST /authors` - Create new author
- `PUT /authors/{id}` - Update author
- `DELETE /authors/{id}` - Delete author
- `GET /tags` - List all tags
- `GET /tags/{id}` - Get specific tag
- `POST /tags` - Create new tag
- `PUT /tags/{id}` - Update tag
- `DELETE /tags/{id}` - Delete tag
- `GET /posts/{id}/comments` - List comments for a post
- `POST /posts/{id}/comments` - Create new comment
- `PUT /comments/{id}` - Update comment
- `DELETE /comments/{id}` - Delete comment
- `POST /media/upload` - Upload media file
- `GET /media/{id}` - Get media file
- `DELETE /media/{id}` - Delete media file
- `GET /search?q={query}` - Search posts, authors, tags

## Getting Started

### Prerequisites
- Go 1.24 or higher
- PostgreSQL database
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/SAMMILLERR/gofr-SCO-blog-service.git
cd gofr-SCO-blog-service
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database configuration
```

4. Run database migrations:
```bash
# Migrations are automatically run on application start
```

5. Build and run the application:
```bash
# Build
go build -o bin/blog-service ./main.go

# Or run directly
go run main.go
```

The service will start on `http://localhost:8080`

### API Documentation

The API is documented using OpenAPI 3.0 specification. You can find the complete API documentation in `swagger.yaml`.

To view the documentation:
- Use any OpenAPI viewer with the `swagger.yaml` file
- Or visit online tools like [Swagger Editor](https://editor.swagger.io/) and paste the content

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run integration tests:
```bash
go test -v ./integration_test.go
```

## Code Quality

### Linting

This project uses `golangci-lint` for code quality checks. The configuration is in `.golangci.yml`.

Install golangci-lint:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Run the linter:
```bash
golangci-lint run
```

Or run with go:
```bash
go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
```

### Code Formatting

Format code using standard Go tools:
```bash
go fmt ./...
goimports -w .
```

## Contributing

1. Create a feature branch from `main`
2. Make your changes following the project structure
3. Add tests for new functionality
4. Run linters and ensure all tests pass:
   ```bash
   golangci-lint run
   go test ./...
   ```
5. Update documentation if needed
6. Submit a pull request

### Development Guidelines

- Follow GoFr framework conventions
- Use dependency injection pattern
- Write comprehensive tests
- Keep handlers thin - delegate business logic to services
- Use proper error handling
- Document all public functions and types

## License

This project is licensed under the MIT License.
