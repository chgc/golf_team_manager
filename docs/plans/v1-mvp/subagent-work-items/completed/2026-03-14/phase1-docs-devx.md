# Subagent Task Proposal

## Basic Information

- Phase: Phase 1
- Area: docs and dev experience
- Proposed task name: phase1-docs-devx
- Related todo id: `phase1-docs-devx`
- Assigned subagent: docs and devx agent

## Goal

補齊 Phase 1 所需的開發文件與 root 快速啟動流程，讓後續 subagent 或開發者可以用一致的方式理解 repo 結構、啟動本地環境，並透過 `just` 快速執行常用指令。

## In Scope

- 建立 `docs\development\local-setup.md`
- 建立 `docs\architecture\repo-structure.md`
- 在 root 建立 `justfile`
- 將常用開發命令整理成 `just` recipes
- 視需要更新 root `README.md`，加入 `just` 的使用方式
- 記錄目前的 pnpm / git worktree / Gin / SQLite 工作流

## Out of Scope

- 不初始化 Angular workspace
- 不初始化 Go module
- 不實作產品功能
- 不加入與目前 repo 狀態無關的複雜自動化

## Dependencies

- `phase1-workspace-scaffold` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `README.md`
- 預計新增的資料夾 / 檔案：
  - `docs\development\local-setup.md`
  - `docs\architecture\repo-structure.md`
  - `justfile`

## Technical Approach

- 先以目前 repo 的真實狀態為基礎，整理 root scaffold、docs 佈局與下一階段工作流
- `justfile` 僅提供與當前階段相符的快速啟動或檢查命令，避免先寫入尚未存在的實作命令
- 命令命名保持簡潔，優先涵蓋：
  - 檢查 repo 狀態
  - 查看計畫 / 文件路徑
  - 後續可擴充前端與後端啟動入口

- 使用的技術與模式：
  - Markdown docs
  - root `justfile`
  - Windows-friendly command conventions
- 依循的規範文件：
  - `docs\plans\v1-mvp\phase-0-conventions\README.md`
  - `docs\plans\v1-mvp\phase-0-conventions\shared-engineering-conventions.md`
  - `docs\plans\v1-mvp\phase-1-foundation\README.md`
  - `docs\plans\v1-mvp\phase-1-foundation\subagents\05-docs-and-devx.md`
- 是否新增依賴：
  - 不新增專案 runtime 依賴
  - `just` 為使用者環境工具，假設可於開發機上安裝 / 使用

## Risks / Open Questions

- 在 frontend / backend 尚未完整初始化前，`justfile` 的 recipes 需避免引用不存在的命令
- 文件若寫得過早過細，後續可能仍需微調
- 需兼顧 Windows 本機工作流

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - 檢查 docs 路徑與 repo 現況一致
  - 檢查 `justfile` recipes 至少可執行當前存在的命令
  - 檢查 README 與 docs 對 `just` 的說明一致
- 完成後如何驗收：
  - 新接手者可以從 root README 與 `justfile` 知道如何開始
  - docs 能正確描述當前 repo 狀態與下一步方向

## Review Status

- Status: approved
- Reviewer: user approval via `proposal LGTM`
- Review notes: Approved to move into the implementation-ready stage.
