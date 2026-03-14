# Subagent Task Proposal

## Basic Information

- Phase: Phase 7
- Area: reservation reporting
- Proposed task name: reservation-report
- Related todo id: `reservation-report`
- Assigned subagent: reservation summary feature agent

## Goal

交付第一個可用的球場預約摘要報表功能，讓 manager 能依單一場次輸出 reservation summary，支援球場預約所需的核心名單與統計資訊。

## In Scope

- 補齊 `GET /api/reports/sessions/{sessionId}/reservation-summary` 正式回應內容
- 依單一 session 產出 reservation summary DTO
- 彙整球員名單、已確認報名人數、預估組數與關鍵場次資訊
- 在 frontend 提供 manager 可檢視的 reservation summary 畫面或區塊
- 對齊 session / registration 狀態規則，避免非適用場次輸出錯誤摘要

## Out of Scope

- 不產出 PDF / Excel 匯出
- 不實作自動分組與 tee time 排程
- 不引入外部報表服務

## Dependencies

- `shared-domain-schema` must be completed first
- `backend-foundation` must be completed first
- `sessions-feature` must be completed first
- `registrations-feature` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `frontend\src\app\features\sessions\`
  - 視需要更新 architecture / development docs
- 預計新增的資料夾 / 檔案：
  - reservation summary response / mapper
  - frontend summary view component
  - 對應測試與 feature handoff docs

## Technical Approach

- 使用的技術與模式：
  - Gin handler / service / repository composition on top of session and registration features
  - frontend summary UI integrated into the session management flow
  - summary output 以目前 v1 可得資料為準，不提前混入 v2 grouping / tee time 資訊
- 依循的規範文件：
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\sessions-feature.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- reservation summary 的欄位需與實際球場預約需求保持一致，避免過多或不足
- 若 session / registration 狀態不符輸出條件，需定義清楚 error 行為

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - reservation summary smoke checks
- 完成後如何驗收：
  - manager 可針對單一 session 取得 reservation summary
  - summary 資料與 registration roster / session detail 一致
  - 非法輸出條件具明確 error handling

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
- Agent review summary:

## Feedback

- Reviewer agent 1:
- Reviewer agent 2:
- Applied proposal updates:
