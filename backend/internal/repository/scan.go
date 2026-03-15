package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/domain"
)

type scanner interface {
	Scan(dest ...any) error
}

func scanPlayer(scanTarget scanner) (domain.Player, error) {
	var (
		player       domain.Player
		status       string
		createdAtRaw string
		updatedAtRaw string
	)

	if err := scanTarget.Scan(
		&player.ID,
		&player.Name,
		&player.Handicap,
		&player.Phone,
		&player.Email,
		&status,
		&player.Notes,
		&createdAtRaw,
		&updatedAtRaw,
	); err != nil {
		return domain.Player{}, err
	}

	createdAt, err := parseTimestamp(createdAtRaw)
	if err != nil {
		return domain.Player{}, err
	}

	updatedAt, err := parseTimestamp(updatedAtRaw)
	if err != nil {
		return domain.Player{}, err
	}

	player.Status = domain.PlayerStatus(status)
	player.CreatedAt = createdAt
	player.UpdatedAt = updatedAt

	return player, nil
}

func scanSession(scanTarget scanner) (domain.Session, error) {
	var (
		session                 domain.Session
		status                  string
		sessionDateRaw          string
		registrationDeadlineRaw string
		createdAtRaw            string
		updatedAtRaw            string
	)

	if err := scanTarget.Scan(
		&session.ID,
		&sessionDateRaw,
		&session.CourseName,
		&session.CourseAddress,
		&session.MaxPlayers,
		&registrationDeadlineRaw,
		&status,
		&session.Notes,
		&createdAtRaw,
		&updatedAtRaw,
	); err != nil {
		return domain.Session{}, err
	}

	sessionDate, err := parseTimestamp(sessionDateRaw)
	if err != nil {
		return domain.Session{}, err
	}

	registrationDeadline, err := parseTimestamp(registrationDeadlineRaw)
	if err != nil {
		return domain.Session{}, err
	}

	createdAt, err := parseTimestamp(createdAtRaw)
	if err != nil {
		return domain.Session{}, err
	}

	updatedAt, err := parseTimestamp(updatedAtRaw)
	if err != nil {
		return domain.Session{}, err
	}

	session.Date = sessionDate
	session.RegistrationDeadline = registrationDeadline
	session.Status = domain.SessionStatus(status)
	session.CreatedAt = createdAt
	session.UpdatedAt = updatedAt

	return session, nil
}

func scanRegistration(scanTarget scanner) (domain.Registration, error) {
	var (
		registration    domain.Registration
		status          string
		registeredAtRaw string
		updatedAtRaw    string
	)

	if err := scanTarget.Scan(
		&registration.ID,
		&registration.PlayerID,
		&registration.SessionID,
		&status,
		&registeredAtRaw,
		&updatedAtRaw,
	); err != nil {
		return domain.Registration{}, err
	}

	registeredAt, err := parseTimestamp(registeredAtRaw)
	if err != nil {
		return domain.Registration{}, err
	}

	updatedAt, err := parseTimestamp(updatedAtRaw)
	if err != nil {
		return domain.Registration{}, err
	}

	registration.Status = domain.RegistrationStatus(status)
	registration.RegisteredAt = registeredAt
	registration.UpdatedAt = updatedAt

	return registration, nil
}

func scanUser(scanTarget scanner) (auth.User, error) {
	var (
		user         auth.User
		playerID     sql.NullString
		role         string
		provider     string
		createdAtRaw string
		updatedAtRaw string
	)

	if err := scanTarget.Scan(
		&user.ID,
		&playerID,
		&user.DisplayName,
		&role,
		&provider,
		&user.Subject,
		&createdAtRaw,
		&updatedAtRaw,
	); err != nil {
		return auth.User{}, err
	}

	createdAt, err := parseTimestamp(createdAtRaw)
	if err != nil {
		return auth.User{}, err
	}

	updatedAt, err := parseTimestamp(updatedAtRaw)
	if err != nil {
		return auth.User{}, err
	}

	if playerID.Valid {
		user.PlayerID = playerID.String
	}

	user.Role = auth.Role(role)
	user.Provider = auth.Provider(provider)
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return user, nil
}

func formatTimestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func parseTimestamp(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("parse timestamp %q", value)
}
