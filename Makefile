.PHONY: sync-claude-rules

# 将 .cursor/rules/*.mdc（SSOT）转为 .claude/rules/*.md；缺失时自动创建 .claude/skills 软链
sync-claude-rules:
	@go run ./script/sync_claude_rules.go