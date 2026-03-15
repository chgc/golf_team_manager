package auth

import "time"

type User struct {
	ID          string
	PlayerID    string
	DisplayName string
	Provider    Provider
	Role        Role
	Subject     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u User) Principal() Principal {
	return Principal{
		DisplayName: u.DisplayName,
		PlayerID:    u.PlayerID,
		Provider:    u.Provider,
		Role:        u.Role,
		Subject:     u.Subject,
		UserID:      u.ID,
	}
}

func (u User) Claims(issuedAt time.Time, ttl time.Duration) Claims {
	return Claims{
		Subject:         u.ID,
		PlayerID:        u.PlayerID,
		Provider:        u.Provider,
		ProviderSubject: u.Subject,
		Role:            u.Role,
		DisplayName:     u.DisplayName,
		IssuedAt:        issuedAt.Unix(),
		ExpiresAt:       issuedAt.Add(ttl).Unix(),
	}
}
