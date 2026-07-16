package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUpdateLogic {
	return &FileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// File
func (l *FileUpdateLogic) FileUpdate(in *iam.FileUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文件ID不能为空"))
	}

	file, err := l.svcCtx.Domain.System.File.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询文件失败", err))
	}

	if in.Name != "" {
		file.Name = in.Name
	}
	if in.Status == 0 || in.Status == 1 {
		file.Status = in.Status
	}

	if err := l.svcCtx.Domain.System.File.Update(l.ctx, file); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新文件记录失败", err))
	}
	return &iam.Empty{}, nil
}
