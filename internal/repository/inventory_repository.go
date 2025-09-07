// internal/repository/inventory_repository.go
package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"hot-coffee/models"
)

type inventoryRepository struct {
	dataDir string
	mutex   sync.RWMutex
}

func NewInventoryRepository(dataDir string) InventoryRepository {
	return &inventoryRepository{
		dataDir: dataDir,
	}
}

func (r *inventoryRepository) getFilePath() string {
	return filepath.Join(r.dataDir, "inventory.json")
}

func (r *inventoryRepository) loadInventoryItems() ([]*models.InventoryItem, error) {
	filePath := r.getFilePath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []*models.InventoryItem{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var items []*models.InventoryItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *inventoryRepository) saveInventoryItems(items []*models.InventoryItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.getFilePath(), data, 0o644)
}

func (r *inventoryRepository) Create(item *models.InventoryItem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	items, err := r.loadInventoryItems()
	if err != nil {
		return err
	}

	items = append(items, item)
	return r.saveInventoryItems(items)
}

func (r *inventoryRepository) GetByID(id string) (*models.InventoryItem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	items, err := r.loadInventoryItems()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.IngredientID == id {
			return item, nil
		}
	}

	return nil, nil
}

func (r *inventoryRepository) GetAll() ([]*models.InventoryItem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.loadInventoryItems()
}

func (r *inventoryRepository) Update(item *models.InventoryItem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	items, err := r.loadInventoryItems()
	if err != nil {
		return err
	}

	for i, existingItem := range items {
		if existingItem.IngredientID == item.IngredientID {
			items[i] = item
			return r.saveInventoryItems(items)
		}
	}

	return nil
}

func (r *inventoryRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	items, err := r.loadInventoryItems()
	if err != nil {
		return err
	}

	for i, item := range items {
		if item.IngredientID == id {
			items = append(items[:i], items[i+1:]...)
			return r.saveInventoryItems(items)
		}
	}

	return nil
}
