# Subagent Task Proposal

## Basic Information

- Phase: Phase 8
- Area: quality assurance and seed data
- Proposed task name: qa-and-seed-data
- Related todo id: `qa-and-seed-data`
- Assigned subagent: qa and seed data agent

## Goal

補齊 v1 MVP demo / 試用前的品質保證與種子資料基線，讓團隊可快速建立可展示環境，並用一致的 smoke path 驗證 players / sessions / registrations / reservation report 主流程。

## In Scope

- 補齊 v1 關鍵流程缺少的 backend / frontend 測試
- 提供可重複建立的 seed data（球員、場次、報名）
- 補齊 demo / smoke check 所需的最小資料初始化流程
- 盤點並補強主要頁面的 loading、empty、error state
- 更新本地開發與 demo 操作文件

## Out of Scope

- 不實作新的 v2 grouping / tee time / websocket 功能
- 不引入外部測試平台或 e2e framework
- 不處理正式 production deployment / CI pipeline

## Dependencies

- `shared-domain-schema` must be completed first
- `backend-foundation` must be completed first
- `frontend-shell` must be completed first
- `auth-foundation` must be completed first
- `players-feature` must be completed first
- `sessions-feature` must be completed first
- `registrations-feature` must be completed first
- `reservation-report` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `backend\cmd\`
  - `frontend\src\app\`
  - `docs\development\`
  - `README.md`
  - 視需要更新 architecture / handoff docs
- 預計新增的資料夾 / 檔案：
  - seed data 相關 bootstrap / command 或 script
  - 補強後的測試檔案
  - demo / smoke check 操作文件
- 後端具名交付：
  - seed data 建立入口（命名待 review 收斂）
  - 補強 players / sessions / registrations / reports handler / service / repository 測試
- 前端具名交付：
  - 補強 `SessionListPage` 與 reservation summary 相關互動 / state 測試
  - 補強首頁或主要 feature page 的 loading / empty / error UX

## Technical Approach

- 使用的技術與模式：
  - 優先沿用既有 Go test、Angular unit test 與 root `justfile` 指令，不新增新的測試框架
  - seed data 以 repo 內可重複執行的本地 bootstrap 方式提供，避免手動 SQL
  - 測試與 seed data 內容只覆蓋 v1 已完成功能，不提前混入 v2+ backlog
  - loading / empty / error state 優先補在既有 feature page 與 data-access flow，不另開新頁面
- 依循的規範文件：
  - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\players-feature.md`
  - `docs\architecture\sessions-feature.md`
  - `docs\architecture\registrations-feature.md`
  - 後續需補入 `docs\architecture\reservation-report.md`（待前一項完成）
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- seed data 內容需在「可 demo」與「不污染正式資料」之間取得平衡，需明確限定為 local/dev 用途
- 若現有頁面已有零散 error handling，需先盤點一致的 UX 邊界，避免補強後造成不同頁面行為不一致
- seed data 建立入口名稱與放置位置需先定義清楚，避免與 migration / app startup 職責混淆

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - seed data bootstrap smoke checks
  - end-to-end manual smoke path covering players -> sessions -> registrations -> reservation report
- 完成後如何驗收：
  - Case 1: 全 repo 既有 build / test 指令保持通過
  - Case 2: 本地可透過文件化流程快速建立 demo seed data
  - Case 3: manager / player 主流程皆可在 seed data 上重現
  - Case 4: 關鍵頁面的 loading / empty / error state 有對應測試或明確 smoke 驗證
  - Case 5: demo / local setup 文件可讓新進成員完成最小可用環境

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
- Agent review summary:

## Feedback

- Reviewer agent 1:
- Reviewer agent 2:
- Applied proposal updates:
