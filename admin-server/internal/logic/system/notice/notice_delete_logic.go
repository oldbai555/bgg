// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	systemrepo "postapocgame/admin-server/internal/repository/system"
)

type NoticeDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeDeleteLogic {
	return &NoticeDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeDeleteLogic) NoticeDelete(req *types.NoticeDeleteReq) (resp *types.Response, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	noticeRepo := systemrepo.NewNoticeRepository(l.svcCtx.Repository)
	if err := noticeRepo.DeleteByID(l.ctx, req.Id); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "删除公告失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "删除成功",
	}, nil
}
