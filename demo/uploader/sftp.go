/**
 * @Author: zjj
 * @Date: 2025/3/20
 * @Desc:
**/

package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
)

// SFTPClient 定义SFTP客户端接口
type SFTPClient interface {
	UploadFile(file *os.File, remotePath string) error
	Close() error
}

// sftpClient 实现SFTP客户端
type sftpClient struct {
	sftp *sftp.Client
}

// NewSFTP 创建一个新的SFTP会话
func NewSFTP(client *ssh.Client) (SFTPClient, error) {
	sftpSession, err := sftp.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("无法创建SFTP会话: %v", err)
	}
	return &sftpClient{sftp: sftpSession}, nil
}

// EnsureRemoteDir 确保远程目录存在，如果不存在则创建
func (s *sftpClient) EnsureRemoteDir(remotePath string) error {
	dir := remotePath[:len(remotePath)-len(filepath.Base(remotePath))] // 提取目录部分
	_, err := s.sftp.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果目录不存在，则递归创建
			return s.sftp.MkdirAll(dir)
		}
		return fmt.Errorf("无法检查远程目录: %v", err)
	}
	return nil
}

// UploadFile 上传文件到远程服务器
func (s *sftpClient) UploadFile(file *os.File, remotePath string) error {
	// 确保远程目录存在
	if err := s.EnsureRemoteDir(remotePath); err != nil {
		return fmt.Errorf("无法确保远程目录: %v", err)
	}

	// 打开远程文件以写入
	remoteFile, err := s.sftp.Create(remotePath)
	if err != nil {
		return fmt.Errorf("无法创建远程文件: %v", err)
	}
	defer remoteFile.Close()

	// 将本地文件内容复制到远程文件
	_, err = io.Copy(remoteFile, file)
	if err != nil {
		return fmt.Errorf("文件上传失败: %v", err)
	}
	return nil
}

// Close 关闭SFTP会话
func (s *sftpClient) Close() error {
	return s.sftp.Close()
}
