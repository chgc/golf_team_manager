package domain

import "time"

type SessionStatus string

const (
	SessionStatusOpen      SessionStatus = "open"
	SessionStatusClosed    SessionStatus = "closed"
	SessionStatusConfirmed SessionStatus = "confirmed"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

type Session struct {
	ID                   string
	Date                 time.Time
	CourseName           string
	CourseAddress        string
	MaxPlayers           int
	RegistrationDeadline time.Time
	Status               SessionStatus
	Notes                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type SessionWriteDTO struct {
	Date                 time.Time     `json:"date"`
	CourseName           string        `json:"courseName"`
	CourseAddress        string        `json:"courseAddress,omitempty"`
	MaxPlayers           int           `json:"maxPlayers"`
	RegistrationDeadline time.Time     `json:"registrationDeadline"`
	Status               SessionStatus `json:"status"`
	Notes                string        `json:"notes,omitempty"`
}

type SessionReadDTO struct {
	ID                   string        `json:"id"`
	Date                 time.Time     `json:"date"`
	CourseName           string        `json:"courseName"`
	CourseAddress        string        `json:"courseAddress,omitempty"`
	MaxPlayers           int           `json:"maxPlayers"`
	RegistrationDeadline time.Time     `json:"registrationDeadline"`
	Status               SessionStatus `json:"status"`
	Notes                string        `json:"notes,omitempty"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
}
