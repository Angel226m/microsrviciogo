// ═══════════════════════════════════════════════════════════════
// Handlers HTTP – Adaptador primario de Usuario (capa de interfaz)
// Traduce peticiones HTTP en llamadas al servicio de aplicación
// ═══════════════════════════════════════════════════════════════
package http

import (
	"encoding/json"
	"net/http"

	"github.com/cloudmart/user-service/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(service port.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes monta las rutas relacionadas con usuarios en el enrutador.
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/auth/register", h.Register)
	r.Post("/api/v1/auth/login", h.Login)
	r.Get("/api/v1/users/me", h.GetProfile)
	r.Put("/api/v1/users/me", h.UpdateProfile)
	r.Get("/api/v1/users/me/addresses", h.ListAddresses)
	r.Post("/api/v1/users/me/addresses", h.AddAddress)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req port.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already registered" {
			status = http.StatusConflict
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"user":    user,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	tokens, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing user id"})
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	id, _ := uuid.Parse(userID)

	var req port.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	id, _ := uuid.Parse(userID)

	addresses, err := h.service.ListAddresses(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"addresses": addresses,
		"count":     len(addresses),
	})
}

func (h *UserHandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	id, _ := uuid.Parse(userID)

	var req port.AddAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	addr, err := h.service.AddAddress(r.Context(), id, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, addr)
}

// ── Helpers ───────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
