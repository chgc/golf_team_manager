# Subagent Task Proposal

## Basic Information

- Phase: Phase 5
- Area: session management
- Proposed task name: sessions-feature
- Related todo id: `sessions-feature`
- Assigned subagent: session management feature agent

## Goal

交付第一個完整的 Manager-facing session feature：場次列表、建立 / 編輯表單、場次詳情、狀態流轉，以及與 backend session API / validation baseline 對齊的前後端流程。

本 work item 明確包含從 foundation baseline 擴充 session detail / update / status management 所需的 backend API 與 frontend page flow，讓 manager 能管理單一場次的完整生命週期。

## In Scope

- 實作場次列表，至少區分 upcoming / history 檢視需求
- 實作建立 / 編輯場次表單
- 實作場次詳情頁
- 實作 manual close / confirm / cancel 等狀態操作
- 補齊 session detail / update flow 所需的 backend routes、repository、service、handler support
- 實作 registration deadline 到期後的 auto-close 檢查機制
- 對齊 max players、deadline、status transition 等驗證規則

## Out of Scope

- 不進入 registration feature 的 player-side 報名 / 請假操作
- 不實作自動分組、tee time 排程、即時同步
- 不產出 reservation report 正式報表內容
- 不在此任務中加入新的 route-level authorization

## Dependencies

- `backend-foundation` must be completed first
- `frontend-shell` must be completed first
- `auth-foundation` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `frontend\src\app\features\sessions\`
  - 視需要更新 architecture / development docs
- 預計新增的資料夾 / 檔案：
  - session detail / form components
  - session repository / service / handler extensions
  - 對應測試與 feature handoff docs

## Technical Approach

- 使用的技術與模式：
  - Angular standalone feature pages + reactive forms
  - Gin handler / service / repository extension on top of backend foundation
  - shared session validation reuse，將 max players、deadline、status transition 規則集中在 backend service / domain validation
  - 以前後端共同維護的 session status contract 表達 `open -> closed -> confirmed -> completed / cancelled`
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\frontend-shell-baseline.md`
  - `docs\architecture\auth-foundation.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- 場次列表需明確定義 upcoming / history 的切分方式，避免前後端各自判斷
- auto-close 檢查機制需與 manual close / confirm / cancel 流程保持一致，避免狀態衝突
- 場次詳情中的即時計算欄位（報名人數、空位、預估組數）需先以目前資料模型可支撐的範圍為主
- 本 feature 不處理 registration UI，但 detail page 的資料結構需支撐後續報名名單整合

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - session feature flow smoke checks
- 完成後如何驗收：
  - manager 可建立 / 編輯 / 檢視 / 關閉 / 確認 / 取消場次
  - deadline 與 status transition validation 一致
  - 場次詳情可顯示核心管理資訊

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
- Agent review summary:

## Feedback

- Reviewer agent 1:
- Reviewer agent 2:
- Applied proposal updates:
