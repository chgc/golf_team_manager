# Subagent Task Proposal

## Basic Information

- Phase: Phase 8
- Area: quality assurance and seed data
- Proposed task name: qa-and-seed-data
- Related todo id: `qa-and-seed-data`
- Assigned subagent: qa and seed data agent

## Goal

補齊 v1 MVP demo / 試用前的品質保證與種子資料基線，讓團隊可快速建立可展示環境，並用一致的 smoke path 驗證 players / sessions / registrations / reservation report 主流程。

## In Scope

- 補齊 v1 關鍵流程缺少的 backend / frontend 測試，固定聚焦於 players / sessions / registrations / reservation report
- 提供可重複建立的 local/dev seed data（球員、場次、報名）
- 補齊 demo / smoke check 所需的最小資料初始化流程
- 盤點並補強 `HomePage`、`PlayerListPage`、`SessionListPage` 的 loading、empty、error state
- 更新本地開發與 demo 操作文件

## Out of Scope

- 不實作新的 v2 grouping / tee time / websocket 功能
- 不引入外部測試平台或 e2e framework
- 不處理正式 production deployment / CI pipeline

## Dependencies

- `shared-domain-schema` must be completed first
- `backend-foundation` must be completed first
- `frontend-shell` must be completed first
- `auth-foundation` must be completed first
- `players-feature` must be completed first
- `sessions-feature` must be completed first
- `registrations-feature` must be completed first
- `reservation-report` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `backend\cmd\`
  - `justfile`
  - `frontend\src\app\`
  - `docs\development\`
  - `README.md`
  - 視需要更新 architecture / handoff docs
- 預計新增的資料夾 / 檔案：
  - `backend\cmd\seed\main.go`
  - 補強後的測試檔案
  - `docs\development\demo-smoke-check.md`
- 後端具名交付：
  - `backend\cmd\seed\main.go`
  - root `justfile` 新增 `backend-seed`
  - seed command 為 local/dev only，執行時會清空並重建目前 `DB_PATH` 指向的 SQLite 資料
  - seed command 執行前需驗證目前環境為 local/dev；若非 dev/stub 模式則直接 exit with error
  - seed command 需可重複執行；同一資料庫重跑後，資料筆數與固定 demo dataset 必須一致
  - 補強 `backend\internal\http\router_test.go` 的 feature smoke / error path coverage
  - 視需要補強 `backend\internal\service\` 測試，固定聚焦 report eligibility、seed compatibility 與主要邊界規則
- 前端具名交付：
  - 補強 `frontend\src\app\features\home\home-page\`
  - 補強 `frontend\src\app\features\players\pages\player-list-page\`
  - 補強 `frontend\src\app\features\sessions\pages\session-list-page\`
  - 補強 `frontend\src\app\features\reports\data-access\`
  - 每個指定頁面至少明確覆蓋 loading / empty / error 三類 state 中與現有功能直接相關的部分
  - `SessionListPage` 額外補 manager / player smoke-related state 驗證
  - `docs\development\local-setup.md`
  - `docs\development\demo-smoke-check.md`

## Technical Approach

- 使用的技術與模式：
  - 優先沿用既有 Go test、Angular unit test 與 root `justfile` 指令，不新增新的測試框架
  - seed data 以 repo 內可重複執行的本地 bootstrap 方式提供，避免手動 SQL
  - 測試與 seed data 內容只覆蓋 v1 已完成功能，不提前混入 v2+ backlog
  - loading / empty / error state 僅補在既有 feature page 與 data-access flow，不另開新頁面
  - seed dataset 固定為 deterministic demo 資料集：
    - 1 位 manager 身分驗證路徑：沿用現有 `dev_stub` 預設 manager principal，不建立獨立 user table seed
    - 1 位 player 身分驗證路徑：固定以 API/debug-header smoke 驗證，不要求瀏覽器 UI 內建角色切換器
    - player smoke contract 固定為：使用 `Invoke-RestMethod` / `curl` 搭配 `X-Debug-Role=player`、`X-Debug-Player-ID=<seed-player-id>`、`X-Debug-Display-Name=<seed-player-name>` 驗證 player 主要 API / state，不在本任務內新增 frontend interceptor 或 `AuthShell` identity toggle UI
    - 6 位 players：5 位 `active`、1 位 `inactive`
    - 4 個 sessions：1 `open`、1 `confirmed`、1 `completed`、1 `cancelled`
    - 7 筆 registrations：固定為 5 筆 `confirmed`、2 筆 `cancelled`，且至少 3 筆 `confirmed` 屬於同一 `confirmed` session，以確保 reservation summary 可穩定產出
  - frontend state matrix 固定為：
    - `HomePage`: shell section 正常顯示與基本 empty-safe rendering（因頁面不直接發 API call，本任務不補 error state）
    - `PlayerListPage`: 空球員清單、API error、載入中狀態
    - `SessionListPage`: 空場次清單、session API error、registration roster empty、reservation summary ineligible / empty / error state
- 依循的規範文件：
  - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\players-feature.md`
  - `docs\architecture\sessions-feature.md`
  - `docs\architecture\registrations-feature.md`
  - `docs\architecture\reservation-report.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- seed data 內容需在「可 demo」與「不污染正式資料」之間取得平衡，需明確限定為 local/dev 用途
- 若現有頁面已有零散 error handling，需先盤點一致的 UX 邊界，避免補強後造成不同頁面行為不一致
- seed command 需與 migration / app startup 完全分離，不可在 `backend-start` 或 app startup 自動執行

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`（於 `backend\` 執行）
  - `just frontend-build`
  - `just frontend-test`
  - `just backend-seed`
  - seed data bootstrap smoke checks
  - end-to-end manual smoke path covering players -> sessions -> registrations -> reservation report
- 完成後如何驗收：
  - Case 1: 全 repo 既有 build / test 指令保持通過
  - Case 2: `just backend-seed` 成功後，資料庫內固定存在 6 位 players、4 個 sessions、7 筆 registrations，且重跑後資料筆數不累加
  - Case 3: manager smoke path 明確可重現：啟動 seed data -> 檢視 players / sessions -> 開啟 `confirmed` session -> 檢查 registration roster -> 取得 reservation summary -> 複製摘要文字
  - Case 4: player smoke path 明確可重現：以 `Invoke-RestMethod` / `curl` 搭配固定 debug headers 呼叫 player principal / session / registration 相關 API，確認 seeded player 的 registration state 與可報名 `open` session 資料正確
  - Case 5: `HomePage`、`PlayerListPage`、`SessionListPage` 的指定 loading / empty / error state 有對應 unit test 或明確 smoke 驗證
  - Case 6: `docs\development\local-setup.md` 與 `docs\development\demo-smoke-check.md` 可讓新進成員完成最小 demo 環境並重現 manager / player smoke path

## Review Status

- Status: approved
- Reviewer: reviewer agents (GPT-5.4 / Claude Sonnet 4.6 across 3 rounds)
- Review notes: 首輪雙 reviewer 皆為 blocking；第二輪一位 approve、一位 blocking；第三輪雙 reviewer 皆為 approve。proposal 已補齊 seed 入口、dev-only guard、deterministic dataset、page state matrix、manager/player smoke contract 與具體驗收矩陣，達到可實作狀態。
- Agent review summary: round 1 blocking + blocking -> round 2 approve + blocking -> round 3 approve + approve

## Feedback

- Reviewer agent 1:
  - proposal 方向正確，但首版仍缺少固定契約：哪些測試要補、哪些頁面 state 要補、seed data 會建立哪些資料、以及 manager / player smoke path 如何被重現，尚不足以直接進入 approved。
  - 已依建議同步最新 `reservation-report` architecture 文件，並補齊 seed command、dataset 規格、auth / demo contract 與具體驗收矩陣。
- Reviewer agent 2:
  - 提案整體方向與相依性合理，但首版最大 blocking 點在於 seed data 入口名稱 / 位置未定案，且測試與頁面補強目標不夠具名，容易讓 subagent 各自解讀。
  - 已依建議固定 `backend\cmd\seed\main.go` + `just backend-seed`、明確列出 6 players / 4 sessions / 7 registrations 的 deterministic dataset，並指定 `HomePage`、`PlayerListPage`、`SessionListPage` 為本輪 page state 補強範圍。
- Reviewer agent 3:
  - 第二輪 review 認為 proposal 已接近可核可，但建議補上 registration 狀態分配、`local-setup.md` 具名交付、dev-only seed 保護要求，以及 `HomePage` 為何不補 error state 的說明。
- Reviewer agent 4:
  - 第二輪 review 指出唯一 remaining blocking 在於 player smoke path 的 auth / demo contract 尚未落到可實作層級；若不補齊，subagent 仍可能自行發明 frontend 身分切換機制。
- Reviewer agent 5:
  - 第三輪 review 確認 proposal 已補齊前幾輪 blocking 的核心缺口：seed 入口、dev-only guard、deterministic dataset、page state matrix、acceptance criteria，以及 manager/player smoke contract 都已明文化，可進入 approved 流程。
- Reviewer agent 6:
  - 第三輪 review 確認 proposal 與目前 repo workflow、Phase 8 目標及既有 architecture 文件一致，且不再要求額外發明 frontend 角色切換機制；僅建議後續在文件中將 dev-only guard 條件與 seed placeholder 實值寫得更具體。
- Applied proposal updates:
  - 固定 seed data 入口為 `backend\cmd\seed\main.go`，並新增 root `just backend-seed` 為對應命令
  - 固定 seed command 為 local/dev only、可重複執行、會清空並重建目前 DB_PATH 指向的 SQLite 資料
  - 明確定義 deterministic demo dataset：6 players、4 sessions、7 registrations，並沿用 `dev_stub` / debug headers 驗證 manager / player 路徑
  - 明確指定要補強的 backend / frontend 檔案與 page matrix，避免 scope 過寬
  - 將驗收條件改寫為具體 manager / player smoke path 與資料筆數檢查
  - 同步更新 Technical Approach 的規範引用，納入已完成的 `docs\architecture\reservation-report.md`
  - 補充 7 筆 registrations 的固定分配、`docs\development\local-setup.md` 具名交付、seed command 的 dev-only guard，以及 `HomePage` 不補 error state 的理由
  - 將 player smoke path 明確收斂為 API/debug-header smoke contract，不在本任務內新增 frontend 角色切換器或 `AuthShell` override UI
