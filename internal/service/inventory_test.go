package service_test

import (
	"context"
	"testing"

	"github.com/muchirisworld/terminal/internal/models"
	"github.com/muchirisworld/terminal/internal/service"
	"github.com/muchirisworld/terminal/internal/store"
)

func TestInventoryService_CreateReceiptWithConversion(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mainStore := store.New(db)
	svc := service.NewInventoryService(mainStore)
	catalogSvc := service.NewCatalogService(mainStore)
	ctx := context.Background()
	orgID := "org_test_inv"

	// 1. Create Product
	p, err := catalogSvc.CreateProduct(ctx, orgID, &models.CreateProductRequest{
		Name:     "Beer",
		BaseUnit: "bottle",
	})
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create Variant
	v, err := catalogSvc.CreateVariant(ctx, orgID, p.ID, &models.CreateVariantRequest{
		SKU:   "BEER-001",
		Price: 5.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 3. Create Conversion (1 case = 24 bottles)
	_, err = svc.UpsertConversion(ctx, orgID, p.ID, &models.UpsertConversionRequest{
		UnitFrom: "case",
		UnitTo:   "bottle",
		Factor:   24.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 4. Receive 2 cases
	event, err := svc.CreateInventoryReceipt(ctx, orgID, v.ID, &models.ReceiptRequest{
		Quantity: 2,
		Unit:     "case",
	})
	if err != nil {
		t.Fatal(err)
	}

	if event.QuantityChange != 48 {
		t.Errorf("expected 48 bottles (2 cases * 24), got %d", event.QuantityChange)
	}

	// 5. Check stock
	stock, err := svc.GetVariantStock(ctx, orgID, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stock.AvailableStock != 48 {
		t.Errorf("expected 48 available stock, got %d", stock.AvailableStock)
	}
}

func TestInventoryService_ReserveInventory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mainStore := store.New(db)
	svc := service.NewInventoryService(mainStore)
	catalogSvc := service.NewCatalogService(mainStore)
	ctx := context.Background()
	orgID := "org_test_res"

	p, _ := catalogSvc.CreateProduct(ctx, orgID, &models.CreateProductRequest{Name: "Water", BaseUnit: "bottle"})
	v, _ := catalogSvc.CreateVariant(ctx, orgID, p.ID, &models.CreateVariantRequest{SKU: "WATER-001", Price: 1.0})

	// Initial stock: 10 bottles
	svc.CreateInventoryAdjustment(ctx, orgID, v.ID, &models.AdjustmentRequest{QuantityChange: 10})

	// 1. Reserve 4 bottles
	res, err := svc.ReserveInventory(ctx, orgID, v.ID, &models.ReservationRequest{
		Quantity: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Quantity != 4 {
		t.Errorf("expected reservation of 4, got %d", res.Quantity)
	}

	// 2. Check stock (Available should be 6)
	stock, _ := svc.GetVariantStock(ctx, orgID, v.ID)
	if stock.AvailableStock != 6 {
		t.Errorf("expected 6 available, got %d", stock.AvailableStock)
	}

	// 3. Try to reserve 7 (should fail)
	_, err = svc.ReserveInventory(ctx, orgID, v.ID, &models.ReservationRequest{Quantity: 7})
	if err == nil {
		t.Error("expected error due to insufficient stock, got nil")
	}

	// 4. Release reservation
	err = svc.ReleaseReservation(ctx, orgID, res.ID)
	if err != nil {
		t.Fatal(err)
	}

	// 5. Check stock (Available should be 10 again)
	stock, _ = svc.GetVariantStock(ctx, orgID, v.ID)
	if stock.AvailableStock != 10 {
		t.Errorf("expected 10 available after release, got %d", stock.AvailableStock)
	}
}
