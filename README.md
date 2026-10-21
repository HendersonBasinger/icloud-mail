# iCloud Mail

自托管的 **iCloud 隐藏邮箱管理、邮箱池与 API 取码服务**，带本地化 Web 控制面板。

**当前版本：v0.2.7** · [下载与发布说明](https://github.com/t508708/icloud-mail/releases/tag/v0.2.7) · [更新日志](CHANGELOG.md) · [安装文档](DELIVERY.md)

## 核心功能

| 功能 | 说明 |
| --- | --- |
| Web 控制面板 | 苹果风格界面、深浅色主题；JS/CSS 本地加载 |
| 隐藏邮箱管理 | 多主号、目录同步、手动批量创建、自动计划、批量删除 |
| 邮箱池 | 项目 Key、自动领取、确认使用、释放、续期与凭据轮换 |
| API 取码 | 单邮箱链接、Bearer API、Pool lease 取码；严格支持 `+tag` 归入已登记根邮箱 |
| 邮件接收 | 按需增量收件、同主号连接复用、MIME 归档、只读 IMAPS |
| 独立凭据 | 每邮箱 API Key、IMAP 密码及 OAuth 凭据，由本网关签发 |

## 优势与已知限制

**适合的场景**：集中管理多个 iCloud 主号、持续补充隐藏邮箱库存，并通过邮箱池或 API 对接取码。控制面板资源本地化；取件按需触发，同主号的多个隐藏邮箱共享同步与连接，减少重复收件。支持已登记根邮箱的 `+tag` 邮件归档，适配使用根邮箱取码的协作项目。

### 自动创建：每个主号每小时最多 23 次

| 通道状态 | 本地自动创建上限（滚动一小时） |
| --- | --- |
| 新、旧通道均正常 | 新通道 19 次＋旧通道 4 次，合计 23 次 |
| 新通道过期，旧通道正常 | 旧通道 5 次；重新登录新通道后自动恢复双通道规则 |
| 旧通道目录认证失效 | 暂停创建，恢复登录后继续；新通道仍需原目录校验 |

**为什么不拉满 20＋5？** 按已有使用观察，新通道常见约 20 次/小时、旧通道约 5 次/小时触限，这些数值并非 Apple 官方配额承诺。双通道自动计划各留一次余量，并将请求分散到时间槽，降低滚动窗口重叠、其他客户端创建等因素导致触顶后整段等待的概率。留余量也不保证免限流；Apple 明确限流后按对应通道冷却，遵守更长的 `Retry-After`。

上述是**创建尝试上限，非保证成功数量**，失败也可能消耗尝试额度；实际新增量受认证、目录确认、主号容量和 Apple 响应影响。手动创建与探测绕过本地自动额度和冷却，但仍受 Apple 上游限制，频繁手动操作也可能影响后续自动创建。详见 [自动创建与通道恢复](docs/AUTO_CREATION_CHANNELS.md)。

### 新通道缺陷：约 8 小时后仍可能需要重新登录

已多次观察到：**新通道即使持续保活成功，仍可能在初次登录约 8 小时后失效**。目前每个主号按 4～6 分钟随机周期独立保活，能刷新短期会话，但尚未证明能突破这一长时失效问题；请勿按全天免登录服务预期使用。这是项目已知缺陷与实测现象，并非已证实的 Apple 官方固定 8 小时上限。

失效后需要管理员重新登录新通道；只要旧通道和目录校验仍正常，已开启的自动计划会降级为每小时最多 5 次，重新登录后恢复 23 次规则。管理连接与邮件收件凭据分别维护，取件还取决于对应收件连接是否正常。长期无人值守需考虑重新登录维护，后台成功续期不等于长期登录保证。详见 [保活记录与限制](docs/APPLE-ACCOUNT-RENEWAL.md)。

**其他限制**：项目依赖非官方 Apple 接口，上游变更、地区差异及账户状态可能影响使用；按需收件的实际耗时受上游网络与邮件投递影响，本地离线测试耗时不代表线上承诺。网关签发的独立 IMAP/OAuth 凭据仅用于本服务，直连 Apple 收件仍需相应主号的 App 专用密码。

## v0.2.7 更新重点

- **登录自动继续**：遇到 Apple 允许跳过的手机号补充提示时，新旧通道自动继续，无需专门打开官网选择“不添加电话号码并继续”。
- **修复登录报错**：解决部分账号验证码已通过，却仍提示“需要完成账户操作”或“服务暂时不可用”的问题。
- **保留设备信任**：登录后继续信任设备；后续连接失败时保留已取得的信任结果，供下次登录复用。是否再次要求验证码仍由 Apple 决定。

本次为登录兼容修复，**尚未解决新通道约 8 小时后可能失效的问题**。强制绑定手机号、条款更新等其他账户待办仍需按 Apple 页面提示处理。

从初版 `v0.1.0` 起的完整新增、调整、修复与验证记录见 [更新日志](CHANGELOG.md)。

## QQ 群

项目交流与部署问题可加入 QQ 群 **1105888476**：[一键加群](https://qm.qq.com/q/6KWBNKKDEk)。

## Docker 部署

需要 Docker Engine、Docker Compose v2 和 Git。GHCR 镜像支持 **linux/amd64** 与 **linux/arm64**，建议至少 2 核、4 GB 内存。

**1. 获取部署配置**

```sh
git clone --depth 1 https://github.com/t508708/icloud-mail.git
cd icloud-mail
cp -n .env.example .env
```

**2. 拉取 `latest` 并启动**

```sh
docker compose pull
docker compose up -d --wait
docker compose ps
```

默认使用 `ghcr.io/t508708/icloud-mail:latest` 与 `ghcr.io/t508708/icloud-mail-postgres:latest`。每次更新执行 `git pull --ff-only && docker compose pull && docker compose up -d --wait`，已有数据库、密钥和邮件归档卷会继续复用。

**3. 登录后台**

默认地址：**http://127.0.0.1:8788/admin/**，用户名 `admin`。查看首次生成的密码：

```sh
docker compose exec -T icloud-api cat /app/keys/admin-password
```

远程服务器可在自己的电脑运行以下命令，再打开上面的本地地址：

```sh
ssh -N -L 8788:127.0.0.1:8788 USER@SERVER
```

公网域名、HTTPS 与宝塔配置见 [公网部署](docs/PUBLIC.md)。镜像已包含前端资源，运行时仍需连接 Apple 或配置的上游 IMAP。

### 从源码构建

需要本机构建时叠加 `compose.build.yaml`：

```sh
docker compose -f compose.yaml -f compose.build.yaml up -d --build --wait
```

需要固定发布或离线安装时，使用 GitHub Release 中带校验文件的镜像包，步骤见 [安装文档](DELIVERY.md#离线与固定版本)。两种方式使用同一控制面板。登录后添加主号、配置收件、连接 Apple 并同步地址，再复制取码链接或将邮箱加入池。直连 iCloud 收件使用 **App 专用密码**；Apple 登录用于地址管理。

升级前按 [安装文档](DELIVERY.md) 备份数据库、密钥与邮件归档，保留 `.env` 和数据卷。

## 文档

[安装与升级](DELIVERY.md) · [公网部署](docs/PUBLIC.md) · [邮箱池与 API](docs/POOL.md) · [OpenAPI](docs/openapi.yaml)

[配置与兼容说明](docs/REFERENCE.md) · [收件设计](docs/MAIL-RECEPTION.md) · [测试与实测记录](docs/RECEIVER-VALIDATION.md)

[自动创建与通道恢复](docs/AUTO_CREATION_CHANNELS.md)

[Apple 登录的手机号补充与设备信任](docs/APPLE-LOGIN-REPAIR.md)

每个部署的 `/docs/` 提供随服务版本更新的公开接入说明；对应源码见 [内置接口文档](internal/httpserver/docs.html)。

## 来源与许可

本项目由 [t508708/icloud-mail](https://github.com/t508708/icloud-mail) 维护，基于 [mangobubu/icloud-api](https://github.com/mangobubu/icloud-api) 开发，保留来源与 Git 历史。

上游目前未提供通用软件许可证，本仓库未另行声明整个代码库的 MIT/GPL 授权。原作者、贡献者与依赖各自的权利保留，详见 [NOTICE.md](NOTICE.md)。本项目为非官方工具，与 Apple 无隶属关系。
