package store

import (
	"context"
	"errors"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

func (s *Store) UpsertUserFromOIDC(ctx context.Context, issuer, subject, email, displayName string) (*User, error) {
	if issuer == "" || subject == "" {
		return nil, errors.New("issuer and subject are required")
	}
	row, err := s.q.UpsertUserFromOIDC(ctx, storedb.UpsertUserFromOIDCParams{
		ID:          newID("usr"),
		Issuer:      issuer,
		Subject:     subject,
		Email:       email,
		DisplayName: displayName,
		CreatedAt:   pgTime(now()),
	})
	if err != nil {
		return nil, err
	}
	return &User{ID: row.ID, Issuer: row.Issuer, Subject: row.Subject, Email: row.Email, DisplayName: row.DisplayName, CreatedAt: fromPgTime(row.CreatedAt), LastLoginAt: fromPgTime(row.LastLoginAt)}, nil
}

func (s *Store) GetUser(ctx context.Context, id string) (*User, error) {
	row, err := s.q.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &User{ID: row.ID, Issuer: row.Issuer, Subject: row.Subject, Email: row.Email, DisplayName: row.DisplayName, CreatedAt: fromPgTime(row.CreatedAt), LastLoginAt: fromPgTime(row.LastLoginAt)}, nil
}

func (s *Store) AddPlatformAdmin(ctx context.Context, userID string) error {
	return s.q.AddPlatformAdmin(ctx, storedb.AddPlatformAdminParams{UserID: userID, CreatedAt: pgTime(now())})
}

func (s *Store) IsPlatformAdmin(ctx context.Context, userID string) (bool, error) {
	return s.q.IsPlatformAdmin(ctx, userID)
}
