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

	"github.com/zeromicro/go-zero/core/logx"
	systemrepo "postapocgame/admin-server/internal/repository/system"
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

// FileDownload 返回文件信息和文件系统路径
func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadReq) (fileInfo *FileDownloadInfo, filePath string, err error) {
	if req == nil || req.Id == 0 {
		return nil, "", errs.New(errs.CodeBadRequest, "文件ID不能为空")
	}

	fileRepo := systemrepo.NewFileRepository(l.svcCtx.Repository)
	file, err := fileRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, "", errs.Wrap(errs.CodeNotFound, "文件不存在", err)
	}

	if !strings.HasPrefix(file.Path, consts.PathFileUploads+"/") {
		return nil, "", errs.New(errs.CodeBadRequest, "文件路径格式不正确")
	}

	// 提取文件名并转换为文件系统路径
	fileName := strings.TrimPrefix(file.Path, consts.PathFileUploads+"/")
	fileSystemPath := filepath.Join(consts.UploadDir, fileName)

	// 检查文件是否存在
	fileStat, err := os.Stat(fileSystemPath)
	if os.IsNotExist(err) {
		l.Errorf("文件不存在，文件系统路径: %s, 访问路径: %s", fileSystemPath, file.Path)
		return nil, "", errs.New(errs.CodeNotFound, "文件不存在")
	}
	if err != nil {
		return nil, "", errs.Wrap(errs.CodeInternalError, "获取文件信息失败", err)
	}

	// 确定文件名（优先使用原始文件名）
	originalName := file.OriginalName
	if originalName == "" {
		originalName = file.Name
	}

	// 返回文件信息和文件系统路径
	fileInfo = &FileDownloadInfo{
		OriginalName: originalName,
		MimeType:     file.MimeType.String,
		Size:         uint64(fileStat.Size()),
	}

	return fileInfo, fileSystemPath, nil
}
