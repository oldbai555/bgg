// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login_log

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogListLogic {
	return &LoginLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogListLogic) LoginLogList(req *types.LoginLogListReq) (resp *types.LoginLogListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	page, pageSize := logicutil.NormalizePage(int64(req.Page), int64(req.PageSize), 20, 100)

	rpcResp, err := l.svcCtx.IamRPC.LoginLogList(l.ctx, &iamclient.LoginLogListRequest{
		Page:      page,
		PageSize:  pageSize,
		UserId:    req.UserId,
		Username:  req.Username,
		Status:    int64(req.Status),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询登录日志列表失败", err)
	}

	items := make([]types.LoginLogItem, 0, len(rpcResp.List))
	for _, log := range rpcResp.List {
		items = append(items, types.LoginLogItem{
			Id:        log.Id,
			UserId:    log.UserId,
			Username:  log.Username,
			IpAddress: log.IpAddress,
			Location:  log.Location,
			Browser:   log.Browser,
			Os:        log.Os,
			UserAgent: log.UserAgent,
			Status:    int(log.Status),
			Message:   log.Message,
			LoginAt:   log.LoginAt,
			LogoutAt:  log.LogoutAt,
			CreatedAt: log.CreatedAt,
		})
	}

	return &types.LoginLogListResp{
		List:     items,
		Total:    rpcResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
