# Subagent Task Proposal

## Basic Information

- Phase: Phase 4
- Area: player management
- Proposed task name: players-feature
- Related todo id: `players-feature`
- Assigned subagent: player management feature agent

## Goal

交付第一個完整的 Manager-facing feature：球員清單、搜尋 / 篩選、建立 / 編輯表單、active / inactive 狀態處理，以及與 backend player API / validation baseline 對齊的前後端流程。

## In Scope

- 實作球員列表頁
- 實作搜尋與狀態篩選
- 實作新增 / 編輯球員表單
- 實作 active / inactive 狀態切換
- 對齊差點、email、duplicate-name warning 等驗證規則

## Out of Scope

- 不進入 session / registration feature
- 不處理 v3 差點歷史或 v4 隊服 / 隊費欄位
- 不實作複雜權限矩陣 beyond manager / player baseline

## Dependencies

- `backend-foundation` must be completed first
- `frontend-shell` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `frontend\src\app\features\players\`
  - 視需要更新 schema / architecture / development docs
- 預計新增的資料夾 / 檔案：
  - player list / form components
  - player repository / service / handler extensions
  - 對應測試與 feature handoff docs

## Technical Approach

- 使用的技術與模式：
  - Angular standalone feature pages + reactive forms
  - Gin handler / service / repository extension on top of backend foundation
  - shared domain validation reuse
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\frontend-shell-baseline.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- 需讓 duplicate-name warning 與 strict validation 區分清楚
- 需確保 inactive player 行為與 future registrations flow 一致
- 需讓 player form 與 list 結構可支撐後續擴充

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - player feature flow smoke checks
- 完成後如何驗收：
  - manager 可建立 / 編輯 / 篩選球員
  - backend 與 frontend validation 一致
  - inactive player 行為明確

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
