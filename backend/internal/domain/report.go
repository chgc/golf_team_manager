package domain

import "time"

type ReservationSummaryPlayerDTO struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

type ReservationSummaryReadDTO struct {
	SessionID            string                        `json:"sessionId"`
	SessionDate          time.Time                     `json:"sessionDate"`
	CourseName           string                        `json:"courseName"`
	CourseAddress        string                        `json:"courseAddress"`
	RegistrationDeadline time.Time                     `json:"registrationDeadline"`
	SessionStatus        SessionStatus                 `json:"sessionStatus"`
	ConfirmedPlayerCount int                           `json:"confirmedPlayerCount"`
	EstimatedGroups      int                           `json:"estimatedGroups"`
	SummaryText          string                        `json:"summaryText"`
	ConfirmedPlayers     []ReservationSummaryPlayerDTO `json:"confirmedPlayers"`
}
