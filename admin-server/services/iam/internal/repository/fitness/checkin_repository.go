package fitness

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const checkinColumns = "id, user_id, checkin_date, template_id, training_type, training_status, exercise_done, exercise_total, meals, steps, step_goal, water_cups, created_at, updated_at, deleted_at"

// CheckinFilter 后台打卡明细筛选；UserID=0 / 日期为空表示不筛选。
type CheckinFilter struct {
	UserID    uint64
	StartDate string
	EndDate   string
}

func (f CheckinFilter) where() sq.And {
	conditions := sq.And{sq.Eq{"deleted_at": 0}}
	if f.UserID > 0 {
		conditions = append(conditions, sq.Eq{"user_id": f.UserID})
	}
	if f.StartDate != "" {
		conditions = append(conditions, sq.GtOrEq{"checkin_date": f.StartDate})
	}
	if f.EndDate != "" {
		conditions = append(conditions, sq.LtOrEq{"checkin_date": f.EndDate})
	}
	return conditions
}

type CheckinRepository interface {
	FindByUserDate(ctx context.Context, userID uint64, date string) (*fitnessmodel.FitnessCheckin, error)
	// ListByUsersRange userIDs 为 nil 表示全部使用者
	ListByUsersRange(ctx context.Context, userIDs []uint64, start, end string) ([]fitnessmodel.FitnessCheckin, error)
	FindPage(ctx context.Context, filter CheckinFilter, page, pageSize int64) ([]fitnessmodel.FitnessCheckin, int64, error)
	Create(ctx context.Context, c *fitnessmodel.FitnessCheckin) error
	Update(ctx context.Context, c *fitnessmodel.FitnessCheckin) error
}

type checkinRepository struct {
	model fitnessmodel.FitnessCheckinModel
	conn  sqlx.SqlConn
}

func NewCheckinRepository(repo *repository.Repository) CheckinRepository {
	return &checkinRepository{model: repo.FitnessCheckinModel, conn: repo.DB}
}

func (r *checkinRepository) FindByUserDate(ctx context.Context, userID uint64, date string) (*fitnessmodel.FitnessCheckin, error) {
	return notFoundAsNil(r.model.FindOneByUserIdCheckinDateDeletedAt(ctx, userID, date, 0))
}

func (r *checkinRepository) ListByUsersRange(ctx context.Context, userIDs []uint64, start, end string) ([]fitnessmodel.FitnessCheckin, error) {
	if userIDs != nil && len(userIDs) == 0 {
		return nil, nil
	}
	conditions := CheckinFilter{StartDate: start, EndDate: end}.where()
	if userIDs != nil {
		conditions = append(conditions, sq.Eq{"user_id": userIDs})
	}
	query, args, err := sq.Select(checkinColumns).From("`fitness_checkin`").Where(conditions).
		OrderBy("user_id ASC", "checkin_date ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []fitnessmodel.FitnessCheckin
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *checkinRepository) FindPage(ctx context.Context, filter CheckinFilter, page, pageSize int64) ([]fitnessmodel.FitnessCheckin, int64, error) {
	var list []fitnessmodel.FitnessCheckin
	total, err := findPage(ctx, r.conn, &list, sq.Select(checkinColumns).From("`fitness_checkin`"), "`fitness_checkin`", filter.where(), page, pageSize, "checkin_date DESC, id DESC")
	return list, total, err
}

func (r *checkinRepository) Create(ctx context.Context, c *fitnessmodel.FitnessCheckin) error {
	_, err := r.model.Insert(ctx, c)
	return err
}

func (r *checkinRepository) Update(ctx context.Context, c *fitnessmodel.FitnessCheckin) error {
	return r.model.Update(ctx, c)
}
