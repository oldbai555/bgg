# bgg 开发交接记录（changelog）

本目录是 `admin-server` + `admin-frontend` **共用**的开发交接记录，覆盖每一次实质性工作会话——起点是 admin-server 单体加固 → 微服务拆分 → 可观测性/CD 三阶段重构（Phase 1-3），但早已不止于此：不管改的是后端还是前端，日常修 bug、做新需求、方案调整都要往里追加条目，不要因为"不是重构任务"或"这次只改前端"就跳过。按时间顺序追加条目，记录做了什么、关键决策、关键文件位置（早年风格对齐的 `docs/后端开发进度.md`/`docs/前端开发进度.md` 已于 2026-07-17 退场归档，见下方 `archive-backend.md`/`archive-frontend.md` 说明）。

**本目录 2026-07-17 前叫 `admin-server/docs/changelog/`，只记 admin-server**：此前的条目（2026-07-10 ~ 2026-07-16）全部是 admin-server 专属内容，历史事实不改写。2026-07-17 起迁到仓库根目录，前后端共用同一份交接记录，同一天多个模块用文件名 `-2`/`-3` 后缀区分，不需要按前后端分文件夹。

原来单文件 `admin-server/docs/progress.md`（20 条日期条目，200KB+）已按日期拆分成本目录下的独立文件，方便按需查阅、避免单文件无限增长。`admin-server/docs/progress.md` 仍然保留（`22-admin-mcp-tool.md` 的 `query_progress` 工具硬编码解析这个文件名，且有几十处历史文档交叉引用它），但只保留最近少量条目，新条目积累到一定数量后批量归档到本目录；`admin-frontend/docs/progress.md` 是前端独立维护的长文档，与本目录不冲突、不合并。

**`archive-backend.md` / `archive-frontend.md` 是一次性批量归档文件，不是常规的按日期条目**：2026-07-17 按「文档分层与生命周期」规则（见 `.cursor/rules/00-workflow.mdc`）退场的 `docs/后端开发进度.md`/`docs/前端开发进度.md`，把 2026-07-10（本目录起点）之前、changelog 未覆盖的历史内容整篇原样搬了过来；2026-07-10 之后与 changelog 重复的部分已删除未保留。这两个文件同样不再维护，只是历史存档，不要往里追加新内容。

**维护方式：只追加，不重写。** 每次实际工作会话结束后，在本目录新增一个日期文件（同一天多条用 `-2`/`-3` 后缀区分），不要回头改写已有条目（除非是修正明确的事实错误）。历史决策即使后来被推翻，也保留原条目 + 新增一条说明推翻原因，不做静默删除。

**本目录同时是写给下一个没有上下文的新会话的交接文档。** 会话开始处理 admin-server 或 admin-frontend 任务前（不限于重构——修 bug、做需求同样适用），先读最新一篇日期文件顶部的「交接摘要」小节，了解当前进度、下一步计划和已知坑，不要凭空重新摸索一遍。

**2026-07-16 起新条目格式**：改用「上线记录」模板，模板文件见 [`TEMPLATE.md`](TEMPLATE.md)（结构：**交接摘要（必填，见下方说明）** + 上线需求/技术优化/线上缺陷三张表 + 逐个版本的「变更内容/上线前准备/依赖服务/服务分布情况/故障回滚」，团队沿用已有的上线记录习惯），实际用例见 [`2026-07-16.md`](2026-07-16.md)。**新增一条 changelog 前先复制 `TEMPLATE.md` 再填内容**，不要从零现编结构。2026-07-16 之前归档的条目保持原有自由格式（叙事体，不做格式转换），不强行套用新模板。

**「交接摘要」小节强制、不可省略**：无论这次会话是否涉及实际上线部署，只要对 admin-server 或 admin-frontend 有实质性改动（新增需求、代码重构、bug 修复、方案调整都算），条目顶部都必须填写「在做什么/已完成/当前状态/下一步计划/踩过的坑」五项。「上线需求/技术优化/线上缺陷」等表格化章节只在涉及实际上线时认真填，纯代码会话可以留空或写"无"，但交接摘要不能空着。

**隐私规范：不写服务器真实公网 IP**。本目录进 git 历史后很难彻底清除，本仓库又是公开仓库，记录部署/服务器相关内容时一律用服务器别名（如 `bgg-dev`，对应本地 `~/.ssh/config` 的 `Host` 别名），真实 IP 不写进任何 changelog 条目；域名（如 `oldbai.top`）不算敏感信息，可以正常写。这条同步写进了 `AGENTS.md`/`.cursor/rules/00-workflow.mdc`/`.claude/rules/00-workflow.md` 的「绝对禁止事项」表。

## 索引

| 文件 | 内容 |
|------|------|
| [archive-backend.md](archive-backend.md) | **批量归档**：admin-server 2025-01 ~ 2026-07-07（本目录起点之前）的历史开发记录，原 `docs/后端开发进度.md` |
| [archive-frontend.md](archive-frontend.md) | **批量归档**：admin-frontend 2025-12 ~ 2026-07-07（本目录起点之前）的历史开发记录，原 `docs/前端开发进度.md` |

