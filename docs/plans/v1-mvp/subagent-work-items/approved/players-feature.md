# Subagent Task Proposal

## Basic Information

- Phase: Phase 4
- Area: player management
- Proposed task name: players-feature
- Related todo id: `players-feature`
- Assigned subagent: player management feature agent

## Goal

交付第一個完整的 Manager-facing feature：球員清單、搜尋 / 篩選、建立 / 編輯表單、active / inactive 狀態處理，以及與 backend player API / validation baseline 對齊的前後端流程。

本 work item 明確包含補齊 player detail / update flow 所需的 backend API，因為目前 foundation 僅有 player list/create route，尚不足以支撐編輯流程。

## In Scope

- 實作球員列表頁
- 實作搜尋與狀態篩選
- 實作新增 / 編輯球員表單
- 實作 active / inactive 狀態切換
- 補齊 player edit flow 所需的 backend `GET /api/players/{playerId}`、`PATCH /api/players/{playerId}` 與相關 repository / service / handler support
- 對齊差點、email、duplicate-name warning 等驗證規則

## Out of Scope

- 不進入 session / registration feature
- 不處理 v3 差點歷史或 v4 隊服 / 隊費欄位
- 不實作複雜權限矩陣 beyond manager / player baseline

## Dependencies

- `backend-foundation` must be completed first
- `frontend-shell` must be completed first
- `auth-foundation` must be completed first

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
  - Gin handler / service / repository extension on top of backend foundation，且本 work item 直接負責補上 `GET /api/players/{playerId}` 與 `PATCH /api/players/{playerId}`
  - shared domain validation reuse，email 以 backend `ValidatePlayerWriteDTO` 為最終權威；invalid email 以 `422 validation_failed` 回傳，frontend 以對應表單驗證同步
  - duplicate-name warning 由 frontend 依目前 player list 做非阻擋式提示，不透過 API error shape 傳遞 warning
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\frontend-shell-baseline.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- duplicate-name 僅做 warning，不做 backend blocking validation；backend 仍允許同名但以 UUID 區分
- duplicate-name warning 由前端在 create/edit form 中即時比對現有 player name，顯示 inline warning 但不阻擋 submit
- inactive player 需保留歷史資料，且未來在 registration feature 中不得出現在新報名選單
- 本 feature 不引入新的 authorization gate；manager-only 屬性先以目前 auth baseline 的 identity context 與 UI scope 表達，不在此任務中補 route-level hard enforcement
- active / inactive 狀態必須可雙向切換，reactivation 為 in-scope 並需有對應測試 / smoke 驗收
- player form 與 list 結構需能支撐後續擴充

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`（既有 quick command）
  - `just frontend-test`（既有 quick command）
  - player feature flow smoke checks
- 完成後如何驗收：
  - manager 可建立 / 編輯 / 篩選球員
  - backend 與 frontend validation 一致，且 API error shape 仍區分 validation / not found / conflict / internal
  - duplicate-name warning 為 warning-only，不阻擋儲存
  - inactive / reactivation 行為明確

## Review Status

- Status: approved
- Reviewer: dual reviewer agents
- Review notes: Auto-approved because two reviewer agents reported no blocking issues.
- Agent review summary: GPT-5.4 approve; Claude Sonnet 4.5 approve.

