# Current Workflow

這份文件整理目前 repo 的開發工作流程，作為根目錄下的快速入口。

## 目的

- 統一目前 planning、review、approval、implementation 的作業順序
- 提供 subagent / 開發者在開始工作前的共用流程說明
- 快速指向常用命令與關鍵目錄

## 目前技術基線

- Frontend: Angular + Angular Material + plain CSS + pnpm
- Backend: Go + Gin + SQLite
- 協作模式: 預設使用 `git worktree`

## 開工前的必要條件

1. 先同步目前規劃文件，確認最新基線：
   - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
   - `docs\plans\v1-mvp\phase-0-conventions\`
   - `docs\plans\v1-mvp\phase-1-foundation\`
2. 規劃 / 規範文件必須先完成並 `commit` + `push`
3. 確認要執行的工作項目已完成 review gate，且符合目前規範

## 標準工作流程

### 1. 建立工作文件

每個 subagent 在開始實作前，都要先建立工作文件，放在：

- `docs\plans\v1-mvp\subagent-work-items\pending\`

工作文件至少應包含：

- 目標與範圍
- 不在範圍內的內容
- 相依性（若有，必須明確標示）
- 預計修改的檔案 / 目錄
- 技術決策與依據
- 風險 / 待確認事項
- 驗收方式

### 2. 等待 review

- 工作文件建立後先等待 review
- 未經 review 核可，不得開始實作
- 若 scope 變更，需更新工作文件並重新 review

### 3. 移到 approved

只有在**使用者明確指示**後，才可將工作文件移到：

- `docs\plans\v1-mvp\subagent-work-items\approved\`

### 4. 先 commit approval，再開始實作

- 文件移到 `approved` 後，必須先將核可狀態 `commit`
- 完成 approval commit 後，才可在 `git worktree` 環境下開始實作
- 實作應維持每個 task / subagent 的 worktree 隔離

### 5. 完成實作後收尾

- 任務完成後，將工作文件移到：
  - `docs\plans\v1-mvp\subagent-work-items\completed\YYYY-MM-DD\`
- 交接時需說明：
  - 本次改動影響的路徑
  - 是否新增依賴、命令或限制
  - 若有偏離規範，需明確說明原因

## 目錄狀態流轉

```text
docs\plans\v1-mvp\subagent-work-items\
├── pending\        # 等待 review
├── approved\       # 已核可且應先完成 approval commit
└── completed\
    └── YYYY-MM-DD\ # 已完成的工作文件
```

## 協作規則

- subagent 預設使用 `git worktree` 模式作業
- 工作應盡量依 task 隔離，避免不同實作互相干擾
- 規劃與文件需能支援多個 worktree 並行使用
- 如果任務 scope 改變，不要直接硬做；先更新文件並重跑 review 流程

## Frontend workflow

- 使用 `pnpm` 管理 frontend 套件
- Angular 相關工作使用 Angular CLI
- Angular 實作以 Angular CLI MCP best practices 為準
- 樣式使用 plain CSS，不使用 SCSS
- 若需要 grid table UI，可使用 `ag-grid community`
- frontend worktree 之間預期透過 pnpm 共用 `node_modules`

## Backend workflow

- 使用 Gin 作為 backend framework
- 不使用 ORM
- 遵循 Google Go style
- 每次修改 Go 程式後執行 `gofmt`
- backend 變更要附帶測試，並保持相關函式可測試

## 常用命令

優先使用 root `justfile`：

```powershell
just status
just worktrees
just plans
just approved
just pending
just completed
just frontend-install
just frontend-build
just frontend-test
just backend-test
just backend-migrate
just backend-start
```

## 目前狀態

- 目前 repo 主要聚焦於 Phase 1 bootstrap 與整體基線建置
- 下一個 ready task 為 `phase1-integration-check`
- backend 預設使用本機 SQLite：
  - `backend\data\golf_team_manager.sqlite`

## 參考來源

- `README.md`
- `docs\development\local-setup.md`
- `docs\plans\v1-mvp\subagent-work-items\README.md`
- `docs\plans\v1-mvp\phase-0-conventions\shared-engineering-conventions.md`
