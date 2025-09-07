package service

import (
	"hot-coffee/internal/repo"
	"hot-coffee/models"
)

type InventoryService struct {
	repo repo.InventoryRepository
}

func NewInventoryService(r repo.InventoryRepository) *InventoryService {
	return &InventoryService{repo: r}
}

func (s *InventoryService) GetAllInventory() ([]models.InventoryItem, error) {
	return s.repo.GetAll()
}

func (s *InventoryService) GetInventoryItem(id string) (*models.InventoryItem, error) {
	return s.repo.GetByID(id)
}

func (s *InventoryService) AddInventoryItem(item models.InventoryItem) error {
	return s.repo.Save(item)
}

func (s *InventoryService) UpdateInventoryItem(item models.InventoryItem) error {
	return s.repo.Update(item)
}

func (s *InventoryService) DeleteInventoryItem(id string) error {
	return s.repo.Delete(id)
}