以下 2026-07-10 ~ 2026-07-16 的条目均为迁移前的 admin-server 专属记录：

| 文件 | 内容 |
|------|------|
| [2026-07-10.md](2026-07-10.md) | 文档集编写完成（Phase 1-3 尚未开始实际代码改动） |
| [2026-07-10-2.md](2026-07-10-2.md) | Phase 1 Week 1 地基工作全部落地 |
| [2026-07-10-3.md](2026-07-10-3.md) | 提交前 Cursor 自动代码审查（Gentleman Guardian Angel）发现问题修复 |
| [2026-07-11.md](2026-07-11.md) | Phase 1 Week 2（`04-domain-iam-chat.md` + `05-domain-task.md`）全部完成 |
| [2026-07-11-2.md](2026-07-11-2.md) | Phase 1 Week 3（`06-domain-blog-video-sdk.md` + `07-domain-monitoring-system-misc.md`）全部完成，含全仓库 Day 5 扫尾 |
| [2026-07-11-3.md](2026-07-11-3.md) | Phase 1 Week 4-5（`08-testing-strategy.md` + `09-ci-cd-and-deployability.md` 后半部分）全部完成 |
| [2026-07-11-4.md](2026-07-11-4.md) | Phase 1 收尾确认——人工冒烟 + 集成测试套件真实验证，发现并修复一个真实缓存 bug |
| [2026-07-11-5.md](2026-07-11-5.md) | Phase 2 启动——`15-service-boundaries.md` 第 1/2/3 项完成的定义全部落地（db/services/ 全量目录重组） |
| [2026-07-11-6.md](2026-07-11-6.md) | task-rpc 拆分 Step 1-2 完成——通用 RPC 脚手架 + `pkg/taskcallback` 契约落地并测试通过 |
| [2026-07-11-7.md](2026-07-11-7.md) | task-rpc 完整拆分收尾——计划 7 步全部完成，真实数据库端到端验证通过 |
| [2026-07-11-8.md](2026-07-11-8.md) | 提交前 Gentleman Guardian Angel 审查发现问题修复 |
| [2026-07-12.md](2026-07-12.md) | sdk-rpc 拆分——Phase 2 第二个服务，`18-service-extraction-runbook.md` checklist 全部完成 |
| [2026-07-12-2.md](2026-07-12-2.md) | 提交前 Cursor 自动代码审查（Gentleman Guardian Angel）发现问题修复 |
| [2026-07-12-3.md](2026-07-12-3.md) | chat-rpc 拆分——Phase 2 第三个服务，首次用到 WS↔gRPC 桥接和 Redis Streams |
| [2026-07-12-4.md](2026-07-12-4.md) | content-rpc 拆分——Phase 2 第四个服务，blog+video 域合并拆分 |
| [2026-07-13.md](2026-07-13.md) | iam-rpc 拆分——Phase 2 最后一个服务，五个服务拆分至此全部落地 |
| [2026-07-13-2.md](2026-07-13-2.md) | Phase 2 收尾遗留问题处理——`VideoUpdate` 字段绑定 bug 修复、多接口端到端验证补充 |
| [2026-07-13-3.md](2026-07-13-3.md) | Phase 3 第一篇——Telemetry + 结构化 JSON 日志接入，六个服务全部落地 |
| [2026-07-13-4.md](2026-07-13-4.md) | Phase 3 第二篇——API 文档生成（`generate-swagger.sh`）完成 |
| [2026-07-13-5.md](2026-07-13-5.md) | Phase 3 第三篇——六容器镜像 + docker-compose 前置工作完成 |
| [2026-07-16.md](2026-07-16.md) | **上线记录格式**：bgg-dev 首次真实部署——清理历史遗留服务、Docker 重装、六服务 docker-compose 混合拓扑上线、admin-frontend + nginx 重构、飞书登录联调 |

以下条目起，本目录迁到仓库根目录，admin-server + admin-frontend 共用：

| 文件 | 内容 |
|------|------|
| [2026-07-17.md](2026-07-17.md) | changelog 目录从 `admin-server/docs/changelog/` 迁到仓库根 `docs/changelog/`，正式确立为前后端共用的强制交接文档，`TEMPLATE.md`/规则文件同步更新 |
| [2026-07-17-2.md](2026-07-17-2.md) | bgg-dev 切到 ghcr.io 镜像部署（`docker-compose.dev-mixed.yml`/`deploy-dev.sh`）+ 新增 `script/branch_from_main.sh` 系列 Git 分支/PR 工作流脚本；踩坑记录：squash merge 断祖先导致 GitHub PR 比较页失效及 `rebase --onto` 修复方法 |
| [2026-07-17-3.md](2026-07-17-3.md) | 文档治理：建立「文档分层与生命周期」规则，退役 `docs/后端开发进度.md`/`docs/前端开发进度.md`/3 份 DDD-lite 草稿，历史内容归档进 `archive-backend.md`/`archive-frontend.md` |
