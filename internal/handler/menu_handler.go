package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type MenuHandler struct {
	service *service.MenuService
}

func NewMenuHandler(service *service.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) HandleMenu(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllMenu(w, r)
	case http.MethodPost:
		h.createMenuItem(w, r)
	default:
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MenuHandler) HandleMenuByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/menu/")
	if id == "" {
		writeErrorJSON(w, "product ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getMenuItem(w, r, id)
	case http.MethodPut:
		h.updateMenuItem(w, r, id)
	case http.MethodDelete:
		h.deleteMenuItem(w, r, id)
	default:
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MenuHandler) getAllMenu(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetMenu()
	if err != nil {
		writeErrorJSON(w, "failed to retrieve menu", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}

func (h *MenuHandler) createMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErrorJSON(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация
	if item.ID == "" {
		writeErrorJSON(w, "product ID is required", http.StatusBadRequest)
		return
	}
	if item.Name == "" {
		writeErrorJSON(w, "product name is required", http.StatusBadRequest)
		return
	}
	if item.Price <= 0 {
		writeErrorJSON(w, "price must be positive", http.StatusBadRequest)
		return
	}

	if err := h.service.AddMenuItem(item); err != nil {
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"status": "created"})
}

func (h *MenuHandler) getMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.service.GetMenuItem(id)
	if err != nil {
		writeErrorJSON(w, "menu item not found", http.StatusNotFound)
		return
	}
	writeJSON(w, item)
}

func (h *MenuHandler) updateMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErrorJSON(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	item.ID = id

	if err := h.service.UpdateMenuItem(item); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorJSON(w, "menu item not found", http.StatusNotFound)
		} else {
			writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	writeJSON(w, map[string]string{"status": "updated"})
}

func (h *MenuHandler) deleteMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteMenuItem(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorJSON(w, "menu item not found", http.StatusNotFound)
		} else {
			writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]string{"status": "deleted"})
}
