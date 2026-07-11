package chat

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	chatmodel "postapocgame/admin-server/internal/model/chat"
)

type ChatUserRepository interface {
	FindByChatID(ctx context.Context, chatID uint64) ([]chatmodel.ChatUser, error)
	FindByUserID(ctx context.Context, userID uint64) ([]chatmodel.ChatUser, error)
	Create(ctx context.Context, chatUser *chatmodel.ChatUser) error
	DeleteByChatIDAndUserID(ctx context.Context, chatID, userID uint64) error
}

type chatUserRepository struct {
	model chatmodel.ChatUserModel
	conn  sqlx.SqlConn
}

func NewChatUserRepository(repo *repository.Repository) ChatUserRepository {
	return &chatUserRepository{model: repo.ChatUserModel, conn: repo.DB}
}

func (r *chatUserRepository) FindByChatID(ctx context.Context, chatID uint64) ([]chatmodel.ChatUser, error) {
	query := `SELECT * FROM chat_user WHERE chat_id = ? ORDER BY joined_at ASC`
	var list []chatmodel.ChatUser
	err := r.conn.QueryRowsCtx(ctx, &list, query, chatID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *chatUserRepository) FindByUserID(ctx context.Context, userID uint64) ([]chatmodel.ChatUser, error) {
	query := `SELECT * FROM chat_user WHERE user_id = ? ORDER BY joined_at ASC`
	var list []chatmodel.ChatUser
	err := r.conn.QueryRowsCtx(ctx, &list, query, userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *chatUserRepository) Create(ctx context.Context, chatUser *chatmodel.ChatUser) error {
	result, err := r.model.Insert(ctx, chatUser)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	chatUser.Id = uint64(id)
	return nil
}

func (r *chatUserRepository) DeleteByChatIDAndUserID(ctx context.Context, chatID, userID uint64) error {
	query := `DELETE FROM chat_user WHERE chat_id = ? AND user_id = ?`
	_, err := r.conn.ExecCtx(ctx, query, chatID, userID)
	return err
}
