# Subagent Task Proposal

## Basic Information

- Phase: Phase 7
- Area: reservation reporting
- Proposed task name: reservation-report
- Related todo id: `reservation-report`
- Assigned subagent: reservation summary feature agent

## Goal

交付第一個可用的球場預約摘要報表功能，讓 manager 能依單一場次輸出 reservation summary，支援球場預約所需的核心名單與統計資訊。

## In Scope

- 補齊 `GET /api/reports/sessions/{sessionId}/reservation-summary` 正式回應內容
- 依單一 session 產出 `ReservationSummaryReadDTO`
- 彙整球員名單、已確認報名人數、預估組數與關鍵場次資訊
- 在既有 `SessionListPage` 的 selected-session detail 區塊中提供 manager 可檢視的 reservation summary 區塊
- 提供 manager 可直接複製的純文字 reservation summary
- 對齊 session / registration 狀態規則，避免非適用場次輸出錯誤摘要

## Out of Scope

- 不產出 PDF / Excel 匯出
- 不實作自動分組與 tee time 排程
- 不引入外部報表服務

## Dependencies

- `shared-domain-schema` must be completed first
- `backend-foundation` must be completed first
- `auth-foundation` must be completed first
- `sessions-feature` must be completed first
- `registrations-feature` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\domain\`
  - `backend\internal\http\`
  - `backend\internal\service\`
  - `frontend\src\app\shared\models\`
  - `frontend\src\app\features\reports\data-access\`
  - `frontend\src\app\features\sessions\pages\session-list-page\`
  - 視需要更新 architecture / development docs
- 預計新增的資料夾 / 檔案：
  - `backend\internal\domain\report.go`
  - `frontend\src\app\features\reports\data-access\reports-api.ts`
  - 對應測試與 feature handoff docs
- 後端具名交付：
  - `APIHandlers.GetReservationSummary()`
  - `ReportService.GetReservationSummary(ctx, sessionID string) (domain.ReservationSummaryReadDTO, error)`
  - `domain.ReservationSummaryReadDTO`
  - `domain.ReservationSummaryPlayerDTO`
  - 將 `router.go` 既有 `NotImplemented` report route 改接 `APIHandlers.GetReservationSummary()`
  - 在共用 HTTP error mapping 中補齊 report-specific error 對應
- 前端具名交付：
  - `ReportsApi.getReservationSummary(sessionId: string)`
  - `ReservationSummaryReadDto`
  - `SessionListPage` selected-session detail summary card、inline state 與 copy action

## Technical Approach

- 使用的技術與模式：
  - Gin handler + `ReportService`，以既有 `SessionRepository.GetByID()`、`RegistrationRepository.ListBySessionID()`、`PlayerRepository.List()` 組合 summary，不在 v1 新增獨立 report repository
  - frontend summary UI 固定整合到既有 `SessionListPage` selected-session detail 區塊，不新增獨立 report route
  - summary output 以目前 v1 可得資料為準，不提前混入 v2 grouping / tee time 資訊
  - `estimatedGroups` 計算規則固定為 `ceil(confirmedPlayerCount / 4)`
  - API success contract 固定為：
    - `sessionId`
    - `sessionDate`
    - `courseName`
    - `courseAddress`（可為空字串，不使用 `null`）
    - `registrationDeadline`（必填，沿用既有 session contract）
    - `sessionStatus`
    - `confirmedPlayerCount`
    - `estimatedGroups`
    - `summaryText`
    - `confirmedPlayers[]`
  - `confirmedPlayers[]` 項目固定為：
    - `playerId`
    - `playerName`
    - 陣列依 `playerName` 升冪排序，與 `summaryText` roster 順序一致
  - 適用場次狀態矩陣固定為：
    - `confirmed` -> 回傳 `200` summary
    - `completed` -> 回傳 `200` summary
    - `open` / `closed` / `cancelled` -> 回傳 `422 session_not_eligible_for_report`
  - error contract 固定為：
    - session 不存在 -> `404 session_not_found`
    - session 狀態不適用 -> `422 session_not_eligible_for_report`
    - 無任何 `confirmed` registration -> `422 reservation_summary_empty`
  - auth / role contract 固定為：
    - endpoint 沿用既有 API auth middleware
    - 僅 `manager` 可取得 reservation summary；非 manager 預期回傳 `403`
  - `summaryText` 純文字模板固定至少包含以下行序：
    - `Session: {sessionDate}`
    - `Course: {courseName}`
    - `Address: {courseAddress or N/A}`
    - `Deadline: {registrationDeadline}`
    - `Status: {sessionStatus}`
    - `Confirmed Players: {confirmedPlayerCount}`
    - `Estimated Groups: {estimatedGroups}`
    - `Roster:`
    - `- {playerName}`（每位 confirmed player 一行，依 player name 升冪排序）
  - frontend rendering matrix 固定為：
    - selected session 為 `confirmed` / `completed` 時呼叫 summary API，成功後顯示 summary card 與 copy button
    - selected session 為 `open` / `closed` / `cancelled` 時不呼叫 summary API，顯示 inline hint：此場次尚不可產生預約摘要
    - backend 回傳 `422 reservation_summary_empty` 時顯示 inline empty state：目前尚無已確認報名，無法產生預約摘要
    - report loading / empty / ineligible state 僅顯示在 summary card 區塊內，不可覆寫既有 page-level generic error banner
  - copy action contract 固定為：
    - 直接複製 backend 回傳的 `summaryText`
    - 使用瀏覽器 Clipboard API
    - 成功 / 失敗狀態以 summary card 區塊內的 inline message 呈現
- 依循的規範文件：
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\sessions-feature.md`
  - `docs\architecture\registrations-feature.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- `summaryText` 文案格式需保持簡潔且穩定，避免 manager 複製後仍需手動大幅編修
- v1 僅輸出已確認名單與預估組數，不提前混入分組、tee time 或外部球場欄位
- 若後續球場預約需要更多欄位，應在後續 phase 以 contract 擴充處理，不在本提案中預留未定義欄位

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-build`
  - `just frontend-test`
  - reservation summary runtime smoke checks
