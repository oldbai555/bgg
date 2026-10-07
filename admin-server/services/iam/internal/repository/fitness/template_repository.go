package fitness

import (
	"context"
	"errors"

	"postapocgame/admin-server/pkg/errs"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const templateColumns = "id, code, name, level, stage, summary, daily_kcal, daily_protein, sort, status, created_at, updated_at, deleted_at"

type TemplateRepository interface {
	FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessTemplate, error)
	FindByCode(ctx context.Context, code string) (*fitnessmodel.FitnessTemplate, error)
	ListAll(ctx context.Context) ([]fitnessmodel.FitnessTemplate, error)
	FindPage(ctx context.Context, page, pageSize int64, keyword string) ([]fitnessmodel.FitnessTemplate, int64, error)
	Create(ctx context.Context, t *fitnessmodel.FitnessTemplate) error
	Update(ctx context.Context, t *fitnessmodel.FitnessTemplate) error
	DeleteByID(ctx context.Context, id uint64) error
}

type templateRepository struct {
	model fitnessmodel.FitnessTemplateModel
	conn  sqlx.SqlConn
}

func NewTemplateRepository(repo *repository.Repository) TemplateRepository {
	return &templateRepository{model: repo.FitnessTemplateModel, conn: repo.DB}
}

func (r *templateRepository) FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessTemplate, error) {
	return notFoundAsNil(r.model.FindOne(ctx, id))
}

func (r *templateRepository) FindByCode(ctx context.Context, code string) (*fitnessmodel.FitnessTemplate, error) {
	return notFoundAsNil(r.model.FindOneByCodeDeletedAt(ctx, code, 0))
}

// ListAll 全部未删除模板（含禁用），模板数量是个位数，统计/名称映射直接全量取。
func (r *templateRepository) ListAll(ctx context.Context) ([]fitnessmodel.FitnessTemplate, error) {
	query, args, err := sq.Select(templateColumns).From("`fitness_template`").
		Where(sq.Eq{"deleted_at": 0}).OrderBy("sort ASC", "id ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []fitnessmodel.FitnessTemplate
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *templateRepository) FindPage(ctx context.Context, page, pageSize int64, keyword string) ([]fitnessmodel.FitnessTemplate, int64, error) {
	conditions := sq.And{sq.Eq{"deleted_at": 0}}
	if keyword != "" {
		pattern := "%" + keyword + "%"
		conditions = append(conditions, sq.Or{sq.Like{"name": pattern}, sq.Like{"code": pattern}})
	}
	var list []fitnessmodel.FitnessTemplate
	total, err := findPage(ctx, r.conn, &list, sq.Select(templateColumns).From("`fitness_template`"), "`fitness_template`", conditions, page, pageSize, "sort ASC, id ASC")
	return list, total, err
}

func (r *templateRepository) Create(ctx context.Context, t *fitnessmodel.FitnessTemplate) error {
	result, err := r.model.Insert(ctx, t)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.Id = uint64(id)
	return nil
}

func (r *templateRepository) Update(ctx context.Context, t *fitnessmodel.FitnessTemplate) error {
	return r.model.Update(ctx, t)
}

func (r *templateRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

// notFoundAsNil 把 model 的 ErrNotFound 统一转换成 (nil, nil)，调用方只判 nil。
func notFoundAsNil[T any](v *T, err error) (*T, error) {
	if errors.Is(err, fitnessmodel.ErrNotFound) {
		return nil, nil
	}
	return v, err
}

// findPage 通用分页：先 COUNT 再取当前页。
func findPage(ctx context.Context, conn sqlx.SqlConn, out any, selectBuilder sq.SelectBuilder, from string, where sq.Sqlizer, page, pageSize int64, orderBy string) (int64, error) {
	countSQL, countArgs, err := sq.Select("COUNT(*)").From(from).Where(where).ToSql()
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var total int64
	if err := conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	listSQL, listArgs, err := selectBuilder.Where(where).OrderBy(orderBy).
		Limit(uint64(pageSize)).Offset(uint64((page - 1) * pageSize)).ToSql()
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := conn.QueryRowsCtx(ctx, out, listSQL, listArgs...); err != nil {
		return 0, err
	}
	return total, nil
}
