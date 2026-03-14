# Release Readiness Checklist

## Scope

This document is the pre-demo / pre-release gate for the current v1 MVP baseline.

Use it to confirm that the repository can still be bootstrapped locally, that the validated command set still works, and that the manager/player smoke paths remain aligned with the deterministic seed dataset.

## Canonical References

- `README.md` for repository overview and top-level quick links
- `docs\development\local-setup.md` for local prerequisites and startup/seed commands
- `docs\development\demo-smoke-check.md` for the full manager/player smoke steps
- `docs\development\v1-handoff-summary.md` for current scope, constraints, and follow-up context

## Pre-Flight Checks

- [ ] Confirm the active work item is already in `approved\` and the approval commit exists before implementation work starts.
- [ ] Confirm required tools are installed: Git, `pnpm`, Node.js, Go, and `just`.
- [ ] Confirm local development is using the expected v1 baseline:
  - frontend: Angular + Angular Material + plain CSS + pnpm
  - backend: Go + Gin + SQLite
- [ ] Confirm `AUTH_MODE=dev_stub` is used for local/demo smoke flows.

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

## Dataset and Environment Checks

- [ ] Seed dataset still contains 6 players, 4 sessions, and 7 registrations.
- [ ] Seed flow remains local/dev only.
- [ ] No release/handoff doc rewrites the smoke steps inconsistently with `demo-smoke-check.md`.
- [ ] No release/handoff doc presents `dev_stub` behavior as production-ready behavior.

## Known Constraints to Reconfirm

- `AUTH_MODE=dev_stub` is still the local baseline; production auth hardening is not part of v1.
- Player smoke remains API/debug-header based rather than an in-app role switcher.
- Local validation assumes SQLite and the repo's `justfile` command set.
- CI/CD and production deployment automation are intentionally out of scope for this release-readiness pass.

## Sign-Off

- [ ] Documentation entry points are consistent across `README.md`, `local-setup.md`, `demo-smoke-check.md`, and `v1-handoff-summary.md`.
- [ ] Validation command list matches the actual repo commands.
- [ ] Known limitations and dev-only constraints are explicit.
- [ ] Team handoff can proceed using `v1-handoff-summary.md`.
