// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FileUpload admin_file 表物理属于 iam-rpc，但文件字节继续由 gateway 直接读写共享
// uploads 卷（和 task/content 拆分时的既有做法一致），只有元数据的按 MD5 去重登记走
// IamRPC.FileRegister（原逻辑迁移自 Domain.System.File.FindByName/Create 那两步）。
func (l *FileUploadLogic) FileUpload(r *http.Request) (resp *types.FileUploadResp, err error) {
	err = r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadRequest, "解析上传文件失败", err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadRequest, "获取上传文件失败", err)
	}
	defer file.Close()

	if err := os.MkdirAll(consts.UploadDir, 0755); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建上传目录失败", err)
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "计算文件MD5失败", err)
	}
	md5Hash := fmt.Sprintf("%x", hash.Sum(nil))

	// 重置文件指针，以便后续读取；Seek 失败会导致下面的 io.Copy 从错误的偏移量开始，
	// 落盘内容和用于命名的 MD5 对不上，必须当错误处理而不是静默忽略。
	if _, err := file.Seek(0, 0); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "重置文件指针失败", err)
	}

	ext := filepath.Ext(header.Filename)
	fileName := md5Hash + ext
	baseURL := l.getStorageBaseURL()

	fileSystemPath := filepath.Join(consts.UploadDir, fileName)
	proxyPath := fmt.Sprintf("%s/%s", consts.PathFileUploads, fileName)

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType([]byte(ext))
	}
	extNoDot := strings.TrimPrefix(ext, ".")

	// 文件已存在（按 MD5 命中）时不需要再写盘，RegisterFile 内部会直接返回已有记录。
	if _, statErr := os.Stat(fileSystemPath); statErr != nil {
		dst, err := os.Create(fileSystemPath)
		if err != nil {
			return nil, errs.Wrap(errs.CodeInternalError, "创建文件失败", err)
		}
		if _, err = io.Copy(dst, file); err != nil {
			dst.Close()
			return nil, errs.Wrap(errs.CodeInternalError, "保存文件失败", err)
		}
		dst.Close()
	}

	fileInfo, err := os.Stat(fileSystemPath)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "获取文件信息失败", err)
	}

	rpcResp, err := l.svcCtx.IamRPC.FileRegister(l.ctx, &iamclient.FileRegisterRequest{
		Name:         fileName,
		OriginalName: header.Filename,
		Path:         proxyPath,
		BaseUrl:      baseURL,
		Size:         uint64(fileInfo.Size()),
		MimeType:     mimeType,
		Ext:          extNoDot,
	})
	if err != nil {
		os.Remove(fileSystemPath)
		return nil, errs.WrapGRPCError("保存文件记录失败", err)
	}

	fullURL := rpcResp.Path
	if rpcResp.BaseUrl != "" && (strings.HasPrefix(rpcResp.BaseUrl, "http://") || strings.HasPrefix(rpcResp.BaseUrl, "https://")) {
		fullURL = rpcResp.BaseUrl + rpcResp.Path
	}

	return &types.FileUploadResp{
		Id:           rpcResp.Id,
		Name:         rpcResp.Name,
		OriginalName: rpcResp.OriginalName,
		Path:         rpcResp.Path,
		BaseUrl:      rpcResp.BaseUrl,
		Url:          fullURL,
		Size:         rpcResp.Size,
		MimeType:     rpcResp.MimeType,
		Ext:          rpcResp.Ext,
	}, nil
}

// getStorageBaseURL 从字典中获取存储 baseURL（字典物理属于 iam-rpc，改成走 DictGet RPC）
func (l *FileUploadLogic) getStorageBaseURL() string {
	rpcResp, err := l.svcCtx.IamRPC.DictGet(l.ctx, &iamclient.DictGetRequest{Code: consts.DictCodeStorageBaseURL})
	if err != nil || len(rpcResp.Items) == 0 {
		l.Errorf("获取存储配置字典失败: %v", err)
		return ""
	}

	baseURL := rpcResp.Items[0].Value
	if baseURL == "" {
		return ""
	}
	return strings.TrimSuffix(baseURL, "/")
}
