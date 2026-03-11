package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	ierrors "github.com/muchirisworld/terminal/internal/ierrors"
	"github.com/muchirisworld/terminal/internal/models"
	"github.com/muchirisworld/terminal/internal/store"
)

type InventoryService struct {
	store *store.Store
}

func NewInventoryService(s *store.Store) *InventoryService {
	return &InventoryService{store: s}
}

func (s *InventoryService) UpsertConversion(ctx context.Context, orgID string, productID uuid.UUID, req *models.UpsertConversionRequest) (*models.UnitConversion, error) {
	return s.store.UpsertConversion(ctx, orgID, productID, req)
}

func (s *InventoryService) ListConversionsByProduct(ctx context.Context, orgID string, productID uuid.UUID) ([]*models.UnitConversion, error) {
	return s.store.ListConversionsByProduct(ctx, orgID, productID)
}

func (s *InventoryService) DeleteConversion(ctx context.Context, orgID string, conversionID uuid.UUID) error {
	return s.store.DeleteConversion(ctx, orgID, conversionID)
}

func (s *InventoryService) CreateInventoryReceipt(ctx context.Context, orgID string, variantID uuid.UUID, req *models.ReceiptRequest) (*models.InventoryEvent, error) {
	variant, err := s.store.GetVariant(ctx, orgID, variantID)
	if err != nil {
		return nil, err
	}
	if variant == nil {
		return nil, fmt.Errorf("variant not found")
	}

	product, err := s.store.GetProduct(ctx, orgID, variant.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	totalBaseQty := req.Quantity
	if req.Unit != product.BaseUnit {
		conv, err := s.store.GetConversion(ctx, orgID, product.ID, req.Unit, product.BaseUnit)
		if err != nil {
			return nil, err
		}
		if conv == nil {
			return nil, fmt.Errorf("conversion from %s to %s not found for product %s", req.Unit, product.BaseUnit, product.Name)
		}

		calculated := req.Quantity * conv.Factor
		// Apply precision rounding if conversion specifies it
		roundingFactor := math.Pow(10, float64(conv.Precision))
		totalBaseQty = math.Round(calculated*roundingFactor) / roundingFactor
	}

	sourceType := models.InventorySourceTypeReceipt
	return s.store.CreateInventoryEvent(ctx, orgID, variantID, models.InventoryEventTypeReceipt, totalBaseQty, &sourceType, req.SourceID, req.Note)
}

func (s *InventoryService) CreateInventoryAdjustment(ctx context.Context, orgID string, variantID uuid.UUID, req *models.AdjustmentRequest) (*models.InventoryEvent, error) {
	sourceType := models.InventorySourceTypeManual
	return s.store.CreateInventoryEvent(ctx, orgID, variantID, models.InventoryEventTypeAdjustment, req.QuantityChange, &sourceType, nil, req.Note)
}

func (s *InventoryService) ReserveInventory(ctx context.Context, orgID string, variantID uuid.UUID, req *models.ReservationRequest) (*models.InventoryReservation, error) {
	var reservation *models.InventoryReservation
	err := s.store.ExecTx(ctx, func(txStore *store.Store) error {
		// 1. Lock variant
		v, err := txStore.GetVariantWithLock(ctx, orgID, variantID)
		if err != nil {
			return err
		}
		if v == nil {
			return fmt.Errorf("variant not found")
		}

		// 2. Check stock
		stock, err := txStore.GetVariantStock(ctx, orgID, variantID)
		if err != nil {
			return err
		}

		if stock.AvailableStock < req.Quantity {
			return &ierrors.InsufficientStockError{Message: fmt.Sprintf("insufficient stock: available %f, requested %f", stock.AvailableStock, req.Quantity)}
		}

		// 3. Create reservation
		res, err := txStore.CreateInventoryReservation(ctx, orgID, variantID, req.Quantity, req.ExpiresAt, req.OrderID)
		if err != nil {
			return err
		}
		reservation = res
		return nil
	})

	return reservation, err
}

func (s *InventoryService) ReleaseReservation(ctx context.Context, orgID string, reservationID uuid.UUID) error {
	return s.store.UpdateReservationStatus(ctx, orgID, reservationID, models.ReservationStatusReleased)
}

func (s *InventoryService) FulfillReservation(ctx context.Context, orgID string, reservationID uuid.UUID) (*models.InventoryEvent, error) {
	var event *models.InventoryEvent
	err := s.store.ExecTx(ctx, func(txStore *store.Store) error {
		// 1. Get reservation
		res, err := txStore.GetReservation(ctx, orgID, reservationID)
		if err != nil {
			return err
		}
		if res == nil {
			return fmt.Errorf("reservation not found")
		}
		if res.Status != models.ReservationStatusActive {
			return fmt.Errorf("reservation is not active")
		}

		// 2. Create fulfillment event (negative quantity change)
		sourceType := models.InventorySourceTypeOrder
		msg := "Order fulfillment"
		e, err := txStore.CreateInventoryEvent(ctx, orgID, res.ProductVariantID, models.InventoryEventTypeFulfillment, -res.Quantity, &sourceType, res.OrderID, &msg)
		if err != nil {
			return err
		}
		event = e

		// 3. Release reservation
		return txStore.UpdateReservationStatus(ctx, orgID, reservationID, models.ReservationStatusReleased)
	})

	return event, err
}

func (s *InventoryService) GetVariantStock(ctx context.Context, orgID string, variantID uuid.UUID) (*models.VariantStock, error) {
	return s.store.GetVariantStock(ctx, orgID, variantID)
}

func (s *InventoryService) ExpireReservations(ctx context.Context) (int64, error) {
	return s.store.ExpireReservations(ctx)
}
