package fitness

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const tipColumns = "id, code, title, content, sort, status, created_at, updated_at, deleted_at"

type TipRepository interface {
	FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessTip, error)
	FindByCode(ctx context.Context, code string) (*fitnessmodel.FitnessTip, error)
	ListEnabled(ctx context.Context) ([]fitnessmodel.FitnessTip, error)
	FindPage(ctx context.Context, page, pageSize int64, keyword string) ([]fitnessmodel.FitnessTip, int64, error)
	Create(ctx context.Context, t *fitnessmodel.FitnessTip) error
	Update(ctx context.Context, t *fitnessmodel.FitnessTip) error
	DeleteByID(ctx context.Context, id uint64) error
}

type tipRepository struct {
	model fitnessmodel.FitnessTipModel
	conn  sqlx.SqlConn
}

func NewTipRepository(repo *repository.Repository) TipRepository {
	return &tipRepository{model: repo.FitnessTipModel, conn: repo.DB}
}

func (r *tipRepository) FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessTip, error) {
	return notFoundAsNil(r.model.FindOne(ctx, id))
}

func (r *tipRepository) FindByCode(ctx context.Context, code string) (*fitnessmodel.FitnessTip, error) {
	return notFoundAsNil(r.model.FindOneByCodeDeletedAt(ctx, code, 0))
}

func (r *tipRepository) ListEnabled(ctx context.Context) ([]fitnessmodel.FitnessTip, error) {
	query, args, err := sq.Select(tipColumns).From("`fitness_tip`").
		Where(sq.Eq{"deleted_at": 0, "status": consts.FitnessStatusEnabled}).OrderBy("sort ASC", "id ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []fitnessmodel.FitnessTip
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *tipRepository) FindPage(ctx context.Context, page, pageSize int64, keyword string) ([]fitnessmodel.FitnessTip, int64, error) {
	conditions := sq.And{sq.Eq{"deleted_at": 0}}
	if keyword != "" {
		pattern := "%" + keyword + "%"
		conditions = append(conditions, sq.Or{sq.Like{"title": pattern}, sq.Like{"content": pattern}})
	}
	var list []fitnessmodel.FitnessTip
	total, err := findPage(ctx, r.conn, &list, sq.Select(tipColumns).From("`fitness_tip`"), "`fitness_tip`", conditions, page, pageSize, "sort ASC, id ASC")
	return list, total, err
}

func (r *tipRepository) Create(ctx context.Context, t *fitnessmodel.FitnessTip) error {
	_, err := r.model.Insert(ctx, t)
	return err
}

func (r *tipRepository) Update(ctx context.Context, t *fitnessmodel.FitnessTip) error {
	return r.model.Update(ctx, t)
}

func (r *tipRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
