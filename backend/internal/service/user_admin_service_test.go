package service

import (
	"context"
	"errors"
	"testing"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
)

func TestUserAdminServiceUpdateRejectsLastManagerDemotion(t *testing.T) {
	users := &stubUserAdminStore{
		getByIDResult: auth.User{
			ID:   "user-1",
			Role: auth.RoleManager,
		},
		countByRoleResult: 1,
	}
	service := NewUserAdminService(users, stubPlayerStore{})

	playerRole := auth.RolePlayer
	_, err := service.Update(context.Background(), "user-1", UserAdminUpdateInput{
		Role: &playerRole,
	})
	if !errors.Is(err, ErrLastManagerDemotionForbidden) {
		t.Fatalf("Update() error = %v, want %v", err, ErrLastManagerDemotionForbidden)
	}

	if users.updateCalled {
		t.Fatal("UpdateRoleAndPlayer() called, want false")
	}
}

func TestUserAdminServiceUpdateRejectsUnknownPlayer(t *testing.T) {
	users := &stubUserAdminStore{
		getByIDResult: auth.User{
			ID:   "user-1",
			Role: auth.RolePlayer,
		},
	}
	service := NewUserAdminService(users, stubPlayerStore{getByIDErr: repository.ErrNotFound})

	playerID := "missing-player"
	_, err := service.Update(context.Background(), "user-1", UserAdminUpdateInput{
		PlayerID: &playerID,
	})
	if !errors.Is(err, ErrPlayerNotFound) {
		t.Fatalf("Update() error = %v, want %v", err, ErrPlayerNotFound)
	}
}

func TestUserAdminServiceUpdateMapsPlayerAlreadyLinkedConflict(t *testing.T) {
	users := &stubUserAdminStore{
		getByIDResult: auth.User{
			ID:   "user-1",
			Role: auth.RolePlayer,
		},
		updateErr: repository.ErrConflict,
	}
	service := NewUserAdminService(users, stubPlayerStore{
		getByIDResult: domain.Player{ID: "player-1"},
	})

	playerID := "player-1"
	_, err := service.Update(context.Background(), "user-1", UserAdminUpdateInput{
		PlayerID: &playerID,
	})
	if !errors.Is(err, ErrPlayerAlreadyLinked) {
		t.Fatalf("Update() error = %v, want %v", err, ErrPlayerAlreadyLinked)
	}
}

func TestUserAdminServiceUpdateAppliesRoleAndPlayer(t *testing.T) {
	managerRole := auth.RoleManager
	playerID := "player-1"
	users := &stubUserAdminStore{
		getByIDResult: auth.User{
			ID:   "user-1",
			Role: auth.RolePlayer,
		},
		updateResult: auth.User{
			ID:       "user-1",
			Role:     auth.RoleManager,
			PlayerID: "player-1",
		},
	}
	service := NewUserAdminService(users, stubPlayerStore{
		getByIDResult: domain.Player{ID: "player-1"},
	})

	updatedUser, err := service.Update(context.Background(), "user-1", UserAdminUpdateInput{
		PlayerID: &playerID,
		Role:     &managerRole,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updatedUser.Role != auth.RoleManager {
		t.Fatalf("updatedUser.Role = %q, want %q", updatedUser.Role, auth.RoleManager)
	}

	if updatedUser.PlayerID != "player-1" {
		t.Fatalf("updatedUser.PlayerID = %q, want %q", updatedUser.PlayerID, "player-1")
	}

	if !users.updateCalled {
		t.Fatal("UpdateRoleAndPlayer() called = false, want true")
	}

	if users.updateRole != auth.RoleManager {
		t.Fatalf("updateRole = %q, want %q", users.updateRole, auth.RoleManager)
	}

	if users.updatePlayerID == nil || *users.updatePlayerID != "player-1" {
		t.Fatalf("updatePlayerID = %v, want %q", users.updatePlayerID, "player-1")
	}
}

type stubUserAdminStore struct {
	countByRoleErr    error
	countByRoleResult int
	getByIDErr        error
	getByIDResult     auth.User
	updateCalled      bool
	updateErr         error
	updatePlayerID    *string
	updateResult      auth.User
	updateRole        auth.Role
}

func (s *stubUserAdminStore) CountByRole(context.Context, auth.Role) (int, error) {
	return s.countByRoleResult, s.countByRoleErr
}

func (s *stubUserAdminStore) GetByID(context.Context, string) (auth.User, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *stubUserAdminStore) GetByProviderSubject(context.Context, auth.Provider, string) (auth.User, error) {
	return auth.User{}, nil
}

func (s *stubUserAdminStore) List(context.Context, repository.UserListFilter) ([]auth.User, error) {
	return nil, nil
}

func (s *stubUserAdminStore) UpdateRoleAndPlayer(
	_ context.Context,
	_ string,
	role auth.Role,
	playerID *string,
) (auth.User, error) {
	s.updateCalled = true
	s.updateRole = role
	s.updatePlayerID = playerID
	return s.updateResult, s.updateErr
}

func (s *stubUserAdminStore) UpsertLineUser(context.Context, auth.UpsertLineUserInput) (auth.User, error) {
	return auth.User{}, nil
}

type stubPlayerStore struct {
	getByIDErr    error
	getByIDResult domain.Player
}

func (s stubPlayerStore) GetByID(context.Context, string) (domain.Player, error) {
	return s.getByIDResult, s.getByIDErr
}
