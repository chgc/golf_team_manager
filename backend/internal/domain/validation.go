package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type ValidationErrors []error

func (v ValidationErrors) Error() string {
	messages := make([]string, 0, len(v))
	for _, err := range v {
		messages = append(messages, err.Error())
	}

	return strings.Join(messages, "; ")
}

func (v ValidationErrors) AsError() error {
	if len(v) == 0 {
		return nil
	}

	return v
}

func ValidatePlayerWriteDTO(input PlayerWriteDTO) error {
	var validationErrors ValidationErrors

	if strings.TrimSpace(input.Name) == "" {
		validationErrors = append(validationErrors, errors.New("name is required"))
	}

	if input.Handicap < 0 || input.Handicap > 54 {
		validationErrors = append(validationErrors, errors.New("handicap must be between 0 and 54"))
	}

	if !isHalfStep(input.Handicap) {
		validationErrors = append(validationErrors, errors.New("handicap must use 0.5 increments"))
	}

	if !isValidPlayerStatus(input.Status) {
		validationErrors = append(validationErrors, fmt.Errorf("invalid player status %q", input.Status))
	}

	if input.Email != "" {
		if _, err := mail.ParseAddress(input.Email); err != nil {
			validationErrors = append(validationErrors, errors.New("email must be valid"))
		}
	}

	return validationErrors.AsError()
}

func ValidateSessionWriteDTO(input SessionWriteDTO) error {
	var validationErrors ValidationErrors

	if input.Date.IsZero() {
		validationErrors = append(validationErrors, errors.New("date is required"))
	}

	if strings.TrimSpace(input.CourseName) == "" {
		validationErrors = append(validationErrors, errors.New("courseName is required"))
	}

	if input.MaxPlayers <= 0 {
		validationErrors = append(validationErrors, errors.New("maxPlayers must be greater than zero"))
	}

	if input.RegistrationDeadline.IsZero() {
		validationErrors = append(validationErrors, errors.New("registrationDeadline is required"))
	}

	if !input.Date.IsZero() && !input.RegistrationDeadline.IsZero() &&
		input.RegistrationDeadline.After(input.Date) {
		validationErrors = append(validationErrors, errors.New("registrationDeadline must be on or before date"))
	}

	if !isValidSessionStatus(input.Status) {
		validationErrors = append(validationErrors, fmt.Errorf("invalid session status %q", input.Status))
	}

	return validationErrors.AsError()
}

func ValidateRegistrationWriteDTO(input RegistrationWriteDTO) error {
	var validationErrors ValidationErrors

	if strings.TrimSpace(input.PlayerID) == "" {
		validationErrors = append(validationErrors, errors.New("playerId is required"))
	}

	if strings.TrimSpace(input.SessionID) == "" {
		validationErrors = append(validationErrors, errors.New("sessionId is required"))
	}

	if !isValidRegistrationStatus(input.Status) {
		validationErrors = append(validationErrors, fmt.Errorf("invalid registration status %q", input.Status))
	}

	return validationErrors.AsError()
}

func isHalfStep(value float64) bool {
	doubled := value * 2
	return doubled == float64(int(doubled))
}

func isValidPlayerStatus(status PlayerStatus) bool {
	switch status {
	case PlayerStatusActive, PlayerStatusInactive:
		return true
	default:
		return false
	}
}

func isValidSessionStatus(status SessionStatus) bool {
	switch status {
	case SessionStatusOpen, SessionStatusClosed, SessionStatusConfirmed, SessionStatusCompleted, SessionStatusCancelled:
		return true
	default:
		return false
	}
}

func isValidRegistrationStatus(status RegistrationStatus) bool {
	switch status {
	case RegistrationStatusConfirmed, RegistrationStatusCancelled:
		return true
	default:
		return false
	}
}
