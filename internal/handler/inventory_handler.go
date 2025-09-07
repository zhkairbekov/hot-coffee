package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) HandleInventory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllInventory(w, r)
	case http.MethodPost:
		h.createInventoryItem(w, r)
	default:
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) HandleInventoryByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/inventory/")
	if id == "" {
		writeErrorJSON(w, "ingredient ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getInventoryItem(w, r, id)
	case http.MethodPut:
		h.updateInventoryItem(w, r, id)
	case http.MethodDelete:
		h.deleteInventoryItem(w, r, id)
	default:
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) getAllInventory(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllInventory()
	if err != nil {
		writeErrorJSON(w, "failed to retrieve inventory", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}

func (h *InventoryHandler) createInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErrorJSON(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация
	if item.IngredientID == "" {
		writeErrorJSON(w, "ingredient ID is required", http.StatusBadRequest)
		return
	}
	if item.Name == "" {
		writeErrorJSON(w, "ingredient name is required", http.StatusBadRequest)
		return
	}
	if item.Unit == "" {
		writeErrorJSON(w, "unit is required", http.StatusBadRequest)
		return
	}
	if item.Quantity < 0 {
		writeErrorJSON(w, "quantity cannot be negative", http.StatusBadRequest)
		return
	}

	if err := h.service.AddInventoryItem(item); err != nil {
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"status": "created"})
}

func (h *InventoryHandler) getInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.service.GetInventoryItem(id)
	if err != nil {
		writeErrorJSON(w, "inventory item not found", http.StatusNotFound)
		return
	}
	writeJSON(w, item)
}

func (h *InventoryHandler) updateInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErrorJSON(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	item.IngredientID = id

	if err := h.service.UpdateInventoryItem(item); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorJSON(w, "inventory item not found", http.StatusNotFound)
		} else {
			writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	writeJSON(w, map[string]string{"status": "updated"})
}

func (h *InventoryHandler) deleteInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteInventoryItem(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorJSON(w, "inventory item not found", http.StatusNotFound)
		} else {
			writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]string{"status": "deleted"})
}
