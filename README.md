# Alfred - a personal finance helper

## API Endpoints

- `GET /health` - Health check endpoint
- `POST /users` - Create a new user
- `GET /users/{username}` - Get user details
- `PUT /users/{username}` - Update user details

## Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run the test script
./test.sh

# Run tests manually
go test ./...
go test -v ./...  # verbose
go test -v -race -coverprofile=coverage.out ./...  # with coverage
```

### Test Structure

- `cmd/api/main_test.go` - Tests for main API endpoints
- `cmd/api/functions_test.go` - Tests for user handler functions
- `cmd/database/database_test.go` - Tests for database functionality

### Automated Testing

- **GitHub Actions**: Tests run automatically on push/PR to main/master branches
- **Local**: Use `make test` or `./test.sh` for quick testing
- **Coverage**: Generate HTML coverage reports with `make test-coverage`

## Development

```bash
# Test and run in development mode
make dev

# Build and run
make build
make run
```

