# Shared Engineering Conventions

## 適用範圍

本文件適用於 frontend、backend、database 與跨模組整合工作。

## 命名與結構

- 目錄與檔名優先使用 `kebab-case`
- TypeScript 型別、Go struct、介面與 exported symbols 使用 `PascalCase`
- 變數與函式使用語言慣例：TypeScript `camelCase`、Go `camelCase` / exported `PascalCase`
- 文件名稱應明確表達用途，不使用模糊名稱如 `misc`、`temp`

## ID 與時間格式

- 主要實體 ID 統一使用字串型 UUID
- API 與儲存層的時間欄位統一使用 UTC
- 對外交換資料時，時間格式使用 ISO 8601 字串

## 資料契約

- domain model、API DTO、資料庫 schema 需明確分層，不混用
- frontend 不直接依賴資料庫 schema 命名
- backend handler 不直接暴露資料庫列結構給前端

## 驗證責任

- frontend：即時表單驗證、基本 UX 錯誤提示
- backend：最終權威驗證、狀態轉換與資料完整性保證
- database：以 constraint 守住最後一道資料一致性邊界

## 錯誤處理

- API 錯誤需提供一致格式
- 需區分 validation、not found、conflict、unauthorized、internal error
- 不允許以成功形狀包裝失敗結果

建議錯誤格式：

```json
{
  "error": {
    "code": "validation_error",
    "message": "maxPlayers must be greater than confirmed registrations"
  }
}
```

## 測試與完成定義

- 新增結構或功能時，至少要有對應 smoke 驗證方式
- 修改契約時，需同步更新相關文件
- 任務完成前需確認：
  - 結構與命名符合規範
  - 文件可讓下一位接手者理解
  - 沒有把臨時方案寫成長期依賴

## Subagent 工作模式

- subagent 在開發功能時，預設使用 `git worktree` 模式進行隔離作業
- 規劃與文件需優先考慮可被多個 worktree 並行使用
- frontend 依賴管理需相容於 worktree 間透過 pnpm 共用 `node_modules`

## Pre-implementation Review Gate

- 在任何實作開始前，當前規劃 / 規範文件必須先 commit 並 push，作為正式基線
- 每個 subagent 在開始實作前，都必須先整理「此次要做的事情」成文件
- 該文件需提交到 `docs\plans\v1-mvp\subagent-work-items\pending\`
- 文件經 review 核可前，不得開始實作
- 只有在使用者明確指示後，文件才可從 `docs\plans\v1-mvp\subagent-work-items\pending\` 移到 `docs\plans\v1-mvp\subagent-work-items\approved\`
- 若 scope 變更，需更新文件並重新 review

## Subagent 交接要求

- 說明本次改動影響哪些路徑
- 說明是否新增依賴、命令或限制
- 若偏離規範，需明確列出偏離點與原因
