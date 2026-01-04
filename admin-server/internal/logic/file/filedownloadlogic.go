// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package file

import (
	"context"
	"fmt"
	"os"
	"strings"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadReq) (resp *types.FileDownloadResp, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文件ID不能为空")
	}

	fileRepo := repository.NewFileRepository(l.svcCtx.Repository)
	file, err := fileRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeNotFound, "文件不存在", err)
	}

	// 检查文件是否存在
	fileSystemPath := file.Path
	if strings.HasPrefix(file.Path, "/uploads/") {
		fileSystemPath = "." + file.Path
	} else if !strings.HasPrefix(file.Path, "./") {
		fileSystemPath = "./" + file.Path
	}

	if _, err := os.Stat(fileSystemPath); os.IsNotExist(err) {
		return nil, errs.New(errs.CodeNotFound, "文件不存在")
	}

	// 返回nginx代理路径（用于前端访问）
	// nginx配置中，/files/uploads/ 路径会代理到后端的 /api/v1/uploads/
	// 文件直接访问路径格式：/files/uploads/xxx（如果文件路径是 /uploads/xxx）
	// 文件下载接口路径格式：/files/download?id=xxx（通过下载接口获取文件）

	var proxyPath string
	if strings.HasPrefix(file.Path, "/uploads/") {
		// 如果是 /uploads/xxx，转换为 /files/uploads/xxx（nginx代理路径）
		fileName := strings.TrimPrefix(file.Path, "/uploads/")
		proxyPath = fmt.Sprintf("/files/uploads/%s", fileName)
	} else if strings.HasPrefix(file.Path, "/files/") {
		// 如果已经是 /files/ 开头，直接使用
		proxyPath = file.Path
	} else {
		// 其他情况，使用下载接口路径（通过nginx代理）
		proxyPath = fmt.Sprintf("/files/download?id=%d", file.Id)
	}

	return &types.FileDownloadResp{
		Url: proxyPath,
	}, nil
}
