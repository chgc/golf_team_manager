# Subagent Task Proposal

## Basic Information

- Phase: Phase 9 follow-up
- Area: auth backend
- Proposed task name: manager-user-admin-api
- Related todo id: `phase9-manager-user-admin-api`
- Assigned subagent: backend auth / API agent

## Goal

新增 manager-only user administration API，讓 manager 可以列出已登入使用者、管理 `users.player_id` linkage，以及升降權 `users.role`，不再依賴手動 SQL。

## In Scope

- 新增 `GET /api/admin/users`
- 新增 `PATCH /api/admin/users/:userId`
- 實作 manager-only authorization
- 支援 `linkState` 與 `role` filter
- 支援更新 `role` 與 `playerId`
- 實作必要 guardrails，例如：
  - `playerId` 必須指向既有 player
  - 最後一位 manager 不可被降權
  - 若採納 proposal 決策，單一 `player` 同時只綁定單一 active auth user
- 補 router / handler / authorization tests
- 更新直接相關的 auth / operation 文件

## Out of Scope

- first-manager bootstrap CLI
- frontend admin UI
- 刪除 auth user
- 完整 audit trail 或細緻 RBAC

## Dependencies

- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-user-admin-foundation.md`
- `docs\architecture\auth-line-sso-implementation-detail.md`
- `docs\development\auth-setup.md`

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\http\router.go`
  - `backend\internal\http\`
  - `backend\internal\http\responses.go`
- 預計新增的資料夾 / 檔案：
  - `backend\internal\http\admin_users_handlers.go`
  - `backend\internal\http\admin_users_handlers_test.go`

## Technical Approach

- 使用的技術與模式：
  - 依 parent proposal 鎖定的 contract，採用最小 admin API surface：
    - `GET /api/admin/users`
    - `PATCH /api/admin/users/:userId`
  - `PATCH` 允許同時更新 `role` 與 `playerId`，並要求單一 transaction 原子完成
  - 僅允許更新 `users.role` 與 `users.player_id`
  - backend 明確回傳 authorization 與 business-rule 失敗，不做 silent fallback
  - manager-only behavior 仍以既有 technical role `manager` 為準
  - 第一版 player selector 不新增新 endpoint，直接重用既有 `GET /api/players`
  - 此 slice 僅在 parent design 與 shared foundation 都通過 approval 後進入 implementation
  - API slice 不擁有 service-layer 檔案；shared service 由 foundation slice 提供
- 依循的規範文件：
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-user-admin-foundation.md`
  - `docs\architecture\auth-line-sso-implementation-detail.md`
  - `docs\development\auth-setup.md`
- 是否新增依賴：
  - 不新增第三方依賴

錯誤回應契約：

- 維持既有 `ErrorResponse` shape：`error.code`, `error.message`, `error.details`
- 此 slice 需補的穩定錯誤碼至少包含：
  - `last_manager_demotion_forbidden`
  - `player_already_linked`
  - `player_not_found`
  - `user_not_found`

## Risks / Open Questions

- 若 manager admin route 需要新增 shell/navigation surfaced behavior，需與 frontend slice 對齊導覽入口命名

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `Set-Location backend; go test ./...`
  - manager / player authorization tests
  - linked / unlinked / promote / demote API cases
- 完成後如何驗收：
  - manager 能透過受保護 API 管理 user role 與 player linkage
  - 最後一位 manager 降權等 guardrail 由 backend 強制保證

## Review Status

- Status: approved
- Reviewer: `agent-6` (explore, claude-haiku-4.5), `agent-7` (general-purpose, claude-sonnet-4.6)
- Review notes: HTTP ownership、PATCH contract、error contract、player selector 來源 已鎖定
- Agent review summary: API slice 依賴 foundation，並可在 foundation 後與 bootstrap slice 平行開發

## Feedback

- Reviewer agent 1:
  - 要求先鎖定 PATCH 原子更新與最後一位 manager 的錯誤契約
- Reviewer agent 2:
  - 指出此 proposal 不應與 bootstrap slice 共享 shared backend repository / docs ownership，並建議明確重用既有 `GET /api/players`
- Applied proposal updates:
  - 新增 shared foundation slice 依賴，移除 shared backend repository ownership
  - 鎖定 PATCH contract、錯誤回應 shape 與 player selector 來源
  - 移除與其他 slice 會衝突的文件 ownership
  - 固定 API slice 僅擁有 HTTP router / handler / handler test 檔案
