package store_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-go-golems/go-go-host/internal/store"
	"github.com/google/uuid"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	dsn := os.Getenv("GO_GO_HOST_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("GO_GO_HOST_TEST_DATABASE_URL is not set; skipping Postgres integration test")
	}
	ctx := context.Background()
	st, err := store.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(st.Close)
	if err := st.ApplyMigrations(ctx); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	return st
}

func TestMigrationsAndCoreRows(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	user, err := st.UpsertUserFromOIDC(ctx, "issuer", "subject", "dev@example.com", "Dev User")
	if err != nil {
		t.Fatalf("upsert user: %v", err)
	}
	suffix := uuid.NewString()[:8]
	org, err := st.CreateOrg(ctx, "parc-"+suffix, "PARC")
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	if err := st.AddMembership(ctx, org.ID, user.ID, store.RoleOrgOwner); err != nil {
		t.Fatalf("add membership: %v", err)
	}
	role, err := st.MembershipRole(ctx, org.ID, user.ID)
	if err != nil {
		t.Fatalf("membership role: %v", err)
	}
	if role != store.RoleOrgOwner {
		t.Fatalf("expected owner role, got %q", role)
	}
	site, err := st.CreateSite(ctx, store.CreateSiteInput{OrgID: org.ID, Slug: "trail-" + suffix, Name: "Trail", PrimaryHost: "trail-" + suffix + ".localhost"})
	if err != nil {
		t.Fatalf("create site: %v", err)
	}
	if err := st.CreateDefaultSiteQuota(ctx, site.ID); err != nil {
		t.Fatalf("create quota: %v", err)
	}
	if _, err := st.GetSiteQuota(ctx, site.ID); err != nil {
		t.Fatalf("get quota: %v", err)
	}
	if _, err := st.InsertAuditEvent(ctx, store.AuditEvent{OrgID: org.ID, ActorType: "user", ActorID: user.ID, Action: "site.create", ResourceType: "site", ResourceID: site.ID}); err != nil {
		t.Fatalf("insert audit: %v", err)
	}
	events, err := st.ListAuditEventsForOrg(ctx, org.ID, 10)
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
}

func TestUpsertUserFromOIDCUpdatesExistingUser(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	subject := "subject-" + uuid.NewString()
	first, err := st.UpsertUserFromOIDC(ctx, "issuer", subject, "old@example.com", "Old")
	if err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	second, err := st.UpsertUserFromOIDC(ctx, "issuer", subject, "new@example.com", "New")
	if err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("expected same user ID, got %s and %s", first.ID, second.ID)
	}
	if second.Email != "new@example.com" || second.DisplayName != "New" {
		t.Fatalf("expected updated profile, got %#v", second)
	}
}
