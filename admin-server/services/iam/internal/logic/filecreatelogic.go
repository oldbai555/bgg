package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCreateLogic {
	return &FileCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// File
func (l *FileCreateLogic) FileCreate(in *iam.FileCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文件名称不能为空"))
	}

	status := in.Status
	if status == 0 {
		status = 1
	}

	file := systemmodel.AdminFile{
		Name:         in.Name,
		OriginalName: in.Name,
		Path:         "",
		BaseUrl:      "",
		Size:         0,
		MimeType:     sql.NullString{Valid: false},
		Ext:          sql.NullString{Valid: false},
		StorageType:  "local",
		Status:       status,
	}

	if err := l.svcCtx.Domain.System.File.Create(l.ctx, &file); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建文件记录失败", err))
	}
	return &iam.Empty{}, nil
}
