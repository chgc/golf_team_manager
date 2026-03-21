# Auth Setup and Operations

## Scope

This document captures the current auth operating model for local development, smoke checks, and handoff. It reflects the LINE SSO flow described in `docs\architecture\auth-line-sso-implementation-detail.md`.

## Related Documents

- `README.md` for top-level repo entry points
- `.env.example` for a non-secret backend env reference
- `docs\development\local-setup.md` for quick-start commands
- `docs\development\demo-smoke-check.md` for step-by-step smoke flows
- `docs\development\release-readiness-checklist.md` for pre-demo / pre-release sign-off

## Auth Mode

LINE mode enables the implemented OAuth + JWT flow.

Required backend env vars:

- `LINE_CLIENT_ID`
- `LINE_CLIENT_SECRET`
- `LINE_REDIRECT_URI`
- `FRONTEND_URL`
- `JWT_SECRET`

Optional backend env vars:

- `JWT_TTL` (defaults to `1h`)

For local development, the backend auto-loads a repository-root `.env` file if present. Values already set in the shell still override `.env`.

Required local frontend runtime config:

```javascript
window.__GTM_AUTH_CONFIG = {
  authMode: 'line',
  backendOrigin: 'http://localhost:8080',
};
```

This runtime config lives in `frontend\public\app-config.js`.

## Local LINE Assumptions

Current local contract:

- frontend origin: `http://localhost:4200`
- backend origin: `http://localhost:8080`
- callback URI: `http://localhost:8080/api/auth/line/callback`
- post-login landing page: `http://localhost:4200/auth/done`

Operational notes:

- start LINE login from `http://localhost:8080/api/auth/line/login`
- do not rely on the Angular `/api/**` proxy to initiate the OAuth redirect
- after login succeeds, normal API traffic can keep using the frontend `/api/**` proxy because the frontend sends `Authorization: Bearer <token>`
- OAuth state and nonce are stored in the `gtm_line_oauth` HttpOnly cookie
- the cookie uses `SameSite=Lax`
- the cookie only becomes `Secure=true` when `LINE_REDIRECT_URI` or `FRONTEND_URL` uses HTTPS

## Login and Logout Summary

### LINE login flow

1. User opens `/login`
2. Frontend button points at the backend auth origin: `/api/auth/line/login`
3. Backend creates OAuth `state` and `nonce`, writes `gtm_line_oauth`, and redirects to LINE
4. LINE redirects back to `/api/auth/line/callback`
5. Backend verifies the callback, upserts the `users` record, signs an app JWT, clears `gtm_line_oauth`, and redirects to `/auth/done#token=...`
6. Frontend stores the JWT in localStorage, calls `GET /api/auth/me`, and routes the user based on the returned principal

### Logout flow

- logout is frontend-only
- the app clears the local JWT and pending redirect state
- there is no backend logout endpoint in the current stateless JWT design

## Pending-Link / Unlinked User Behavior

The implemented LINE contract creates new authenticated users as:

- `auth_provider = line`
- `role = player`
- `player_id = NULL`

User-visible behavior:

- `GET /api/auth/me` returns `200`
- `playerId` is omitted for an unlinked player
- the frontend treats `role=player` plus missing `playerId` as authenticated-but-unlinked
- the app redirects that user to `/auth/pending-link`
- logout remains available from the pending-link page

Current manager tooling:

- manager user administration UI is now available at `/admin/users`
- backend manager user administration API is available at:
  - `GET /api/admin/users`
  - `PATCH /api/admin/users/:userId`
- pending-link remains a manager-assisted flow; unlinked end users still cannot self-link inside the app

## Bootstrap the First Manager

After Phase 9 foundation and bootstrap CLI are available, local / operator setup should stop relying on manual SQLite edits to assign the first manager.

Recommended flow:

1. Let the target user complete one LINE login so the backend creates the `users` row.
2. If you need the internal `user-id`, list current users from the repository root:

   ```powershell
   just list-users
   ```

   Optional filters:

   ```powershell
   just list-users -- --link-state unlinked
   just list-users -- --role manager
   ```

3. Run the helper from the repository root:

   ```powershell
   just promote-manager -- --user-id <user-id>
   ```

   For the common LINE-subject flow, the shortcut also supports:

   ```powershell
   just promote-manager -- --subject <line-subject>
   ```

4. Optional: set the initial player link at the same time:

   ```powershell
   just promote-manager -- --subject <line-subject> --player-id <player-id>
   ```

5. If needed, the Node helpers and backend CLI remain available:

   ```powershell
   node scripts/list-users.mjs
   ```

   ```powershell
   node scripts/promote-manager.mjs --user-id <user-id>
   ```

   ```powershell
   Set-Location backend
   ```

   ```powershell
   go run ./cmd/admin promote-user --user-id <user-id>
   ```

Command rules:

- use either `--user-id` or `--provider` with `--subject`
- use `list-users` if you need to discover the generated `user-id`
- the helper script defaults `--subject` lookups to `--provider line`
- the command only updates an existing `users` row; it never creates a user
- running it for an already-manager user is an idempotent no-op
- invalid user / player lookups return a non-zero exit

## Validation Notes

### LINE validation

- seed the deterministic dataset first if you need it
- start backend with LINE env config against the same SQLite file
- confirm login redirects through LINE and returns to `/auth/done#token=...`
- confirm `/api/auth/me` succeeds after callback
- confirm a new LINE user lands on `/auth/pending-link`
- confirm a manager can open `/admin/users`, link that account to a player, and optionally promote/demote role
- confirm the linked user reaches the normal protected app after the next login
- confirm logout returns the browser to `/login`

## Rollback / Fallback Guidance

If LINE credentials, callback registration, or the LINE provider are unavailable, local protected flows are blocked until LINE configuration is restored.

### JWT / callback caveats

- changing `JWT_SECRET` invalidates any stored browser token and requires users to log in again
- expired or invalid JWTs are cleared by the frontend, which falls back to the login flow on the next protected-route visit
- if callback or frontend origins move to HTTPS, re-check the cookie behavior because `Secure=true` will become active

## Non-Goals / Current Caveats

- no secrets belong in the repo; `.env.example` uses placeholders only
- refresh tokens are not implemented
- end-user self-service account linking is not implemented
- production secret-management and deployment automation are outside this document's scope
