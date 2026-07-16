// Package repository 从 internal/repository/task/task_repository.go 原样搬迁而来，唯一的
// 结构性改动是构造函数从吃单体的 *repository.Repository（一个聚合了全部 9 个业务域 Model 的
// 大句柄）改成只吃 task-rpc 自己需要的 taskmodel.AdminTaskModel + sqlx.SqlConn——task-rpc 从
// 第一天起只有 admin_task 一张表，不该也不能继续持有指向其它域的句柄。
package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/task/internal/consts"
	taskmodelpkg "postapocgame/admin-server/services/task/internal/model/task"
)

// TaskRepository 任务仓储接口
// 说明：对 admin_task 表的常用访问方法统一从这里封装，供 Logic 和调度器使用
type TaskRepository interface {
	// Create 创建任务
	Create(ctx context.Context, task *taskmodelpkg.AdminTask) (uint64, error)
	// FindOne 根据主键ID查询任务
	FindOne(ctx context.Context, id uint64) (*taskmodelpkg.AdminTask, error)
	// FindPage 分页查询任务列表
	FindPage(ctx context.Context, page, pageSize int64, filters *TaskQueryFilter) ([]taskmodelpkg.AdminTask, int64, error)
	// UpdateStatus 更新任务状态（只更新 status、started_at、finished_at、updated_at）
	UpdateStatus(ctx context.Context, id uint64, status int64, startedAt, finishedAt int64) error
	// UpdateResult 更新任务结果（result、error_message、status、finished_at、updated_at）
	UpdateResult(ctx context.Context, id uint64, status int64, result, errorMessage string, finishedAt int64) error
	// FindRecent 查询最近的任务（按创建时间倒序）
	FindRecent(ctx context.Context, limit int64, userId uint64) ([]taskmodelpkg.AdminTask, error)
	// FindPendingAsync 扫描待执行的异步任务（execution_type=2, status=1, scheduled_at=0），供调度器使用
	FindPendingAsync(ctx context.Context, limit int64) ([]taskmodelpkg.AdminTask, error)
	// FindPendingScheduled 扫描待执行的定时任务（execution_type=2, status=1, 0<scheduled_at<=now），供调度器使用
	FindPendingScheduled(ctx context.Context, limit int64, now int64) ([]taskmodelpkg.AdminTask, error)
}

// TaskQueryFilter 任务查询过滤条件
type TaskQueryFilter struct {
	Name          string
	Type          int64
	ExecutionType int64
	Status        int64
	UserId        uint64
	StartTime     int64
	EndTime       int64
}

type taskRepository struct {
	model taskmodelpkg.AdminTaskModel
	conn  sqlx.SqlConn
}

// NewTaskRepository 创建任务仓储实现
func NewTaskRepository(model taskmodelpkg.AdminTaskModel, conn sqlx.SqlConn) TaskRepository {
	return &taskRepository{model: model, conn: conn}
}

// Create 创建任务
func (r *taskRepository) Create(ctx context.Context, task *taskmodelpkg.AdminTask) (uint64, error) {
	result, err := r.model.Insert(ctx, task)
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "创建任务失败", err)
	}
	if result == nil {
		return 0, errs.Wrap(errs.CodeBadDB, "创建任务失败：返回结果为空", nil)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "获取任务ID失败", err)
	}
	return uint64(id), nil
}

// FindOne 根据ID查询任务
func (r *taskRepository) FindOne(ctx context.Context, id uint64) (*taskmodelpkg.AdminTask, error) {
	task, err := r.model.FindOne(ctx, id)
	if err != nil {
		if err == taskmodelpkg.ErrNotFound {
			return nil, errs.Wrap(errs.CodeNotFound, "任务不存在", err)
		}
		return nil, errs.Wrap(errs.CodeBadDB, "查询任务失败", err)
	}
	return task, nil
}

// FindPage 分页查询任务列表
func (r *taskRepository) FindPage(ctx context.Context, page, pageSize int64, filters *TaskQueryFilter) ([]taskmodelpkg.AdminTask, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
	}

	if filters != nil && filters.Name != "" {
		conditions = append(conditions, sq.Like{"name": "%" + filters.Name + "%"})
	}
	if filters != nil && filters.Type > 0 {
		conditions = append(conditions, sq.Eq{"type": filters.Type})
	}
	if filters != nil && filters.ExecutionType > 0 {
		conditions = append(conditions, sq.Eq{"execution_type": filters.ExecutionType})
	}
	if filters != nil && filters.Status > 0 {
		conditions = append(conditions, sq.Eq{"status": filters.Status})
	}
	if filters != nil && filters.UserId > 0 {
		conditions = append(conditions, sq.Eq{"user_id": filters.UserId})
	}
	if filters != nil {
		if filters.StartTime > 0 {
			conditions = append(conditions, sq.GtOrEq{"created_at": filters.StartTime})
		}
		if filters.EndTime > 0 {
			conditions = append(conditions, sq.LtOrEq{"created_at": filters.EndTime})
		}
	}

	sqlStr, args, err := sq.Select("*").
		From("`admin_task`").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "任务列表SQL生成有误", err)
	}

	var list []taskmodelpkg.AdminTask
	if err := r.conn.QueryRowsCtx(ctx, &list, sqlStr, args...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "任务列表查询失败", err)
	}

	countSql, countArgs, err := sq.Select("COUNT(*)").
		From("`admin_task`").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "任务总数SQL生成有误", err)
	}

	var total int64
	if err := r.conn.QueryRowCtx(ctx, &total, countSql, countArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "任务总数查询失败", err)
	}

	return list, total, nil
}

