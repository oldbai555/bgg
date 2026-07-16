// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"net/http"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login iam 域已拆分成独立服务，业务逻辑（用户查找/密码校验/生成 token/登录日志/未读公告
// 通知）整段搬进了 services/iam/internal/logic/loginlogic.go，这里只剩：从 httpReq 取出
// IP/UA 显式传给 iam-rpc（gateway 侧已经不持有 *http.Request 之外的用户上下文）+ 映射响应。
func (l *LoginLogic) Login(req *types.LoginReq, httpReq *http.Request) (resp *types.TokenPair, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	clientIP := ""
	userAgent := ""
	if httpReq != nil {
		clientIP = getClientIPFromRequest(httpReq)
		userAgent = httpReq.UserAgent()
	}

	rpcResp, err := l.svcCtx.IamRPC.Login(l.ctx, &iamclient.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		ClientIp:  clientIP,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("登录失败", err)
	}

	return &types.TokenPair{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
	}, nil
}

// getClientIPFromRequest 获取客户端 IP 地址
func getClientIPFromRequest(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
