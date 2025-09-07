# Hot Coffee - Coffee Shop Management System

A RESTful API service for managing coffee shop operations including orders, menu items, and inventory.

## Features

- **Order Management**: Create, update, delete, and close orders
- **Menu Management**: Manage coffee shop menu items with ingredients
- **Inventory Management**: Track ingredient stock levels
- **Automatic Inventory Deduction**: Stock is automatically updated when orders are processed
- **Reports**: Get total sales and popular items analytics
- **JSON File Storage**: All data persisted in JSON files
- **Layered Architecture**: Clean separation between presentation, business logic, and data layers

## Quick Start

### Build the application
```bash
go build -o hot-coffee .
```

### Run with default settings
```bash
./hot-coffee
```

### Run with custom port and data directory
```bash
./hot-coffee --port 3000 --dir ./my-data
```

### Show help
```bash
./hot-coffee --help
```

## API Endpoints

### Orders
- `POST /orders` - Create a new order
- `GET /orders` - Get all orders
- `GET /orders/{id}` - Get specific order
- `PUT /orders/{id}` - Update order
- `DELETE /orders/{id}` - Delete order
- `POST /orders/{id}/close` - Close order

### Menu Items
- `POST /menu` - Add menu item
- `GET /menu` - Get all menu items
- `GET /menu/{id}` - Get specific menu item
- `PUT /menu/{id}` - Update menu item
- `DELETE /menu/{id}` - Delete menu item

### Inventory
- `POST /inventory` - Add inventory item
- `GET /inventory` - Get all inventory items
- `GET /inventory/{id}` - Get specific inventory item
- `PUT /inventory/{id}` - Update inventory item
- `DELETE /inventory/{id}` - Delete inventory item

### Reports
- `GET /reports/total-sales` - Get total sales amount
- `GET /reports/popular-items` - Get popular menu items

## Example Usage

### 1. Add Inventory Items
```bash
curl -X POST http://localhost:8080/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "ingredient_id": "espresso_shot",
    "name": "Espresso Shot",
    "quantity": 500,
    "unit": "shots"
  }'

curl -X POST http://localhost:8080/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "ingredient_id": "milk",
    "name": "Milk",
    "quantity": 5000,
    "unit": "ml"
  }'
```

### 2. Add Menu Items
```bash
curl -X POST http://localhost:8080/menu \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "latte",
    "name": "Caffe Latte",
    "description": "Espresso with steamed milk",
    "price": 3.50,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      },
      {
        "ingredient_id": "milk",
        "quantity": 200
      }
    ]
  }'
```

### 3. Create Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "items": [
      {
        "product_id": "latte",
        "quantity": 2
      }
    ]
  }'
```

### 4. Close Order
```bash
curl -X POST http://localhost:8080/orders/{order_id}/close
```

### 5. Get Reports
```bash
# Total sales
curl http://localhost:8080/reports/total-sales

# Popular items
curl http://localhost:8080/reports/popular-items
```

## Project Structure

```
hot-coffee/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handler/               # HTTP handlers (Presentation Layer)
│   │   ├── order_handler.go
│   │   ├── menu_handler.go
│   │   ├── inventory_handler.go
│   │   ├── reports_handler.go
│   │   └── utils.go
│   ├── service/               # Business logic (Service Layer)
│   │   ├── interfaces.go
│   │   ├── order_service.go
│   │   ├── menu_service.go
│   │   ├── inventory_service.go
│   │   └── reports_service.go
│   └── repository/            # Data access (Repository Layer)
│       ├── interfaces.go
│       ├── order_repository.go
│       ├── menu_repository.go
│       └── inventory_repository.go
├── models/                    # Data models
│   ├── order.go
│   ├── menu_item.go
│   ├── inventory_item.go
│   └── reports.go
├── data/                      # JSON data files (created automatically)
│   ├── orders.json
│   ├── menu_items.json
│   └── inventory.json
├── go.mod
└── README.md
```

## Architecture

The application follows a three-layered architecture:

1. **Presentation Layer (Handlers)**: Handles HTTP requests/responses and input validation
2. **Business Logic Layer (Services)**: Contains core business logic and rules
3. **Data Access Layer (Repositories)**: Manages data persistence using JSON files

## Data Storage

All data is stored in JSON files within the data directory:
- `orders.json` - Customer orders
- `menu_items.json` - Menu items with ingredients
- `inventory.json` - Ingredient inventory

## Error Handling

The application returns appropriate HTTP status codes:
- `200 OK` - Successful GET requests
- `201 Created` - Successful POST requests
- `204 No Content` - Successful DELETE requests
- `400 Bad Request` - Invalid input
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Unexpected errors

## Logging

Uses Go's `log/slog` package for structured logging with contextual information.
