package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloudmart/notification-service/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct{ svc port.NotificationService }

func NewHandler(svc port.NotificationService) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/notifications", h.Send)
	r.Get("/notifications/{id}", h.GetByID)
	r.Get("/notifications/user/{userID}", h.ListByUser)
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	var req port.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid payload")
		return
	}
	n, err := h.svc.Send(r.Context(), req)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, n)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid ID")
		return
	}
	n, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) ListByUser(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	notifs, total, err := h.svc.ListByUser(r.Context(), uid, page, limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": notifs, "total": total, "page": page, "limit": limit,
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