// UpdateStatus 更新任务状态
// 注意：必须走 r.model.Update（而不是裸 squirrel + r.conn.ExecCtx），否则不会失效
// FindOne 的 Redis 缓存（cache:adminTask:id:<id>），调度器改完状态后轮询读到的仍是旧缓存。
// 这是 Phase 1 阶段真实修过的一个生产 bug（见 docs/progress.md 2026-07-11 续三条目），
// 搬迁到 task-rpc 时必须保留这个写法，不能退化回裸 SQL。
func (r *taskRepository) UpdateStatus(ctx context.Context, id uint64, status int64, startedAt, finishedAt int64) error {
	task, err := r.model.FindOne(ctx, id)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务状态失败：查询任务不存在", err)
	}

	task.Status = status
	if status == consts.TaskStatusRunning && startedAt > 0 {
		task.StartedAt = startedAt
	}
	if (status == consts.TaskStatusCompleted || status == consts.TaskStatusFailed) && finishedAt > 0 {
		task.FinishedAt = finishedAt
	}

	if err := r.model.Update(ctx, task); err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务状态失败", err)
	}

	return nil
}

// UpdateResult 更新任务结果
// 同 UpdateStatus，走 r.model.Update 以正确失效缓存。
func (r *taskRepository) UpdateResult(ctx context.Context, id uint64, status int64, result, errorMessage string, finishedAt int64) error {
	task, err := r.model.FindOne(ctx, id)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务结果失败：查询任务不存在", err)
	}

	task.Status = status
	task.Result = sql.NullString{String: result, Valid: result != ""}
	task.ErrorMessage = errorMessage
	if finishedAt > 0 {
		task.FinishedAt = finishedAt
	}

	if err := r.model.Update(ctx, task); err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务结果失败", err)
	}

	return nil
}

// FindRecent 查询最近的任务
func (r *taskRepository) FindRecent(ctx context.Context, limit int64, userId uint64) ([]taskmodelpkg.AdminTask, error) {
	if limit <= 0 {
		limit = 10
	}

	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
	}

	if userId > 0 {
		conditions = append(conditions, sq.Eq{"user_id": userId})
	}

	sqlStr, args, err := sq.Select("*").
		From("`admin_task`").
		Where(conditions).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "最近任务SQL生成有误", err)
	}

	var list []taskmodelpkg.AdminTask
	if err := r.conn.QueryRowsCtx(ctx, &list, sqlStr, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "最近任务查询失败", err)
	}

	return list, nil
}

// FindPendingAsync 扫描待执行的异步任务
func (r *taskRepository) FindPendingAsync(ctx context.Context, limit int64) ([]taskmodelpkg.AdminTask, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"execution_type": consts.TaskExecutionTypeAsync},
		sq.Eq{"status": consts.TaskStatusPending},
		sq.Eq{"scheduled_at": 0},
	}

	sqlStr, args, err := sq.Select("*").
		From("`admin_task`").
		Where(conditions).
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "异步任务SQL生成失败", err)
	}

	var tasks []taskmodelpkg.AdminTask
	if err := r.conn.QueryRowsCtx(ctx, &tasks, sqlStr, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询异步任务失败", err)
	}
	return tasks, nil
}

// FindPendingScheduled 扫描待执行的定时任务
func (r *taskRepository) FindPendingScheduled(ctx context.Context, limit int64, now int64) ([]taskmodelpkg.AdminTask, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"execution_type": consts.TaskExecutionTypeAsync},
		sq.Eq{"status": consts.TaskStatusPending},
		sq.Gt{"scheduled_at": 0},
		sq.LtOrEq{"scheduled_at": now},
	}

	sqlStr, args, err := sq.Select("*").
		From("`admin_task`").
		Where(conditions).
		OrderBy("scheduled_at ASC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "定时任务SQL生成失败", err)
	}

	var tasks []taskmodelpkg.AdminTask
	if err := r.conn.QueryRowsCtx(ctx, &tasks, sqlStr, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询定时任务失败", err)
	}
	return tasks, nil
}
