#!/usr/bin/env bash
# MCP vue-lsp：在 admin-frontend workspace 内通过 node 启动 @vue/language-server
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 1
exec node "node_modules/@vue/language-server/bin/vue-language-server.js" --stdio
