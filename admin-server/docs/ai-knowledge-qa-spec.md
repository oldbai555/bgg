# Spec: AI 知识库问答（bgg 项目 AI 能力接入 · Phase 1）

## Objective

给 bgg 项目补一个此前完全空白的业务类型：AI 能力接入。Phase 1 做「博客知识库问答」——用户提一个问题，系统对已发布的 blog 文章做向量检索，把最相关的文章片段喂给本地 LLM，生成基于真实文章内容的回答。Phase 2（本 spec 不覆盖，留作后续迭代）再把同一套检索+生成能力接入 chat 模块，做成群聊/私聊里的 AI 机器人。

验证目标：跑通"文本 → 向量 → 检索 → 生成"这条 RAG 最小闭环，且全程不依赖任何外部付费 API。

## 关键假设（需确认）

1. **本机 Redis 8.8.0 未加载 RediSearch 模块**（`FT._LIST` 返回 unknown command），只加载了 `vectorset` 模块（Redis 8 原生 Vector Sets 特性）。因此向量存储/检索必须走 `VADD`/`VSIM`/`VREM` 等 Vector Set 命令，**不是** MCP 工具名字暗示的 RediSearch 哈希索引（`FT.CREATE`/`FT.SEARCH`）那一套。go-redis v9.16 若无 Vector Set 的类型化 API，需用 `UniversalClient.Do(ctx, "VADD", ...)` 发原始命令。
2. **Ollama 已就绪**：`brew install ollama` 已装好并通过 `brew services start ollama` 常驻（`http://localhost:11434`），embedding 用 `bge-m3`、对话用 `qwen2.5:7b`（本机 Apple M4 + 16GB 内存，7B Q4 量化模型 + bge-m3 合计约 5-6GB，按需加载不会同时占满内存）。两个模型拉取中，完成后即可开始联调。
3. Phase 1 **不新建 MySQL 表**：知识库来源直接实时调用现有 `ContentRPC`（`PublicBlogArticleList`/`Detail` 或对应的已登录态列表接口）取文章，向量和文章片段都存在 Redis 里（Vector Set 的 `SETATTR` 可以给每个向量元素挂 JSON 属性，用来存文章片段原文，省掉一张表）。这意味着这次不走 `generate-sql.sh` 脚手架，符合"聊天/WebSocket 类偏离单表 CRUD 的模块手写"的项目约定——但**新增 API 端点仍然要走 `admin.api` 编辑 + 用户执行 `generate-api.sh` 生成 handler 骨架**，不能跳过（项目规则里 `generate-*.sh` 一律用户亲自执行，没有"因为是原型就免除"这一条）。
4. 交互方式确认为普通 HTTP 请求-响应（不做流式打字机效果）。

## Tech Stack

- 后端：go-zero（沿用 gateway 现有 `internal/handler`、`internal/logic` 分层），新增 group `ai/knowledge_qa`
- 向量存储：本机 Redis 8.8.0 Vector Sets（`VADD`/`VSIM`），走 go-redis v9.16 的 `Do()` 原始命令
- LLM/Embedding：本地 Ollama REST API（`/api/embeddings`、`/api/chat`），Go 侧用标准库 `net/http` 直接调，不引入额外 SDK
- 知识来源：现有 `ContentRPC`（content-rpc 的 blog 文章相关方法），不新建表、不碰 content 服务的 repository/model

## Project Structure（新增/改动）

```
admin-server/api/admin.api                                   → 新增 group: ai/knowledge_qa（Ask + Reindex 两个接口的类型+服务块）
admin-server/internal/handler/ai/knowledge_qa/                → goctl 生成（用户跑 generate-api.sh 后产出，不手写）
admin-server/internal/logic/ai/knowledge_qa/
  ├── knowledge_qa_ask_logic.go                               → 手写：embed 问题 → VSIM 检索 → 拼 prompt → 调 Ollama chat → 返回
  └── knowledge_qa_reindex_logic.go                           → 手写：拉取已发布文章 → 分片 → embed → VADD 写入 Redis
admin-server/internal/ollama/client.go                        → 手写：Ollama HTTP 客户端（Embeddings/Chat 两个方法）
admin-server/internal/svc/servicecontext.go                   → 新增 OllamaClient 字段 + 复用现有共享 Redis
admin-server/etc/admin-api.yaml                                → 新增 Ollama 配置段（BaseURL/EmbedModel/ChatModel）
```

## Code Style（示例）

