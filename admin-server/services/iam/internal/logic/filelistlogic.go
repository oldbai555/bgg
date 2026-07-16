package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileListLogic) FileList(in *iam.FileListRequest) (*iam.FileListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.System.File.FindPage(l.ctx, in.Page, in.PageSize, in.Name)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询文件列表失败", err))
	}

	items := make([]*iam.FileItem, 0, len(list))
	for _, f := range list {
		items = append(items, &iam.FileItem{
			Id:           f.Id,
			Name:         f.Name,
			OriginalName: f.OriginalName,
			Path:         f.Path,
			BaseUrl:      f.BaseUrl,
			Status:       f.Status,
			CreatedAt:    f.CreatedAt,
		})
	}

	return &iam.FileListResponse{Total: total, List: items}, nil
}
