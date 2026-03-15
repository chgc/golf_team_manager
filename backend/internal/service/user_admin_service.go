package service

import (
	"context"
	"errors"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
)

type UserAdminStore interface {
	CountByRole(ctx context.Context, role auth.Role) (int, error)
	GetByID(ctx context.Context, userID string) (auth.User, error)
	GetByProviderSubject(ctx context.Context, provider auth.Provider, subject string) (auth.User, error)
	List(ctx context.Context, filter repository.UserListFilter) ([]auth.User, error)
	UpdateRoleAndPlayer(ctx context.Context, userID string, role auth.Role, playerID *string) (auth.User, error)
}

type UserAdminPlayerStore interface {
	GetByID(ctx context.Context, playerID string) (domain.Player, error)
}

type UserAdminUpdateInput struct {
	ClearPlayerLink bool
	PlayerID        *string
	Role            *auth.Role
}

type UserAdminService struct {
	users   UserAdminStore
	players UserAdminPlayerStore
}

func NewUserAdminService(users UserAdminStore, players UserAdminPlayerStore) *UserAdminService {
	return &UserAdminService{
		users:   users,
		players: players,
	}
}

func (s *UserAdminService) GetByID(ctx context.Context, userID string) (auth.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return auth.User{}, ErrUserNotFound
		}

		return auth.User{}, err
	}

	return user, nil
}

func (s *UserAdminService) GetByProviderSubject(
	ctx context.Context,
	provider auth.Provider,
	subject string,
) (auth.User, error) {
	user, err := s.users.GetByProviderSubject(ctx, provider, subject)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return auth.User{}, ErrUserNotFound
		}

		return auth.User{}, err
	}

	return user, nil
}

func (s *UserAdminService) List(ctx context.Context, filter repository.UserListFilter) ([]auth.User, error) {
	return s.users.List(ctx, filter)
}

func (s *UserAdminService) Update(
	ctx context.Context,
	userID string,
	input UserAdminUpdateInput,
) (auth.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return auth.User{}, ErrUserNotFound
		}

		return auth.User{}, err
	}

	nextRole := user.Role
	if input.Role != nil {
		nextRole = *input.Role
	}

	nextPlayerID := user.PlayerID
	if input.ClearPlayerLink {
		nextPlayerID = ""
	}

	if input.PlayerID != nil {
		if _, err := s.players.GetByID(ctx, *input.PlayerID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return auth.User{}, ErrPlayerNotFound
			}

			return auth.User{}, err
		}

		nextPlayerID = *input.PlayerID
	}

	if user.Role == auth.RoleManager && nextRole != auth.RoleManager {
		managerCount, err := s.users.CountByRole(ctx, auth.RoleManager)
		if err != nil {
			return auth.User{}, err
		}

		if managerCount <= 1 {
			return auth.User{}, ErrLastManagerDemotionForbidden
		}
	}

	var nullablePlayerID *string
	if nextPlayerID != "" {
		nullablePlayerID = &nextPlayerID
	}

	updatedUser, err := s.users.UpdateRoleAndPlayer(ctx, userID, nextRole, nullablePlayerID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return auth.User{}, ErrUserNotFound
		case errors.Is(err, repository.ErrConflict):
			return auth.User{}, ErrPlayerAlreadyLinked
		default:
			return auth.User{}, err
		}
	}

	return updatedUser, nil
}
