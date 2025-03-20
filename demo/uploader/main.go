package main

import (
	"embed"
	"fmt"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

type UploadJson struct {
	ServerIP   string `json:"server_ip"`
	ServerUser string `json:"server_user"`
	ServerPass string `json:"server_pass"`
	RemoteDir  string `json:"remote_dir"`
}

//go:embed uploader.json
var configFile embed.FS

func initConfig() {
	configPath := path.Join(utils.GetCurDir(), "uploader.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configData, err := configFile.ReadFile("uploader.json")
		if err != nil {
			fmt.Printf("err:%v", err)
			return
		}
		err = os.WriteFile(configPath, configData, 0644)
		if err != nil {
			fmt.Printf("err:%v", err)
			return
		}
	}
}

func main() {
	// initConfig()
	//bytes, err := os.ReadFile(path.Join(utils.GetCurDir(), "uploader.json"))
	//if err != nil {
	//	fmt.Printf("err:%v", err)
	//	return
	//}
	bytes, err := configFile.ReadFile("uploader.json")
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	var uploadJson UploadJson
	err = json.Unmarshal(bytes, &uploadJson)
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	// 检查命令行参数是否正确
	if len(os.Args) < 2 {
		fmt.Println("请提供要上传的文件路径")
		return
	}
	filePath := os.Args[1]
	filePath = filepath.ToSlash(filePath)

	// 创建SSH客户端配置
	config := &ssh.ClientConfig{
		User: uploadJson.ServerUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(uploadJson.ServerPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境中应使用更安全的方式验证主机密钥
		Timeout:         10 * time.Second,
	}
	// 连接服务器
	client, err := ssh.Dial("tcp", uploadJson.ServerIP+":22", config)
	if err != nil {
		log.Fatalf("无法连接服务器: %v", err)
	}
	defer client.Close()

	remoteFilePath := uploaderBySftp(client, filePath, uploadJson.RemoteDir)
	fmt.Printf("文件 %s 已成功上传到服务器路径 %s\n", path.Base(filePath), remoteFilePath)
}

func uploaderBySftp(client *ssh.Client, filePath string, remoteDir string) string {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return ""
	}
	defer file.Close()

	// 创建SFTP会话
	sftp, err := NewSFTP(client)
	if err != nil {
		log.Fatalf("无法创建SFTP会话: %v", err)
	}
	defer sftp.Close()

	// 上传文件
	remoteFilePath := path.Join(remoteDir, path.Base(filePath))
	err = sftp.UploadFile(file, remoteFilePath)
	if err != nil {
		log.Fatalf("文件上传失败: %v", err)
	}
	return remoteFilePath
}

func uploadToHttp(client *ssh.Client, remoteFilePath string) {
	// 构造 curl 命令，调用服务器上的 HTTP 服务上传文件
	curlCommand := fmt.Sprintf("sh /home/work/package/supervisor.sh upload %s", remoteFilePath)
	// 创建 SSH 会话
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return
	}
	defer session.Close()
	// 执行命令
	output, err := session.CombinedOutput(curlCommand)
	if err != nil {
		log.Fatalf("文件上传失败: %v, 输出: %s", err, string(output))
	}
}
