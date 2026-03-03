package store

import (
	"context"
	"database/sql"

	"github.com/muchirisworld/terminal/internal/models"
)

// Create creates a new user in the database.
func (s *Store) CreateUser(ctx context.Context, userRequest *models.UserRequest) (*models.User, error) {
	row := s.dbtx.QueryRowContext(ctx,
		`INSERT INTO users (
			 id, name, email, email_verified, image
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, name, email, email_verified, image, created_at, updated_at`,
		userRequest.ID, userRequest.Name, userRequest.Email, userRequest.EmailVerified, userRequest.Image,
	)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}
