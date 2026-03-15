# Auth LINE SSO Implementation Detail

## Scope

This document translates the external `Golf Team Manager - 09 認證設計` note into a repository-specific implementation detail that fits the current Go + Gin + SQLite + Angular baseline.

## Goals

- keep Angular isolated from direct LINE API calls
- preserve `dev_stub` as the default local-development mode
- introduce a production-ready `line` auth mode without rewriting the existing principal contract
- keep `users` as the canonical authentication source and link business-facing player data through `users.player_id`
- let `/api/auth/me` remain the single frontend bootstrap endpoint for current identity

## Non-Goals

- full production hardening beyond the current repo scope
- multi-provider auth beyond `dev_stub` and `line`
- approval UI for linking newly authenticated users to players
- server-side session storage for app sessions

## Canonical Data Model

### Authentication source of truth

`users` remains the auth source of truth.

Responsibilities:

- `users`
  - auth provider metadata
  - provider subject (`LINE user id` or dev subject)
  - current display name snapshot
  - role used by authorization
  - optional link to `players`
- `players`
  - golf business profile and operational data
  - remains independent from auth provider concerns

This keeps provider identity out of `players` and matches the existing schema in `backend\migrations\000003_create_auth_users.sql`.

### Required user states

Initial implementation should support these states without adding a second auth source:

1. linked manager
   - `users.role = manager`
   - `users.player_id` may be null or linked, depending on whether the manager also has a player profile
2. linked player
   - `users.role = player`
   - `users.player_id` points at the managed player record
3. unlinked authenticated user
   - `users.role = player`
   - `users.player_id IS NULL`
   - authenticated successfully, but cannot perform player-bound operations until linked

The unlinked state replaces the external note's simpler "create player immediately" behavior and fits the current repo structure better.

## Principal and JWT Contract

### Backend principal

The backend should preserve the current `auth.Principal` shape:

```json
{
  "displayName": "王小明",
  "playerId": "player-123",
  "provider": "line",
  "role": "player",
  "subject": "U1234567890abcdef",
  "userId": "user-456"
}
```

Field mapping:

- `displayName` -> `users.display_name`
- `playerId` -> `users.player_id`
- `provider` -> `users.auth_provider`
- `role` -> `users.role`
- `subject` -> `users.provider_subject`
- `userId` -> `users.id`

### App JWT payload

The app JWT should carry the same identity contract used by the principal:

```json
{
  "sub": "user-456",
  "player_id": "player-123",
  "provider": "line",
  "provider_subject": "U1234567890abcdef",
  "role": "player",
  "display_name": "王小明",
  "exp": 1234567890,
  "iat": 1234560000
}
```

Rules:

- `sub` stores the internal user id, not the LINE subject
- `player_id` is optional
- `provider_subject` is required for traceability and parity with the existing principal
- `display_name` is a cache field for UI bootstrap convenience

## Auth Modes

### `AUTH_MODE=dev_stub`

Keep the current middleware-driven behavior:

- requests are authenticated by `DevelopmentAuth`
- `/api/auth/me` returns the injected principal
- frontend can bootstrap immediately without a login screen

### `AUTH_MODE=line`

Switch to JWT-backed auth:

- `/api/auth/line/login` initiates the LINE OAuth redirect
- `/api/auth/line/callback` completes the flow, upserts the `users` record, signs an app JWT, and redirects to frontend `/auth/done#token=...`
- business APIs read the JWT from `Authorization: Bearer <token>`
- `/api/auth/me` derives `auth.Principal` from the validated JWT

The router should choose the auth middleware by mode so frontend code can keep using `/api/auth/me` regardless of environment.

## Backend Endpoints

### `GET /api/auth/me`

- remains the frontend bootstrap endpoint
- returns `401` when no valid JWT is present in `line` mode
- returns `200` with principal in `dev_stub` mode
- returns `200` for an authenticated but unlinked LINE user, with `playerId` omitted

### Unlinked-user contract

For the first LINE auth slice, `playerId` being absent is the only supported signal that the authenticated user is not yet linked to a golf-player record.

Backend contract:

- `/api/auth/me`
  - returns `200`
  - returns a normal `auth.Principal`
  - omits `playerId`
- player-bound business operations
  - return `403`
  - use a stable API error code such as `player_link_required`
- manager-only behavior
  - continues to key off `role = manager`
  - must not assume `playerId` is present

Frontend contract:

- treat `principal.role === 'player' && !principal.playerId` as `authenticated-unlinked`
- allow login completion and shell bootstrap
- redirect feature navigation for unlinked players to a dedicated pending state page or equivalent blocking UX
- keep logout available
- do not treat the missing `playerId` state as an authentication failure

### `GET /api/auth/line/login`

- generate random `state` and `nonce`
- write them into an HttpOnly cookie with short expiration
- redirect to LINE authorize endpoint

### `GET /api/auth/line/callback`

- validate `state` against cookie
- exchange `code` for LINE tokens
- verify the ID token and `nonce`
- resolve or create `users` by `(auth_provider, provider_subject)`
- preserve existing `player_id` if the user record already exists
- default newly created users to:
  - `auth_provider = line`
  - `role = player`
  - `player_id = NULL`
  - `display_name = <LINE display name>`
- sign app JWT
- clear the state cookie
- redirect to `<FRONTEND_URL>/auth/done#token=<app_jwt>`

