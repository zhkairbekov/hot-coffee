	package repo

	import (
		"encoding/json"
		"errors"
		"os"
		"path/filepath"
		"sync"

		"hot-coffee/models"
	)

	type MenuRepository interface {
		GetAll() ([]models.MenuItem, error)
		GetByID(id string) (*models.MenuItem, error)
		Save(item models.MenuItem) error
		Update(item models.MenuItem) error
		Delete(id string) error
	}

	type menuRepositoryJSON struct {
		filePath string
		mu       sync.Mutex
	}

	func NewMenuRepository(dataDir string) MenuRepository {
		return &menuRepositoryJSON{
			filePath: filepath.Join(dataDir, "menu_items.json"),
		}
	}


	func (r *menuRepositoryJSON) readAll() ([]models.MenuItem, error) {
		r.mu.Lock()
		defer r.mu.Unlock()

		if _, err := os.Stat(r.filePath); errors.Is(err, os.ErrNotExist) {
			return []models.MenuItem{}, nil
		}

		data, err := os.ReadFile(r.filePath)
		if err != nil {
			return nil, err
		}

		var items []models.MenuItem
		if len(data) == 0 {
			return []models.MenuItem{}, nil
		}
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, err
		}
		return items, nil
	}

	func (r *menuRepositoryJSON) writeAll(items []models.MenuItem) error {
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

	// методы

	func (r *menuRepositoryJSON) GetAll() ([]models.MenuItem, error) {
		return r.readAll()
	}

	func (r *menuRepositoryJSON) GetByID(id string) (*models.MenuItem, error) {
		items, err := r.readAll()
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			if item.ID == id {
				return &item, nil
			}
		}
		return nil, errors.New("menu item not found")
	}

	func (r *menuRepositoryJSON) Save(item models.MenuItem) error {
		items, err := r.readAll()
		if err != nil {
			return err
		}
		for _, it := range items {
			if it.ID == item.ID {
				return errors.New("menu item already exists")
			}
		}
		items = append(items, item)
		return r.writeAll(items)
	}

	func (r *menuRepositoryJSON) Update(item models.MenuItem) error {
		items, err := r.readAll()
		if err != nil {
			return err
		}
		found := false
		for i, it := range items {
			if it.ID == item.ID {
				items[i] = item
				found = true
				break
			}
		}
		if !found {
			return errors.New("menu item not found")
		}
		return r.writeAll(items)
	}

	func (r *menuRepositoryJSON) Delete(id string) error {
		items, err := r.readAll()
		if err != nil {
			return err
		}
		newItems := make([]models.MenuItem, 0, len(items))
		found := false
		for _, it := range items {
			if it.ID == id {
				found = true
				continue
			}
			newItems = append(newItems, it)
		}
		if !found {
			return errors.New("menu item not found")
		}
		return r.writeAll(newItems)
	}
