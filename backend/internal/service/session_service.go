package service

import (
	"context"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
)

type SessionService struct {
	repository repository.SessionRepository
}

func NewSessionService(repository repository.SessionRepository) *SessionService {
	return &SessionService{repository: repository}
}

func (s *SessionService) Create(ctx context.Context, input domain.SessionWriteDTO) (domain.SessionReadDTO, error) {
	if err := domain.ValidateSessionWriteDTO(input); err != nil {
		return domain.SessionReadDTO{}, err
	}

	session, err := s.repository.Create(ctx, input)
	if err != nil {
		return domain.SessionReadDTO{}, err
	}

	return mapSession(session), nil
}

func (s *SessionService) GetByID(ctx context.Context, sessionID string) (domain.Session, error) {
	return s.repository.GetByID(ctx, sessionID)
}

func (s *SessionService) List(ctx context.Context) ([]domain.SessionReadDTO, error) {
	sessions, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]domain.SessionReadDTO, 0, len(sessions))
	for _, session := range sessions {
		results = append(results, mapSession(session))
	}

	return results, nil
}

func mapSession(session domain.Session) domain.SessionReadDTO {
	return domain.SessionReadDTO{
		ID:                   session.ID,
		Date:                 session.Date,
		CourseName:           session.CourseName,
		CourseAddress:        session.CourseAddress,
		MaxPlayers:           session.MaxPlayers,
		RegistrationDeadline: session.RegistrationDeadline,
		Status:               session.Status,
		Notes:                session.Notes,
		CreatedAt:            session.CreatedAt,
		UpdatedAt:            session.UpdatedAt,
	}
}
