/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package tool

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ListFile(baseFolder string, handleFile func(path string, info os.FileInfo)) error {
	return filepath.Walk(baseFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			handleFile(path, info)
		}
		return nil
	})
}

func GetFileMd5(file *os.File) string {
	_, _ = file.Seek(0, 0)
	md5h := md5.New()
	_, _ = io.Copy(md5h, file)
	sum := fmt.Sprintf("%x", md5h.Sum(nil))
	return sum
}

func ToSlash(path string) string {
	return filepath.ToSlash(path)
}

func FormSlash(path string) string {
	return filepath.FromSlash(path)
}

//ToSlash(path string) string
//	功能：此函数将路径字符串 path 中的所有平台特定的分隔符（由全局变量 Separator 表示）替换为正斜杠 '/'。这对于需要将路径标准化为通用（通常是 Unix 风格）路径格式的场景非常有用，比如在网络传输或跨平台兼容的场景下。
//	示例：如果在 Windows 平台上，path 是 "C:\Users\Example"，并且 Separator 是 \，则调用 ToSlash(path) 后，结果会是 "C:/Users/Example"。
//FromSlash(path string) string
//	功能：与 ToSlash 相反，此函数将路径字符串 path 中的所有正斜杠 '/' 替换为平台特定的分隔符（同样由 Separator 表示）。这在需要将通用路径格式转换为特定于当前操作系统路径格式时很有用，比如在准备保存文件或访问本地文件系统资源时。
//	示例：如果 path 是 "C:/Users/Example"，并且假设当前平台是 Windows（因此 Separator 是 \），调用 FromSlash(path) 后，结果会是 "C:\Users\Example"。
//总结
//	ToSlash 用于将路径转换为使用正斜杠的通用格式。
//	FromSlash 则将使用正斜杠的路径转换回当前操作系统所使用的路径分隔符格式。
//	这两个函数在处理跨平台路径问题时特别有用，尤其是在编写需要兼容多种操作系统（如 Windows、Linux、macOS）的程序时。
