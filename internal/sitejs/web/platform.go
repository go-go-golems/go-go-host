package web

import (
	"context"
	"net/http"
)

type platformContextKey struct{}

type PlatformContext struct {
	RequestID    string
	OrgID        string
	SiteID       string
	DeploymentID string
	Host         string
}

func WithPlatformContext(r *http.Request, platform PlatformContext) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), platformContextKey{}, platform))
}

func PlatformFromContext(ctx context.Context) (PlatformContext, bool) {
	platform, ok := ctx.Value(platformContextKey{}).(PlatformContext)
	return platform, ok
}
