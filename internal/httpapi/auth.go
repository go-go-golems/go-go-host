package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-go-golems/go-go-host/internal/store"
)

type principal struct {
	User *store.User
}

type principalContextKey struct{}

func withPrincipal(ctx context.Context, p principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, p)
}

func principalFromContext(ctx context.Context) (principal, bool) {
	p, ok := ctx.Value(principalContextKey{}).(principal)
	return p, ok
}

var errUnauthenticated = errors.New("unauthenticated")

func authMiddleware(next http.Handler, authn *oidcAuthenticator, devAuthEnabled bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authn == nil || authn.st == nil {
			writeError(w, http.StatusServiceUnavailable, "store is not configured")
			return
		}
		var user *store.User
		var err error
		if devAuthEnabled {
			user, err = authenticateDev(r, authn.st, authn.cfg.DevPlatformAdminSubjects)
		} else {
			user, err = authn.authenticate(r)
		}
		if err != nil {
			status := http.StatusUnauthorized
			if err != errUnauthenticated {
				status = http.StatusInternalServerError
			}
			writeError(w, status, err.Error())
			return
		}
		next.ServeHTTP(w, r.WithContext(withPrincipal(r.Context(), principal{User: user})))
	})
}

func authenticateDev(r *http.Request, st *store.Store, platformAdminSubjects []string) (*store.User, error) {
	subject := r.Header.Get("X-Go-Go-Host-User")
	if subject == "" {
		subject = "dev-user"
	}
	email := r.Header.Get("X-Go-Go-Host-Email")
	if email == "" {
		email = subject + "@dev.local"
	}
	name := r.Header.Get("X-Go-Go-Host-Name")
	if name == "" {
		name = subject
	}
	user, err := st.UpsertUserFromOIDC(r.Context(), "dev", subject, email, name)
	if err != nil {
		return nil, err
	}
	for _, adminSubject := range platformAdminSubjects {
		if adminSubject == subject {
			if err := st.AddPlatformAdmin(r.Context(), user.ID); err != nil {
				return nil, err
			}
			break
		}
	}
	return user, nil
}

func requirePrincipal(r *http.Request) (principal, error) {
	p, ok := principalFromContext(r.Context())
	if !ok || p.User == nil {
		return principal{}, errUnauthenticated
	}
	return p, nil
}
