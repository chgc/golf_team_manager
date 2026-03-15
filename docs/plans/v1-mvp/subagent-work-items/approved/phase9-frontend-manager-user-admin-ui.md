# Subagent Task Proposal

## Basic Information

- Phase: Phase 9 follow-up
- Area: frontend admin UX
- Proposed task name: manager-user-admin-ui
- Related todo id: `phase9-frontend-manager-user-admin-ui`
- Assigned subagent: frontend Angular agent

## Goal

新增 manager-only 的 `/admin/users` 頁面，讓 manager 可以查看 unlinked accounts、管理 player linkage，以及升降權 user role。

## In Scope

- 新增 manager-only route，例如 `/admin/users`
- 新增對應 page / component / service integration
- 呈現 unlinked 與 linked user 管理畫面
- 支援 link to player、unlink、promote to manager、demote to player
- 顯示最小必要欄位：
  - display name
  - auth provider
  - masked provider subject（若採納）
  - current role
  - current player link
  - created / updated timestamp
- 依 backend API 錯誤狀態顯示明確 UX feedback
- 補直接相關的 frontend unit tests

## Out of Scope

- first-manager bootstrap CLI
- backend admin API implementation
- 修改 pending-link 頁面成自助 linking 流程
- 大範圍設計系統重做

## Dependencies

- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-user-admin-api.md`
- `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-user-admin-foundation.md`
- `docs\architecture\auth-line-sso-implementation-detail.md`

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `frontend\src\app\`
  - 直接相關 routing / auth / admin feature 檔案
  - `docs\development\demo-smoke-check.md`
- 預計新增的資料夾 / 檔案：
  - manager user admin feature 相關 Angular 檔案與測試

## Technical Approach

- 使用的技術與模式：
  - manager-only UI 以 parent proposal 鎖定的 admin API contract 為前提
  - 不把 manager 操作塞進 `/auth/pending-link`
  - 若 backend API 尚未實作，可先用 mock/stub 完成頁面骨架與互動
  - 第一版 player selector 直接重用既有 `GET /api/players`
  - merge-ready 整合仍依賴 Admin API slice 已落地
  - 保持既有 auth bootstrap 與 route guard 模式，不引入新的角色系統
  - 此 slice 僅在 parent design 通過 approval 後進入 implementation；實際整合另依賴 API slice
- 依循的規範文件：
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-user-admin-api.md`
  - `docs\architecture\auth-line-sso-implementation-detail.md`
- 是否新增依賴：
  - 不新增第三方依賴

## Risks / Open Questions

- 若現有 app shell 沒有合適 admin 導覽入口，可能需要一併調整 manager navigation
- player 選擇器 UX 需要確認是否使用現有資料來源即可支撐
- provider subject 是否應顯示遮罩後片段，需兼顧可辨識性與資訊暴露

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `just frontend-build`
  - `just frontend-test`
  - manager 路由與 admin page 互動的直接相關測試
- 完成後如何驗收：
  - manager 可以在 UI 內完成 user linkage 與角色管理
  - unlinked player 仍維持既有 `/auth/pending-link` 阻擋流程

## Review Status

- Status: approved
- Reviewer: `agent-6` (explore, claude-haiku-4.5), `agent-7` (general-purpose, claude-sonnet-4.6)
- Review notes: UI ownership、player selector、demo smoke doc ownership、API integration gate 已鎖定
- Agent review summary: UI slice 可在 API contract 固定後部分平行開發，最終整合仍依賴 API slice

## Feedback

- Reviewer agent 1:
  - 建議 frontend 僅在 API contract 固定後部分平行開發，並對齊穩定錯誤碼
- Reviewer agent 2:
  - 指出 demo smoke doc 應由 UI slice 單獨持有，並建議直接重用既有 `GET /api/players`
- Applied proposal updates:
  - 新增 shared foundation 依賴與完整 template 欄位
  - 明確保留 `docs\development\demo-smoke-check.md` 給 UI slice 單獨持有
  - 鎖定第一版 player selector 來源與 UI / API 的整合關係
