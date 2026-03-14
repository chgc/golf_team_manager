package domain

import "time"

type RegistrationStatus string

const (
	RegistrationStatusConfirmed RegistrationStatus = "confirmed"
	RegistrationStatusCancelled RegistrationStatus = "cancelled"
)

type Registration struct {
	ID           string
	PlayerID     string
	SessionID    string
	Status       RegistrationStatus
	RegisteredAt time.Time
	UpdatedAt    time.Time
}

type RegistrationWriteDTO struct {
	PlayerID  string             `json:"playerId"`
	SessionID string             `json:"sessionId"`
	Status    RegistrationStatus `json:"status"`
}

type RegistrationStatusUpdateDTO struct {
	Status RegistrationStatus `json:"status"`
}

type RegistrationReadDTO struct {
	ID           string             `json:"id"`
	PlayerID     string             `json:"playerId"`
	SessionID    string             `json:"sessionId"`
	Status       RegistrationStatus `json:"status"`
	RegisteredAt time.Time          `json:"registeredAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}
