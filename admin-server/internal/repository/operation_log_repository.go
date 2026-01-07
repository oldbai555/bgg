package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"
)

type OperationLogRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.AdminOperationLog, error)
	FindPage(ctx context.Context, page, pageSize int64, userId uint64, username, operationType, operationObject, method, startTime, endTime string) ([]model.AdminOperationLog, int64, error)
	Create(ctx context.Context, log *model.AdminOperationLog) error
	// 批量创建（用于异步写入）
	BatchCreate(ctx context.Context, logs []*model.AdminOperationLog) error
}

type operationLogRepository struct {
	model model.AdminOperationLogModel
	conn  sqlx.SqlConn
}

func NewOperationLogRepository(repo *Repository) OperationLogRepository {
	return &operationLogRepository{model: repo.AdminOperationLogModel, conn: repo.DB}
}

func (r *operationLogRepository) FindByID(ctx context.Context, id uint64) (*model.AdminOperationLog, error) {
	return r.model.FindOne(ctx, id)
}

func (r *operationLogRepository) FindPage(ctx context.Context, page, pageSize int64, userId uint64, username, operationType, operationObject, method, startTime, endTime string) ([]model.AdminOperationLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	conditions := sq.And{sq.Eq{"deleted_at": 0}}

	if userId > 0 {
		conditions = append(conditions, sq.Eq{"user_id": userId})
	}
	if username != "" {
		conditions = append(conditions, sq.Like{"username": "%" + username + "%"})
	}
	if operationType != "" {
		conditions = append(conditions, sq.Eq{"operation_type": operationType})
	}
	if operationObject != "" {
		conditions = append(conditions, sq.Eq{"operation_object": operationObject})
	}
	if method != "" {
		conditions = append(conditions, sq.Eq{"method": method})
	}
	if startTime != "" {
		// 解析时间字符串为时间戳
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			conditions = append(conditions, sq.GtOrEq{"created_at": t.Unix()})
		}
	}
	if endTime != "" {
		// 解析时间字符串为时间戳
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			conditions = append(conditions, sq.LtOrEq{"created_at": t.Unix()})
		}
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("admin_operation_log").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// 查询列表
	var list []model.AdminOperationLog
	listSQL, listArgs, err := sq.Select("*").
		From("admin_operation_log").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *operationLogRepository) Create(ctx context.Context, log *model.AdminOperationLog) error {
	// go-zero 生成的 Model 会自动处理 created_at 和 updated_at
	_, err := r.model.Insert(ctx, log)
	return err
}

func (r *operationLogRepository) BatchCreate(ctx context.Context, logs []*model.AdminOperationLog) error {
	if len(logs) == 0 {
		return nil
	}

	// 批量插入（使用事务或循环插入）
	for _, log := range logs {
		_, err := r.model.Insert(ctx, log)
		if err != nil {
			return err
		}
	}
	return nil
}
