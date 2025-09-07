// internal/handler/inventory_handler.go
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Warn("Invalid JSON in create inventory item request", "error", err)
		writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validateInventoryItem(&item); err != nil {
		slog.Warn("Inventory item validation failed", "error", err)
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.inventoryService.CreateInventoryItem(&item); err != nil {
		slog.Error("Failed to create inventory item", "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		slog.Error("Failed to get all inventory items", "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	item, err := h.inventoryService.GetInventoryItemByID(id)
	if err != nil {
		slog.Error("Failed to get inventory item", "itemID", id, "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if item == nil {
		writeErrorResponse(w, "Inventory item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Warn("Invalid JSON in update inventory item request", "error", err)
		writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	item.IngredientID = id
	if err := validateInventoryItem(&item); err != nil {
		slog.Warn("Inventory item validation failed", "error", err)
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.inventoryService.UpdateInventoryItem(&item); err != nil {
		slog.Error("Failed to update inventory item", "itemID", id, "error", err)
		if err.Error() == "inventory item not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
		} else {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if err := h.inventoryService.DeleteInventoryItem(id); err != nil {
		slog.Error("Failed to delete inventory item", "itemID", id, "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
