# Subagent Brief — Backend Bootstrap

## 任務目標

建立 Go 後端骨架，使服務具備可啟動、可擴充、可承接資料庫與後續 REST API 的基礎能力。

此任務必須先閱讀 `..\..\phase-0-conventions\backend-go-conventions.md`。

## 範圍內

- 初始化 Go module
- 建立 `cmd` / `internal` 結構
- 建立 Gin router、config、app startup
- 提供 `/health` 或等效 health endpoint
- 建立統一錯誤回應與 middleware 放置位置

## 範圍外

- 不實作完整 player / session / registration CRUD
- 不正式串接 OAuth
- 不做 WebSocket
- 不引入 ORM library

## 建議檔案方向

- `backend\cmd\api\main.go`
- `backend\internal\http\`
- `backend\internal\config\`
- `backend\internal\app\`
- `backend\internal\middleware\`

## 建議步驟

1. 決定 module path
2. 建立 server 啟動入口與設定載入
3. 加入 Gin router 與基本 middleware
4. 建立 health endpoint
5. 預留 DB、auth、handlers 的 package 位置
6. 驗證 `go test ./...` 與服務可啟動

## 設計限制

- 以後續 REST API 與 middleware 擴充為前提
- 保持 `main` 精簡，將組裝放在 `internal`
- 避免 Phase 1 就把業務與 transport 緊耦合
- 遵循 Google Go style guide
- 每次編輯後執行 `gofmt`
- baseline 交付必須包含可通過的測試
- 使用 Gin 作為 HTTP framework
- 不使用 ORM，保留 raw SQL / repository 路徑

## 驗收標準

- 後端可啟動
- health endpoint 可回應
- 專案結構清楚，後續 handler / service / repository 可自然擴充

## 交接備註

- DB agent 之後會接手 SQLite 與 migration，請保留乾淨的注入點
