package misc

import (
	"context"
	"postapocgame/admin-server/services/iam/internal/repository"

	miscmodel "postapocgame/admin-server/services/iam/internal/model/misc"
	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DailyShortSentenceRepository interface {
	FindByID(ctx context.Context, id uint64) (*miscmodel.DailyShortSentence, error)
	FindPage(ctx context.Context, page, pageSize int64, keyword string, sentenceType int64) ([]miscmodel.DailyShortSentence, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, sentence *miscmodel.DailyShortSentence) error
	Update(ctx context.Context, sentence *miscmodel.DailyShortSentence) error
}

type dailyShortSentenceRepository struct {
	model miscmodel.DailyShortSentenceModel
	conn  sqlx.SqlConn
}

func NewDailyShortSentenceRepository(repo *repository.Repository) DailyShortSentenceRepository {
	return &dailyShortSentenceRepository{model: repo.DailyShortSentenceModel, conn: repo.DB}
}

func (r *dailyShortSentenceRepository) FindByID(ctx context.Context, id uint64) (*miscmodel.DailyShortSentence, error) {
	return r.model.FindOne(ctx, id)
}

func (r *dailyShortSentenceRepository) FindPage(ctx context.Context, page, pageSize int64, keyword string, sentenceType int64) ([]miscmodel.DailyShortSentence, int64, error) {
	// 如果有关键词或类型筛选，需要自定义查询
	if keyword != "" || sentenceType > 0 {
		return r.findPageWithFilter(ctx, page, pageSize, keyword, sentenceType)
	}
	// 否则使用生成的分页方法
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *dailyShortSentenceRepository) findPageWithFilter(ctx context.Context, page, pageSize int64, keyword string, sentenceType int64) ([]miscmodel.DailyShortSentence, int64, error) {
	conditions := sq.And{sq.Eq{"deleted_at": 0}}

	if keyword != "" {
		keywordPattern := "%" + keyword + "%"
		conditions = append(conditions, sq.Or{
			sq.Like{"content": keywordPattern},
			sq.Like{"literature_author": keywordPattern},
		})
	}

	if sentenceType > 0 {
		conditions = append(conditions, sq.Eq{"type": sentenceType})
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`daily_short_sentence`").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	listSQL, listArgs, err := sq.Select("id", "type", "content", "img", "literature_author", "convert_img", "created_at", "updated_at", "deleted_at").
		From("`daily_short_sentence`").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []miscmodel.DailyShortSentence
	err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *dailyShortSentenceRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *dailyShortSentenceRepository) Create(ctx context.Context, sentence *miscmodel.DailyShortSentence) error {
	result, err := r.model.Insert(ctx, sentence)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sentence.Id = uint64(id)
	return nil
}

func (r *dailyShortSentenceRepository) Update(ctx context.Context, sentence *miscmodel.DailyShortSentence) error {
	return r.model.Update(ctx, sentence)
}
