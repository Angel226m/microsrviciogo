package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloudmart/inventory-service/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	svc port.InventoryService
}

func NewHandler(svc port.InventoryService) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/stock/{productID}", h.GetStock)
	r.Get("/stock", h.ListStock)
	r.Post("/stock/{productID}/reserve", h.Reserve)
	r.Post("/stock/{productID}/release", h.Release)
	r.Post("/stock/{productID}/restock", h.Restock)
	r.Put("/stock/{productID}", h.UpdateStock)
	r.Get("/movements/{productID}", h.GetMovements)
}

func (h *Handler) GetStock(w http.ResponseWriter, r *http.Request) {
	pid, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	stock, err := h.svc.GetStock(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stock)
}

func (h *Handler) ListStock(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	stocks, total, err := h.svc.ListStock(r.Context(), page, limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  stocks,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) Reserve(w http.ResponseWriter, r *http.Request) {
	pid, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	var body struct {
		Quantity    int    `json:"quantity"`
		ReferenceID string `json:"reference_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid payload")
		return
	}
	refID, err := uuid.Parse(body.ReferenceID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid reference_id")
		return
	}
	if err := h.svc.Reserve(r.Context(), pid, body.Quantity, refID); err != nil {
		writeErr(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reserved"})
}

func (h *Handler) Release(w http.ResponseWriter, r *http.Request) {
	pid, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	var body struct {
		Quantity    int    `json:"quantity"`
		ReferenceID string `json:"reference_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid payload")
		return
	}
	refID, err := uuid.Parse(body.ReferenceID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid reference_id")
		return
	}
	if err := h.svc.ReleaseReservation(r.Context(), pid, body.Quantity, refID); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "released"})
}

func (h *Handler) Restock(w http.ResponseWriter, r *http.Request) {
	pid, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	var body struct {
		Quantity int    `json:"quantity"`
		Notes    string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := h.svc.Restock(r.Context(), pid, body.Quantity); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "restocked"})
}

func (h *Handler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	pid, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	var req port.UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid payload")
		return
	}
	updated, err := h.svc.UpdateStock(r.Context(), pid, req)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) GetMovements(w http.ResponseWriter, r *http.Request) {
	pid, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	movements, err := h.svc.GetMovements(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, movements)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
