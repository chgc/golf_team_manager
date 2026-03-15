package auth

import (
	"errors"
	"testing"
	"time"
)

func TestHMACTokenManagerRoundTrip(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	manager := &HMACTokenManager{
		secret: []byte("test-secret"),
		now:    func() time.Time { return now },
	}

	token, err := manager.Sign(Claims{
		Subject:         "user-1",
		PlayerID:        "player-1",
		Provider:        ProviderLINEOAuth,
		ProviderSubject: "line-user-1",
		Role:            RolePlayer,
		DisplayName:     "王小明",
		IssuedAt:        now.Unix(),
		ExpiresAt:       now.Add(time.Hour).Unix(),
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	claims, err := manager.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	principal := claims.Principal()
	if principal.UserID != "user-1" {
		t.Fatalf("principal.UserID = %q, want %q", principal.UserID, "user-1")
	}

	if principal.PlayerID != "player-1" {
		t.Fatalf("principal.PlayerID = %q, want %q", principal.PlayerID, "player-1")
	}

	if principal.Subject != "line-user-1" {
		t.Fatalf("principal.Subject = %q, want %q", principal.Subject, "line-user-1")
	}
}

func TestHMACTokenManagerRejectsExpiredToken(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	manager := &HMACTokenManager{
		secret: []byte("test-secret"),
		now:    func() time.Time { return now.Add(2 * time.Hour) },
	}

	token, err := manager.Sign(Claims{
		Subject:         "user-1",
		Provider:        ProviderLINEOAuth,
		ProviderSubject: "line-user-1",
		Role:            RolePlayer,
		DisplayName:     "王小明",
		IssuedAt:        now.Unix(),
		ExpiresAt:       now.Add(time.Hour).Unix(),
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	_, err = manager.Validate(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("Validate() error = %v, want %v", err, ErrExpiredToken)
	}
}
