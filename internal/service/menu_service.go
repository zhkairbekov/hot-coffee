package service

import (
	"hot-coffee/internal/repo"
	"hot-coffee/models"
)

type MenuService struct {
	repo repo.MenuRepository 
}

func NewMenuService(r repo.MenuRepository) *MenuService {
	return &MenuService{repo: r}
}

func (s *MenuService) GetMenu() ([]models.MenuItem, error) {
	return s.repo.GetAll()
}

func (s *MenuService) GetMenuItem(id string) (*models.MenuItem, error) {
	return s.repo.GetByID(id)
}

func (s *MenuService) AddMenuItem(item models.MenuItem) error {
	return s.repo.Save(item)
}

func (s *MenuService) UpdateMenuItem(item models.MenuItem) error {
	return s.repo.Update(item)
}

func (s *MenuService) DeleteMenuItem(id string) error {
	return s.repo.Delete(id)
}

