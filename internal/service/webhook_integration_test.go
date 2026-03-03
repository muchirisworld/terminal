package service_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/muchirisworld/terminal/internal/service"
	"github.com/muchirisworld/terminal/internal/store"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// Provide a default that matches local dev if possible, but skip if we can't connect
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to database: %v", err)
	}

	// Just checking if we can ping
	if err := db.Ping(); err != nil {
		t.Skipf("Skipping integration test: cannot ping database: %v", err)
	}

	// Truncate tables for a clean slate
	_, err = db.Exec(`
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`)
	if err != nil {
		// If tables don't exist, we might skip or fail. We assume migrations are run.
		t.Skipf("Failed to truncate tables (are migrations run?): %v", err)
	}

	return db
}

func TestWebhookService_DuplicateEvent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := service.NewWebhookService(store.New(db), logger)
	ctx := context.Background()

	payload := []byte(`{
		"instance_id": "evt_dup1",
		"type": "user.created",
		"data": {
			"id": "user_dup",
			"first_name": "Test",
			"last_name": "User",
			"email_addresses": [{"email_address": "dup@example.com"}],
			"image_url": ""
		}
	}`)

	// Process first time
	err := svc.Process(ctx, "evt_dup1", payload)
	if err != nil {
		t.Fatalf("expected no error on first process, got: %v", err)
	}

	// Process second time (Duplicate)
	err = svc.Process(ctx, "evt_dup1", payload)
	if err != nil {
		t.Fatalf("expected no error on duplicate process, got: %v", err)
	}

	// Verify only 1 user was created
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = 'user_dup'").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 user, got %d", count)
	}
}

func TestWebhookService_MembershipBeforeOrganization(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := service.NewWebhookService(store.New(db), logger)
	ctx := context.Background()

	// 1. Create a user first so we can attach them to a membership later
	userPayload := []byte(`{
		"instance_id": "evt_user_1",
		"type": "user.created",
		"data": {
			"id": "user_mem",
			"first_name": "Mem",
			"last_name": "User",
			"email_addresses": [{"email_address": "mem@example.com"}],
			"image_url": ""
		}
	}`)
	if err := svc.Process(ctx, "evt_user_1", userPayload); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// 2. Try to create membership before org
	membershipPayload := []byte(`{
		"instance_id": "evt_mem_1",
		"type": "organizationMembership.created",
		"data": {
			"id": "mem_1",
			"role": "org:admin",
			"organization": {"id": "org_1"},
			"public_user_data": {"user_id": "user_mem"}
		}
	}`)

	err := svc.Process(ctx, "evt_mem_1", membershipPayload)
	if err == nil {
		t.Fatalf("expected error processing membership before organization, got nil")
	}

	// Ensure webhook_events recorded the error
	var attempts int
	var lastError sql.NullString
	err = db.QueryRow("SELECT attempts, last_error FROM webhook_events WHERE instance_id = 'evt_mem_1'").Scan(&attempts, &lastError)
	if err != nil {
		t.Fatalf("failed to query webhook_events: %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
	if !lastError.Valid || lastError.String == "" {
		t.Errorf("expected last_error to be set")
	}

	// 3. Create the organization
	orgPayload := []byte(`{
		"instance_id": "evt_org_1",
		"type": "organization.created",
		"data": {
			"id": "org_1",
			"name": "Test Org",
			"slug": "test-org",
			"image_url": "",
			"created_by": "user_mem"
		}
	}`)
	if err := svc.Process(ctx, "evt_org_1", orgPayload); err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}

	// 4. Retry membership creation
	err = svc.Process(ctx, "evt_mem_1", membershipPayload)
	if err != nil {
		t.Fatalf("expected no error on retry after org is created, got: %v", err)
	}

	// Verify membership was created
	var memCount int
	err = db.QueryRow("SELECT COUNT(*) FROM memberships WHERE org_id = 'org_1' AND user_id = 'user_mem'").Scan(&memCount)
	if err != nil {
		t.Fatal(err)
	}
	if memCount != 1 {
		t.Errorf("expected 1 membership, got %d", memCount)
	}

	// Verify webhook event is now processed
	var processedAt sql.NullTime
	err = db.QueryRow("SELECT processed_at FROM webhook_events WHERE instance_id = 'evt_mem_1'").Scan(&processedAt)
	if err != nil {
		t.Fatal(err)
	}
	if !processedAt.Valid {
		t.Errorf("expected processed_at to be set on success")
	}
}

func TestWebhookService_SuccessfulProjectionUpdatesRows(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := service.NewWebhookService(store.New(db), logger)
	ctx := context.Background()

	// 1. Create User
	userPayload := []byte(`{
		"instance_id": "evt_proj_user",
		"type": "user.created",
		"data": {
			"id": "user_proj",
			"first_name": "Proj",
			"last_name": "User",
			"email_addresses": [{"email_address": "proj@example.com"}],
			"image_url": ""
		}
	}`)
	if err := svc.Process(ctx, "evt_proj_user", userPayload); err != nil {
		t.Fatalf("failed to process user: %v", err)
	}

	var email string
	err := db.QueryRow("SELECT email FROM users WHERE id = 'user_proj'").Scan(&email)
	if err != nil {
		t.Fatalf("failed to query user: %v", err)
	}
	if email != "proj@example.com" {
		t.Errorf("expected proj@example.com, got %s", email)
	}

	// 2. Update User
	userUpdatePayload := []byte(`{
		"instance_id": "evt_proj_user_upd",
		"type": "user.updated",
		"data": {
			"id": "user_proj",
			"first_name": "Proj",
			"last_name": "Updated",
			"email_addresses": [{"email_address": "proj-updated@example.com"}],
			"image_url": "http://image"
		}
	}`)
	if err := svc.Process(ctx, "evt_proj_user_upd", userUpdatePayload); err != nil {
		t.Fatalf("failed to process user update: %v", err)
	}

	var name, updatedEmail, image string
	err = db.QueryRow("SELECT name, email, image FROM users WHERE id = 'user_proj'").Scan(&name, &updatedEmail, &image)
	if err != nil {
		t.Fatalf("failed to query user after update: %v", err)
	}
	if name != "Proj Updated" || updatedEmail != "proj-updated@example.com" || image != "http://image" {
		t.Errorf("user update failed to project correctly")
	}

	// 3. Delete User
	userDeletePayload := []byte(`{
		"instance_id": "evt_proj_user_del",
		"type": "user.deleted",
		"data": {
			"id": "user_proj"
		}
	}`)
	if err := svc.Process(ctx, "evt_proj_user_del", userDeletePayload); err != nil {
		t.Fatalf("failed to process user delete: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = 'user_proj'").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 users after delete, got %d", count)
	}
}
