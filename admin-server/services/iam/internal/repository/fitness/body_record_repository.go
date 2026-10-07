package fitness

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const bodyRecordColumns = "id, user_id, record_date, weight_kg, waist_cm, created_at, updated_at, deleted_at"

type BodyRecordRepository interface {
	FindByUserDate(ctx context.Context, userID uint64, date string) (*fitnessmodel.FitnessBodyRecord, error)
	// ListByUsersRange 按使用者、日期正序；start 为空表示不限开始日期，userIDs 为 nil 表示全部使用者
	ListByUsersRange(ctx context.Context, userIDs []uint64, start, end string) ([]fitnessmodel.FitnessBodyRecord, error)
	Create(ctx context.Context, b *fitnessmodel.FitnessBodyRecord) error
	Update(ctx context.Context, b *fitnessmodel.FitnessBodyRecord) error
}

type bodyRecordRepository struct {
	model fitnessmodel.FitnessBodyRecordModel
	conn  sqlx.SqlConn
}

func NewBodyRecordRepository(repo *repository.Repository) BodyRecordRepository {
	return &bodyRecordRepository{model: repo.FitnessBodyRecordModel, conn: repo.DB}
}

func (r *bodyRecordRepository) FindByUserDate(ctx context.Context, userID uint64, date string) (*fitnessmodel.FitnessBodyRecord, error) {
	return notFoundAsNil(r.model.FindOneByUserIdRecordDateDeletedAt(ctx, userID, date, 0))
}

func (r *bodyRecordRepository) ListByUsersRange(ctx context.Context, userIDs []uint64, start, end string) ([]fitnessmodel.FitnessBodyRecord, error) {
	if userIDs != nil && len(userIDs) == 0 {
		return nil, nil
	}
	conditions := sq.And{sq.Eq{"deleted_at": 0}}
	if userIDs != nil {
		conditions = append(conditions, sq.Eq{"user_id": userIDs})
	}
	if start != "" {
		conditions = append(conditions, sq.GtOrEq{"record_date": start})
	}
	if end != "" {
		conditions = append(conditions, sq.LtOrEq{"record_date": end})
	}
	query, args, err := sq.Select(bodyRecordColumns).From("`fitness_body_record`").Where(conditions).
		OrderBy("user_id ASC", "record_date ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []fitnessmodel.FitnessBodyRecord
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *bodyRecordRepository) Create(ctx context.Context, b *fitnessmodel.FitnessBodyRecord) error {
	_, err := r.model.Insert(ctx, b)
	return err
}

func (r *bodyRecordRepository) Update(ctx context.Context, b *fitnessmodel.FitnessBodyRecord) error {
	return r.model.Update(ctx, b)
}
