# Hot Coffee - Makefile
# Coffee Shop Management System

# Variables
BINARY_NAME=hot-coffee
CMD_PATH=./cmd
DATA_DIR=./data
SCRIPTS_DIR=./scripts

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=gofumpt

# Default target
all: fmt-check test build

# Build the application
build:
	make clean
	make init-data
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build completed: ./$(BINARY_NAME)"

# Run the application
run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BINARY_NAME) --port 8080 --dir $(DATA_DIR)

# Run with custom port
run-port: build
	@echo "Starting $(BINARY_NAME) on port $(PORT)..."
	$(BINARY_NAME) --port $(PORT) --dir $(DATA_DIR)

# Initialize sample data
.PHONY: init-data
init-data:
	@echo "Initializing sample data..."
	@chmod +x $(SCRIPTS_DIR)/init_data.sh
	@$(SCRIPTS_DIR)/init_data.sh
	@echo "Sample data initialized"

# Format code
fmt:
	@echo "Formatting code..."
	@if command -v $(GOFMT) > /dev/null; then \
		$(GOFMT) -w .; \
	else \
		echo "gofumpt not found. Install with: go install mvdan.cc/gofumpt@latest"; \
		exit 1; \
	fi
	@echo "Code formatted"

# Check formatting
fmt-check:
	@echo "Checking code formatting..."
	@if command -v $(GOFMT) > /dev/null; then \
		if [ -n "$$($(GOFMT) -l .)" ]; then \
			echo "Code is not formatted. Run 'make fmt' to fix."; \
			$(GOFMT) -l .; \
			exit 1; \
		fi; \
	else \
		echo "gofumpt not found. Please install: go install mvdan.cc/gofumpt@latest"; \
		exit 1; \
	fi
	@echo "Code formatting OK"


# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BINARY_NAME)
	@rm -rf $(DATA_DIR)
	@echo "Clean completed"

# Show help
.PHONY: help
help:
	@echo "Hot Coffee - Coffee Shop Management System"
	@echo "Available targets:"
	@echo ""
	@echo "Build:"
	@echo "  build         - Build the application"
	@echo ""
	@echo "Run:"
	@echo "  run           - Build and run the application"
	@echo "  run-port      - Run with custom port (make run-port PORT=3000)"
	@echo ""
	@echo "Data:"
	@echo "  init-data     - Initialize sample data"
	@echo ""
	@echo "Code quality:"
	@echo "  fmt           - Format code with gofumpt"
	@echo "  fmt-check     - Check code formatting"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean         - Clean build artifacts and /data"
	@echo ""
	@echo "Info:"
	@echo "  help          - Show this help message"