Ollama 客户端方法签名风格，对齐项目现有 client 封装习惯（参考 `services/*/internal/repository` 里 error wrap 方式）：

```go
type Client struct {
    baseURL    string
    httpClient *http.Client
}

func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
    // POST {baseURL}/api/embeddings
}

func (c *Client) Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
    // POST {baseURL}/api/chat，非流式（stream: false）
}
```

Redis Vector Set 原始命令封装：

```go
// VADD key FP32 <blob> element [SETATTR json]
err := redisClient.Do(ctx, "VADD", "ai:kb:blog", "VALUES", dim, vectorBytes, elementID, "SETATTR", attrJSON).Err()
// VSIM key VALUES dim <blob> COUNT k WITHATTRIBS
res := redisClient.Do(ctx, "VSIM", "ai:kb:blog", "VALUES", dim, queryVectorBytes, "COUNT", topK, "WITHATTRIBS")
```

## Testing Strategy

- 单测：Ollama 客户端用 httptest mock 服务端，验证请求体/响应解析正确（不依赖真实 Ollama 进程）
- 集成验证（手工，非自动化）：本机启动 Ollama + Redis 后跑通一次完整链路，用真实已发布文章验证
- 不追求覆盖率指标，Phase 1 是原型验证性质

## Success Criteria

1. 对一篇已发布博客文章执行 Reindex 后，Redis 里能查到对应向量元素（`VCARD ai:kb:blog` > 0）
2. 针对文章明确覆盖的内容提问，Ask 接口返回的回答包含该文章的关键信息
3. 问一个知识库里完全没有的问题，回答明确表示"未找到相关内容"，而不是编造答案（RAG 防幻觉的核心验收点）
4. 全链路验证过程不产生任何外部付费 API 调用记录

## Boundaries

- **Always**：新增 API 遵循 `admin.api` 现有规范（group snake_case、query/可选字段加 `optional`、中间件顺序 Performance→RateLimit→Auth→Permission→OperationLog）；Ollama BaseURL/模型名放配置文件，不硬编码
- **Ask first**：`admin.api` 改完之后执行 `generate-api.sh`（必须用户亲自跑）；确定具体拉取哪几个 Ollama 模型、本机资源是否够用；Phase 2 何时启动、是否需要拆独立 `ai-rpc` 服务
- **Never**：不直接访问 content 服务的 repository/model 取文章数据（必须走 `ContentRPC` 现有接口）；不把 Ollama 调用失败静默吞掉导致返回假答案；不引入任何需要 API Key 的云端模型（违反"零外部费用"这个 Phase 1 的核心目标）

## Plan（技术实施步骤，待你 review）

**预研已完成的两项风险验证**（这两点原来是"计划里的风险"，已经用真实环境验证过，不再是假设）：

1. **Vector Set 命令语法已验证**：在本机 Redis 上用 `redis-cli` 跑通了 `VADD ai:test:probe VALUES 3 1.0 0.0 0.0 elem1 SETATTR '{"text":"hello"}'` → `VSIM ... WITHATTRIBS` → 正确返回 `elem1` + `{"text":"hello"}`，`VCARD`/`VEMB` 也符合预期，测试 key 已清理。Go 侧按这个语法拼 `Do()` 命令没问题。
2. **ContentRPC 现状已确认**：`BlogArticleList`/`PublicBlogArticleList` 两个列表接口**都不返回正文**（只有摘要），只有按 ID 查询的 `BlogArticleDetail`/`PublicBlogArticleDetail` 才带 `content` 全文字段。所以 Reindex 逻辑必须是"先拿已发布文章 ID 列表，再逐篇调 Detail 取全文"，不存在一次性批量拿全文的接口——文章数量不多的情况下这个 N+1 调用可以接受，量大了再考虑加个批量 RPC。

**实施步骤（按依赖顺序）**：

