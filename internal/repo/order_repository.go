package repo

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"hot-coffee/models"
)

type OrderRepository interface {
	GetAll() ([]models.Order, error)
	GetByID(id string) (*models.Order, error)
	Save(order models.Order) error
	Update(order models.Order) error
	Delete(id string) error
}

type orderRepositoryJSON struct {
	filePath string
	mu       sync.Mutex
}

func NewOrderRepository(dataDir string) OrderRepository {
	return &orderRepositoryJSON{
		filePath: filepath.Join(dataDir, "orders.json"),
	}
}

func (r *orderRepositoryJSON) readAll() ([]models.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := os.Stat(r.filePath); errors.Is(err, os.ErrNotExist) {
		return []models.Order{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	if len(data) == 0 {
		return []models.Order{}, nil
	}
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepositoryJSON) writeAll(orders []models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}

	tmp := r.filePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, r.filePath)
}

func (r *orderRepositoryJSON) GetAll() ([]models.Order, error) {
	return r.readAll()
}

func (r *orderRepositoryJSON) GetByID(id string) (*models.Order, error) {
	orders, err := r.readAll()
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		if order.ID == id {
			return &order, nil
		}
	}
	return nil, errors.New("order not found")
}

func (r *orderRepositoryJSON) Save(order models.Order) error {
	orders, err := r.readAll()
	if err != nil {
		return err
	}
	for _, o := range orders {
		if o.ID == order.ID {
			return errors.New("order already exists")
		}
	}
	orders = append(orders, order)
	return r.writeAll(orders)
}

func (r *orderRepositoryJSON) Update(order models.Order) error {
	orders, err := r.readAll()
	if err != nil {
		return err
	}
	found := false
	for i, o := range orders {
		if o.ID == order.ID {
			orders[i] = order
			found = true
			break
		}
	}
	if !found {
		return errors.New("order not found")
	}
	return r.writeAll(orders)
}

func (r *orderRepositoryJSON) Delete(id string) error {
	orders, err := r.readAll()
	if err != nil {
		return err
	}
	newOrders := make([]models.Order, 0, len(orders))
	found := false
	for _, o := range orders {
		if o.ID == id {
			found = true
			continue
		}
		newOrders = append(newOrders, o)
	}
	if !found {
		return errors.New("order not found")
	}
	return r.writeAll(newOrders)
}
