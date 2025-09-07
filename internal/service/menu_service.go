// internal/service/menu_service.go
package service

import (
	"errors"
	"log/slog"

	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type menuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

func (s *menuService) CreateMenuItem(item *models.MenuItem) error {
	if err := s.menuRepo.Create(item); err != nil {
		slog.Error("Failed to create menu item", "error", err)
		return err
	}

	slog.Info("Menu item created", "itemID", item.ID, "name", item.Name)
	return nil
}

func (s *menuService) GetMenuItemByID(id string) (*models.MenuItem, error) {
	item, err := s.menuRepo.GetByID(id)
	if err != nil {
		slog.Error("Failed to get menu item", "itemID", id, "error", err)
		return nil, err
	}
	return item, nil
}

func (s *menuService) GetAllMenuItems() ([]*models.MenuItem, error) {
	items, err := s.menuRepo.GetAll()
	if err != nil {
		slog.Error("Failed to get all menu items", "error", err)
		return nil, err
	}
	return items, nil
}

func (s *menuService) UpdateMenuItem(item *models.MenuItem) error {
	existing, err := s.menuRepo.GetByID(item.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("menu item not found")
	}

	if err := s.menuRepo.Update(item); err != nil {
		slog.Error("Failed to update menu item", "itemID", item.ID, "error", err)
		return err
	}

	slog.Info("Menu item updated", "itemID", item.ID, "name", item.Name)
	return nil
}

func (s *menuService) DeleteMenuItem(id string) error {
	if err := s.menuRepo.Delete(id); err != nil {
		slog.Error("Failed to delete menu item", "itemID", id, "error", err)
		return err
	}

	slog.Info("Menu item deleted", "itemID", id)
	return nil
}
