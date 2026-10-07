package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessStatsOverviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessStatsOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessStatsOverviewLogic {
	return &FitnessStatsOverviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessStatsOverviewLogic) FitnessStatsOverview(in *iam.FitnessStatsOverviewRequest) (*iam.FitnessStatsOverviewResponse, error) {
	ov, err := l.svcCtx.Domain.Fitness.Overview(l.ctx, in.Date)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessStatsOverviewResponse{
		Date:              ov.Date,
		MemberCount:       int64(ov.MemberCount),
		TodayCheckinCount: int64(ov.TodayCheckinCount),
		WeekAvgRate:       ov.WeekAvgRate,
	}
	for _, r := range ov.Ranking {
		resp.StreakRanking = append(resp.StreakRanking, &iam.FitnessRankItem{
			UserId: r.Member.UserId, Nickname: r.Member.Nickname, Avatar: r.Member.Avatar,
			StreakDays: int64(r.Summary.StreakDays), WeekRate: r.Summary.WeekRate,
		})
	}
	for _, d := range ov.DailyCounts {
		resp.DailyCounts = append(resp.DailyCounts, &iam.FitnessDailyCount{Date: d.Date, Count: int64(d.Count)})
	}
	return resp, nil
}
