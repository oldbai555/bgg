package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/services/iam/iam"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileRegisterLogic {
	return &FileRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FileRegister 供 gateway FileUpload：文件字节由 gateway 直接写共享 uploads 卷（和
// task/content 拆分时的既有做法一致），这里只做 admin_file 元数据的按 MD5 去重登记，
// 原逻辑迁移自 internal/logic/system/file/file_upload_logic.go。
func (l *FileRegisterLogic) FileRegister(in *iam.FileRegisterRequest) (*iam.FileRegisterResponse, error) {
	if existing, err := l.svcCtx.Domain.System.File.FindByName(l.ctx, in.Name); err == nil && existing != nil {
		return &iam.FileRegisterResponse{
			Id:             existing.Id,
			Name:           existing.Name,
			OriginalName:   existing.OriginalName,
			Path:           existing.Path,
			BaseUrl:        existing.BaseUrl,
			Size:           existing.Size,
			MimeType:       existing.MimeType.String,
			Ext:            existing.Ext.String,
			AlreadyExisted: true,
		}, nil
	}

	fileModel := systemmodel.AdminFile{
		Name:         in.Name,
		OriginalName: in.OriginalName,
		Path:         in.Path,
		BaseUrl:      in.BaseUrl,
		Size:         in.Size,
		MimeType:     sql.NullString{String: in.MimeType, Valid: in.MimeType != ""},
		Ext:          sql.NullString{String: in.Ext, Valid: in.Ext != ""},
		StorageType:  "local",
		Status:       1,
	}

	if err := l.svcCtx.Domain.System.File.Create(l.ctx, &fileModel); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &iam.FileRegisterResponse{
		Id:           fileModel.Id,
		Name:         fileModel.Name,
		OriginalName: fileModel.OriginalName,
		Path:         fileModel.Path,
		BaseUrl:      fileModel.BaseUrl,
		Size:         fileModel.Size,
		MimeType:     in.MimeType,
		Ext:          in.Ext,
	}, nil
}
