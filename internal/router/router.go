// internal/router/router.go
package router

import (
	"net/http"

	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

func NewRouter(menuService *service.MenuService, orderService *service.OrderService, inventoryService *service.InventoryService, reportService *service.ReportService) http.Handler {
	mux := http.NewServeMux()

	// Создаем обработчики
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	reportHandler := handler.NewReportHandler(reportService)

	// Health check
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Menu endpoints
	mux.HandleFunc("/menu", menuHandler.HandleMenu)
	mux.HandleFunc("/menu/", menuHandler.HandleMenuByID)

	// Order endpoints
	mux.HandleFunc("/orders", orderHandler.HandleOrders)
	mux.HandleFunc("/orders/", orderHandler.HandleOrderByID)

	// Inventory endpoints
	mux.HandleFunc("/inventory", inventoryHandler.HandleInventory)
	mux.HandleFunc("/inventory/", inventoryHandler.HandleInventoryByID)

	// Report endpoints
	mux.HandleFunc("/reports/total-sales", reportHandler.HandleTotalSales)
	mux.HandleFunc("/reports/popular-items", reportHandler.HandlePopularItems)

	return mux
}
