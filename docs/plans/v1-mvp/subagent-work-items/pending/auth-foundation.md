# Subagent Task Proposal

## Basic Information

- Phase: Phase 3
- Area: auth foundation
- Proposed task name: auth-foundation
- Related todo id: `auth-foundation`
- Assigned subagent: auth foundation agent

## Goal

建立可支撐 Player / Manager 差異化操作的 auth baseline，先以 development-friendly stub identity 與 role model 開路，同時保留未來 LINE OAuth 接軌點。

## In Scope

- 定義 user / player relationship baseline
- 建立 development stub identity flow
- 建立 role model（manager / player）
- 保留 LINE OAuth 所需欄位與接軌點
- 對齊 frontend shell 與 backend API foundation 的 auth extension points

## Out of Scope

- 不完成正式 LINE OAuth 整合
- 不完成完整 session-based 或 token-based production auth hardening
- 不進入 feature-complete authorization rules

## Dependencies

- `shared-domain-schema` must be completed first

## Planned Changes

- 預計修改的資料夾 / 檔案：
  - `backend\internal\`
  - `frontend\src\app\`
  - 視需要更新 schema / architecture 文件
- 預計新增的資料夾 / 檔案：
  - auth baseline 相關 service / model / docs
  - dev stub identity wiring

## Technical Approach

- 使用的技術與模式：
  - backend auth abstraction with dev stub
  - frontend session/identity shell state
  - future LINE OAuth compatibility
- 依循的規範文件：
  - `docs\plans\v1-mvp\golf-team-manager-implementation-plan.md`
  - `docs\architecture\shared-domain-baseline.md`
  - `docs\architecture\backend-api-foundation.md`
  - `docs\architecture\frontend-shell-baseline.md`
- 是否新增依賴：
  - 原則上不新增；先以低依賴 stub flow 為主

## Risks / Open Questions

- 需避免把 dev stub 寫成未來正式 auth 的阻礙
- 需讓 player identity 與 player entity 的關聯足夠清楚
- 需讓 role-based 行為差異能支撐後續 feature work

## Validation Plan

- 會執行哪些 build / test / smoke checks：
  - `go test ./...`
  - `just frontend-test`
  - auth stub flow smoke checks
- 完成後如何驗收：
  - manager / player baseline 可被前後端辨識
  - 後續 features 可依角色擴充
  - 未來 LINE OAuth 有明確接點

## Review Status

- Status: pending-review
- Reviewer:
- Review notes:
