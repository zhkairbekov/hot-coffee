# hot-coffee — Coffee Shop Management System

**In short:** A REST backend for a coffee shop built with Go (standard library only). Data is stored in JSON files, the project follows a three-layer architecture (Handlers → Services → Repositories), and logging is done using `log/slog`.

---

## Features

* **Orders**: Create, retrieve, update, delete, and close orders.
* **Menu**: CRUD operations for menu items (with ingredients and prices).
* **Inventory**: CRUD operations for ingredients (units and stock levels).
* **Reports**: Total sales and popular menu items.
* **Inventory Management**: Automatically checks and deducts stock when processing orders.

## Architecture

Three-layer design:

* **Presentation (handlers)** — HTTP endpoints, request validation, and response codes.
* **Business (services)** — Business logic, aggregations, and inventory rules.
* **Data Access (repositories / dal)** — Reads and writes JSON files (stored in the `data/` directory).

```
hot-coffee/
├── cmd/
│   └── main.go                 # Entry point, starts HTTP server
├── internal/
│   ├── handler/                # HTTP handlers (net/http)
│   ├── service/                # Business logic
│   └── dal/                    # JSON repositories
├── models/                     # Data models
├── data/                       # *.json files (created at first run)
├── go.mod
└── go.sum
```

## Data Models (simplified)

```go
// models/order.go
type Order struct {
    ID           string      `json:"order_id"`
    CustomerName string      `json:"customer_name"`
    Items        []OrderItem `json:"items"`
    Status       string      `json:"status"` // open|closed
    CreatedAt    string      `json:"created_at"` // RFC3339
}

type OrderItem struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}
```

```go
// models/menu_item.go
type MenuItem struct {
  ID          string               `json:"product_id"`
  Name        string               `json:"name"`
  Description string               `json:"description"`
  Price       float64              `json:"price"`
  Ingredients []MenuItemIngredient `json:"ingredients"`
}

type MenuItemIngredient struct {
  IngredientID string  `json:"ingredient_id"`
  Quantity     float64 `json:"quantity"`
}
```

```go
// models/inventory_item.go
type InventoryItem struct {
  IngredientID string  `json:"ingredient_id"`
  Name         string  `json:"name"`
  Quantity     float64 `json:"quantity"`
  Unit         string  `json:"unit"` // g|ml|shots|...
}
```

## Build and Run

Requirements: Go 1.22+, standard library only.

```bash
go build -o hot-coffee .
./hot-coffee --port 8080 --dir ./data
```

### Help

```bash
./hot-coffee --help
```

Output:

```
Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>]
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.
```

## Data Storage

Stored in the `data/` directory:

* `orders.json`
* `menu_items.json`
* `inventory.json`

Each file contains an array of objects and is created automatically.

## Logging

Uses Go's `log/slog` with Info/Warn/Error levels and context (IDs, paths, status codes, errors).

```go
slog.Info("order created", "orderID", id)
slog.Error("inventory update failed", "err", err)
```

## API Endpoints

### Orders

* `POST   /orders`
* `GET    /orders`
* `GET    /orders/{id}`
* `PUT    /orders/{id}`
* `DELETE /orders/{id}`
* `POST   /orders/{id}/close`

### Menu

* `POST   /menu`
* `GET    /menu`
* `GET    /menu/{id}`
* `PUT    /menu/{id}`
* `DELETE /menu/{id}`

### Inventory

* `POST   /inventory`
* `GET    /inventory`
* `GET    /inventory/{id}`
* `PUT    /inventory/{id}`
* `DELETE /inventory/{id}`

### Reports

* `GET /reports/total-sales` — total revenue from closed orders.
* `GET /reports/popular-items` — top-selling menu items.


## Inventory Business Rules

1. When creating or closing an order, the service checks stock for all ingredients.
2. If successful, it deducts the required amounts (`quantity * ingredient quantity`) from `inventory.json`.
3. If insufficient, the API returns 400 with details about the missing ingredient.

## Error Handling

* `200 OK` — successful GET/PUT/DELETE
* `201 Created` — successful POST
* `400 Bad Request` — validation, wrong format, or not enough stock
* `404 Not Found` — resource not found
* `500 Internal Server Error` — unexpected errors

Example:

```json
{ "error": "Insufficient inventory for ingredient 'Milk'. Required: 200ml, Available: 150ml." }
```

## User Quick Start

1. Clone the repository.
2. Run `go build -o hot-coffee .` in the root.
3. Start the server: `./hot-coffee --port 8080 --dir ./data`.
4. Add menu items and inventory via API.
5. Create orders and test reports.
