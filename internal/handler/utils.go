// internal/handler/utils.go
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"hot-coffee/models"
)

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

func validateOrder(order *models.Order) error {
	if strings.TrimSpace(order.CustomerName) == "" {
		return errors.New("customer name is required")
	}

	if len(order.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	for _, item := range order.Items {
		if strings.TrimSpace(item.ProductID) == "" {
			return errors.New("product ID is required for all items")
		}
		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than 0")
		}
	}

	return nil
}

func validateMenuItem(item *models.MenuItem) error {
	if strings.TrimSpace(item.ID) == "" {
		return errors.New("product ID is required")
	}
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("name is required")
	}
	if item.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if len(item.Ingredients) == 0 {
		return errors.New("menu item must have at least one ingredient")
	}

	for _, ingredient := range item.Ingredients {
		if strings.TrimSpace(ingredient.IngredientID) == "" {
			return errors.New("ingredient ID is required")
		}
		if ingredient.Quantity <= 0 {
			return errors.New("ingredient quantity must be greater than 0")
		}
	}

	return nil
}

func validateInventoryItem(item *models.InventoryItem) error {
	if strings.TrimSpace(item.IngredientID) == "" {
		return errors.New("ingredient ID is required")
	}
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(item.Unit) == "" {
		return errors.New("unit is required")
	}
	if item.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	return nil
}
