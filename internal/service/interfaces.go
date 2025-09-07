// internal/service/interfaces.go
package service

import "hot-coffee/models"

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrderByID(id string) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id string) error
	CloseOrder(id string) error
}

type MenuService interface {
	CreateMenuItem(item *models.MenuItem) error
	GetMenuItemByID(id string) (*models.MenuItem, error)
	GetAllMenuItems() ([]*models.MenuItem, error)
	UpdateMenuItem(item *models.MenuItem) error
	DeleteMenuItem(id string) error
}

type InventoryService interface {
	CreateInventoryItem(item *models.InventoryItem) error
	GetInventoryItemByID(id string) (*models.InventoryItem, error)
	GetAllInventoryItems() ([]*models.InventoryItem, error)
	UpdateInventoryItem(item *models.InventoryItem) error
	DeleteInventoryItem(id string) error
}

type ReportsService interface {
	GetTotalSales() (*models.TotalSalesResponse, error)
	GetPopularItems() (*models.PopularItemsResponse, error)
}
