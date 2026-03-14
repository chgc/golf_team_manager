# Auth Foundation

## Scope

This document records the current authentication baseline for development and future LINE OAuth compatibility.

## Backend Baseline

- `AUTH_MODE` defaults to `dev_stub`
- request middleware injects a development principal on every request
- `/api/auth/me` returns the current principal
- debug headers can override the development identity:
  - `X-Debug-Display-Name`
  - `X-Debug-Role`
  - `X-Debug-Subject`
  - `X-Debug-Player-ID`

## Schema Baseline

The `users` table now reserves:

- role (`manager` / `player`)
- auth provider (`dev_stub` / `line`)
- provider subject
- optional `player_id` link

## Frontend Baseline

- `AuthShell` exposes signal-based identity state
- `AuthApi` reserves the HTTP boundary for `/api/auth/me`
- the app toolbar shows the current development identity

## Future Direction

- replace the dev-stub middleware with a real LINE OAuth integration
- bind authenticated users to player records where appropriate
- add route or action-level authorization once feature flows expand
