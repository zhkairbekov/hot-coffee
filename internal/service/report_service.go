package service

import (
	"sort"

	"hot-coffee/internal/repo"
	"hot-coffee/models"
)

type ReportService struct {
	orderRepo repo.OrderRepository
	menuRepo  repo.MenuRepository
}

type PopularItem struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	TotalOrders int     `json:"total_orders"`
	TotalSales  float64 `json:"total_sales"`
}

func NewReportService(orderRepo repo.OrderRepository, menuRepo repo.MenuRepository) *ReportService {
	return &ReportService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

func (s *ReportService) GetTotalSales() (float64, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return 0, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return 0, err
	}

	// Создаем карту для быстрого поиска цен товаров
	menuMap := make(map[string]*models.MenuItem)
	for i, item := range menuItems {
		menuMap[item.ID] = &menuItems[i]
	}

	var totalSales float64
	for _, order := range orders {
		if order.Status == "closed" { // Считаем только закрытые заказы
			for _, orderItem := range order.Items {
				if menuItem, exists := menuMap[orderItem.ProductID]; exists {
					totalSales += menuItem.Price * float64(orderItem.Quantity)
				}
			}
		}
	}

	return totalSales, nil
}

func (s *ReportService) GetPopularItems() ([]PopularItem, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Создаем карту для быстрого поиска информации о товарах
	menuMap := make(map[string]*models.MenuItem)
	for i, item := range menuItems {
		menuMap[item.ID] = &menuItems[i]
	}

	// Считаем популярность товаров
	popularity := make(map[string]*PopularItem)

	for _, order := range orders {
		if order.Status == "closed" { // Считаем только закрытые заказы
			for _, orderItem := range order.Items {
				if menuItem, exists := menuMap[orderItem.ProductID]; exists {
					if _, itemExists := popularity[orderItem.ProductID]; !itemExists {
						popularity[orderItem.ProductID] = &PopularItem{
							ProductID: orderItem.ProductID,
							Name:      menuItem.Name,
						}
					}

					popularity[orderItem.ProductID].TotalOrders += orderItem.Quantity
					popularity[orderItem.ProductID].TotalSales += menuItem.Price * float64(orderItem.Quantity)
				}
			}
		}
	}

	// Преобразуем карту в слайс и сортируем по количеству заказов
	var result []PopularItem
	for _, item := range popularity {
		result = append(result, *item)
	}

	// Сортируем по убыванию количества заказов
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalOrders > result[j].TotalOrders
	})

	return result, nil
}
