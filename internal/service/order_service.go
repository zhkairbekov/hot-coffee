package service

import (
	"fmt"
	"log/slog"
	"time"

	"hot-coffee/internal/repo"
	"hot-coffee/models"
)

type OrderService struct {
	orderRepo     repo.OrderRepository
	menuRepo      repo.MenuRepository
	inventoryRepo repo.InventoryRepository
}

func NewOrderService(orderRepo repo.OrderRepository, menuRepo repo.MenuRepository, inventoryRepo repo.InventoryRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *OrderService) GetOrder(id string) (*models.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) CreateOrder(order models.Order) error {
	// Генерируем ID и устанавливаем время создания
	if order.ID == "" {
		order.ID = fmt.Sprintf("order_%d", time.Now().UnixNano())
	}
	order.CreatedAt = time.Now().Format(time.RFC3339)
	order.Status = "open"

	// Проверяем наличие всех продуктов в меню и рассчитываем необходимые ингредиенты
	requiredIngredients := make(map[string]float64)

	for _, orderItem := range order.Items {
		menuItem, err := s.menuRepo.GetByID(orderItem.ProductID)
		if err != nil {
			slog.Error("menu item not found", "productID", orderItem.ProductID, "err", err)
			return fmt.Errorf("product not found: %s", orderItem.ProductID)
		}

		// Рассчитываем общее количество ингредиентов для данного товара
		for _, ingredient := range menuItem.Ingredients {
			requiredIngredients[ingredient.IngredientID] += ingredient.Quantity * float64(orderItem.Quantity)
		}
	}

	// Проверяем наличие ингредиентов в инвентаре
	if err := s.checkInventoryAvailability(requiredIngredients); err != nil {
		return err
	}

	// Списываем ингредиенты из инвентаря
	updates := make(map[string]float64)
	for ingredientID, quantity := range requiredIngredients {
		updates[ingredientID] = -quantity
	}

	if err := s.inventoryRepo.UpdateQuantities(updates); err != nil {
		slog.Error("failed to update inventory", "err", err)
		return fmt.Errorf("failed to update inventory: %v", err)
	}

	// Сохраняем заказ
	if err := s.orderRepo.Save(order); err != nil {
		slog.Error("failed to save order", "orderID", order.ID, "err", err)
		return err
	}

	slog.Info("order created", "orderID", order.ID, "customerName", order.CustomerName)
	return nil
}

func (s *OrderService) UpdateOrder(order models.Order) error {
	existingOrder, err := s.orderRepo.GetByID(order.ID)
	if err != nil {
		return err
	}

	// Если заказ уже закрыт, запрещаем изменения
	if existingOrder.Status == "closed" {
		return fmt.Errorf("cannot update closed order")
	}

	order.CreatedAt = existingOrder.CreatedAt
	return s.orderRepo.Update(order)
}

func (s *OrderService) DeleteOrder(id string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Если заказ был открыт, возвращаем ингредиенты в инвентарь
	if order.Status == "open" {
		if err := s.returnIngredientsToInventory(order); err != nil {
			slog.Warn("failed to return ingredients to inventory", "orderID", id, "err", err)
		}
	}

	return s.orderRepo.Delete(id)
}

func (s *OrderService) CloseOrder(id string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status == "closed" {
		return fmt.Errorf("order already closed")
	}

	order.Status = "closed"
	slog.Info("order closed", "orderID", id)
	return s.orderRepo.Update(*order)
}

func (s *OrderService) checkInventoryAvailability(requiredIngredients map[string]float64) error {
	inventoryItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check inventory: %v", err)
	}

	inventoryMap := make(map[string]*models.InventoryItem)
	for i, item := range inventoryItems {
		inventoryMap[item.IngredientID] = &inventoryItems[i]
	}

	for ingredientID, requiredQuantity := range requiredIngredients {
		inventoryItem, exists := inventoryMap[ingredientID]
		if !exists {
			return fmt.Errorf("ingredient not found in inventory: %s", ingredientID)
		}

		if inventoryItem.Quantity < requiredQuantity {
			return fmt.Errorf("insufficient inventory for ingredient '%s'. Required: %.2f%s, Available: %.2f%s",
				inventoryItem.Name, requiredQuantity, inventoryItem.Unit, inventoryItem.Quantity, inventoryItem.Unit)
		}
	}

	return nil
}

func (s *OrderService) returnIngredientsToInventory(order *models.Order) error {
	requiredIngredients := make(map[string]float64)

	for _, orderItem := range order.Items {
		menuItem, err := s.menuRepo.GetByID(orderItem.ProductID)
		if err != nil {
			continue // Пропускаем, если товар не найден
		}

		for _, ingredient := range menuItem.Ingredients {
			requiredIngredients[ingredient.IngredientID] += ingredient.Quantity * float64(orderItem.Quantity)
		}
	}

	// Возвращаем ингредиенты (положительные значения)
	updates := make(map[string]float64)
	for ingredientID, quantity := range requiredIngredients {
		updates[ingredientID] = quantity
	}

	return s.inventoryRepo.UpdateQuantities(updates)
}
