# Repository Structure

## Purpose

This document explains the current repository layout for `golf_team_manager` during Phase 1 bootstrap work.

## Current Layout

```text
.
├── backend\
│   └── .gitkeep
├── docs\
│   ├── architecture\
│   │   └── repo-structure.md
│   ├── development\
│   │   └── local-setup.md
│   └── plans\
│       └── v1-mvp\
├── frontend\
│   └── .gitkeep
├── .gitignore
├── justfile
└── README.md
```

## Directory Responsibilities

### `frontend\`

- Angular workspace root
- Managed with `pnpm`
- Generated and extended with Angular CLI
- Uses plain CSS instead of SCSS

### `backend\`

- Go backend root
- Gin-based HTTP service
- No ORM usage
- Must follow Google Go style guidance, `gofmt`, and test requirements

### `docs\architecture\`

- Architecture notes and structure-oriented documentation

### `docs\development\`

- Local setup instructions
- Developer workflow notes
- Quick-start guidance

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