- 完成後如何驗收：
  - Case 1: `confirmed` session 回傳 `200`，response 內 `confirmedPlayerCount`、`estimatedGroups`、`confirmedPlayers[]` 與現有 registration roster 一致
  - Case 2: `completed` session 回傳 `200`，summary 可供 manager 檢視與複製
  - Case 3: `open` / `closed` / `cancelled` session 回傳 `422 session_not_eligible_for_report`
  - Case 4: 不存在的 `sessionId` 回傳 `404 session_not_found`
  - Case 5: 適用狀態但無任何 `confirmed` registration 時回傳 `422 reservation_summary_empty`
  - Case 6: frontend 在既有 `SessionListPage` selected-session detail 區塊正確顯示 summary card 與 copy action，且內容與 backend response 一致
  - Case 7: frontend 在 `open` / `closed` / `cancelled` session 僅顯示 inline ineligible hint，不觸發 page-level generic error banner
  - Case 8: copied text 與 `summaryText` 模板完全一致，包含固定欄位順序與 roster 行格式

## Review Status

- Status: approved
- Reviewer: reviewer agents (GPT-5.4 / Claude Sonnet 4.6 across 3 rounds)
- Review notes: 首輪雙 reviewer 皆為 blocking；第二輪一位 approve、一位 blocking；第三輪雙 reviewer 皆為 approve。proposal 已補齊 DTO contract、frontend rendering matrix、`summaryText` 模板、auth / role contract 與 error mapping，達到可實作狀態。
- Agent review summary: round 1 blocking + blocking -> round 2 approve + blocking -> round 3 approve + approve

## Feedback

- Reviewer agent 1:
  - 提案方向正確，但原版本未明確定義 reservation summary contract、session 適用狀態、前端整合落點與可測試驗收矩陣，因此尚不足以直接核可。
  - 已依建議補齊 `ReservationSummaryReadDTO` 欄位、`estimatedGroups` 規則、`SessionListPage` selected-session detail 整合位置，以及 success / error scenario matrix。
- Reviewer agent 2:
  - 原版本缺少具名後端交付與明確 error contract，且 `GET /api/reports/sessions/{sessionId}/reservation-summary` 的 response body 欄位與 frontend 載點不夠具體，會導致 subagent 各自解讀。
  - 已依建議明確指定 `APIHandlers.GetReservationSummary()`、`ReportService.GetReservationSummary()`、`domain.ReservationSummaryReadDTO` / `ReservationSummaryPlayerDTO`，並固定 manager summary UI 嵌入現有 `SessionListPage`。
- Reviewer agent 3:
  - 第二輪 review 認為 proposal 已可核可，但建議額外明文化 `summaryText` 模板、欄位空值契約、manager-only auth 與 copy action 行為，以降低實作歧義。
- Reviewer agent 4:
  - 第二輪 review 指出剩餘 blocking 在於 frontend rendering matrix 尚未固定，以及 `summaryText` 缺少可驗證的固定格式；若不補齊，容易讓 summary card 狀態處理污染既有 page-level error UX。
- Reviewer agent 5:
  - 第三輪 review 確認 proposal 已補齊前兩輪卡點：DTO 欄位、前端整合落點、error matrix、frontend rendering states 與 `summaryText` contract 都已明文化，具備進入 approved 階段的條件。
- Reviewer agent 6:
  - 第三輪 review 確認 `summaryText` 的空值 fallback、frontend inline state 邊界與 copy action contract 已足夠清楚；僅建議補充 `confirmedPlayers[]` 排序與非 manager 的預期 HTTP 狀態碼。
- Applied proposal updates:
  - 新增 `auth-foundation` 相依性，並引用 `docs\architecture\registrations-feature.md`
  - 補齊 reservation summary DTO 欄位、player roster 欄位與 `summaryText`
  - 固定 `estimatedGroups = ceil(confirmedPlayerCount / 4)`
  - 固定可輸出狀態為 `confirmed` / `completed`，並定義 `404 session_not_found`、`422 session_not_eligible_for_report`、`422 reservation_summary_empty`
  - 明確指定 backend / frontend 具名交付與 `SessionListPage` selected-session detail 整合落點
  - 將驗收條件改寫為具體 success / error scenario matrix
  - 補齊 frontend rendering matrix，明確規定何時呼叫 / 跳過 summary API，以及 ineligible / empty 狀態需使用 inline state，不可污染 page-level generic error UX
  - 固定 `summaryText` 欄位順序、roster 行格式、copy action 行為，以及共用 HTTP error mapping 交付
  - 補充 `confirmedPlayers[]` 依 `playerName` 升冪排序，並明確記錄非 manager 預期回傳 `403`
