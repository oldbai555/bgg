#!/usr/bin/env node
import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js'
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js'
import { z } from 'zod'
import {
  handleGetComponent,
  handleGetPatterns,
  handleListComponents,
} from './tools.js'

const server = new McpServer({
  name: 'admin-frontend-ui',
  version: '0.1.0',
})

server.tool(
  'ui_list_components',
  '列出 admin-frontend 项目内 UI 组件（common/layout/blog），含路径与一句话说明',
  {
    category: z
      .enum(['common', 'layout', 'blog'])
      .optional()
      .describe('按分类筛选，不传则返回全部'),
  },
  async ({ category }) => ({
    content: [{ type: 'text', text: handleListComponents(category) }],
  }),
)

server.tool(
  'ui_get_component',
  '获取项目 UI 组件文档：README、关联类型（如 D2Table 的 table.ts）',
  {
    name: z.string().describe('组件名，如 D2Table、IcpFooter、BlogHeader'),
  },
  async ({ name }) => ({
    content: [{ type: 'text', text: handleGetComponent(name) }],
  }),
)

server.tool(
  'ui_get_patterns',
  '获取前端开发约定片段：dict / permission / api / public-pages',
  {
    topic: z
      .enum(['dict', 'permission', 'api', 'public-pages'])
      .optional()
      .describe('约定主题，不传则返回索引'),
  },
  async ({ topic }) => ({
    content: [{ type: 'text', text: handleGetPatterns(topic) }],
  }),
)

async function main() {
  const transport = new StdioServerTransport()
  await server.connect(transport)
}

main().catch((err) => {
  console.error('frontend-ui MCP failed:', err)
  process.exit(1)
})
