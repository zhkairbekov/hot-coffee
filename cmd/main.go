// cmd/main.go
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"hot-coffee/internal/repo"
	"hot-coffee/internal/router"
	"hot-coffee/internal/service"
	"hot-coffee/pkg"
)

var (
	port = flag.Int("port", 8080, "Port number")
	dir  = flag.String("dir", "data", "Path to the data directory")
	help = flag.Bool("help", false, "Show help message")
)

func main() {
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	if err := os.MkdirAll(*dir, 0755); err != nil {
		slog.Error("failed to create data directory", "err", err)
		os.Exit(1)
	}

	if *port < 1 || *port > 65535 {
		slog.Error("incorrect port number", "port", *port)
		os.Exit(1)
	}

	// Создаем репозитории
	menuRepo := repo.NewMenuRepository(*dir)
	orderRepo := repo.NewOrderRepository(*dir)
	inventoryRepo := repo.NewInventoryRepository(*dir)

	// Создаем сервисы
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo, menuRepo, inventoryRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	reportService := service.NewReportService(orderRepo, menuRepo)

	// Создаем роутер
	r := router.NewRouter(menuService, orderService, inventoryService, reportService)

	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	pkg.Graceful(srv, 5*time.Second)

	slog.Info("starting server", "addr", addr, "dataDir", *dir)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>]
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.
`)
}
