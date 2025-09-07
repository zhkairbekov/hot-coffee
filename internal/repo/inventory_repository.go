package repo

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"hot-coffee/models"
)

type InventoryRepository interface {
	GetAll() ([]models.InventoryItem, error)
	GetByID(id string) (*models.InventoryItem, error)
	Save(item models.InventoryItem) error
	Update(item models.InventoryItem) error
	Delete(id string) error
	UpdateQuantities(updates map[string]float64) error
}

type inventoryRepositoryJSON struct {
	filePath string
	mu       sync.Mutex
}

func NewInventoryRepository(dataDir string) InventoryRepository {
	return &inventoryRepositoryJSON{
		filePath: filepath.Join(dataDir, "inventory.json"),
	}
}

func (r *inventoryRepositoryJSON) readAll() ([]models.InventoryItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := os.Stat(r.filePath); errors.Is(err, os.ErrNotExist) {
		return []models.InventoryItem{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var items []models.InventoryItem
	if len(data) == 0 {
		return []models.InventoryItem{}, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *inventoryRepositoryJSON) writeAll(items []models.InventoryItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	tmp := r.filePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, r.filePath)
}

func (r *inventoryRepositoryJSON) GetAll() ([]models.InventoryItem, error) {
	return r.readAll()
}

func (r *inventoryRepositoryJSON) GetByID(id string) (*models.InventoryItem, error) {
	items, err := r.readAll()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.IngredientID == id {
			return &item, nil
		}
	}
	return nil, errors.New("inventory item not found")
}

func (r *inventoryRepositoryJSON) Save(item models.InventoryItem) error {
	items, err := r.readAll()
	if err != nil {
		return err
	}
	for _, it := range items {
		if it.IngredientID == item.IngredientID {
			return errors.New("inventory item already exists")
		}
	}
	items = append(items, item)
	return r.writeAll(items)
}

func (r *inventoryRepositoryJSON) Update(item models.InventoryItem) error {
	items, err := r.readAll()
	if err != nil {
		return err
	}
	found := false
	for i, it := range items {
		if it.IngredientID == item.IngredientID {
			items[i] = item
			found = true
			break
		}
	}
	if !found {
		return errors.New("inventory item not found")
	}
	return r.writeAll(items)
}

func (r *inventoryRepositoryJSON) Delete(id string) error {
	items, err := r.readAll()
	if err != nil {
		return err
	}
	newItems := make([]models.InventoryItem, 0, len(items))
	found := false
	for _, it := range items {
		if it.IngredientID == id {
			found = true
			continue
		}
		newItems = append(newItems, it)
	}
	if !found {
		return errors.New("inventory item not found")
	}
	return r.writeAll(newItems)
}

func (r *inventoryRepositoryJSON) UpdateQuantities(updates map[string]float64) error {
	items, err := r.readAll()
	if err != nil {
		return err
	}

	for i, item := range items {
		if delta, exists := updates[item.IngredientID]; exists {
			items[i].Quantity += delta
			if items[i].Quantity < 0 {
				items[i].Quantity = 0
			}
		}
	}

	return r.writeAll(items)
}
