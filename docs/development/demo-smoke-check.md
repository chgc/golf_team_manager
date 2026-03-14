# Demo Smoke Check

## Scope

This document records the current v1 MVP demo flow after running local seed data.

## Related Documents

- `docs\development\local-setup.md` for prerequisites and command entry points
- `docs\development\release-readiness-checklist.md` for the final validation gate before demo/handoff
- `docs\development\v1-handoff-summary.md` for current scope, limitations, and follow-up context

## Prerequisites

- `AUTH_MODE=dev_stub`
- `DB_PATH` points to a local SQLite file
- `just backend-seed` has completed successfully
- backend API is running through `just backend-start`
- frontend app is running through `just frontend-start` with local `/api/**` proxying to the backend

## Seeded Demo Dataset

- Players: `6`
  - Active players: `5`
  - Inactive players: `1`
- Sessions: `4`
  - `session-open`
  - `session-confirmed`
  - `session-completed`
  - `session-cancelled`
- Registrations: `7`
  - Confirmed: `5`
  - Cancelled: `2`

Fixed player smoke identity:

- `playerId`: `player-ben`
- `displayName`: `Ben Lin`

## Manager Smoke Path

1. Run `just backend-seed`
2. Start the backend with `just backend-start`
3. Open the frontend and browse to `/players`
4. Confirm the seeded players render, including inactive player `Frank Ho`
5. Browse to `/sessions`
6. Open `Sunrise Valley Golf Club`
7. Confirm the registration roster and reservation summary render
8. Use the copy action and confirm the summary text copies successfully

## Player API Smoke Path

Use PowerShell:

```powershell
$headers = @{
  'X-Debug-Role' = 'player'
  'X-Debug-Player-ID' = 'player-ben'
  'X-Debug-Display-Name' = 'Ben Lin'
}

Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/auth/me' -Headers $headers
Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/sessions' -Headers $headers
Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/sessions/session-open/registrations' -Headers $headers
```

Expected outcome:

- `/api/auth/me` returns `role=player` and `playerId=player-ben`
- the open session remains visible
- the player registration state for `session-open` is present and aligned with the seed data

## Notes

- The seed command is local/dev only and must not be used outside `AUTH_MODE=dev_stub`
- Player smoke is intentionally defined as API/debug-header smoke, not an in-app role switcher
- This document is the source of truth for demo smoke details; release/handoff docs should link here rather than duplicate the full walkthrough
