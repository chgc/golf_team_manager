# Shared Domain Baseline

## Scope

This document captures the Phase 2 shared-domain baseline for `Player`, `Session`, and `Registration`.

## Domain Entities

### Player

- `id`
- `name`
- `handicap`
- `phone`
- `email`
- `status`
- `notes`
- `createdAt`
- `updatedAt`

Rules:

- `handicap` must stay within `0` to `54`
- UI and API validation should enforce `0.5` increments
- duplicate names are allowed, but IDs remain unique
- inactive players stay in history and should not appear in new-registration selection flows

### Session

- `id`
- `date`
- `courseName`
- `courseAddress`
- `maxPlayers`
- `registrationDeadline`
- `status`
- `notes`
- `createdAt`
- `updatedAt`

Rules:

- `status` lifecycle: `open -> closed -> confirmed -> completed`
- `cancelled` remains an allowed terminal branch
- `registrationDeadline` must be on or before `date`
- `maxPlayers` must be positive; multiples of 4 are recommended but not required in schema

### Registration

- `id`
- `playerId`
- `sessionId`
- `status`
- `registeredAt`
- `updatedAt`

Rules:

- `status` is `confirmed` or `cancelled`
- the baseline model uses a single row per `playerId + sessionId`
- duplicate registration for the same player/session pair is blocked by a unique constraint

## SQLite Tables

The baseline migration creates:

- `players`
- `sessions`
- `registrations`
- `schema_migrations`
- `app_metadata`

Foreign keys:

- `registrations.player_id -> players.id`
- `registrations.session_id -> sessions.id`

## Derived Values

These values should be computed at query / handler / frontend level instead of stored columns:

- confirmed registration count
- remaining seats
- required group count

## API Boundary Baseline

### Players

- `GET /api/players`
- `POST /api/players`
- `GET /api/players/{playerId}`
- `PATCH /api/players/{playerId}`

### Sessions

- `GET /api/sessions`
- `POST /api/sessions`
- `GET /api/sessions/{sessionId}`
- `PATCH /api/sessions/{sessionId}`

### Registrations

- `GET /api/sessions/{sessionId}/registrations`
- `POST /api/sessions/{sessionId}/registrations`
- `PATCH /api/registrations/{registrationId}`

### Reports

- `GET /api/reports/sessions/{sessionId}/reservation-summary`

## Go Package Baseline

The shared-domain code now lives under:

- `backend\internal\domain\`

The package defines:

- core domain structs
- read/write DTO baselines
- validation helpers for create/update payloads

## Future Extension Points

- `Group` stays out of the schema until v2
- auth / role ownership rules stay out of this phase and will be handled by `auth-foundation`
- feature handlers and repositories should consume these contracts instead of redefining entity shapes
