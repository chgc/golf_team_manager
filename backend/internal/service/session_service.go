package service

import (
	"context"
	"fmt"
	"time"

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

func (s *SessionService) GetByID(ctx context.Context, sessionID string) (domain.SessionReadDTO, error) {
	session, err := s.repository.GetByID(ctx, sessionID)
	if err != nil {
		return domain.SessionReadDTO{}, err
	}

	reconciledSession, err := s.reconcileAutoClose(ctx, session)
	if err != nil {
		return domain.SessionReadDTO{}, err
	}

	return mapSession(reconciledSession), nil
}

func (s *SessionService) List(ctx context.Context) ([]domain.SessionReadDTO, error) {
	sessions, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]domain.SessionReadDTO, 0, len(sessions))
	for _, session := range sessions {
		reconciledSession, err := s.reconcileAutoClose(ctx, session)
		if err != nil {
			return nil, err
		}

		results = append(results, mapSession(reconciledSession))
	}

	return results, nil
}

func (s *SessionService) Update(
	ctx context.Context,
	sessionID string,
	input domain.SessionWriteDTO,
) (domain.SessionReadDTO, error) {
	if err := domain.ValidateSessionWriteDTO(input); err != nil {
		return domain.SessionReadDTO{}, err
	}

	currentSession, err := s.repository.GetByID(ctx, sessionID)
	if err != nil {
		return domain.SessionReadDTO{}, err
	}

	reconciledSession, err := s.reconcileAutoClose(ctx, currentSession)
	if err != nil {
		return domain.SessionReadDTO{}, err
	}

	if err := validateSessionStatusTransition(reconciledSession.Status, input.Status); err != nil {
		return domain.SessionReadDTO{}, err
	}

	session, err := s.repository.Update(ctx, sessionID, input)
	if err != nil {
		return domain.SessionReadDTO{}, err
	}

	return mapSession(session), nil
}

func (s *SessionService) reconcileAutoClose(ctx context.Context, session domain.Session) (domain.Session, error) {
	if session.Status != domain.SessionStatusOpen {
		return session, nil
	}

	now := time.Now().UTC()
	if !session.RegistrationDeadline.Before(now) {
		return session, nil
	}

	closedSession, err := s.repository.Update(ctx, session.ID, domain.SessionWriteDTO{
		Date:                 session.Date,
		CourseName:           session.CourseName,
		CourseAddress:        session.CourseAddress,
		MaxPlayers:           session.MaxPlayers,
		RegistrationDeadline: session.RegistrationDeadline,
		Status:               domain.SessionStatusClosed,
		Notes:                session.Notes,
	})
	if err != nil {
		return domain.Session{}, err
	}

	return closedSession, nil
}

func validateSessionStatusTransition(currentStatus domain.SessionStatus, nextStatus domain.SessionStatus) error {
	if currentStatus == nextStatus {
		return nil
	}

	switch currentStatus {
	case domain.SessionStatusOpen:
		if nextStatus == domain.SessionStatusClosed || nextStatus == domain.SessionStatusCancelled {
			return nil
		}
	case domain.SessionStatusClosed:
		if nextStatus == domain.SessionStatusConfirmed || nextStatus == domain.SessionStatusCancelled {
			return nil
		}
	case domain.SessionStatusConfirmed:
		if nextStatus == domain.SessionStatusCompleted || nextStatus == domain.SessionStatusCancelled {
			return nil
		}
	}

	return domain.ValidationErrors{
		fmt.Errorf("invalid session status transition %q -> %q", currentStatus, nextStatus),
	}
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
