/**
 * @Author: zjj
 * @Date: 2025/2/24
 * @Desc:
**/

package miniosdk

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"io"
	"net/url"
	"path/filepath"
	"time"
)

type Client struct {
	cli *minio.Client
}

func NewClient(endpoint, accessKey, secretAccessKey string) (*Client, error) {
	var cli Client
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: false, // 是否使用https进行通信
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	cli.cli = minioClient
	return &cli, nil
}

// FPutObject 从本地读入文件并上传
func (c *Client) FPutObject(bucketName, filePath string) (string, error) {
	baseName := filepath.Base(filepath.ToSlash(filePath))
	UserMetadata := map[string]string{
		"origin_name": baseName,
	}
	// objectSize可设置为-1，表示不确定文件大小，但是-1会预分配比较大的内存。
	// 将文件ContentType为二进制类型，之后点击这个文件链接会自动触发下载
	uploadInfo, err := c.cli.FPutObject(context.Background(), bucketName, baseName, filePath, minio.PutObjectOptions{ContentType: "application/octet-stream", UserMetadata: UserMetadata})
	if err != nil {
		return "", lberr.Wrap(err)
	}
	log.Infof("Successfully FPutObject bytes: %v", uploadInfo)
	return uploadInfo.Key, nil
}

// UploadNetIO IO流上传
func (c *Client) UploadNetIO(bucketName, fileName string, reader io.Reader) (string, error) {
	uploadInfo, err := c.cli.PutObject(context.Background(), bucketName, fileName, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return "", lberr.Wrap(err)
	}
	log.Infof("Successfully UploadNetIO bytes: %v", uploadInfo)
	return uploadInfo.Key, nil
}

func (c *Client) Download(bucketName, objectName string, reader io.Reader) error {
	object, err := c.cli.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return lberr.Wrap(err)
	}
	defer object.Close()

	_, err = io.ReadAll(reader)
	if err != nil {
		return lberr.Wrap(err)
	}
	return nil
}

func (c *Client) FGetObject(bucketName, objectName, filePath string) error {
	// 整个文件下载和保存到指定目录，适合文件下载，如下载pdf文件
	err := c.cli.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return lberr.Wrap(err)
	}
	return nil
}

func (c *Client) PreSignedPutObject(bucketName, objectName string) (string, error) {
	u, err := c.cli.PresignedPutObject(context.Background(), bucketName, objectName, time.Minute)
	if err != nil {
		return "", lberr.Wrap(err)
	}
	return u.String(), nil
}

func (c *Client) PreSignedGetObject(bucketName, objectName string) (string, error) {
	u, err := c.cli.PresignedGetObject(context.Background(), bucketName, objectName, time.Minute, url.Values{})
	if err != nil {
		return "", lberr.Wrap(err)
	}
	return u.String(), nil
}
