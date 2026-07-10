.PHONY: sync-claude-rules sync-claude-mcp-check sync-claude-mcp-approve sync-claude-mcp-import engram-sync engram-sync-export engram-sync-import engram-sync-status engram-sync-push engram-sync-pull setup-ai setup-ai-check

# 将 .cursor/rules/*.mdc（SSOT）转为 .claude/rules/*.md；缺失时自动创建 .claude/skills 软链
sync-claude-rules:
	@go run ./script/sync_claude_rules.go

# Claude Code 项目 MCP：.mcp.json 为团队 SSOT（已提交 git），见 docs/AI工具链上手.md
sync-claude-mcp-check:
	@bash ./script/sync_claude_mcp.sh check

sync-claude-mcp-approve:
	@bash ./script/sync_claude_mcp.sh approve

# 维护者：从 ~/.cursor/mcp.json 导入并规范化路径，更新 .mcp.json 后需 commit
sync-claude-mcp-import:
	@bash ./script/sync_claude_mcp.sh import-cursor

# Engram 跨设备记忆同步（需已安装 engram CLI）
engram-sync engram-sync-export:
	@bash ./script/engram_sync.sh export

engram-sync-import:
	@bash ./script/engram_sync.sh import

engram-sync-status:
	@bash ./script/engram_sync.sh status

# 离开当前设备：导出记忆并 commit .engram/（之后手动 git push）
engram-sync-push:
	@bash ./script/engram_sync.sh push

# 换设备开始开发：git pull 并导入记忆
engram-sync-pull:
	@bash ./script/engram_sync.sh pull

# 新设备 / 新维护者：初始化 Gentle-AI + CodeGraph + Engram（Cursor + Claude Code 插件，见 docs/AI工具链上手.md）
setup-ai:
	@bash ./script/setup_ai_toolchain.sh init

setup-ai-check:
	@bash ./script/setup_ai_toolchain.sh check