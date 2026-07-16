package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileGetMetaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileGetMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileGetMetaLogic {
	return &FileGetMetaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FileGetMeta 供 gateway FileDownload：gateway 拿到 path 后自己转换成 uploads 卷内的
// 文件系统路径并读取字节，这里只返回 admin_file 的元数据。
func (l *FileGetMetaLogic) FileGetMeta(in *iam.FileGetMetaRequest) (*iam.FileGetMetaResponse, error) {
	file, err := l.svcCtx.Domain.System.File.FindByID(l.ctx, in.Id)
	if err != nil {
		return &iam.FileGetMetaResponse{Exists: false}, nil
	}

	originalName := file.OriginalName
	if originalName == "" {
		originalName = file.Name
	}

	return &iam.FileGetMetaResponse{
		Exists:       true,
		OriginalName: originalName,
		MimeType:     file.MimeType.String,
		Path:         file.Path,
	}, nil
}
