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

type PlayerRepository interface {
	Create(ctx context.Context, input domain.PlayerWriteDTO) (domain.Player, error)
	GetByID(ctx context.Context, playerID string) (domain.Player, error)
	List(ctx context.Context) ([]domain.Player, error)
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

func (r *SQLitePlayerRepository) List(ctx context.Context) ([]domain.Player, error) {
	rows, err := r.database.QueryContext(
		ctx,
		`SELECT id, name, handicap, phone, email, status, notes, created_at, updated_at
		FROM players
		ORDER BY created_at ASC`,
	)
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
