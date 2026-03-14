# Subagent Task Proposal

## Basic Information

- Phase: Phase 2
- Area: frontend shell
- Proposed task name: frontend-shell
- Related todo id: `frontend-shell`
- Assigned subagent: frontend shell agent

## Goal

在既有 Angular shell 基線上，承接 shared domain schema 與 backend API foundation，整理出 feature-aligned 的頁面結構、service 邊界與基本資料流模式，讓 players / sessions / registrations 畫面能在一致前端架構上開發。

## In Scope

- 對齊 frontend route / feature structure 與 shared domain vocabulary
- 建立 API service / data-access baseline
- 建立 shared UI / layout 邊界與頁面骨架
- 為 players / sessions / registrations 預留清楚 feature 入口

## Out of Scope

- 不完成完整 players / sessions / registrations UI
- 不引入 SCSS 或額外 UI framework
- 不實作正式 auth flow

## Dependencies

- `shared-domain-schema` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `frontend\src\app\`
  - 視需要更新 `README.md` 與開發文件
- 預計新增的資料夾 / 檔案：
  - feature-aligned services / models / data-access baseline
  - shared UI shell / page scaffold
  - route wiring 與對應測試

## Technical Approach

- 使用的技術與模式：
  - Angular CLI generated workspace
  - standalone components
  - feature-first routing and service structure
  - plain CSS only
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\plans\v1-mvp\phase-0-conventions\frontend-angular-conventions.md`
- 是否新增依賴：
  - 原則上不新增；若未來有 grid table 需求才考慮 ag-grid community

## Risks / Open Questions

- 需避免在 shell 階段就把 feature UI 做到過深
- 需讓 data-access 邊界與 backend API foundation 對齊
- 需保持 Angular CLI / MCP best practices 與 plain CSS 約束

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `just frontend-build`
  - `just frontend-test`
  - route / shell rendering checks
- 完成後如何驗收：
  - frontend 架構可承接 feature pages
  - data-access 邊界清楚
  - 與 backend/shared-domain vocabulary 一致

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
