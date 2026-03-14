# Subagent Work Items and Review Gate

## 目的

在任何 subagent 開始實作前，先以文件方式整理工作內容，讓 review 有明確標的，避免直接開工造成 scope 漂移或與既有規範衝突。

## 流程

1. 主規劃 / 規範文件先完成並 commit + push
2. subagent 在開始實作前，先建立一份工作文件
3. 工作文件放在 `pending\`
4. 等待 review
5. 只有在**使用者明確指示**後，才可將文件移到 `approved\`
6. 文件移到 `approved\` 後，需先將這次核可結果 commit
7. 完成 commit 後，才可在 `git worktree` 環境下開始實作

## 目錄

```text
subagent-work-items\
├── README.md
├── templates\
│   └── subagent-task-template.md
├── pending\
└── approved\
```

## 命名建議

檔名使用：

`<phase>-<area>-<short-task-name>.md`

例如：

- `phase1-frontend-bootstrap-shell.md`
- `phase1-backend-bootstrap-gin-server.md`
- `phase2-domain-schema-player-session.md`

## 最低內容要求

每份工作文件至少包含：

- 目標與範圍
- 不在範圍內的內容
- 相依性（若有，需明確列出）
- 預計修改檔案 / 目錄
- 技術決策與依據規範
- 風險 / 待確認事項
- 驗收方式

## 注意事項

- 未經 review 核可，不得開始實作
- 除非使用者明確指示，不能自動將文件從 `pending\` 移到 `approved\`
- 文件進入 `approved\` 後，未完成 commit 前不得開始實作
- 若工作範圍變更，需更新文件並重新 review
- 文件內容應與 `Phase 0`、`Phase 1` 規範保持一致
