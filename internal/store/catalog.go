package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/muchirisworld/terminal/internal/models"
)

func (s *Store) CreateProduct(ctx context.Context, orgID string, req *models.CreateProductRequest) (*models.Product, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`INSERT INTO products (organization_id, name, description, base_unit)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, organization_id, name, description, base_unit, status, created_at, updated_at`,
		orgID, req.Name, req.Description, req.BaseUnit,
	)

	var p models.Product
	err := row.Scan(&p.ID, &p.OrganizationID, &p.Name, &p.Description, &p.BaseUnit, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (s *Store) UpdateProduct(ctx context.Context, orgID string, productID uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error) {
	query := `UPDATE products SET `
	args := []interface{}{}
	count := 1

	if req.Name != nil {
		query += fmt.Sprintf("name = $%d, ", count)
		args = append(args, *req.Name)
		count++
	}
	if req.Description != nil {
		query += fmt.Sprintf("description = $%d, ", count)
		args = append(args, *req.Description)
		count++
	}
	if req.Status != nil {
		query += fmt.Sprintf("status = $%d, ", count)
		args = append(args, *req.Status)
		count++
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND organization_id = $%d RETURNING id, organization_id, name, description, base_unit, status, created_at, updated_at", count, count+1)
	args = append(args, productID, orgID)

	var p models.Product
	err := s.dbtx.QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.OrganizationID, &p.Name, &p.Description, &p.BaseUnit, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (s *Store) GetProduct(ctx context.Context, orgID string, productID uuid.UUID) (*models.Product, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`SELECT id, organization_id, name, description, base_unit, status, created_at, updated_at
		 FROM products WHERE id = $1 AND organization_id = $2`,
		productID, orgID,
	)

	var p models.Product
	err := row.Scan(&p.ID, &p.OrganizationID, &p.Name, &p.Description, &p.BaseUnit, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (s *Store) ListProducts(ctx context.Context, orgID string, limit, offset int) ([]*models.Product, error) {
	rows, err := s.dbtx.QueryContext(ctx,
		`SELECT id, organization_id, name, description, base_unit, status, created_at, updated_at
		 FROM products WHERE organization_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.OrganizationID, &p.Name, &p.Description, &p.BaseUnit, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (s *Store) DeleteProduct(ctx context.Context, orgID string, productID uuid.UUID) error {
	_, err := s.dbtx.ExecContext(ctx,
		`DELETE FROM products WHERE id = $1 AND organization_id = $2`,
		productID, orgID,
	)
	return err
}

func (s *Store) CreateVariant(ctx context.Context, orgID string, productID uuid.UUID, req *models.CreateVariantRequest) (*models.ProductVariant, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`INSERT INTO product_variants (organization_id, product_id, sku, barcode, price, cost, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, organization_id, product_id, sku, barcode, price, cost, is_active, created_at, updated_at`,
		orgID, productID, req.SKU, req.Barcode, req.Price, req.Cost, req.IsActive,
	)

	var v models.ProductVariant
	err := row.Scan(&v.ID, &v.OrganizationID, &v.ProductID, &v.SKU, &v.Barcode, &v.Price, &v.Cost, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (s *Store) UpdateVariant(ctx context.Context, orgID string, variantID uuid.UUID, req *models.UpdateVariantRequest) (*models.ProductVariant, error) {
	query := `UPDATE product_variants SET `
	args := []interface{}{}
	count := 1

	if req.SKU != nil {
		query += fmt.Sprintf("sku = $%d, ", count)
		args = append(args, *req.SKU)
		count++
	}
	if req.Barcode != nil {
		query += fmt.Sprintf("barcode = $%d, ", count)
		args = append(args, *req.Barcode)
		count++
	}
	if req.Price != nil {
		query += fmt.Sprintf("price = $%d, ", count)
		args = append(args, *req.Price)
		count++
	}
	if req.Cost != nil {
		query += fmt.Sprintf("cost = $%d, ", count)
		args = append(args, *req.Cost)
		count++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf("is_active = $%d, ", count)
		args = append(args, *req.IsActive)
		count++
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND organization_id = $%d RETURNING id, organization_id, product_id, sku, barcode, price, cost, is_active, created_at, updated_at", count, count+1)
	args = append(args, variantID, orgID)

	var v models.ProductVariant
	err := s.dbtx.QueryRowContext(ctx, query, args...).Scan(&v.ID, &v.OrganizationID, &v.ProductID, &v.SKU, &v.Barcode, &v.Price, &v.Cost, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (s *Store) GetVariant(ctx context.Context, orgID string, variantID uuid.UUID) (*models.ProductVariant, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`SELECT id, organization_id, product_id, sku, barcode, price, cost, is_active, created_at, updated_at
		 FROM product_variants WHERE id = $1 AND organization_id = $2`,
		variantID, orgID,
	)

	var v models.ProductVariant
	err := row.Scan(&v.ID, &v.OrganizationID, &v.ProductID, &v.SKU, &v.Barcode, &v.Price, &v.Cost, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &v, err
}

func (s *Store) ListVariantsByProduct(ctx context.Context, orgID string, productID uuid.UUID) ([]*models.ProductVariant, error) {
	rows, err := s.dbtx.QueryContext(ctx,
		`SELECT id, organization_id, product_id, sku, barcode, price, cost, is_active, created_at, updated_at
		 FROM product_variants WHERE organization_id = $1 AND product_id = $2 ORDER BY created_at ASC`,
		orgID, productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := make([]*models.ProductVariant, 0)
	for rows.Next() {
		var v models.ProductVariant
		if err := rows.Scan(&v.ID, &v.OrganizationID, &v.ProductID, &v.SKU, &v.Barcode, &v.Price, &v.Cost, &v.IsActive, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		variants = append(variants, &v)
	}
	return variants, nil
}

func (s *Store) DeleteVariant(ctx context.Context, orgID string, variantID uuid.UUID) error {
	_, err := s.dbtx.ExecContext(ctx,
		`DELETE FROM product_variants WHERE id = $1 AND organization_id = $2`,
		variantID, orgID,
	)
	return err
}
