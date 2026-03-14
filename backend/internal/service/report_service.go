package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
)

type ReportService struct {
	playerRepository       repository.PlayerRepository
	registrationRepository repository.RegistrationRepository
	sessionRepository      repository.SessionRepository
}

func NewReportService(
	playerRepository repository.PlayerRepository,
	sessionRepository repository.SessionRepository,
	registrationRepository repository.RegistrationRepository,
) *ReportService {
	return &ReportService{
		playerRepository:       playerRepository,
		registrationRepository: registrationRepository,
		sessionRepository:      sessionRepository,
	}
}

func (s *ReportService) GetReservationSummary(
	ctx context.Context,
	sessionID string,
) (domain.ReservationSummaryReadDTO, error) {
	session, err := s.sessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.ReservationSummaryReadDTO{}, ErrSessionReportNotFound
		}

		return domain.ReservationSummaryReadDTO{}, err
	}

	if session.Status != domain.SessionStatusConfirmed && session.Status != domain.SessionStatusCompleted {
		return domain.ReservationSummaryReadDTO{}, ErrReservationSummaryNotEligible
	}

	registrations, err := s.registrationRepository.ListBySession(ctx, sessionID)
	if err != nil {
		return domain.ReservationSummaryReadDTO{}, err
	}

	players, err := s.playerRepository.List(ctx, repository.PlayerListFilter{})
	if err != nil {
		return domain.ReservationSummaryReadDTO{}, err
	}

	playerNamesByID := make(map[string]string, len(players))
	for _, player := range players {
		playerNamesByID[player.ID] = player.Name
	}

	confirmedPlayers := make([]domain.ReservationSummaryPlayerDTO, 0, len(registrations))
	for _, registration := range registrations {
		if registration.Status != domain.RegistrationStatusConfirmed {
			continue
		}

		playerName := registration.PlayerID
		if mappedName, ok := playerNamesByID[registration.PlayerID]; ok {
			playerName = mappedName
		}

		confirmedPlayers = append(confirmedPlayers, domain.ReservationSummaryPlayerDTO{
			PlayerID:   registration.PlayerID,
			PlayerName: playerName,
		})
	}

	if len(confirmedPlayers) == 0 {
		return domain.ReservationSummaryReadDTO{}, ErrReservationSummaryEmpty
	}

	sort.Slice(confirmedPlayers, func(left, right int) bool {
		if confirmedPlayers[left].PlayerName == confirmedPlayers[right].PlayerName {
			return confirmedPlayers[left].PlayerID < confirmedPlayers[right].PlayerID
		}

		return confirmedPlayers[left].PlayerName < confirmedPlayers[right].PlayerName
	})

	estimatedGroups := (len(confirmedPlayers) + 3) / 4

	return domain.ReservationSummaryReadDTO{
		SessionID:            session.ID,
		SessionDate:          session.Date,
		CourseName:           session.CourseName,
		CourseAddress:        session.CourseAddress,
		RegistrationDeadline: session.RegistrationDeadline,
		SessionStatus:        session.Status,
		ConfirmedPlayerCount: len(confirmedPlayers),
		EstimatedGroups:      estimatedGroups,
		SummaryText:          buildSummaryText(session, len(confirmedPlayers), estimatedGroups, confirmedPlayers),
		ConfirmedPlayers:     confirmedPlayers,
	}, nil
}

func buildSummaryText(
	session domain.Session,
	confirmedPlayerCount int,
	estimatedGroups int,
	confirmedPlayers []domain.ReservationSummaryPlayerDTO,
) string {
	addressLine := session.CourseAddress
	if strings.TrimSpace(addressLine) == "" {
		addressLine = "N/A"
	}

	lines := []string{
		fmt.Sprintf("Session: %s", formatSummaryTimestamp(session.Date)),
		fmt.Sprintf("Course: %s", session.CourseName),
		fmt.Sprintf("Address: %s", addressLine),
		fmt.Sprintf("Deadline: %s", formatSummaryTimestamp(session.RegistrationDeadline)),
		fmt.Sprintf("Status: %s", session.Status),
		fmt.Sprintf("Confirmed Players: %d", confirmedPlayerCount),
		fmt.Sprintf("Estimated Groups: %d", estimatedGroups),
		"Roster:",
	}

	for _, player := range confirmedPlayers {
		lines = append(lines, fmt.Sprintf("- %s", player.PlayerName))
	}

	return strings.Join(lines, "\n")
}

func formatSummaryTimestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}
