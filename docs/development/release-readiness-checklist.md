# Release Readiness Checklist

## Scope

This document is the pre-demo / pre-release gate for the current v1 MVP baseline.

Use it to confirm that the repository can still be bootstrapped locally, that the validated command set still works, and that the manager/player smoke paths remain aligned with the deterministic seed dataset.

## Canonical References

- `README.md` for repository overview and top-level quick links
- `docs\development\local-setup.md` for local prerequisites and startup/seed commands
- `docs\development\auth-setup.md` for auth-mode setup, LINE env requirements, and fallback guidance
- `docs\development\demo-smoke-check.md` for the full manager/player smoke steps
- `docs\development\v1-handoff-summary.md` for current scope, constraints, and follow-up context

## Pre-Flight Checks

- [ ] Confirm the active work item is already in `approved\` and the approval commit exists before implementation work starts.
- [ ] Confirm required tools are installed: Git, `pnpm`, Node.js, Go, and `just`.
- [ ] Confirm local development is using the expected v1 baseline:
  - frontend: Angular + Angular Material + plain CSS + pnpm
  - backend: Go + Gin + SQLite
- [ ] Confirm the intended auth mode is explicit before smoke validation:
  - `dev_stub` for seed-driven manager/player demo smoke
  - `line` for LINE OAuth + JWT smoke
- [ ] Confirm `just backend-seed` is only executed in `AUTH_MODE=dev_stub`.
- [ ] Confirm no auth secrets are committed; LINE values stay in local env or deployment-secret storage only.

## Command Validation

Run these commands from the repository root unless stated otherwise:

```powershell
just backend-test
just frontend-build
just frontend-test
$env:AUTH_MODE = 'dev_stub'
just backend-seed
```

Expected outcome:

- `just backend-test` passes
- `just frontend-build` passes
- `just frontend-test` passes
- `just backend-seed` succeeds and rebuilds the deterministic local demo dataset
- `just frontend-start` can drive the manager smoke path with `/api/**` requests proxied to the local backend

For LINE-mode validation, also confirm the local runtime contract before starting the apps:

- [ ] backend env vars are present: `LINE_CLIENT_ID`, `LINE_CLIENT_SECRET`, `LINE_REDIRECT_URI`, `FRONTEND_URL`, `JWT_SECRET`
- [ ] `frontend\public\app-config.js` is set to `authMode: 'line'` with `backendOrigin: 'http://127.0.0.1:8080'`
- [ ] LINE callback registration matches `http://127.0.0.1:8080/api/auth/line/callback`

## Smoke Validation

### Manager Smoke

- [ ] Follow `docs\development\demo-smoke-check.md#manager-smoke-path`.
- [ ] Confirm the seeded players list renders, including inactive player `Frank Ho`.
- [ ] Confirm the session detail view for `Sunrise Valley Golf Club` renders both the registration roster and reservation summary.
- [ ] Confirm the copy-summary action works.

### Player Smoke

- [ ] Follow `docs\development\demo-smoke-check.md#player-api-smoke-path`.
- [ ] Confirm `/api/auth/me` returns the seeded player principal when debug headers are provided.
- [ ] Confirm the open session remains visible to the player principal.
- [ ] Confirm the player's registration state for `session-open` matches the seeded dataset.

### LINE Auth Smoke

- [ ] Follow `docs\development\demo-smoke-check.md#line-auth-smoke-path`.
- [ ] Confirm the login flow starts from the backend auth origin rather than the Angular `/api/**` proxy.
- [ ] Confirm callback completion stores a JWT and `/api/auth/me` returns the authenticated principal.
- [ ] Confirm a newly authenticated, unlinked LINE user lands on `/auth/pending-link`.
- [ ] Confirm logout clears the local JWT and returns the browser to `/login`.
- [ ] If a linked-user smoke is required, document the manual `users.player_id` update that was used because no manager-linking UI exists yet.

## Dataset and Environment Checks

- [ ] Seed dataset still contains 6 players, 4 sessions, and 7 registrations.
- [ ] Seed flow remains local/dev only.
- [ ] No release/handoff doc rewrites the smoke steps inconsistently with `demo-smoke-check.md`.
- [ ] No release/handoff doc presents `dev_stub` behavior as production-ready behavior.
- [ ] Local LINE guidance still matches the documented origins: frontend `http://localhost:4200`, backend `http://127.0.0.1:8080`.
- [ ] Auth callback cookies still use the documented local assumptions (`SameSite=Lax`; `Secure` only when frontend or callback uses HTTPS).

## Known Constraints to Reconfirm

- `AUTH_MODE=dev_stub` is still the seed/demo baseline; `line` mode is the auth integration path.
- Player smoke remains API/debug-header based rather than an in-app role switcher.
- Newly authenticated LINE users remain blocked on `/auth/pending-link` until an operator or future manager flow links `users.player_id`.
- Logout is frontend-only because the current app session is a stateless JWT stored locally.
- Expired or invalid LINE JWTs should fall back to the login flow after the frontend clears the token.
- Local validation assumes SQLite and the repo's `justfile` command set.
- CI/CD and production deployment automation are intentionally out of scope for this release-readiness pass.

## Sign-Off

- [ ] Documentation entry points are consistent across `README.md`, `local-setup.md`, `auth-setup.md`, `demo-smoke-check.md`, and `v1-handoff-summary.md`.
- [ ] Validation command list matches the actual repo commands.
- [ ] Known limitations and dev-only constraints are explicit.
- [ ] Team handoff can proceed using `v1-handoff-summary.md`.
