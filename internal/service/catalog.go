package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchirisworld/terminal/internal/models"
	"github.com/muchirisworld/terminal/internal/store"
)

type CatalogService struct {
	store *store.Store
}

func NewCatalogService(s *store.Store) *CatalogService {
	return &CatalogService{store: s}
}

func (s *CatalogService) CreateProduct(ctx context.Context, orgID string, req *models.CreateProductRequest) (*models.Product, error) {
	return s.store.CreateProduct(ctx, orgID, req)
}

func (s *CatalogService) UpdateProduct(ctx context.Context, orgID string, productID uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error) {
	return s.store.UpdateProduct(ctx, orgID, productID, req)
}

func (s *CatalogService) GetProduct(ctx context.Context, orgID string, productID uuid.UUID) (*models.Product, error) {
	return s.store.GetProduct(ctx, orgID, productID)
}

func (s *CatalogService) ListProducts(ctx context.Context, orgID string, limit, offset int) ([]*models.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.store.ListProducts(ctx, orgID, limit, offset)
}

func (s *CatalogService) DeleteProduct(ctx context.Context, orgID string, productID uuid.UUID) error {
	return s.store.DeleteProduct(ctx, orgID, productID)
}

func (s *CatalogService) CreateVariant(ctx context.Context, orgID string, productID uuid.UUID, req *models.CreateVariantRequest) (*models.ProductVariant, error) {
	return s.store.CreateVariant(ctx, orgID, productID, req)
}

func (s *CatalogService) UpdateVariant(ctx context.Context, orgID string, variantID uuid.UUID, req *models.UpdateVariantRequest) (*models.ProductVariant, error) {
	return s.store.UpdateVariant(ctx, orgID, variantID, req)
}

func (s *CatalogService) GetVariant(ctx context.Context, orgID string, variantID uuid.UUID) (*models.ProductVariant, error) {
	return s.store.GetVariant(ctx, orgID, variantID)
}

func (s *CatalogService) ListVariantsByProduct(ctx context.Context, orgID string, productID uuid.UUID) ([]*models.ProductVariant, error) {
	return s.store.ListVariantsByProduct(ctx, orgID, productID)
}

func (s *CatalogService) ArchiveProduct(ctx context.Context, orgID string, productID uuid.UUID) (*models.Product, error) {
	status := models.ProductStatusArchived
	return s.store.UpdateProduct(ctx, orgID, productID, &models.UpdateProductRequest{Status: &status})
}

func (s *CatalogService) DeleteVariant(ctx context.Context, orgID string, variantID uuid.UUID) error {
	return s.store.DeleteVariant(ctx, orgID, variantID)
}
