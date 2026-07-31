#!/usr/bin/env bash
# 本机 Ollama 服务管理脚本（AI 知识库问答 ai/knowledge_qa 用，模型/端点配置见
# admin-server/etc/admin-api.yaml 的 Ollama 段，详见 admin-server/docs/ai-knowledge-qa-spec.md）。
set -euo pipefail

OLLAMA_HOST="${OLLAMA_HOST:-http://127.0.0.1:11434}"
EMBED_MODEL="${OLLAMA_EMBED_MODEL:-bge-m3}"
CHAT_MODEL="${OLLAMA_CHAT_MODEL:-qwen2.5:7b}"

usage() {
  cat <<EOF
用法: $(basename "$0") <command>

命令:
  start     启动 Ollama 后台服务（brew services start ollama）
  stop      停止 Ollama 后台服务
  restart   重启 Ollama 后台服务
  status    查看服务运行状态 + 已装模型列表
  models    列出已安装模型
  pull      拉取本项目需要的模型（${EMBED_MODEL} + ${CHAT_MODEL}）
  test      冒烟测试：跑一次 embeddings + chat，验证服务可正常响应
  logs      跟踪 Homebrew 服务日志（Ctrl+C 退出）

环境变量:
  OLLAMA_HOST         默认 ${OLLAMA_HOST}
  OLLAMA_EMBED_MODEL   默认 ${EMBED_MODEL}
  OLLAMA_CHAT_MODEL    默认 ${CHAT_MODEL}
EOF
}

require_ollama() {
  if ! command -v ollama >/dev/null 2>&1; then
    echo "错误: 未找到 ollama 可执行文件，先执行 brew install ollama" >&2
    exit 1
  fi
}

cmd_start() {
  brew services start ollama
}

cmd_stop() {
  brew services stop ollama
}

cmd_restart() {
  brew services restart ollama
}

cmd_status() {
  brew services info ollama
  echo
  echo "已安装模型:"
  ollama list
}

cmd_models() {
  ollama list
}

cmd_pull() {
  ollama pull "$EMBED_MODEL"
  ollama pull "$CHAT_MODEL"
}

cmd_test() {
  echo "== 测试 embeddings（模型: ${EMBED_MODEL}）=="
  curl -sf "${OLLAMA_HOST}/api/embeddings" \
    -d "{\"model\":\"${EMBED_MODEL}\",\"prompt\":\"你好，世界\"}" \
    | python3 -c 'import json,sys; d=json.load(sys.stdin); print("向量维度:", len(d["embedding"]))'

  echo
  echo "== 测试 chat（模型: ${CHAT_MODEL}）=="
  curl -sf "${OLLAMA_HOST}/api/chat" \
    -d "{\"model\":\"${CHAT_MODEL}\",\"messages\":[{\"role\":\"user\",\"content\":\"用一句话介绍你自己\"}],\"stream\":false}" \
    | python3 -c 'import json,sys; d=json.load(sys.stdin); print(d["message"]["content"])'
}

cmd_logs() {
  local log_file
  log_file="$(brew services info ollama --json 2>/dev/null | python3 -c 'import json,sys; print(json.load(sys.stdin)[0].get("log_path") or "")' 2>/dev/null || true)"
  if [[ -n "$log_file" && -f "$log_file" ]]; then
    tail -f "$log_file"
  else
    echo "找不到日志路径，尝试默认位置: /opt/homebrew/var/log/ollama.log" >&2
    tail -f /opt/homebrew/var/log/ollama.log
  fi
}

require_ollama

case "${1:-}" in
  start) cmd_start ;;
  stop) cmd_stop ;;
  restart) cmd_restart ;;
  status) cmd_status ;;
  models) cmd_models ;;
  pull) cmd_pull ;;
  test) cmd_test ;;
  logs) cmd_logs ;;
  *) usage; exit 1 ;;
esac
