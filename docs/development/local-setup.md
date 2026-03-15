# Local Setup

## Scope

This document covers the current local-development workflow for the v1 MVP baseline, including auth-mode setup, seed data, startup flow, and the canonical links to demo, release-readiness, and handoff guidance.

## Required Tools

- Git
- `pnpm`
- Node.js
- Go
- `just`

## Current Workflow

### 1. Sync planning state

- Review the current plan under `docs\plans\v1-mvp\`
- Confirm the relevant work item is in `approved\`
- Confirm that the approved state has already been committed

### 2. Use quick commands

The repository root includes a `justfile` for common commands.

Examples:

```powershell
just status
just worktrees
just plans
just approved
just pending
just completed
just frontend-install
just frontend-build
just frontend-test
just backend-test
just backend-migrate
just backend-seed
just backend-start
```

Validation results for the completed Phase 1 baseline are tracked in:

- `docs\development\phase-1-validation.md`

Release/demo/handoff references:

- `docs\development\auth-setup.md`
- `docs\development\demo-smoke-check.md`
- `docs\development\release-readiness-checklist.md`
- `docs\development\v1-handoff-summary.md`

### 3. Choose the auth mode

The backend reads environment variables from the shell or process environment. `.env.example` is a reference file only; it is not auto-loaded by the repo commands.

#### Default local mode: `dev_stub`

`dev_stub` is the default backend mode and the only mode that supports `just backend-seed`.

PowerShell example:

```powershell
$env:AUTH_MODE = 'dev_stub'
$env:DB_PATH = 'data\golf_team_manager.sqlite'
just backend-seed
just backend-start
just frontend-start
```

Optional `dev_stub` overrides:

- `AUTH_DEV_DEFAULT_ROLE` (`manager` or `player`)
- `AUTH_DEV_DEFAULT_NAME`
- `AUTH_DEV_DEFAULT_SUBJECT`
- `AUTH_DEV_DEFAULT_USER_ID`
- `AUTH_DEV_DEFAULT_PLAYER_ID`

In `dev_stub`, the frontend bootstraps through `GET /api/auth/me` without a login redirect.

#### LINE mode: `line`

For local LINE SSO, set the backend env vars first:

```powershell
$env:AUTH_MODE = 'line'
$env:LINE_CLIENT_ID = '<line-channel-id>'
$env:LINE_CLIENT_SECRET = '<line-channel-secret>'
$env:LINE_REDIRECT_URI = 'http://127.0.0.1:8080/api/auth/line/callback'
$env:FRONTEND_URL = 'http://localhost:4200'
$env:JWT_SECRET = '<local-dev-jwt-secret>'
$env:JWT_TTL = '1h'
just backend-start
```

Then switch the frontend runtime config in `frontend\public\app-config.js` to match LINE mode before starting the Angular dev server:

```javascript
window.__GTM_AUTH_CONFIG = {
  authMode: 'line',
  backendOrigin: 'http://127.0.0.1:8080',
};
```

Start the frontend after the runtime config is updated:

```powershell
just frontend-start
```

Local LINE assumptions:

- frontend origin: `http://localhost:4200`
- backend origin: `http://127.0.0.1:8080`
- callback URI: `http://127.0.0.1:8080/api/auth/line/callback`
- post-login landing route: `http://localhost:4200/auth/done`

Use the backend origin for LINE login initiation. The Angular `/api/**` proxy is still valid for authenticated API traffic after the JWT is stored, but it is not the default entrypoint for local OAuth login.

### 4. Start from a deterministic local dataset

If you need the seeded players/sessions dataset, seed it while `AUTH_MODE=dev_stub`, then restart the backend in `line` mode against the same SQLite file:

```powershell
$env:AUTH_MODE = 'dev_stub'
just backend-seed

$env:AUTH_MODE = 'line'
just backend-start
```

This is useful for local auth smoke checks because the seed command intentionally fails outside `dev_stub`.

### 5. Subagent implementation flow

1. Create a proposal in `pending\`
2. Wait for user review
3. Move to `approved\` only after explicit user instruction
4. Commit the approved state
5. Implement the task
6. Move the task file to `completed\YYYY-MM-DD\`

## Current Notes

- Frontend implementation uses Angular CLI and `pnpm`
- Frontend styling uses plain CSS
- Backend implementation uses Go + Gin and no ORM
- Approved task files are implementation-ready only after the approval commit exists
- Frontend quick-start commands are exposed through the root `justfile`
- Backend smoke-test and startup commands are exposed through the root `justfile`
- Backend currently defaults to `backend\data\golf_team_manager.sqlite` and can run migrations with `just backend-migrate`
- Demo seed data can be rebuilt with `just backend-seed`
- `backend-seed` is local/dev only and currently requires `AUTH_MODE=dev_stub`
- `just frontend-start` now serves the Angular app with `/api/**` proxied to the local backend for smoke/demo use
- frontend auth runtime mode is controlled by `frontend\public\app-config.js`
- local LINE mode requires backend env vars plus a frontend runtime switch from `dev_stub` to `line`
- new LINE users are created as authenticated but unlinked players and land on `/auth/pending-link` until a manager links them to a player record
- logout is frontend-only token removal; there is no backend logout endpoint in the current stateless JWT flow
- The current manager/player demo path is documented in `docs\development\demo-smoke-check.md`
- auth-specific operations, smoke expectations, and fallback notes live in `docs\development\auth-setup.md`
- Pre-demo / pre-release gate checks live in `docs\development\release-readiness-checklist.md`
- Current scope / constraints / follow-up handoff notes live in `docs\development\v1-handoff-summary.md`

## Document Boundaries

- Use this file for local prerequisites, common commands, and startup/seed entry points.
- Use `auth-setup.md` for auth-mode details, login/logout flow, pending-link behavior, and rollback/fallback guidance.
- Use `demo-smoke-check.md` for the deterministic dataset and step-by-step smoke paths.
- Use `release-readiness-checklist.md` for the final validation checklist before demo or handoff.
- Use `v1-handoff-summary.md` for scope summary, dev-only constraints, and post-MVP follow-up notes.
