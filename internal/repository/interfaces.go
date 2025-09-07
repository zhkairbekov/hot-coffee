// internal/repository/interfaces.go
package repository

import "hot-coffee/models"

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id string) (*models.Order, error)
	GetAll() ([]*models.Order, error)
	Update(order *models.Order) error
	Delete(id string) error
}

type MenuRepository interface {
	Create(item *models.MenuItem) error
	GetByID(id string) (*models.MenuItem, error)
	GetAll() ([]*models.MenuItem, error)
	Update(item *models.MenuItem) error
	Delete(id string) error
}

type InventoryRepository interface {
	Create(item *models.InventoryItem) error
	GetByID(id string) (*models.InventoryItem, error)
	GetAll() ([]*models.InventoryItem, error)
	Update(item *models.InventoryItem) error
	Delete(id string) error
}
