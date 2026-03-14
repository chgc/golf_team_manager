# Subagent Task Proposal

## Basic Information

- Phase: Phase 8
- Area: release readiness
- Proposed task name: release-readiness
- Related todo id: `release-readiness`
- Assigned subagent: release readiness agent

## Goal

在 v1 MVP 功能與 seed data 完成後，補齊 release readiness / demo handoff / operator-facing 文件，讓團隊可用單一且不互相重複的文件集合完成本地啟動、demo smoke、pre-release 檢查與交接。

## In Scope

- 整理 v1 MVP 的 release readiness checklist，作為 demo 前與交接前的總入口
- 補齊 demo / handoff / operator-facing 文件缺口，但只限於 v1 已存在能力與流程
- 盤點並收斂既有文件角色，避免 `README.md`、`local-setup.md`、`demo-smoke-check.md` 與新文件內容重複
- 收斂目前已知限制、dev-only 約束、手動 smoke path 與 post-MVP follow-up 清單
- 驗證目前 repo 的既有 build / test / seed / smoke 指令仍可作為 release readiness baseline

## Out of Scope

- 不新增新功能
- 不實作 CI/CD pipeline
- 不處理 production deployment automation
- 不新增 e2e framework、外部監控或正式環境 deployment guide
- 不擴寫 architecture 決策文件，除非僅需補充連結或簡短 cross-reference

## Dependencies

- `qa-and-seed-data` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `README.md`
  - `docs\development\local-setup.md`
  - `docs\development\demo-smoke-check.md`
  - `WORKFLOW.md`（若需要補 release readiness / handoff 入口連結）
- 預計新增的資料夾 / 檔案：
  - `docs\development\release-readiness-checklist.md`
  - `docs\development\v1-handoff-summary.md`
- 文件角色定義：
  - `README.md`：repo 首頁與 quick links，只保留高層入口，不重複詳細步驟
  - `docs\development\local-setup.md`：本地開發環境、常用命令、seed / startup 入口
  - `docs\development\demo-smoke-check.md`：deterministic dataset、manager/player smoke path、demo 操作細節
  - `docs\development\release-readiness-checklist.md`：pre-demo / pre-release checklist、驗證矩陣、已知限制與 sign-off 項目
  - `docs\development\v1-handoff-summary.md`：交接摘要、操作入口、dev-only 限制、後續 follow-up 與 ownership 提示

## Technical Approach

- 使用的技術與模式：
  - 以文件與既有 smoke command 為主，不新增執行環境依賴
  - 以 v1 已完成功能為範圍，不提前混入 v2+ backlog
  - 以「單一文件單一責任」方式整理文件：release checklist 負責總覽與 gate，既有文件保留細節操作
  - checklist / handoff 文件優先引用既有命令與 smoke path，不複製同一段操作內容到多處
  - 若需更新 workflow / README，只補充入口與文件定位，不把詳細 release 步驟塞回 root 文件
- 依循的規範文件：
  - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
  - `WORKFLOW.md`
  - `README.md`
  - `docs\development\local-setup.md`
  - `docs\development\demo-smoke-check.md`
  - `docs\architecture\qa-and-seed-data.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- release readiness / handoff 文件若直接複製 `local-setup.md` 或 `demo-smoke-check.md` 內容，後續維護容易漂移；需以 cross-link + ownership 邊界解決
- 需避免把 dev-only smoke / seed 限制誤寫成 production-ready 操作建議
- 若 review 認為 `v1-handoff-summary.md` 命名不夠貼近團隊習慣，可在 re-review 時調整名稱，但需保留其「交接摘要」職責不變

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `just backend-test`
  - `just frontend-build`
  - `just frontend-test`
  - `just backend-seed`
  - 依 `docs\development\demo-smoke-check.md` 執行 manager smoke path
  - 依 `docs\development\demo-smoke-check.md` 執行 player debug-header API smoke path
- 完成後如何驗收：
  - Case 1: `README.md`、`local-setup.md`、`demo-smoke-check.md`、`release-readiness-checklist.md`、`v1-handoff-summary.md` 的角色分工明確，reviewer 可指出每份文件的唯一主要用途
  - Case 2: 新進成員可從 `README.md` 快速找到 local setup、demo smoke、release checklist 與 handoff summary 的 canonical 入口
  - Case 3: `release-readiness-checklist.md` 明確列出 release/demo 前需確認的命令、人工 smoke path、已知限制與 sign-off 項目，且引用既有文件而非大段重複
  - Case 4: `v1-handoff-summary.md` 明確整理目前版本範圍、dev-only 條件、demo/operator 入口與 follow-up backlog 邊界
  - Case 5: `just backend-test`、`just frontend-build`、`just frontend-test`、`just backend-seed` 保持通過，且 checklist 中引用的命令與 repo 實際命令一致
  - Case 6: manager / player smoke path 可依 `demo-smoke-check.md` 重現，且新 checklist / handoff 文件對這兩條 smoke path 的描述不與現有文件衝突

## Review Status

- Status: approved
- Reviewer: reviewer agents (GPT-5.4 / Claude Sonnet 4.6 across 2 rounds)
- Review notes: 首輪 review 為 blocking，主因是 deliverables、文件邊界與 validation plan 過於模糊；更新後重新進行第二輪雙 reviewer review，兩位 reviewer 皆確認 proposal 已補齊具名交付檔案、文件角色分工、具體 validation commands 與 outcome-based acceptance matrix，且與 `WORKFLOW.md`、`README.md`、`local-setup.md`、`demo-smoke-check.md`、`qa-and-seed-data.md` 一致，已達 implementation-ready。
- Agent review summary: round 1 blocking -> round 2 approve + approve

## Feedback

- Reviewer agent 1:
  - proposal 方向正確且相依性已滿足，但首版仍未達 implementation-ready。主要 blocking 在於 planned changes 未具名、文件邊界未定義、validation plan 仍是 placeholder，且 handoff 交付內容過於寬泛。
- Reviewer agent 2:
  - 已要求 proposal 明確列出新舊文件分工、具體 build/test/smoke commands、可 review 的 acceptance matrix，以及是否真的需要碰 architecture docs。若不補齊，進入 approved 後仍容易各自解讀。
- Reviewer agent 3:
  - 第二輪 review 確認 proposal 已無 material blocking issue。新文件已具名，文件 ownership 邊界清楚，validation plan 與 acceptance criteria 可獨立驗證，且 dependency `qa-and-seed-data` 已滿足，可進入 approved。
- Reviewer agent 4:
  - 第二輪 review 確認 proposal 已符合 review gate 要求。scope 維持在 v1 文件整備，不擴張到新 runtime / deployment 工作；release readiness checklist 與 handoff summary 的職責分離也足夠清楚，可開始 approval 流程。
- Applied proposal updates:
  - 將新文件具名固定為 `docs\development\release-readiness-checklist.md` 與 `docs\development\v1-handoff-summary.md`
  - 明確定義 `README.md`、`local-setup.md`、`demo-smoke-check.md`、release checklist、handoff summary 的文件邊界與角色
  - 從 planned changes 移除預設的 `docs\architecture\` 廣泛修改範圍，僅保留必要 cross-reference 的可能性
  - 補齊具體 validation commands：`just backend-test`、`just frontend-build`、`just frontend-test`、`just backend-seed`，以及 manager/player smoke path
  - 將 acceptance 改寫為可 review 的 outcome-based matrix，聚焦 canonical entry points、文件不重複、dev-only 限制明文化與 handoff 交付完整性
