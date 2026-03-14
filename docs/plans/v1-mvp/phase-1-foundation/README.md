# Phase 1 Breakdown — 專案骨架與開發基礎

## 目的

將 `Phase 1` 拆成可獨立分配的工作包，讓後續可由多個 subagent 並行處理，同時維持一致的技術方向、目錄慣例與驗收標準。

本階段所有工作都必須以 `..\phase-0-conventions\README.md` 為前提。

## Phase 1 目標

讓目前幾乎空白的 repo 具備以下能力：

- 有清楚的 top-level repo 結構
- 前端 Angular 20 專案可啟動並顯示基礎 shell
- 後端 Go 服務可啟動並提供 health endpoint
- SQLite 連線與 migration 機制可運作
- 文件足以讓後續功能開發接手

## Phase 1 不做的事

- 不實作球員、場次、報名等完整業務功能
- 不做 LINE OAuth 正式串接
- 不做 WebSocket
- 不進入 v2 分組功能

## 建議技術決策

- Frontend：Angular 20、standalone API、Angular Material、plain CSS、pnpm
- Backend：Go、`gin` framework、分層為 `cmd` / `internal` / `migrations`，且不使用 ORM
- SQLite driver：優先使用 pure-Go driver，避免 Windows 本機開發被 CGO 阻塞
- Migration：使用 SQL migration files + 輕量 runner，不在 Phase 1 引入過重框架

### 強制規範

- Frontend 必須使用 **Angular CLI**
- Frontend 必須以 **Angular CLI MCP best practices** 為準
- Frontend style 必須使用 **plain CSS**，不使用 SCSS
- Frontend 套件管理必須使用 **pnpm**
- 若有 grid table 顯示需求，可使用 **ag-grid community**
- Backend 必須遵循 **Google Go style guide**
- Backend 每次編輯後都必須執行 `gofmt`
- Backend 變更必須包含測試，並確保受影響函式可被測試驗證
- Backend framework 必須使用 **Gin**
- Backend 不得使用 **ORM library**

## 建議目錄結構

```text
.
├── frontend\
├── backend\
│   ├── cmd\
│   ├── internal\
│   ├── migrations\
│   └── data\
├── docs\
│   ├── plans\
│   ├── architecture\
│   └── development\
└── .gitignore
```

## Workstream 拆解

| Workstream | 目標 | 建議 subagent | 詳細文件 |
|---|---|---|---|
| W1 | 建立 repo 骨架與共用慣例 | general-purpose / task | `subagents\01-workspace-scaffold.md` |
| W2 | 建立 Angular app shell | general-purpose / task | `subagents\02-frontend-bootstrap.md` |
| W3 | 建立 Go API 骨架 | general-purpose / task | `subagents\03-backend-bootstrap.md` |
| W4 | 建立 SQLite 與 migration 基線 | general-purpose / task | `subagents\04-sqlite-migrations.md` |
| W5 | 建立文件與開發流程 | general-purpose | `subagents\05-docs-and-devx.md` |

## 建議執行順序

1. W1 `workspace-scaffold`
2. W2 `frontend-bootstrap` 與 W3 `backend-bootstrap` 可並行
3. W4 `sqlite-migrations` 在 W3 完成後接續
4. W5 `docs-and-devx` 在 W1 完成後即可進行，並在 W2/W3/W4 結束後補齊細節
5. 最後由主 agent 做整合驗收

## Subagent 分配建議

### A. Repo Scaffold Agent

- 先建立 top-level 結構、ignore 規則、共用命名與基礎說明
- 產出後，其他 subagent 才能在穩定結構上工作

### B. Frontend Agent

- 專注 Angular CLI 產生專案與 Material shell
- 不自行決定後端 API 格式，只先做 app shell、routing、layout、core/shared/features 邊界

### C. Backend Agent

- 專注 Go server、router、config、health endpoint 與 package layout
- 不搶先實作完整業務 handler

### D. DB Agent

- 專注 SQLite driver、DB 初始化、migration runner、最小 smoke schema
- 不在 Phase 1 就建立完整業務 schema

### E. Docs / DevEx Agent

- 整理本地啟動方式、架構邊界、Phase 1 決策與後續交接須知

## 共用約束

- 一律使用 Windows-friendly 路徑與命令
- 盡量採低依賴方案，避免在骨架階段過度複雜化
- 任何新增結構都要與後續 `shared-domain-schema` 相容
- Frontend 與 backend 在 Phase 1 只需 smoke-level 可執行，不需完整串接
- subagent 預設以 `git worktree` 模式作業，Phase 1 輸出需相容於多 worktree 並行開發
- frontend 依賴管理需支援透過 pnpm 在 worktree 間共用 `node_modules`
- 實作開始前，需先完成 docs commit / push，並讓各 subagent 先提交工作文件等待 review

## 驗收標準

- `frontend` 可成功安裝依賴並啟動 / build
- `backend` 可成功啟動並提供 health endpoint
- migration runner 可對 SQLite 檔建立至少一個 smoke migration
- 文件說明清楚各目錄用途、啟動方式與下一階段接口

## 整合驗收清單

- Top-level 結構不互相衝突
- Angular 專案符合 standalone + feature-first 邊界
- Go 專案可支撐後續 API、middleware、auth、repository 擴充
- SQLite 方案在 Windows 本機可行
- 文件與實際結構一致

## 後續關聯

完成 Phase 1 後，下一步接續：

- `shared-domain-schema`
- `auth-foundation`
- `players-feature`
- `sessions-feature`
- `registrations-feature`

因此 Phase 1 的輸出要特別注意：

- 不要把業務邏輯綁死在 bootstrap code
- 保留清楚的 domain / transport / persistence 邊界
- 若由 subagent 接手實作，先在 `..\subagent-work-items\pending\` 建立提案文件
