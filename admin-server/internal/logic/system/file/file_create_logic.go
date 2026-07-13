// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCreateLogic {
	return &FileCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileCreateLogic) FileCreate(req *types.FileCreateReq) error {
	if req == nil || req.Name == "" {
		return errs.New(errs.CodeBadRequest, "文件名称不能为空")
	}

	_, err := l.svcCtx.IamRPC.FileCreate(l.ctx, &iamclient.FileCreateRequest{
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("创建文件记录失败", err)
	}
	return nil
}
