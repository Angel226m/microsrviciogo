// ═══════════════════════════════════════════════════════════════
// PostgreSQL Repository – Product driven adapter
// ═══════════════════════════════════════════════════════════════
package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudmart/product-service/internal/domain/model"
	"github.com/cloudmart/product-service/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewProductPostgresRepo(pool *pgxpool.Pool) port.ProductRepository {
	return &productPostgresRepo{pool: pool}
}

func (r *productPostgresRepo) List(ctx context.Context, filter model.ProductFilter) ([]model.Product, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, "is_active = true")

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}

	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIdx))
		args = append(args, *filter.MinPrice)
		argIdx++
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIdx))
		args = append(args, *filter.MaxPrice)
		argIdx++
	}

	if filter.Brand != "" {
		conditions = append(conditions, fmt.Sprintf("brand = $%d", argIdx))
		args = append(args, filter.Brand)
		argIdx++
	}

	if filter.IsFeatured != nil {
		conditions = append(conditions, fmt.Sprintf("is_featured = $%d", argIdx))
		args = append(args, *filter.IsFeatured)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products.items %s", where)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Sort
	orderBy := "created_at DESC"
	if filter.SortBy != "" {
		dir := "ASC"
		if filter.SortDir == "desc" {
			dir = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filter.SortBy, dir)
	}

	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT id, sku, name, slug, description, short_description, price, compare_at_price,
		       category_id, brand, tags, images, thumbnail_url, is_active, is_featured,
		       rating_avg, rating_count, created_at, updated_at
		FROM products.items %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		where, orderBy, argIdx, argIdx+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Slug, &p.Description, &p.ShortDescription,
			&p.Price, &p.CompareAtPrice, &p.CategoryID, &p.Brand, &p.Tags, &p.Images,
			&p.ThumbnailURL, &p.IsActive, &p.IsFeatured, &p.RatingAvg, &p.RatingCount,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	return products, total, nil
}

func (r *productPostgresRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	query := `
		SELECT id, sku, name, slug, description, short_description, price, compare_at_price,
		       category_id, brand, tags, images, thumbnail_url, is_active, is_featured,
		       rating_avg, rating_count, created_at, updated_at
		FROM products.items WHERE id = $1`

	p := &model.Product{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Slug, &p.Description, &p.ShortDescription,
		&p.Price, &p.CompareAtPrice, &p.CategoryID, &p.Brand, &p.Tags, &p.Images,
		&p.ThumbnailURL, &p.IsActive, &p.IsFeatured, &p.RatingAvg, &p.RatingCount,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *productPostgresRepo) FindBySlug(ctx context.Context, slug string) (*model.Product, error) {
	query := `
		SELECT id, sku, name, slug, description, short_description, price, compare_at_price,
		       category_id, brand, tags, images, thumbnail_url, is_active, is_featured,
		       rating_avg, rating_count, created_at, updated_at
		FROM products.items WHERE slug = $1 AND is_active = true`

	p := &model.Product{}
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Slug, &p.Description, &p.ShortDescription,
		&p.Price, &p.CompareAtPrice, &p.CategoryID, &p.Brand, &p.Tags, &p.Images,
		&p.ThumbnailURL, &p.IsActive, &p.IsFeatured, &p.RatingAvg, &p.RatingCount,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *productPostgresRepo) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products.items (id, sku, name, slug, description, short_description, price,
		    compare_at_price, category_id, brand, tags, images, thumbnail_url, is_active, is_featured,
		    created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`

	_, err := r.pool.Exec(ctx, query,
		product.ID, product.SKU, product.Name, product.Slug, product.Description,
		product.ShortDescription, product.Price, product.CompareAtPrice, product.CategoryID,
		product.Brand, product.Tags, product.Images, product.ThumbnailURL,
		product.IsActive, product.IsFeatured, product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *productPostgresRepo) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products.items SET name=$2, slug=$3, description=$4, short_description=$5,
		    price=$6, compare_at_price=$7, brand=$8, tags=$9, images=$10, thumbnail_url=$11,
		    is_active=$12, is_featured=$13, updated_at=$14 WHERE id=$1`

	_, err := r.pool.Exec(ctx, query,
		product.ID, product.Name, product.Slug, product.Description, product.ShortDescription,
		product.Price, product.CompareAtPrice, product.Brand, product.Tags, product.Images,
		product.ThumbnailURL, product.IsActive, product.IsFeatured, product.UpdatedAt,
	)
	return err
}

func (r *productPostgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE products.items SET is_active = false WHERE id = $1", id)
	return err
}

// ── Category Repository ───────────────────────────────────────────────

type categoryPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryPostgresRepo(pool *pgxpool.Pool) port.CategoryRepository {
	return &categoryPostgresRepo{pool: pool}
}

func (r *categoryPostgresRepo) List(ctx context.Context) ([]model.Category, error) {
	query := `SELECT id, name, slug, description, image_url, parent_id, sort_order, is_active, created_at
		FROM products.categories WHERE is_active = true ORDER BY sort_order`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ImageURL,
			&c.ParentID, &c.SortOrder, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryPostgresRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	query := `SELECT id, name, slug, description, image_url, parent_id, sort_order, is_active, created_at
		FROM products.categories WHERE id = $1`

	c := &model.Category{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Slug, &c.Description,
		&c.ImageURL, &c.ParentID, &c.SortOrder, &c.IsActive, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ── Review Repository ─────────────────────────────────────────────────

type reviewPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewReviewPostgresRepo(pool *pgxpool.Pool) port.ReviewRepository {
	return &reviewPostgresRepo{pool: pool}
}

func (r *reviewPostgresRepo) FindByProductID(ctx context.Context, productID uuid.UUID) ([]model.Review, error) {
	query := `SELECT id, product_id, user_id, rating, title, body, is_verified, created_at
		FROM products.reviews WHERE product_id = $1 ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var rv model.Review
		if err := rows.Scan(&rv.ID, &rv.ProductID, &rv.UserID, &rv.Rating,
			&rv.Title, &rv.Body, &rv.IsVerified, &rv.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	return reviews, nil
}

func (r *reviewPostgresRepo) Create(ctx context.Context, review *model.Review) error {
	query := `INSERT INTO products.reviews (id, product_id, user_id, rating, title, body, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, query, review.ID, review.ProductID, review.UserID,
		review.Rating, review.Title, review.Body, review.CreatedAt)
	return err
}

func (r *reviewPostgresRepo) GetAverageRating(ctx context.Context, productID uuid.UUID) (float64, int, error) {
	var avg float64
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COALESCE(AVG(rating),0), COUNT(*) FROM products.reviews WHERE product_id = $1", productID).Scan(&avg, &count)
	return avg, count, err
}
