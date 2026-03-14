package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/config"
	appdb "github.com/chgc/golf_team_manager/backend/internal/db"
)

type seededPlayer struct {
	id       string
	name     string
	handicap float64
	phone    string
	email    string
	status   string
	notes    string
}

type seededSession struct {
	id                   string
	date                 string
	courseName           string
	courseAddress        string
	maxPlayers           int
	registrationDeadline string
	status               string
	notes                string
}

type seededRegistration struct {
	id           string
	playerID     string
	sessionID    string
	status       string
	registeredAt string
}

const seedTimestamp = "2026-03-14T00:00:00Z"

var seedPlayers = []seededPlayer{
	{id: "player-alice", name: "Alice Chen", handicap: 10.5, phone: "0911-000-001", email: "alice@example.com", status: "active", notes: "Seed player for reservation summary"},
	{id: "player-ben", name: "Ben Lin", handicap: 14, phone: "0911-000-002", email: "ben@example.com", status: "active", notes: "Open-session player"},
	{id: "player-cathy", name: "Cathy Wang", handicap: 18.5, phone: "0911-000-003", email: "cathy@example.com", status: "active", notes: "Confirmed-session player"},
	{id: "player-daniel", name: "Daniel Wu", handicap: 9, phone: "0911-000-004", email: "daniel@example.com", status: "active", notes: "Completed-session player"},
	{id: "player-ella", name: "Ella Tsai", handicap: 22, phone: "0911-000-005", email: "ella@example.com", status: "active", notes: "Cancelled registration player"},
	{id: "player-frank", name: "Frank Ho", handicap: 30, phone: "0911-000-006", email: "frank@example.com", status: "inactive", notes: "Inactive seed player"},
}

var seedSessions = []seededSession{
	{id: "session-open", date: "2026-04-12T08:00:00Z", courseName: "Open Ridge Golf Club", courseAddress: "Taoyuan", maxPlayers: 8, registrationDeadline: "2026-04-05T23:59:00Z", status: "open", notes: "Open session for player smoke"},
	{id: "session-confirmed", date: "2026-04-19T08:00:00Z", courseName: "Sunrise Valley Golf Club", courseAddress: "Taipei", maxPlayers: 8, registrationDeadline: "2026-04-10T23:59:00Z", status: "confirmed", notes: "Confirmed session for summary smoke"},
	{id: "session-completed", date: "2026-03-08T08:00:00Z", courseName: "Harbor Links", courseAddress: "Kaohsiung", maxPlayers: 8, registrationDeadline: "2026-03-01T23:59:00Z", status: "completed", notes: "Completed session in history"},
	{id: "session-cancelled", date: "2026-04-26T08:00:00Z", courseName: "Pine Hills", courseAddress: "Hsinchu", maxPlayers: 8, registrationDeadline: "2026-04-18T23:59:00Z", status: "cancelled", notes: "Cancelled session for history view"},
}

var seedRegistrations = []seededRegistration{
	{id: "registration-1", playerID: "player-alice", sessionID: "session-confirmed", status: "confirmed", registeredAt: "2026-04-01T09:00:00Z"},
	{id: "registration-2", playerID: "player-cathy", sessionID: "session-confirmed", status: "confirmed", registeredAt: "2026-04-01T09:05:00Z"},
	{id: "registration-3", playerID: "player-ella", sessionID: "session-confirmed", status: "confirmed", registeredAt: "2026-04-01T09:10:00Z"},
	{id: "registration-4", playerID: "player-daniel", sessionID: "session-completed", status: "confirmed", registeredAt: "2026-02-25T09:00:00Z"},
	{id: "registration-5", playerID: "player-ben", sessionID: "session-open", status: "confirmed", registeredAt: "2026-04-02T09:00:00Z"},
	{id: "registration-6", playerID: "player-frank", sessionID: "session-open", status: "cancelled", registeredAt: "2026-04-02T09:15:00Z"},
	{id: "registration-7", playerID: "player-ella", sessionID: "session-cancelled", status: "cancelled", registeredAt: "2026-04-03T09:00:00Z"},
}

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if cfg.Auth.Mode != "dev_stub" {
		log.Fatalf("backend seed requires AUTH_MODE=dev_stub; current value: %s", cfg.Auth.Mode)
	}

	database, err := appdb.Open(cfg.DB)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := appdb.RunMigrations(context.Background(), database); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	if err := reseedDatabase(context.Background(), database); err != nil {
		log.Fatalf("reseed database: %v", err)
	}

	log.Printf("seeded %d players, %d sessions, and %d registrations into %s", len(seedPlayers), len(seedSessions), len(seedRegistrations), cfg.DB.Path)
}

func reseedDatabase(ctx context.Context, database *sql.DB) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}

	defer tx.Rollback()

	for _, statement := range []string{
		"DELETE FROM registrations",
		"DELETE FROM users",
		"DELETE FROM sessions",
		"DELETE FROM players",
	} {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("exec %q: %w", statement, err)
		}
	}

	if err := insertPlayers(ctx, tx); err != nil {
		return err
	}

	if err := insertSessions(ctx, tx); err != nil {
		return err
	}

	if err := insertRegistrations(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}

func insertPlayers(ctx context.Context, tx *sql.Tx) error {
	for _, player := range seedPlayers {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO players (id, name, handicap, phone, email, status, notes, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			player.ID(),
			player.name,
			player.handicap,
			player.phone,
			player.email,
			player.status,
			player.notes,
			seedTimestamp,
			seedTimestamp,
		); err != nil {
			return fmt.Errorf("insert player %s: %w", player.id, err)
		}
	}

	return nil
}

func insertSessions(ctx context.Context, tx *sql.Tx) error {
	for _, session := range seedSessions {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO sessions (id, session_date, course_name, course_address, max_players, registration_deadline, status, notes, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			session.id,
			session.date,
			session.courseName,
			session.courseAddress,
			session.maxPlayers,
			session.registrationDeadline,
			session.status,
			session.notes,
			seedTimestamp,
			seedTimestamp,
		); err != nil {
			return fmt.Errorf("insert session %s: %w", session.id, err)
		}
	}

	return nil
}

func insertRegistrations(ctx context.Context, tx *sql.Tx) error {
	for _, registration := range seedRegistrations {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO registrations (id, player_id, session_id, status, registered_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)`,
			registration.id,
			registration.playerID,
			registration.sessionID,
			registration.status,
			registration.registeredAt,
			registration.registeredAt,
		); err != nil {
			return fmt.Errorf("insert registration %s: %w", registration.id, err)
		}
	}

	return nil
}

func (p seededPlayer) ID() string {
	return p.id
}

func init() {
	time.Local = time.UTC
}
