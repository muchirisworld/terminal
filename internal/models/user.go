package models

import "time"

// User represents a user in the database.
type User struct {
	ID            string    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"email"`
	EmailVerified bool      `db:"email_verified" json:"email_verified"`
	Image         *string   `db:"image" json:"image"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type UserRequest struct {
	ID            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	Email         string `db:"email" json:"email"`
	EmailVerified bool   `db:"email_verified" json:"email_verified"`
	Image         string `db:"image" json:"image"`
}