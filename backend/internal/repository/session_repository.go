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

type SessionRepository interface {
	Create(ctx context.Context, input domain.SessionWriteDTO) (domain.Session, error)
	GetByID(ctx context.Context, sessionID string) (domain.Session, error)
	List(ctx context.Context) ([]domain.Session, error)
}

type SQLiteSessionRepository struct {
	database *sql.DB
}

func NewSQLiteSessionRepository(database *sql.DB) *SQLiteSessionRepository {
	return &SQLiteSessionRepository{database: database}
}

func (r *SQLiteSessionRepository) Create(ctx context.Context, input domain.SessionWriteDTO) (domain.Session, error) {
	now := time.Now().UTC()
	session := domain.Session{
		ID:                   uuid.NewString(),
		Date:                 input.Date.UTC(),
		CourseName:           input.CourseName,
		CourseAddress:        input.CourseAddress,
		MaxPlayers:           input.MaxPlayers,
		RegistrationDeadline: input.RegistrationDeadline.UTC(),
		Status:               input.Status,
		Notes:                input.Notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	_, err := r.database.ExecContext(
		ctx,
		`INSERT INTO sessions (id, session_date, course_name, course_address, max_players, registration_deadline, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		session.ID,
		formatTimestamp(session.Date),
		session.CourseName,
		session.CourseAddress,
		session.MaxPlayers,
		formatTimestamp(session.RegistrationDeadline),
		session.Status,
		session.Notes,
		formatTimestamp(session.CreatedAt),
		formatTimestamp(session.UpdatedAt),
	)
	if err != nil {
		return domain.Session{}, fmt.Errorf("insert session: %w", err)
	}

	return session, nil
}

func (r *SQLiteSessionRepository) GetByID(ctx context.Context, sessionID string) (domain.Session, error) {
	row := r.database.QueryRowContext(
		ctx,
		`SELECT id, session_date, course_name, course_address, max_players, registration_deadline, status, notes, created_at, updated_at
		FROM sessions
		WHERE id = ?`,
		sessionID,
	)

	session, err := scanSession(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Session{}, ErrNotFound
		}

		return domain.Session{}, fmt.Errorf("select session by id: %w", err)
	}

	return session, nil
}

func (r *SQLiteSessionRepository) List(ctx context.Context) ([]domain.Session, error) {
	rows, err := r.database.QueryContext(
		ctx,
		`SELECT id, session_date, course_name, course_address, max_players, registration_deadline, status, notes, created_at, updated_at
		FROM sessions
		ORDER BY session_date ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]domain.Session, 0)
	for rows.Next() {
		session, err := scanSession(rows)
		if err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}

	return sessions, nil
}
