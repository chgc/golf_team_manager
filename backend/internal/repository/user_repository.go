package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/google/uuid"
)

type UserRepository interface {
	UpsertLineUser(ctx context.Context, input auth.UpsertLineUserInput) (auth.User, error)
}

type SQLiteUserRepository struct {
	database *sql.DB
}

func NewSQLiteUserRepository(database *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{database: database}
}

func (r *SQLiteUserRepository) UpsertLineUser(ctx context.Context, input auth.UpsertLineUserInput) (auth.User, error) {
	existing, err := r.getByProviderSubject(ctx, auth.ProviderLINEOAuth, input.Subject)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return auth.User{}, err
	}

	now := time.Now().UTC()
	if errors.Is(err, ErrNotFound) {
		user := auth.User{
			ID:          uuid.NewString(),
			DisplayName: input.DisplayName,
			Provider:    auth.ProviderLINEOAuth,
			Role:        auth.RolePlayer,
			Subject:     input.Subject,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		_, err = r.database.ExecContext(
			ctx,
			`INSERT INTO users (id, player_id, display_name, role, auth_provider, provider_subject, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			user.ID,
			nil,
			user.DisplayName,
			user.Role,
			user.Provider,
			user.Subject,
			formatTimestamp(user.CreatedAt),
			formatTimestamp(user.UpdatedAt),
		)
		if err != nil {
			return auth.User{}, fmt.Errorf("insert user: %w", err)
		}

		return user, nil
	}

	existing.DisplayName = input.DisplayName
	existing.UpdatedAt = now

	_, err = r.database.ExecContext(
		ctx,
		`UPDATE users
		SET display_name = ?, updated_at = ?
		WHERE id = ?`,
		existing.DisplayName,
		formatTimestamp(existing.UpdatedAt),
		existing.ID,
	)
	if err != nil {
		return auth.User{}, fmt.Errorf("update user: %w", err)
	}

	return existing, nil
}

func (r *SQLiteUserRepository) getByProviderSubject(ctx context.Context, provider auth.Provider, subject string) (auth.User, error) {
	row := r.database.QueryRowContext(
		ctx,
		`SELECT id, player_id, display_name, role, auth_provider, provider_subject, created_at, updated_at
		FROM users
		WHERE auth_provider = ? AND provider_subject = ?`,
		provider,
		subject,
	)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, ErrNotFound
		}

		return auth.User{}, fmt.Errorf("select user by provider subject: %w", err)
	}

	return user, nil
}
