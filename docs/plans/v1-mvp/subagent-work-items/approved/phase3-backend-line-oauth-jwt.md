# Subagent Task Proposal

## Basic Information

- Phase: Phase 3 follow-up
- Area: backend auth
- Proposed task name: line-oauth-jwt
- Related todo id: `auth-review-gate`
- Assigned subagent: backend auth agent

## Goal

在 backend 落地 LINE OAuth login/callback flow、auth user upsert、app JWT 簽發與驗證 middleware，並維持 `dev_stub` 與 `line` mode 可切換。

## In Scope

- 新增 `line` mode 所需 config 與 fail-fast 驗證
- 新增 `GET /api/auth/line/login`
- 新增 `GET /api/auth/line/callback`
- 新增 JWT signer / validator abstraction
- 新增 router / auth middleware 依 `AUTH_MODE` 切換邏輯
- 新增 JWT auth middleware 並維持 `/api/auth/me`
- 明確定義並測試 `line` mode 下 `/api/auth/me` 的 `401` / unlinked-user 行為
- 修正或正式化 dev stub principal 中 `UserID` 與 `Subject` 的 mapping 行為
- 為 LINE auth flow、新 middleware 與 config 加上 Go 測試

## Out of Scope

- 不完成 frontend login UX
- 不實作 refresh token storage 與 rotation
- 不實作 manager UI 來連結未綁定 user 與 player

## Dependencies

- `phase3-auth-implementation-detail-alignment.md`

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\config\`
  - `backend\internal\http\`
  - `backend\internal\http\middleware\`
  - `backend\internal\auth\`
  - `backend\migrations\`（僅在確認 schema 補欄位需要時）
- 預計新增的資料夾 / 檔案：
  - backend auth service / JWT helper / LINE client abstraction
  - 對應測試檔

## Technical Approach

- 使用的技術與模式：
  - Gin handlers + middleware
  - stateless app JWT with bearer auth
  - provider adapter pattern for LINE token exchange / verification
  - `users` table as canonical auth source
- 依循的規範文件：
  - `docs\architecture\auth-line-sso-implementation-detail.md`
  - `docs\architecture\auth-foundation.md`
  - `docs\plans\v1-mvp\phase-0-conventions\backend-go-conventions.md`
- 是否新增依賴：
  - 盡量避免；如需 JWT library，需在 proposal review 中明確說明

## Risks / Open Questions

- 需要明確處理 `dev_stub` 與 `line` mode router wiring，避免破壞現有 smoke flow
- LINE callback 對 unknown user 的行為必須與 implementation detail 完全一致
- 若新增 JWT library，需確認與現有 backend 依賴策略相容
- 目前 migration 已具備主要 auth user 欄位，若提案需要 schema 變更，必須先證明是 implementation detail 明確要求而非預設補強

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `Set-Location backend; go test ./...`
  - local login redirect / callback smoke check
  - `/api/auth/me` in both `dev_stub` and `line` mode
- 完成後如何驗收：
  - backend 能在 `line` mode 完成登入、簽發 JWT、驗證 JWT 並回傳 principal
  - `dev_stub` mode 維持可用

## Review Status

- Status: approved
- Reviewer: reviewer agent 1, reviewer agent 2
- Review notes: First-pass review requested explicit AUTH_MODE router wiring, `/api/auth/me` semantics in `line` mode, and dev-stub principal mapping coverage. Proposal was updated and approved on second pass.
- Agent review summary: Approved by two independent reviewer agents after backend auth-mode behavior and migration expectations were made explicit.

## Feedback

- Reviewer agent 1: no blocking issues. Proposal is review-ready with clear backend scope, explicit dependency on the implementation-detail task, concrete validation, and acceptable out-of-scope boundaries.
- Reviewer agent 2: Blocking: proposal 需明確包含 AUTH_MODE 對應的 router/middleware 切換與 line mode 下 /api/auth/me 的 401 行為，否則會和目前 router 永遠套 DevelopmentAuth 的實作不一致。另外請把現有 dev stub principal 中 UserID=Subject 的基線問題納入處理/測試；users schema 目前已具備必要欄位，migration 應視為非預期而非預設。
- Applied proposal updates: expanded scope to include AUTH_MODE router wiring, explicit `/api/auth/me` line-mode behavior, dev-stub principal mapping review, and clarified that schema migration is exception-based rather than assumed.
