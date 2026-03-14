# Subagent Task Proposal

## Basic Information

- Phase: Phase 1
- Area: backend bootstrap
- Proposed task name: phase1-backend-bootstrap
- Related todo id: `phase1-backend-bootstrap`
- Assigned subagent: backend bootstrap agent

## Goal

建立 Go + Gin 後端骨架，使 `backend\` 具備可啟動、可測試、可承接 SQLite 與後續 REST API 的基礎能力，並符合 Google Go style guide、`gofmt` 與測試要求。

## In Scope

- 初始化 Go module
- 建立 `cmd` / `internal` 結構
- 建立 Gin router、基本 middleware 與 app startup
- 提供 `/health` 或等效健康檢查 endpoint
- 建立 config 與 package layout 基線
- 保留後續 DB / auth / handlers 的乾淨注入點

## Out of Scope

- 不實作完整 player / session / registration CRUD
- 不正式串接 OAuth
- 不做 WebSocket
- 不引入 ORM library
- 不在此階段建立完整 SQLite schema

## Dependencies

- `phase1-workspace-scaffold` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\` 內 Go 專案檔案
  - 視需要更新 root `README.md` 的 backend 啟動指令說明
- 預計新增的資料夾 / 檔案：
  - `backend\cmd\api\main.go`
  - `backend\internal\app\`
  - `backend\internal\config\`
  - `backend\internal\http\`
  - `backend\internal\middleware\`
  - Go module 與 baseline test 檔案

## Technical Approach

- 使用 Gin 作為 HTTP framework
- 不使用 ORM，保留 raw SQL / repository 路徑
- 將啟動組裝邏輯集中在 `internal`
- 讓 `/health` 與基礎測試先建立 smoke-level 可運作基線
- 依 Google Go style guide 實作並於編輯後執行 `gofmt`

- 使用的技術與模式：
  - Go
  - Gin
  - cmd/internal project layout
  - raw SQL-friendly architecture
- 依循的規範文件：
  - `docs\plans\v1-mvp\phase-0-conventions\README.md`
  - `docs\plans\v1-mvp\phase-0-conventions\backend-go-conventions.md`
  - `docs\plans\v1-mvp\phase-1-foundation\README.md`
  - `docs\plans\v1-mvp\phase-1-foundation\subagents\03-backend-bootstrap.md`
- 是否新增依賴：
  - Gin
  - 不引入 ORM library

## Risks / Open Questions

- 需在保持最小骨架的前提下，為後續 SQLite migration 預留清楚擴充點
- 需避免過早把 transport、config、DB lifecycle 寫死
- 需確保 baseline 測試與專案結構足以支撐後續演進

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - 啟動服務並確認 `/health` 可回應
  - 確認 `gofmt` 已套用到所有新增 / 修改的 Go 檔案
- 完成後如何驗收：
  - `backend` 可成功啟動
  - 有健康檢查端點
  - 結構清楚且可承接 DB / auth / handlers 擴充

## Review Status

- Status: approved
- Reviewer: user approval via `phase1-backend-bootstrap.md LGTM`
- Review notes: Approved to move into the implementation-ready stage.

