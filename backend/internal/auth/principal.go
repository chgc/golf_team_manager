package auth

type Role string

const (
	RoleManager Role = "manager"
	RolePlayer  Role = "player"
)

type Provider string

const (
	ProviderLINEOAuth Provider = "line"
)

type Principal struct {
	DisplayName string   `json:"displayName"`
	PlayerID    string   `json:"playerId,omitempty"`
	Provider    Provider `json:"provider"`
	Role        Role     `json:"role"`
	Subject     string   `json:"subject"`
	UserID      string   `json:"userId"`
}
