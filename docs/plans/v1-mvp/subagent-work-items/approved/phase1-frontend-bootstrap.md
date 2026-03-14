# Subagent Task Proposal

## Basic Information

- Phase: Phase 1
- Area: frontend bootstrap
- Proposed task name: phase1-frontend-bootstrap
- Related todo id: `phase1-frontend-bootstrap`
- Assigned subagent: frontend bootstrap agent

## Goal

建立 Angular 20 前端骨架，使 `frontend\` 成為可用 `pnpm` 管理、可由 Angular CLI 擴充、並符合 Angular CLI MCP best practices 的 app shell。

## In Scope

- 使用 Angular CLI 初始化 `frontend`
- 使用 `pnpm` 作為 frontend 套件管理與腳本執行方式
- 使用 plain CSS，不使用 SCSS
- 加入 Angular Material
- 建立 app shell、routing、首頁 placeholder
- 建立 `core` / `shared` / `features` 邊界
- 確保結構相容於 `git worktree` 工作流與 pnpm 共用 `node_modules`

## Out of Scope

- 不實作球員 / 場次 / 報名等真實功能
- 不綁定後端 API
- 不做認證流程
- 無 grid table 需求前，不引入 `ag-grid community`

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `frontend\` 內 Angular workspace 檔案
  - 視需要更新 root `README.md` 的 frontend 啟動指令說明
- 預計新增的資料夾 / 檔案：
  - `frontend\src\app\core\`
  - `frontend\src\app\shared\`
  - `frontend\src\app\features\`
  - `frontend\src\styles.css`
  - Angular CLI 產生的 workspace 設定檔

## Technical Approach

- 使用 Angular CLI 產生 Angular 20 workspace
- 套件安裝與腳本執行使用 `pnpm`
- style 設定維持 CSS-only
- app shell 優先建立最小導覽與路由占位，不提前耦合功能模組
- 依 Angular CLI MCP best practices 使用 standalone、feature-first 與 lazy-load-friendly 結構

- 使用的技術與模式：
  - Angular CLI
  - pnpm
  - Angular Material
  - plain CSS
- 依循的規範文件：
  - `docs\plans\v1-mvp\phase-0-conventions\README.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
  - `docs\plans\v1-mvp\phase-1-foundation\README.md`
  - `docs\plans\v1-mvp\phase-1-foundation\subagents\02-frontend-bootstrap.md`
- 是否新增依賴：
  - Angular framework 與 Angular Material
  - 不主動新增其他 UI / state libraries

## Risks / Open Questions

- 需確認 Angular CLI + pnpm 在目前 repo/worktree 佈局下的初始化方式
- 需避免產出預設 SCSS 設定
- 需確保 workspace 結構與後續 feature 模組邊界一致

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `pnpm install`
  - `pnpm build`
  - `pnpm test` baseline
  - 檢查 styles 與 Angular 設定為 CSS
- 完成後如何驗收：
  - `frontend` 可啟動 / build
  - 有基礎 app shell 與 routes
  - 後續 feature 可在 `core` / `shared` / `features` 結構上直接擴充

## Review Status

- Status: approved
- Reviewer: user approval via `LGTM`
- Review notes: Approved to move into the implementation-ready stage.
