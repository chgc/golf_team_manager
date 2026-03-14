# Subagent Task Proposal

## Basic Information

- Phase: Phase 5
- Area: session management
- Proposed task name: sessions-feature
- Related todo id: `sessions-feature`
- Assigned subagent: session management feature agent

## Goal

交付第一個完整的 Manager-facing session feature：場次列表、建立 / 編輯表單、場次詳情、狀態流轉，以及與 backend session API / validation baseline 對齊的前後端流程。

本 work item 明確包含從 foundation baseline 擴充 session detail / update / status management 所需的 backend API 與 frontend page flow，讓 manager 能管理單一場次的完整生命週期。

## In Scope

- 實作場次列表，並明確以 `session date >= today` 視為 upcoming、`session date < today` 視為 history；`cancelled` 場次不論日期皆歸入 history
- 實作建立 / 編輯場次表單
- 實作場次詳情頁，顯示 session 自身欄位與已報名人數 / 剩餘名額 / 預估組數等統計，但不顯示完整報名者名單
- 實作 manual close / confirm / cancel / complete 等狀態操作，皆為 manager 手動觸發，並依明確狀態流轉矩陣進行 backend 驗證
- 補齊 session detail / update flow 所需的 backend `GET /api/sessions/{sessionId}`、`PATCH /api/sessions/{sessionId}`、repository、service、handler support
- 實作 registration deadline 到期後的 auto-close 檢查機制，且不引入背景排程元件
- 對齊 max players、deadline、status transition 等驗證規則

## Out of Scope

- 不進入 registration feature 的 player-side 報名 / 請假操作
- 不實作自動分組、tee time 排程、即時同步
- 不產出 reservation report 正式報表內容
- 不在此任務中加入新的 route-level authorization

## Dependencies

