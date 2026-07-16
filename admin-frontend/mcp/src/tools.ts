import { readFileSync, existsSync } from 'node:fs'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { dirname } from 'node:path'
import { findComponent, listByCategory } from './catalog.js'

const __dirname = dirname(fileURLToPath(import.meta.url))
/** admin-frontend 根目录（mcp/dist -> mcp -> admin-frontend） */
export const FRONTEND_ROOT = join(__dirname, '..', '..')
/** 仓库根目录 */
export const REPO_ROOT = join(FRONTEND_ROOT, '..')

function readText(path: string): string | null {
  if (!existsSync(path)) {
    return null
  }
  return readFileSync(path, 'utf-8')
}

export function handleListComponents(category?: string): string {
  const items = listByCategory(category)
  if (items.length === 0) {
    return category
      ? `未找到分类 "${category}" 下的组件。可用分类：common、layout、blog`
      : '组件目录为空'
  }
  const lines = items.map(
    (c) => `- **${c.name}** (${c.category}): ${c.summary}\n  路径: \`${c.path}\``,
  )
  const header = category
    ? `## 组件列表（${category}）\n\n共 ${items.length} 个\n\n`
    : `## 项目 UI 组件列表\n\n共 ${items.length} 个（分类：common / layout / blog）\n\n`
  return header + lines.join('\n\n')
}

export function handleGetComponent(name: string): string {
  const entry = findComponent(name)
  if (!entry) {
    return `未找到组件 "${name}"。请用 ui_list_components 查看可用组件。`
  }

  const parts: string[] = [
    `# ${entry.name}`,
    '',
    `- **分类**: ${entry.category}`,
    `- **路径**: \`${entry.path}\``,
    `- **说明**: ${entry.summary}`,
    '',
  ]

  if (entry.readme) {
    const readmePath = join(FRONTEND_ROOT, entry.readme)
    const readme = readText(readmePath)
    if (readme) {
      parts.push('## 文档', '', readme, '')
    }
  }

  if (entry.relatedTypes) {
    for (const rel of entry.relatedTypes) {
      const typePath = join(FRONTEND_ROOT, rel)
      const content = readText(typePath)
      if (content) {
        parts.push(`## 关联类型 (\`${rel}\`)`, '', '```typescript', content, '```', '')
      }
    }
  }

  if (entry.name === 'D2Table' && !entry.readme) {
    parts.push(
      '## 提示',
      '',
      'D2Table 文档见 `src/components/common/README.md`；完整示例参考 `src/views/system/RoleList.vue`。',
    )
  }

  return parts.join('\n')
}

const PATTERNS: Record<string, string> = {
  dict: `## 字典约定

- 所有下拉/单选/复选选项必须来自字典，禁止硬编码
- 使用 \`useDictOptions\` + \`stores/dict.ts\`
- 新增字典 code 须加入 \`REQUIRED_DICT_CODES\`（登录后批量加载）
- 字典 value 从 1 开始；筛选参数用 0 表示「全部/不筛选」

\`\`\`ts
import { useDictOptions } from '@/composables/useDictOptions'
const { options: statusOptions } = useDictOptions('notice_status')
\`\`\``,

  permission: `## 权限约定

- 按钮：\`v-permission="'module:action'"\`（如 \`blog_tag:create\`）
- 路由：\`meta.permission\` 与动态菜单联动
- 列表操作列权限：D2Table 不内置权限，需自定义列或隐藏 haveEdit/haveDetail`,

  api: `## API 层约定

- 生成入口：\`admin-server/scripts/generate-ts.sh\`（**用户执行**，AI 不得代替）
- 产物：\`src/api/generated/\`，**禁止手改**
- 业务代码只从 \`src/api/*.ts\` 二次封装层导入
- \`package.json\` 的 \`api:gen\` 已失效，勿用
- 时间字段：后端 int64 秒级时间戳，前端展示层格式化`,

  'public-pages': `## 公共展示页约定（views/public、components/blog）

- 根类名：\`public-list-page\` 或 \`public-detail-page\`
- 列表：\`@import '@/styles/public-list.scss'\`；暖色渐变 + 白色卡片
- 详情：\`@import '@/styles/public-detail.scss'\`；居中白卡片 max-width ~800px
- 断点：\`@media (max-width: 768px)\`
- 列表进详情前：分页/搜索/滚动写入 sessionStorage，返回时恢复
- 埋点：统一用 \`MetricReporter\`（勿各写 metricApi.report）
- 底部：\`IcpFooter\` 必挂`,
}

export function handleGetPatterns(topic?: string): string {
  if (!topic) {
    const keys = Object.keys(PATTERNS)
    const previews = keys.map((k) => `### ${k}\n\n${PATTERNS[k].split('\n').slice(0, 3).join('\n')}...`)
    return [
      '## 前端约定索引',
      '',
      '可用 topic：`dict`、`permission`、`api`、`public-pages`',
      '',
      ...previews,
    ].join('\n\n')
  }

  const key = topic.trim().toLowerCase()
  const pattern = PATTERNS[key]
  if (!pattern) {
    return `未知 topic "${topic}"。可用：${Object.keys(PATTERNS).join('、')}`
  }
  return pattern
}
