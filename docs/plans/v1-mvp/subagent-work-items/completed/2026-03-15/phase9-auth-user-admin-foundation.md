# Subagent Task Proposal

## Basic Information

- Phase: Phase 9 follow-up
- Area: auth backend
- Proposed task name: user-admin-foundation
- Related todo id: `phase9-auth-user-admin-foundation`
- Assigned subagent: backend auth / repository agent

## Goal

建立 Phase 9 共用的 backend foundation，集中處理 manager user-admin 功能所需的 repository / service / migration / error contract，避免 bootstrap CLI 與 admin API worktree 同時改同一批核心 backend 檔案。

## In Scope

- 擴充 `users` 相關 repository / service 能力，支援：
  - 查找 user
  - 列出 user admin view 所需資料
  - 原子更新 `role` 與 `player_id`
- 實作 `player_id` 1:1 linking 的資料層策略
- 定義 shared business errors，供 HTTP layer 與 CLI 共用
- 規劃並實作必要 migration / index 變更
- 補直接相關 backend tests

## Out of Scope

- first-manager CLI command 入口
- manager-only HTTP routing / handlers
- frontend admin UI
- audit trail 系統化設計

## Dependencies

- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
- `docs\architecture\auth-line-sso-implementation-detail.md`
- 既有 `users` schema 與 migration baseline

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\migrations\`
  - `backend\internal\repository\user_repository.go`
  - `backend\internal\repository\user_repository_test.go`
- 預計新增的資料夾 / 檔案：
  - `backend\internal\service\user_admin_service.go`
  - `backend\internal\service\user_admin_service_test.go`
  - Phase 9 admin foundation 直接相關 migration helper 檔案

## Technical Approach

- 使用的技術與模式：
  - 在 shared backend slice 集中處理 repository / service / migration / error contract
  - `PATCH` 對應的 role / player link 更新由 foundation 提供單一原子操作
  - 1:1 player linkage 優先由 DB index + service/repository conflict handling 共同保證
  - 此 slice 是 child proposals 中的先行 backend foundation；需在 parent design doc 核可後優先落地
  - service-layer ownership 固定由 foundation slice 持有，提供 API 與 CLI 共用能力
- 依循的規範文件：
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
  - `docs\architecture\auth-line-sso-implementation-detail.md`
  - `WORKFLOW.md`
- 是否新增依賴：
  - 不新增第三方依賴

## Risks / Open Questions

- SQLite partial unique index 與既有本機資料相容性需確認
- shared error 型別要放在 `service`、`auth`，或更聚焦的 package，需要依 repo 既有 pattern 決定

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `Set-Location backend; go test ./...`
- 完成後如何驗收：
  - Bootstrap CLI 與 Admin API 都能依賴同一套 shared user-admin foundation，而不需重複修改核心 repository / migration 檔案

## Review Status

- Status: approved
- Reviewer: `agent-6` (explore, claude-haiku-4.5), `agent-7` (general-purpose, claude-sonnet-4.6)
- Review notes: shared backend ownership 已明確，無剩餘 blocker
- Agent review summary: foundation slice 先行落地 shared repository / service / migration / error contract，之後再支撐 bootstrap 與 API 平行實作

## Feedback

- Reviewer agent 1:
  - 建議把 `player_id` 唯一性、PATCH 原子性、錯誤契約在 child proposals 前先鎖定
- Reviewer agent 2:
  - 指出 bootstrap 與 API proposal 同時修改 repository / auth 核心檔案，應先切出 shared foundation
- Applied proposal updates:
  - 新增 shared backend foundation proposal，作為 bootstrap / API 的先行相依項
  - 固定 foundation 擁有 `user_admin_service.go` / `user_admin_service_test.go`
  - 移除過於寬泛的 `backend\internal\auth\` ownership 宣告
