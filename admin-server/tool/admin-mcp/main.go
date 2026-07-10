// admin-mcp 是本仓库自建的项目专属 MCP server，通过 stdio 暴露给 Cursor/Claude Code，
// 覆盖三类能力：封装 scripts/generate-*.sh、项目约定查询、进度查询。设计见
// admin-server/docs/22-admin-mcp-tool.md，本文件只负责组装，具体实现在 internal/ 下三个包。
package main

import (
	"log"
	"os"
	"path/filepath"

	"admin-mcp/internal/conventiontools"
	"admin-mcp/internal/progresstools"
	"admin-mcp/internal/scripttools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	repoRoot, dataDir, docsDir, err := resolveDirs()
	if err != nil {
		log.Fatalf("解析目录失败: %v", err)
	}

	s := server.NewMCPServer("admin-mcp", "0.1.0")

	scripttools.Register(s, repoRoot)
	conventiontools.Register(s, dataDir)
	progresstools.Register(s, docsDir)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("admin-mcp server 退出: %v", err)
	}
}

// resolveDirs 从可执行文件自身路径推导出三个固定目录，不依赖进程启动时的工作目录
// （Cursor/Claude Code 启动 MCP server 子进程时的 cwd 不一定是 admin-server/）：
//   - repoRoot: admin-server/（二进制路径 tool/admin-mcp/bin/admin-mcp 向上 4 层）
//   - dataDir:  tool/admin-mcp/data（本模块自带的种子数据目录，向上 1 层的 data 子目录）
//   - docsDir:  admin-server/docs
func resolveDirs() (repoRoot, dataDir, docsDir string, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", "", "", err
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", "", "", err
	}

	binDir := filepath.Dir(exePath)      // .../tool/admin-mcp/bin
	adminMcpDir := filepath.Dir(binDir)  // .../tool/admin-mcp
	toolDir := filepath.Dir(adminMcpDir) // .../tool
	repoRoot = filepath.Dir(toolDir)     // .../admin-server
	dataDir = filepath.Join(adminMcpDir, "data")
	docsDir = filepath.Join(repoRoot, "docs")
	return repoRoot, dataDir, docsDir, nil
}
