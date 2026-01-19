// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task_public

import (
	"context"
	"strconv"
	"time"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"

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

func (l *TaskRecentLogic) TaskRecent(req *types.TaskRecentReq) (resp *types.TaskRecentResp, err error) {
	// 获取当前用户（仅要求登录，不做权限校验）
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	limit := l.getRecentTaskLimit(req.Limit)

	// 查询最近的任务（只查询当前用户的任务）
	taskRepo := repository.NewTaskRepository(l.svcCtx.Repository)
	tasks, err := taskRepo.FindRecent(l.ctx, limit, user.UserID)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询最近任务失败", err)
	}

	// 转换为响应结构
	list := make([]types.TaskItem, 0, len(tasks))
	for _, task := range tasks {
		item := types.TaskItem{
			Id:            task.Id,
			Name:          task.Name,
			TaskType:      task.Type,
			ExecutionType: task.ExecutionType,
			Status:        task.Status,
			UserId:        task.UserId,
			ScheduledAt:   task.ScheduledAt,
			StartedAt:     task.StartedAt,
			FinishedAt:    task.FinishedAt,
			CreatedAt:     task.CreatedAt,
			ErrorMessage:  task.ErrorMessage,
		}

		// 处理可空字段
		if task.Params.Valid {
			item.Params = task.Params.String
		}
		if task.Result.Valid {
			item.Result = task.Result.String
		}

		list = append(list, item)
	}

	return &types.TaskRecentResp{
		List: list,
	}, nil
}

// getRecentTaskLimit 获取最近任务数量限制，优先使用业务缓存，其次字典，兜底 10
func (l *TaskRecentLogic) getRecentTaskLimit(requestLimit int64) int64 {
	if requestLimit > 0 {
		return requestLimit
	}

	const (
		cacheKey        = "task:config:recent_limit"
		defaultLimit    = int64(10)
		cacheExpireSecs = 600
	)

	var cached int64
	if err := l.svcCtx.Repository.BusinessCache.Get(l.ctx, cacheKey, &cached); err == nil && cached > 0 {
		return cached
	}

	limit := defaultLimit
	dictTypeRepo := repository.NewDictTypeRepository(l.svcCtx.Repository)
	dictType, err := dictTypeRepo.FindByCode(l.ctx, "task_config")
	if err == nil && dictType != nil {
		dictItemRepo := repository.NewDictItemRepository(l.svcCtx.Repository)
		items, err := dictItemRepo.FindByTypeID(l.ctx, dictType.Id)
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if item.Label == "最近任务数量" {
					if parsedLimit, parseErr := strconv.ParseInt(item.Value, 10, 64); parseErr == nil && parsedLimit > 0 {
						limit = parsedLimit
						break
					}
				}
			}
		}
	}

	_ = l.svcCtx.Repository.BusinessCache.Set(l.ctx, cacheKey, limit, cacheExpireSecs+int(time.Now().Unix()%60))
	return limit
}
