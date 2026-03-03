package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/muchirisworld/terminal/internal/store"
)

type clerkHandlerFunc func(ctx context.Context, s *store.Store, evt Event) error

var clerkHandlers = map[string]clerkHandlerFunc{
	"user.created":                    handleUserUpsert,
	"user.updated":                    handleUserUpsert,
	"user.deleted":                    handleUserDeleted,
	"organization.created":            handleOrganizationUpsert,
	"organization.updated":            handleOrganizationUpsert,
	"organization.deleted":            handleOrganizationDeleted,
	"organizationMembership.created":  handleMembershipUpsert,
	"organizationMembership.updated":  handleMembershipUpsert,
	"organizationMembership.deleted":  handleMembershipDeleted,
	"organizationInvitation.created":  handleInvitationUpsert,
	"organizationInvitation.accepted": handleInvitationUpsert,
	"organizationInvitation.revoked":  handleInvitationRevoked,
}

// GetClerkHandler returns the handler for the given event type.
func GetClerkHandler(eventType string) (clerkHandlerFunc, bool) {
	h, ok := clerkHandlers[eventType]
	return h, ok
}

// Structs mapping to Clerk payload formats
type UserPayload struct {
	ID           string `json:"id"`
	EmailAddress []struct {
		EmailAddress string `json:"email_address"`
	} `json:"email_addresses"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ImageURL  string `json:"image_url"`
}

func handleUserUpsert(ctx context.Context, s *store.Store, evt Event) error {
	var payload UserPayload
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}

	name := strings.TrimSpace(payload.FirstName + " " + payload.LastName)
	if name == "" {
		name = "Unnamed User"
	}
	email := ""
	if len(payload.EmailAddress) > 0 {
		email = payload.EmailAddress[0].EmailAddress
	}

	return s.UpsertUser(ctx, payload.ID, name, email, payload.ImageURL)
}

func handleUserDeleted(ctx context.Context, s *store.Store, evt Event) error {
	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}

	return s.DeleteUser(ctx, payload.ID)
}

type OrganizationPayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ImageURL  string `json:"image_url"`
	CreatedBy string `json:"created_by"`
}

func handleOrganizationUpsert(ctx context.Context, s *store.Store, evt Event) error {
	var payload OrganizationPayload
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}

	// Out-of-order check: user must exist if created_by is set
	if payload.CreatedBy != "" {
		exists, err := s.UserExists(ctx, payload.CreatedBy)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("dependency not satisfied: user %s does not exist yet", payload.CreatedBy)
		}
	}

	return s.UpsertOrganization(ctx, payload.ID, payload.Name, payload.Slug, payload.ImageURL, payload.CreatedBy)
}

func handleOrganizationDeleted(ctx context.Context, s *store.Store, evt Event) error {
	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}
	return s.DeleteOrganization(ctx, payload.ID)
}

type MembershipPayload struct {
	ID           string `json:"id"`
	Role         string `json:"role"`
	Organization struct {
		ID string `json:"id"`
	} `json:"organization"`
	PublicUserData struct {
		UserID string `json:"user_id"`
	} `json:"public_user_data"`
}

func mapRole(clerkRole string) string {
	if strings.HasSuffix(clerkRole, "admin") {
		return "admin"
	}
	return "member"
}

func handleMembershipUpsert(ctx context.Context, s *store.Store, evt Event) error {
	var payload MembershipPayload
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}

	orgID := payload.Organization.ID
	userID := payload.PublicUserData.UserID
	role := mapRole(payload.Role)

	orgExists, err := s.OrganizationExists(ctx, orgID)
	if err != nil {
		return err
	}
	if !orgExists {
		return fmt.Errorf("dependency not satisfied: org %s does not exist", orgID)
	}

	userExists, err := s.UserExists(ctx, userID)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf("dependency not satisfied: user %s does not exist", userID)
	}

	return s.UpsertMembership(ctx, orgID, userID, role)
}

func handleMembershipDeleted(ctx context.Context, s *store.Store, evt Event) error {
	var payload MembershipPayload
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}

	return s.DeleteMembership(ctx, payload.Organization.ID, payload.PublicUserData.UserID)
}

type InvitationPayload struct {
	ID             string `json:"id"`
	EmailAddress   string `json:"email_address"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	Status         string `json:"status"`
}

func mapInvitationStatus(status string) string {
	switch status {
	case "pending":
		return "pending"
	case "accepted":
		return "accepted"
	case "revoked":
		return "revoked"
	default:
		return "pending" // Default to pending
	}
}

func handleInvitationUpsert(ctx context.Context, s *store.Store, evt Event) error {
	var payload InvitationPayload
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}

	orgExists, err := s.OrganizationExists(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}
	if !orgExists {
		return fmt.Errorf("dependency not satisfied: org %s does not exist", payload.OrganizationID)
	}

	role := mapRole(payload.Role)
	status := mapInvitationStatus(payload.Status)

	return s.UpsertInvitation(ctx, payload.ID, payload.OrganizationID, payload.EmailAddress, role, status)
}

func handleInvitationRevoked(ctx context.Context, s *store.Store, evt Event) error {
	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(evt.Data, &payload); err != nil {
		return err
	}

	return s.RevokeInvitation(ctx, payload.ID)
}
