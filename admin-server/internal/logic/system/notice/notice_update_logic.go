// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeUpdateLogic {
	return &NoticeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeUpdateLogic) NoticeUpdate(req *types.NoticeUpdateReq) (resp *types.Response, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	_, err = l.svcCtx.IamRPC.NoticeUpdate(l.ctx, &iamclient.NoticeUpdateRequest{
		Id:          req.Id,
		Title:       req.Title,
		Content:     req.Content,
		NoticeType:  req.NoticeType,
		Status:      req.Status,
		PublishTime: req.PublishTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("更新公告失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "更新成功",
	}, nil
}
