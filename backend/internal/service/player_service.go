package service

import (
	"context"

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

func (s *PlayerService) List(ctx context.Context) ([]domain.PlayerReadDTO, error) {
	players, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]domain.PlayerReadDTO, 0, len(players))
	for _, player := range players {
		results = append(results, mapPlayer(player))
	}

	return results, nil
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
