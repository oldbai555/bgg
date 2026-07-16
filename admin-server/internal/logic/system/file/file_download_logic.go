// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"
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

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FileDownloadInfo 文件下载信息
type FileDownloadInfo struct {
	OriginalName string
	MimeType     string
	Size         uint64
}

// FileDownload 元数据查询走 IamRPC.FileGetMeta（admin_file 表物理属于 iam-rpc），
// 文件字节仍由 gateway 直接从共享 uploads 卷读取（和 FileUpload 落盘的做法对应）。
func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadReq) (fileInfo *FileDownloadInfo, filePath string, err error) {
	if req == nil || req.Id == 0 {
		return nil, "", errs.New(errs.CodeBadRequest, "文件ID不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.FileGetMeta(l.ctx, &iamclient.FileGetMetaRequest{Id: req.Id})
	if err != nil {
		return nil, "", errs.WrapGRPCError("查询文件记录失败", err)
	}
	if !rpcResp.Exists {
		return nil, "", errs.New(errs.CodeNotFound, "文件不存在")
	}

	if !strings.HasPrefix(rpcResp.Path, consts.PathFileUploads+"/") {
		return nil, "", errs.New(errs.CodeBadRequest, "文件路径格式不正确")
	}

	// 提取文件名并转换为文件系统路径
	fileName := strings.TrimPrefix(rpcResp.Path, consts.PathFileUploads+"/")
	fileSystemPath := filepath.Join(consts.UploadDir, fileName)

	// 检查文件是否存在
	fileStat, err := os.Stat(fileSystemPath)
	if os.IsNotExist(err) {
		l.Errorf("文件不存在，文件系统路径: %s, 访问路径: %s", fileSystemPath, rpcResp.Path)
		return nil, "", errs.New(errs.CodeNotFound, "文件不存在")
	}
	if err != nil {
		return nil, "", errs.Wrap(errs.CodeInternalError, "获取文件信息失败", err)
	}

	// 确定文件名（优先使用原始文件名）
	originalName := rpcResp.OriginalName

	// 返回文件信息和文件系统路径
	fileInfo = &FileDownloadInfo{
		OriginalName: originalName,
		MimeType:     rpcResp.MimeType,
		Size:         uint64(fileStat.Size()),
	}

	return fileInfo, fileSystemPath, nil
}
