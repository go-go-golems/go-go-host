package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-go-golems/go-go-host/internal/store"
	"github.com/google/uuid"
)

func TestRuntimeStatusPersistenceAndReconciliation(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	suffix := uuid.NewString()[:8]
	user, err := st.UpsertUserFromOIDC(ctx, "issuer", "runtime-"+suffix, "runtime@example.com", "Runtime")
	if err != nil {
		t.Fatalf("upsert user: %v", err)
	}
	org, err := st.CreateOrg(ctx, "runtime-org-"+suffix, "Runtime Org")
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	if err := st.AddMembership(ctx, org.ID, user.ID, store.RoleOrgOwner); err != nil {
		t.Fatalf("add membership: %v", err)
	}
	site, err := st.CreateSite(ctx, store.CreateSiteInput{OrgID: org.ID, Slug: "runtime-site-" + suffix, Name: "Runtime Site", PrimaryHost: "runtime-site-" + suffix + ".localhost"})
	if err != nil {
		t.Fatalf("create site: %v", err)
	}

	startedAt := time.Now().UTC().Add(-time.Minute)
	if err := st.UpsertRuntimeStatus(ctx, store.RuntimeStatus{
		SiteID:        site.ID,
		OrgID:         org.ID,
		DeploymentID:  "dep_test",
		Hosts:         []string{site.PrimaryHost},
		Status:        "ready",
		StartedAt:     startedAt,
		RequestsTotal: 7,
		ErrorsTotal:   1,
	}); err != nil {
		t.Fatalf("upsert runtime status: %v", err)
	}
	got, err := st.GetRuntimeStatus(ctx, site.ID)
	if err != nil {
		t.Fatalf("get runtime status: %v", err)
	}
	if got.Status != "ready" || got.RequestsTotal != 7 || got.ErrorsTotal != 1 || len(got.Hosts) != 1 {
		t.Fatalf("unexpected runtime status: %#v", got)
	}
	if err := st.ReconcileStaleRuntimeStatuses(ctx); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	reconciled, err := st.GetRuntimeStatus(ctx, site.ID)
	if err != nil {
		t.Fatalf("get reconciled status: %v", err)
	}
	if reconciled.Status != "stopped" {
		t.Fatalf("expected stopped after reconciliation, got %#v", reconciled)
	}
	if reconciled.LastError == "" {
		t.Fatalf("expected reconciliation message")
	}
}
