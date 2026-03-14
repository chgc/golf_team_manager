# Golf Team Manager 執行計畫

## 問題描述

依 `CKNotepads/01 Projects/Golf team manager` 下的功能分析與規格，為目前幾乎空白的 `golf_team_manager` repo 制定一份可直接執行的實作計畫，並以 **v1 MVP 先落地、v2+ 逐步擴充** 為主線。

## 目前狀態分析

- 本地 repo 目前僅有 `.git`，尚未建立前端、後端、資料庫、文件或測試基礎設施。
- 規格文件已明確定義核心領域模型：`Player`、`Session`、`Registration`、`Group(v2)`。
- Tech Stack 指向 **Angular v20 + Angular Material + Go + SQLite**，並以 **LINE OAuth** 為正式認證方案。
- MVP 功能聚焦在：
  - 球員管理
  - 場次管理
  - 報名 / 請假流程
  - 球場預約摘要輸出
- v2 起才進入自動分組、Tee Time、WebSocket 即時更新；v3/v4 為後續 backlog。

## 規格來源

- `Golf Team Manager - 00 設計總覽.md`
- `Golf Team Manager - 01 功能規劃.md`
- `Golf Team Manager - 02 資料模型.md`
- `Golf Team Manager - 03 報名流程.md`
- `Golf Team Manager - 04 球員管理.md`
- `Golf Team Manager - 05 場次管理.md`
- `Golf Team Manager - 06 分組與排程.md`
- `Golf Team Manager - 07 球場預約報表.md`
- `Golf Team Manager - 08 Tech Stack.md`

## 規劃假設

- **主要交付主線**：以 `v1 MVP` 為第一個可交付版本。
- **技術主線**：以正式架構 `Angular + Go + SQLite` 規劃，不以純前端 localStorage 版作為主要目標。
- **認證策略**：保留 `LINE OAuth` 為正式需求，但在開發期允許以 dev stub / mock identity 降低整合阻塞。
- **範圍控管**：`v2/v3/v4` 只做架構預留與 backlog 拆解，不納入第一波交付。

## docs 歸檔策略

依「**實作階段 / 版本**」分類，文件歸檔在：

- `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
- `docs\plans\v1-mvp\phase-0-conventions\README.md`
- `docs\plans\v1-mvp\phase-1-foundation\README.md`
- `docs\plans\v1-mvp\subagent-work-items\README.md`

後續若繼續延伸：

- `docs\plans\v2-grouping\`
- `docs\plans\v3-notifications-history\`
- `docs\plans\v4-admin-extensions\`

## 建議實作方式

採 **單一 repo、前後端分層**：

- `frontend\`：Angular 20 + Angular Material + plain CSS + pnpm
- `backend\`：Go API + SQLite
- `docs\`：規劃、架構、API 與操作文件

先完成共享領域模型、API 邊界與資料表，再分別落前端功能模組，避免 UI 與後端資料結構脫節。

## 執行階段

### Phase 0 — 開發規範定義

目標：先定義前後端與 cross-cutting 開發規範，讓後續 subagent 與實作者能在同一套工程準則下工作，避免骨架建立後再返工。

詳細文件：

- `docs\plans\v1-mvp\phase-0-conventions\README.md`
- `docs\plans\v1-mvp\phase-0-conventions\shared-engineering-conventions.md`
- `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- `docs\plans\v1-mvp\phase-0-conventions\backend-go-conventions.md`

工作內容：

1. 定義 shared engineering conventions（命名、ID、時間格式、錯誤格式、驗證責任、Definition of Done）
2. 定義 frontend Angular conventions
3. 定義 backend Go conventions
4. 定義 Phase 1 之後的 handoff / subagent 使用方式
5. 定義 subagent 在實作前的文件提案與 review gate 流程

完成條件：

