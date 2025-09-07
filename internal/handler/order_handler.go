// internal/handler/order_handler.go
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Warn("Invalid JSON in create order request", "error", err)
		writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validateOrder(&order); err != nil {
		slog.Warn("Order validation failed", "error", err)
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.orderService.CreateOrder(&order); err != nil {
		slog.Error("Failed to create order", "error", err)
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		slog.Error("Failed to get all orders", "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		slog.Error("Failed to get order", "orderID", id, "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if order == nil {
		writeErrorResponse(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Warn("Invalid JSON in update order request", "error", err)
		writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	order.ID = id
	if err := validateOrder(&order); err != nil {
		slog.Warn("Order validation failed", "error", err)
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.orderService.UpdateOrder(&order); err != nil {
		slog.Error("Failed to update order", "orderID", id, "error", err)
		if err.Error() == "order not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
		} else {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	if err := h.orderService.DeleteOrder(id); err != nil {
		slog.Error("Failed to delete order", "orderID", id, "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrorResponse(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	if err := h.orderService.CloseOrder(id); err != nil {
		slog.Error("Failed to close order", "orderID", id, "error", err)
		if err.Error() == "order not found" {
			writeErrorResponse(w, err.Error(), http.StatusNotFound)
		} else {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
