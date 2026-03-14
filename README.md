# Golf Team Manager

Golf Team Manager is a greenfield repository for a golf-team management system built with:

- frontend: Angular + Angular Material + plain CSS + pnpm
- backend: Go + Gin + SQLite

This repository currently focuses on planning and Phase 1 bootstrap work.

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
- `just backend-start`

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

After the workspace scaffold, the next ready tasks are expected to bootstrap:

- frontend Angular shell
- backend Gin server
- SQLite migration baseline
