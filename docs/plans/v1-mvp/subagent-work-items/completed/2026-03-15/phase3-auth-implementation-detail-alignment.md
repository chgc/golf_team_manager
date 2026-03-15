# Subagent Task Proposal

## Basic Information

- Phase: Phase 3 follow-up
- Area: auth
- Proposed task name: implementation-detail-alignment
- Related todo id: `auth-architecture-detail`
- Assigned subagent: auth architecture agent

## Goal

將 `CKNotepads` 的 `09 認證設計` 轉換成符合本 repo 現況的 implementation detail，明確定義 `users` / `players` 分工、JWT principal contract、`dev_stub` 與 `line` mode 共存策略，以及 unknown LINE user 的 MVP 行為。

## In Scope

- 撰寫 repo 內 auth implementation detail 文件
- 定義 canonical auth data model
- 定義 backend principal 與 JWT payload mapping
- 定義 frontend `/api/auth/me` bootstrap 與 `/auth/done` callback contract
- 定義 authenticated-but-unlinked user 的 backend/frontend contract
- 定義 local `line` mode 的 host/cookie/CORS 假設
- 明確標示第一波不納入的項目

## Out of Scope

- 不直接修改 backend 或 frontend 生產程式碼
- 不建立真正的 LINE OAuth client
- 不完成 refresh token 與 account-linking UI

## Dependencies

- 無；此 task 為其他 auth work items 的上游依據

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `docs\architecture\auth-foundation.md`
  - `docs\plans\v1-mvp\subagent-work-items\pending\`
- 預計新增的資料夾 / 檔案：
  - `docs\architecture\auth-line-sso-implementation-detail.md`

## Technical Approach

- 使用的技術與模式：
  - 文件導向 architecture alignment
  - 以既有 `users` schema 為 auth source
  - 以現有 `auth.Principal` contract 為前後端共用 identity contract
- 依循的規範文件：
  - `WORKFLOW.md`
  - `docs\plans\v1-mvp\subagent-work-items\README.md`
  - `docs\architecture\auth-foundation.md`
  - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
- 是否新增依賴：
  - 否

## Risks / Open Questions

- 若 implementation detail 沒先定義 unknown LINE user 行為，後續 backend/frontend 會各自假設不同流程
- 需要避免外部文件中的 `players.line_uid` 思路與 repo 現有 `users` schema 衝突
- local `line` mode 若仍沿用 dev proxy 啟動 login，state cookie 與 callback origin 可能不一致

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - 文件交叉比對，確認與目前程式碼、schema、workflow 一致
- 完成後如何驗收：
  - 後續 backend/frontend/auth-ops proposals 都能直接引用此文件
  - canonical model、JWT contract、mode strategy 明確可執行

## Review Status

- Status: approved
- Reviewer: reviewer agent 1, reviewer agent 2
- Review notes: First-pass blocking feedback was incorporated into the implementation detail and proposal scope. Second-pass review confirmed no blocking issues remain.
- Agent review summary: Approved by two independent reviewer agents after clarifying the unlinked-user contract and local line-mode host/cookie assumptions.

## Feedback

- Reviewer agent 1: no blocking issues. Proposal matches the workflow/template, scope is clear, dependencies are explicit, and validation is sufficient for a docs-only architecture-alignment task.
- Reviewer agent 2: Blocking: 請把「已登入但尚未 linked 的 user（player_id = NULL）」契約寫死，包含 backend 回應語意 / HTTP status / frontend UX；另外也要明確定義 local LINE mode 的 host/cookie/CORS 假設。這個 upstream 文件若不先鎖定，backend/frontend 會各自補不同流程。
- Applied proposal updates: added explicit scope for unlinked-user contract and local line-mode host/cookie assumptions; updated the implementation detail document to lock those decisions before downstream work starts.
