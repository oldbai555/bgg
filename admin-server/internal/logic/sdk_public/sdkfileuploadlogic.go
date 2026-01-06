// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk_public

import (
	"context"
	"net/http"

	"postapocgame/admin-server/internal/logic/file"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkFileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkFileUploadLogic {
	return &SdkFileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkFileUploadLogic) SdkFileUpload(r *http.Request) (resp *types.SdkFileUploadResp, err error) {
	if r == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求不能为空")
	}
	fileLogic := file.NewFileUploadLogic(l.ctx, l.svcCtx)
	fileResp, err := fileLogic.FileUpload(r)
	if err != nil {
		return nil, err
	}
	return &types.SdkFileUploadResp{
		FileId: fileResp.Id,
		Url:    fileResp.Url,
		Name:   fileResp.Name,
	}, nil
}
