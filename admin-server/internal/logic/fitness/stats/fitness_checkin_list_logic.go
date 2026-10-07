// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package stats

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessCheckinListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessCheckinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessCheckinListLogic {
	return &FitnessCheckinListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessCheckinListLogic) FitnessCheckinList(req *types.FitnessCheckinListReq) (resp *types.FitnessCheckinListResp, err error) {
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)
	rpcResp, err := l.svcCtx.IamRPC.FitnessCheckinList(l.ctx, &iamclient.FitnessCheckinListRequest{
		Page:      req.Page,
		PageSize:  req.PageSize,
		UserId:    req.UserId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询打卡明细失败", err)
	}
	items := make([]types.FitnessCheckinRow, 0, len(rpcResp.List))
	for _, r := range rpcResp.List {
		items = append(items, types.FitnessCheckinRow{
			Id:             r.Id,
			Date:           r.Date,
			UserId:         r.UserId,
			Nickname:       r.Nickname,
			TemplateName:   r.TemplateName,
			TrainingType:   r.TrainingType,
			TrainingStatus: r.TrainingStatus,
			ExerciseDone:   r.ExerciseDone,
			ExerciseTotal:  r.ExerciseTotal,
			MealOnPlan:     r.MealOnPlan,
			MealOther:      r.MealOther,
			MealSkipped:    r.MealSkipped,
			Steps:          r.Steps,
			StepGoal:       r.StepGoal,
			WaterCups:      r.WaterCups,
			Score:          r.Score,
			UpdatedAt:      r.UpdatedAt,
		})
	}
	return &types.FitnessCheckinListResp{Total: rpcResp.Total, List: items}, nil
}
