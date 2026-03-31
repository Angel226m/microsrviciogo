// ═══════════════════════════════════════════════════════════════
// Handlers HTTP – Adaptador primario de Pedidos
// Traduce peticiones HTTP en llamadas al servicio de aplicación
// ═══════════════════════════════════════════════════════════════
package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloudmart/order-service/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	service port.OrderService
}

func NewOrderHandler(service port.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/orders", h.Create)
	r.Get("/api/v1/orders", h.List)
	r.Get("/api/v1/orders/{id}", h.GetByID)
	r.Put("/api/v1/orders/{id}/cancel", h.Cancel)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(r.Header.Get("X-User-ID"))

	var req port.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	req.UserID = userID

	order, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(r.Header.Get("X-User-ID"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	orders, total, err := h.service.ListByUser(r.Context(), userID, page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"orders": orders, "total": total})
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	order, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "order not found"})
		return
	}
	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	userID, _ := uuid.Parse(r.Header.Get("X-User-ID"))

	if err := h.service.Cancel(r.Context(), id, userID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "order cancelled"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
