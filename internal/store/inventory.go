package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/muchirisworld/terminal/internal/models"
)

func (s *Store) UpsertConversion(ctx context.Context, orgID string, productID uuid.UUID, req *models.UpsertConversionRequest) (*models.UnitConversion, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`INSERT INTO unit_conversions (organization_id, product_id, unit_from, unit_to, factor, precision)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (organization_id, product_id, unit_from, unit_to)
		 DO UPDATE SET factor = EXCLUDED.factor, precision = EXCLUDED.precision, updated_at = now()
		 RETURNING id, organization_id, product_id, unit_from, unit_to, factor, precision, created_at, updated_at`,
		orgID, productID, req.UnitFrom, req.UnitTo, req.Factor, req.Precision,
	)

	var c models.UnitConversion
	err := row.Scan(&c.ID, &c.OrganizationID, &c.ProductID, &c.UnitFrom, &c.UnitTo, &c.Factor, &c.Precision, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (s *Store) ListConversionsByProduct(ctx context.Context, orgID string, productID uuid.UUID) ([]*models.UnitConversion, error) {
	rows, err := s.dbtx.QueryContext(ctx,
		`SELECT id, organization_id, product_id, unit_from, unit_to, factor, precision, created_at, updated_at
		 FROM unit_conversions WHERE organization_id = $1 AND product_id = $2`,
		orgID, productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversions := make([]*models.UnitConversion, 0)
	for rows.Next() {
		var c models.UnitConversion
		if err := rows.Scan(&c.ID, &c.OrganizationID, &c.ProductID, &c.UnitFrom, &c.UnitTo, &c.Factor, &c.Precision, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		conversions = append(conversions, &c)
	}
	return conversions, nil
}

func (s *Store) GetConversion(ctx context.Context, orgID string, productID uuid.UUID, unitFrom, unitTo string) (*models.UnitConversion, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`SELECT id, organization_id, product_id, unit_from, unit_to, factor, precision, created_at, updated_at
		 FROM unit_conversions WHERE organization_id = $1 AND product_id = $2 AND unit_from = $3 AND unit_to = $4`,
		orgID, productID, unitFrom, unitTo,
	)

	var c models.UnitConversion
	err := row.Scan(&c.ID, &c.OrganizationID, &c.ProductID, &c.UnitFrom, &c.UnitTo, &c.Factor, &c.Precision, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *Store) DeleteConversion(ctx context.Context, orgID string, conversionID uuid.UUID) error {
	_, err := s.dbtx.ExecContext(ctx,
		`DELETE FROM unit_conversions WHERE id = $1 AND organization_id = $2`,
		conversionID, orgID,
	)
	return err
}

func (s *Store) CreateInventoryEvent(ctx context.Context, orgID string, variantID uuid.UUID, eventType models.InventoryEventType, qty int64, sourceType *models.InventorySourceType, sourceID *uuid.UUID, note *string) (*models.InventoryEvent, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`INSERT INTO inventory_events (organization_id, product_variant_id, event_type, quantity_change, source_type, source_id, note)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, organization_id, product_variant_id, event_type, quantity_change, source_type, source_id, note, created_at`,
		orgID, variantID, eventType, qty, sourceType, sourceID, note,
	)

	var e models.InventoryEvent
	err := row.Scan(&e.ID, &e.OrganizationID, &e.ProductVariantID, &e.EventType, &e.QuantityChange, &e.SourceType, &e.SourceID, &e.Note, &e.CreatedAt)
	return &e, err
}

func (s *Store) CreateInventoryReservation(ctx context.Context, orgID string, variantID uuid.UUID, qty int64, expiresAt *time.Time, orderID *uuid.UUID) (*models.InventoryReservation, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`INSERT INTO inventory_reservations (organization_id, product_variant_id, quantity, expires_at, order_id, status)
		 VALUES ($1, $2, $3, $4, $5, 'active')
		 RETURNING id, organization_id, product_variant_id, order_id, quantity, status, expires_at, created_at, released_at`,
		orgID, variantID, qty, expiresAt, orderID,
	)

	var r models.InventoryReservation
	err := row.Scan(&r.ID, &r.OrganizationID, &r.ProductVariantID, &r.OrderID, &r.Quantity, &r.Status, &r.ExpiresAt, &r.CreatedAt, &r.ReleasedAt)
	return &r, err
}

func (s *Store) GetVariantStock(ctx context.Context, orgID string, variantID uuid.UUID) (*models.VariantStock, error) {
	query := `
		SELECT
			COALESCE((SELECT SUM(quantity_change) FROM inventory_events WHERE product_variant_id = $1 AND organization_id = $2), 0) as total_stock,
			COALESCE((SELECT SUM(quantity) FROM inventory_reservations WHERE product_variant_id = $1 AND organization_id = $2 AND status = 'active' AND (expires_at IS NULL OR expires_at > now())), 0) as reserved_stock
	`
	var stock models.VariantStock
	err := s.dbtx.QueryRowContext(ctx, query, variantID, orgID).Scan(&stock.TotalStock, &stock.ReservedStock)
	if err != nil {
		return nil, err
	}
	stock.AvailableStock = stock.TotalStock - stock.ReservedStock
	return &stock, nil
}

func (s *Store) GetVariantWithLock(ctx context.Context, orgID string, variantID uuid.UUID) (*models.ProductVariant, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`SELECT id, organization_id, product_id, sku, barcode, price, cost, is_active, created_at, updated_at
		 FROM product_variants WHERE id = $1 AND organization_id = $2 FOR UPDATE`,
		variantID, orgID,
	)

	var v models.ProductVariant
	err := row.Scan(&v.ID, &v.OrganizationID, &v.ProductID, &v.SKU, &v.Barcode, &v.Price, &v.Cost, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &v, err
}

func (s *Store) UpdateReservationStatus(ctx context.Context, orgID string, reservationID uuid.UUID, status models.ReservationStatus) error {
	var releasedAt *time.Time
	if status == models.ReservationStatusReleased {
		now := time.Now()
		releasedAt = &now
	}

	_, err := s.dbtx.ExecContext(ctx,
		`UPDATE inventory_reservations SET status = $1, released_at = $2 WHERE id = $3 AND organization_id = $4`,
		status, releasedAt, reservationID, orgID,
	)
	return err
}

func (s *Store) GetReservation(ctx context.Context, orgID string, reservationID uuid.UUID) (*models.InventoryReservation, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`SELECT id, organization_id, product_variant_id, order_id, quantity, status, expires_at, created_at, released_at
		 FROM inventory_reservations WHERE id = $1 AND organization_id = $2`,
		reservationID, orgID,
	)

	var r models.InventoryReservation
	err := row.Scan(&r.ID, &r.OrganizationID, &r.ProductVariantID, &r.OrderID, &r.Quantity, &r.Status, &r.ExpiresAt, &r.CreatedAt, &r.ReleasedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &r, err
}

func (s *Store) ExpireReservations(ctx context.Context) (int64, error) {
	res, err := s.dbtx.ExecContext(ctx,
		`UPDATE inventory_reservations SET status = 'expired' WHERE status = 'active' AND expires_at < now()`,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
