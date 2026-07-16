// sync_claude_rules 将 .cursor/rules/*.mdc（SSOT）转为 Claude Code 可读的 .claude/rules/*.md。
//
// Cursor frontmatter: description, globs, alwaysApply
// Claude frontmatter: paths（由 globs 映射）、alwaysApply: false（路径懒加载）
//
// alwaysApply: true 或无 globs 的规则不输出 frontmatter（会话级无条件加载）。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	cursorRulesDir     = ".cursor/rules"
	cursorSkillsDir    = ".cursor/skills"
	claudeDir          = ".claude"
	claudeRulesDir     = ".claude/rules"
	claudeSkillsLink   = ".claude/skills"
	claudeSkillsTarget = "../.cursor/skills"
)

func main() {
	check := flag.Bool("check", false, "校验 .claude/rules 是否与 .cursor/rules 同步")
	flag.Parse()

	root, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync_claude_rules: %v\n", err)
		os.Exit(1)
	}

	mdcFiles, err := filepath.Glob(filepath.Join(root, cursorRulesDir, "*.mdc"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync_claude_rules: glob: %v\n", err)
		os.Exit(1)
	}
	if len(mdcFiles) == 0 {
		fmt.Fprintf(os.Stderr, "sync_claude_rules: no .mdc files in %s\n", cursorRulesDir)
		os.Exit(1)
	}
	sort.Strings(mdcFiles)

	outDir := filepath.Join(root, claudeRulesDir)
	expected := make(map[string][]byte, len(mdcFiles))

	for _, mdcPath := range mdcFiles {
		base := strings.TrimSuffix(filepath.Base(mdcPath), ".mdc")
		outPath := filepath.Join(outDir, base+".md")

		raw, err := os.ReadFile(mdcPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sync_claude_rules: read %s: %v\n", mdcPath, err)
			os.Exit(1)
		}

		content, err := convertMDCToClaudeRule(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sync_claude_rules: convert %s: %v\n", mdcPath, err)
			os.Exit(1)
		}
		expected[outPath] = content
	}

	if *check {
		if err := checkSynced(expected); err != nil {
			fmt.Fprintf(os.Stderr, "rules-check: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("rules-check: ok")
		return
	}

	if err := ensureClaudeSkillsSymlink(root); err != nil {
		fmt.Fprintf(os.Stderr, "sync_claude_rules: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "sync_claude_rules: mkdir %s: %v\n", outDir, err)
		os.Exit(1)
	}

	for outPath, content := range expected {
		if err := os.WriteFile(outPath, content, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "sync_claude_rules: write %s: %v\n", outPath, err)
			os.Exit(1)
		}
	}

	// 删除已无对应 .mdc 的陈旧 .md
	existing, err := filepath.Glob(filepath.Join(outDir, "*.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync_claude_rules: glob out: %v\n", err)
		os.Exit(1)
	}
	expectedSet := make(map[string]struct{}, len(expected))
	for p := range expected {
		expectedSet[p] = struct{}{}
	}
	for _, p := range existing {
		if _, ok := expectedSet[p]; !ok {
			if err := os.Remove(p); err != nil {
				fmt.Fprintf(os.Stderr, "sync_claude_rules: remove stale %s: %v\n", p, err)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("sync-claude-rules: wrote %d files to %s\n", len(expected), claudeRulesDir)
}

func ensureClaudeSkillsSymlink(root string) error {
	skillsSrc := filepath.Join(root, cursorSkillsDir)
	if _, err := os.Stat(skillsSrc); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s not found; cannot link %s", cursorSkillsDir, claudeSkillsLink)
		}
		return fmt.Errorf("stat %s: %w", cursorSkillsDir, err)
	}

	if err := os.MkdirAll(filepath.Join(root, claudeDir), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", claudeDir, err)
	}

	linkPath := filepath.Join(root, claudeSkillsLink)
	fi, err := os.Lstat(linkPath)
	if err == nil {
		if fi.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("%s exists but is not a symlink", claudeSkillsLink)
		}
		if sameSymlinkTarget(linkPath, skillsSrc) {
			return nil
		}
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("remove stale %s: %w", claudeSkillsLink, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("lstat %s: %w", claudeSkillsLink, err)
	}

	if err := os.Symlink(claudeSkillsTarget, linkPath); err != nil {
		return fmt.Errorf("symlink %s -> %s: %w", claudeSkillsLink, claudeSkillsTarget, err)
	}
	fmt.Printf("sync-claude-rules: linked %s -> %s\n", claudeSkillsLink, claudeSkillsTarget)
	return nil
}

func sameSymlinkTarget(linkPath, wantTarget string) bool {
	current, err := os.Readlink(linkPath)
	if err != nil {
		return false
	}
	if current == claudeSkillsTarget {
		return true
	}
	resolved, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		return false
	}
	want, err := filepath.EvalSymlinks(wantTarget)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(filepath.Dir(linkPath), want)
	if err != nil {
		return false
	}
	return resolved == want || current == rel
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, cursorRulesDir)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cannot find %s from cwd", cursorRulesDir)
		}
		dir = parent
	}
}

