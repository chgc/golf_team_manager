# Subagent Task Proposal

## Basic Information

- Phase: Phase 9 follow-up
- Area: auth + player-management
- Proposed task name: manager-player-linking
- Related todo id: `draft-admin-design-doc`
- Assigned subagent: auth / operations design agent

## Goal

定義一份可 review 的設計文件，回答目前 LINE login 完成後，系統要如何「設定第一個 admin / manager」以及後續如何由 manager 在系統內管理 player linkage 與使用者權限，而不再依賴手動修改 SQLite 的 `users.player_id` / `users.role`。

這份文件的重點不是直接實作，而是先把 bootstrap、權限模型、API / UI 範圍、風險與 rollout 順序鎖定，讓後續實作 work item 可以直接承接。

## Current State Summary

- repo 目前已經有兩種角色：`manager` 與 `player`
- `users` 是 auth source of truth；球隊業務資料仍在 `players`
- LINE 新登入使用者會自動建立 `users` row，預設 `role=player` 且 `player_id=NULL`
- frontend 會把 `role=player` 且沒有 `playerId` 的使用者導向 `/auth/pending-link`
- 目前沒有 manager-linking UI 或 API；linked-user smoke 仍需手動更新 `users.player_id`
- `manager` 權限目前已被 backend / frontend 使用，例如保留給 manager 的報表與 manager-only 操作

## In Scope

- 明確定義「admin」在本 repo 中對應的授權模型
- 設計第一個 manager 的 bootstrap 流程
- 設計 manager 在後台管理「已登入使用者」與「player record link」的操作模型
- 設計是否允許 manager 升降權其他使用者為 manager
- 定義推薦的 backend API 邊界、frontend 頁面邊界、資料更新規則與驗收方式
- 定義 rollout 順序，讓功能可分階段落地

## Out of Scope

- 本 proposal 不直接修改 production code
- 不新增第三種正式角色（例如 `super_admin`）到現有 schema
- 不導入完整 IAM / RBAC 系統
- 不處理多球隊 / 多租戶 / organization-scoped admin
- 不做 LINE provider 以外的邀請制 onboarding

## Dependencies

