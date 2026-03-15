# Subagent Task Proposal

## Basic Information

- Phase: Phase 9 follow-up
- Area: auth backend
- Proposed task name: bootstrap-first-manager
- Related todo id: `phase9-bootstrap-first-manager`
- Assigned subagent: backend auth / operator tooling agent

## Goal

新增一個 repository 內可操作的 backend CLI / admin command，讓 operator 能把已完成 LINE login 的既有 `users` row 提升為第一位 `manager`，取代手動改 SQLite。

## In Scope

- 新增 backend operator command entry point
- 支援用 `user_id` 或 `(auth_provider, provider_subject)` 查找既有 user
- 將既有 user 提升為 `role=manager`
- 第一版即支援可選的 `--player-id`
- 對找不到 user、重複參數、無效參數提供明確失敗訊息
- 補齊 command 層與資料更新邏輯測試
- 更新直接相關的操作文件

## Out of Scope

- manager-only HTTP API
- frontend admin UI
- 新增第三種角色
- audit trail 系統化設計

## Dependencies

- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-user-admin-foundation.md`
- `docs\architecture\auth-line-sso-implementation-detail.md`
- 既有 `users` schema（`role`, `player_id`, `auth_provider`, `provider_subject`）

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\cmd\admin\`
  - `docs\development\auth-setup.md`
- 預計新增的資料夾 / 檔案：
  - `backend\cmd\admin\` 下的 first-manager command 相關檔案
  - 直接相關測試檔案

## Technical Approach

- 使用的技術與模式：
  - 使用 repository 內的 Go command，而不是公開 HTTP endpoint
  - command path 固定放在 `backend\cmd\admin\`
  - command 僅允許更新既有 `users` row，不自動建立 user
  - command 依賴 shared foundation slice 提供的查詢 / 更新能力，不直接擁有 shared repository 檔案
  - command 對已是 `manager` 的 target 採 idempotent no-op，仍回傳成功
  - 此 slice 僅在 parent design 與 shared foundation 都通過 approval 後進入 implementation
- 依循的規範文件：
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-user-admin-foundation.md`
  - `docs\architecture\auth-line-sso-implementation-detail.md`
  - `docs\development\auth-setup.md`
- 是否新增依賴：
  - 不新增第三方依賴

預期 command contract：

- 成功 promote：exit `0`
- user 已是 manager：exit `0`，回報 no-op
- user 不存在 / player 不存在 / 參數無效：非 `0`
- command 僅在 target user 已完成 LINE login、系統內已有 `users` row 時可成功
- `--player-id` 在第一版即為可選參數；若提供，需與 role 更新同一筆原子完成

## Risks / Open Questions

- operator output 是否需要同時印出更新前 / 更新後狀態，需在 implementation 時依既有 CLI 風格定稿

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `Set-Location backend; go test ./...`
  - command-level smoke，驗證成功 promote 與找不到 user 時的明確失敗
- 完成後如何驗收：
  - operator 不需手動改 SQLite，就能把既有 LINE user 設成第一位 manager

## Review Status

- Status: approved
- Reviewer: `agent-6` (explore, claude-haiku-4.5), `agent-7` (general-purpose, claude-sonnet-4.6)
- Review notes: command path、idempotency、exit behavior、`--player-id` scope 已鎖定
- Agent review summary: bootstrap slice 可在 foundation 合併後，與 API slice 平行進行

## Feedback

- Reviewer agent 1:
  - 要求補 bootstrap command 的 success / failure / idempotency contract
- Reviewer agent 2:
  - 指出此 proposal 不應與 API slice 共享 `user_repository.go` 等核心 backend 檔案 ownership
- Applied proposal updates:
  - 新增對 shared foundation slice 的依賴
  - 移除 shared backend 檔案 ownership，讓此 slice 聚焦 command 入口與操作文件
  - 固定 command path 為 `backend\cmd\admin\`
  - 明確補上 command contract 與 template 欄位
  - 鎖定 `--player-id` 為第一版可選參數
