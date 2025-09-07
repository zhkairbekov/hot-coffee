package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"hot-coffee/internal/handler"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/service"
)

const (
	defaultPort = 8080
	defaultDir  = "./data"
)

func main() {
	var (
		port     = flag.Int("port", defaultPort, "Port number")
		dataDir  = flag.String("dir", defaultDir, "Path to the data directory")
		showHelp = flag.Bool("help", false, "Show this screen")
	)

	flag.Parse()

	if *showHelp {
		printUsage()
		return
	}

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		slog.Error("Failed to create data directory", "error", err)
		os.Exit(1)
	}

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(*dataDir)
	menuRepo := repository.NewMenuRepository(*dataDir)
	inventoryRepo := repository.NewInventoryRepository(*dataDir)

	// Initialize services
	orderService := service.NewOrderService(orderRepo, menuRepo, inventoryRepo)
	menuService := service.NewMenuService(menuRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	reportsService := service.NewReportsService(orderRepo, menuRepo)

	// Initialize handlers
	orderHandler := handler.NewOrderHandler(orderService)
	menuHandler := handler.NewMenuHandler(menuService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	reportsHandler := handler.NewReportsHandler(reportsService)

	// Setup routes
	mux := http.NewServeMux()

	// Order routes
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrder)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.UpdateOrder)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrder)
	mux.HandleFunc("POST /orders/{id}/close", orderHandler.CloseOrder)

	// Menu routes
	mux.HandleFunc("POST /menu", menuHandler.CreateMenuItem)
	mux.HandleFunc("GET /menu", menuHandler.GetAllMenuItems)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetMenuItem)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.UpdateMenuItem)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteMenuItem)

	// Inventory routes
	mux.HandleFunc("POST /inventory", inventoryHandler.CreateInventoryItem)
	mux.HandleFunc("GET /inventory", inventoryHandler.GetAllInventoryItems)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.GetInventoryItem)
	mux.HandleFunc("PUT /inventory/{id}", inventoryHandler.UpdateInventoryItem)
	mux.HandleFunc("DELETE /inventory/{id}", inventoryHandler.DeleteInventoryItem)

	// Reports routes
	mux.HandleFunc("GET /reports/total-sales", reportsHandler.GetTotalSales)
	mux.HandleFunc("GET /reports/popular-items", reportsHandler.GetPopularItems)

	addr := ":" + strconv.Itoa(*port)
	slog.Info("Starting server", "port", *port, "data_dir", *dataDir)

	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Coffee Shop Management System")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  hot-coffee [--port <N>] [--dir <S>]")
	fmt.Println("  hot-coffee --help")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help       Show this screen.")
	fmt.Println("  --port N     Port number.")
	fmt.Println("  --dir S      Path to the data directory.")
}
