package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/google/uuid"
)

type RegistrationRepository interface {
	CountConfirmedBySession(ctx context.Context, sessionID string) (int, error)
	Create(ctx context.Context, input domain.RegistrationWriteDTO) (domain.Registration, error)
	GetByID(ctx context.Context, registrationID string) (domain.Registration, error)
	ListBySession(ctx context.Context, sessionID string) ([]domain.Registration, error)
	UpdateStatus(ctx context.Context, registrationID string, status domain.RegistrationStatus) (domain.Registration, error)
}

type SQLiteRegistrationRepository struct {
	database *sql.DB
}

func NewSQLiteRegistrationRepository(database *sql.DB) *SQLiteRegistrationRepository {
	return &SQLiteRegistrationRepository{database: database}
}

func (r *SQLiteRegistrationRepository) CountConfirmedBySession(ctx context.Context, sessionID string) (int, error) {
	var count int
	if err := r.database.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM registrations WHERE session_id = ? AND status = ?`,
		sessionID,
		domain.RegistrationStatusConfirmed,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count confirmed registrations: %w", err)
	}

	return count, nil
}

func (r *SQLiteRegistrationRepository) Create(ctx context.Context, input domain.RegistrationWriteDTO) (domain.Registration, error) {
	now := time.Now().UTC()
	registration := domain.Registration{
		ID:           uuid.NewString(),
		PlayerID:     input.PlayerID,
		SessionID:    input.SessionID,
		Status:       input.Status,
		RegisteredAt: now,
		UpdatedAt:    now,
	}

	_, err := r.database.ExecContext(
		ctx,
		`INSERT INTO registrations (id, player_id, session_id, status, registered_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		registration.ID,
		registration.PlayerID,
		registration.SessionID,
		registration.Status,
		formatTimestamp(registration.RegisteredAt),
		formatTimestamp(registration.UpdatedAt),
	)
	if err != nil {
		if isSQLiteConstraintError(err) {
			return domain.Registration{}, ErrConflict
		}

		return domain.Registration{}, fmt.Errorf("insert registration: %w", err)
	}

	return registration, nil
}

func (r *SQLiteRegistrationRepository) ListBySession(ctx context.Context, sessionID string) ([]domain.Registration, error) {
	rows, err := r.database.QueryContext(
		ctx,
		`SELECT id, player_id, session_id, status, registered_at, updated_at
		FROM registrations
		WHERE session_id = ?
		ORDER BY registered_at ASC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("list registrations by session: %w", err)
	}
	defer rows.Close()

	registrations := make([]domain.Registration, 0)
	for rows.Next() {
		registration, err := scanRegistration(rows)
		if err != nil {
			return nil, fmt.Errorf("scan registration: %w", err)
		}

		registrations = append(registrations, registration)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate registrations: %w", err)
	}

	return registrations, nil
}

func (r *SQLiteRegistrationRepository) GetByID(ctx context.Context, registrationID string) (domain.Registration, error) {
	row := r.database.QueryRowContext(
		ctx,
		`SELECT id, player_id, session_id, status, registered_at, updated_at
		FROM registrations
		WHERE id = ?`,
		registrationID,
	)

	registration, err := scanRegistration(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Registration{}, ErrNotFound
		}

		return domain.Registration{}, fmt.Errorf("select registration by id: %w", err)
	}

	return registration, nil
}

func (r *SQLiteRegistrationRepository) UpdateStatus(
	ctx context.Context,
	registrationID string,
	status domain.RegistrationStatus,
) (domain.Registration, error) {
	now := time.Now().UTC()
	result, err := r.database.ExecContext(
		ctx,
		`UPDATE registrations
		SET status = ?, updated_at = ?
		WHERE id = ?`,
		status,
		formatTimestamp(now),
		registrationID,
	)
	if err != nil {
		return domain.Registration{}, fmt.Errorf("update registration status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.Registration{}, fmt.Errorf("rows affected for registration update: %w", err)
	}

	if rowsAffected == 0 {
		return domain.Registration{}, ErrNotFound
	}

	return r.GetByID(ctx, registrationID)
}
