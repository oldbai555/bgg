# admin-frontend 重构过程记录（Phase 1-3）

> 本文档是本轮 `admin-frontend` 重构 Phase 1-3 期间的唯一过程记录（阶段/周次/关键决策），只追加不重写。与仓库根目录 `docs/前端开发进度.md`（跨越整个项目生命周期的功能级进度索引）是两份不同的记录，不要互相替代——分工原则见 `00-refactor-overview.md` 第 4 节。

---

## 2026-07-13：文档规划阶段完成

**What**：产出 `admin-frontend/docs/00~09.md` 全部 10 篇任务书文档，风格对齐 `admin-server/docs/00-refactor-overview.md` 的先例。`admin-frontend/src/` 下没有任何代码改动，本次交付是文档。

**核心决策**（详见各文档，此处只记摘要，不复制正文）：
- 域目录重组对齐后端 9 域（`iam/system/monitoring/misc/blog+video→content/chat/sdk/task`），而不是 5 个部署服务粒度。
- API 层：8 个域各建一个手写 wrapper，视图禁止直接 import `api/generated/` 里的请求函数。
- 引入 vitest，核心逻辑（stores/composables/request 拦截器/纯函数 utils）补测试，不设覆盖率门槛。
- 后台管理界面视觉：设计令牌化 + 精修（稳健路线，不做大幅重塑）。
- 公共页面视觉：完全重构为响应式优先的企业级方案，废弃"小程序风格"单一 768px 断点契约。
- 暗色模式：本轮做成全面支持（后台 + 公共页均需适配）。
- 规则/文档同步纳入每个 Phase 的收尾动作，不是收尾才做（见 `09-rules-and-docs-sync-checklist.md`），已额外核实脚手架模板 `list_page.vue.tpl` 存在与新 API 规范的直接冲突（第 47 行硬编码 import generated），需要在 Phase 1 尽早修模板，避免后续新生成代码持续违规。

**Why**：用户要求 admin-frontend 参照已完成三阶段重构的 admin-server，进行"架构 + 视觉/UX 一起重做"的大规模重构；项目未上线，无兼容性负担，用户明确表态"时间足够，可以放心大胆重构，要改得彻底"。

**下一步**：从下次会话开始按 `00-refactor-overview.md` 第 3 节的 Phase 1 Week 1 顺序实际动代码（域目录重组 + API wrapper 全覆盖），完成后回来本文档追加条目。
