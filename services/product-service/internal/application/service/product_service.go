// ═══════════════════════════════════════════════════════════════
// Application Layer – Product Service (use cases)
// ═══════════════════════════════════════════════════════════════
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudmart/product-service/internal/domain/model"
	"github.com/cloudmart/product-service/internal/domain/port"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidRating   = errors.New("rating must be between 1 and 5")
)

type productService struct {
	productRepo  port.ProductRepository
	categoryRepo port.CategoryRepository
	reviewRepo   port.ReviewRepository
	cache        port.CacheRepository
	events       port.EventPublisher
}

func NewProductService(
	productRepo port.ProductRepository,
	categoryRepo port.CategoryRepository,
	reviewRepo port.ReviewRepository,
	cache port.CacheRepository,
	events port.EventPublisher,
) port.ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		reviewRepo:   reviewRepo,
		cache:        cache,
		events:       events,
	}
}

func (s *productService) List(ctx context.Context, filter model.ProductFilter) (*model.ProductListResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 50 {
		filter.Limit = 20
	}

	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	totalPages := (total + filter.Limit - 1) / filter.Limit

	return &model.ProductListResult{
		Products:   products,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) GetBySlug(ctx context.Context, slug string) (*model.Product, error) {
	product, err := s.productRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *productService) Create(ctx context.Context, req port.CreateProductRequest) (*model.Product, error) {
	slug := generateSlug(req.Name)

	product := &model.Product{
		ID:               uuid.New(),
		SKU:              req.SKU,
		Name:             req.Name,
		Slug:             slug,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		CompareAtPrice:   req.CompareAtPrice,
		CategoryID:       req.CategoryID,
		Brand:            req.Brand,
		Tags:             req.Tags,
		Images:           req.Images,
		Attributes:       req.Attributes,
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if len(product.Images) > 0 {
		product.ThumbnailURL = product.Images[0]
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	_ = s.events.Publish(ctx, "product.created", map[string]interface{}{
		"product_id": product.ID,
		"sku":        product.SKU,
		"name":       product.Name,
		"price":      product.Price,
	})

	return product, nil
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, req port.UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if req.Name != "" {
		product.Name = req.Name
		product.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CompareAtPrice != nil {
		product.CompareAtPrice = req.CompareAtPrice
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if len(req.Images) > 0 {
		product.Images = req.Images
		product.ThumbnailURL = req.Images[0]
	}
	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf("product:%s", id.String()))

	return product, nil
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *productService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.categoryRepo.List(ctx)
}

func (s *productService) GetReviews(ctx context.Context, productID uuid.UUID) ([]model.Review, error) {
	return s.reviewRepo.FindByProductID(ctx, productID)
}

func (s *productService) AddReview(ctx context.Context, req port.AddReviewRequest) (*model.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, ErrInvalidRating
	}

	review := &model.Review{
		ID:        uuid.New(),
		ProductID: req.ProductID,
		UserID:    req.UserID,
		Rating:    req.Rating,
		Title:     req.Title,
		Body:      req.Body,
		CreatedAt: time.Now(),
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("create review: %w", err)
	}

	return review, nil
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	return slug
}
