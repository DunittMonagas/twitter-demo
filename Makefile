.PHONY: generate-mocks test test-coverage clean-mocks

# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	@mkdir -p internal/mocks
	@$(HOME)/go/bin/mockgen -source=internal/infrastructure/repository/user.go -destination=internal/mocks/mock_user_repository.go -package=mocks
	@$(HOME)/go/bin/mockgen -source=internal/usecase/user.go -destination=internal/mocks/mock_user_usecase.go -package=mocks
	@echo "Mocks generated successfully"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...

# Clean generated mocks
clean-mocks:
	@echo "Cleaning mocks..."
	@rm -rf internal/mocks
	@echo "Mocks removed"
