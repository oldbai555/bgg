// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package my

import (
	"context"

	"postapocgame/admin-server/internal/logic/fitness/fitnessconv"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyDayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMyDayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyDayLogic {
	return &FitnessMyDayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMyDayLogic) FitnessMyDay(req *types.FitnessMyDayReq) (resp *types.FitnessMyDayResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.IamRPC.FitnessMyDay(l.ctx, &iamclient.FitnessMyDayRequest{UserId: user.UserID, Date: req.Date})
	if err != nil {
		return nil, errs.WrapGRPCError("查询当日计划失败", err)
	}
	return &types.FitnessMyDayResp{
		Date:         rpcResp.Date,
		Today:        rpcResp.Today,
		Weekday:      rpcResp.Weekday,
		CanCheckin:   rpcResp.CanCheckin,
		HasPlan:      rpcResp.HasPlan,
		TemplateId:   rpcResp.TemplateId,
		TemplateName: rpcResp.TemplateName,
		IsOverride:   rpcResp.IsOverride,
		Content:      fitnessconv.DayContent(rpcResp.Content),
		Tips:         fitnessconv.TipItems(rpcResp.Tips),
		Checkin:      fitnessconv.CheckinItem(rpcResp.Checkin),
		Suggestion:   fitnessconv.SuggestionItem(rpcResp.Suggestion),
	}, nil
}
