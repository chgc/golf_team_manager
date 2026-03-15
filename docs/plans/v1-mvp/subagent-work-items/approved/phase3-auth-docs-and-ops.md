# Subagent Task Proposal

## Basic Information

- Phase: Phase 3 follow-up
- Area: auth docs and ops
- Proposed task name: auth-docs-and-ops
- Related todo id: `auth-approved-handoff`
- Assigned subagent: auth ops agent

## Goal

補齊 auth implementation 落地後需要的開發文件、環境變數範本、local setup、demo smoke 檢查與 handoff 說明，確保後續 subagent 與開發者能穩定使用新 auth flow。

## In Scope

- 更新 local setup 文件中的 auth 啟動方式
- 新增或更新 `.env.example` 需求說明
- 更新 demo/smoke 文件中的 auth 驗證步驟
- 記錄 `dev_stub` 與 `line` mode 切換方式
- 記錄 local `line` mode 的 env、host、cookie/CORS 與 callback 設定
- 記錄 known limitations 與 rollback / fallback 做法

## Out of Scope

- 不實作 backend 或 frontend 主要 auth 邏輯
- 不處理雲端 deployment secrets management

## Dependencies

- `phase3-auth-implementation-detail-alignment.md`
- `phase3-backend-line-oauth-jwt.md`（需已完成實作與驗證）
- `phase3-frontend-auth-flow.md`（需已完成實作與驗證）

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `README.md`
  - `.env.example`
  - `docs\development\local-setup.md`
  - `docs\development\demo-smoke-check.md`
  - `docs\development\release-readiness-checklist.md`
- 預計新增的資料夾 / 檔案：
  - 視需要新增 auth-specific setup notes

## Technical Approach

- 使用的技術與模式：
  - docs-first operational handoff
  - 以 `.env.example` 描述必要設定，避免 secrets 進版控
  - 對齊 backend/frontend 實際落地行為
- 依循的規範文件：
  - `WORKFLOW.md`
  - `docs\architecture\auth-line-sso-implementation-detail.md`
  - `docs\plans\v1-mvp\subagent-work-items\README.md`
- 是否新增依賴：
  - 否

## Risks / Open Questions

- 若 backend/frontend 最終 contract 與 implementation detail 有差異，文件容易失真
- `.env.example` 與 local setup 必須同步，不然 demo flow 會失敗
- 若未在 backend/frontend 驗證完成後再更新文件，release/demo 文件會提前固化錯誤流程

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - 文件命令與實際命令交叉驗證
  - auth smoke checklist walkthrough
- 完成後如何驗收：
  - 新開發者可依文件啟用 `dev_stub` 或 `line` mode
  - demo / smoke 文件足以覆蓋新的 auth flow

## Review Status

- Status: approved
- Reviewer: reviewer agent 1, reviewer agent 2
- Review notes: First-pass review required docs sequencing to depend on completed-and-validated backend/frontend auth work and asked for local line-mode operational guidance. Proposal was updated and approved on second pass.
- Agent review summary: Approved by two independent reviewer agents; only a minor `.env.example` planning nit remained and was incorporated before approval.

## Feedback

- Reviewer agent 1: blocking issue — please make workflow sequencing explicit. This task must depend on backend/frontend auth implementation being completed and validated (not only on the proposal docs existing), so the docs reflect the final landed contract and operational flow.
- Reviewer agent 2: Blocking: docs/ops 任務應依賴 backend/frontend auth 實作完成且 smoke 驗證通過後再落地，否則文件容易和最終行為脫節。請在 scope 中明確納入 local line mode 的 env、host、cookie/CORS 與 fallback/rollback 指引。
- Applied proposal updates: tightened dependencies to completed-and-validated backend/frontend work, added local line-mode operational guidance, and made rollout/fallback documentation explicit.
