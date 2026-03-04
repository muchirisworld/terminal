package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/muchirisworld/terminal/internal/models"
	"github.com/muchirisworld/terminal/internal/service"
)

type CatalogHandler struct {
	service *service.CatalogService
}

func NewCatalogHandler(service *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

func getOrgID(r *http.Request) string {
	return r.Header.Get("X-Organization-ID")
}

func (h *CatalogHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	if orgID == "" {
		http.Error(w, "missing organization id", http.StatusBadRequest)
		return
	}

	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.service.CreateProduct(r.Context(), orgID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *CatalogHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	id, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.service.UpdateProduct(r.Context(), orgID, id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *CatalogHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	id, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	p, err := h.service.GetProduct(r.Context(), orgID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if p == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	products, err := h.service.ListProducts(r.Context(), orgID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *CatalogHandler) ArchiveProduct(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	id, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	p, err := h.service.ArchiveProduct(r.Context(), orgID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *CatalogHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	productID, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req models.CreateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v, err := h.service.CreateVariant(r.Context(), orgID, productID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *CatalogHandler) ListVariantsByProduct(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	productID, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	variants, err := h.service.ListVariantsByProduct(r.Context(), orgID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variants)
}

func (h *CatalogHandler) GetVariant(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	id, err := uuid.Parse(chi.URLParam(r, "variantID"))
	if err != nil {
		http.Error(w, "invalid variant id", http.StatusBadRequest)
		return
	}

	v, err := h.service.GetVariant(r.Context(), orgID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if v == nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *CatalogHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	orgID := getOrgID(r)
	id, err := uuid.Parse(chi.URLParam(r, "variantID"))
	if err != nil {
		http.Error(w, "invalid variant id", http.StatusBadRequest)
		return
	}

	var req models.UpdateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v, err := h.service.UpdateVariant(r.Context(), orgID, id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
