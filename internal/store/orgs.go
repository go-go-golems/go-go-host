package store

import (
	"context"
	"errors"
	"fmt"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateOrg(ctx context.Context, slug, name string) (*Org, error) {
	if slug == "" || name == "" {
		return nil, errors.New("slug and name are required")
	}
	row, err := s.q.CreateOrg(ctx, storedb.CreateOrgParams{ID: newID("org"), Slug: slug, Name: name, CreatedAt: pgTime(now())})
	if err != nil {
		return nil, fmt.Errorf("insert org: %w", err)
	}
	return &Org{ID: row.ID, Slug: row.Slug, Name: row.Name, CreatedAt: fromPgTime(row.CreatedAt)}, nil
}

func (s *Store) GetOrg(ctx context.Context, id string) (*Org, error) {
	row, err := s.q.GetOrg(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Org{ID: row.ID, Slug: row.Slug, Name: row.Name, CreatedAt: fromPgTime(row.CreatedAt)}, nil
}

func (s *Store) AddMembership(ctx context.Context, orgID, userID, role string) error {
	if role != RoleOrgOwner && role != RoleOrgDeveloper && role != RoleOrgViewer {
		return fmt.Errorf("invalid role %q", role)
	}
	return s.q.AddMembership(ctx, storedb.AddMembershipParams{OrgID: orgID, UserID: userID, Role: role, CreatedAt: pgTime(now())})
}

func (s *Store) MembershipRole(ctx context.Context, orgID, userID string) (string, error) {
	role, err := s.q.GetMembershipRole(ctx, storedb.GetMembershipRoleParams{OrgID: orgID, UserID: userID})
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return role, err
}

func (s *Store) ListOrgsForUser(ctx context.Context, userID string) ([]Org, error) {
	rows, err := s.q.ListOrgsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	orgs := make([]Org, 0, len(rows))
	for _, row := range rows {
		orgs = append(orgs, Org{ID: row.ID, Slug: row.Slug, Name: row.Name, CreatedAt: fromPgTime(row.CreatedAt)})
	}
	return orgs, nil
}

func (s *Store) ListMembershipsForUser(ctx context.Context, userID string) ([]OrgMembership, error) {
	rows, err := s.q.ListMembershipsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	memberships := make([]OrgMembership, 0, len(rows))
	for _, row := range rows {
		memberships = append(memberships, OrgMembership{OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Role: row.Role, CreatedAt: fromPgTime(row.CreatedAt)})
	}
	return memberships, nil
}
