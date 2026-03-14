# Players Feature

## Scope

This document records the first feature-complete player management slice delivered on top of the existing shared domain, backend API foundation, frontend shell, and auth foundation.

## Backend Delivery

- player list now supports `query` and `status` filters
- `GET /api/players/{playerId}` returns a single player record
- `PATCH /api/players/{playerId}` supports the edit and activate/deactivate flow
- validation continues to reuse the shared player write DTO rules, including handicap range, status validity, and email format checks

## Frontend Delivery

- the players page now provides:
  - search by player name
  - active / inactive status filtering
  - create player form
  - edit player form
  - active / inactive toggle actions
- duplicate player names surface as a warning only and do not block saving
- the page keeps manager-focused scope without adding new route-level authorization

## Validation Coverage

- `go test ./...`
- `just frontend-build`
- `just frontend-test`
- runtime smoke checks for create, detail, update, and filtered list flows through `/api/players`

## Follow-up Notes

- inactive players must remain in historical records and will later be excluded from new registration choices
- richer player history and administrative fields remain later-phase work
