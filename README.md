# Golf Team Manager

Golf Team Manager is a golf-team management system repository built with:

- frontend: Angular + Angular Material + plain CSS + pnpm
- backend: Go + Gin + SQLite

This repository currently includes the completed v1 foundation baseline plus player, session, registration, reservation summary reporting, and LINE SSO auth.

## Repository Layout

```text
.
├── backend\             # Go service scaffold and future API implementation
├── docs\                # Planning, architecture, and development documentation
│   ├── architecture\    # Architecture decisions and diagrams
│   ├── development\     # Local workflow and developer guidance
│   └── plans\           # Phase plans, conventions, and subagent work items
└── frontend\            # Angular workspace scaffold and future UI implementation
```

## Development Workflow

### Before implementation

1. Planning and governance docs must be committed and pushed first.
2. Each subagent must prepare a task proposal under:
   - `docs\plans\v1-mvp\subagent-work-items\pending\`
3. The proposal must be reviewed and approved before implementation starts.
4. Approved proposals live under:
   - `docs\plans\v1-mvp\subagent-work-items\approved\`
5. After a proposal is moved to `approved`, commit that approval state before starting implementation.
6. After implementation is completed, move the task document into:
   - `docs\plans\v1-mvp\subagent-work-items\completed\<date>\`
7. An implementation task is only fully closed after the code/docs changes are committed **and** the task document has been moved to `completed` in a committed follow-up state.

## Development Documentation

Use these docs as the canonical entry points for day-to-day development, demo prep, and handoff:

- `docs\development\local-setup.md` for local bootstrap, common commands, and seed/startup entry points
- `docs\development\auth-setup.md` for auth-mode switching, LINE local setup, login/logout flow, pending-link behavior, and rollback guidance
- `docs\development\demo-smoke-check.md` for the deterministic dataset and manager/player smoke paths
- `docs\development\release-readiness-checklist.md` for the pre-demo / pre-release validation gate
- `docs\development\v1-handoff-summary.md` for the current v1 scope, constraints, and follow-up handoff notes
- `WORKFLOW.md` for the proposal / review / approval / implementation process

### Subagent collaboration

- Subagents should work in `git worktree` mode by default.
- Work should stay isolated per task to keep planning, review, and implementation clean.
- If a task scope changes, update the proposal doc and re-run review before continuing.
- Implementation in a worktree starts only after the approved task document has been committed.

### Frontend workflow

- Use `pnpm` as the frontend package manager.
- Use Angular CLI for workspace and code generation tasks.
- Follow Angular CLI MCP best practices for Angular implementation work.
- Keep frontend styling in plain CSS.
- If a grid table UI is needed, `ag-grid community` is an allowed option.
- Frontend worktrees are expected to share dependencies through the pnpm workflow.

### Quick commands

Use the root `justfile` as the primary quick-start entry point for common repository commands.

Examples:

- `just status`
- `just worktrees`
- `just plans`
- `just approved`
- `just pending`
- `just frontend-install`
- `just frontend-build`
- `just frontend-test`
- `just backend-test`
- `just backend-migrate`
- `just backend-seed`
- `just backend-start`

### Auth at a glance

- Local/dev and production now use LINE OAuth + JWT only.
   - backend requires `LINE_CLIENT_ID`, `LINE_CLIENT_SECRET`, `LINE_REDIRECT_URI`, `FRONTEND_URL`, and `JWT_SECRET`
   - frontend runtime auth mode is configured through `frontend\public\app-config.js` (set to `line`)
   - local LINE login starts from the backend origin (`http://localhost:8080/api/auth/line/login`), not the Angular dev-server proxy

Copy `.env.example` to a repository-root `.env` for local backend settings, and use `docs\development\auth-setup.md` for the full local auth workflow.

### Backend workflow

- Use Gin as the backend framework.
- Do not use an ORM library.
- Follow Google Go style guidelines.
- Run `gofmt` after every Go edit.
- Include tests for backend changes and keep affected functions testable.

## Current Phase

The current approved plan is tracked under:

- `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
- `docs\plans\v1-mvp\phase-0-conventions\`
- `docs\plans\v1-mvp\phase-1-foundation\`

## Next Steps

The current v1 baseline has been validated across:

- Angular frontend build and unit tests
- Go backend unit tests
- SQLite migration bootstrap
- backend startup and `/health` smoke check

The current release-readiness and handoff references live under:

- `docs\development\auth-setup.md`
- `docs\development\release-readiness-checklist.md`
- `docs\development\v1-handoff-summary.md`

The backend currently defaults to a local SQLite database at `backend\data\golf_team_manager.sqlite` and runs baseline migrations on startup.