### `POST /api/auth/refresh`

Deferred from the first implementation slice.

Reason:

- the external note marks it optional
- there is no existing refresh-token storage model in the repo
- initial implementation can use a short-lived access token with explicit re-login

### `POST /api/auth/logout`

Not required for the first slice when app auth is stateless JWT in localStorage.

Initial logout behavior is frontend-only token removal.

## Local Development Host and Cookie Assumptions

The current frontend dev flow uses Angular dev-server proxying `/api/**` to `http://127.0.0.1:8080`.

That proxy is acceptable for normal API traffic, but it is not the default path for starting LINE login in local development because OAuth `state` and `nonce` cookies must round-trip on the backend callback origin.

Local `line` mode assumptions:

- backend public origin: `http://127.0.0.1:8080`
- frontend public origin: `http://localhost:4200`
- LINE callback URI: `http://127.0.0.1:8080/api/auth/line/callback`
- frontend post-login landing page: `http://localhost:4200/auth/done`

Implications:

- the login button in local `line` mode should navigate to the backend origin directly, not a proxied relative `/api/auth/line/login`
- normal authenticated API calls from the SPA may still use the existing `/api/**` proxy after the JWT is stored
- auth docs and local setup must call out this distinction explicitly

If the team later unifies frontend and backend under one public origin, this assumption can be simplified.

## Backend Components

### Config

Extend `backend\internal\config\config.go` with:

- `LINE_CLIENT_ID`
- `LINE_CLIENT_SECRET`
- `LINE_REDIRECT_URI`
- `FRONTEND_URL`
- `JWT_SECRET`
- optional `JWT_TTL`

Rules:

- `line` mode must fail fast when required vars are missing
- `dev_stub` mode must not require LINE or JWT config

### Persistence and services

Expected backend additions:

- LINE OAuth client abstraction
- JWT signer / validator abstraction
- auth user repository or focused auth data access layer
- mode-aware auth middleware
- callback service that maps LINE claims to `users`

### Middleware strategy

- keep `DevelopmentAuth` for `dev_stub`
- add JWT middleware for `line`
- centralize principal injection so handlers continue reading `PrincipalFromContext`

### Unknown-user behavior

Do not auto-create a `players` record in the first auth slice.

Instead:

- create or update the `users` row
- leave `player_id` null until a later linking step or manager action
- expose the missing link through `/api/auth/me`

This avoids creating incomplete golf-player records from OAuth profile data alone.

## Frontend Implementation Detail

### Auth shell responsibilities

`AuthShell` should evolve from a static signal wrapper into the frontend auth state coordinator.

Responsibilities:

- bootstrap current principal through `/api/auth/me`
- expose auth status: loading, authenticated, unauthenticated
- persist and read the JWT from localStorage
- clear invalid or expired tokens
- support development-friendly behavior in `dev_stub`

### Bootstrap and route-render ownership

The frontend must have a single owner for auth bootstrap.

First-slice contract:

- `AuthShell.initialize()` runs once during app startup
- initialization behavior:
  - in `dev_stub`, call `/api/auth/me` immediately and mark the app authenticated on success
  - in `line`, only call `/api/auth/me` when a local JWT exists
- while initialization is in progress, protected routes must not render feature pages yet
- when `/api/auth/me` returns `401` in `line` mode:
  - clear the stored token
  - mark the session unauthenticated
  - let the auth guard redirect to login
- guards must depend on the initialized auth state instead of racing against the first HTTP call

### Required frontend pieces

1. login entry
   - login page or login panel
   - in `line` mode, button redirects to the backend auth origin for `/api/auth/line/login`
2. callback completion route
   - `/auth/done`
   - reads `#token=...`
   - stores token
   - refreshes principal
   - redirects to the default post-login route
3. HTTP interceptor
   - adds `Authorization` header when token exists
4. auth guard
   - protects routes that require an authenticated principal
   - optionally distinguishes manager-only routes later
5. logout action
   - clears local token and auth shell state

### Route behavior

Initial route policy:

- public:
  - login route
  - auth-done route
- authenticated-unlinked player:
  - pending-link route or equivalent blocking page
- authenticated:
  - existing app routes

Manager-only route guards are deferred until manager-only pages become more explicit.

### Expiration handling

The frontend may decode JWT payload to check `exp` for UX decisions, but signature validation stays backend-only.

On expired or invalid token:

- remove local token
- mark session unauthenticated
- redirect to login page when a protected route is visited

## Security Notes

- keep OAuth `state` and `nonce` in HttpOnly cookies
- use `SameSite=Lax`
- use `Secure=true` outside purely local HTTP development
- send the app JWT to the frontend via URL fragment, not query string
- never persist LINE access tokens in the database
- validate JWT signature and expiration on every protected API request

## Validation Expectations

Implementation work derived from this document should validate at least:

- `go test ./...` in `backend\`
- `just frontend-build`
- `just frontend-test`
- manual smoke checks for:
  - `dev_stub` bootstrap via `/api/auth/me`
  - LINE login redirect flow
  - callback success and token persistence
  - unauthenticated redirect behavior

## Planned Work Item Boundaries

This document is intended to feed the following work items:

1. auth implementation detail alignment
2. backend LINE OAuth and JWT
3. frontend auth flow integration
4. auth docs and local-ops updates
