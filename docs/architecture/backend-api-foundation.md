# Backend API Foundation

## Scope

This document records the current backend API foundation that sits on top of the shared domain schema baseline.

## Layering

The backend now follows a simple layered structure:

- `internal\repository\` for SQLite persistence
- `internal\service\` for validation-aware application logic
- `internal\http\` for Gin route registration and transport handlers

## Current Route Baseline

### Health

- `GET /health`

### Players

- `GET /api/players`
- `GET /api/players/{playerId}`
- `POST /api/players`
- `PATCH /api/players/{playerId}`
- `GET /api/players` supports `query` and `status` filters

### Sessions

- `GET /api/sessions`
- `GET /api/sessions/{sessionId}`
- `POST /api/sessions`
- `PATCH /api/sessions/{sessionId}`

### Registrations

- `GET /api/sessions/{sessionId}/registrations`
- `POST /api/sessions/{sessionId}/registrations`

### Reports

- `GET /api/reports/sessions/{sessionId}/reservation-summary`
  - currently reserved and returns `501 not_implemented`

## Error Response Shape

All handled API errors use:

```json
{
  "error": {
    "code": "validation_failed",
    "message": "request validation failed",
    "details": ["name is required"]
  }
}
```

Current error codes include:

- `validation_failed`
- `not_found`
- `conflict`
- `player_inactive`
- `session_not_open`
- `session_capacity_full`
- `not_implemented`
- `internal_error`

## Middleware Baseline

- Gin logger
- Gin recovery
- request ID middleware with `X-Request-ID`

## Current Boundaries

- player management now supports manager-facing list, detail, filtering, create, and update flows
- session management now supports manager-facing list, detail, create, update, status transitions, and auto-close reconciliation
- registration flows still remain at foundation scope
- auth / role enforcement is intentionally not implemented yet
- richer report generation remains a later phase
