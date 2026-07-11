package iam

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	chatdomain "postapocgame/admin-server/internal/domain/chat"
	iammodel "postapocgame/admin-server/internal/model/iam"
	"postapocgame/admin-server/internal/repository"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserDomainService 承载"建用户"这类跨越 IAM 自身表 + 需要触发 Chat 域初始化的编排逻辑。
type UserDomainService struct {
	repo       *repository.Repository
	onboarding chatdomain.Onboarding // 窄接口，不 import internal/repository/chat
}

func NewUserDomainService(repo *repository.Repository, onboarding chatdomain.Onboarding) *UserDomainService {
	return &UserDomainService{repo: repo, onboarding: onboarding}
}

type CreateUserInput struct {
	Username, Nickname, Password, Avatar, Signature string
	DepartmentId                                    uint64
	Status                                           int64
}

// CreateUser 建用户（用户名唯一性校验 + 密码加密 + 落库包在事务里）。
// Chat 初始化异步尽力而为，失败不回滚用户创建，这是产品既定语义。
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
		}
		if err := userRepo.Create(ctx, user); err != nil {
			return errs.Wrap(errs.CodeInternalError, "创建用户失败", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Chat 初始化：异步、尽力而为，不阻塞建用户请求、失败不回滚（产品既定语义，见 04-domain-iam-chat.md）。
	go func(newUserID uint64) {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("Chat onboarding 发生 panic: userId=%d, err=%v", newUserID, r)
			}
		}()
		if err := s.onboarding.InitNewUser(context.Background(), newUserID); err != nil {
			logx.Errorf("初始化新用户聊天数据失败: userId=%d, err=%v", newUserID, err)
		}
	}(user.Id)

	return user, nil
}
