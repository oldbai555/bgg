package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\branches\\20241226 C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\trunk
func main() {
	if len(os.Args) < 3 {
		fmt.Println("请提供要Merge的文件路径")
		return
	}
	branch := os.Args[1]
	branch = filepath.FromSlash(branch)
	trunk := os.Args[2]
	trunk = filepath.FromSlash(trunk)
	err := DoMerge(trunk, branch, false)
	if err != nil {
		os.Exit(1)
		return
	}
}
