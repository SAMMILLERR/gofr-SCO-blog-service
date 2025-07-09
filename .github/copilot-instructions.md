<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# GoFr Blog Service - Copilot Instructions

This is a GoFr-based blog service project. Please follow these guidelines when generating code:

## Framework Guidelines
- Use GoFr framework patterns and best practices
- Follow GoFr's context-based handler pattern: `func(ctx *gofr.Context) (interface{}, error)`
- Use GoFr's built-in database, logging, and middleware features
- Follow RESTful API design principles

## Code Style
- Use Go's standard naming conventions (CamelCase for exported, camelCase for unexported)
- Add comprehensive error handling
- Include proper struct tags for JSON serialization
- Use dependency injection pattern for services
- Write unit tests for all business logic

## Database
- Use GoFr's built-in database abstraction
- Include proper database migrations
- Use prepared statements for security
- Implement proper indexing for search functionality

## Security
- Validate all input data
- Implement proper authentication and authorization
- Use secure file upload practices
- Sanitize HTML content when needed

## Project Structure
- Keep handlers thin - delegate business logic to services
- Use models for data structures
- Implement proper separation of concerns
- Include comprehensive testing at all levels

## API Design
- Use consistent response formats
- Implement proper HTTP status codes
- Include pagination for list endpoints
- Add proper API versioning (v1 prefix)
- Document all endpoints clearly
