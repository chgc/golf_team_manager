# Sessions Feature

## Scope

This document records the first feature-complete session management slice delivered on top of the shared domain, backend API foundation, frontend shell, auth foundation, and players feature baseline.

## Backend Delivery

- `GET /api/sessions/{sessionId}` now returns a single session record
- `PATCH /api/sessions/{sessionId}` supports editing and manager-driven status transitions
- session status transitions are validated against the approved matrix:
  - `open -> closed`
  - `open -> cancelled`
  - `closed -> confirmed`
  - `closed -> cancelled`
  - `confirmed -> completed`
  - `confirmed -> cancelled`
- expired `open` sessions auto-close during service-level reconciliation on session reads and status mutation entrypoints

## Frontend Delivery

- the sessions page now provides:
  - upcoming / history session views
  - create / edit session form
  - session detail card
  - manager status actions for close / confirm / cancel / complete
- session detail shows:
  - date
  - course
  - deadline
  - status
  - confirmed player count
  - remaining spots
  - estimated groups
  - notes
- registration counts are derived from the existing registrations API so the page stays aligned with current backend boundaries

## Validation Coverage

- `go test ./...`
- `just frontend-build`
- `just frontend-test`
- runtime smoke checks for session create, detail, update, valid / invalid status transitions, and expired-deadline auto-close

## Follow-up Notes

- registration roster management remains a separate feature phase
- future report generation can build on the same session detail and registration summary foundations
