package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"

	ierrors "github.com/muchirisworld/terminal/internal/ierrors"
	"github.com/muchirisworld/terminal/internal/models"
)

const MAX_BYTES = 1 << 20 // 1MB

// UserService is the interface for the user service.
type UserService interface {
	Create(ctx context.Context, req *models.UserRequest) (*models.User, error)
}

// UserHandler is the handler for the user model.
type UserHandler struct {
	service UserService
	logger  *slog.Logger
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// Create is the handler for creating a new user.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, MAX_BYTES)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
		}
		return
	}

		var userReq models.UserRequest
		if err := json.Unmarshal(body, &userReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := h.service.Create(r.Context(), &userReq)
	if err != nil {
		var validationErr *ierrors.ValidationError
		if errors.As(err, &validationErr) {
			http.Error(w, validationErr.Message, http.StatusBadRequest)
			return
		}
		h.logger.Error("failed to create user", "err", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.logger.Error("failed to encode user", "err", err)
	}
}
