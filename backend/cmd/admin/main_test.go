package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
	"github.com/chgc/golf_team_manager/backend/internal/service"
)

func TestRunPromoteUserRequiresLookupMode(t *testing.T) {
	err := run(context.Background(), []string{"promote-user"}, &bytes.Buffer{}, &bytes.Buffer{}, &stubUserAdminService{})
	if err == nil || !strings.Contains(err.Error(), "either --user-id or --provider with --subject is required") {
		t.Fatalf("run() error = %v, want lookup requirement", err)
	}
}

func TestRunListUsersPrintsUsersTable(t *testing.T) {
	stdout := &bytes.Buffer{}
	service := &stubUserAdminService{
		listResult: []auth.User{
			{
				ID:          "user-1",
				DisplayName: "Alice",
				Role:        auth.RoleManager,
				Provider:    auth.ProviderLINEOAuth,
				Subject:     "subject-1",
				PlayerID:    "player-1",
			},
			{
				ID:          "user-2",
				DisplayName: "Bob",
				Role:        auth.RolePlayer,
				Provider:    auth.ProviderLINEOAuth,
				Subject:     "subject-2",
			},
		},
	}

	err := run(context.Background(), []string{"list-users"}, stdout, &bytes.Buffer{}, service)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "USER_ID") || !strings.Contains(got, "user-1") || !strings.Contains(got, "subject-2") {
		t.Fatalf("stdout = %q, want users table", got)
	}

	if !strings.Contains(got, "\t") && !strings.Contains(got, "  ") {
		t.Fatalf("stdout = %q, want tabular output", got)
	}
}

func TestRunListUsersAppliesFilters(t *testing.T) {
	service := &stubUserAdminService{}

	err := run(
		context.Background(),
		[]string{"list-users", "--role", "manager", "--link-state", "unlinked"},
		&bytes.Buffer{},
		&bytes.Buffer{},
		service,
	)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if service.listFilter.Role != auth.RoleManager {
		t.Fatalf("listFilter.Role = %q, want %q", service.listFilter.Role, auth.RoleManager)
	}

	if service.listFilter.LinkState != "unlinked" {
		t.Fatalf("listFilter.LinkState = %q, want %q", service.listFilter.LinkState, "unlinked")
	}
}

func TestRunListUsersRejectsUnsupportedRole(t *testing.T) {
	err := run(
		context.Background(),
		[]string{"list-users", "--role", "coach"},
		&bytes.Buffer{},
		&bytes.Buffer{},
		&stubUserAdminService{},
	)
	if err == nil || !strings.Contains(err.Error(), "unsupported role") {
		t.Fatalf("run() error = %v, want unsupported role message", err)
	}
}

func TestRunPromoteUserRejectsMixedLookupModes(t *testing.T) {
	err := run(
		context.Background(),
		[]string{"promote-user", "--user-id", "user-1", "--provider", "line", "--subject", "subject-1"},
		&bytes.Buffer{},
		&bytes.Buffer{},
		&stubUserAdminService{},
	)
	if err == nil || !strings.Contains(err.Error(), "use either --user-id or --provider with --subject") {
		t.Fatalf("run() error = %v, want mixed lookup validation", err)
	}
}

func TestRunPromoteUserPromotesByUserID(t *testing.T) {
	stdout := &bytes.Buffer{}
	service := &stubUserAdminService{
		getByIDResult: auth.User{
			ID:          "user-1",
			DisplayName: "Alice",
			Role:        auth.RolePlayer,
		},
		updateResult: auth.User{
			ID:          "user-1",
			DisplayName: "Alice",
			Role:        auth.RoleManager,
		},
	}

	err := run(
		context.Background(),
		[]string{"promote-user", "--user-id", "user-1"},
		stdout,
		&bytes.Buffer{},
		service,
	)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if !service.updateCalled {
		t.Fatal("Update() called = false, want true")
	}

	if got := stdout.String(); !strings.Contains(got, "promoted user user-1 (Alice) to manager") {
		t.Fatalf("stdout = %q, want promote message", got)
	}
}

