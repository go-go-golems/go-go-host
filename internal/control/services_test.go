package control_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/store"
	"github.com/google/uuid"
)

func newCore(t *testing.T) (*control.Core, *store.Store, context.Context) {
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
		t.Fatalf("migrations: %v", err)
	}
	cfg := config.Default()
	cfg.BaseDomain = "example.test"
	return control.NewCoreWithStore(cfg, st), st, ctx
}

func TestOrgOwnerCanCreateSiteAndViewerCannot(t *testing.T) {
	core, st, ctx := newCore(t)
	suffix := uuid.NewString()[:8]
	owner, _ := st.UpsertUserFromOIDC(ctx, "issuer", "owner-"+suffix, "owner@example.com", "Owner")
	viewer, _ := st.UpsertUserFromOIDC(ctx, "issuer", "viewer-"+suffix, "viewer@example.com", "Viewer")
	org, err := core.Orgs.CreateOrg(ctx, owner.ID, "parc-"+suffix, "PARC")
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	if err := st.AddMembership(ctx, org.ID, viewer.ID, store.RoleOrgViewer); err != nil {
		t.Fatalf("add viewer: %v", err)
	}
	site, err := core.Sites.CreateSite(ctx, owner.ID, org.ID, "trail-"+suffix, "Trail")
	if err != nil {
		t.Fatalf("owner create site: %v", err)
	}
	if site.PrimaryHost != "trail-"+suffix+".example.test" {
		t.Fatalf("unexpected host %q", site.PrimaryHost)
	}
	if _, err := core.Sites.CreateSite(ctx, viewer.ID, org.ID, "viewer-site", "Viewer Site"); !errors.Is(err, control.ErrPermissionDenied) {
		t.Fatalf("expected permission denied for viewer, got %v", err)
	}
}

func TestCrossOrgSiteListDenied(t *testing.T) {
	core, st, ctx := newCore(t)
	suffix := uuid.NewString()[:8]
	alice, _ := st.UpsertUserFromOIDC(ctx, "issuer", "alice-"+suffix, "alice@example.com", "Alice")
	bob, _ := st.UpsertUserFromOIDC(ctx, "issuer", "bob-"+suffix, "bob@example.com", "Bob")
	org, err := core.Orgs.CreateOrg(ctx, alice.ID, "alice-org-"+suffix, "Alice Org")
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	if _, err := core.Sites.CreateSite(ctx, alice.ID, org.ID, "trail-"+suffix, "Trail"); err != nil {
		t.Fatalf("create site: %v", err)
	}
	if _, err := core.Sites.ListSites(ctx, bob.ID, org.ID); !errors.Is(err, control.ErrPermissionDenied) {
		t.Fatalf("expected cross-org permission denied, got %v", err)
	}
}
