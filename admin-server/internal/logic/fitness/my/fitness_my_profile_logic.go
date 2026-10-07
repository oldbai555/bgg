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

type FitnessMyProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMyProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyProfileLogic {
	return &FitnessMyProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMyProfileLogic) FitnessMyProfile() (resp *types.FitnessMyProfileResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.IamRPC.FitnessMyProfile(l.ctx, &iamclient.FitnessMyProfileRequest{UserId: user.UserID})
	if err != nil {
		return nil, errs.WrapGRPCError("查询个人信息失败", err)
	}
	templates := make([]types.FitnessTemplateBrief, 0, len(rpcResp.Templates))
	for _, t := range rpcResp.Templates {
		templates = append(templates, *fitnessconv.TemplateBrief(t))
	}
	resp = &types.FitnessMyProfileResp{
		Today:              rpcResp.Today,
		Nickname:           rpcResp.Nickname,
		Avatar:             rpcResp.Avatar,
		HeightCm:           rpcResp.HeightCm,
		TargetWeightKg:     rpcResp.TargetWeightKg,
		TrainingPreference: rpcResp.TrainingPreference,
		LatestWeightKg:     rpcResp.LatestWeightKg,
		StreakDays:         rpcResp.StreakDays,
		WeekRate:           rpcResp.WeekRate,
		CurrentTemplate:    fitnessconv.TemplateBrief(rpcResp.CurrentTemplate),
		Suggestion:         fitnessconv.SuggestionItem(rpcResp.Suggestion),
		Templates:          templates,
	}
	if p := rpcResp.PendingSwitch; p != nil {
		resp.PendingSwitch = &types.FitnessPendingSwitch{TemplateId: p.TemplateId, TemplateName: p.TemplateName, EffectiveDate: p.EffectiveDate}
	}
	return resp, nil
}
