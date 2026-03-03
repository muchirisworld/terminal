package store

import (
	"context"
)

// UpsertUser inserts or updates a user.
func (s *Store) UpsertUser(ctx context.Context, id, name, email, image string) error {
	_, err := s.dbtx.ExecContext(ctx, `
		INSERT INTO users (id, name, email, email_verified, image, created_at, updated_at)
		VALUES ($1, $2, $3, true, $4, now(), now())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			image = EXCLUDED.image,
			updated_at = now()
	`, id, name, email, image)
	return err
}

// DeleteUser deletes a user.
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, err := s.dbtx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

// UserExists checks if a user exists.
func (s *Store) UserExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.dbtx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	return exists, err
}

// UpsertOrganization inserts or updates an organization.
func (s *Store) UpsertOrganization(ctx context.Context, id, name, slug, imageUrl, createdBy string) error {
	_, err := s.dbtx.ExecContext(ctx, `
		INSERT INTO organizations (id, name, slug, image_url, created_by, is_active)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), true)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			slug = EXCLUDED.slug,
			image_url = EXCLUDED.image_url,
			is_active = true
	`, id, name, slug, imageUrl, createdBy)
	return err
}

// DeleteOrganization marks an organization as inactive.
func (s *Store) DeleteOrganization(ctx context.Context, id string) error {
	_, err := s.dbtx.ExecContext(ctx, "UPDATE organizations SET is_active = false WHERE id = $1", id)
	return err
}

// OrganizationExists checks if an organization exists.
func (s *Store) OrganizationExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.dbtx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM organizations WHERE id = $1)", id).Scan(&exists)
	return exists, err
}

// UpsertMembership inserts or updates a membership.
func (s *Store) UpsertMembership(ctx context.Context, orgID, userID, role string) error {
	_, err := s.dbtx.ExecContext(ctx, `
		INSERT INTO memberships (org_id, user_id, role, is_active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (user_id, org_id) DO UPDATE SET
			role = EXCLUDED.role,
			is_active = true
	`, orgID, userID, role)
	return err
}

// DeleteMembership marks a membership as inactive.
func (s *Store) DeleteMembership(ctx context.Context, orgID, userID string) error {
	_, err := s.dbtx.ExecContext(ctx, "UPDATE memberships SET is_active = false WHERE org_id = $1 AND user_id = $2", orgID, userID)
	return err
}

// UpsertInvitation inserts or updates an invitation.
func (s *Store) UpsertInvitation(ctx context.Context, id, orgID, email, role, status string) error {
	_, err := s.dbtx.ExecContext(ctx, `
		INSERT INTO invitations (id, org_id, email, role, status, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			role = EXCLUDED.role
	`, id, orgID, email, role, status)
	return err
}

// RevokeInvitation marks an invitation as revoked.
func (s *Store) RevokeInvitation(ctx context.Context, id string) error {
	_, err := s.dbtx.ExecContext(ctx, "UPDATE invitations SET status = 'revoked', is_active = false WHERE id = $1", id)
	return err
}
