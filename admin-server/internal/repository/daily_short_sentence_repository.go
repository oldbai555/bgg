package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model"
)

type DailyShortSentenceRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.DailyShortSentence, error)
	FindPage(ctx context.Context, page, pageSize int64, keyword string, sentenceType int64) ([]model.DailyShortSentence, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, sentence *model.DailyShortSentence) error
	Update(ctx context.Context, sentence *model.DailyShortSentence) error
}

type dailyShortSentenceRepository struct {
	model model.DailyShortSentenceModel
	conn  sqlx.SqlConn
}

func NewDailyShortSentenceRepository(repo *Repository) DailyShortSentenceRepository {
	return &dailyShortSentenceRepository{model: repo.DailyShortSentenceModel, conn: repo.DB}
}

func (r *dailyShortSentenceRepository) FindByID(ctx context.Context, id uint64) (*model.DailyShortSentence, error) {
	return r.model.FindOne(ctx, id)
}

func (r *dailyShortSentenceRepository) FindPage(ctx context.Context, page, pageSize int64, keyword string, sentenceType int64) ([]model.DailyShortSentence, int64, error) {
	// 如果有关键词或类型筛选，需要自定义查询
	if keyword != "" || sentenceType > 0 {
		return r.findPageWithFilter(ctx, page, pageSize, keyword, sentenceType)
	}
	// 否则使用生成的分页方法
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *dailyShortSentenceRepository) findPageWithFilter(ctx context.Context, page, pageSize int64, keyword string, sentenceType int64) ([]model.DailyShortSentence, int64, error) {
	var whereConditions []string
	var args []interface{}

	whereConditions = append(whereConditions, "deleted_at = 0")

	if keyword != "" {
		whereConditions = append(whereConditions, "(content LIKE ? OR literature_author LIKE ?)")
		keywordPattern := "%" + keyword + "%"
		args = append(args, keywordPattern, keywordPattern)
	}

	if sentenceType > 0 {
		whereConditions = append(whereConditions, "type = ?")
		args = append(args, sentenceType)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// 查询总数
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `daily_short_sentence` WHERE %s", whereClause)
	err := r.conn.QueryRowCtx(ctx, &total, countQuery, args...)
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

	query := fmt.Sprintf("SELECT id, type, content, img, literature_author, convert_img, created_at, updated_at, deleted_at FROM `daily_short_sentence` WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, pageSize, offset)

	var list []model.DailyShortSentence
	err = r.conn.QueryRowsCtx(ctx, &list, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *dailyShortSentenceRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *dailyShortSentenceRepository) Create(ctx context.Context, sentence *model.DailyShortSentence) error {
	_, err := r.model.Insert(ctx, sentence)
	return err
}

func (r *dailyShortSentenceRepository) Update(ctx context.Context, sentence *model.DailyShortSentence) error {
	return r.model.Update(ctx, sentence)
}
