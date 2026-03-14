# Subagent Task Proposal

## Basic Information

- Phase: Phase 1
- Area: database bootstrap
- Proposed task name: phase1-db-bootstrap
- Related todo id: `phase1-db-bootstrap`
- Assigned subagent: sqlite migration bootstrap agent

## Goal

建立 SQLite 連線與 migration 基線，使 `backend\` 可在空白環境下初始化資料庫，並為後續正式 schema / repository 實作提供穩定延伸點。

## In Scope

- 選定 pure-Go SQLite driver
- 建立 DB 連線封裝與 lifecycle
- 建立 migration 檔案目錄與命名規則
- 建立 migration runner 或等效初始化流程
- 建立 smoke migration 與測試驗證

## Out of Scope

- 不建立完整 player / session / registration 正式 schema
- 不加入 seed data
- 不引入 ORM library
- 不處理正式部署與備份策略

## Dependencies

- `phase1-workspace-scaffold` must be completed first
- `phase1-backend-bootstrap` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\cmd\api\main.go`
  - `backend\internal\config\`
  - `backend\internal\app\`
  - root `README.md` 或 `docs\development\local-setup.md`（若需補充 DB 初始化指令）
- 預計新增的資料夾 / 檔案：
  - `backend\internal\db\`
  - `backend\migrations\`
  - `backend\data\` 或等效本機資料目錄
  - migration runner 與對應測試檔案

## Technical Approach

- 使用的技術與模式：
  - Go
  - Gin 既有 backend skeleton
  - pure-Go SQLite driver
  - file-based SQL migrations
- 依循的規範文件：
  - `docs\plans\v1-mvp\phase-0-conventions\README.md`
  - `docs\plans\v1-mvp\phase-0-conventions\backend-go-conventions.md`
  - `docs\plans\v1-mvp\phase-1-foundation\README.md`
  - `docs\plans\v1-mvp\phase-1-foundation\subagents\04-sqlite-migrations.md`
- 是否新增依賴：
  - SQLite driver
  - migration helper library（若評估比自寫 runner 更穩定且仍保持簡潔）

## Risks / Open Questions

- 需避免把 migration 與 app startup 綁得過死，導致後續 command / CI 使用不便
- 需確保 Windows 本機環境可穩定建立 SQLite 檔案
- 需與既有 backend config / app lifecycle 對齊，避免重複管理路徑與初始化責任

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - 建立空白 SQLite 檔並成功執行 migration
  - 啟動 backend 並確認 DB initialization / migration 流程可正常工作
  - 確認 `gofmt` 已套用到所有新增 / 修改的 Go 檔案
- 完成後如何驗收：
  - 空白環境可從零建立 SQLite 檔與 baseline schema
  - migration 可重複執行
  - 後續 Phase 2 schema 工作可直接在既有 migration 結構上演進

## Review Status

- Status: approved
- Reviewer: user approval via `lgtm` 
- Review notes: Approved to enter implementation after approval commit.

