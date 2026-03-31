// ═══════════════════════════════════════════════════════════════
// Handlers HTTP – Adaptador primario de Producto
// Traduce peticiones HTTP en llamadas al servicio de aplicación
// ═══════════════════════════════════════════════════════════════
package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloudmart/product-service/internal/domain/model"
	"github.com/cloudmart/product-service/internal/domain/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service port.ProductService
}

func NewProductHandler(service port.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/products", h.List)
	r.Get("/api/v1/products/{slug}", h.GetBySlug)
	r.Post("/api/v1/products", h.Create)
	r.Put("/api/v1/products/{id}", h.Update)
	r.Delete("/api/v1/products/{id}", h.Delete)
	r.Get("/api/v1/categories", h.ListCategories)
	r.Get("/api/v1/products/{id}/reviews", h.GetReviews)
	r.Post("/api/v1/products/{id}/reviews", h.AddReview)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := model.ProductFilter{
		Search:  q.Get("search"),
		Brand:   q.Get("brand"),
		SortBy:  q.Get("sort_by"),
		SortDir: q.Get("sort_dir"),
		Page:    1,
		Limit:   20,
	}

	if p, err := strconv.Atoi(q.Get("page")); err == nil {
		filter.Page = p
	}
	if l, err := strconv.Atoi(q.Get("limit")); err == nil {
		filter.Limit = l
	}
	if min, err := strconv.ParseFloat(q.Get("min_price"), 64); err == nil {
		filter.MinPrice = &min
	}
	if max, err := strconv.ParseFloat(q.Get("max_price"), 64); err == nil {
		filter.MaxPrice = &max
	}
	if catID := q.Get("category_id"); catID != "" {
		if id, err := uuid.Parse(catID); err == nil {
			filter.CategoryID = &id
		}
	}
	if feat := q.Get("featured"); feat == "true" {
		t := true
		filter.IsFeatured = &t
	}

	result, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	product, err := h.service.GetBySlug(r.Context(), slug)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req port.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	product, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}

	var req port.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	product, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"categories": categories})
}

func (h *ProductHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}

	reviews, err := h.service.GetReviews(r.Context(), productID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"reviews": reviews, "count": len(reviews)})
}

func (h *ProductHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	productID, _ := uuid.Parse(chi.URLParam(r, "id"))
	userID, _ := uuid.Parse(r.Header.Get("X-User-ID"))

	var req struct {
		Rating int    `json:"rating"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	review, err := h.service.AddReview(r.Context(), port.AddReviewRequest{
		ProductID: productID,
		UserID:    userID,
		Rating:    req.Rating,
		Title:     req.Title,
		Body:      req.Body,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, review)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
