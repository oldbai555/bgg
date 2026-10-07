// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package stats

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessStatsOverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessStatsOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessStatsOverviewLogic {
	return &FitnessStatsOverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessStatsOverviewLogic) FitnessStatsOverview(req *types.FitnessStatsOverviewReq) (resp *types.FitnessStatsOverviewResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.FitnessStatsOverview(l.ctx, &iamclient.FitnessStatsOverviewRequest{Date: req.Date})
	if err != nil {
		return nil, errs.WrapGRPCError("查询打卡概览失败", err)
	}
	ranking := make([]types.FitnessRankItem, 0, len(rpcResp.StreakRanking))
	for _, r := range rpcResp.StreakRanking {
		ranking = append(ranking, types.FitnessRankItem{UserId: r.UserId, Nickname: r.Nickname, Avatar: r.Avatar, StreakDays: r.StreakDays, WeekRate: r.WeekRate})
	}
	daily := make([]types.FitnessDailyCount, 0, len(rpcResp.DailyCounts))
	for _, d := range rpcResp.DailyCounts {
		daily = append(daily, types.FitnessDailyCount{Date: d.Date, Count: d.Count})
	}
	return &types.FitnessStatsOverviewResp{
		Date:              rpcResp.Date,
		MemberCount:       rpcResp.MemberCount,
		TodayCheckinCount: rpcResp.TodayCheckinCount,
		WeekAvgRate:       rpcResp.WeekAvgRate,
		StreakRanking:     ranking,
		DailyCounts:       daily,
	}, nil
}
