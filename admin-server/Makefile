.PHONY: init wire help

# 首次环境：安装 wire CLI（与 gameapi init 对齐）
init:
	go install github.com/google/wire/cmd/wire@latest

# wire 生成目录（单入口；若未来拆 cmd/task 可扩展为多个）
WIRE_CMD_DIRS := ./internal/wire
# 将 .cursor/rules/*.mdc（SSOT）转为 .claude/rules/*.md；缺失时自动创建 .claude/skills 软链
sync-claude-rules:
	@go run ../script/sync_claude_rules.go
# 改 providers.go / wire.go 后执行，重生 wire_gen.go
wire:
	@echo "Generating wire code..."
	@go get github.com/google/wire/cmd/wire@v0.6.0 2>/dev/null || true
	@for dir in $(WIRE_CMD_DIRS); do go run github.com/google/wire/cmd/wire $$dir; done
	@go mod tidy

help:
	@echo "Available commands:"
	@echo "  make init   - Install wire CLI"
	@echo "  make wire   - Regenerate internal/wire/wire_gen.go"
