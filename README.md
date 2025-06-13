# 🚀 Go API Starter Template

A production-ready Go API starter template with standardized response structure, authentication, role-based permissions, and external API integration capabilities.

## ✨ Features

- **🔐 JWT Authentication** - Secure token-based authentication
- **👥 Role-Based Access Control (RBAC)** - Flexible permission system
- **🌐 External API Integration** - Built-in HTTP client with retry logic
- **📊 Standardized Response Structure** - Consistent API responses
- **🔄 API Versioning Support** - Easy version management
- **🗄️ Database Migration** - Structured database schema management
- **📝 Comprehensive Logging** - Configurable logging levels
- **⚙️ Environment Configuration** - Flexible configuration management
- **🛡️ Middleware Support** - Authentication and API middleware

## 🛠️ Installation

### 1. Clone the Project
```bash
git clone <repository-url>
cd <project-directory>
```

### 2. Initialize Go Module
```bash
go mod init <nama_projek>
go mod tidy
```

### 3. Update Project Name
Replace all occurrences of `response-std` with your `<nama_projek>` throughout the codebase.

### 4. Environment Setup
```bash
cp .env.example .env
```

Edit `.env` file with your configuration:
```env
APP_PORT=5220
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database_name

JWT_SECRET=your_jwt_secret_key

# API Configuration
API_VERSION=v1
API_BASE_URL=http://localhost:5220/api/v1

# External API Configuration
EXTERNAL_API_BASE_URL=https://api.example.com
EXTERNAL_API_KEY=your_api_key_here

# Logging & Performance
REQUEST_TIMEOUT=5S
MAX_RETRIES=3
RETRY_DELAY=200MS
ENABLE_LOGGING=TRUE
LOG_LEVEL=debug
ENVIRONMENT=development
```

## 📦 Prerequisites

### Install golang-migrate
Before running the application, you need to install `golang-migrate`:

1. Download the latest stable version from [golang-migrate releases](https://github.com/golang-migrate/migrate/releases)
2. Place `migrate.exe` in your `GOPATH/bin/` directory
3. Verify installation:
```bash
migrate -version
```
If you see a version number, you're ready to proceed!

## 🏃‍♂️ Running the Application

### 1. Run Database Migrations
```bash
make migrate-up
# or
go run cmd/migrate/migrate.go up
```

### 2. Start the Server
```bash
go run main.go
```

Your API will be available at `http://localhost:5220/api/v1`

## 🏗️ Project Structure

```
├── cmd/                    # Command line utilities
├── config/                 # Configuration files
├── core/                   # Core application logic
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Custom middleware
│   ├── models/           # Data models
│   ├── response/         # Response utilities
│   ├── router/           # Route registry
│   └── services/         # Business logic services
├── v1/                    # API version 1
│   ├── controllers/      # HTTP controllers
│   ├── database/         # Migrations and seeds
│   ├── middleware/       # Version-specific middleware
│   └── routes/           # API routes
└── v2/                    # API version 2 (future)
```

## 📚 Usage Examples

### 🔧 Standardized Response Structure

The template provides a consistent response format:

```json
{
  "status": "success|error",
  "code": 200,
  "message": "Operation successful",
  "data": {...}
}
```

### 🎯 Response Methods Available

```go
// Success Responses
response.Success(c, "Operation successful", data)      // 200
response.Created(c, "Resource created", data)          // 201
response.Accepted(c, "Request accepted", data)         // 202
response.NoContent(c)                                  // 204

// Error Responses
response.Error(c, 500, "Internal server error")       // Custom error
response.BadRequest(c, "Invalid input")                // 400
response.Unauthorized(c, "Access denied")              // 401
response.Forbidden(c, "Permission denied")             // 403
response.NotFound(c, "Resource not found")             // 404
```

### 🎯 Controller Implementation Example

```go
func (a *AuthController) Me(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        response.Unauthorized(c, "Unauthenticated")
        return
    }

    u, ok := user.(models.User)
    if !ok {
        response.NotFound(c, "User not found")
        return
    }

    data := gin.H{
        "id":    u.ID,
        "name":  u.Name,
        "email": u.Email,
        "roles": u.Roles,
    }

    response.Success(c, "User fetched successfully!", data)
}
```

## 🌐 External API Integration

### 1. Router Configuration (main.go)
```go
router.GET("/external/users/:id", apiHandler.GetExternalUser)
```

### 2. Handler Implementation (api_handler.go)
```go
func (h *APIHandler) GetExternalUser(c *gin.Context) {
    // Extract parameters
    userID := c.Param("id")
    authToken := c.GetHeader("Authorization")

    // Call service
    user, err := h.userService.GetUserFromExternalAPI(userID, authToken)
    if err != nil {
        response.Error(c, 500, err.Error())
        return
    }

    // Return standardized response
    response.Success(c, "User fetched successfully", user)
}
```

### 3. Service Implementation (user_service.go)
```go
func (s *UserService) GetUserFromExternalAPI(userID string, token string) (*User, error) {
    // Prepare request
    url := fmt.Sprintf("%s/users/%s", s.config.ExternalAPIURL, userID)
    
    apiReq := &models.APIRequest{
        Method:  "GET",
        URL:     url,
        Headers: map[string]string{"Authorization": token},
    }

    // Execute request
    res := s.apiClient.ExecuteRequest(apiReq)
    if !res.Success {
        return nil, fmt.Errorf("external API error: %s", res.Error)
    }

    // Parse response
    var user User
    if err := mapstructure.Decode(res.Data, &user); err != nil {
        return nil, err
    }

    return &user, nil
}
```

### 4. API Call Example
```http
GET /api/v1/external/users/123
Authorization: Bearer your_token_here
```

## 🔄 API Versioning

The template supports API versioning through environment configuration:

```env
API_VERSION=v1  # Change to v2, v3, etc.
```

This automatically routes requests to the appropriate version directory (`v1/`, `v2/`, etc.).

## 🗄️ Database Migrations (using Makefile)

### Available Commands
```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create name=create_new_table

# Check migration status
make migrate-status
```

### Migration Files Structure
```
v1/database/migrations/
├── 20250614000100_create_users_table.up.sql
├── 20250614000100_create_users_table.down.sql
├── 20250614000101_create_personal_access_tokens_table.up.sql
└── 20250614000101_create_personal_access_tokens_table.down.sql
```

## 🔐 Authentication & Authorization

### JWT Token Authentication
- Login endpoint generates JWT tokens
- Middleware validates tokens on protected routes
- User information is injected into request context

### Role-Based Access Control
- Users can have multiple roles
- Roles contain specific permissions
- Middleware checks permissions before allowing access

## 🚦 Middleware

### Available Middleware
- **Authentication Middleware**: Validates JWT tokens
- **API Middleware**: Handles external API requests
- **CORS Middleware**: Configures cross-origin requests
- **Logging Middleware**: Logs request/response details

## 📊 Logging

Configurable logging with multiple levels:
- `debug`: Detailed debugging information
- `info`: General information
- `warn`: Warning messages
- `error`: Error messages
- `fatal`: Fatal errors that cause program termination
- `panic`: Panic-level errors

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

---

**Happy coding! 🎉**

For issues and questions, please create an issue in the repository.