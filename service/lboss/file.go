/**
 * @Author: zjj
 * @Date: 2024/5/7
 * @Desc:
**/

package main

import (
	"github.com/oldbai555/bgg/service/lboss/compress"
	"github.com/oldbai555/lbtool/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unsafe"
)

var (
	BaseStoragePath   = path.Join(utils.GetCurDir(), "storage")
	BaseTemplatesPath = path.Join(utils.GetCurDir(), "templates", "*")
)

const (
	MaxMultipartMemory = 1024 * 1024 * 512 // 最大支持512MB
)

type File struct {
	Name      string `json:"name"`
	ReName    string `json:"rename"`
	Path      string `json:"path"`
	Md5       string `json:"md5"`
	Size      int64  `json:"size"`
	TimeStamp int64  `json:"timeStamp"`
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func listFile(baseFolder string) (fileList []string, err error) {
	err = filepath.Walk(baseFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			fileList = append(fileList, strings.Replace(path, "\\", "/", -1))
		}
		return nil
	})
	return fileList, err
}

func syncFileIndex() ([]string, error) {
	fileList, err := listFile(BaseStoragePath)
	if err != nil {
		return nil, err
	}
	var sUrlList []string
	for i := 0; i < len(fileList); i++ {
		savePath := fileList[i]
		sUrl := compress.GenShortUrl(compress.CharsetRandomAlphanumeric, savePath, func(url, keyword string) bool {
			data, _ := dbConn.Get([]byte(keyword), nil)
			if data == nil {
				return true
			}
			return false
		})
		err = dbConn.Put([]byte(sUrl), []byte(savePath), nil)
		sUrlList = append(sUrlList, sUrl)
	}
	return sUrlList, nil
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
