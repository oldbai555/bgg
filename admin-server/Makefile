.PHONY: init wire help

# 首次环境：安装 wire CLI（与 gameapi init 对齐）
init:
	go install github.com/google/wire/cmd/wire@latest

# wire 生成目录（单入口；若未来拆 cmd/task 可扩展为多个）
WIRE_CMD_DIRS := ./internal/wire

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
