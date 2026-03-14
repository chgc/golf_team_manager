# Subagent Brief — Workspace Scaffold

## 任務目標

建立 repo 的頂層骨架與共用慣例，讓後續 frontend、backend、docs 工作流都能在一致的基線上進行。

## 範圍內

- 建立 top-level 資料夾
- 補齊 `.gitignore`
- 建立最小 root README 或導覽文件
- 定義命名與結構慣例
- 預留 `docs` 下的架構 / 開發文件位置
- 規劃可相容於 `git worktree` 模式的協作基線

## 範圍外

- 不建立完整業務功能
- 不決定 API schema 細節
- 不處理 LINE OAuth

## 建議輸出

- `frontend\`
- `backend\`
- `docs\architecture\`
- `docs\development\`
- root `.gitignore`
- root `README.md`
- 若需要，補充 pnpm / worktree 的 root 級設定說明

## 建議步驟

1. 檢查 repo 現況與既有檔案
2. 建立頂層資料夾與基礎導覽
3. 補上 Windows / Node / Go / SQLite 相關 ignore 規則
4. 在 README 說明 repo 角色分工與開發入口
5. 確認結構不會阻礙 Angular CLI、pnpm 與 Go module 後續初始化
6. 確認結構適合 subagent 透過 git worktree 並行作業

## 驗收標準

- 後續 subagent 不需重排 top-level 結構
- ignore 規則足以避免常見暫存檔進版控
- root README 能快速說明 repo 佈局
- 基線結構能支援 git worktree + pnpm 共用依賴的工作流

## 交接備註

- 將結構決策寫清楚，避免 frontend / backend subagent 各自建立不同慣例
