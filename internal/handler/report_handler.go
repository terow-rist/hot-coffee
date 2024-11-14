package handler

import (
	"hot-coffee/internal/service"
	"log/slog"
	"net/http"
)

type ReportsHandler struct {
	reportsService *service.ReportsService
}

func NewReportsHandler(reportsService *service.ReportsService) *ReportsHandler {
	return &ReportsHandler{reportsService: reportsService}
}

func (h *ReportsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	slog.Info("Received report request", slog.String("method", r.Method), slog.String("path", path))

	switch r.Method {
	case http.MethodGet:
		if path == "/reports/total-sales" {
			h.GetTotalSales(w, r)
		} else if path == "/reports/popular-items" {
			h.GetPopularItems(w, r)
		} else {
			respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetTotalSales handles the /reports/total-sales endpoint
func (h *ReportsHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	totalSales, err := h.reportsService.GetTotalSales()
	if err != nil {
		slog.Error("Failed to get total sales", slog.String("error", err.Error()))
		respondWithError(w, "Failed to calculate total sales", http.StatusInternalServerError)
		return
	}

	slog.Info("Total sales calculated successfully", slog.Float64("total_sales", totalSales))
	respondWithJSON(w, map[string]float64{"total_sales": totalSales}, http.StatusOK)
}

// GetPopularItems handles the /reports/popular-items endpoint
func (h *ReportsHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.reportsService.GetPopularItems()
	if err != nil {
		slog.Error("Failed to get popular items", slog.String("error", err.Error()))
		respondWithError(w, "Failed to get popular items", http.StatusInternalServerError)
		return
	}

	slog.Info("Popular items fetched successfully")
	respondWithJSON(w, popularItems, http.StatusOK)
}