- `shared-domain-schema` must be completed first
- `backend-foundation` must be completed first
- `frontend-shell` must be completed first
- `auth-foundation` must be completed first
- `players-feature` should be used as the implementation template for backend handler extension and frontend signal/form patterns

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `frontend\src\app\features\sessions\`
  - 視需要更新 architecture / development docs
- 預計新增的資料夾 / 檔案：
  - session detail / form components
  - session repository / service / handler extensions
  - 對應測試與 feature handoff docs
- backend 預計補齊：
  - `GET /api/sessions/{sessionId}`（repository.GetByID / service.GetByID 已存在，補齊 handler 與 route 綁定）
  - `PATCH /api/sessions/{sessionId}`
  - `SessionService.Update()` 與 status-transition validation
  - `SessionRepository.Update()` 與對應 SQL update

## Technical Approach

- 使用的技術與模式：
  - Angular standalone feature pages + reactive forms
  - Gin handler / service / repository extension on top of backend foundation
  - shared session validation reuse，將 max players、deadline、status transition 規則集中在 backend service / domain validation
  - 以前後端共同維護的 session status contract 與明確狀態流轉矩陣：
    - `open -> closed`
    - `open -> cancelled`
    - `closed -> confirmed`
    - `closed -> cancelled`
    - `confirmed -> completed`
    - `confirmed -> cancelled`
    - `completed` / `cancelled` 視為終態，不允許再轉換
  - auto-close 採同步 service-level reconciliation：在 session list / detail read 與相關 status mutation entrypoints 先將 `registrationDeadline < now` 且仍為 `open` 的場次轉為 `closed`，不使用 scheduler / cron / background worker
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\frontend-shell-baseline.md`
  - `docs\architecture\auth-foundation.md`
  - `docs\architecture\players-feature.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- auto-close 檢查機制需與 manual close / confirm / cancel 流程保持一致，避免 deadline reconciliation 與手動狀態操作互相覆蓋
- 場次詳情中的即時計算欄位（報名人數、空位、預估組數）需先以目前資料模型可支撐的範圍為主
- 本 feature 不處理 registration UI，但 detail page 的資料結構需支撐後續報名名單整合
- upcoming / history 與 cancelled 場次的呈現規則已在本 proposal 固定，實作時需保持前後端一致
- session date 與 registrationDeadline 的先後順序驗證需先確認 shared-domain-schema 是否已涵蓋；若尚未涵蓋，則由本 feature service 層補齊

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - session feature flow smoke checks
- 完成後如何驗收：
  - manager 可建立 / 編輯 / 檢視 / 關閉 / 確認 / 取消 / 完成場次
  - deadline 與 status transition validation 一致
  - 場次詳情可顯示日期、球場、截止日、狀態、報名人數、剩餘名額、預估組數、備註
  - invalid status transition 回傳 validation error
  - session feature smoke checks 至少覆蓋 create -> edit -> detail -> close / confirm / cancel / complete，以及 expired deadline auto-close 行為

## Review Status

- Status: approved
- Reviewer: GPT-5.4, Claude Sonnet 4.6
- Review notes: Initial dual review found blocking proposal gaps. Proposal was revised, then approved by both reviewer agents on re-review.
- Agent review summary: Initial pass = blocking from both reviewer agents; re-review = GPT-5.4 approve, Claude Sonnet 4.6 approve.

## Feedback

- Reviewer agent 1:
  - 本 proposal 的整體方向與 Phase 5 目標大致一致，但目前仍不建議核可，主因是缺少可直接實作的 backend 交付清單。
  - 建議比照 `players-feature` 的寫法，明確把 `GET /api/sessions/{sessionId}`、`PATCH /api/sessions/{sessionId}` 與對應 repository / service / handler 納入 in-scope，避免實作者自行猜測 foundation 缺口。
  - `auto-close` 目前只有需求描述，尚未定義觸發方式與一致性策略；若不先寫清楚，實作時很容易出現把 side effect 放在 read path、或引入未規劃基礎設施的問題。
  - `upcoming / history` 的切分規則建議補成明確契約，避免前後端各自依 `date` 或 `registrationDeadline` 做不同判斷。
  - `status transition` 建議補成明確狀態機表述，至少寫清楚 `open -> closed -> confirmed` 與 `cancelled` 是否可由哪些狀態進入、是否允許回復。
  - 驗收條件建議從高層描述改成可測試行為，尤其要補上 detail 頁核心欄位、invalid transition、以及 auto-close 的 smoke check。
- Reviewer agent 2:
  - 缺少 `shared-domain-schema` 依賴宣告。本 feature 直接依賴 shared domain 的 Session struct、DTO、status enum 與驗證邏輯，應明確列為前置依賴。
  - 建議將 `players-feature` 列為參考模板依賴，讓 subagent 知道要對齊既有 backend handler extension 與 frontend signal / reactive form 模式。
  - 狀態流轉矩陣、auto-close 觸發時機、upcoming / history 切分條件目前都不夠明確，尚不足以直接進入實作。
  - 場次詳情頁應明確界定本次只顯示 session 自身欄位與統計，不延伸到完整報名名單。
  - `SessionService` / `SessionRepository` 的已知缺口應在 Planned Changes 中直接點名，避免實作者遺漏。
- Reviewer agent 2 (re-review):
  - `SessionRepository.GetByID()` 與 `sessionService.GetByID()` 已在 backend foundation 中實作；`GET /api/sessions/{sessionId}` 所需補齊的部分僅為 handler 方法與 route 綁定，實作前應先確認現有程式碼，避免重複實作整條堆疊。
  - `confirmed -> completed` 狀態轉換在技術方案的矩陣中已列出，但 In Scope 原本只明確提及 close / confirm / cancel；建議統一視為 manager 手動觸發操作，並在驗收 smoke check 中補上此轉換案例。
  - Session date 與 registrationDeadline 的先後順序約束需在實作前確認是否已由 shared-domain-schema 覆蓋，否則應在 service 層補齊。
- Reviewer agent 1 (re-review):
  - 本次修訂已補齊先前 review 指出的核心缺口：相依性、backend 補齊範圍、狀態流轉矩陣、upcoming/history 規則、detail page 顯示邊界與 auto-close 模型都已明確，proposal 已達可核可狀態。
  - 前後端交付內容已足夠具體：backend 明確點名 `GET /api/sessions/{sessionId}`、`PATCH /api/sessions/{sessionId}`、service / repository 補齊；frontend 也清楚涵蓋列表、建立 / 編輯、詳情與狀態操作流程。
  - 驗收條件已從高層描述收斂為可測試行為，能對應 create -> edit -> detail -> close / confirm / cancel / complete 與 expired deadline auto-close 的 smoke path。
  - 與目前 foundation 的一致性良好；實作時只需依 proposal 既定契約維持前後端對 `today`、deadline 與 status transition 的判斷一致即可。
- Applied proposal updates:
  - 補上 `shared-domain-schema` 與 `players-feature` 依賴說明。
  - 明確將 `GET /api/sessions/{sessionId}`、`PATCH /api/sessions/{sessionId}`、`SessionService.Update()`、`SessionRepository.Update()` 納入交付清單。
  - 固定 upcoming / history 規則、detail page 顯示範圍，以及可接受的 status transition 矩陣。
  - 將 auto-close 執行模型明確為同步 service-level reconciliation，且不引入排程元件。
  - 將驗收條件改寫成可測試的具體行為。
  - 補充 `complete` 為 in-scope manager action，並將對應 smoke check 納入驗收。
  - 補註 `GET /api/sessions/{sessionId}` 的既有 service / repository 前提，避免重工。
  - 補上 session date 與 registrationDeadline 先後順序驗證的確認要求。
