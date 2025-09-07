# Hot Coffee Project - Complete Implementation

## Project Structure
```
hot-coffee/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── handler/                   # HTTP handlers (Presentation Layer)
│   │   ├── order_handler.go
│   │   ├── menu_handler.go
│   │   ├── inventory_handler.go
│   │   ├── reports_handler.go
│   │   └── utils.go
│   ├── service/                   # Business logic (Service Layer)
│   │   ├── interfaces.go
│   │   ├── order_service.go
│   │   ├── menu_service.go
│   │   ├── inventory_service.go
│   │   └── reports_service.go
│   └── repository/                # Data access (Repository Layer)
│       ├── interfaces.go
│       ├── order_repository.go
│       ├── menu_repository.go
│       └── inventory_repository.go
├── models/                        # Data models
│   ├── order.go
│   ├── menu_item.go
│   ├── inventory_item.go
│   └── reports.go
├── tests/                         # Integration tests
│   └── integration_test.go
├── scripts/                       # Utility scripts
│   ├── init_data.sh              # Initialize sample data
│   ├── test_api.sh               # Test API endpoints
│   └── deploy.sh                 # Deployment script
├── docker/                        # Docker configuration
│   ├── Dockerfile
│   └── docker-compose.yml
├── data/                          # JSON data files (auto-created)
│   ├── orders.json
│   ├── menu_items.json
│   └── inventory.json
├── go.mod                         # Go module file
├── go.sum                         # Go dependencies (auto-generated)
├── README.md                      # Project documentation
├── Makefile                       # Build automation
├── .gitignore                     # Git ignore rules
└── config.example.json            # Configuration example
```

## How to Set Up and Run

### 1. Initialize Go Module
```bash
# Create project directory
mkdir hot-coffee
cd hot-coffee

# Initialize Go module
go mod init hot-coffee
```

### 2. Create Directory Structure
```bash
# Create directories
mkdir -p cmd
mkdir -p internal/handler
mkdir -p internal/service  
mkdir -p internal/repository
mkdir -p models
mkdir -p tests
mkdir -p scripts
mkdir -p docker
mkdir -p data
```

### 3. Copy the Code Files
Copy all the provided code into their respective files according to the structure above.

### 4. Fix Package Declarations
Each file should have the correct package declaration:
- Files in `cmd/` should have `package main`
- Files in `internal/handler/` should have `package handler`
- Files in `internal/service/` should have `package service`
- Files in `internal/repository/` should have `package repository`
- Files in `models/` should have `package models`
- Files in `tests/` should have `package tests`

### 5. Update Import Statements
The main.go file should import:
```go
import (
    "hot-coffee/internal/handler"
    "hot-coffee/internal/repository"
    "hot-coffee/internal/service"
    "hot-coffee/models"
)
```

### 6. Build and Run
```bash
# Build the application
go build -o hot-coffee ./cmd

# Run with default settings (port 8080, data dir ./data)
./hot-coffee

# Run with custom settings
./hot-coffee --port 3000 --dir ./my-data

# Show help
./hot-coffee --help
```

### 7. Initialize Sample Data (Optional)
```bash
# Make script executable
chmod +x scripts/init_data.sh

# Run initialization script
./scripts/init_data.sh
```

### 8. Test the API
```bash
# Make script executable
chmod +x scripts/test_api.sh

# Run API tests (make sure server is running first)
./scripts/test_api.sh
```

## Key Features Implemented

✅ **Three-Layered Architecture**
- Presentation Layer (HTTP Handlers)
- Business Logic Layer (Services)
- Data Access Layer (Repositories)

✅ **Complete REST API**
- Orders: Create, Read, Update, Delete, Close
- Menu Items: Full CRUD operations
- Inventory: Full CRUD operations
- Reports: Total sales, Popular items

✅ **JSON File Storage**
- Separate files for each entity
- Thread-safe operations with mutexes
- Automatic file creation

✅ **Inventory Management**
- Automatic inventory deduction on order creation
- Inventory validation before processing orders
- Detailed error messages for insufficient stock

✅ **Comprehensive Logging**
- Structured logging with slog
- Contextual information in all logs
- Error tracking and debugging support

✅ **Input Validation**
- Request validation for all endpoints
- Proper HTTP status codes
- Meaningful error messages

✅ **Command Line Interface**
- Configurable port and data directory
- Help information
- Proper error handling for startup issues

## Testing Examples

### Create Inventory Items
```bash
curl -X POST http://localhost:8080/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "ingredient_id": "espresso_shot",
    "name": "Espresso Shot", 
    "quantity": 500,
    "unit": "shots"
  }'
```

### Create Menu Items
```bash
curl -X POST http://localhost:8080/menu \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "latte",
    "name": "Caffe Latte",
    "description": "Espresso with steamed milk",
    "price": 3.50,
    "ingredients": [
      {"ingredient_id": "espresso_shot", "quantity": 1},
      {"ingredient_id": "milk", "quantity": 200}
    ]
  }'
```

### Create Orders
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "items": [
      {"product_id": "latte", "quantity": 2}
    ]
  }'
```

### Get Reports
```bash
# Total sales
curl http://localhost:8080/reports/total-sales

# Popular items
curl http://localhost:8080/reports/popular-items
```

## Important Notes

1. **Package Structure**: Make sure each file has the correct package declaration
2. **Import Paths**: Update import paths to match your module name
3. **Go Version**: Requires Go 1.23+ for the new HTTP routing features
4. **File Permissions**: Make shell scripts executable with `chmod +x`
5. **Data Directory**: The application will automatically create the data directory
6. **Thread Safety**: All repository operations are thread-safe using mutexes

This implementation fully satisfies all the requirements in the technical specification:
- ✅ Three-layered architecture
- ✅ REST API with all required endpoints
- ✅ JSON file storage
- ✅ Inventory management with automatic deduction
- ✅ Comprehensive logging
- ✅ Error handling and validation
- ✅ Command-line interface
- ✅ Aggregations and reports

The code follows Go best practices and is ready for compilation and deployment!