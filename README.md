# Twitter Demo

A Twitter demo application built with Go, using clean architecture with layer separation (Repository, Usecase, Controller).

## Project Structure

```
twitter-demo/
├── cmd/                          # Application entry points
│   ├── read-api/
│   ├── write-api/
│   └── worker/
├── internal/
│   ├── config/                   # Configurations
│   ├── domain/                   # Domain entities
│   ├── infrastructure/
│   │   └── repository/          # Data access layer
│   ├── usecase/                 # Business logic
│   ├── interfaces/
│   │   ├── controller/          # HTTP controllers
│   │   └── dto/                 # Data Transfer Objects
│   └── mocks/                   # Generated mocks for testing
├── pkg/                         # Shared packages
└── database/                    # Database scripts
```

## Testing

This project includes comprehensive unit tests with mocks for all layers:

### Testing Dependencies

- **testify**: More readable assertions
- **go-sqlmock**: SQL query mocking
- **gomock**: Interface mock generation

### Running Tests

```bash
# All tests
go test ./...

# Tests with verbose mode
go test -v ./...

# Tests with coverage
go test -cover ./...

# Tests for specific layer
go test -v ./internal/infrastructure/repository/
go test -v ./internal/usecase/
go test -v ./internal/interfaces/controller/
```

### Generating Mocks

Mocks are automatically generated using `mockgen`:

```bash
# Generate all mocks
make generate-mocks

# Clean generated mocks
make clean-mocks
```

### Test Coverage

Current tests cover:

- **Repository Layer**: Tests with sqlmock for database operations
  - SelectByID (success and not found cases)
  - SelectByEmail
  - Insert (success and error cases)
  
- **Usecase Layer**: Tests with mock repository for business logic
  - CreateUser (success, duplicate email, duplicate username)
  - UpdateUser (success, user not found)
  - GetUserByID (success, error)
  
- **Controller Layer**: Tests with mock usecase for HTTP endpoints
  - GetUserByID (200, 400, 500)
  - CreateUser (201, 400, 500)
  - UpdateUser (200)
  - GetAllUsers (200)

## Available Commands

```bash
# Testing
make test              # Run all tests
make test-coverage     # Run tests with coverage
make generate-mocks    # Generate mocks
make clean-mocks       # Clean generated mocks
```