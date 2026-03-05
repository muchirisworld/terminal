package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusArchived ProductStatus = "archived"
)

type Product struct {
	ID             uuid.UUID     `json:"id"`
	OrganizationID string        `json:"organization_id"`
	Name           string        `json:"name"`
	Description    *string       `json:"description"`
	BaseUnit       string        `json:"base_unit"`
	Status         ProductStatus `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type ProductVariant struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProductID      uuid.UUID `json:"product_id"`
	SKU            string    `json:"sku"`
	Barcode        *string   `json:"barcode"`
	Price          float64   `json:"price"` // Using float64 for simplicity in JSON, numeric in DB
	Cost           *float64  `json:"cost"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProductImage struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProductID      uuid.UUID `json:"product_id"`
	ImageKey       string    `json:"image_key"`
	Position       int       `json:"position"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	BaseUnit    string  `json:"base_unit"`
}

type UpdateProductRequest struct {
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	Status      *ProductStatus `json:"status"`
}

type CreateVariantRequest struct {
	SKU      string   `json:"sku"`
	Barcode  *string  `json:"barcode"`
	Price    float64  `json:"price"`
	Cost     *float64 `json:"cost"`
	IsActive bool     `json:"is_active"`
}

type UpdateVariantRequest struct {
	SKU      *string  `json:"sku"`
	Barcode  *string  `json:"barcode"`
	Price    *float64 `json:"price"`
	Cost     *float64 `json:"cost"`
	IsActive *bool    `json:"is_active"`
}