- 前端與後端都有明確可執行的開發規範
- 後續 Phase 1 與 feature work 可直接引用 Phase 0 文件
- Angular 開發規範明確要求使用 Angular CLI 與 Angular CLI MCP best practices，並使用 plain CSS；前端套件管理使用 pnpm；若有 grid table 顯示需求，可使用 ag-grid community
- Go 開發規範明確要求遵循 Google Go style guide、編輯後執行 `gofmt`、補齊測試並確保受影響函式可被測試驗證，且 backend framework 使用 Gin、禁止使用 ORM library
- 進入實作前，規劃文件需先 commit 並 push；各 subagent 需先提交工作文件、待 review 後才可開工

### Phase 1 — 專案骨架與開發基礎

目標：讓 repo 從空白狀態進入可開發、可執行、可測試的基礎狀態。

詳細拆解與 subagent handoff 文件：`docs\plans\v1-mvp\phase-1-foundation\README.md`

前置條件：遵循 `Phase 0` 已定義的工程規範，尤其是 frontend 必須使用 Angular CLI 並依 Angular CLI MCP best practices 開發；backend 必須使用 Gin、不得使用 ORM library，並遵循 Google Go style guide、編輯後執行 `gofmt`，納入可通過的測試。

工作內容：

1. 建立 repo 基礎結構：`frontend\`、`backend\`、`docs\`
2. 使用 pnpm 初始化 Angular 20 專案，使用 CSS styles，加入 Angular Material、路由與基本 layout
3. 初始化 Go module、Gin router、設定檔載入與 API 啟動流程
4. 建立 SQLite 連線與 migration 機制
5. 建立共享文件：
   - 架構說明
   - API 規格草稿
   - 本地開發方式

完成條件：

- 前端可啟動並顯示基礎 shell
- 後端可啟動並成功連接 SQLite
- migration 可建立基本資料表

### Phase 2 — 領域模型、資料表與 API 契約

目標：先把 `Player` / `Session` / `Registration` 的資料模型、驗證規則與 API 邊界定清楚。

工作內容：

1. 依規格建立 SQLite schema
2. 定義 Go models、DTO、validation rules
3. 規劃 REST API：
   - `players`
   - `sessions`
   - `registrations`
   - `reports`
4. 定義狀態轉換規則：
   - `Session`: `open -> closed -> confirmed -> completed / cancelled`
   - `Registration`: `confirmed / cancelled`
5. 補上資料完整性規則：
   - 同球員同場次不可重複報名
   - 名額不可低於已報名人數
   - 差點限制 `0–54`

完成條件：

- migration 與 model 一致
- API 契約可支持 v1 所有主要畫面
- 關鍵商業規則有對應測試

### Phase 3 — 認證與角色基線

目標：建立可支撐 Player / Manager 差異化操作的存取模型。

工作內容：

1. 設計使用者與球員的關聯方式
2. 先落地開發用 mock / stub 身分
3. 保留 LINE OAuth 所需欄位（如 `line_uid`）
4. 建立最小權限控制：
   - Manager：完整 CRUD 與名單調整
   - Player：僅可瀏覽場次與操作自己的報名

完成條件：

- 前後端都能區分 Manager / Player 行為
- 認證替換為 LINE OAuth 時不需重寫主要業務流程

### Phase 4 — 球員管理（v1）

目標：完成球員名冊維護能力。

工作內容：

1. 球員清單、搜尋、狀態篩選
2. 新增 / 編輯球員表單
3. 停用 / 重新啟用球員
4. 邊界處理：
   - 差點非法值驗證
   - 同名警告
   - 停用球員不再出現在新報名選單

完成條件：

- Manager 可完整維護球員資料
- 停用不影響歷史資料

### Phase 5 — 場次管理（v1）

目標：完成場次建立、查詢、編輯、取消與狀態管理。

工作內容：

1. 場次列表：即將到來 / 歷史
2. 建立 / 編輯場次表單
3. 場次詳情頁與即時計算欄位
4. 手動關閉報名、確認名單、取消場次
5. 截止日自動轉 `closed` 的檢查機制

完成條件：

- Manager 可從列表管理完整場次生命週期
- 場次詳情能顯示人數、空位與預計組數

### Phase 6 — 報名系統（v1）

目標：打通球員報名 / 請假與 Manager 手動調整流程。

工作內容：

1. 球員查看場次與剩餘名額
2. `我要報名` / `取消報名`
3. 額滿鎖定、截止鎖定、取消二次確認
4. Manager 手動新增 / 移除報名
5. 防重複報名與容量檢查

完成條件：

- 球員可自助完成報名與請假
- Manager 可處理截止後特殊情況

### Phase 7 — 球場預約摘要（v1）

目標：交付第一個具體業務價值輸出。

工作內容：

1. 報名統計與組數計算
2. 純文字摘要產生器
3. 一鍵複製
4. 列印版樣式
5. Session 確認後的摘要輸出流程串接

完成條件：

- Manager 可從場次詳情一鍵取得可貼到 LINE / 電話紀錄的預約摘要

### Phase 8 — 品質保證、種子資料與上線前準備

目標：讓 v1 MVP 達到可 demo / 可試用的穩定程度。

工作內容：

1. 建立前端單元測試與關鍵互動測試
2. 建立後端 handler / service / repository 測試
3. 建立 seed data（球員、場次、報名）
4. 補齊錯誤訊息、空狀態與 loading 狀態
5. 整理本地開發與部署前置文件

完成條件：

- 關鍵流程可被測試覆蓋
- 團隊可快速建立 demo 環境

## v2+ Backlog（後續版本）

### v2 — 分組與排程

- 蛇形分組 / 隨機分組
- Tee Time 自動排程
- 拖曳調整組別
- WebSocket 即時推送報名人數更新
- 球場預約摘要擴充分組資訊

### v3 — 通知與歷史

- 報名開放通知
- 截止提醒
- 分組公布
- 出席歷史與出席率統計
- 差點追蹤

### v4 — 管理延伸

- 隊費管理
- 隊服尺寸登記與統計

## 不納入第一波交付

- 線上金流
- 成績記錄 / 計分
- 多球隊管理
- 球場 API 直接串接

## Todo 清單

1. `repo-bootstrap`
2. `shared-domain-schema`
3. `auth-foundation`
4. `frontend-shell`
5. `backend-foundation`
6. `players-feature`
7. `sessions-feature`
8. `registrations-feature`
9. `reservation-report`
10. `qa-and-seed-data`
11. `v2-grouping-roadmap`
12. `v3-v4-backlog`

## 相依關係

- `repo-bootstrap` → 所有後續工作
- `shared-domain-schema` 依賴 `repo-bootstrap`
- `frontend-shell`、`backend-foundation`、`auth-foundation` 依賴 `shared-domain-schema`
- `players-feature`、`sessions-feature` 依賴 `frontend-shell` 與 `backend-foundation`
- `registrations-feature` 依賴 `players-feature`、`sessions-feature`、`auth-foundation`
- `reservation-report` 依賴 `registrations-feature`
- `qa-and-seed-data` 依賴 v1 功能完成
- `v2-grouping-roadmap` 依賴 `registrations-feature`

## 風險與注意事項

- **最大風險是認證整合**：LINE OAuth 需要外部設定，容易拖慢第一版交付，建議以 dev stub 先完成主流程。
- **目前 repo 為空**：需先處理腳手架與開發慣例，不是單純 feature 開發。
- **v1 / v2 邊界要守住**：MVP 只需輸出組數，不要過早把拖曳分組與 WebSocket 混進第一版。
- **狀態流轉需集中管理**：避免前端與後端各自判斷造成不一致。

## 建議下一步

若要開始落地，建議先依 `Phase 0` 規範建立共識，再從 `repo-bootstrap` + `shared-domain-schema` 兩項一起展開，先把 repo、資料表、API 契約與前端模組邊界定住，再進入功能實作。
