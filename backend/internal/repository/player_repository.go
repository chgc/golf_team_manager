package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/google/uuid"
)

type PlayerRepository interface {
	Create(ctx context.Context, input domain.PlayerWriteDTO) (domain.Player, error)
	GetByID(ctx context.Context, playerID string) (domain.Player, error)
	List(ctx context.Context, filter PlayerListFilter) ([]domain.Player, error)
	Update(ctx context.Context, playerID string, input domain.PlayerWriteDTO) (domain.Player, error)
}

type PlayerListFilter struct {
	Query  string
	Status domain.PlayerStatus
}

type SQLitePlayerRepository struct {
	database *sql.DB
}

func NewSQLitePlayerRepository(database *sql.DB) *SQLitePlayerRepository {
	return &SQLitePlayerRepository{database: database}
}

func (r *SQLitePlayerRepository) Create(ctx context.Context, input domain.PlayerWriteDTO) (domain.Player, error) {
	now := time.Now().UTC()
	player := domain.Player{
		ID:        uuid.NewString(),
		Name:      input.Name,
		Handicap:  input.Handicap,
		Phone:     input.Phone,
		Email:     input.Email,
		Status:    input.Status,
		Notes:     input.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := r.database.ExecContext(
		ctx,
		`INSERT INTO players (id, name, handicap, phone, email, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		player.ID,
		player.Name,
		player.Handicap,
		player.Phone,
		player.Email,
		player.Status,
		player.Notes,
		formatTimestamp(player.CreatedAt),
		formatTimestamp(player.UpdatedAt),
	)
	if err != nil {
		return domain.Player{}, fmt.Errorf("insert player: %w", err)
	}

	return player, nil
}

func (r *SQLitePlayerRepository) GetByID(ctx context.Context, playerID string) (domain.Player, error) {
	row := r.database.QueryRowContext(
		ctx,
		`SELECT id, name, handicap, phone, email, status, notes, created_at, updated_at
		FROM players
		WHERE id = ?`,
		playerID,
	)

	player, err := scanPlayer(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Player{}, ErrNotFound
		}

		return domain.Player{}, fmt.Errorf("select player by id: %w", err)
	}

	return player, nil
}

func (r *SQLitePlayerRepository) List(ctx context.Context, filter PlayerListFilter) ([]domain.Player, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(
		`SELECT id, name, handicap, phone, email, status, notes, created_at, updated_at
		FROM players`,
	)

	args := make([]any, 0, 2)
	conditions := make([]string, 0, 2)
	if strings.TrimSpace(filter.Query) != "" {
		conditions = append(conditions, "name LIKE ? COLLATE NOCASE")
		args = append(args, "%"+strings.TrimSpace(filter.Query)+"%")
	}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY created_at ASC")

	rows, err := r.database.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}
	defer rows.Close()

	players := make([]domain.Player, 0)
	for rows.Next() {
		player, err := scanPlayer(rows)
		if err != nil {
			return nil, fmt.Errorf("scan player: %w", err)
		}

		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate players: %w", err)
	}

	return players, nil
}

func (r *SQLitePlayerRepository) Update(
	ctx context.Context,
	playerID string,
	input domain.PlayerWriteDTO,
) (domain.Player, error) {
	now := time.Now().UTC()
	result, err := r.database.ExecContext(
		ctx,
		`UPDATE players
		SET name = ?, handicap = ?, phone = ?, email = ?, status = ?, notes = ?, updated_at = ?
		WHERE id = ?`,
		input.Name,
		input.Handicap,
		input.Phone,
		input.Email,
		input.Status,
		input.Notes,
		formatTimestamp(now),
		playerID,
	)
	if err != nil {
		return domain.Player{}, fmt.Errorf("update player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.Player{}, fmt.Errorf("rows affected for player update: %w", err)
	}

	if rowsAffected == 0 {
		return domain.Player{}, ErrNotFound
	}

	return r.GetByID(ctx, playerID)
}
