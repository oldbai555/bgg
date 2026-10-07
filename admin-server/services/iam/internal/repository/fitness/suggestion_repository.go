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

const suggestionColumns = "id, user_id, week_start, rule_code, from_template_id, to_template_id, reason, status, decided_at, created_at, updated_at, deleted_at"

type SuggestionRepository interface {
	FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessSuggestion, error)
	FindByUserWeek(ctx context.Context, userID uint64, weekStart string) (*fitnessmodel.FitnessSuggestion, error)
	ListPendingByUser(ctx context.Context, userID uint64) ([]fitnessmodel.FitnessSuggestion, error)
	Create(ctx context.Context, s *fitnessmodel.FitnessSuggestion) error
	Update(ctx context.Context, s *fitnessmodel.FitnessSuggestion) error
}

type suggestionRepository struct {
	model fitnessmodel.FitnessSuggestionModel
	conn  sqlx.SqlConn
}

func NewSuggestionRepository(repo *repository.Repository) SuggestionRepository {
	return &suggestionRepository{model: repo.FitnessSuggestionModel, conn: repo.DB}
}

func (r *suggestionRepository) FindByID(ctx context.Context, id uint64) (*fitnessmodel.FitnessSuggestion, error) {
	return notFoundAsNil(r.model.FindOne(ctx, id))
}

func (r *suggestionRepository) FindByUserWeek(ctx context.Context, userID uint64, weekStart string) (*fitnessmodel.FitnessSuggestion, error) {
	return notFoundAsNil(r.model.FindOneByUserIdWeekStartDeletedAt(ctx, userID, weekStart, 0))
}

func (r *suggestionRepository) ListPendingByUser(ctx context.Context, userID uint64) ([]fitnessmodel.FitnessSuggestion, error) {
	query, args, err := sq.Select(suggestionColumns).From("`fitness_suggestion`").
		Where(sq.Eq{"deleted_at": 0, "user_id": userID, "status": consts.FitnessSuggestionPending}).OrderBy("week_start ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []fitnessmodel.FitnessSuggestion
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *suggestionRepository) Create(ctx context.Context, s *fitnessmodel.FitnessSuggestion) error {
	result, err := r.model.Insert(ctx, s)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	s.Id = uint64(id)
	return nil
}

func (r *suggestionRepository) Update(ctx context.Context, s *fitnessmodel.FitnessSuggestion) error {
	return r.model.Update(ctx, s)
}
