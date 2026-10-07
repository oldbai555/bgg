package fitness

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	"postapocgame/admin-server/services/iam/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// MemberRow 使用者档案 + admin_user 展示信息（同在 admin 库，直接 JOIN）。
type MemberRow struct {
	UserId             uint64  `db:"user_id"`
	Username           string  `db:"username"`
	Nickname           string  `db:"nickname"`
	Avatar             string  `db:"avatar"`
	HeightCm           float64 `db:"height_cm"`
	TargetWeightKg     float64 `db:"target_weight_kg"`
	TrainingPreference int64   `db:"training_preference"`
	CreatedAt          int64   `db:"created_at"`
}

const memberColumns = "p.user_id, u.username, u.nickname, u.avatar, p.height_cm, p.target_weight_kg, p.training_preference, p.created_at"

type ProfileRepository interface {
	FindByUserID(ctx context.Context, userID uint64) (*fitnessmodel.FitnessProfile, error)
	Create(ctx context.Context, p *fitnessmodel.FitnessProfile) error
	Update(ctx context.Context, p *fitnessmodel.FitnessProfile) error
	FindMember(ctx context.Context, userID uint64) (*MemberRow, error)
	FindMemberPage(ctx context.Context, page, pageSize int64, keyword string) ([]MemberRow, int64, error)
	// ListMembers 全部使用者（统计/排行用；使用者规模是公司内部人数级别）
	ListMembers(ctx context.Context, userIDs []uint64) ([]MemberRow, error)
	// UserNicknames 任意 admin_user 的展示名（切换记录的操作人可能是没有档案的管理员），昵称为空回落用户名
	UserNicknames(ctx context.Context, userIDs []uint64) (map[uint64]string, error)
}

type profileRepository struct {
	model fitnessmodel.FitnessProfileModel
	conn  sqlx.SqlConn
}

func NewProfileRepository(repo *repository.Repository) ProfileRepository {
	return &profileRepository{model: repo.FitnessProfileModel, conn: repo.DB}
}

func (r *profileRepository) FindByUserID(ctx context.Context, userID uint64) (*fitnessmodel.FitnessProfile, error) {
	return notFoundAsNil(r.model.FindOneByUserIdDeletedAt(ctx, userID, 0))
}

func (r *profileRepository) Create(ctx context.Context, p *fitnessmodel.FitnessProfile) error {
	_, err := r.model.Insert(ctx, p)
	return err
}

func (r *profileRepository) Update(ctx context.Context, p *fitnessmodel.FitnessProfile) error {
	return r.model.Update(ctx, p)
}

func memberSelect() sq.SelectBuilder {
	return sq.Select(memberColumns).From("`fitness_profile` p").
		Join("`admin_user` u ON u.id = p.user_id AND u.deleted_at = 0")
}

func (r *profileRepository) FindMember(ctx context.Context, userID uint64) (*MemberRow, error) {
	query, args, err := memberSelect().Where(sq.Eq{"p.deleted_at": 0, "p.user_id": userID}).Limit(1).ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []MemberRow
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (r *profileRepository) FindMemberPage(ctx context.Context, page, pageSize int64, keyword string) ([]MemberRow, int64, error) {
	conditions := sq.And{sq.Eq{"p.deleted_at": 0}}
	if keyword != "" {
		pattern := "%" + keyword + "%"
		conditions = append(conditions, sq.Or{sq.Like{"u.nickname": pattern}, sq.Like{"u.username": pattern}})
	}
	from := "`fitness_profile` p JOIN `admin_user` u ON u.id = p.user_id AND u.deleted_at = 0"
	var list []MemberRow
	total, err := findPage(ctx, r.conn, &list, memberSelect(), from, conditions, page, pageSize, "p.id DESC")
	return list, total, err
}

func (r *profileRepository) ListMembers(ctx context.Context, userIDs []uint64) ([]MemberRow, error) {
	conditions := sq.And{sq.Eq{"p.deleted_at": 0}}
	if userIDs != nil {
		if len(userIDs) == 0 {
			return nil, nil
		}
		conditions = append(conditions, sq.Eq{"p.user_id": userIDs})
	}
	query, args, err := memberSelect().Where(conditions).OrderBy("p.id ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var list []MemberRow
	if err := r.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *profileRepository) UserNicknames(ctx context.Context, userIDs []uint64) (map[uint64]string, error) {
	out := make(map[uint64]string, len(userIDs))
	if len(userIDs) == 0 {
		return out, nil
	}
	query, args, err := sq.Select("id", "username", "nickname").From("`admin_user`").
		Where(sq.Eq{"id": userIDs}).ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	var rows []struct {
		Id       uint64 `db:"id"`
		Username string `db:"username"`
		Nickname string `db:"nickname"`
	}
	if err := r.conn.QueryRowsCtx(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	for _, row := range rows {
		name := row.Nickname
		if name == "" {
			name = row.Username
		}
		out[row.Id] = name
	}
	return out, nil
}
