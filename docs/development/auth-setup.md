# Auth Setup and Operations

## Scope

This document captures the current auth operating model for local development, smoke checks, and handoff. It reflects the implemented `dev_stub` and LINE SSO flows described in `docs\architecture\auth-line-sso-implementation-detail.md`.

## Related Documents

- `README.md` for top-level repo entry points
- `.env.example` for a non-secret backend env reference
- `docs\development\local-setup.md` for quick-start commands
- `docs\development\demo-smoke-check.md` for step-by-step smoke flows
- `docs\development\release-readiness-checklist.md` for pre-demo / pre-release sign-off

## Auth Modes

### `AUTH_MODE=dev_stub`

Current default local mode.

- backend injects a development principal through middleware
- frontend bootstraps immediately from `GET /api/auth/me`
- `just backend-seed` only works in this mode
- debug-header overrides remain available for API smoke checks

Supported local overrides:

- `AUTH_DEV_DEFAULT_ROLE`
- `AUTH_DEV_DEFAULT_NAME`
- `AUTH_DEV_DEFAULT_SUBJECT`
- `AUTH_DEV_DEFAULT_USER_ID`
- `AUTH_DEV_DEFAULT_PLAYER_ID`

### `AUTH_MODE=line`

LINE mode enables the implemented OAuth + JWT flow.

Required backend env vars:

- `LINE_CLIENT_ID`
- `LINE_CLIENT_SECRET`
- `LINE_REDIRECT_URI`
- `FRONTEND_URL`
- `JWT_SECRET`

Optional backend env vars:

- `JWT_TTL` (defaults to `1h`)

Required local frontend runtime config:

```javascript
window.__GTM_AUTH_CONFIG = {
  authMode: 'line',
  backendOrigin: 'http://127.0.0.1:8080',
};
```

This runtime config lives in `frontend\public\app-config.js`.

## Local LINE Assumptions

Current local contract:

- frontend origin: `http://localhost:4200`
- backend origin: `http://127.0.0.1:8080`
- callback URI: `http://127.0.0.1:8080/api/auth/line/callback`
- post-login landing page: `http://localhost:4200/auth/done`

Operational notes:

- start LINE login from `http://127.0.0.1:8080/api/auth/line/login`
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

Current limitation:

- there is no manager-linking UI or API in this repo yet
- a linked-user LINE smoke requires a manual operator update to `users.player_id`

## Validation Notes

### Dev-stub validation

- use `just backend-seed`
- use `just backend-start`
- use `just frontend-start`
- validate manager and player smoke through `docs\development\demo-smoke-check.md`

### LINE validation

- seed the deterministic dataset first in `dev_stub` if you need it
- restart the backend in `line` mode against the same SQLite file
- confirm login redirects through LINE and returns to `/auth/done#token=...`
- confirm `/api/auth/me` succeeds after callback
- confirm a new LINE user lands on `/auth/pending-link`
- confirm logout returns the browser to `/login`

## Rollback / Fallback Guidance

### Local fallback to `dev_stub`

If LINE credentials, callback registration, or the LINE provider are not available during local work:

1. set `AUTH_MODE=dev_stub`
2. restore `frontend\public\app-config.js` to:

   ```javascript
   window.__GTM_AUTH_CONFIG = {
     authMode: 'dev_stub',
     backendOrigin: 'http://127.0.0.1:8080',
   };
   ```

3. restart backend and frontend

This restores the local demo/bootstrap path without needing LINE credentials.

### JWT / callback caveats

- changing `JWT_SECRET` invalidates any stored browser token and requires users to log in again
- expired or invalid JWTs are cleared by the frontend, which falls back to the login flow on the next protected-route visit
- if callback or frontend origins move to HTTPS, re-check the cookie behavior because `Secure=true` will become active

## Non-Goals / Current Caveats

- no secrets belong in the repo; `.env.example` uses placeholders only
- refresh tokens are not implemented
- manager-driven account linking is not implemented
- production secret-management and deployment automation are outside this document's scope
