# Demo Smoke Check

## Scope

This document records the current v1 MVP demo flow after running local seed data.

## Related Documents

- `docs\development\local-setup.md` for prerequisites and command entry points
- `docs\development\auth-setup.md` for auth-mode setup, local LINE assumptions, and auth operations notes
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

Invoke-RestMethod -Uri 'http://localhost:8080/api/auth/me' -Headers $headers
Invoke-RestMethod -Uri 'http://localhost:8080/api/sessions' -Headers $headers
Invoke-RestMethod -Uri 'http://localhost:8080/api/sessions/session-open/registrations' -Headers $headers
```

Expected outcome:

- `/api/auth/me` returns `role=player` and `playerId=player-ben`
- the open session remains visible
- the player registration state for `session-open` is present and aligned with the seed data

## Notes

- The seed command is local/dev only and must not be used outside `AUTH_MODE=dev_stub`
- Player smoke is intentionally defined as API/debug-header smoke, not an in-app role switcher
- This document is the source of truth for demo smoke details; release/handoff docs should link here rather than duplicate the full walkthrough

## LINE Auth Smoke Path

Use this path after the deterministic dataset is already available locally. Seed first in `dev_stub`, then restart the backend in `line` mode against the same SQLite database.

### LINE-mode prerequisites

- `frontend\public\app-config.js` is set to:

  ```javascript
  window.__GTM_AUTH_CONFIG = {
    authMode: 'line',
    backendOrigin: 'http://localhost:8080',
  };
  ```

- backend env vars are set:

  ```powershell
  $env:AUTH_MODE = 'line'
  $env:LINE_CLIENT_ID = '<line-channel-id>'
  $env:LINE_CLIENT_SECRET = '<line-channel-secret>'
  $env:LINE_REDIRECT_URI = 'http://localhost:8080/api/auth/line/callback'
  $env:FRONTEND_URL = 'http://localhost:4200'
  $env:JWT_SECRET = '<local-dev-jwt-secret>'
  $env:JWT_TTL = '1h'
  ```

- backend API is running through `just backend-start`
- frontend app is running through `just frontend-start`
- the LINE developer app callback is registered as `http://localhost:8080/api/auth/line/callback`

### New / unlinked LINE user smoke

1. Open `http://localhost:4200/login`
2. Click **Continue with LINE**
3. Confirm the browser navigates to `http://localhost:8080/api/auth/line/login` before redirecting to LINE
4. Complete LINE sign-in
5. Confirm the backend redirects to `http://localhost:4200/auth/done#token=...`
6. Confirm the app lands on `/auth/pending-link`
7. Confirm the page explains that the account is authenticated but not yet linked to a player record
8. Click **Logout**
9. Confirm the app returns to `/login`

Expected outcome:

- backend creates or updates a `users` row with `auth_provider=line`
- `/api/auth/me` succeeds after callback and returns `role=player` with `playerId` omitted for an unlinked user
- protected app routes remain blocked by the pending-link page until the user is linked
- logout removes the local JWT and requires a new login

### Linked-user follow-up smoke

Use a manager account plus the new admin UI instead of manual SQLite edits.

1. Complete the unlinked-user smoke above once so the target LINE account exists in `users`
2. Bootstrap or sign in as a manager account
3. Open `http://localhost:4200/admin/users`
4. Confirm the target account appears in **Unlinked Accounts**
5. Choose the intended player in **Player link**
6. Optional: change **Role** to `manager` if this user should also manage the team
7. Click **Save changes**
8. Confirm the success message appears and the user moves into **Linked Accounts**
9. Log out the manager
10. Sign in again as the linked LINE user

Expected linked-user outcome:

- the manager can complete linkage without leaving the app
- `/api/auth/me` returns the linked `playerId`
- the linked user is redirected into the normal protected app routes instead of `/auth/pending-link`
- if the role was promoted to `manager`, the linked user can also open `/admin/users`