| # | 步骤 | 产出 | 依赖 |
|---|---|---|---|
| 1 | 编辑 `admin.api`：新增 group `ai/knowledge_qa`，定义 `KnowledgeQaAskReq{question string}`/`Resp{answer string, sources []string}`、`KnowledgeQaReindexReq{}`/`Resp{articleCount, chunkCount int64}`，两个 POST 接口，中间件走 Auth+Permission（管理员功能，不对公众开放） | admin.api diff | 无 |
| 2 | **你执行** `generate-api.sh` | handler 骨架 `internal/handler/ai/knowledge_qa/` | 步骤1 |
| 3 | 手写 `internal/ollama/client.go`：`Embed(ctx, text) ([]float32, error)`、`Chat(ctx, system, user string) (string, error)`，httptest mock 单测 | 可独立开发 | 无（可与步骤4并行） |
| 4 | 手写 `internal/vectorstore/blogindex.go`：`Add(ctx, elementID string, vec []float32, attr map[string]any) error`、`Search(ctx, vec []float32, topK int) ([]SearchResult, error)`，封装 VADD/VSIM 的 `Do()` 调用+ WITHATTRIBS 结果解析 | 可独立开发 | 无（可与步骤3并行） |
| 5 | `servicecontext.go` 新增 `OllamaClient`/复用现有共享 `Redis` 字段；`etc/admin-api.yaml` 新增 Ollama 配置段（BaseURL/EmbedModel=bge-m3/ChatModel=qwen2.5:7b） | 组装依赖 | 步骤3、4 |
| 6 | 实现 `knowledge_qa_reindex_logic.go`：`ContentRPC.BlogArticleList`（status=上架）拿 ID → 循环 `BlogArticleDetail` 拿 title+content → 按"先分段再定长切分，50 字重叠"分片 → 逐片 `Embed` → `vectorstore.Add`（element id = `<articleID>:<chunkIndex>`，attrib 存 `{article_id, title, chunk_text}`） | Reindex 接口可用 | 步骤2、5 |
| 7 | 实现 `knowledge_qa_ask_logic.go`：`Embed(question)` → `vectorstore.Search(topK=5)` → 若无结果直接返回"未找到相关内容"（不调用 Chat，防止空检索还硬编答案）→ 否则拼 system prompt（"只根据以下资料回答，资料未提及的内容如实说不知道"）+ 检索片段 + 问题 → `Chat` → 返回 answer + 命中文章标题去重列表 | Ask 接口可用 | 步骤2、5 |
| 8 | 手工联调：跑 Reindex → 跑 Ask（覆盖问题/无关问题两种）→ 对照 Success Criteria 三条逐一验证 | Phase 1 验收 | 步骤6、7 |

**待你 review 的点**：中间件顺序（步骤1，这是管理员功能不对外开放，对不对？）、Ollama 调用失败时的降级行为（步骤7，比如 Ollama 没启动时 Ask 接口应该报什么错，而不是 panic）。其余步骤技术路径明确，review 通过后我会拆成 Tasks 逐条落地。

## 已决定的技术细节（不再是待确认项，仅记录取舍）

- **分片策略**：先按 `\n\n` 切段落，单段超过 ~500 字再按固定长度二次切分，切片间保留约 50 字重叠，避免关键信息被切断在片段边界。
- **Reindex 触发方式**：Phase 1 只做手动触发（管理员调用 Reindex 接口），不联动文章发布事件自动增量更新——自动化留给 Phase 2 或后续迭代，避免 Phase 1 引入额外的事件耦合。
- **文档定稿位置**：`admin-server/docs/ai-knowledge-qa-spec.md`，作为独立阶段性规划文档维护，专项结束（Phase 1 验证通过并转normal功能或废弃）后按项目文档生命周期规则退场/归档。
- **Vector Set 落地实现**：项目锁定的 go-zero v1.9.3 的 `*redis.Redis` 没有通用 `Do`/`DoCtx`（v1.10+ 才有），实际走的是 `EvalCtx` 跑一段 `redis.call(cmd, KEYS[1], unpack(ARGV))` 的 Lua 脚本转发 VADD/VSIM，已用 redis-cli 验证 VADD/VSIM 未被标记 noscript、可以被 EVAL 调用。
- **RBAC 初始化 SQL**：`db/services/iam/knowledge_qa/init_knowledge_qa.sql`，仿 `db/services/iam/demo/init_demo.sql` 先例——Phase 1 没有前端页面，不建 `admin_menu`，只插入两个权限（`ai_knowledge_qa:ask`/`ai_knowledge_qa:reindex`）+ 对应接口 + 超级管理员角色关联，需要用户亲自执行（DB SQL 一律用户执行）。
- **网关超时**：`etc/admin-api.yaml` 全局 `Timeout` 从 30s 调到 120s（Ask 调本地 7B 模型可能超过原来的 30s）。这是网关级全局超时，会影响所有路由；Reindex 全量重建大量文章时仍可能超出 120s，Phase 1 测试请控制在少量文章，量大了应该把 Reindex 迁到已有的 task-rpc 异步任务体系，而不是继续调大这个全局超时。
