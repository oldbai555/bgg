// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeCreateLogic {
	return &NoticeCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeCreateLogic) NoticeCreate(req *types.NoticeCreateReq) (resp *types.Response, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err = l.svcCtx.IamRPC.NoticeCreate(l.ctx, &iamclient.NoticeCreateRequest{
		Title:          req.Title,
		Content:        req.Content,
		NoticeType:     req.NoticeType,
		Status:         req.Status,
		PublishTime:    req.PublishTime,
		OperatorUserId: user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建公告失败", err)
	}

	return &types.Response{
		Code:    0,
		Message: "创建成功",
	}, nil
}
