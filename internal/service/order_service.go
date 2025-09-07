// internal/service/order_service.go
package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type orderService struct {
	orderRepo     repository.OrderRepository
	menuRepo      repository.MenuRepository
	inventoryRepo repository.InventoryRepository
}

func NewOrderService(orderRepo repository.OrderRepository, menuRepo repository.MenuRepository, inventoryRepo repository.InventoryRepository) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *orderService) CreateOrder(order *models.Order) error {
	// Generate order ID
	order.ID = generateID()
	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Check if all products exist and validate inventory
	if err := s.validateAndDeductInventory(order); err != nil {
		return err
	}

	if err := s.orderRepo.Create(order); err != nil {
		slog.Error("Failed to create order", "error", err)
		return err
	}

	slog.Info("Order created", "orderID", order.ID, "customer", order.CustomerName)
	return nil
}

func (s *orderService) GetOrderByID(id string) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		slog.Error("Failed to get order", "orderID", id, "error", err)
		return nil, err
	}
	return order, nil
}

func (s *orderService) GetAllOrders() ([]*models.Order, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		slog.Error("Failed to get all orders", "error", err)
		return nil, err
	}
	return orders, nil
}

func (s *orderService) UpdateOrder(order *models.Order) error {
	existing, err := s.orderRepo.GetByID(order.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("order not found")
	}

	if err := s.orderRepo.Update(order); err != nil {
		slog.Error("Failed to update order", "orderID", order.ID, "error", err)
		return err
	}

	slog.Info("Order updated", "orderID", order.ID)
	return nil
}

func (s *orderService) DeleteOrder(id string) error {
	if err := s.orderRepo.Delete(id); err != nil {
		slog.Error("Failed to delete order", "orderID", id, "error", err)
		return err
	}

	slog.Info("Order deleted", "orderID", id)
	return nil
}

func (s *orderService) CloseOrder(id string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	order.Status = "closed"
	if err := s.orderRepo.Update(order); err != nil {
		slog.Error("Failed to close order", "orderID", id, "error", err)
		return err
	}

	slog.Info("Order closed", "orderID", id)
	return nil
}

func (s *orderService) validateAndDeductInventory(order *models.Order) error {
	requiredIngredients := make(map[string]float64)

	// Calculate total ingredients needed
	for _, orderItem := range order.Items {
		menuItem, err := s.menuRepo.GetByID(orderItem.ProductID)
		if err != nil {
			return err
		}
		if menuItem == nil {
			return fmt.Errorf("product not found: %s", orderItem.ProductID)
		}

		for _, ingredient := range menuItem.Ingredients {
			requiredIngredients[ingredient.IngredientID] += ingredient.Quantity * float64(orderItem.Quantity)
		}
	}

	// Check inventory and deduct ingredients
	for ingredientID, requiredQty := range requiredIngredients {
		inventoryItem, err := s.inventoryRepo.GetByID(ingredientID)
		if err != nil {
			return err
		}
		if inventoryItem == nil {
			return fmt.Errorf("ingredient not found in inventory: %s", ingredientID)
		}

		if inventoryItem.Quantity < requiredQty {
			return fmt.Errorf("insufficient inventory for ingredient '%s'. Required: %.2f%s, Available: %.2f%s",
				inventoryItem.Name, requiredQty, inventoryItem.Unit, inventoryItem.Quantity, inventoryItem.Unit)
		}

		// Deduct from inventory
		inventoryItem.Quantity -= requiredQty
		if err := s.inventoryRepo.Update(inventoryItem); err != nil {
			return err
		}
	}

	return nil
}
