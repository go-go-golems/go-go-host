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

type SiteCapability struct {
	SiteID     string
	Capability string
	Enabled    bool
	ConfigJSON []byte
	UpdatedAt  time.Time
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
	CanActivate     bool
	AllowedChannels []string
	AllowedPaths    []string
	ExpiresAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type AgentEnrollmentToken struct {
	TokenHash string
	AgentID   string
	OrgID     string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UsedAt    time.Time
}

type AgentKey struct {
	ID         string
	AgentID    string
	PublicKey  string
	Status     string
	CreatedAt  time.Time
	RevokedAt  time.Time
	LastUsedAt time.Time
}

type DeployRun struct {
	ID                string
	SiteID            string
	ActorType         string
	ActorID           string
	AgentID           string
	RequestedByUserID string
	Status            string
	AllowedActions    []string
	AllowedChannels   []string
	AllowedPaths      []string
	UploadTokenHash   string
	ExpiresAt         time.Time
	CreatedAt         time.Time
	FinishedAt        time.Time
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

	AgentKeyStatusActive  = "active"
	AgentKeyStatusRevoked = "revoked"

	AgentEnrollmentTokenStatusActive = "active"
	AgentEnrollmentTokenStatusUsed   = "used"

	DeployRunStatusPending   = "pending"
	DeployRunStatusUploading = "uploading"
	DeployRunStatusCompleted = "completed"
	DeployRunStatusRejected  = "rejected"
)
