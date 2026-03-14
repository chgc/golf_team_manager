# Backend Go Conventions

## 強制規則

- 後端開發 **必須遵循 Google Go style guide**
- 每次編輯 Go 程式後 **必須執行 `gofmt`**
- backend 變更 **必須包含測試**
- 受影響函式必須可被測試驗證；不得把核心邏輯寫成難以隔離測試的形式
- backend framework **必須使用 Gin**
- **不得使用 ORM library**

參考來源：

- Google Go Style Guide: `https://google.github.io/styleguide/go/guide`
- Google Go Best Practices: `https://google.github.io/styleguide/go/best-practices`

## 工具與結構

- backend 使用 Go 實作
- HTTP router 使用 `gin`
- package layout 以 `cmd` + `internal` 為主
- migrations 與資料檔案目錄需明確獨立

建議結構：

```text
backend\
├── cmd\
├── internal\
│   ├── app\
│   ├── config\
│   ├── db\
│   ├── http\
│   ├── middleware\
│   └── services\
└── migrations\
```

## 分層原則

- `cmd` 負責啟動
- `http` 負責 Gin handlers、request/response mapping
- `services` 負責業務邏輯
- `db` / repository 負責資料存取
- 不讓 handler 直接操作 SQL 細節

## 契約與型別

- domain model、request DTO、response DTO 分開定義
- 不直接把 database row struct 當 response model
- 明確處理 nullable 欄位，不靠隱含零值混過去

## Config 與執行環境

- config 來源應集中管理
- 本地開發所需的 DB 路徑、port、mode 等設定需可調整
- 避免散落在多個 package 直接讀 environment variables

## 錯誤與 HTTP 行為

- API 錯誤格式遵循 shared conventions
- 使用明確 HTTP status codes
- 不吞錯，不用 broad success fallback 掩蓋錯誤
- health endpoint 與 migration 錯誤要清楚可觀測

## Database 與 Migration

- SQLite 連線生命週期集中管理
- migration 檔案需可重複執行
- Phase 1 先做 smoke migration，不急著建完整業務 schema
- 優先選擇 Windows 友善的 pure-Go SQLite 方案
- 資料存取使用原生 SQL / query layer，不引入 ORM

## 測試要求

- 至少具備 `go test ./...` 的 baseline
- Phase 1 至少覆蓋啟動、health、migration smoke path
- 新增或修改的函式需設計成可測試，並以對應測試驗證行為
- 不能把重要邏輯藏在難以 mock、難以注入或難以驗證的流程中

## Subagent 交付要求

- 說明新增的 package 與責任邊界
- 說明啟動命令、`gofmt` 執行方式與驗收方式
- 說明測試覆蓋哪些受影響行為
- 若新增第三方套件，需說明理由與替代方案取捨
