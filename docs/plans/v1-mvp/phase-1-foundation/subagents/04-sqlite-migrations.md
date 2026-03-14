# Subagent Brief — SQLite and Migrations Bootstrap

## 任務目標

建立 SQLite 初始化與 migration 機制，讓後續 schema 工作能在穩定基線上演進。

## 範圍內

- 選定 SQLite driver
- 建立 DB 連線封裝
- 建立 migration 檔案目錄
- 建立 migration runner 或等效初始化流程
- 建立 smoke migration 驗證整體流程

## 範圍外

- 不建立完整業務 schema
- 不加入 seed data 的完整內容
- 不處理正式部署策略

## 建議檔案方向

- `backend\migrations\`
- `backend\internal\db\`
- `backend\internal\db\migrate.go`
- `backend\data\` 或等效本機資料目錄

## 建議步驟

1. 選擇 Windows 友善的 pure-Go SQLite driver
2. 建立 DB 開啟、關閉、錯誤處理流程
3. 定義 migration 檔命名規則
4. 建立第一個 smoke migration
5. 讓 backend 啟動時可選擇執行 migration 或提供獨立命令
6. 驗證空白環境可成功初始化 DB

## 設計限制

- 以簡單可維護為優先
- migration 需可重複執行且不破壞既有資料
- 不在 Phase 1 過度設計 ORM 層

## 驗收標準

- 可從零建立 SQLite 檔與基本 schema
- migration 流程清楚且可重複操作
- 後續 `shared-domain-schema` 可直接新增正式 schema

## 交接備註

- 與 backend agent 對齊 config / path / lifecycle，避免雙方各自管理 DB 初始化