- `docs\architecture\auth-line-sso-implementation-detail.md`
- `docs\development\auth-setup.md`
- `docs\development\demo-smoke-check.md`
- 既有 `users` schema（`role`, `player_id`, `auth_provider`, `provider_subject`）

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `docs\plans\v1-mvp\subagent-work-items\pending\`
- 預計新增的資料夾 / 檔案：
  - `docs\plans\v1-mvp\subagent-work-items\pending\phase9-auth-manager-player-linking.md`

## Technical Approach

### 1. Terminology: use product label "Admin", keep technical role `manager`

建議在這個 repo 內先不要新增 `admin` 角色值。

原因：

- 目前 schema、principal contract、backend authorization、frontend labels 都已經以 `manager` / `player` 為基線
- 本次需求真正缺的是「如何把某個 LINE user 設為 manager」與「manager 如何管理 player linkage」，不是新的權限層級
- 若直接新增 `admin` role，會擴大 schema、JWT、middleware、frontend guard 與 seed / test 的修改面

推薦做法：

- product / UX 層面可以把 manager 文案顯示為 `Admin` 或 `Manager`
- backend / DB / JWT 仍維持 `role = manager`
- 後續若真的需要區分 operator 與 super-admin，再另開獨立 proposal 處理

### 2. Bootstrap the first manager explicitly

目前最大的 gap 不是一般 player link，而是「系統裡第一個 manager 從哪裡來」。

候選方案：

1. 直接手動改 SQLite
   - 優點：最快
   - 缺點：不可審計、容易操作錯、依賴 DB 工具、不適合長期流程
2. 環境變數 allowlist（例如啟動時把某些 LINE subject 自動升為 manager）
   - 優點：實作快
   - 缺點：部署與 local/dev 行為耦合、subject 變更難追蹤、容易造成隱性授權
3. 一次性 bootstrap CLI / admin command
   - 優點：顯式、可寫入文件、可保留最小權限、比直接 SQL 安全
   - 缺點：需要新增一個小型 operator command

推薦方案：**一次性 bootstrap CLI / admin command**

建議行為：

- 使用者先完成一次 LINE login，讓系統建立 `users` row
- operator 使用 repository 內的管理命令，依 `user_id` 或 `(auth_provider, provider_subject)` 將該帳號提升為 `manager`
- 之後由這位 manager 在產品 UI 內處理日常 link / promote 操作

建議指令形式（後續 implementation 再定稿）：

```powershell
Set-Location backend
go run ./cmd/admin promote-user --user-id <user-id> --role manager
```

可接受替代參數：

- `--provider line --subject <line-subject>`
- `--player-id <player-id>`（可選）

這個 bootstrap command 應該：

- 僅更新既有 `users` row，不自動創建假 user
- 找不到 user 時明確失敗
- 寫出可讀的 operator output
- 與 app runtime 分離，不暴露為公開 HTTP endpoint

### 3. Steady-state management model

一旦至少有一位 manager，日常管理應進入 manager-only UI + API，而不是繼續用 operator SQL。

推薦將需求拆成兩個 manager 能力：

1. **Account linking**
   - 查看已完成 LINE login 但尚未 linked 的使用者
   - 將該使用者連到既有 `players` record
   - 必要時解除 link 或改連到另一位 player
2. **Access management**
   - 查看已登入的使用者清單與其目前角色
   - 將特定使用者升為 `manager`
   - 視需求允許降回 `player`

### 4. Backend API shape

建議不要新增太多細碎 endpoint；先用一組 manager-only user-admin API 支撐後台頁面。

推薦最小 API：

- `GET /api/admin/users`
  - manager only
  - 支援 query filter：
    - `linkState=unlinked|linked|all`
    - `role=manager|player|all`
- `PATCH /api/admin/users/:userId`
  - manager only
  - 可更新欄位：
    - `role`
    - `playerId`

推薦 request/response 原則：

- 僅允許更新 `users.role` 與 `users.player_id`
- `playerId` 設為 `null` 代表 unlink
- 更新前需檢查 `playerId` 是否存在
- 若同一個 `player` 只能綁定單一使用者，需在這個 proposal 的 implementation slice 明確加上唯一性規則；若暫時不限制，也要在文件寫清楚

推薦在 implementation 時補一個明確決策：

- **建議同一位 player 同時只綁一個 active auth user**

理由：

- 避免一個球員資料被多人共用登入身份
- 比較符合真實世界的個人 LINE 帳號使用情境
- 能降低 manager UI 的歧義

若採用這個規則，則 backend 在 link 時應拒絕把某個 `player_id` 指派給第二個不同 `user`

### 5. Frontend UI shape

建議把這個能力放進 manager-only 的 admin 頁面，而不是塞進 pending-link 頁面。

原因：

- pending-link 頁面屬於被阻擋的 player 視角
- linking / promote 是 manager 的操作，不是 unlinked user 自助完成的流程

推薦最小 UI：

- 新增 manager-only route，例如：
  - `/admin/users`
- 頁面至少包含兩個區塊：
  - Unlinked accounts
  - User access / role management

每個 user row 顯示：

- display name
- auth provider
- provider subject（可視情況只顯示部分遮罩）
- current role
- current player link
- created / updated timestamp

主要操作：

- link to player
- unlink player
- promote to manager
- demote to player

### 6. Guardrails and business rules

建議在設計階段先鎖定以下 guardrails：

- manager 不能刪除 auth user，只能調整角色與 linkage
- manager 不應能把自己降權成最後一位非-manager，否則系統可能失去管理入口
- 至少要保證系統永遠存在一位 manager；若要降權最後一位 manager，backend 應拒絕
- newly authenticated LINE user 預設仍維持 `role=player` 與 `player_id=NULL`
- pending-link UX 仍保留，直到 manager 完成 linking

### 7. Rollout plan

推薦分三個 slice：

1. **Bootstrap slice**
   - 新增 operator CLI / admin command
   - 可將既有 LINE user 提升為 first manager
2. **Manager admin API slice**
   - 新增 manager-only user listing / update API
   - 補 backend tests 與 authorization checks
3. **Manager admin UI slice**
   - 新增 `/admin/users`
   - 完成 linking / role management UX

這樣可以先解決「誰是第一個 manager」的阻塞問題，再逐步補齊產品內操作。

## Risks / Open Questions

- UI 文案要用 `Admin`、`Manager`，還是兩者並存？技術上建議維持 `manager`
- 是否允許 manager 同時沒有 `player_id`？目前 architecture 文件允許
- 是否允許 manager 將自己降權？建議只在系統仍有其他 manager 時允許
- `player_id` 是否要強制單一 user 對應？本 proposal 建議是
- 是否需要 audit trail（誰把誰升權 / link）？MVP 可以先不做，但若要正式營運，後續應補

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - 本文件先做 repo 現況對齊檢查，確認與既有 auth architecture / local setup / smoke docs 一致
- 完成後如何驗收：
  - reviewer 可以從這份文件直接回答：
    - 第一個 manager 怎麼建立
    - manager 後續怎麼管理 unlinked user
    - 技術上是否需要新增 `admin` role
    - 後續 implementation 應拆成哪些 slice

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
- Agent review summary:

## Feedback

- Reviewer agent 1:
- Reviewer agent 2:
- Applied proposal updates:
