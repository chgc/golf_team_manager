# Subagent Task Proposal

## Basic Information

- Phase: Phase 1
- Area: workspace scaffold
- Proposed task name: phase1-workspace-scaffold
- Related todo id: `phase1-workspace-scaffold`
- Assigned subagent: workspace scaffold agent

## Goal

建立 repo 的第一層骨架與共用開發基線，讓後續 frontend、backend、docs 與其他 subagent 可以在一致、可維護、相容於 `git worktree` 的結構上展開工作。

## In Scope

- 建立 top-level 資料夾：`frontend\`、`backend\`、`docs\architecture\`、`docs\development\`
- 建立 root `.gitignore`
- 建立 root `README.md`
- 補充 repo 層級的開發說明，包含：
  - frontend 使用 `pnpm`
  - subagent 使用 `git worktree`
  - frontend worktree 透過 pnpm 共用 `node_modules`
- 定義 root 層級的協作與結構原則，避免後續子任務各自發散

## Out of Scope

- 不初始化 Angular workspace
- 不初始化 Go module
- 不建立 SQLite migration
- 不實作任何產品功能
- 不決定 API 契約細節

## Dependencies

- None

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `README.md`
  - `.gitignore`
- 預計新增的資料夾 / 檔案：
  - `frontend\`
  - `backend\`
  - `docs\architecture\`
  - `docs\development\`
  - 視需要新增 root 級 workflow 說明檔，但優先整合進 `README.md`

## Technical Approach

- 使用最小化 root scaffold，先確保後續 Angular CLI、pnpm、Gin/Go module 與 SQLite 相關工作不會互相卡住
- `README.md` 負責說明 repo 結構、主要工作流與各子目錄用途
- `.gitignore` 需涵蓋：
  - Windows 暫存檔
  - Node / pnpm 常見產物
  - Go build 輸出
  - SQLite / local data 暫存檔
  - worktree 開發中常見不應進版控的檔案
- 結構上需保留未來可在不同 worktree 中並行建立 `frontend` / `backend` / docs 的彈性

- 使用的技術與模式：
  - Git repository top-level scaffold
  - pnpm-friendly frontend workflow
  - git worktree-friendly collaboration layout
- 依循的規範文件：
  - `docs\plans\v1-mvp\phase-0-conventions\README.md`
  - `docs\plans\v1-mvp\phase-0-conventions\shared-engineering-conventions.md`
  - `docs\plans\v1-mvp\phase-1-foundation\README.md`
  - `docs\plans\v1-mvp\phase-1-foundation\subagents\01-workspace-scaffold.md`
- 是否新增依賴：
  - 否

## Risks / Open Questions

- 若 root workflow 說明寫得太簡略，後續 subagent 可能仍會各自補自己的做法，造成重複
- 若 `.gitignore` 不完整，後續 Angular / Go 初始化後容易產生雜訊檔案
- 此階段先不決定 monorepo 工具，只建立最小可延伸結構

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - 檢查目錄結構是否正確建立
  - 檢查 `README.md` 是否清楚說明 repo 結構與工作流
  - 檢查 `.gitignore` 是否覆蓋主要暫存與輸出檔案類型
  - 檢查目前輸出是否與後續 Angular CLI / pnpm / Go module 初始化相容
- 完成後如何驗收：
  - 後續 `phase1-frontend-bootstrap` 與 `phase1-backend-bootstrap` 不需重排 root 結構即可直接接手
  - root 文件足以作為 worktree 協作與 frontend pnpm workflow 的共同入口

## Review Status

- Status: approved
- Reviewer: user approval via `go`
- Review notes: Proceed with the workspace scaffold implementation.
