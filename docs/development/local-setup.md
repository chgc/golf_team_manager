# Local Setup

## Scope

This document covers the current local-development workflow for the v1 MVP baseline, including seed data, startup flow, and the canonical links to demo, release-readiness, and handoff guidance.

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

- `docs\development\demo-smoke-check.md`
- `docs\development\release-readiness-checklist.md`
- `docs\development\v1-handoff-summary.md`

### 3. Subagent implementation flow

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
- The current manager/player demo path is documented in `docs\development\demo-smoke-check.md`
- Pre-demo / pre-release gate checks live in `docs\development\release-readiness-checklist.md`
- Current scope / constraints / follow-up handoff notes live in `docs\development\v1-handoff-summary.md`

## Document Boundaries

- Use this file for local prerequisites, common commands, and startup/seed entry points.
- Use `demo-smoke-check.md` for the deterministic dataset and step-by-step smoke paths.
- Use `release-readiness-checklist.md` for the final validation checklist before demo or handoff.
- Use `v1-handoff-summary.md` for scope summary, dev-only constraints, and post-MVP follow-up notes.
