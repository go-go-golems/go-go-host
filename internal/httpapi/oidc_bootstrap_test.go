package httpapi

import (
	"testing"

	"github.com/go-go-golems/go-go-host/internal/config"
)

func TestTokenMatchesClientFromAudienceOrAuthorizedParty(t *testing.T) {
	clientID := "go-go-host-dashboard"
	if !tokenMatchesClient(clientID, []string{clientID}, oidcClaims{}) {
		t.Fatalf("expected token audience to match client")
	}
	if !tokenMatchesClient(clientID, nil, oidcClaims{Audience: []string{"other", clientID}}) {
		t.Fatalf("expected claims audience to match client")
	}
	if !tokenMatchesClient(clientID, nil, oidcClaims{AuthorizedParty: clientID}) {
		t.Fatalf("expected authorized party to match client")
	}
	if tokenMatchesClient(clientID, []string{"other"}, oidcClaims{AuthorizedParty: "other"}) {
		t.Fatalf("did not expect unrelated token to match client")
	}
	if tokenMatchesClient("", []string{clientID}, oidcClaims{AuthorizedParty: clientID}) {
		t.Fatalf("empty configured client id must not match")
	}
}

func TestTokenMatchesAnyAcceptedClient(t *testing.T) {
	accepted := []string{"go-go-host-dashboard", "go-go-host-cli"}
	if !tokenMatchesAnyClient(accepted, nil, oidcClaims{AuthorizedParty: "go-go-host-cli"}) {
		t.Fatalf("expected CLI authorized party to match accepted clients")
	}
	if !tokenMatchesAnyClient(accepted, []string{"go-go-host-cli"}, oidcClaims{}) {
		t.Fatalf("expected CLI audience to match accepted clients")
	}
	if tokenMatchesAnyClient(accepted, []string{"unknown"}, oidcClaims{AuthorizedParty: "unknown"}) {
		t.Fatalf("did not expect unknown client to match accepted clients")
	}
}

func TestShouldBootstrapPlatformAdminFromSubjectEmailAndRoles(t *testing.T) {
	cfg := config.Default()
	cfg.PlatformAdminOIDCSubjects = []string{"admin-sub"}
	cfg.PlatformAdminEmails = []string{"ops@example.com"}
	cfg.PlatformAdminOIDCRoles = []string{"go-go-host-admin", "platform-admins"}

	if !shouldBootstrapPlatformAdmin(cfg, "admin-sub", oidcClaims{}) {
		t.Fatalf("expected matching subject to bootstrap admin")
	}
	if !shouldBootstrapPlatformAdmin(cfg, "other", oidcClaims{Email: "OPS@example.com"}) {
		t.Fatalf("expected matching email to bootstrap admin")
	}
	if !shouldBootstrapPlatformAdmin(cfg, "other", oidcClaims{RealmAccess: oidcRoleClaim{Roles: []string{"go-go-host-admin"}}}) {
		t.Fatalf("expected matching realm role to bootstrap admin")
	}
	if !shouldBootstrapPlatformAdmin(cfg, "other", oidcClaims{Groups: []string{"/platform-admins"}}) {
		t.Fatalf("expected matching group to bootstrap admin")
	}
	if !shouldBootstrapPlatformAdmin(cfg, "other", oidcClaims{ResourceAccess: map[string]oidcResourceRoleClaim{"go-go-host-dashboard": {Roles: []string{"go-go-host-admin"}}}}) {
		t.Fatalf("expected matching client role to bootstrap admin")
	}
	if shouldBootstrapPlatformAdmin(cfg, "other", oidcClaims{Email: "user@example.com", RealmAccess: oidcRoleClaim{Roles: []string{"viewer"}}}) {
		t.Fatalf("did not expect unrelated user to bootstrap admin")
	}
}
