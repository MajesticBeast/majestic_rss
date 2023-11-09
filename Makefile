.DEFAULT_GOAL := run

.PHONY: fmt vet build clean migrate

MIGRATION_FLAG_FILE = .migration_done
IMAGE_NAME := "rss_agg:latest"

migrate: $(MIGRATION_FLAG_FILE)
	@echo "Running migrations..."

$(MIGRATION_FLAG_FILE):
	@echo "Running migrations..."
	@goose -dir "./sql/schema" postgres $(DBCONN) up
	@touch $(MIGRATION_FLAG_FILE)

fmt:
	@echo "--> Formatting code..."
	@go fmt ./...

vet: fmt
	@echo "--> Vetting code..."
	@go vet ./...

build: vet
	@echo "--> Building Docker image..."
	@docker build -t $(IMAGE_NAME) .

clean: build
	@echo "--> Cleaning up..."
	@go clean