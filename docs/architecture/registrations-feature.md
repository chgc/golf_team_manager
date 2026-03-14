# Registrations Feature

## Scope

This document records the first feature-complete registration workflow slice delivered on top of the shared domain, backend API foundation, sessions feature, and auth baseline.

## Backend Delivery

- `PATCH /api/registrations/{registrationId}` now supports registration status updates
- registration status remains a two-state v1 contract:
  - `confirmed`
  - `cancelled`
- leave / absence is intentionally modeled as `cancelled` in v1 rather than introducing a third state
- registration business validation remains consistent for:
  - duplicate registration
  - inactive player
  - session not open
  - session capacity full
- cancelled registrations can be restored to confirmed when the player is active, the session is open, and capacity is available

## Frontend Delivery

- registration flow is integrated into the existing `SessionListPage` selected-session detail area
- the standalone `/registrations` route and `registration-list-page` shell are retired
- manager flow now supports:
  - add player to session
  - cancel registration
  - restore cancelled registration
- player flow now supports:
  - register for session
  - cancel / take leave
  - register again after cancellation

## Validation Coverage

- `go test ./...`
- `just frontend-build`
- `just frontend-test`
- runtime smoke checks for register, cancel, and restore flows

## Follow-up Notes

- route-level authorization is still intentionally deferred; role boundaries remain expressed through the current auth identity context and frontend UI scope
- reservation summary reporting is the next phase that will consume the session + registration foundations
