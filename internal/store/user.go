package store

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/muchirisworld/terminal/internal/models"
)

// UserStore is the repository for the user model.
type UserStore struct {
	db *sqlx.DB
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

// Create creates a new user in the database.
func (s *UserStore) Create(ctx context.Context, userRequest *models.UserRequest) (*models.User, error) {
	rows, err := s.db.NamedQueryContext(ctx,
		`INSERT INTO users (
			 id, name, email, email_verified, image
		) VALUES (
			:id, :name, :email, :email_verified, :image
		) RETURNING *`,
		userRequest,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	if err := rows.StructScan(&user); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}
