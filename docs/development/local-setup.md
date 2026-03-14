# Local Setup

## Scope

This document covers the current Phase 1 local-development workflow for repository bootstrap work.

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
just backend-start
```

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

## Future Expansion

This document should be updated after:

- Angular workspace bootstrap
- Go backend bootstrap
- SQLite migration bootstrap

At that point, concrete frontend and backend startup commands should be added.
