package fitness

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const templateDayColumns = "id, template_id, weekday, plan_date, training_type, training_time, step_goal, warmup, exercises, stretch, training_note, meals, tip, created_at, updated_at, deleted_at"

type TemplateDayRepository interface {
	FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessTemplateDay, error)
	FindByKey(ctx context.Context, templateID uint64, weekday int64, planDate string) (*fitnessmodel.FitnessTemplateDay, error)
	// ListByTemplates 指定模板的全部周模板行 + 覆盖行；覆盖行只取 [overrideFrom, overrideTo]，两者为空则不限
	ListByTemplates(ctx context.Context, templateIDs []uint64, overrideFrom, overrideTo string) ([]fitnessmodel.FitnessTemplateDay, error)
	Create(ctx context.Context, d *fitnessmodel.FitnessTemplateDay) error
	Update(ctx context.Context, d *fitnessmodel.FitnessTemplateDay) error
	DeleteByID(ctx context.Context, id uint64) error
}

type templateDayRepository struct {
	model fitnessmodel.FitnessTemplateDayModel
	conn  sqlx.SqlConn
}

func NewTemplateDayRepository(repo *repository.Repository) TemplateDayRepository {
	return &templateDayRepository{model: repo.FitnessTemplateDayModel, conn: repo.DB}
}

func (r *templateDayRepository) FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessTemplateDay, error) {
	return notFoundAsNil(r.model.FindOne(ctx, id))
}

func (r *templateDayRepository) FindByKey(ctx context.Context, templateID uint64, weekday int64, planDate string) (*fitnessmodel.FitnessTemplateDay, error) {
	return notFoundAsNil(r.model.FindOneByTemplateIdWeekdayPlanDateDeletedAt(ctx, templateID, weekday, planDate, 0))
}

func (r *templateDayRepository) ListByTemplates(ctx context.Context, templateIDs []uint64, overrideFrom, overrideTo string) ([]fitnessmodel.FitnessTemplateDay, error) {
	if len(templateIDs) == 0 {
		return nil, nil
	}
	overrideCond := sq.And{sq.Eq{"weekday": 0}}
	if overrideFrom != "" {
		overrideCond = append(overrideCond, sq.GtOrEq{"plan_date": overrideFrom})
	}
	if overrideTo != "" {
		overrideCond = append(overrideCond, sq.LtOrEq{"plan_date": overrideTo})
	}
	query, args, err := sq.Select(templateDayColumns).From("`fitness_template_day`").
		Where(sq.And{
			sq.Eq{"deleted_at": 0},
			sq.Eq{"template_id": templateIDs},
			sq.Or{sq.Gt{"weekday": 0}, overrideCond},
		}).
		OrderBy("template_id ASC", "weekday DESC", "plan_date ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []fitnessmodel.FitnessTemplateDay
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *templateDayRepository) Create(ctx context.Context, d *fitnessmodel.FitnessTemplateDay) error {
	result, err := r.model.Insert(ctx, d)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	d.Id = uint64(id)
	return nil
}

func (r *templateDayRepository) Update(ctx context.Context, d *fitnessmodel.FitnessTemplateDay) error {
	return r.model.Update(ctx, d)
}

func (r *templateDayRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
