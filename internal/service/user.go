package service

import (
	"context"

	ierrors "github.com/muchirisworld/terminal/internal/ierrors"
	"github.com/muchirisworld/terminal/internal/models"
)

// UserStore is the interface for the user store.
type UserStore interface {
	Create(ctx context.Context, user *models.UserRequest) (*models.User, error)
}

// UserService is the service for the user model.
type UserService struct {
	store UserStore
}

// NewUserService creates a new UserService.
func NewUserService(store UserStore) *UserService {
	return &UserService{store: store}
}

// Create creates a new user.
func (s *UserService) Create(ctx context.Context, req *models.UserRequest) (*models.User, error) {
	if req.Name == "" {
		return nil, &ierrors.ValidationError{Message: "name cannot be empty"}
	}

	if req.Email == "" {
		return nil, &ierrors.ValidationError{Message: "email cannot be empty"}
	}

	user, err := s.store.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
