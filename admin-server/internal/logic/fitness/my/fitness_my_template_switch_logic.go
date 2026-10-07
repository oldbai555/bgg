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

type FitnessMyTemplateSwitchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMyTemplateSwitchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyTemplateSwitchLogic {
	return &FitnessMyTemplateSwitchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMyTemplateSwitchLogic) FitnessMyTemplateSwitch(req *types.FitnessMyTemplateSwitchReq) (resp *types.FitnessTemplateSwitchResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.IamRPC.FitnessMyTemplateSwitch(l.ctx, &iamclient.FitnessMyTemplateSwitchRequest{UserId: user.UserID, TemplateId: req.TemplateId})
	if err != nil {
		return nil, errs.WrapGRPCError("切换模板失败", err)
	}
	return &types.FitnessTemplateSwitchResp{EffectiveDate: rpcResp.EffectiveDate}, nil
}
