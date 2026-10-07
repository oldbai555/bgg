// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package stats

import (
	"context"

	"postapocgame/admin-server/internal/logic/fitness/fitnessconv"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMemberDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMemberDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMemberDetailLogic {
	return &FitnessMemberDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMemberDetailLogic) FitnessMemberDetail(req *types.FitnessMemberDetailReq) (resp *types.FitnessMemberDetailResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.FitnessMemberDetail(l.ctx, &iamclient.FitnessMemberDetailRequest{
		UserId:    req.UserId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询使用者打卡详情失败", err)
	}
	calendar := make([]types.FitnessCalendarDay, 0, len(rpcResp.Calendar))
	for _, d := range rpcResp.Calendar {
		calendar = append(calendar, types.FitnessCalendarDay{Date: d.Date, Counted: d.Counted, Checked: d.Checked, Score: d.Score, TrainingType: d.TrainingType})
	}
	return &types.FitnessMemberDetailResp{
		Member:      fitnessconv.MemberItem(rpcResp.Member),
		RangeRate:   rpcResp.RangeRate,
		Calendar:    calendar,
		BodyRecords: fitnessconv.BodyItems(rpcResp.BodyRecords),
		Switches:    fitnessconv.SwitchItems(rpcResp.Switches),
	}, nil
}
