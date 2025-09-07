// internal/service/reports_service.go
package service

import (
	"log/slog"

	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type reportsService struct {
	orderRepo repository.OrderRepository
	menuRepo  repository.MenuRepository
}

func NewReportsService(orderRepo repository.OrderRepository, menuRepo repository.MenuRepository) ReportsService {
	return &reportsService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

func (s *reportsService) GetTotalSales() (*models.TotalSalesResponse, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		slog.Error("Failed to get orders for total sales", "error", err)
		return nil, err
	}

	var totalSales float64
	for _, order := range orders {
		if order.Status == "closed" {
			for _, orderItem := range order.Items {
				menuItem, err := s.menuRepo.GetByID(orderItem.ProductID)
				if err != nil {
					continue
				}
				if menuItem != nil {
					totalSales += menuItem.Price * float64(orderItem.Quantity)
				}
			}
		}
	}

	return &models.TotalSalesResponse{TotalSales: totalSales}, nil
}

func (s *reportsService) GetPopularItems() (*models.PopularItemsResponse, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		slog.Error("Failed to get orders for popular items", "error", err)
		return nil, err
	}

	itemCounts := make(map[string]int)
	for _, order := range orders {
		if order.Status == "closed" {
			for _, orderItem := range order.Items {
				itemCounts[orderItem.ProductID] += orderItem.Quantity
			}
		}
	}

	var popularItems []models.PopularItem
	for productID, count := range itemCounts {
		menuItem, err := s.menuRepo.GetByID(productID)
		if err != nil || menuItem == nil {
			continue
		}

		popularItems = append(popularItems, models.PopularItem{
			ProductID:   productID,
			Name:        menuItem.Name,
			TotalOrders: count,
		})
	}

	return &models.PopularItemsResponse{Items: popularItems}, nil
}
