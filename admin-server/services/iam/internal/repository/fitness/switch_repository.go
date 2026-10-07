package fitness

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const switchColumns = "id, user_id, from_template_id, to_template_id, effective_date, source, reason, operator_id, created_at, updated_at, deleted_at"

type SwitchRepository interface {
	// ListByUsers 指定使用者的全部切换记录，按生效日期、ID 正序
	ListByUsers(ctx context.Context, userIDs []uint64) ([]fitnessmodel.FitnessTemplateSwitch, error)
	FindPage(ctx context.Context, userID uint64, page, pageSize int64) ([]fitnessmodel.FitnessTemplateSwitch, int64, error)
	Create(ctx context.Context, s *fitnessmodel.FitnessTemplateSwitch) error
	DeleteByID(ctx context.Context, id uint64) error
}

type switchRepository struct {
	model fitnessmodel.FitnessTemplateSwitchModel
	conn  sqlx.SqlConn
}

func NewSwitchRepository(repo *repository.Repository) SwitchRepository {
	return &switchRepository{model: repo.FitnessTemplateSwitchModel, conn: repo.DB}
}

func (r *switchRepository) ListByUsers(ctx context.Context, userIDs []uint64) ([]fitnessmodel.FitnessTemplateSwitch, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	query, args, err := sq.Select(switchColumns).From("`fitness_template_switch`").
		Where(sq.Eq{"deleted_at": 0, "user_id": userIDs}).
		OrderBy("user_id ASC", "effective_date ASC", "id ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []fitnessmodel.FitnessTemplateSwitch
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *switchRepository) FindPage(ctx context.Context, userID uint64, page, pageSize int64) ([]fitnessmodel.FitnessTemplateSwitch, int64, error) {
	conditions := sq.And{sq.Eq{"deleted_at": 0}}
	if userID > 0 {
		conditions = append(conditions, sq.Eq{"user_id": userID})
	}
	var list []fitnessmodel.FitnessTemplateSwitch
	total, err := findPage(ctx, r.conn, &list, sq.Select(switchColumns).From("`fitness_template_switch`"), "`fitness_template_switch`", conditions, page, pageSize, "id DESC")
	return list, total, err
}

func (r *switchRepository) Create(ctx context.Context, s *fitnessmodel.FitnessTemplateSwitch) error {
	_, err := r.model.Insert(ctx, s)
	return err
}

func (r *switchRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
