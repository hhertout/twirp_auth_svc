all: build

proto_install:
	@echo "Installing protoc..."
	@go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

generate:
	@echo "Generating protobuf files for /rpc$(file)..."
	@protoc --go_out=. --twirp_out=. rpc$(file)

migration-generate:
	@echo "Generate migration in ./migrations"
	@touch migrations/`date '+%Y-%m-%d_%s'`_migration.sql

build:
	@echo "Building..."
	@go build -o tmp/main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

update:
	@echo "Updating dependencies..."
	@go get -u ./...

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean