# Subagent Task Proposal

## Basic Information

- Phase: Phase 6
- Area: registration management
- Proposed task name: registrations-feature
- Related todo id: `registrations-feature`
- Assigned subagent: registration feature agent

## Goal

交付第一個完整的 registration feature：球員報名 / 取消 / 請假流程、manager 手動調整名單、與 session / player 狀態規則對齊的前後端流程，讓 v1 的核心報名作業可完整運作。

## In Scope

- 實作 player-facing 場次報名 / 取消 / 請假操作：
  - 報名使用既有 `POST /api/sessions/{sessionId}/registrations`
  - 取消 / 請假皆在 v1 明確對應為將既有 registration 狀態改為 `cancelled`
  - player 只能操作自己的 registration，不新增 route-level authorization，先沿用目前 auth baseline 的 UI / identity context 表達
- 實作 manager-facing 報名名單調整流程：
  - manager 可代球員建立 registration
  - manager 可將既有 registration 在 `confirmed` / `cancelled` 間切換，以表達取消、請假、恢復報名等調整
  - manager 不可越過既有 business constraints：inactive player 不可報名、session 非 `open` 不可新增 confirmed registration、session capacity 不可超量、duplicate registration 不可建立
- 對齊 inactive player、session closed / confirmed / cancelled、session capacity、duplicate registration 等限制規則
- 補齊 registration mutation flow 所需的 backend `PATCH /api/registrations/{registrationId}`、repository、service、handler support
- 整合既有 `GET /api/sessions/{sessionId}/registrations`、`POST /api/sessions/{sessionId}/registrations`
- 補齊前端 registration feature pages / 區塊、表單與既有 `session-list-page` 的選取 detail 區塊整合
- 在本任務內移除獨立 `/registrations` route，避免與 session detail 內嵌 registration flow 重複

## Out of Scope

- 不實作自動分組、tee time 排程、WebSocket 即時同步
- 不產出 reservation report 正式報表
- 不在此任務中導入正式 LINE OAuth

## Dependencies

- `shared-domain-schema` must be completed first
- `backend-foundation` must be completed first
- `frontend-shell` must be completed first
- `auth-foundation` must be completed first
- `players-feature` must be completed first
- `sessions-feature` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `frontend\src\app\features\sessions\`
  - `frontend\src\app\features\registrations\`
  - `frontend\src\app\app.routes.ts`
  - 視需要更新 architecture / development docs
- 預計新增的資料夾 / 檔案：
  - registration components / helpers integrated into existing session detail area
  - registration repository / service / handler extensions
  - 對應測試與 feature handoff docs
- backend 預計補齊：
  - `PATCH /api/registrations/{registrationId}`
  - `RegistrationService.UpdateStatus()`（或等效命名）以處理 cancel / leave / restore flow
  - `RegistrationRepository.UpdateStatus()` 與對應 SQL update
  - 既有 `POST /api/sessions/{sessionId}/registrations` flow 的 business validation 對齊確認：duplicate registration、inactive player、session not open、capacity full
- frontend 預計補齊：
  - 既有 `frontend\src\app\features\sessions\pages\session-list-page\` 中的 detail 區塊內嵌 registration roster（manager view）
  - player-facing registration action UI，放在既有 sessions route / selected-session detail flow 內，而非新增新的 detail route 或獨立 `/registrations` route
  - registration data-access update method 與對應 signal / form state
  - 移除或退場既有 `frontend\src\app\features\registrations\pages\registration-list-page\` 與 `/registrations` route，避免 dead code / route conflict

## Technical Approach

- 使用的技術與模式：
  - Angular standalone feature pages + reactive forms / signals
  - Gin handler / service / repository extension on top of backend foundation
  - shared registration validation reuse，集中處理 duplicate registration、session capacity、player inactive、session not open 等規則
  - frontend integration target 固定為既有 `SessionListPage` 的 selected-session detail 區塊；本任務不新增 `session-detail-page`，而是在現有 sessions feature route 中擴充 registration UI
  - registration status contract（v1 不新增新的 domain status）：
    - 新增報名：建立 `confirmed` registration
    - player 取消：`confirmed -> cancelled`
    - player 請假：在 v1 明確等同於 `confirmed -> cancelled`，僅為 UI / wording 差異，不新增新狀態
    - manager 恢復報名：`cancelled -> confirmed`
    - manager 取消 / 請假調整：`confirmed -> cancelled`
  - registration mutation matrix：
    - `POST /api/sessions/{sessionId}/registrations`：建立 `confirmed`
    - `PATCH /api/registrations/{registrationId}`：在 `confirmed` / `cancelled` 間切換
    - player 僅能操作自己的 registration；manager 可操作 roster 中任一 registration
  - manager actions and limits：
    - manager 可代球員建立報名、取消報名、恢復已取消報名
    - manager 仍不可忽略 inactive player、capacity full、session not open 等既有 backend validation
    - auth 仍沿用目前 auth baseline 的 UI / identity context，不在本任務新增 route-level authorization
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\players-feature.md`
  - `docs\architecture\sessions-feature.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- player-facing 與 manager-facing flow 會共享 session detail 頁面，但需在 UI 上清楚區分可用操作，避免單一頁面承載過多誤導性操作
- 既有 `registration-list-page` 與 `/registrations` route 將在本任務退場，需確保移除後不留下失效導覽或 dead code
- session capacity、inactive player、duplicate registration 等規則需維持前後端一致
- v1 明確不新增 registration 第三種狀態；若後續要區分請假與取消，需回到 shared-domain-schema 擴充

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - registration feature flow smoke checks
- 完成後如何驗收：
  - player 可對 open session 建立報名，成功後 registration status 為 `confirmed`
  - player 可取消 / 請假自己的報名，成功後 registration status 為 `cancelled`
  - manager 可在 roster 中代球員建立報名、取消報名、恢復已取消報名
  - duplicate registration 會回傳 `conflict`
  - inactive player 會回傳 `player_inactive`
  - session 非 `open` 時新增或恢復 confirmed registration 會回傳 `session_not_open`
  - session capacity 已滿時新增或恢復 confirmed registration 會回傳 `session_capacity_full`
  - business validation 與 API error shape 維持一致
  - smoke path 至少覆蓋：player register -> player cancel/leave -> manager restore -> duplicate -> inactive player -> non-open session -> capacity full

## Review Status

- Status: approved
- Reviewer: Claude Sonnet 4.6, GPT-5.4
- Review notes: Initial dual review found blocking proposal gaps. Proposal was revised three times, then approved by both reviewer agents on final re-review.
- Agent review summary: Initial pass = blocking from both reviewer agents; second pass = mixed approve/blocking; final re-review = GPT-5.4 approve, Claude Sonnet 4.6 approve.

## Feedback

- Reviewer agent 1:
  - 本 proposal 的整體方向與 Phase 6 目標一致，相依性宣告完整，out-of-scope 邊界清晰，但尚未達到可核可狀態，主因是 registration 狀態契約、backend 具名交付清單、以及 player-facing / manager-facing 操作邊界三者均未明確定義。
  - Risks 區塊第三點（「狀態契約需在 proposal 細化後再進入實作」）本身即為 blocking 指標，proposal 在解決此問題前不應進入 approved 狀態。
  - 建議參照 sessions-feature 在第一次 review 後的修訂模式：將 Registration 的初始狀態、可允許的狀態轉換矩陣（含「請假」的資料語意）、以及 manager override 的操作契約補入 Technical Approach，固定後再重送 review。
  - 「請假」應在 Technical Approach 或 In Scope 中明確說明其在 Registration 資料模型中的對應，例如：是否視為 `cancelled`？是否有額外 flag？若需新增狀態，應說明對 shared-domain-baseline 的影響。
  - Backend 交付清單應比照 sessions-feature 的寫法，在 Planned Changes 中逐一具名：至少應包含 `PATCH /api/registrations/{registrationId}`、對應的 `RegistrationService` 方法、`RegistrationRepository` 方法，以及現有 `POST /api/sessions/{sessionId}/registrations` handler 是否需要補齊 business validation。
  - Player-facing 操作（報名 / 取消 / 請假）與 manager-facing 操作（名單調整 / override）應在 In Scope 中分別條列，並說明各操作使用哪個 endpoint、前端是在哪個 page / component 觸發。
  - 驗收條件應改寫為具體可測試行為，例如：player 完成報名後 registration status 為何、取消後 status 如何轉換、manager override 的回應與 error path 如何驗收。
  - 前端 pages 應在 Planned Changes 中具名，避免 subagent 自行決定 routing 架構。
- Reviewer agent 2:
  - 目前 proposal 方向正確，但尚未達到可核可狀態；主要缺口是 registration 狀態契約仍未定義清楚，實作時容易各自解讀。
  - 請明確說明 `報名`、`請假`、`取消`、`manager override` 分別對應到哪些 domain status 與 API 操作，並補上狀態流轉矩陣。
  - 請不要再以「視需要補齊 update / cancel flow」描述 backend 範圍，需直接寫清楚是否實作 registration update endpoint，以及對應的 repository / service / handler / frontend API 交付。
  - `manager-facing 報名名單調整流程` 目前過於籠統；需明確列出 manager 可執行的操作、不可越過的限制規則，以及是否仍沿用目前 auth baseline 的 UI / identity context，而不新增 route-level authorization。
  - 驗收條件需從高層描述收斂成可測試行為，至少列出 register、cancel、duplicate、inactive player、session not open、capacity full、manager adjustment 等 smoke / acceptance cases。
- Reviewer agent 2 (re-review):
  - 本次修訂已補齊前一輪 review 的核心缺口：scope、相依性、backend / frontend 交付清單、registration status contract、以及 manager / player 操作邊界都已明確，proposal 已達可核可狀態。
  - v1 的 registration 契約已固定為僅使用 `confirmed` / `cancelled`，並明確定義 player 的報名 / 取消 / 請假與 manager 的建立 / 取消 / 恢復名單調整，足以作為前後端一致實作依據。
  - 驗收條件已收斂為可測試的 success / error scenario matrix，涵蓋 register、cancel / leave、manager restore、duplicate、inactive player、session not open、capacity full，符合 repo review gate 對可驗收性的要求。
  - proposal 也已清楚維持目前 auth baseline：僅以既有 UI / identity context 表達角色邊界，不在此任務額外引入 route-level authorization，與既有 workflow 一致。
- Reviewer agent 2 (final re-review):
  - 本次修訂已明確將前端整合載體固定為既有 `SessionListPage` 的 selected-session detail 區塊，成功消除先前「不存在 session detail page」的執行落點風險。
  - 提案也已明確宣告退場既有獨立 `/registrations` route 與 `registration-list-page`，可避免後續實作出現 route conflict、重複入口或 dead code。
  - 目前 proposal 的 scope、狀態契約、前後端交付清單、manager / player 邊界與驗收矩陣已可支撐進入 approved 流程，符合 repo review gate 要求。
- Reviewer agent 1 (re-review):
  - 本次 review 確認先前兩個 blocking 問題（狀態契約、backend 交付清單、player / manager 操作邊界、驗收條件）均已妥善修正，提案品質明顯提升。
  - 發現新的 blocking 問題：sessions feature 目前只存在 `session-list-page`，並不存在 session detail page；提案在前端整合方向上缺乏可執行的操作對象，需補充說明整合載體。
  - 發現第二個 blocking 問題：`frontend\\src\\app\\features\\registrations\\pages\\registration-list-page\\` 已存在，但提案聲明不新增獨立 `/registrations` route，卻未說明此現有頁面的處置方式，會造成 route conflict 或 dead code 風險。
  - Backend 補齊說明、驗收條件、狀態轉換矩陣與 manager 可操作範圍均已達標，可維持現狀。
- Reviewer agent 1 (final re-review):
  - 本輪確認先前 reviewer agent 1 提出的兩個 blocking 問題均已在提案文字中妥善處理：整合載體固定為現有 `SessionListPage` selected-session detail 區塊，且已明確退場 `registration-list-page` 元件與 `/registrations` route。
  - 已透過實際查閱 `app.routes.ts` 與 frontend 目錄結構確認：`session-list-page`、`registration-list-page` 及 `/registrations` route 均確實存在於程式庫，提案描述與現況一致。
  - `registrations\data-access\registrations-api.ts` 雖未逐字點名，但 proposal 已明確保留並更新 registration data-access 層；實作時不應誤刪此目錄。
  - 驗收條件已涵蓋所有關鍵 success / error path，符合 repo review gate 對可測試性的要求。
  - v1 僅使用 `confirmed` / `cancelled` 兩種狀態，且請假明確等同 `cancelled` 的契約定義清晰，足以作為前後端一致實作依據。
- Applied proposal updates:
  - 固定 v1 registration status contract：僅使用 `confirmed` / `cancelled`，其中請假明確等同於 `cancelled`。
  - 明確列出 `PATCH /api/registrations/{registrationId}`、`RegistrationService.UpdateStatus()`、`RegistrationRepository.UpdateStatus()` 與前端 data-access 交付。
  - 將 player-facing 與 manager-facing 操作、可做事項與限制規則拆開描述。
  - 固定 registration UI 以既有 `SessionListPage` 的 selected-session detail 區塊為主，不新增新的 session detail route。
  - 明確將 `frontend\src\app\features\sessions\` 與 `frontend\src\app\app.routes.ts` 納入修改面，並指定移除既有獨立 `/registrations` route / `registration-list-page`。
  - 將驗收條件改寫成具體的 success / error scenario matrix。
  - 保留並更新既有 `registrations\data-access\registrations-api.ts`，不移除 data-access 層。
