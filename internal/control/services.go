package control

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-go-golems/go-go-host/internal/store"
)

var slugRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)

type OrgService struct{ store *store.Store }

type SiteService struct {
	store      *store.Store
	baseDomain string
}

func (s *OrgService) CreateOrg(ctx context.Context, actorUserID, slug, name string) (*store.Org, error) {
	if s.store == nil {
		return nil, errors.New("store is not configured")
	}
	if !slugRE.MatchString(slug) {
		return nil, fmt.Errorf("invalid org slug %q", slug)
	}
	org, err := s.store.CreateOrg(ctx, slug, name)
	if err != nil {
		return nil, err
	}
	if err := s.store.AddMembership(ctx, org.ID, actorUserID, store.RoleOrgOwner); err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: org.ID, ActorType: "user", ActorID: actorUserID, Action: "org.create", ResourceType: "org", ResourceID: org.ID})
	return org, nil
}

func (s *OrgService) EnsureRole(ctx context.Context, userID, orgID string, allowed ...string) error {
	role, err := s.store.MembershipRole(ctx, orgID, userID)
	if err != nil {
		return err
	}
	for _, a := range allowed {
		if role == a {
			return nil
		}
	}
	return ErrPermissionDenied
}

func (s *SiteService) CreateSite(ctx context.Context, actorUserID, orgID, slug, name string) (*store.Site, error) {
	if s.store == nil {
		return nil, errors.New("store is not configured")
	}
	if !slugRE.MatchString(slug) {
		return nil, fmt.Errorf("invalid site slug %q", slug)
	}
	role, err := s.store.MembershipRole(ctx, orgID, actorUserID)
	if err != nil {
		return nil, err
	}
	if role != store.RoleOrgOwner && role != store.RoleOrgDeveloper {
		return nil, ErrPermissionDenied
	}
	host := slug
	baseDomain := strings.Trim(s.baseDomain, ".")
	if baseDomain != "" && baseDomain != "localhost" {
		host = slug + "." + baseDomain
	} else if baseDomain == "localhost" {
		host = slug + ".localhost"
	}
	site, err := s.store.CreateSite(ctx, store.CreateSiteInput{OrgID: orgID, Slug: slug, Name: name, PrimaryHost: host})
	if err != nil {
		return nil, err
	}
	if err := s.store.CreateDefaultSiteQuota(ctx, site.ID); err != nil {
		return nil, err
	}
	if err := s.store.CreateDefaultSiteCapabilities(ctx, site.ID); err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: orgID, ActorType: "user", ActorID: actorUserID, Action: "site.create", ResourceType: "site", ResourceID: site.ID})
	return site, nil
}

func (s *SiteService) ListSites(ctx context.Context, actorUserID, orgID string) ([]store.Site, error) {
	role, err := s.store.MembershipRole(ctx, orgID, actorUserID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, ErrPermissionDenied
	}
	return s.store.ListSitesByOrg(ctx, orgID)
}

var ErrPermissionDenied = errors.New("permission denied")
