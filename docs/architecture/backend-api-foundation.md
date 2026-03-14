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
- `POST /api/players`

### Sessions

- `GET /api/sessions`
- `POST /api/sessions`

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

- Create/list flows exist for foundation purposes, but feature-complete CRUD is still deferred
- auth / role enforcement is intentionally not implemented yet
- richer report generation remains a later phase
