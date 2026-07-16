// Package exec 提供跑 shell 脚本 + 自动确认 + git diff 收集产物列表的共享逻辑，
// scripttools 的 6 个 tool 共用，避免每个 tool 重复实现同一套子进程管理代码。
package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// ScriptResult 是 RunWithAutoConfirm 的返回结构，序列化为 MCP tool 的结构化响应。
type ScriptResult struct {
	Success        bool     `json:"success"`
	ExitCode       int      `json:"exit_code"`
	Stdout         string   `json:"stdout"`
	Stderr         string   `json:"stderr"`
	GeneratedFiles []string `json:"generated_files"`
}

// RunWithAutoConfirm 跑一个 generate-*.sh 脚本，向 stdin 喂 "y\n" 自动通过脚本自带的
// 交互式确认提示（子进程默认没有 TTY，read -p 会读到 EOF 直接判定"未确认"并退出，这里
// 显式提供确认输入，把 10-dev-execution-and-review-points.md 已经拍板的"开发期 AI 可以
// 直接执行 generate-*.sh"这条政策落到代码里，不是绕开政策）。
//
// GeneratedFiles 用 git status 前后对比推导，而不是硬编码每个脚本的输出文件名规律——
// 脚本内部命名规则变了也不用同步改这里；已知 generate-sql.sh 的 rm -f sqlgen 之后立刻判断
// $? 拿到的是 rm 的退出码而不是 sqlgen 本身的（12-scripts-standardization.md 已经点出这个
// bug，按该文档的任务修复，本工具不做绕过式的特殊处理），所以 ExitCode/Success 字段在这一个
// 脚本上不完全可信，调用方应该优先看 GeneratedFiles 是否非空。
func RunWithAutoConfirm(repoRoot, scriptPath string, args []string, watchDir string) (*ScriptResult, error) {
	before, err := gitStatusSnapshot(repoRoot, watchDir)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(scriptPath, args...)
	cmd.Dir = repoRoot
	cmd.Stdin = bytes.NewBufferString("y\n")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	after, err := gitStatusSnapshot(repoRoot, watchDir)
	if err != nil {
		return nil, err
	}

	result := &ScriptResult{
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		GeneratedFiles: diffFiles(before, after),
	}
	if exitErr, ok := runErr.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	} else if runErr == nil {
		result.ExitCode = 0
	} else {
		// 脚本本身没找到/没有执行权限等，属于工具环境错误，不是业务失败。
		return nil, runErr
	}
	result.Success = result.ExitCode == 0
	return result, nil
}

// gitStatusSnapshot 对 watchDir 跑一次 git status --porcelain=v1，返回"状态码 路径"
// 到路径的映射（用状态码 + 路径整行做 key，方便 diffFiles 精确识别新增/变化的条目）。
func gitStatusSnapshot(repoRoot, watchDir string) (map[string]struct{}, error) {
	cmd := exec.Command("git", "status", "--porcelain=v1", "--", watchDir)
	cmd.Dir = repoRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git status 失败: %w, stderr=%s", err, stderr.String())
	}

	lines := strings.Split(stdout.String(), "\n")
	snapshot := make(map[string]struct{}, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		snapshot[line] = struct{}{}
	}
	return snapshot, nil
}

// diffFiles 计算 after 相对 before 新增的 git status 条目，提取出路径部分并排序返回。
// porcelain=v1 每行形如 "XY path" 或 "XY orig -> path"（rename），这里统一取最后一段路径。
func diffFiles(before, after map[string]struct{}) []string {
	var added []string
	for line := range after {
		if _, existed := before[line]; existed {
			continue
		}
		added = append(added, extractPath(line))
	}
	sort.Strings(added)
	return added
}

// extractPath 从一行 "XY path" 或 "XY orig -> path" 中提取路径部分。
func extractPath(line string) string {
	if len(line) < 3 {
		return strings.TrimSpace(line)
	}
	rest := strings.TrimSpace(line[3:])
	if idx := strings.Index(rest, " -> "); idx >= 0 {
		return rest[idx+len(" -> "):]
	}
	return rest
}
