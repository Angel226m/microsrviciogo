package http

import (
	"encoding/json"
	"net/http"

	"github.com/cloudmart/payment-service/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct{ service port.PaymentService }

func NewPaymentHandler(svc port.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: svc}
}

func (h *PaymentHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/payments", h.Process)
	r.Get("/api/v1/payments/{id}", h.GetByID)
}

func (h *PaymentHandler) Process(w http.ResponseWriter, r *http.Request) {
	var req port.ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	req.UserID, _ = uuid.Parse(r.Header.Get("X-User-ID"))

	tx, err := h.service.ProcessPayment(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusPaymentRequired, map[string]interface{}{"error": err.Error(), "transaction": tx})
		return
	}
	writeJSON(w, http.StatusCreated, tx)
}

func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	tx, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "transaction not found"})
		return
	}
	writeJSON(w, http.StatusOK, tx)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
