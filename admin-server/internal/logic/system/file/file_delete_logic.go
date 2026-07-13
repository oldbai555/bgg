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

type FileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDeleteLogic) FileDelete(req *types.FileDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "文件ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.FileDelete(l.ctx, &iamclient.FileDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除文件记录失败", err)
	}
	return nil
}
