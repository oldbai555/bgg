package progresstools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseChecklistEntriesAgainstRealFile(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("../../../../docs", "14-production-deployment-checklist.md"))
	if err != nil {
		t.Fatalf("读取 14-production-deployment-checklist.md 失败: %v", err)
	}

	entries := parseChecklistEntries(string(raw))
	if len(entries) == 0 {
		t.Fatal("parseChecklistEntries() 未解析出任何条目")
	}

	first := entries[0]
	if first.Index != 1 {
		t.Errorf("first.Index = %d, want 1", first.Index)
	}
	if first.Title != "JWT 密钥改用环境变量注入" {
		t.Errorf("first.Title = %q, want %q", first.Title, "JWT 密钥改用环境变量注入")
	}
	if first.Status != "已就绪，待执行" {
		t.Errorf("first.Status = %q, want %q", first.Status, "已就绪，待执行")
	}
	if first.Trigger == "" || first.Action == "" || first.Verification == "" {
		t.Errorf("first entry has empty field: trigger=%q action=%q verification=%q", first.Trigger, first.Action, first.Verification)
	}
}

func TestParseChecklistEntriesSynthetic(t *testing.T) {
	content := "# 标题\n\n" +
		"### 1 · 示例条目\n\n" +
		"**触发条件**：改动 A\n\n" +
		"**部署时要做什么**：\n1. 步骤一\n2. 步骤二\n\n" +
		"**如何验证生效**：检查日志\n\n" +
		"**状态**：`已执行`（2026-01-01 已上线）。\n\n" +
		"### 2 · 第二条\n\n" +
		"**触发条件**：改动 B\n\n" +
		"**部署时要做什么**：步骤\n\n" +
		"**如何验证生效**：检查\n\n" +
		"**状态**：`TBD`。\n"

	entries := parseChecklistEntries(content)
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}
	if entries[0].Status != "已执行" {
		t.Errorf("entries[0].Status = %q, want 已执行", entries[0].Status)
	}
	if entries[1].Status != "TBD" {
		t.Errorf("entries[1].Status = %q, want TBD", entries[1].Status)
	}
}

func TestQueryDeploymentChecklistPendingOnlyFilters(t *testing.T) {
	content := "### 1 · A\n\n**触发条件**：x\n\n**部署时要做什么**：x\n\n**如何验证生效**：x\n\n**状态**：`已执行`。\n\n" +
		"### 2 · B\n\n**触发条件**：x\n\n**部署时要做什么**：x\n\n**如何验证生效**：x\n\n**状态**：`已就绪，待执行`。\n"
	entries := parseChecklistEntries(content)

	pending := 0
	for _, e := range entries {
		if e.Status != "已执行" {
			pending++
		}
	}
	if pending != 1 {
		t.Fatalf("pending count = %d, want 1", pending)
	}
}
