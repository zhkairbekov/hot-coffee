// tests/integration_test.go (basic integration test example)
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"hot-coffee/internal/handler"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/service"
	"hot-coffee/models"
)

func setupTestServer(t *testing.T) *httptest.Server {
	// Create temporary directory for test data
	tempDir := t.TempDir()

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(tempDir)
	menuRepo := repository.NewMenuRepository(tempDir)
	inventoryRepo := repository.NewInventoryRepository(tempDir)

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
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	mux.HandleFunc("POST /menu", menuHandler.CreateMenuItem)
	mux.HandleFunc("GET /menu", menuHandler.GetAllMenuItems)
	mux.HandleFunc("POST /inventory", inventoryHandler.CreateInventoryItem)
	mux.HandleFunc("GET /inventory", inventoryHandler.GetAllInventoryItems)
	mux.HandleFunc("GET /reports/total-sales", reportsHandler.GetTotalSales)

	return httptest.NewServer(mux)
}

func TestBasicAPIFlow(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	client := &http.Client{}

	// 1. Add inventory item
	inventoryItem := models.InventoryItem{
		IngredientID: "test_ingredient",
		Name:         "Test Ingredient",
		Quantity:     100,
		Unit:         "units",
	}

	body, _ := json.Marshal(inventoryItem)
	req, _ := http.NewRequest("POST", server.URL+"/inventory", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// 2. Add menu item
	menuItem := models.MenuItem{
		ID:          "test_product",
		Name:        "Test Product",
		Description: "A test product",
		Price:       5.00,
		Ingredients: []models.MenuItemIngredient{
			{
				IngredientID: "test_ingredient",
				Quantity:     10,
			},
		},
	}

	body, _ = json.Marshal(menuItem)
	req, _ = http.NewRequest("POST", server.URL+"/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// 3. Create order
	order := models.Order{
		CustomerName: "Test Customer",
		Items: []models.OrderItem{
			{
				ProductID: "test_product",
				Quantity:  2,
			},
		},
	}

	body, _ = json.Marshal(order)
	req, _ = http.NewRequest("POST", server.URL+"/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// 4. Check inventory was updated
	req, _ = http.NewRequest("GET", server.URL+"/inventory", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var inventory []*models.InventoryItem
	json.NewDecoder(resp.Body).Decode(&inventory)

	if len(inventory) != 1 {
		t.Errorf("Expected 1 inventory item, got %d", len(inventory))
	}

	// Should have 80 units left (100 - 2*10)
	if inventory[0].Quantity != 80 {
		t.Errorf("Expected inventory quantity 80, got %f", inventory[0].Quantity)
	}
}

func TestMain(m *testing.M) {
	// Setup test environment
	os.Exit(m.Run())
}

// Performance and load testing example
func TestOrderCreationPerformance(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	client := &http.Client{}

	// Setup inventory and menu first
	// (similar to previous test setup...)

	// Measure performance of creating multiple orders
	orderCount := 100
	for i := 0; i < orderCount; i++ {
		order := models.Order{
			CustomerName: "Performance Test Customer",
			Items: []models.OrderItem{
				{
					ProductID: "test_product",
					Quantity:  1,
				},
			},
		}

		body, _ := json.Marshal(order)
		req, _ := http.NewRequest("POST", server.URL+"/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
	}

	t.Logf("Successfully created %d orders", orderCount)
}
