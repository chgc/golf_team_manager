# Subagent Brief — Frontend Bootstrap

## 任務目標

建立 Angular 20 前端骨架，使專案具備可啟動、可 build、可承接後續 feature modules 的 app shell。

此任務必須先閱讀 `..\..\phase-0-conventions\frontend-angular-conventions.md`。

## 範圍內

- 使用 Angular CLI 初始化 `frontend`
- 啟用 routing、CSS、standalone
- 加入 Angular Material
- 建立 `core` / `shared` / `features` 邊界
- 建立基礎 layout、首頁 placeholder、空路由骨架
- 使用 pnpm 管理 frontend 相依套件

## 範圍外

- 不實作球員 / 場次 / 報名真實 UI
- 不綁定真實 API
- 不做認證流程

## 建議檔案方向

- `frontend\src\app\core\`
- `frontend\src\app\shared\`
- `frontend\src\app\features\`
- `frontend\src\app\app.routes.ts`
- `frontend\src\styles.css`

## 建議步驟

1. 用 Angular CLI 建立 app
2. 用 pnpm 管理安裝與腳本執行
3. 安裝與設定 Angular Material
4. 建立主 layout：toolbar、side nav 或 top nav
5. 設定首頁與基礎路由占位
6. 建立後續可延伸的資料夾慣例與 shared UI placeholder
7. 驗證 `pnpm` 下的 build 與 test baseline 可跑通

## 設計限制

- 使用 standalone APIs
- 採 feature-first 結構，不用 NgModule-first 老式組織
- Phase 1 僅需 shell，不需 state management 複雜方案
- 以 Angular CLI MCP best practices 作為準則來源
- 後續元件、service、route scaffold 優先使用 Angular CLI 指令建立
- style 使用 plain CSS，不使用 SCSS
- 有 grid table 顯示需求時，才引入 ag-grid community
- 需相容於 subagent 以 git worktree 模式開發，並透過 pnpm 共用 `node_modules`

## 驗收標準

- `frontend` 可啟動
- 有清楚 layout 與 routing 基線
- 後續 feature 可直接在 `features` 下擴充

## 交接備註

- 對 backend API 僅保留 service placeholder，不要硬寫 endpoint 細節
