// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/task/taskclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskRecentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskRecentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskRecentLogic {
	return &TaskRecentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TaskRecent 薄胶水：原来"缓存优先、字典兜底、硬编码兜底"的 getRecentTaskLimit 逻辑
// 已经按 15-service-boundaries.md 第 5 节末尾的建议下沉成 task-rpc 自己的静态配置
// （services/task/etc/task.yaml 的 RecentTaskLimit 字段），不做 RPC 查询系统字典
// （该值几乎不变，跨服务查一次字典表换来的收益极小，见该节原文分析）。req.Limit=0
// 时传给 task-rpc 的也是 0，由 task-rpc 侧决定用静态配置兜底。
func (l *TaskRecentLogic) TaskRecent(req *types.TaskRecentReq) (resp *types.TaskRecentResp, err error) {
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	rpcResp, err := l.svcCtx.TaskRPC.TaskRecent(l.ctx, &taskclient.TaskRecentRequest{
		Limit:  req.Limit,
		UserId: user.UserID,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询最近任务失败", err)
	}

	list := make([]types.TaskItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, types.TaskItem{
			Id:            item.Id,
			Name:          item.Name,
			TaskType:      item.TaskType,
			ExecutionType: item.ExecutionType,
			Status:        item.Status,
			UserId:        item.UserId,
			ScheduledAt:   item.ScheduledAt,
			StartedAt:     item.StartedAt,
			FinishedAt:    item.FinishedAt,
			CreatedAt:     item.CreatedAt,
			Params:        item.Params,
			Result:        item.Result,
			ErrorMessage:  item.ErrorMessage,
		})
	}

	return &types.TaskRecentResp{
		List: list,
	}, nil
}
