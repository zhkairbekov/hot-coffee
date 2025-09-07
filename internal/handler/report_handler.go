package handler

import (
	"net/http"

	"hot-coffee/internal/service"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleTotalSales(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	totalSales, err := h.service.GetTotalSales()
	if err != nil {
		writeErrorJSON(w, "failed to calculate total sales", http.StatusInternalServerError)
		return
	}

	response := map[string]float64{
		"total_sales": totalSales,
	}
	writeJSON(w, response)
}

func (h *ReportHandler) HandlePopularItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorJSON(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	popularItems, err := h.service.GetPopularItems()
	if err != nil {
		writeErrorJSON(w, "failed to get popular items", http.StatusInternalServerError)
		return
	}

	writeJSON(w, popularItems)
}
