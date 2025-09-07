// internal/handler/reports_handler.go
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/service"
)

type ReportsHandler struct {
	reportsService service.ReportsService
}

func NewReportsHandler(reportsService service.ReportsService) *ReportsHandler {
	return &ReportsHandler{
		reportsService: reportsService,
	}
}

func (h *ReportsHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	totalSales, err := h.reportsService.GetTotalSales()
	if err != nil {
		slog.Error("Failed to get total sales", "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalSales)
}

func (h *ReportsHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.reportsService.GetPopularItems()
	if err != nil {
		slog.Error("Failed to get popular items", "error", err)
		writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(popularItems)
}
