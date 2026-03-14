# v1 Handoff Summary

## Scope

This document summarizes the current v1 MVP delivery, the canonical operator/developer entry points, and the main constraints that must be understood before demo, handoff, or follow-up planning.

## Delivered v1 Surface

The repository currently includes:

- player management flow
- session management flow
- registration flow
- reservation summary/report flow
- deterministic local seed data for demo/smoke usage
- development documentation for setup, smoke, release readiness, and workflow governance

## Canonical Entry Points

- `README.md` for repository overview and quick links
- `WORKFLOW.md` for proposal/review/approval/implementation rules
- `docs\development\local-setup.md` for local bootstrap and common commands
- `docs\development\demo-smoke-check.md` for detailed smoke execution
- `docs\development\release-readiness-checklist.md` for pre-demo / pre-release validation

## Operator / Demo Path Summary

1. Use `docs\development\local-setup.md` to confirm tooling and command entry points.
2. Set local auth to `AUTH_MODE=dev_stub`.
3. Run `just backend-seed` to rebuild the deterministic demo dataset.
4. Run the manager and player smoke paths from `docs\development\demo-smoke-check.md`.
5. Use `docs\development\release-readiness-checklist.md` as the final gate before demo or handoff.

## Current Constraints and Dev-Only Rules

- Local/demo auth currently depends on `AUTH_MODE=dev_stub`.
- Debug headers can override the development principal for player smoke checks.
- `backend-seed` is local/dev only and must not be treated as a production data path.
- Player smoke is intentionally API/debug-header based, not a frontend identity switcher.
- Local validation assumes the repo's `justfile` command set and a local SQLite database.
- `just frontend-start` is expected to proxy `/api/**` traffic to the local backend during demo/local smoke use.

## Known Limitations

- No CI/CD pipeline is included in this release-readiness scope.
- No production deployment automation is included in this release-readiness scope.
- Production LINE OAuth integration remains future work; v1 local validation uses the dev-stub baseline.
- v1 does not include v2 grouping/tee-time workflows or v3/v4 notification/history/admin extensions.

## Follow-Up Backlog Boundary

The main implementation-plan backlog after v1 remains:

- v2: grouping and scheduling improvements such as snake/random grouping, tee-time scheduling, drag-and-drop adjustments, WebSocket updates, and richer reservation summaries
- v3: notifications and history, including registration notifications, reminders, grouping announcements, attendance history, and handicap tracking
- v4: management extensions such as dues and team-shirt tracking

Reference: `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`

## Handoff Notes

- Keep release-readiness docs focused on validation and sign-off, not on duplicating detailed smoke steps.
- Keep `demo-smoke-check.md` as the single source of truth for the deterministic dataset and smoke walkthrough.
- If the local command set, auth baseline, or seed dataset changes, update all affected entry-point docs together.
