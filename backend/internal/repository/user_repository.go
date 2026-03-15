package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/google/uuid"
)

type UserRepository interface {
	CountByRole(ctx context.Context, role auth.Role) (int, error)
	GetByID(ctx context.Context, userID string) (auth.User, error)
	GetByProviderSubject(ctx context.Context, provider auth.Provider, subject string) (auth.User, error)
	List(ctx context.Context, filter UserListFilter) ([]auth.User, error)
	UpdateRoleAndPlayer(ctx context.Context, userID string, role auth.Role, playerID *string) (auth.User, error)
	UpsertLineUser(ctx context.Context, input auth.UpsertLineUserInput) (auth.User, error)
}

type UserLinkState string

const (
	UserLinkStateLinked   UserLinkState = "linked"
	UserLinkStateUnlinked UserLinkState = "unlinked"
)

type UserListFilter struct {
	LinkState UserLinkState
	Role      auth.Role
}

type SQLiteUserRepository struct {
	database *sql.DB
}

func NewSQLiteUserRepository(database *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{database: database}
}

func (r *SQLiteUserRepository) UpsertLineUser(ctx context.Context, input auth.UpsertLineUserInput) (auth.User, error) {
	existing, err := r.GetByProviderSubject(ctx, auth.ProviderLINEOAuth, input.Subject)
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

func (r *SQLiteUserRepository) CountByRole(ctx context.Context, role auth.Role) (int, error) {
	var count int
	if err := r.database.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM users WHERE role = ?`,
		role,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count users by role: %w", err)
	}

	return count, nil
}

func (r *SQLiteUserRepository) GetByID(ctx context.Context, userID string) (auth.User, error) {
	row := r.database.QueryRowContext(
		ctx,
		`SELECT id, player_id, display_name, role, auth_provider, provider_subject, created_at, updated_at
		FROM users
		WHERE id = ?`,
		userID,
	)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, ErrNotFound
		}

		return auth.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return user, nil
}

func (r *SQLiteUserRepository) GetByProviderSubject(ctx context.Context, provider auth.Provider, subject string) (auth.User, error) {
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

func (r *SQLiteUserRepository) List(ctx context.Context, filter UserListFilter) ([]auth.User, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(
		`SELECT id, player_id, display_name, role, auth_provider, provider_subject, created_at, updated_at
		FROM users`,
	)

	args := make([]any, 0, 2)
	conditions := make([]string, 0, 2)
	if filter.Role != "" {
		conditions = append(conditions, "role = ?")
		args = append(args, filter.Role)
	}

	switch filter.LinkState {
	case UserLinkStateLinked:
		conditions = append(conditions, "player_id IS NOT NULL")
	case UserLinkStateUnlinked:
		conditions = append(conditions, "player_id IS NULL")
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY created_at ASC")

	rows, err := r.database.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]auth.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *SQLiteUserRepository) UpdateRoleAndPlayer(
	ctx context.Context,
	userID string,
	role auth.Role,
	playerID *string,
) (auth.User, error) {
	now := time.Now().UTC()

	var nullablePlayerID any
	if playerID != nil {
		nullablePlayerID = *playerID
	}

	result, err := r.database.ExecContext(
		ctx,
		`UPDATE users
		SET player_id = ?, role = ?, updated_at = ?
		WHERE id = ?`,
		nullablePlayerID,
		role,
		formatTimestamp(now),
		userID,
	)
	if err != nil {
		if isSQLiteConstraintError(err) {
			return auth.User{}, ErrConflict
		}

		return auth.User{}, fmt.Errorf("update user role and player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return auth.User{}, fmt.Errorf("rows affected for user update: %w", err)
	}

	if rowsAffected == 0 {
		return auth.User{}, ErrNotFound
	}

	return r.GetByID(ctx, userID)
}
