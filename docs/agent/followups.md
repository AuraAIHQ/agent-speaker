# Follow-ups ledger（append-only · 永不删行 · 提交进仓库）

> pilot 的 review triage 把「真问题但不阻塞（B）」和延后项记在这里。
> 主线 task 全部完成后，由 `pilot run` 批量合成一个 cleanup PR 做掉，逐条标 [x] done=PR#n。
> `- [ ]`=OPEN，`- [x]`=DONE。GitHub PR comment 是永久兜底。

- [ ] FU-1 · B · src=PR#33 · 2026-09-03 · install.sh 的 curl|bash 安装路径拉 releases/latest/download/hyphae-${PLATFORM}.tar.gz,但仓库没有任何 release 自动化上传过二进制资产(v0.26.0 release 是空的,改名前对 agent-speaker 也是同样 404)。需要补一条 GitHub Actions release workflow 打包上传二进制,再验证 install.sh 真的能跑通
- [ ] FU-2 · B · src=PR#33 · 2026-09-03 · internal/messaging/outbox.go 的 SaveOutbox 仍用固定临时文件名 outbox.json.tmp,两个并发写者会 O_TRUNC 到同一 inode、按各自 offset 写入,os.Rename 发布的是已损坏的 inode(PR#33 并发轮实测 3000 次并发对中 11% 产出不可解析文件)。它比 keystore 暴露得多:daemon 每 60s 重试 tick + 每次 agent msg 都写,且全非交互。改用 os.CreateTemp;同时考虑 rename 前留一份 .bak
- [ ] FU-3 · B · src=PR#33-peer-handoff · 2026-09-10 · internal/messaging/outbox.go 的重试路径 StoreOutgoingMessage(&event, ..., event.Content, true) 把密文当明文存进 plaintext 列——event 是从 entry.EventJSON 反解出来的已加密(可能还 zstd 压缩过)事件。单独发生时是脏数据;若 daemon 已先存过该 id 的解密明文,重试会用密文覆盖掉它(新的 UPSERT 只挡空值,挡不住非空的错值,见 TestPlaintextGuardDoesNotCoverNonEmptyOverwrite)。修法在源头:重试路径要么解密后再存,要么存空 plaintext 让 UPSERT 的保留分支生效
- [ ] FU-4 · B · src=PR#33-peer-handoff · 2026-09-10 · daemon 的 watchOneRelay 对 RelayConnect 失败是裸 return 0,relay 连不上和确实没有新消息在 stdout 里都打印 'Watching... (no new messages)',下游无法区分死连接和安静;且全进程没有 status 子命令/心跳文件,拿不到 last-seen。Agent24 因此只能在桥侧自造活性信号,而活性判定本该由持有连接的一方给出。建议:连接失败要有区别于'无新消息'的输出,并考虑加一个 status/心跳出口
