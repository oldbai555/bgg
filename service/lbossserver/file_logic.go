/**
 * @Author: zjj
 * @Date: 2025/2/26
 * @Desc:
**/

package lbossserver

import (
	"github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/bgg/service/lboss"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/micro/uctx"
	"io"
	"net/url"
	"path"
	"strings"
)

// 保存到数据库
func saveFileToOrm(ctx uctx.IUCtx, file *lboss.ModelFile) error {
	if file == nil {
		log.Errorf("file is nil")
		return nil
	}

	err := OrmFile.NewBaseScope().Create(ctx, &file)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func uploadToOss(nCtx uctx.IUCtx, fileSize int64, filePath, fileName, reFileName string, reader io.Reader) (string, error) {
	p, err := minIoSDK.UploadNetIO(constant.BucketByPublic, path.Join(filePath, reFileName), reader)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", lberr.Wrap(err)
	}
	err = saveFileToOrm(nCtx, &lboss.ModelFile{
		Size:   fileSize,
		Name:   fileName,
		Rename: reFileName,
		Path:   strings.TrimLeft(p, "/"),
		Bucket: constant.BucketByPublic,
		Domain: constant.DOMAIN,
	})
	if err != nil {
		return "", lberr.Wrap(err)
	}
	result, err := url.JoinPath(constant.DOMAIN, p)
	if err != nil {
		return "", lberr.Wrap(err)
	}
	return result, nil
}
