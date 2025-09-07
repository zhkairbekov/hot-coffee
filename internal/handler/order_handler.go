package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllOrders(w, r)
	case http.MethodPost:
		h.createOrder(w, r)
	default:
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) HandleOrderByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		writeErrorJSON(w, "order ID is required", http.StatusBadRequest)
		return
	}

	orderID := parts[0]

	// Проверяем, есть ли дополнительные части пути
	if len(parts) > 1 && parts[1] == "close" {
		if r.Method == http.MethodPost {
			h.closeOrder(w, r, orderID)
		} else {
			writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Обычные CRUD операции
	switch r.Method {
	case http.MethodGet:
		h.getOrder(w, r, orderID)
	case http.MethodPut:
		h.updateOrder(w, r, orderID)
	case http.MethodDelete:
		h.deleteOrder(w, r, orderID)
	default:
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		writeErrorJSON(w, "failed to retrieve orders", http.StatusInternalServerError)
		return
	}
	writeJSON(w, orders)
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeErrorJSON(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация
	if order.CustomerName == "" {
		writeErrorJSON(w, "customer name is required", http.StatusBadRequest)
		return
	}

	if len(order.Items) == 0 {
		writeErrorJSON(w, "order must contain at least one item", http.StatusBadRequest)
		return
	}

	for i, item := range order.Items {
		if item.ProductID == "" {
			writeErrorJSON(w, "product ID is required for all items", http.StatusBadRequest)
			return
		}
		if item.Quantity <= 0 {
			writeErrorJSON(w, "quantity must be positive", http.StatusBadRequest)
			return
		}
		order.Items[i] = item
	}

	if err := h.service.CreateOrder(order); err != nil {
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"status": "created", "order_id": order.ID})
}

func (h *OrderHandler) getOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	order, err := h.service.GetOrder(orderID)
	if err != nil {
		writeErrorJSON(w, "order not found", http.StatusNotFound)
		return
	}
	writeJSON(w, order)
}

func (h *OrderHandler) updateOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeErrorJSON(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	order.ID = orderID

	if err := h.service.UpdateOrder(order); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorJSON(w, "order not found", http.StatusNotFound)
		} else {
			writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	writeJSON(w, map[string]string{"status": "updated"})
}

func (h *OrderHandler) deleteOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	if err := h.service.DeleteOrder(orderID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorJSON(w, "order not found", http.StatusNotFound)
		} else {
			writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]string{"status": "deleted"})
}

func (h *OrderHandler) closeOrder(w http.ResponseWriter, r *http.Request, orderID string) {
	if err := h.service.CloseOrder(orderID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorJSON(w, "order not found", http.StatusNotFound)
		} else {
			writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	writeJSON(w, map[string]string{"status": "closed"})
}
