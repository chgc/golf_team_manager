package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/config"
	appdb "github.com/chgc/golf_team_manager/backend/internal/db"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
	"github.com/chgc/golf_team_manager/backend/internal/service"
)

type userAdminService interface {
	GetByID(ctx context.Context, userID string) (auth.User, error)
	GetByProviderSubject(ctx context.Context, provider auth.Provider, subject string) (auth.User, error)
	Update(ctx context.Context, userID string, input service.UserAdminUpdateInput) (auth.User, error)
}

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	database, err := appdb.Open(cfg.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := appdb.RunMigrations(context.Background(), database); err != nil {
		fmt.Fprintf(os.Stderr, "run migrations: %v\n", err)
		os.Exit(1)
	}

	userRepository := repository.NewSQLiteUserRepository(database)
	playerRepository := repository.NewSQLitePlayerRepository(database)
	adminService := service.NewUserAdminService(userRepository, playerRepository)

	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr, adminService); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer, adminService userAdminService) error {
	if len(args) == 0 {
		return fmt.Errorf("expected subcommand: promote-user")
	}

	switch args[0] {
	case "promote-user":
		return runPromoteUser(ctx, args[1:], stdout, stderr, adminService)
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func runPromoteUser(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	adminService userAdminService,
) error {
	flagSet := flag.NewFlagSet("promote-user", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	var (
		playerID string
		provider string
		subject  string
		userID   string
	)
	flagSet.StringVar(&playerID, "player-id", "", "optional player id to link while promoting")
	flagSet.StringVar(&provider, "provider", "", "auth provider lookup value")
	flagSet.StringVar(&subject, "subject", "", "provider subject lookup value")
	flagSet.StringVar(&userID, "user-id", "", "internal user id lookup value")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(flagSet.Args(), " "))
	}

	trimmedUserID := strings.TrimSpace(userID)
	trimmedProvider := strings.TrimSpace(provider)
	trimmedSubject := strings.TrimSpace(subject)
	trimmedPlayerID := strings.TrimSpace(playerID)

	lookupByUserID := trimmedUserID != ""
	lookupByProviderSubject := trimmedProvider != "" || trimmedSubject != ""
	switch {
	case lookupByUserID && lookupByProviderSubject:
		return fmt.Errorf("use either --user-id or --provider with --subject")
	case !lookupByUserID && !lookupByProviderSubject:
		return fmt.Errorf("either --user-id or --provider with --subject is required")
	case trimmedProvider == "" && trimmedSubject != "":
		return fmt.Errorf("--provider is required when --subject is provided")
	case trimmedProvider != "" && trimmedSubject == "":
		return fmt.Errorf("--subject is required when --provider is provided")
	}

	user, err := lookupUser(ctx, adminService, trimmedUserID, trimmedProvider, trimmedSubject)
	if err != nil {
		return err
	}

	if user.Role == auth.RoleManager && (trimmedPlayerID == "" || user.PlayerID == trimmedPlayerID) {
		fmt.Fprintf(stdout, "user %s (%s) is already a manager\n", user.ID, user.DisplayName)
		return nil
	}

	managerRole := auth.RoleManager
	updateInput := service.UserAdminUpdateInput{Role: &managerRole}
	if trimmedPlayerID != "" {
		updateInput.PlayerID = &trimmedPlayerID
	}

	updatedUser, err := adminService.Update(ctx, user.ID, updateInput)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPlayerAlreadyLinked):
			return fmt.Errorf("player %q is already linked to another user", trimmedPlayerID)
		case errors.Is(err, service.ErrPlayerNotFound):
			return fmt.Errorf("player not found: %s", trimmedPlayerID)
		case errors.Is(err, service.ErrUserNotFound):
			return fmt.Errorf("user not found")
		default:
			return fmt.Errorf("promote user: %w", err)
		}
	}

	if trimmedPlayerID != "" {
		fmt.Fprintf(
			stdout,
			"promoted user %s (%s) to manager and linked player %s\n",
			updatedUser.ID,
			updatedUser.DisplayName,
			trimmedPlayerID,
		)
		return nil
	}

	fmt.Fprintf(stdout, "promoted user %s (%s) to manager\n", updatedUser.ID, updatedUser.DisplayName)
	return nil
}

func lookupUser(
	ctx context.Context,
	adminService userAdminService,
	userID string,
	provider string,
	subject string,
) (auth.User, error) {
	if userID != "" {
		user, err := adminService.GetByID(ctx, userID)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				return auth.User{}, fmt.Errorf("user not found: %s", userID)
			}

			return auth.User{}, fmt.Errorf("lookup user by id: %w", err)
		}

		return user, nil
	}

	normalizedProvider, err := parseProvider(provider)
	if err != nil {
		return auth.User{}, err
	}

	user, err := adminService.GetByProviderSubject(ctx, normalizedProvider, subject)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return auth.User{}, fmt.Errorf("user not found for provider=%s subject=%s", provider, subject)
		}

		return auth.User{}, fmt.Errorf("lookup user by provider subject: %w", err)
	}

	return user, nil
}

func parseProvider(value string) (auth.Provider, error) {
	switch auth.Provider(strings.TrimSpace(value)) {
	case auth.ProviderLINEOAuth:
		return auth.ProviderLINEOAuth, nil
	default:
		return "", fmt.Errorf("unsupported provider %q", value)
	}
}