type mdcFrontmatter struct {
	Globs       string
	AlwaysApply *bool
}

func parseMDC(data []byte) (mdcFrontmatter, string, error) {
	s := string(data)
	if !strings.HasPrefix(s, "---\n") {
		return mdcFrontmatter{}, strings.TrimLeft(s, "\n"), nil
	}

	rest := s[len("---\n"):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return mdcFrontmatter{}, "", fmt.Errorf("unclosed frontmatter")
	}

	fmRaw := rest[:end]
	body := strings.TrimLeft(rest[end+len("\n---"):], "\n")

	var fm mdcFrontmatter
	for _, line := range strings.Split(fmRaw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		switch key {
		case "globs":
			fm.Globs = val
		case "alwaysApply":
			switch strings.ToLower(val) {
			case "true":
				t := true
				fm.AlwaysApply = &t
			case "false":
				f := false
				fm.AlwaysApply = &f
			}
		}
	}

	return fm, body, nil
}

func convertMDCToClaudeRule(raw []byte) ([]byte, error) {
	fm, body, err := parseMDC(raw)
	if err != nil {
		return nil, err
	}

	// 常驻规则或无路径范围：Claude 无条件加载，仅输出正文。
	if fm.AlwaysApply != nil && *fm.AlwaysApply {
		return []byte(body), nil
	}
	if strings.TrimSpace(fm.Globs) == "" {
		return []byte(body), nil
	}

	// 路径懒加载：paths 用 CSV 单行（Claude Code 解析器兼容格式）。
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("alwaysApply: false\n")
	b.WriteString("paths: ")
	b.WriteString(fm.Globs)
	b.WriteString("\n---\n\n")
	b.WriteString(body)
	return []byte(b.String()), nil
}

func checkSynced(expected map[string][]byte) error {
	var stale []string
	for outPath, want := range expected {
		got, err := os.ReadFile(outPath)
		if err != nil {
			if os.IsNotExist(err) {
				stale = append(stale, fmt.Sprintf("missing %s", rel(outPath)))
				continue
			}
			return err
		}
		if string(got) != string(want) {
			stale = append(stale, fmt.Sprintf("drift %s (run make sync-claude-rules)", rel(outPath)))
		}
	}

	outDir := filepath.Dir(firstKey(expected))
	existing, err := filepath.Glob(filepath.Join(outDir, "*.md"))
	if err != nil {
		return err
	}
	for _, p := range existing {
		if _, ok := expected[p]; !ok {
			stale = append(stale, fmt.Sprintf("stale %s (no matching .mdc)", rel(p)))
		}
	}

	if len(stale) > 0 {
		sort.Strings(stale)
		return fmt.Errorf("out of sync:\n  %s", strings.Join(stale, "\n  "))
	}
	return nil
}

func firstKey(m map[string][]byte) string {
	for k := range m {
		return k
	}
	return ""
}

func rel(path string) string {
	if wd, err := os.Getwd(); err == nil {
		if r, err := filepath.Rel(wd, path); err == nil {
			return r
		}
	}
	return path
}
