package chat

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/pkg/errs"
	chatmodel "postapocgame/admin-server/internal/model/chat"
)

type ChatMessageRepository interface {
	FindByID(ctx context.Context, id uint64) (*chatmodel.ChatMessage, error)
	FindPage(ctx context.Context, page, pageSize int64, roomId string, userId uint64) ([]chatmodel.ChatMessage, int64, error)
	FindByChatID(ctx context.Context, page, pageSize int64, chatId uint64) ([]chatmodel.ChatMessage, int64, error)
	FindPrivateMessages(ctx context.Context, page, pageSize int64, currentUserId, targetUserId uint64) ([]chatmodel.ChatMessage, int64, error)
	Create(ctx context.Context, message *chatmodel.ChatMessage) error
	Update(ctx context.Context, message *chatmodel.ChatMessage) error
	DeleteByID(ctx context.Context, id uint64) error
}

type chatMessageRepository struct {
	model chatmodel.ChatMessageModel
	conn  sqlx.SqlConn
}

func NewChatMessageRepository(repo *repository.Repository) ChatMessageRepository {
	return &chatMessageRepository{model: repo.ChatMessageModel, conn: repo.DB}
}

func (r *chatMessageRepository) FindByID(ctx context.Context, id uint64) (*chatmodel.ChatMessage, error) {
	return r.model.FindOne(ctx, id)
}

func (r *chatMessageRepository) FindPage(ctx context.Context, page, pageSize int64, roomId string, userId uint64) ([]chatmodel.ChatMessage, int64, error) {
	// 构建查询条件
	where := sq.And{sq.Eq{"deleted_at": 0}}

	if roomId != "" {
		where = append(where, sq.Eq{"room_id": roomId, "to_user_id": 0})
	} else if userId > 0 {
		// 私聊：查询与指定用户相关的消息（包括发送和接收）
		// 注意：这里需要结合当前用户ID来过滤，但当前方法没有当前用户ID参数
		// 暂时保持原逻辑，前端需要额外过滤
		where = append(where, sq.Or{
			sq.Eq{"from_user_id": userId},
			sq.Eq{"to_user_id": userId},
		})
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`chat_message`").Where(where).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listSQL, listArgs, err := sq.Select("*").
		From("`chat_message`").
		Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []chatmodel.ChatMessage
	err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *chatMessageRepository) FindByChatID(ctx context.Context, page, pageSize int64, chatId uint64) ([]chatmodel.ChatMessage, int64, error) {
	// 根据 chatId 查询消息，如果 chatId == 0，则查询所有消息
	where := sq.And{sq.Eq{"deleted_at": 0}}
	if chatId > 0 {
		where = append(where, sq.Eq{"chat_id": chatId})
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`chat_message`").Where(where).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listSQL, listArgs, err := sq.Select("*").
		From("`chat_message`").
		Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []chatmodel.ChatMessage
	err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *chatMessageRepository) FindPrivateMessages(ctx context.Context, page, pageSize int64, currentUserId, targetUserId uint64) ([]chatmodel.ChatMessage, int64, error) {
	// 查询当前用户和指定用户之间的私聊消息
	where := sq.And{
		sq.Or{
			sq.And{sq.Eq{"from_user_id": currentUserId}, sq.Eq{"to_user_id": targetUserId}},
			sq.And{sq.Eq{"from_user_id": targetUserId}, sq.Eq{"to_user_id": currentUserId}},
		},
		sq.Eq{"deleted_at": 0},
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`chat_message`").Where(where).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listSQL, listArgs, err := sq.Select("*").
		From("`chat_message`").
		Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []chatmodel.ChatMessage
	err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *chatMessageRepository) Create(ctx context.Context, message *chatmodel.ChatMessage) error {
	result, err := r.model.Insert(ctx, message)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	message.Id = uint64(id)
	return nil
}

func (r *chatMessageRepository) Update(ctx context.Context, message *chatmodel.ChatMessage) error {
	return r.model.Update(ctx, message)
}

func (r *chatMessageRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
