# Subagent Task Proposal

## Basic Information

- Phase: Phase 6
- Area: registration management
- Proposed task name: registrations-feature
- Related todo id: `registrations-feature`
- Assigned subagent: registration feature agent

## Goal

交付第一個完整的 registration feature：球員報名 / 請假流程、manager 手動調整名單、與 session / player 狀態規則對齊的前後端流程，讓 v1 的核心報名作業可完整運作。

## In Scope

- 實作 player-facing 場次報名 / 取消 / 請假操作
- 實作 manager-facing 報名名單調整流程
- 對齊 inactive player、session closed / confirmed / cancelled 等限制規則
- 整合既有 `GET /api/sessions/{sessionId}/registrations`、`POST /api/sessions/{sessionId}/registrations` 並視需要補齊 update / cancel flow
- 補齊前端 registration feature pages、表單與 detail 整合

## Out of Scope

- 不實作自動分組、tee time 排程、WebSocket 即時同步
- 不產出 reservation report 正式報表
- 不在此任務中導入正式 LINE OAuth

## Dependencies

- `shared-domain-schema` must be completed first
- `backend-foundation` must be completed first
- `frontend-shell` must be completed first
- `auth-foundation` must be completed first
- `players-feature` must be completed first
- `sessions-feature` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `frontend\src\app\features\registrations\`
  - 視需要更新 architecture / development docs
- 預計新增的資料夾 / 檔案：
  - registration feature pages / components
  - registration repository / service / handler extensions
  - 對應測試與 feature handoff docs

## Technical Approach

- 使用的技術與模式：
  - Angular standalone feature pages + reactive forms / signals
  - Gin handler / service / repository extension on top of backend foundation
  - shared registration validation reuse，集中處理 duplicate registration、session capacity、player inactive、session not open 等規則
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\players-feature.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- player-facing 與 manager-facing flow 的 UI 邊界需明確定義，避免單一頁面承載過多操作
- session capacity、inactive player、duplicate registration 等規則需維持前後端一致
- 報名 / 請假 / manager override 的狀態契約需在 proposal 細化後再進入實作

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - registration feature flow smoke checks
- 完成後如何驗收：
  - player 可完成報名 / 取消或請假
  - manager 可查看並調整名單
  - business validation 與 API error shape 維持一致

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
- Agent review summary:

## Feedback

- Reviewer agent 1:
- Reviewer agent 2:
- Applied proposal updates:
