# Copilot Instructions for Portfolio API

## Project Overview
This is a backend API for a personal portfolio application built with Go. It manages user profiles, education, experience, projects, and authentication. The API provides RESTful endpoints for portfolio data management with JWT-based authentication.

**Tech Stack:**
- Language: Go 1.26.1
- Web Framework: Gin
- Database: PostgreSQL with GORM ORM
- Authentication: JWT (golang-jwt)
- API Documentation: Swagger/OpenAPI
- File Upload: Cloudinary
- CORS: gin-contrib/cors

## Architecture & Project Structure

### Directory Organization
- **main.go**: Entry point with Swagger configuration and server setup
- **config/**: Database connection and configuration
- **controllers/**: Business logic organized by domain (auth, education, experience, profile, projects)
- **models/**: Database models and schema setup
- **routes/**: API route definitions and setup
- **middleware/**: Authentication and request middleware
- **docs/**: Auto-generated Swagger documentation

### Controller Organization
Each feature domain has its own package under `controllers/`:
- `auth/`: Login and registration
- `profile/`: Profile data retrieval and updates
- `education/`: Education history CRUD operations
- `experience/`: Work experience CRUD operations
- `projects/`: Portfolio projects CRUD operations with file uploads

## Code Conventions

### Naming Conventions
- **Packages**: Lowercase, descriptive names (e.g., `auth`, `education`, `controllers`)
- **Functions**: PascalCase for exported functions (e.g., `Login`, `GetProfile`, `CreateEducation`)
- **Variables**: camelCase for unexported, PascalCase for exported
- **Database Models**: PascalCase (e.g., `User`, `Education`, `Experience`, `Project`)
- **Files**: Lowercase with hyphens for multi-word names OR lowercase without spaces

### Function Patterns
- Handler functions follow the Gin pattern: `func FunctionName(c *gin.Context)`
- Handlers should use `c.JSON()` for success responses and `c.JSON()` with error status for failures
- Use `c.BindJSON()` for request body parsing with validation
- Include HTTP status codes: 200 (OK), 201 (Created), 400 (Bad Request), 401 (Unauthorized), 404 (Not Found), 500 (Internal Server Error)

### Error Handling
- Return appropriate HTTP status codes
- Include meaningful error messages in JSON responses
- Log errors for debugging but don't expose internal details to clients
- Use consistent error response structure

### Database & Models
- Use GORM for all database operations
- Models are defined in `models/setup.go`
- Include proper struct tags for database columns: `gorm:"column:column_name"` and JSON: `json:"json_name"`
- Always use transactions for multi-step operations
- Implement soft deletes where appropriate using GORM's DeletedAt field

## API Design Guidelines

### Routing Structure
- All routes should be under `/api/v1` prefix
- Organize routes by resource (e.g., `/api/v1/projects`, `/api/v1/education`)
- Use HTTP methods correctly:
  - GET: Retrieve data
  - POST: Create resources
  - PUT/PATCH: Update resources
  - DELETE: Remove resources

### Authentication
- Protected routes use `middleware.AuthMiddleware()` 
- JWT token should be passed in `Authorization: Bearer <token>` header
- Admin routes are under `/api/v1/admin` prefix
- Member routes are under `/api/v1/member` prefix

### Response Format
Responses should be consistent JSON with appropriate status codes:
```go
// Success response example
c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data": dataObject,
    "message": "Operation successful",
})

// Error response example
c.JSON(http.StatusBadRequest, gin.H{
    "success": false,
    "error": "Error message",
})
```

## Development Patterns

### Controller Implementation
1. Create file in appropriate subdirectory under `controllers/`
2. Define handler function with `*gin.Context` parameter
3. Parse and validate input using `c.BindJSON()`
4. Perform business logic with database operations
5. Return JSON response with appropriate status code

### New Feature Addition Steps
1. Define database model in `models/setup.go` with GORM tags
2. Create new controller package under `controllers/`
3. Implement handler functions in controller files
4. Register routes in `routes/routes.go`
5. Add Swagger documentation comments on handler functions
6. Test endpoints before committing

### File Upload (Cloudinary)
- Use `projects/upload.go` as reference for file upload patterns
- Configure Cloudinary credentials via environment variables
- Handle upload errors gracefully
- Return upload response with public URL

## Environment Variables
Essential environment variables (typically in `.env`):
- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT token signing
- `ALLOWED_ORIGINS`: CORS allowed origins (comma-separated)
- `CLOUDINARY_URL`: Cloudinary API credentials
- Database-specific: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`

## Middleware Usage
- `middleware.AuthMiddleware()`: Validates JWT token and extracts user info
- Always use for protected routes that require authentication
- Apply at route group level using `group.Use(middleware.AuthMiddleware())`

## API Documentation
- Swagger documentation is automatically generated from comments
- Add swagger annotations to handler functions
- Generated docs available at `/swagger/index.html`
- Keep swagger comments updated when modifying endpoints

## Testing & Debugging
- Use consistent logging patterns from existing code
- Include meaningful error messages in logs
- Test with various input validation scenarios
- Verify CORS settings for frontend integration (default: http://localhost:5173)

## Common Tasks

### Adding a New CRUD Endpoint
1. Create model in `models/setup.go`
2. Create controller file (e.g., `controllers/feature/get.go`)
3. Implement handler functions
4. Add routes to `routes/routes.go`
5. Update Swagger documentation

### Modifying Database Schema
1. Update model in `models/setup.go`
2. Run migrations through GORM
3. Test database operations thoroughly

### Adding Authentication Requirement
1. Add route to appropriate auth group (`/admin` or `/member`)
2. Apply `middleware.AuthMiddleware()` to route group
3. Extract user info from context if needed

## Best Practices
- Keep controller logic focused and testable
- Use dependency injection patterns for database access
- Validate all user inputs before processing
- Use environment variables for configuration
- Follow Go fmt standards (run `go fmt ./...`)
- Handle concurrent requests safely (Gin is thread-safe)
- Use meaningful variable and function names
- Add comments for complex business logic
- Keep database queries efficient (avoid N+1 queries)
- Always validate JWT tokens on protected routes
