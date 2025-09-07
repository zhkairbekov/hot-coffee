// internal/service/inventory_service.go
package service

import (
	"errors"
	"log/slog"

	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *inventoryService) CreateInventoryItem(item *models.InventoryItem) error {
	if err := s.inventoryRepo.Create(item); err != nil {
		slog.Error("Failed to create inventory item", "error", err)
		return err
	}

	slog.Info("Inventory item created", "itemID", item.IngredientID, "name", item.Name)
	return nil
}

func (s *inventoryService) GetInventoryItemByID(id string) (*models.InventoryItem, error) {
	item, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		slog.Error("Failed to get inventory item", "itemID", id, "error", err)
		return nil, err
	}
	return item, nil
}

func (s *inventoryService) GetAllInventoryItems() ([]*models.InventoryItem, error) {
	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		slog.Error("Failed to get all inventory items", "error", err)
		return nil, err
	}
	return items, nil
}

func (s *inventoryService) UpdateInventoryItem(item *models.InventoryItem) error {
	existing, err := s.inventoryRepo.GetByID(item.IngredientID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("inventory item not found")
	}

	if err := s.inventoryRepo.Update(item); err != nil {
		slog.Error("Failed to update inventory item", "itemID", item.IngredientID, "error", err)
		return err
	}

	slog.Info("Inventory item updated", "itemID", item.IngredientID, "name", item.Name)
	return nil
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	if err := s.inventoryRepo.Delete(id); err != nil {
		slog.Error("Failed to delete inventory item", "itemID", id, "error", err)
		return err
	}

	slog.Info("Inventory item deleted", "itemID", id)
	return nil
}
