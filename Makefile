.PHONY: build run clean test

build:
	go build -o hot-coffee .

run: build
	./hot-coffee

clean:
	rm -f hot-coffee
	rm -rf data/

test:
	go test ./...

fmt:
	gofumpt -w .

lint:
	golangci-lint run

help:
	@echo "Available commands:"
	@echo "  build  - Build the application"
	@echo "  run    - Build and run the application"
	@echo "  clean  - Remove binary and data files"
	@echo "  test   - Run tests"
	@echo "  fmt    - Format code with gofumpt"
	@echo "  lint   - Run linter"
	@echo "  help   - Show this help message"