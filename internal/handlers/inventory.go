package handlers

import (
	"encoding/json"
	"errors"
	"github.com/muchirisworld/terminal/internal/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	ierrors "github.com/muchirisworld/terminal/internal/ierrors"
	"github.com/muchirisworld/terminal/internal/models"
	"github.com/muchirisworld/terminal/internal/service"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) UpsertConversion(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orgID := authCtx.OrgID
	productID, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req models.UpsertConversionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := h.service.UpsertConversion(r.Context(), orgID, productID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (h *InventoryHandler) ListConversionsByProduct(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orgID := authCtx.OrgID
	productID, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	conversions, err := h.service.ListConversionsByProduct(r.Context(), orgID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversions)
}

func (h *InventoryHandler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orgID := authCtx.OrgID
	variantID, err := uuid.Parse(chi.URLParam(r, "variantID"))
	if err != nil {
		http.Error(w, "invalid variant id", http.StatusBadRequest)
		return
	}

	var req models.ReceiptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.service.CreateInventoryReceipt(r.Context(), orgID, variantID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *InventoryHandler) CreateAdjustment(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orgID := authCtx.OrgID
	variantID, err := uuid.Parse(chi.URLParam(r, "variantID"))
	if err != nil {
		http.Error(w, "invalid variant id", http.StatusBadRequest)
		return
	}

	var req models.AdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.service.CreateInventoryAdjustment(r.Context(), orgID, variantID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *InventoryHandler) ReserveInventory(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orgID := authCtx.OrgID
	variantID, err := uuid.Parse(chi.URLParam(r, "variantID"))
	if err != nil {
		http.Error(w, "invalid variant id", http.StatusBadRequest)
		return
	}

	var req models.ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.service.ReserveInventory(r.Context(), orgID, variantID, &req)
	if err != nil {
		var insErr *ierrors.InsufficientStockError
		if errors.As(err, &insErr) {
			http.Error(w, insErr.Message, http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *InventoryHandler) ReleaseReservation(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orgID := authCtx.OrgID
	reservationID, err := uuid.Parse(chi.URLParam(r, "reservationID"))
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}

	err = h.service.ReleaseReservation(r.Context(), orgID, reservationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) GetVariantStock(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	orgID := authCtx.OrgID
	variantID, err := uuid.Parse(chi.URLParam(r, "variantID"))
	if err != nil {
		http.Error(w, "invalid variant id", http.StatusBadRequest)
		return
	}

	stock, err := h.service.GetVariantStock(r.Context(), orgID, variantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}
