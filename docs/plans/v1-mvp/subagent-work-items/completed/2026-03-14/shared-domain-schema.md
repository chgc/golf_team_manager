# Subagent Task Proposal

## Basic Information

- Phase: Phase 2
- Area: shared domain schema
- Proposed task name: shared-domain-schema
- Related todo id: `shared-domain-schema`
- Assigned subagent: domain schema and API contract agent

## Goal

將已核准的產品規格轉成可落地的 shared domain baseline，包含 SQLite schema、Go models / DTOs、validation rules 與 API boundary 草稿，讓後續 players / sessions / registrations feature work 能在一致契約上展開。

## In Scope

- 定義 Player / Session / Registration 的資料表與關聯
- 定義 Go models、DTO 與 validation baseline
- 定義主要 REST API boundary 草稿
- 明確保留 future Group(v2) 的 schema extension point
- 補齊相關文件與驗證測試策略

## Out of Scope

- 不實作完整前後端 CRUD 畫面或 handler
- 不完成正式 auth / LINE OAuth 整合
- 不進入 v2 grouping 邏輯

## Dependencies

- `repo-bootstrap` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\migrations\`
  - `backend\internal\`
  - `docs\architecture\`
  - 視需要更新 `README.md` 與開發文件
- 預計新增的資料夾 / 檔案：
  - shared domain schema / API contract 相關文件
  - Go model / DTO / validation baseline 檔案
  - 後續 feature 可重用的 migration / contract baseline

## Technical Approach

- 使用的技術與模式：
  - SQLite schema-first modeling
  - Go model / DTO separation
  - validation rules aligned with business constraints
  - REST API contract drafting before feature implementation
- 依循的規範文件：
  - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
  - `docs\plans\v1-mvp\phase-0-conventions\README.md`
  - `docs\development\phase-1-validation.md`
- 是否新增依賴：
  - 原則上不新增；優先延用現有 Go / SQLite baseline

## Risks / Open Questions

- 需避免在 Phase 2 就過度綁死未來 v2 grouping 模型
- 需明確區分 persistence model、API DTO 與 validation responsibilities
- 需讓 migration 命名與 domain vocabulary 一開始就一致

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - migration apply / re-apply validation
  - contract / schema consistency review
- 完成後如何驗收：
  - Player / Session / Registration baseline schema 明確
  - 後續 feature work 有一致的資料契約可依循
  - 關鍵 business rules 可被實作與測試

## Review Status

- Status: approved
- Reviewer: user approval via `lgtm`
- Review notes: Approved to enter implementation after approval commit.

