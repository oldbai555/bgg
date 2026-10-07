// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package member

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMemberAssignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMemberAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMemberAssignLogic {
	return &FitnessMemberAssignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMemberAssignLogic) FitnessMemberAssign(req *types.FitnessMemberAssignReq) (resp *types.FitnessTemplateSwitchResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.IamRPC.FitnessMemberAssign(l.ctx, &iamclient.FitnessMemberAssignRequest{
		UserId:     req.UserId,
		TemplateId: req.TemplateId,
		Reason:     req.Reason,
		OperatorId: user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("指定模板失败", err)
	}
	return &types.FitnessTemplateSwitchResp{EffectiveDate: rpcResp.EffectiveDate}, nil
}
