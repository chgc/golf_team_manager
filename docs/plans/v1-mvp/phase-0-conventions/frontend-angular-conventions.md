# Frontend Angular Conventions

## 強制規則

- 前端開發 **必須使用 Angular CLI**
- 前端開發 **必須以 Angular CLI MCP 提供的 best practices 為準**
- 若有 Angular 版本上下文，應優先讀取對應 workspace 的 Angular best practices
- style 使用 **plain CSS**，不使用 SCSS
- 若有 grid table 顯示需求，可使用 **ag-grid community**
- 前端套件管理 **必須使用 pnpm**

## 工具與初始化

- 使用 Angular CLI 建立與擴充 workspace、component、service、guard、route scaffold
- 不手工拼湊初始 Angular 專案結構來取代 CLI
- Angular UI 套件以 Angular Material 為預設
- 預設 style 副檔名使用 `.css`
- 只有在有 grid table 顯示需求時才引入 `ag-grid-community`
- 安裝、執行與腳本管理以 `pnpm` 指令為主

## Angular 架構原則

- 使用 standalone APIs
- 採 feature-first 結構
- feature route 優先使用 lazy loading
- 共享能力放在 `core` / `shared`，避免 feature 間直接互相耦合

建議結構：

```text
src\app\
├── core\
├── shared\
└── features\
```

## 元件與狀態

- 元件保持單一職責
- 預設使用 `ChangeDetectionStrategy.OnPush`
- 使用 `input()` / `output()`，避免舊式 decorator 風格
- local state 優先使用 signals
- derived state 使用 `computed()`
- 不使用 `mutate`

## Template 規範

- 使用 Angular 原生控制流：`@if`、`@for`、`@switch`
- 不使用 `*ngIf`、`*ngFor`、`*ngSwitch` 作為新程式碼預設
- 不使用 `ngClass`，改用 `class` bindings
- 不使用 `ngStyle`，改用 `style` bindings
- template 保持簡潔，不在模板中塞入複雜邏輯

## 表單與資料流

- 表單預設使用 Reactive Forms
- service 負責資料存取與 API 互動
- component 負責畫面協調與使用者互動
- 在 Phase 1 僅建立 shell 與 placeholder，不要過早耦合真實 API 契約

## 相依注入與服務

- 使用 `inject()` 而非 constructor injection 作為預設
- singleton service 使用 `providedIn: 'root'`
- service 應維持單一責任

## Accessibility 與品質

- 必須通過 AXE 檢查
- 必須符合 WCAG AA 最低要求
- 靜態圖片使用 `NgOptimizedImage`

## Phase 1 特別要求

- `frontend` 必須由 Angular CLI 初始化
- routing、CSS、Angular Material 與 app shell 必須在 Phase 1 建好
- Phase 1 僅建立 layout、routes、shared conventions、placeholder pages
- 初始化與後續開發流程需相容於 `git worktree` 模式，並可透過 pnpm 共用 `node_modules`

## Subagent 交付要求

- 說明使用了哪些 Angular CLI 指令
- 說明使用了哪些 `pnpm` 指令
- 說明是否依據 Angular CLI MCP best practices 調整結構
- 說明是否維持 CSS-only 設定，以及是否引入 ag-grid community
- 若遇到規範衝突，先以 Angular CLI MCP best practices 為優先
