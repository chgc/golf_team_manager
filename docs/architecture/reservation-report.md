# Reservation Report Feature

## Scope

This document records the first feature-complete reservation summary reporting slice delivered on top of the sessions and registrations workflow.

## Backend Delivery

- `GET /api/reports/sessions/{sessionId}/reservation-summary` now returns a manager-only reservation summary
- the report response includes:
  - session date
  - course name
  - course address
  - registration deadline
  - session status
  - confirmed player count
  - estimated groups
  - copy-ready `summaryText`
  - sorted `confirmedPlayers`
- reservation summary output is currently limited to `confirmed` and `completed` sessions
- report-specific errors now distinguish:
  - `404 session_not_found`
  - `422 session_not_eligible_for_report`
  - `422 reservation_summary_empty`
  - `403 forbidden` for non-manager access

## Frontend Delivery

- the reservation summary UI is integrated into the existing `SessionListPage` selected-session detail area
- manager flow now supports:
  - loading the reservation summary for eligible sessions
  - showing inline ineligible / empty / error states without polluting the page-level error banner
  - copying the backend-provided `summaryText` through the Clipboard API
- the summary card reuses the backend response as the source of truth for:
  - confirmed player count
  - estimated groups
  - roster order
  - copy-ready summary text

## Validation Coverage

- `go test ./...`
- `just frontend-build`
- `just frontend-test`
- runtime smoke checks for player create, session create, registration create, session confirm, and reservation summary fetch

## Follow-up Notes

- the current summary text is intentionally plain-text only and does not yet include PDF / Excel export
- tee time, grouping, and richer reservation formatting remain future-phase work
