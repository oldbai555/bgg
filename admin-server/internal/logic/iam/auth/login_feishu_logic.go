package auth

import (
	"context"
	"net/http"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginFeishuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginFeishuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginFeishuLogic {
	return &LoginFeishuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// LoginFeishu 同 Login 的薄胶水模式：真正的飞书 code 换用户信息/自动建号/签发 token
// 全部在 services/iam/internal/logic/loginfeishulogic.go，这里只转发 IP/UA + 映射响应。
func (l *LoginFeishuLogic) LoginFeishu(req *types.LoginFeishuReq, httpReq *http.Request) (resp *types.TokenPair, err error) {
	if req == nil || req.Code == "" {
		return nil, errs.New(errs.CodeBadRequest, "缺少飞书授权 code")
	}

	clientIP := ""
	userAgent := ""
	if httpReq != nil {
		clientIP = getClientIPFromRequest(httpReq)
		userAgent = httpReq.UserAgent()
	}

	rpcResp, err := l.svcCtx.IamRPC.LoginFeishu(l.ctx, &iamclient.LoginFeishuRequest{
		Code:      req.Code,
		State:     req.State,
		ClientIp:  clientIP,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("飞书登录失败", err)
	}

	return &types.TokenPair{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
	}, nil
}
