# Repository Structure

## Purpose

This document explains the current repository layout for `golf_team_manager` after the Phase 1 foundation baseline was completed and validated.

## Current Layout

```text
.
├── backend\
│   ├── cmd\
│   │   ├── api\
│   │   └── migrate\
│   ├── data\
│   ├── internal\
│   │   ├── app\
│   │   ├── config\
│   │   ├── db\
│   │   ├── domain\
│   │   ├── http\
│   │   ├── repository\
│   │   └── service\
│   └── migrations\
├── docs\
│   ├── architecture\
│   │   └── repo-structure.md
│   ├── development\
│   │   ├── local-setup.md
│   │   └── phase-1-validation.md
│   └── plans\
│       └── v1-mvp\
├── frontend\
│   ├── src\
│   │   └── app\
│   │       ├── core\
│   │       ├── features\
│   │       └── shared\
│   ├── public\
│   └── angular.json
├── .gitignore
├── WORKFLOW.md
├── justfile
└── README.md
```

## Directory Responsibilities

### `frontend\`

- Angular workspace root
- Managed with `pnpm`
- Generated and extended with Angular CLI
- Uses plain CSS instead of SCSS
- Shared models and feature data-access services now define the frontend shell baseline

### `backend\`

- Go backend root
- Gin-based HTTP service
- SQLite config, connection, and migration baseline
- Shared domain structs, DTOs, and validation live under `internal\domain\`
- API foundation layers now include repository, service, and Gin handler packages
- No ORM usage
- Must follow Google Go style guidance, `gofmt`, and test requirements

### `docs\architecture\`

- Architecture notes and structure-oriented documentation

### `docs\development\`

- Local setup instructions
- Developer workflow notes
- Quick-start guidance
- Validation records for completed phases

### `docs\plans\`

- Planning documents
- Conventions and workflow rules
- Subagent task proposals and lifecycle folders

## Work Item Lifecycle

```text
pending\  ->  approved\  ->  completed\YYYY-MM-DD\
```

- `pending\`: waiting for review
- `approved\`: explicitly approved and committed, ready for implementation
- `completed\YYYY-MM-DD\`: implemented and archived by completion date

## Notes

- Subagent work is designed around `git worktree`
- The root `justfile` is the quick command entry point
- Planning and governance docs define the source-of-truth workflow
- Phase 1 foundation has been validated; Phase 2 should now build on the existing frontend, backend, and SQLite baseline
