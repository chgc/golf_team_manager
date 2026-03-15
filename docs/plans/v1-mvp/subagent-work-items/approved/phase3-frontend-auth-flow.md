# Subagent Task Proposal

## Basic Information

- Phase: Phase 3 follow-up
- Area: frontend auth
- Proposed task name: frontend-auth-flow
- Related todo id: `auth-review-gate`
- Assigned subagent: frontend auth agent

## Goal

在 frontend 補齊 login entry、`/auth/done` callback route、JWT local storage、auth bootstrap、interceptor、guard 與 logout UX，讓 Angular app 能接上 backend 的 LINE SSO flow。

## In Scope

- 補齊 auth shell 狀態管理
- 定義 app startup 的 auth bootstrap / loading ownership
- 新增 login route / view
- 新增 `/auth/done` route / component
- 新增 auth interceptor
- 新增 auth guard
- 定義 unauthenticated `401` redirect 與 authenticated-but-unlinked user UX
- 將既有 app routes 轉為依 auth state 保護
- 為新增 auth flow 補上 Angular 單元測試

## Out of Scope

- 不直接整合 LINE SDK
- 不新增複雜的 account-linking UI
- 不處理 refresh token 續期 UX

## Dependencies

- `phase3-auth-implementation-detail-alignment.md`
- `phase3-backend-line-oauth-jwt.md`（需先完成 callback 與 `/api/auth/me` contract 的實作驗證）

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `frontend\src\app\core\auth\`
  - `frontend\src\app\app.routes.ts`
  - `frontend\src\app\app.ts`
  - `frontend\src\app\shared\models\auth.models.ts`（若 contract 需擴充）
- 預計新增的資料夾 / 檔案：
  - login page / auth-done page / guard / interceptor / tests

## Technical Approach

- 使用的技術與模式：
  - Angular standalone components
  - signal-driven auth shell
  - `HttpInterceptorFn`
  - route guard based on authenticated principal
  - token persistence in localStorage for first slice
- 依循的規範文件：
  - `docs\architecture\auth-line-sso-implementation-detail.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 否；使用 Angular 內建能力完成

## Risks / Open Questions

- 需要避免在 `dev_stub` mode 強制顯示 login 頁面，破壞現有開發體驗
- unknown linked-state user 的 UX 必須與 backend principal contract 一致
- route bootstrap 與 toolbar 顯示需避免初始化閃爍
- local `line` mode login 需直接前往 backend auth origin，而不是沿用相對 `/api` proxy

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `just frontend-build`
  - `just frontend-test`
  - login / callback / logout smoke flow
- 完成後如何驗收：
  - frontend 能接收 JWT、保護路由、透過 `/api/auth/me` 顯示目前 principal
  - `dev_stub` mode 仍可直接進入 app

## Review Status

- Status: approved
- Reviewer: reviewer agent 1, reviewer agent 2
- Review notes: First-pass review requested explicit bootstrap ownership, 401 redirect handling, unlinked-user UX, and stronger backend dependency. Proposal was updated and approved on second pass.
- Agent review summary: Approved by two independent reviewer agents after frontend auth bootstrap and protected-route behavior were clarified.

## Feedback

- Reviewer agent 1: no blocking issues. Proposal is review-ready, dependencies and frontend scope are clear, and the validation plan covers build/test plus auth smoke flow expectations.
- Reviewer agent 2: Blocking: proposal 需補齊 auth bootstrap/loading contract：何時呼叫 /api/auth/me、在結果回來前如何處理 route render、401 時如何導回 login；同時要定義 authenticated-but-unlinked user 的 UX。依賴建議提高到 backend callback + /api/auth/me 行為已實作/驗證，而不只是文件 contract。
- Applied proposal updates: added explicit scope for bootstrap ownership, 401 redirect behavior, unlinked-user UX, strengthened backend dependency, and documented the need to use backend auth origin for local line-mode login.
