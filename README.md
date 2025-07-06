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
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── handlers/             # HTTP handlers
│   │   ├── posts.go
│   │   ├── authors.go
│   │   ├── tags.go
│   │   ├── comments.go
│   │   └── media.go
│   ├── models/              # Data models
│   │   ├── post.go
│   │   ├── author.go
│   │   ├── tag.go
│   │   ├── comment.go
│   │   └── media.go
│   ├── services/            # Business logic
│   │   ├── post_service.go
│   │   ├── author_service.go
│   │   ├── tag_service.go
│   │   ├── comment_service.go
│   │   ├── media_service.go
│   │   └── search_service.go
│   └── database/            # Database layer
│       ├── migrations/
│       └── db.go
├── api/
│   └── docs/                # API documentation
├── uploads/                 # Media files storage
├── tests/                   # Unit and integration tests
├── docker/
│   └── Dockerfile
├── .env.example
├── go.mod
├── go.sum
└── README.md
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

### Posts
- `GET /api/v1/posts` - List all posts
- `GET /api/v1/posts/{id}` - Get specific post
- `POST /api/v1/posts` - Create new post
- `PUT /api/v1/posts/{id}` - Update post
- `DELETE /api/v1/posts/{id}` - Delete post

### Authors
- `GET /api/v1/authors` - List all authors
- `GET /api/v1/authors/{id}` - Get specific author
- `POST /api/v1/authors` - Create new author
- `PUT /api/v1/authors/{id}` - Update author
- `DELETE /api/v1/authors/{id}` - Delete author

### Tags
- `GET /api/v1/tags` - List all tags
- `GET /api/v1/tags/{id}` - Get specific tag
- `POST /api/v1/tags` - Create new tag
- `PUT /api/v1/tags/{id}` - Update tag
- `DELETE /api/v1/tags/{id}` - Delete tag

### Comments
- `GET /api/v1/posts/{id}/comments` - List comments for a post
- `POST /api/v1/posts/{id}/comments` - Create new comment
- `PUT /api/v1/comments/{id}` - Update comment
- `DELETE /api/v1/comments/{id}` - Delete comment

### Media
- `POST /api/v1/media/upload` - Upload media file
- `GET /api/v1/media/{id}` - Get media file
- `DELETE /api/v1/media/{id}` - Delete media file

### Search
- `GET /api/v1/search?q={query}` - Search posts, authors, tags
- `GET /api/v1/search/posts?q={query}` - Search posts only

## Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Add tests
4. Run linters
5. Submit a pull request

## License

This project is licensed under the MIT License.
