package models

import (
	"time"

	"github.com/google/uuid"
)

type InventoryEventType string

const (
	InventoryEventTypeReceipt     InventoryEventType = "purchase_receipt"
	InventoryEventTypeFulfillment InventoryEventType = "order_fulfillment"
	InventoryEventTypeAdjustment  InventoryEventType = "adjustment"
	InventoryEventTypeReturn      InventoryEventType = "return"
	InventoryEventTypeCaseBreak   InventoryEventType = "case_break"
)

type InventorySourceType string

const (
	InventorySourceTypeReceipt InventorySourceType = "receipt"
	InventorySourceTypeOrder   InventorySourceType = "order"
	InventorySourceTypeManual  InventorySourceType = "manual"
	InventorySourceTypeSystem  InventorySourceType = "system"
)

type ReservationStatus string

const (
	ReservationStatusActive   ReservationStatus = "active"
	ReservationStatusReleased ReservationStatus = "released"
	ReservationStatusExpired  ReservationStatus = "expired"
)

type UnitConversion struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProductID      uuid.UUID `json:"product_id"`
	UnitFrom       string    `json:"unit_from"`
	UnitTo         string    `json:"unit_to"`
	Factor         float64   `json:"factor"`
	Precision      int       `json:"precision"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type InventoryEvent struct {
	ID               uuid.UUID             `json:"id"`
	OrganizationID   string                `json:"organization_id"`
	ProductVariantID uuid.UUID             `json:"product_variant_id"`
	EventType        InventoryEventType    `json:"event_type"`
	QuantityChange   int64                 `json:"quantity_change"`
	SourceType       *InventorySourceType `json:"source_type"`
	SourceID         *uuid.UUID            `json:"source_id"`
	Note             *string               `json:"note"`
	CreatedAt        time.Time             `json:"created_at"`
}

type InventoryReservation struct {
	ID               uuid.UUID          `json:"id"`
	OrganizationID   string             `json:"organization_id"`
	ProductVariantID uuid.UUID          `json:"product_variant_id"`
	OrderID          *uuid.UUID         `json:"order_id"`
	Quantity         int64              `json:"quantity"`
	Status           ReservationStatus `json:"status"`
	ExpiresAt        *time.Time         `json:"expires_at"`
	CreatedAt        time.Time          `json:"created_at"`
	ReleasedAt       *time.Time         `json:"released_at"`
}

type UpsertConversionRequest struct {
	UnitFrom  string  `json:"unit_from"`
	UnitTo    string  `json:"unit_to"`
	Factor    float64 `json:"factor"`
	Precision int     `json:"precision"`
}

type ReceiptRequest struct {
	Quantity int64   `json:"quantity"`
	Unit     string  `json:"unit"`
	SourceID *uuid.UUID `json:"source_id"`
	Note     *string    `json:"note"`
}

type AdjustmentRequest struct {
	QuantityChange int64   `json:"quantity_change"`
	Note           *string `json:"note"`
}

type ReservationRequest struct {
	Quantity  int64      `json:"quantity"`
	ExpiresAt *time.Time `json:"expires_at"`
	OrderID   *uuid.UUID `json:"order_id"`
}

type VariantStock struct {
	TotalStock     int64 `json:"total_stock"`
	ReservedStock  int64 `json:"reserved_stock"`
	AvailableStock int64 `json:"available_stock"`
}
