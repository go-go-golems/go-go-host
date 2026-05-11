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
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
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
	idToken, err := verifier.Verify(r.Context(), token)
	if err != nil {
		return nil, fmt.Errorf("verify id token: %w", err)
	}
	var claims oidcClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("decode id token claims: %w", err)
	}
	displayName := claims.Name
	if displayName == "" {
		displayName = claims.PreferredUsername
	}
	if displayName == "" {
		displayName = idToken.Subject
	}
	return a.st.UpsertUserFromOIDC(r.Context(), issuer, idToken.Subject, claims.Email, displayName)
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
	a.verifier = provider.Verifier(&oidc.Config{ClientID: clientID})
	return a.verifier, nil
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