func TestRunPromoteUserPromotesByProviderSubjectAndLinksPlayer(t *testing.T) {
	stdout := &bytes.Buffer{}
	service := &stubUserAdminService{
		getByProviderSubjectResult: auth.User{
			ID:          "user-1",
			DisplayName: "Alice",
			Role:        auth.RolePlayer,
		},
		updateResult: auth.User{
			ID:          "user-1",
			DisplayName: "Alice",
			Role:        auth.RoleManager,
			PlayerID:    "player-1",
		},
	}

	err := run(
		context.Background(),
		[]string{"promote-user", "--provider", "line", "--subject", "subject-1", "--player-id", "player-1"},
		stdout,
		&bytes.Buffer{},
		service,
	)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if service.updatePlayerID == nil || *service.updatePlayerID != "player-1" {
		t.Fatalf("updatePlayerID = %v, want %q", service.updatePlayerID, "player-1")
	}

	if got := stdout.String(); !strings.Contains(got, "linked player player-1") {
		t.Fatalf("stdout = %q, want linked player message", got)
	}
}

func TestRunPromoteUserTreatsExistingManagerAsNoOp(t *testing.T) {
	stdout := &bytes.Buffer{}
	service := &stubUserAdminService{
		getByIDResult: auth.User{
			ID:          "user-1",
			DisplayName: "Alice",
			Role:        auth.RoleManager,
		},
	}

	err := run(
		context.Background(),
		[]string{"promote-user", "--user-id", "user-1"},
		stdout,
		&bytes.Buffer{},
		service,
	)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if service.updateCalled {
		t.Fatal("Update() called = true, want false")
	}

	if got := stdout.String(); !strings.Contains(got, "already a manager") {
		t.Fatalf("stdout = %q, want no-op message", got)
	}
}

func TestRunPromoteUserReturnsNotFoundError(t *testing.T) {
	service := &stubUserAdminService{getByIDErr: service.ErrUserNotFound}

	err := run(
		context.Background(),
		[]string{"promote-user", "--user-id", "missing-user"},
		&bytes.Buffer{},
		&bytes.Buffer{},
		service,
	)
	if err == nil || !strings.Contains(err.Error(), "user not found: missing-user") {
		t.Fatalf("run() error = %v, want user not found message", err)
	}
}

func TestRunPromoteUserReturnsPlayerConflictError(t *testing.T) {
	service := &stubUserAdminService{
		getByIDResult: auth.User{
			ID:          "user-1",
			DisplayName: "Alice",
			Role:        auth.RolePlayer,
		},
		updateErr: service.ErrPlayerAlreadyLinked,
	}

	err := run(
		context.Background(),
		[]string{"promote-user", "--user-id", "user-1", "--player-id", "player-1"},
		&bytes.Buffer{},
		&bytes.Buffer{},
		service,
	)
	if err == nil || !strings.Contains(err.Error(), "player \"player-1\" is already linked") {
		t.Fatalf("run() error = %v, want player conflict message", err)
	}
}

type stubUserAdminService struct {
	getByIDErr                 error
	getByIDResult              auth.User
	getByProviderSubjectErr    error
	getByProviderSubjectResult auth.User
	listErr                    error
	listFilter                 repository.UserListFilter
	listResult                 []auth.User
	updateCalled               bool
	updateErr                  error
	updatePlayerID             *string
	updateResult               auth.User
}

func (s *stubUserAdminService) GetByID(context.Context, string) (auth.User, error) {
	return s.getByIDResult, s.getByIDErr
}

func (s *stubUserAdminService) GetByProviderSubject(context.Context, auth.Provider, string) (auth.User, error) {
	return s.getByProviderSubjectResult, s.getByProviderSubjectErr
}

func (s *stubUserAdminService) List(_ context.Context, filter repository.UserListFilter) ([]auth.User, error) {
	s.listFilter = filter
	return s.listResult, s.listErr
}

func (s *stubUserAdminService) Update(
	_ context.Context,
	_ string,
	input service.UserAdminUpdateInput,
) (auth.User, error) {
	s.updateCalled = true
	s.updatePlayerID = input.PlayerID
	return s.updateResult, s.updateErr
}

func TestParseProviderRejectsUnsupportedValues(t *testing.T) {
	_, err := parseProvider("github")
	if err == nil || !strings.Contains(err.Error(), "unsupported provider") {
		t.Fatalf("parseProvider() error = %v, want unsupported provider", err)
	}
}

func TestLookupUserByProviderSubjectWrapsUnexpectedErrors(t *testing.T) {
	service := &stubUserAdminService{getByProviderSubjectErr: errors.New("boom")}

	_, err := lookupUser(context.Background(), service, "", "line", "subject-1")
	if err == nil || !strings.Contains(err.Error(), "lookup user by provider subject") {
		t.Fatalf("lookupUser() error = %v, want wrapped provider lookup error", err)
	}
}
