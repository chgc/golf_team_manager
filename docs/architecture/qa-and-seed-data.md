# QA and Seed Data

## Scope

This document records the Phase 8 quality-assurance and seed-data baseline added on top of the completed v1 MVP features.

## Backend Delivery

- `backend\cmd\seed\main.go` now provides a local/dev-only seed command
- root `justfile` now exposes `just backend-seed`
- the seed command:
  - requires `AUTH_MODE=dev_stub`
  - clears and rebuilds the current SQLite database
  - produces a deterministic demo dataset
- the seeded dataset currently includes:
  - 6 players
  - 4 sessions
  - 7 registrations

## Frontend Delivery

- `HomePage` now has explicit empty-safe rendering
- `PlayerListPage` now surfaces a loading state in addition to existing empty/error behavior
- `SessionListPage` now surfaces a session-list loading state in addition to the existing summary/roster state coverage
- unit tests now cover the key loading / empty / error paths called out in the approved proposal

## Documentation Delivery

- `docs\development\local-setup.md` now includes seed-data setup guidance
- `docs\development\demo-smoke-check.md` now records the deterministic demo dataset and manager/player smoke paths

## Validation Coverage

- `go test ./...`
- `just frontend-build`
- `just frontend-test`
- `just backend-seed`
- runtime smoke checks for:
  - manager reservation summary path
  - player debug-header API path

## Follow-up Notes

- player smoke remains intentionally defined as API/debug-header smoke rather than an in-app role switcher
- release-readiness and final handoff consolidation remain the next planning step
