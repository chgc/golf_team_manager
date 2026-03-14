# Subagent Task Proposal

## Basic Information

- Phase: Phase 2
- Area: backend foundation
- Proposed task name: backend-foundation
- Related todo id: `backend-foundation`
- Assigned subagent: backend API foundation agent

## Goal

建立承接 shared domain baseline 的後端 API foundation，包含 repository / service / handler 邊界、錯誤回應格式、基礎 middleware 與第一批 domain-aligned route wiring，讓後續 players / sessions / registrations feature work 能在一致的 Gin API 基線上展開。

## In Scope

- 建立 repository / service / handler package baseline
- 建立 API error response shape 與 request validation flow
- 建立第一批 `/api` route wiring
- 串接 shared domain models / DTO / validation baseline
- 保留 auth integration extension point

## Out of Scope

- 不完成完整 players / sessions / registrations CRUD
- 不實作正式 LINE OAuth 流程
- 不進入前端畫面開發

## Dependencies

- `shared-domain-schema` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\cmd\api\`
  - `backend\internal\http\`
  - `backend\internal\`
  - 視需要更新開發與架構文件
- 預計新增的資料夾 / 檔案：
  - repository / service / handler baseline 檔案
  - API error / response helper
  - 第一批 route registration 與測試

## Technical Approach

- 使用的技術與模式：
  - Gin
  - repository / service / handler layering
  - shared domain DTO validation at transport boundary
  - explicit JSON error response format
- 依循的規範文件：
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\development\phase-1-validation.md`
  - `docs\plans\v1-mvp\phase-0-conventions\backend-go-conventions.md`
- 是否新增依賴：
  - 原則上不新增

## Risks / Open Questions

- 需避免在 foundation 階段就把 handlers 寫成 feature-complete CRUD
- 需讓未來 auth middleware 能無痛接上
- 需讓 repository 邊界與 migration/schema vocabulary 一致

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - backend startup smoke check
  - API route / response baseline test
- 完成後如何驗收：
  - API foundation 可承接 feature handlers
  - validation / error handling 基線清楚
  - shared domain schema 已被後端 transport 層正確引用

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
