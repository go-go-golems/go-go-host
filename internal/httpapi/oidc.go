package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type oidcAuthenticator struct {
	cfg config.Config
	st  *store.Store

	mu       sync.Mutex
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

type oidcClaims struct {
	Audience          []string                         `json:"aud"`
	AuthorizedParty   string                           `json:"azp"`
	Email             string                           `json:"email"`
	PreferredUsername string                           `json:"preferred_username"`
	Name              string                           `json:"name"`
	RealmAccess       oidcRoleClaim                    `json:"realm_access"`
	ResourceAccess    map[string]oidcResourceRoleClaim `json:"resource_access"`
	Groups            []string                         `json:"groups"`
}

type oidcRoleClaim struct {
	Roles []string `json:"roles"`
}

type oidcResourceRoleClaim struct {
	Roles []string `json:"roles"`
}

func (a *oidcAuthenticator) authenticate(r *http.Request) (*store.User, error) {
	if a.st == nil {
		return nil, fmt.Errorf("store is not configured")
	}
	issuer := strings.TrimSpace(a.cfg.OIDCIssuer)
	clientID := strings.TrimSpace(a.cfg.OIDCClientID)
	if issuer == "" || clientID == "" {
		return nil, fmt.Errorf("oidcIssuer and oidcClientId are required when devAuth is false")
	}
	token := bearerToken(r)
	if token == "" {
		return nil, errUnauthenticated
	}
	verifier, err := a.getVerifier(r.Context(), issuer, clientID)
	if err != nil {
		return nil, err
	}
	verifiedToken, err := verifier.Verify(r.Context(), token)
	if err != nil {
		return nil, fmt.Errorf("verify oidc bearer token: %w", err)
	}
	var claims oidcClaims
	if err := verifiedToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("decode oidc bearer token claims: %w", err)
	}
	if !tokenMatchesClient(clientID, verifiedToken.Audience, claims) {
		return nil, fmt.Errorf("verify oidc bearer token: expected audience or authorized party %q", clientID)
	}
	displayName := claims.Name
	if displayName == "" {
		displayName = claims.PreferredUsername
	}
	if displayName == "" {
		displayName = verifiedToken.Subject
	}
	user, err := a.st.UpsertUserFromOIDC(r.Context(), issuer, verifiedToken.Subject, claims.Email, displayName)
	if err != nil {
		return nil, err
	}
	if shouldBootstrapPlatformAdmin(a.cfg, verifiedToken.Subject, claims) {
		alreadyAdmin, err := a.st.IsPlatformAdmin(r.Context(), user.ID)
		if err != nil {
			return nil, err
		}
		if err := a.st.AddPlatformAdmin(r.Context(), user.ID); err != nil {
			return nil, err
		}
		if !alreadyAdmin {
			_, err := a.st.InsertAuditEvent(r.Context(), store.AuditEvent{
				OrgID:        "",
				ActorType:    "system",
				ActorID:      "oidc-bootstrap",
				Action:       "platform_admin.bootstrap",
				ResourceType: "user",
				ResourceID:   user.ID,
				IPAddress:    r.RemoteAddr,
				UserAgent:    r.UserAgent(),
				MetadataJSON: fmt.Sprintf(`{"issuer":%q,"subject":%q,"email":%q}`, issuer, verifiedToken.Subject, claims.Email),
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return user, nil
}

func shouldBootstrapPlatformAdmin(cfg config.Config, subject string, claims oidcClaims) bool {
	for _, configured := range cfg.PlatformAdminOIDCSubjects {
		if strings.EqualFold(strings.TrimSpace(configured), strings.TrimSpace(subject)) {
			return true
		}
	}
	for _, configured := range cfg.PlatformAdminEmails {
		if strings.EqualFold(strings.TrimSpace(configured), strings.TrimSpace(claims.Email)) {
			return true
		}
	}
	if len(cfg.PlatformAdminOIDCRoles) == 0 {
		return false
	}
	actualRoles := map[string]struct{}{}
	for _, role := range claims.RealmAccess.Roles {
		actualRoles[strings.TrimSpace(role)] = struct{}{}
	}
	for _, role := range claims.Groups {
		actualRoles[strings.Trim(strings.TrimSpace(role), "/")] = struct{}{}
	}
	for _, resource := range claims.ResourceAccess {
		for _, role := range resource.Roles {
			actualRoles[strings.TrimSpace(role)] = struct{}{}
		}
	}
	for _, configured := range cfg.PlatformAdminOIDCRoles {
		if _, ok := actualRoles[strings.Trim(strings.TrimSpace(configured), "/")]; ok {
			return true
		}
	}
	return false
}

func (a *oidcAuthenticator) getVerifier(ctx context.Context, issuer, clientID string) (*oidc.IDTokenVerifier, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.verifier != nil {
		return a.verifier, nil
	}
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("discover OIDC provider: %w", err)
	}
	a.provider = provider
	// Keycloak access tokens issued to public browser clients may not include
	// the client ID in aud unless an audience mapper is configured. They do
	// include azp (authorized party), while ID tokens include aud. Verify the
	// signature/issuer/expiry here, then enforce aud-or-azp locally so the API
	// can accept both dashboard ID tokens and API access tokens without
	// weakening issuer or expiry checks.
	a.verifier = provider.Verifier(&oidc.Config{ClientID: clientID, SkipClientIDCheck: true})
	return a.verifier, nil
}

func tokenMatchesClient(clientID string, tokenAudience []string, claims oidcClaims) bool {
	clientID = strings.TrimSpace(clientID)
	for _, audience := range tokenAudience {
		if strings.EqualFold(strings.TrimSpace(audience), clientID) {
			return true
		}
	}
	for _, audience := range claims.Audience {
		if strings.EqualFold(strings.TrimSpace(audience), clientID) {
			return true
		}
	}
	return strings.EqualFold(strings.TrimSpace(claims.AuthorizedParty), clientID)
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
