// internal/repository/menu_repository.go
package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"hot-coffee/models"
)

type menuRepository struct {
	dataDir string
	mutex   sync.RWMutex
}

func NewMenuRepository(dataDir string) MenuRepository {
	return &menuRepository{
		dataDir: dataDir,
	}
}

func (r *menuRepository) getFilePath() string {
	return filepath.Join(r.dataDir, "menu_items.json")
}

func (r *menuRepository) loadMenuItems() ([]*models.MenuItem, error) {
	filePath := r.getFilePath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []*models.MenuItem{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var items []*models.MenuItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *menuRepository) saveMenuItems(items []*models.MenuItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.getFilePath(), data, 0o644)
}

func (r *menuRepository) Create(item *models.MenuItem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	items, err := r.loadMenuItems()
	if err != nil {
		return err
	}

	items = append(items, item)
	return r.saveMenuItems(items)
}

func (r *menuRepository) GetByID(id string) (*models.MenuItem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	items, err := r.loadMenuItems()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}

	return nil, nil
}

func (r *menuRepository) GetAll() ([]*models.MenuItem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.loadMenuItems()
}

func (r *menuRepository) Update(item *models.MenuItem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	items, err := r.loadMenuItems()
	if err != nil {
		return err
	}

	for i, existingItem := range items {
		if existingItem.ID == item.ID {
			items[i] = item
			return r.saveMenuItems(items)
		}
	}

	return nil
}

func (r *menuRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	items, err := r.loadMenuItems()
	if err != nil {
		return err
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			return r.saveMenuItems(items)
		}
	}

	return nil
}
