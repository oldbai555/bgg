package chat

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	chatmodel "postapocgame/admin-server/internal/model/chat"
)

type ChatRepository interface {
	FindByID(ctx context.Context, id uint64) (*chatmodel.Chat, error)
	FindByUserID(ctx context.Context, userID uint64) ([]chatmodel.Chat, error)
	FindUsersByChatID(ctx context.Context, chatID uint64) ([]chatmodel.ChatUser, error)
	Create(ctx context.Context, chat *chatmodel.Chat) error
	Update(ctx context.Context, chat *chatmodel.Chat) error
	DeleteByID(ctx context.Context, id uint64) error
	FindPrivateChatByUserIDs(ctx context.Context, userID1, userID2 uint64) (*chatmodel.Chat, error)
	FindGroups(ctx context.Context, page, pageSize int64, name string) ([]chatmodel.Chat, int64, error)
	CountMembersByChatID(ctx context.Context, chatID uint64) (int64, error)
}

type chatRepository struct {
	model chatmodel.ChatModel
	conn  sqlx.SqlConn
}

func NewChatRepository(repo *repository.Repository) ChatRepository {
	return &chatRepository{model: repo.ChatModel, conn: repo.DB}
}

func (r *chatRepository) FindByID(ctx context.Context, id uint64) (*chatmodel.Chat, error) {
	return r.model.FindOne(ctx, id)
}

func (r *chatRepository) FindByUserID(ctx context.Context, userID uint64) ([]chatmodel.Chat, error) {
	// 通过chat_user关联表查询用户参与的所有聊天
	query := `
		SELECT c.* 
		FROM chat c
		INNER JOIN chat_user cu ON c.id = cu.chat_id
		WHERE cu.user_id = ? AND c.deleted_at = 0
		ORDER BY c.created_at DESC
	`
	var list []chatmodel.Chat
	err := r.conn.QueryRowsCtx(ctx, &list, query, userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *chatRepository) FindUsersByChatID(ctx context.Context, chatID uint64) ([]chatmodel.ChatUser, error) {
	// 查询聊天中的所有用户
	query := `
		SELECT * 
		FROM chat_user
		WHERE chat_id = ?
		ORDER BY joined_at ASC
	`
	var list []chatmodel.ChatUser
	err := r.conn.QueryRowsCtx(ctx, &list, query, chatID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *chatRepository) Create(ctx context.Context, chat *chatmodel.Chat) error {
	_, err := r.model.Insert(ctx, chat)
	return err
}

func (r *chatRepository) Update(ctx context.Context, chat *chatmodel.Chat) error {
	return r.model.Update(ctx, chat)
}

func (r *chatRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *chatRepository) FindPrivateChatByUserIDs(ctx context.Context, userID1, userID2 uint64) (*chatmodel.Chat, error) {
	// 查找两个用户之间的私聊（type=1）
	// 私聊必须包含且仅包含这两个用户
	query := `
		SELECT c.* 
		FROM chat c
		INNER JOIN chat_user cu1 ON c.id = cu1.chat_id AND cu1.user_id = ?
		INNER JOIN chat_user cu2 ON c.id = cu2.chat_id AND cu2.user_id = ?
		WHERE c.type = 1 AND c.deleted_at = 0
		GROUP BY c.id
		HAVING COUNT(DISTINCT cu.user_id) = 2
		LIMIT 1
	`
	// 修复：使用子查询来正确计算用户数
	query = `
		SELECT c.* 
		FROM chat c
		WHERE c.type = 1 AND c.deleted_at = 0
		AND EXISTS (SELECT 1 FROM chat_user cu1 WHERE cu1.chat_id = c.id AND cu1.user_id = ?)
		AND EXISTS (SELECT 1 FROM chat_user cu2 WHERE cu2.chat_id = c.id AND cu2.user_id = ?)
		AND (SELECT COUNT(*) FROM chat_user cu WHERE cu.chat_id = c.id) = 2
		LIMIT 1
	`
	var chat chatmodel.Chat
	err := r.conn.QueryRowCtx(ctx, &chat, query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// FindGroups 查询群组列表（分页、搜索）
func (r *chatRepository) FindGroups(ctx context.Context, page, pageSize int64, name string) ([]chatmodel.Chat, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	// 构建查询条件（使用 squirrel）
	conditions := sq.And{
		sq.Eq{"type": 2},
		sq.Eq{"deleted_at": 0},
	}
	if name != "" {
		conditions = append(conditions, sq.Like{"name": "%" + name + "%"})
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").
		From("`chat`").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// 查询列表
	listSQL, listArgs, err := sq.Select("*").
		From("`chat`").
		Where(conditions).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []chatmodel.Chat
	if err := r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// CountMembersByChatID 统计群组成员数量
func (r *chatRepository) CountMembersByChatID(ctx context.Context, chatID uint64) (int64, error) {
	query := "SELECT COUNT(*) FROM `chat_user` WHERE chat_id = ?"
	var count int64
	err := r.conn.QueryRowCtx(ctx, &count, query, chatID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
