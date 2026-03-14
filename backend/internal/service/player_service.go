package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
)

type PlayerService struct {
	repository repository.PlayerRepository
}

func NewPlayerService(repository repository.PlayerRepository) *PlayerService {
	return &PlayerService{repository: repository}
}

func (s *PlayerService) Create(ctx context.Context, input domain.PlayerWriteDTO) (domain.PlayerReadDTO, error) {
	if err := domain.ValidatePlayerWriteDTO(input); err != nil {
		return domain.PlayerReadDTO{}, err
	}

	player, err := s.repository.Create(ctx, input)
	if err != nil {
		return domain.PlayerReadDTO{}, err
	}

	return mapPlayer(player), nil
}

func (s *PlayerService) GetByID(ctx context.Context, playerID string) (domain.PlayerReadDTO, error) {
	player, err := s.repository.GetByID(ctx, playerID)
	if err != nil {
		return domain.PlayerReadDTO{}, err
	}

	return mapPlayer(player), nil
}

func (s *PlayerService) List(
	ctx context.Context,
	query string,
	status string,
) ([]domain.PlayerReadDTO, error) {
	filter := repository.PlayerListFilter{Query: query}
	if status != "" {
		switch normalizedStatus := domain.PlayerStatus(strings.TrimSpace(status)); normalizedStatus {
		case domain.PlayerStatusActive, domain.PlayerStatusInactive:
			filter.Status = normalizedStatus
		default:
			return nil, domain.ValidationErrors{fmt.Errorf("invalid player status filter %q", status)}
		}
	}

	players, err := s.repository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := make([]domain.PlayerReadDTO, 0, len(players))
	for _, player := range players {
		results = append(results, mapPlayer(player))
	}

	return results, nil
}

func (s *PlayerService) Update(
	ctx context.Context,
	playerID string,
	input domain.PlayerWriteDTO,
) (domain.PlayerReadDTO, error) {
	if err := domain.ValidatePlayerWriteDTO(input); err != nil {
		return domain.PlayerReadDTO{}, err
	}

	player, err := s.repository.Update(ctx, playerID, input)
	if err != nil {
		return domain.PlayerReadDTO{}, err
	}

	return mapPlayer(player), nil
}

func mapPlayer(player domain.Player) domain.PlayerReadDTO {
	return domain.PlayerReadDTO{
		ID:        player.ID,
		Name:      player.Name,
		Handicap:  player.Handicap,
		Phone:     player.Phone,
		Email:     player.Email,
		Status:    player.Status,
		Notes:     player.Notes,
		CreatedAt: player.CreatedAt,
		UpdatedAt: player.UpdatedAt,
	}
}
