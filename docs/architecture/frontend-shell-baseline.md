# Frontend Shell Baseline

## Scope

This document records the current Angular shell baseline after aligning frontend structure with the shared domain schema and backend API foundation.

## Route Baseline

- `/` → home landing shell
- `/players` → player feature landing page
- `/sessions` → session feature landing page
- `/registrations` → registration feature landing page

All feature routes are lazy loaded with standalone components.

## Current Frontend Boundaries

### Shared models

- `src/app/shared/models/domain.models.ts`

The shared model file mirrors the backend DTO vocabulary for:

- players
- sessions
- registrations

### Data-access services

- `PlayersApi`
- `SessionsApi`
- `RegistrationsApi`

These services use `HttpClient` and provide the first typed HTTP boundary for future feature work.

### Shell responsibilities

- keep navigation and route structure stable
- provide feature landing pages without prematurely implementing full UI workflows
- preserve Angular CLI / standalone / plain CSS constraints

## Current Gaps

- no auth state or role-aware navigation yet
- no reactive forms yet for create/edit workflows
- no list/table or detail-page UX yet

## Next Recommended Work

- `auth-foundation`
