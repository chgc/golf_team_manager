# Phase 0 — 開發規範定義

## 目的

在正式開始 `Phase 1` 骨架建置前，先定義 shared、frontend、backend 的工程規範，讓後續 subagent 可以在一致的設計邊界、命名方式與交付標準下工作。

## 產出文件

- `shared-engineering-conventions.md`
- `frontend-angular-conventions.md`
- `backend-go-conventions.md`

## Phase 0 的定位

- `Phase 0` 不直接產生產品功能
- `Phase 0` 的任務是降低後續返工與整併成本
- `Phase 1` 之後所有工作都應以這裡的規範作為預設準則

## 共用要求

- 規範優先於局部實作偏好
- 後續 subagent 若需偏離規範，必須在交付說明中明確說明理由
- 新增規範時，應同步更新主計畫與對應 handoff 文件
- 在任何 implementation 開始前，規劃 / 規範文件必須先 commit 並 push 到遠端，作為可審查基線

## Angular 特別要求

前端 Angular 開發有以下強制規則：

1. 必須使用 **Angular CLI**
2. 必須使用 **Angular CLI MCP 提供的 best practices** 作為前端開發準則
3. 在實際 Angular workspace 建立後，應優先讀取 workspace 對應版本的 Angular best practices，再進行實作
4. style 使用 **plain CSS**，不使用 SCSS
5. 若有 grid table 顯示需求，可使用 **ag-grid community**
6. 前端套件管理使用 **pnpm**

## Go 特別要求

後端 Go 開發有以下強制規則：

1. 必須遵循 **Google Go style guide**
2. 每次編輯 Go 程式後，必須執行 `gofmt`
3. 必須包含測試
4. 受影響函式必須可被測試驗證，不能留下明顯不可測的核心邏輯
5. backend framework 使用 **Gin**
6. **禁止使用 ORM library**

## Subagent 使用方式

- 接手 frontend 任務前，先閱讀 `frontend-angular-conventions.md`
- 接手 backend 任務前，先閱讀 `backend-go-conventions.md`
- 涉及 API 契約、錯誤格式、時間欄位、測試要求時，先閱讀 `shared-engineering-conventions.md`
- subagent 開發功能時，預設使用 **git worktree** 模式
- subagent 在開始實作前，必須先提交工作文件供 review，核可後才可開工

## 與 Phase 1 的關係

`Phase 1` 的 repo 骨架、Angular app shell、Go server、SQLite migration 與開發文件，都必須遵循 `Phase 0` 已定義規範。
