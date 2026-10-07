// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"net/http"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessLoginFeishuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessLoginFeishuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessLoginFeishuLogic {
	return &FitnessLoginFeishuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessLoginFeishuLogic) FitnessLoginFeishu(req *types.FitnessLoginFeishuReq, httpReq *http.Request) (resp *types.TokenPair, err error) {
	if req == nil || req.Code == "" {
		return nil, errs.New(errs.CodeBadRequest, "缺少飞书授权 code")
	}
	scene := consts.FeishuSceneFitnessOAuth
	if req.Mode == consts.FitnessLoginModeH5 {
		scene = consts.FeishuSceneFitnessH5
	}

	clientIP, userAgent := "", ""
	if httpReq != nil {
		clientIP = logicutil.ClientIP(httpReq)
		userAgent = httpReq.UserAgent()
	}

	rpcResp, err := l.svcCtx.IamRPC.LoginFeishu(l.ctx, &iamclient.LoginFeishuRequest{
		Code:      req.Code,
		State:     req.State,
		ClientIp:  clientIP,
		UserAgent: userAgent,
		Scene:     scene,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("飞书登录失败", err)
	}
	return &types.TokenPair{AccessToken: rpcResp.AccessToken, RefreshToken: rpcResp.RefreshToken}, nil
}
