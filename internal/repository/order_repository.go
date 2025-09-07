// internal/repository/order_repository.go
package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"hot-coffee/models"
)

type orderRepository struct {
	dataDir string
	mutex   sync.RWMutex
}

func NewOrderRepository(dataDir string) OrderRepository {
	return &orderRepository{
		dataDir: dataDir,
	}
}

func (r *orderRepository) getFilePath() string {
	return filepath.Join(r.dataDir, "orders.json")
}

func (r *orderRepository) loadOrders() ([]*models.Order, error) {
	filePath := r.getFilePath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []*models.Order{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var orders []*models.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) saveOrders(orders []*models.Order) error {
	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.getFilePath(), data, 0o644)
}

func (r *orderRepository) Create(order *models.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	orders, err := r.loadOrders()
	if err != nil {
		return err
	}

	orders = append(orders, order)
	return r.saveOrders(orders)
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	orders, err := r.loadOrders()
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}

	return nil, nil
}

func (r *orderRepository) GetAll() ([]*models.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.loadOrders()
}

func (r *orderRepository) Update(order *models.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	orders, err := r.loadOrders()
	if err != nil {
		return err
	}

	for i, existingOrder := range orders {
		if existingOrder.ID == order.ID {
			orders[i] = order
			return r.saveOrders(orders)
		}
	}

	return nil
}

func (r *orderRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	orders, err := r.loadOrders()
	if err != nil {
		return err
	}

	for i, order := range orders {
		if order.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			return r.saveOrders(orders)
		}
	}

	return nil
}
