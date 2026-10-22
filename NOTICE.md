# 来源与许可说明

[简体中文](#来源与许可说明) · [English](#english)

本仓库 `t508708/icloud-mail` 是 [mangobubu/icloud-api](https://github.com/mangobubu/icloud-api) 的派生版本，保留上游来源与 Git 历史。在其基础上增加和调整了邮箱池、按需取件、Apple Account 连接管理、批量操作、Web 界面与部署工具。

## 许可现状

本仓库选择 [MIT License](LICENSE) 作为维护者有权许可内容的发布协议，`LICENSE` 保留标准 MIT 原文。该选择无需向 MIT 或其他机构申请、审批或付费。

该许可仅在授权者实际拥有相应版权或明确授权的范围内生效。本说明不主张维护者拥有上游代码或其他贡献者的版权，也不替代这些权利人的授权。此前对维护者拥有“原创增量”的笼统表述未经逐项权属核验，本次予以更正。

上游已有代码、其他权利人的贡献及第三方依赖保留各自权利和许可状态。根目录存在 MIT 文本或平台识别出 MIT 标签，均不代表整个派生代码库已取得统一 MIT 授权。

2026-10-22 重新核查上游 `main` 提交 [`866a0f059200dad544a2f759f6a3acde31a24c1f`](https://github.com/mangobubu/icloud-api/tree/866a0f059200dad544a2f759f6a3acde31a24c1f)：GitHub 仓库元数据未识别许可证，完整文件树未发现 `LICENSE` / `LICENCE` / `COPYING` / `NOTICE` 文件。此前 v0.2.0 核查记录对应提交 [`2d5d570d8892c7de84a679d7d4e2bbc0360a759a`](https://github.com/mangobubu/icloud-api/tree/2d5d570d8892c7de84a679d7d4e2bbc0360a759a)。这些提交仅用于固定许可核查依据，不表示本项目已合入其全部功能。

上游目前未声明 MIT、Apache、GPL 等通用软件许可证，原作者与其他贡献者保留各自代码的权利。源码公开本身不代表整个代码库已授予商业使用或再分发许可；需要额外授权时应向相应权利人确认。后续取得明确授权时，再更新本说明和相应许可范围。

本项目现于独立 GitHub 仓库 [t508708/icloud-mail](https://github.com/t508708/icloud-mail) 维护；迁移前的 GitHub fork 留档于 [t508708/icloud-mail-fork-archive-20260924](https://github.com/t508708/icloud-mail-fork-archive-20260924)。仓库关系调整保留本项目 Git 历史、上游来源和上述许可说明，不改变各部分代码原有的权利归属。公开仓库只同步经过检查的项目源码；当前服务器的账号、邮件、凭据、证书、运行日志和本机运维历史均未作为发布内容。

## 参考项目

- [q1953258942/iCloud-Privacy-Mail](https://github.com/q1953258942/iCloud-Privacy-Mail)：Apple Account 协议流程参考。
- [cquxl/iCloud-Alias-Hub](https://github.com/cquxl/iCloud-Alias-Hub)：连接续期与管理交互参考；其 `NOTICE.md` 对自身授权另有说明。
- [pyicloud](https://github.com/picklepete/pyicloud)、[ImapFlow](https://github.com/postalsys/imapflow)：会话持久化与 IMAP 同步工程参考。

参考依据和固定版本见 `docs/APPLE-ACCOUNT-RENEWAL.md`、`docs/APPLE-ACCOUNT-RISK-REVIEW.md`、`docs/MAIL-RECEPTION.md`。引用某项目并不表示该项目的许可证覆盖本仓库。

Go 与 npm 依赖按各自许可证分发，版本记录于 `go.mod`、`go.sum` 和 `web/package-lock.json`。本次改动后的镜像构建在 `/usr/share/licenses/icloud-mail/` 同时包含本项目 `LICENSE` 与本说明；应用镜像另保留 Go 标准库许可证和构建依赖中的 LICENSE / NOTICE 等文件及版本索引。基础镜像组件沿用其发行版的授权与来源，本仓库未重新许可这些依赖。本项目为非官方自托管工具，与 Apple Inc. 无隶属或背书关系。

## English

### Origin and license scope

`t508708/icloud-mail` is derived from [mangobubu/icloud-api](https://github.com/mangobubu/icloud-api), with upstream attribution and Git history preserved. It adds and adapts mailbox pools, on-demand retrieval, Apple Account connection management, bulk operations, the Web interface, and deployment tooling.

This repository selects the [MIT License](LICENSE) for material the maintainer has the right to license. `LICENSE` contains the standard MIT text. Adopting MIT requires no application, approval, or payment to MIT or any other institution.

The grant is effective only to the extent that the grantor holds the relevant copyright or explicit permission to license it. This notice does not assert that the maintainer owns upstream code or other contributors' copyright, and does not replace permission from those rights holders. The earlier blanket description of maintainer-owned “original contributions” had not been verified contribution by contribution and is corrected here.

Pre-existing upstream code, other parties' contributions, and third-party dependencies retain their respective rights and licensing status. The presence of an MIT file, or a hosting platform's MIT label, does not establish MIT permission for the entire combined codebase.

As of the 2026-10-22 review, no general-purpose software license was identified for upstream commit [`866a0f059200dad544a2f759f6a3acde31a24c1f`](https://github.com/mangobubu/icloud-api/tree/866a0f059200dad544a2f759f6a3acde31a24c1f): GitHub reported no license, and the complete file tree contained no `LICENSE`, `LICENCE`, `COPYING`, or `NOTICE` file. The earlier v0.2.0 review covered [`2d5d570d8892c7de84a679d7d4e2bbc0360a759a`](https://github.com/mangobubu/icloud-api/tree/2d5d570d8892c7de84a679d7d4e2bbc0360a759a). These references document the review, not incorporation of every upstream change.

Upstream authors and other contributors retain their rights. Public source access alone does not grant commercial-use or redistribution rights for the combined codebase. Additional permissions must come from the relevant rights holders. This notice and the applicable license scope will be updated if explicit permission is obtained.

The project is independently maintained at [t508708/icloud-mail](https://github.com/t508708/icloud-mail). Its former GitHub fork is archived at [t508708/icloud-mail-fork-archive-20260924](https://github.com/t508708/icloud-mail-fork-archive-20260924). Changing the repository relationship does not change ownership or licensing. Published source excludes deployment accounts, mail, credentials, certificates, runtime logs, and private operational history.

### References and third-party components

- [q1953258942/iCloud-Privacy-Mail](https://github.com/q1953258942/iCloud-Privacy-Mail): Apple Account protocol flow reference.
- [cquxl/iCloud-Alias-Hub](https://github.com/cquxl/iCloud-Alias-Hub): session renewal and management interaction reference; its own `NOTICE.md` describes its licensing.
- [pyicloud](https://github.com/picklepete/pyicloud) and [ImapFlow](https://github.com/postalsys/imapflow): session persistence and IMAP synchronization references.

Reference details and pinned versions are recorded in `docs/APPLE-ACCOUNT-RENEWAL.md`, `docs/APPLE-ACCOUNT-RISK-REVIEW.md`, and `docs/MAIL-RECEPTION.md`. Citing a project does not extend its license to this repository.

Go and npm dependencies retain their respective licenses; versions are recorded in `go.mod`, `go.sum`, and `web/package-lock.json`. Builds using the current Dockerfile include this notice and `LICENSE` under `/usr/share/licenses/icloud-mail/`. The application image additionally includes the Go standard library license, dependency license and notice files, and version indexes. Base-image components retain their distribution licenses. This repository does not relicense those components.

This is an unofficial self-hosted tool, with no affiliation with or endorsement by Apple Inc.
