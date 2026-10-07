// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package my

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMySuggestionDecideLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMySuggestionDecideLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMySuggestionDecideLogic {
	return &FitnessMySuggestionDecideLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMySuggestionDecideLogic) FitnessMySuggestionDecide(req *types.FitnessMySuggestionDecideReq) error {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	_, err := l.svcCtx.IamRPC.FitnessMySuggestionDecide(l.ctx, &iamclient.FitnessMySuggestionDecideRequest{UserId: user.UserID, Id: req.Id, Accept: req.Accept})
	if err != nil {
		return errs.WrapGRPCError("处理模板建议失败", err)
	}
	return nil
}
