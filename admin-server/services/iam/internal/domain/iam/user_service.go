package iam

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/repository"
	iamrepo "postapocgame/admin-server/services/iam/internal/repository/iam"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	zredis "github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// StreamChatUserCreated 必须和 services/chat/internal/consumer/chat_user_created_consumer.go
// 里的同名常量保持一致，两边各自维护一份（16-rpc-conventions.md 第 6 节"直接复制不共享"）。
const StreamChatUserCreated = "stream:chat.user.created"

// chatUserCreatedEvent 与 chat-rpc 侧 chatUserCreatedEvent 字段一一对应。
type chatUserCreatedEvent struct {
	UserID    uint64 `json:"userId"`
	CreatedAt int64  `json:"createdAt"`
}

// UserDomainService 承载"建用户"这类跨越 IAM 自身表 + 需要触发 Chat 域初始化的编排逻辑。
type UserDomainService struct {
	repo  *repository.Repository
	redis *zredis.Redis
}

func NewUserDomainService(repo *repository.Repository, redis *zredis.Redis) *UserDomainService {
	return &UserDomainService{repo: repo, redis: redis}
}

type CreateUserInput struct {
	Username, Nickname, Password, Avatar, Signature string
	DepartmentId                                    uint64
	Status                                          int64
}

// CreateUser 建用户（用户名唯一性校验 + 密码加密 + 落库包在事务里）。
// Chat 初始化异步尽力而为，失败不回滚用户创建，这是产品既定语义（04-domain-iam-chat.md）。
// chat 域拆分成独立服务（chat-rpc）之后，触发方式从进程内 goroutine 直调
// chatdomain.Onboarding.InitNewUser 改成发布 stream:chat.user.created 事件，见
// 17-async-eventing.md 第 2.1 节、services/chat/internal/consumer/
// chat_user_created_consumer.go——语义完全不变（失败只记日志，不影响 CreateUser 返回成功），
// 只是触发机制从"进程内直调"换成"跨进程 Streams"。
func (s *UserDomainService) CreateUser(ctx context.Context, in CreateUserInput) (*iammodel.AdminUser, error) {
	if in.Username == "" || in.Password == "" {
		return nil, errs.New(errs.CodeBadRequest, "用户名和密码不能为空")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "密码加密失败", err)
	}

	user := &iammodel.AdminUser{
		Username:     in.Username,
		Nickname:     in.Nickname,
		PasswordHash: string(hash),
		Avatar:       in.Avatar,
		Signature:    in.Signature,
		DepartmentId: in.DepartmentId,
		Status:       in.Status,
	}

	err = s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		userRepo := iamrepo.NewUserRepository(txRepo)
		if _, err := userRepo.FindByUsername(ctx, in.Username); err == nil {
			return errs.New(errs.CodeBadRequest, "用户名已存在")
		} else if !errors.Is(err, sqlx.ErrNotFound) {
			return errs.Wrap(errs.CodeInternalError, "查询用户名失败", err)
		}
		if err := userRepo.Create(ctx, user); err != nil {
			return errs.Wrap(errs.CodeInternalError, "创建用户失败", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.publishChatUserCreated(ctx, user.Id)

	return user, nil
}

// publishChatUserCreated 发布失败同样只记日志，不回滚用户创建——与拆分前
// "goroutine 直调 + recover 兜底"的既有语义一致，只是失败面从"onboarding 内部报错"换成
// "XAdd 失败"。XAdd 本身是同步快速调用，不需要再包一层 goroutine。
func (s *UserDomainService) publishChatUserCreated(ctx context.Context, newUserID uint64) {
	payload, err := json.Marshal(chatUserCreatedEvent{UserID: newUserID, CreatedAt: time.Now().Unix()})
	if err != nil {
		logx.Errorf("序列化 chat.user.created 事件失败: userId=%d, err=%v", newUserID, err)
		return
	}
	if _, err := s.redis.XAddCtx(ctx, StreamChatUserCreated, false, "*", []string{"payload", string(payload)}); err != nil {
		logx.Errorf("发布 %s 事件失败: userId=%d, err=%v", StreamChatUserCreated, newUserID, err)
	}
}
