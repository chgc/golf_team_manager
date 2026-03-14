# Subagent Task Proposal

## Basic Information

- Phase: Phase 1
- Area: integration validation
- Proposed task name: phase1-integration-check
- Related todo id: `phase1-integration-check`
- Assigned subagent: phase 1 integration validation agent

## Goal

對已完成的 workspace、frontend、backend、sqlite migration 基線做一次整合驗收，確認 Phase 1 交付物在同一 repo 中可一致運作，並補齊交接文件中的驗證結果與已知限制。

## In Scope

- 驗證 frontend build / test
- 驗證 backend test / migrate / startup health check
- 驗證 root `justfile` 快速指令與文件一致
- 盤點 Phase 1 已完成輸出與 remaining gaps
- 補齊整合驗收結論與 handoff notes

## Out of Scope

- 不新增 Phase 2 正式業務 schema
- 不實作 players / sessions / registrations feature
- 不引入新的基礎框架或重構已完成的 bootstrap 結構

## Dependencies

- `phase1-frontend-bootstrap` must be completed first
- `phase1-backend-bootstrap` must be completed first
- `phase1-db-bootstrap` must be completed first
- `phase1-docs-devx` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `README.md`
  - `docs\development\local-setup.md`
  - `docs\architecture\repo-structure.md`
  - 視需要更新 `justfile`
- 預計新增的資料夾 / 檔案：
  - 視需要新增整合驗收紀錄文件

## Technical Approach

- 使用的技術與模式：
  - root `justfile` workflow validation
  - Angular CLI / pnpm validation
  - Go / Gin / SQLite validation
  - docs-to-implementation consistency check
- 依循的規範文件：
  - `docs\plans\v1-mvp\phase-0-conventions\README.md`
  - `docs\plans\v1-mvp\phase-1-foundation\README.md`
  - `docs\development\local-setup.md`
  - `docs\architecture\repo-structure.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- 需避免把整合驗收變成 Phase 2 規劃或功能開發
- 需確保文件中的命令與實際 repo 狀態完全一致
- 若發現 cross-workstream 缺口，需只補最小必要修正

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `just frontend-build`
  - `just frontend-test`
  - `just backend-test`
  - `just backend-migrate`
  - 啟動 backend 並確認 `/health`
- 完成後如何驗收：
  - Phase 1 foundation outputs 可被完整驗證
  - 文件、指令、實際 repo 結構一致
  - Phase 2 可在清楚的 foundation baseline 上開始

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
