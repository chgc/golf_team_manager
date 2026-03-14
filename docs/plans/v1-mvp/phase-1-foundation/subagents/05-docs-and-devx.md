# Subagent Brief — Docs and Dev Experience

## 任務目標

補齊 Phase 1 所需文件與開發流程說明，讓後續 subagent 或人類開發者可快速接手。

## 範圍內

- 本地開發啟動說明
- 目錄用途說明
- 技術決策與 Phase 1 邊界說明
- 前後端與 DB 的驗收命令整理

## 範圍外

- 不重寫產品規格
- 不進入 v2+ 規劃

## 建議輸出

- `docs\development\local-setup.md`
- `docs\architecture\repo-structure.md`
- 視需要更新 root `README.md`

## 建議步驟

1. 盤點 workspace、frontend、backend、db 實際輸出
2. 寫成本地啟動順序
3. 記錄常用命令與驗證方式
4. 記錄 Phase 1 不做的內容，避免後續 scope 漂移

## 驗收標準

- 新接手的 agent 能知道從哪裡開始
- 文件與實際檔案結構一致
- 有明確下一步指向 `shared-domain-schema`

## 交接備註

- 此工作最好在其他 Phase 1 workstream 初步完成後再收尾一次
