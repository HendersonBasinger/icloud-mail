# 来源与许可说明

本仓库 `t508708/icloud-mail` 是 [mangobubu/icloud-api](https://github.com/mangobubu/icloud-api) 的派生版本，保留上游来源与 Git 历史。在其基础上增加和调整了邮箱池、按需取件、Apple Account 连接管理、批量操作、Web 界面与部署工具。

## 许可现状

本项目维护者拥有权利且有权授权的原创代码和改动采用 [MIT License](LICENSE)，版权人为 `t508708`。MIT 的适用范围仅为上述原创增量，不重新许可上游已有代码、其他权利人拥有的贡献或第三方依赖，也不授予整个派生代码库统一的 MIT 授权。

2026-10-22 重新核查上游 `main` 提交 [`866a0f059200dad544a2f759f6a3acde31a24c1f`](https://github.com/mangobubu/icloud-api/tree/866a0f059200dad544a2f759f6a3acde31a24c1f)：GitHub 仓库元数据未识别许可证，完整文件树未发现 `LICENSE` / `LICENCE` / `COPYING` / `NOTICE` 文件。此前 v0.2.0 核查记录对应提交 [`2d5d570d8892c7de84a679d7d4e2bbc0360a759a`](https://github.com/mangobubu/icloud-api/tree/2d5d570d8892c7de84a679d7d4e2bbc0360a759a)。这些提交仅用于固定许可核查依据，不表示本项目已合入其全部功能。

上游目前未声明 MIT、Apache、GPL 等通用软件许可证，原作者与其他贡献者保留各自代码的权利。源码公开本身不代表整个代码库已授予商业使用或再分发许可；需要额外授权时应向相应权利人确认。后续取得明确授权时，再更新本说明和相应许可范围。

本项目现于独立 GitHub 仓库 [t508708/icloud-mail](https://github.com/t508708/icloud-mail) 维护；迁移前的 GitHub fork 留档于 [t508708/icloud-mail-fork-archive-20260924](https://github.com/t508708/icloud-mail-fork-archive-20260924)。仓库关系调整保留本项目 Git 历史、上游来源和上述许可说明，不改变各部分代码原有的权利归属。公开仓库只同步经过检查的项目源码；当前服务器的账号、邮件、凭据、证书、运行日志和本机运维历史均未作为发布内容。

## 参考项目

- [q1953258942/iCloud-Privacy-Mail](https://github.com/q1953258942/iCloud-Privacy-Mail)：Apple Account 协议流程参考。
- [cquxl/iCloud-Alias-Hub](https://github.com/cquxl/iCloud-Alias-Hub)：连接续期与管理交互参考；其 `NOTICE.md` 对自身授权另有说明。
- [pyicloud](https://github.com/picklepete/pyicloud)、[ImapFlow](https://github.com/postalsys/imapflow)：会话持久化与 IMAP 同步工程参考。

参考依据和固定版本见 `docs/APPLE-ACCOUNT-RENEWAL.md`、`docs/APPLE-ACCOUNT-RISK-REVIEW.md`、`docs/MAIL-RECEPTION.md`。引用某项目并不表示该项目的许可证覆盖本仓库。

Go 与 npm 依赖按各自许可证分发，版本记录于 `go.mod`、`go.sum` 和 `web/package-lock.json`。本次改动后的镜像构建在 `/usr/share/licenses/icloud-mail/` 同时包含本项目 `LICENSE` 与本说明；应用镜像另保留 Go 标准库许可证和构建依赖中的 LICENSE / NOTICE 等文件及版本索引。基础镜像组件沿用其发行版的授权与来源，本仓库未重新许可这些依赖。本项目为非官方自托管工具，与 Apple Inc. 无隶属或背书关系。
