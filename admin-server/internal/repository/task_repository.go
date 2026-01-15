package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"
)

// TaskRepository 任务仓储接口
// 说明：对 admin_task 表的常用访问方法统一从这里封装，供 Logic 和调度器使用
type TaskRepository interface {
	// Create 创建任务
	Create(ctx context.Context, task *model.AdminTask) (uint64, error)
	// FindOne 根据主键ID查询任务
	FindOne(ctx context.Context, id uint64) (*model.AdminTask, error)
	// FindPage 分页查询任务列表
	FindPage(ctx context.Context, page, pageSize int64, filters *TaskQueryFilter) ([]model.AdminTask, int64, error)
	// UpdateStatus 更新任务状态（只更新 status、started_at、finished_at、updated_at）
	UpdateStatus(ctx context.Context, id uint64, status int64, startedAt, finishedAt int64) error
	// UpdateResult 更新任务结果（result、error_message、status、finished_at、updated_at）
	UpdateResult(ctx context.Context, id uint64, status int64, result, errorMessage string, finishedAt int64) error
	// FindRecent 查询最近的任务（按创建时间倒序）
	FindRecent(ctx context.Context, limit int64, userId uint64) ([]model.AdminTask, error)
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
	model model.AdminTaskModel
	repo  *Repository
}

// NewTaskRepository 创建任务仓储实现
func NewTaskRepository(repo *Repository) TaskRepository {
	return &taskRepository{
		model: repo.AdminTaskModel,
		repo:  repo,
	}
}

// Create 创建任务
func (r *taskRepository) Create(ctx context.Context, task *model.AdminTask) (uint64, error) {
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
func (r *taskRepository) FindOne(ctx context.Context, id uint64) (*model.AdminTask, error) {
	task, err := r.model.FindOne(ctx, id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errs.Wrap(errs.CodeNotFound, "任务不存在", err)
		}
		return nil, errs.Wrap(errs.CodeBadDB, "查询任务失败", err)
	}
	return task, nil
}

// FindPage 分页查询任务列表
func (r *taskRepository) FindPage(ctx context.Context, page, pageSize int64, filters *TaskQueryFilter) ([]model.AdminTask, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 基础条件：未删除
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
	}

	// 任务名称模糊查询
	if filters != nil && filters.Name != "" {
		conditions = append(conditions, sq.Like{"name": "%" + filters.Name + "%"})
	}

	// 任务类型筛选（>0 才追加条件，对应字典值）
	if filters != nil && filters.Type > 0 {
		conditions = append(conditions, sq.Eq{"type": filters.Type})
	}

	// 执行类型筛选
	if filters != nil && filters.ExecutionType > 0 {
		conditions = append(conditions, sq.Eq{"execution_type": filters.ExecutionType})
	}

	// 状态筛选
	if filters != nil && filters.Status > 0 {
		conditions = append(conditions, sq.Eq{"status": filters.Status})
	}

	// 用户筛选
	if filters != nil && filters.UserId > 0 {
		conditions = append(conditions, sq.Eq{"user_id": filters.UserId})
	}

	// 创建时间范围
	if filters != nil {
		if filters.StartTime > 0 {
			conditions = append(conditions, sq.GtOrEq{"created_at": filters.StartTime})
		}
		if filters.EndTime > 0 {
			conditions = append(conditions, sq.LtOrEq{"created_at": filters.EndTime})
		}
	}

	// 查询列表
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

	var list []model.AdminTask
	if err := r.repo.DB.QueryRowsCtx(ctx, &list, sqlStr, args...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "任务列表查询失败", err)
	}

	// 查询总数
	countSql, countArgs, err := sq.Select("COUNT(*)").
		From("`admin_task`").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "任务总数SQL生成有误", err)
	}

	var total int64
	if err := r.repo.DB.QueryRowCtx(ctx, &total, countSql, countArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "任务总数查询失败", err)
	}

	return list, total, nil
}

// UpdateStatus 更新任务状态
func (r *taskRepository) UpdateStatus(ctx context.Context, id uint64, status int64, startedAt, finishedAt int64) error {
	builder := sq.Update("`admin_task`").
		Set("`status`", status).
		Set("`updated_at`", sq.Expr("UNIX_TIMESTAMP()")).
		Where(sq.Eq{"id": id, "deleted_at": 0})

	// 根据状态设置开始/结束时间
	if status == consts.TaskStatusRunning && startedAt > 0 {
		builder = builder.Set("`started_at`", startedAt)
	}
	if (status == consts.TaskStatusCompleted || status == consts.TaskStatusFailed) && finishedAt > 0 {
		builder = builder.Set("`finished_at`", finishedAt)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务状态SQL生成有误", err)
	}

	_, err = r.repo.DB.ExecCtx(ctx, sqlStr, args...)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务状态失败", err)
	}

	return nil
}

// UpdateResult 更新任务结果
func (r *taskRepository) UpdateResult(ctx context.Context, id uint64, status int64, result, errorMessage string, finishedAt int64) error {
	builder := sq.Update("`admin_task`").
		Set("`status`", status).
		Set("`result`", result).
		Set("`error_message`", errorMessage).
		Set("`updated_at`", sq.Expr("UNIX_TIMESTAMP()")).
		Where(sq.Eq{"id": id, "deleted_at": 0})

	if finishedAt > 0 {
		builder = builder.Set("`finished_at`", finishedAt)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务结果SQL生成有误", err)
	}

	_, err = r.repo.DB.ExecCtx(ctx, sqlStr, args...)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新任务结果失败", err)
	}

	return nil
}

// FindRecent 查询最近的任务
func (r *taskRepository) FindRecent(ctx context.Context, limit int64, userId uint64) ([]model.AdminTask, error) {
	if limit <= 0 {
		limit = 10
	}

	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
	}

	// 如果指定用户，则只查该用户的任务
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

	var list []model.AdminTask
	if err := r.repo.DB.QueryRowsCtx(ctx, &list, sqlStr, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "最近任务查询失败", err)
	}

	return list, nil
}
