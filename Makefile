TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost:5433/invitation_service_test?sslmode=disable

.PHONY: build
build:
	go build -o bin/invitation-service ./cmd/invitation-service

.PHONY: run
run:
	docker-compose up --build

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@go test -count=1 ./internal/handler/... ./internal/service/...
	@echo "Unit tests done."

.PHONY: test-integration
test-integration:
	@echo "Starting test database..."
	@docker-compose up -d --wait postgres-test
	@echo "Running integration tests..."
	@TEST_DATABASE_URL=$(TEST_DATABASE_URL) go test -count=1 ./internal/repository/... ; \
	EXIT_CODE=$$? ; \
	docker-compose rm -sf postgres-test ; \
	echo "Integration tests done." ; \
	exit $$EXIT_CODE

.PHONY: test
test: test-unit test-integration

.PHONY: down
down:
	docker-compose down

.PHONY: db-conn
db-conn:
	docker-compose exec postgres psql -U postgres -d invitation_service
