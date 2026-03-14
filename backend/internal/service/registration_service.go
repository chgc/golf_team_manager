package service

import (
	"context"
	"errors"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
)

type RegistrationService struct {
	playerRepository       repository.PlayerRepository
	sessionRepository      repository.SessionRepository
	registrationRepository repository.RegistrationRepository
}

func NewRegistrationService(
	playerRepository repository.PlayerRepository,
	sessionRepository repository.SessionRepository,
	registrationRepository repository.RegistrationRepository,
) *RegistrationService {
	return &RegistrationService{
		playerRepository:       playerRepository,
		sessionRepository:      sessionRepository,
		registrationRepository: registrationRepository,
	}
}

func (s *RegistrationService) Create(
	ctx context.Context,
	input domain.RegistrationWriteDTO,
) (domain.RegistrationReadDTO, error) {
	if err := domain.ValidateRegistrationWriteDTO(input); err != nil {
		return domain.RegistrationReadDTO{}, err
	}

	player, err := s.playerRepository.GetByID(ctx, input.PlayerID)
	if err != nil {
		return domain.RegistrationReadDTO{}, err
	}

	if player.Status != domain.PlayerStatusActive {
		return domain.RegistrationReadDTO{}, ErrPlayerInactive
	}

	session, err := s.sessionRepository.GetByID(ctx, input.SessionID)
	if err != nil {
		return domain.RegistrationReadDTO{}, err
	}

	if session.Status != domain.SessionStatusOpen {
		return domain.RegistrationReadDTO{}, ErrSessionNotOpen
	}

	confirmedCount, err := s.registrationRepository.CountConfirmedBySession(ctx, input.SessionID)
	if err != nil {
		return domain.RegistrationReadDTO{}, err
	}

	if input.Status == domain.RegistrationStatusConfirmed && confirmedCount >= session.MaxPlayers {
		return domain.RegistrationReadDTO{}, ErrSessionCapacityFull
	}

	registration, err := s.registrationRepository.Create(ctx, input)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.RegistrationReadDTO{}, repository.ErrConflict
		}

		return domain.RegistrationReadDTO{}, err
	}

	return mapRegistration(registration), nil
}

func (s *RegistrationService) ListBySession(
	ctx context.Context,
	sessionID string,
) ([]domain.RegistrationReadDTO, error) {
	if _, err := s.sessionRepository.GetByID(ctx, sessionID); err != nil {
		return nil, err
	}

	registrations, err := s.registrationRepository.ListBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	results := make([]domain.RegistrationReadDTO, 0, len(registrations))
	for _, registration := range registrations {
		results = append(results, mapRegistration(registration))
	}

	return results, nil
}

func (s *RegistrationService) UpdateStatus(
	ctx context.Context,
	registrationID string,
	input domain.RegistrationStatusUpdateDTO,
) (domain.RegistrationReadDTO, error) {
	if err := domain.ValidateRegistrationStatusUpdateDTO(input); err != nil {
		return domain.RegistrationReadDTO{}, err
	}

	registration, err := s.registrationRepository.GetByID(ctx, registrationID)
	if err != nil {
		return domain.RegistrationReadDTO{}, err
	}

	if registration.Status == input.Status {
		return mapRegistration(registration), nil
	}

	switch {
	case registration.Status == domain.RegistrationStatusConfirmed && input.Status == domain.RegistrationStatusCancelled:
		updatedRegistration, err := s.registrationRepository.UpdateStatus(ctx, registrationID, input.Status)
		if err != nil {
			return domain.RegistrationReadDTO{}, err
		}

		return mapRegistration(updatedRegistration), nil
	case registration.Status == domain.RegistrationStatusCancelled && input.Status == domain.RegistrationStatusConfirmed:
		player, err := s.playerRepository.GetByID(ctx, registration.PlayerID)
		if err != nil {
			return domain.RegistrationReadDTO{}, err
		}

		if player.Status != domain.PlayerStatusActive {
			return domain.RegistrationReadDTO{}, ErrPlayerInactive
		}

		session, err := s.sessionRepository.GetByID(ctx, registration.SessionID)
		if err != nil {
			return domain.RegistrationReadDTO{}, err
		}

		if session.Status != domain.SessionStatusOpen {
			return domain.RegistrationReadDTO{}, ErrSessionNotOpen
		}

		confirmedCount, err := s.registrationRepository.CountConfirmedBySession(ctx, registration.SessionID)
		if err != nil {
			return domain.RegistrationReadDTO{}, err
		}

		if confirmedCount >= session.MaxPlayers {
			return domain.RegistrationReadDTO{}, ErrSessionCapacityFull
		}

		updatedRegistration, err := s.registrationRepository.UpdateStatus(ctx, registrationID, input.Status)
		if err != nil {
			return domain.RegistrationReadDTO{}, err
		}

		return mapRegistration(updatedRegistration), nil
	default:
		return domain.RegistrationReadDTO{}, domain.ValidationErrors{
			errors.New("invalid registration status transition"),
		}
	}
}

func mapRegistration(registration domain.Registration) domain.RegistrationReadDTO {
	return domain.RegistrationReadDTO{
		ID:           registration.ID,
		PlayerID:     registration.PlayerID,
		SessionID:    registration.SessionID,
		Status:       registration.Status,
		RegisteredAt: registration.RegisteredAt,
		UpdatedAt:    registration.UpdatedAt,
	}
}
