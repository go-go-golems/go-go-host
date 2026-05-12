package control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-go-golems/go-go-host/internal/store"
)

var slugRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)
var configKeyRE = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.-]{0,63}$`)
var hostnameRE = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)+$`)

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

func (s *SiteService) ListConfig(ctx context.Context, actorUserID, siteID string) ([]store.SiteConfigItem, error) {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper, store.RoleOrgViewer); err != nil {
		return nil, err
	}
	return s.store.ListSiteConfig(ctx, siteID)
}

func (s *SiteService) UpsertConfig(ctx context.Context, actorUserID, siteID, key string, valueJSON []byte) error {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper); err != nil {
		return err
	}
	if !configKeyRE.MatchString(key) {
		return fmt.Errorf("invalid config key %q", key)
	}
	if !json.Valid(valueJSON) {
		return fmt.Errorf("invalid config JSON for %q", key)
	}
	if err := s.store.UpsertSiteConfig(ctx, siteID, key, valueJSON); err != nil {
		return err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: actorUserID, Action: "site.config.upsert", ResourceType: "site", ResourceID: site.ID, MetadataJSON: fmt.Sprintf(`{"key":%q}`, key)})
	return nil
}

func (s *SiteService) DeleteConfig(ctx context.Context, actorUserID, siteID, key string) error {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper); err != nil {
		return err
	}
	if err := s.store.DeleteSiteConfig(ctx, siteID, key); err != nil {
		return err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: actorUserID, Action: "site.config.delete", ResourceType: "site", ResourceID: site.ID, MetadataJSON: fmt.Sprintf(`{"key":%q}`, key)})
	return nil
}

func (s *SiteService) ListCapabilities(ctx context.Context, actorUserID, siteID string) ([]store.SiteCapability, error) {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper, store.RoleOrgViewer); err != nil {
		return nil, err
	}
	return s.store.ListSiteCapabilities(ctx, siteID)
}

func (s *SiteService) UpsertCapability(ctx context.Context, actorUserID, siteID, capability string, enabled bool, configJSON []byte) error {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner); err != nil {
		return err
	}
	if strings.TrimSpace(capability) == "" {
		return errors.New("capability is required")
	}
	if len(configJSON) == 0 {
		configJSON = []byte("{}")
	}
	if !json.Valid(configJSON) {
		return fmt.Errorf("invalid capability config JSON for %q", capability)
	}
	if err := s.store.UpsertSiteCapability(ctx, siteID, capability, enabled, configJSON); err != nil {
		return err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: actorUserID, Action: "site.capability.update", ResourceType: "site", ResourceID: site.ID, MetadataJSON: fmt.Sprintf(`{"capability":%q,"enabled":%t}`, capability, enabled)})
	return nil
}

func (s *SiteService) ListDomains(ctx context.Context, actorUserID, siteID string) ([]store.SiteDomain, error) {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper, store.RoleOrgViewer); err != nil {
		return nil, err
	}
	return s.store.ListSiteDomains(ctx, siteID)
}

func (s *SiteService) AddDomain(ctx context.Context, actorUserID, siteID, hostname string) (*store.SiteDomain, error) {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper); err != nil {
		return nil, err
	}
	hostname = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(hostname), "."))
	if !hostnameRE.MatchString(hostname) {
		return nil, fmt.Errorf("invalid hostname %q", hostname)
	}
	domain, err := s.store.CreateSiteDomain(ctx, siteID, hostname)
	if err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: actorUserID, Action: "site.domain.add", ResourceType: "site_domain", ResourceID: domain.ID, MetadataJSON: fmt.Sprintf(`{"hostname":%q}`, hostname)})
	return domain, nil
}

func (s *SiteService) VerifyDomain(ctx context.Context, actorUserID, siteID, domainID string) (*store.SiteDomain, error) {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper); err != nil {
		return nil, err
	}
	domain, err := s.store.GetSiteDomain(ctx, domainID)
	if err != nil {
		return nil, err
	}
	if domain.SiteID != siteID {
		return nil, ErrPermissionDenied
	}
	domain, err = s.store.VerifySiteDomain(ctx, domainID)
	if err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: actorUserID, Action: "site.domain.verify", ResourceType: "site_domain", ResourceID: domain.ID, MetadataJSON: fmt.Sprintf(`{"hostname":%q,"mode":"manual-placeholder"}`, domain.Hostname)})
	return domain, nil
}

func (s *SiteService) DeleteDomain(ctx context.Context, actorUserID, siteID, domainID string) error {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return err
	}
	if err := s.ensureSiteRole(ctx, actorUserID, site.OrgID, store.RoleOrgOwner, store.RoleOrgDeveloper); err != nil {
		return err
	}
	domain, err := s.store.GetSiteDomain(ctx, domainID)
	if err != nil {
		return err
	}
	if domain.SiteID != siteID {
		return ErrPermissionDenied
	}
	if err := s.store.DeleteSiteDomain(ctx, domainID); err != nil {
		return err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: actorUserID, Action: "site.domain.delete", ResourceType: "site_domain", ResourceID: domain.ID, MetadataJSON: fmt.Sprintf(`{"hostname":%q}`, domain.Hostname)})
	return nil
}

func (s *SiteService) ensureSiteRole(ctx context.Context, userID, orgID string, allowed ...string) error {
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

var ErrPermissionDenied = errors.New("permission denied")
