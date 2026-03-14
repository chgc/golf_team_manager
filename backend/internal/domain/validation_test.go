package domain

import (
	"testing"
	"time"
)

func TestValidatePlayerWriteDTOAcceptsValidPayload(t *testing.T) {
	err := ValidatePlayerWriteDTO(PlayerWriteDTO{
		Name:     "王大明",
		Handicap: 12.5,
		Email:    "wang@example.com",
		Status:   PlayerStatusActive,
	})
	if err != nil {
		t.Fatalf("ValidatePlayerWriteDTO() error = %v", err)
	}
}

func TestValidatePlayerWriteDTORejectsInvalidValues(t *testing.T) {
	err := ValidatePlayerWriteDTO(PlayerWriteDTO{
		Name:     " ",
		Handicap: 12.3,
		Email:    "not-an-email",
		Status:   PlayerStatus("unknown"),
	})
	if err == nil {
		t.Fatal("ValidatePlayerWriteDTO() error = nil, want error")
	}
}

func TestValidateSessionWriteDTOAcceptsValidPayload(t *testing.T) {
	sessionDate := time.Date(2026, time.April, 5, 8, 0, 0, 0, time.UTC)
	deadline := time.Date(2026, time.March, 29, 23, 59, 0, 0, time.UTC)

	err := ValidateSessionWriteDTO(SessionWriteDTO{
		Date:                 sessionDate,
		CourseName:           "台北高爾夫球場",
		MaxPlayers:           12,
		RegistrationDeadline: deadline,
		Status:               SessionStatusOpen,
	})
	if err != nil {
		t.Fatalf("ValidateSessionWriteDTO() error = %v", err)
	}
}

func TestValidateSessionWriteDTORejectsInvalidValues(t *testing.T) {
	sessionDate := time.Date(2026, time.April, 5, 8, 0, 0, 0, time.UTC)
	deadline := sessionDate.Add(24 * time.Hour)

	err := ValidateSessionWriteDTO(SessionWriteDTO{
		Date:                 sessionDate,
		CourseName:           "",
		MaxPlayers:           0,
		RegistrationDeadline: deadline,
		Status:               SessionStatus("unknown"),
	})
	if err == nil {
		t.Fatal("ValidateSessionWriteDTO() error = nil, want error")
	}
}

func TestValidateRegistrationWriteDTOAcceptsValidPayload(t *testing.T) {
	err := ValidateRegistrationWriteDTO(RegistrationWriteDTO{
		PlayerID:  "player-1",
		SessionID: "session-1",
		Status:    RegistrationStatusConfirmed,
	})
	if err != nil {
		t.Fatalf("ValidateRegistrationWriteDTO() error = %v", err)
	}
}

func TestValidateRegistrationWriteDTORejectsInvalidValues(t *testing.T) {
	err := ValidateRegistrationWriteDTO(RegistrationWriteDTO{
		PlayerID:  "",
		SessionID: "",
		Status:    RegistrationStatus("unknown"),
	})
	if err == nil {
		t.Fatal("ValidateRegistrationWriteDTO() error = nil, want error")
	}
}
