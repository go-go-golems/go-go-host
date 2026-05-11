package store

import "time"

type User struct {
	ID          string
	Issuer      string
	Subject     string
	Email       string
	DisplayName string
	CreatedAt   time.Time
	LastLoginAt time.Time
}

type Org struct {
	ID        string
	Slug      string
	Name      string
	CreatedAt time.Time
}

type Membership struct {
	OrgID     string
	UserID    string
	Role      string
	CreatedAt time.Time
}

type OrgMembership struct {
	OrgID     string
	OrgSlug   string
	OrgName   string
	Role      string
	CreatedAt time.Time
}

type Site struct {
	ID                 string
	OrgID              string
	Slug               string
	Name               string
	PrimaryHost        string
	Status             string
	ActiveDeploymentID string
	CreatedAt          time.Time
}

type SiteQuota struct {
	SiteID           string
	BundleMaxBytes   int64
	DBSoftMaxBytes   int64
	DBHardMaxBytes   int64
	RequestTimeoutMS int
	UpdatedAt        time.Time
}

type RuntimeStatus struct {
	SiteID        string
	OrgID         string
	DeploymentID  string
	Hosts         []string
	Status        string
	StartedAt     time.Time
	LastError     string
	RequestsTotal int64
	ErrorsTotal   int64
	UpdatedAt     time.Time
}

type Agent struct {
	ID              string
	OrgID           string
	Name            string
	Status          string
	CreatedByUserID string
	CreatedAt       time.Time
	LastSeenAt      time.Time
}

type AgentSiteGrant struct {
	AgentID         string
	SiteID          string
	CanDeploy       bool
	CanRollback     bool
	AllowedChannels []string
	AllowedPaths    []string
	ExpiresAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type AuditEvent struct {
	ID           string
	OrgID        string
	ActorType    string
	ActorID      string
	Action       string
	ResourceType string
	ResourceID   string
	IPAddress    string
	UserAgent    string
	MetadataJSON string
	CreatedAt    time.Time
}

const (
	RoleOrgOwner     = "org_owner"
	RoleOrgDeveloper = "org_developer"
	RoleOrgViewer    = "org_viewer"

	SiteStatusProvisioning = "provisioning"
	SiteStatusActive       = "active"
	SiteStatusSuspended    = "suspended"

	AgentStatusActive  = "active"
	AgentStatusRevoked = "revoked"
)
