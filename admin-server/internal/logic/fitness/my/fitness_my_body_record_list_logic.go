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

type FitnessMyBodyRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMyBodyRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyBodyRecordListLogic {
	return &FitnessMyBodyRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMyBodyRecordListLogic) FitnessMyBodyRecordList(req *types.FitnessMyBodyRecordListReq) (resp *types.FitnessBodyRecordListResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.IamRPC.FitnessMyBodyRecordList(l.ctx, &iamclient.FitnessMyBodyRecordListRequest{UserId: user.UserID, Days: req.Days})
	if err != nil {
		return nil, errs.WrapGRPCError("查询身体数据失败", err)
	}
	return &types.FitnessBodyRecordListResp{List: fitnessconv.BodyItems(rpcResp.List)}, nil
}
