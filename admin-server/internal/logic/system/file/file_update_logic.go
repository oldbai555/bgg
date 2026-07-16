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

type FileUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUpdateLogic {
	return &FileUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUpdateLogic) FileUpdate(req *types.FileUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "文件ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.FileUpdate(l.ctx, &iamclient.FileUpdateRequest{
		Id:     req.Id,
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("更新文件记录失败", err)
	}
	return nil
}
