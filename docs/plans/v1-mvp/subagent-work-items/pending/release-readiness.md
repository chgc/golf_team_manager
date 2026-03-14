# Subagent Task Proposal

## Basic Information

- Phase: Phase 8
- Area: release readiness
- Proposed task name: release-readiness
- Related todo id: `release-readiness`
- Assigned subagent: release readiness agent

## Goal

在 v1 MVP 功能與 seed data 完成後，補齊 demo / handoff / pre-release 檢查清單，讓團隊可以用一致流程驗證並交接目前版本。

## In Scope

- 整理 v1 MVP 的 release readiness checklist
- 盤點 demo 前必要的命令、文件與 smoke path
- 補齊 handoff / operator-facing 文件缺口
- 收斂目前已知限制、dev-only 約束與後續 follow-up 清單

## Out of Scope

- 不新增新功能
- 不實作 CI/CD pipeline
- 不處理 production deployment automation

## Dependencies

- `qa-and-seed-data` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `README.md`
  - `docs\development\`
  - `docs\architecture\`
  - 視需要更新 root workflow / handoff docs
- 預計新增的資料夾 / 檔案：
  - release readiness checklist / handoff docs（待 review 收斂）

## Technical Approach

- 使用的技術與模式：
  - 以文件與既有 smoke command 為主，不新增執行環境依賴
  - 以 v1 已完成功能為範圍，不提前混入 v2+ backlog
- 依循的規範文件：
  - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
  - `WORKFLOW.md`
  - `README.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- release readiness 文件需避免與 local setup / demo smoke 文件重複，需明確定義各文件角色

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - 以 review 後收斂
- 完成後如何驗收：
  - 以 review 後收斂

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
- Agent review summary:

## Feedback

- Reviewer agent 1:
- Reviewer agent 2:
- Applied proposal updates:
