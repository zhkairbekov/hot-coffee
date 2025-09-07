// internal/handler/menu_handler.go
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type MenuHandler struct {
	menuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Warn("Invalid JSON in create menu item request", "error", err)
		writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validateMenuItem(&item); err != nil {
		slog.Warn("Menu item validation failed", "error", err)
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.menuService.CreateMenuItem(&item); err != nil {
		slog.Error("Failed to create menu item", "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.menuService.GetAllMenuItems()
	if err != nil {
		slog.Error("Failed to get all menu items", "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Menu item ID is required", http.StatusBadRequest)
		return
	}

	item, err := h.menuService.GetMenuItemByID(id)
	if err != nil {
		slog.Error("Failed to get menu item", "itemID", id, "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if item == nil {
		writeErrorResponse(w, "Menu item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Menu item ID is required", http.StatusBadRequest)
		return
	}

	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Warn("Invalid JSON in update menu item request", "error", err)
		writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	item.ID = id
	if err := validateMenuItem(&item); err != nil {
		slog.Warn("Menu item validation failed", "error", err)
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.menuService.UpdateMenuItem(&item); err != nil {
		slog.Error("Failed to update menu item", "itemID", id, "error", err)
		if err.Error() == "menu item not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
		} else {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Menu item ID is required", http.StatusBadRequest)
		return
	}

	if err := h.menuService.DeleteMenuItem(id); err != nil {
		slog.Error("Failed to delete menu item", "itemID", id, "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
