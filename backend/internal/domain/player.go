package domain

import "time"

type PlayerStatus string

const (
	PlayerStatusActive   PlayerStatus = "active"
	PlayerStatusInactive PlayerStatus = "inactive"
)

type Player struct {
	ID        string
	Name      string
	Handicap  float64
	Phone     string
	Email     string
	Status    PlayerStatus
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PlayerWriteDTO struct {
	Name     string       `json:"name"`
	Handicap float64      `json:"handicap"`
	Phone    string       `json:"phone,omitempty"`
	Email    string       `json:"email,omitempty"`
	Status   PlayerStatus `json:"status"`
	Notes    string       `json:"notes,omitempty"`
}

type PlayerReadDTO struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Handicap  float64      `json:"handicap"`
	Phone     string       `json:"phone,omitempty"`
	Email     string       `json:"email,omitempty"`
	Status    PlayerStatus `json:"status"`
	Notes     string       `json:"notes,omitempty"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}
