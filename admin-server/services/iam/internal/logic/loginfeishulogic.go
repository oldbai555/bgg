package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/feishu"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/consts"
	iamdomain "postapocgame/admin-server/services/iam/internal/domain/iam"
	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginFeishuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginFeishuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginFeishuLogic {
	return &LoginFeishuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginFeishuLogic) LoginFeishu(in *iam.LoginFeishuRequest) (*iam.TokenPair, error) {
	if in == nil || in.Code == "" {
		recordLoginLog(l.svcCtx, 0, "", in.GetClientIp(), in.GetUserAgent(), "缺少飞书授权 code", false)
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "缺少飞书授权 code"))
	}

	client := feishu.NewClient(l.svcCtx.Config.Feishu.AppId, l.svcCtx.Config.Feishu.AppSecret, l.svcCtx.Config.Feishu.RedirectUri)
	userInfo, err := client.ExchangeUserInfo(l.ctx, in.Code)
	if err != nil {
		recordLoginLog(l.svcCtx, 0, "", in.ClientIp, in.UserAgent, "飞书授权失败", false)
		return nil, toGRPCStatus(errs.Wrap(errs.CodeUnauthorized, "飞书授权失败", err))
	}

	user, err := l.findOrCreateUser(userInfo)
	if err != nil {
		recordLoginLog(l.svcCtx, 0, userInfo.Name, in.ClientIp, in.UserAgent, "飞书账号建号/绑定失败", false)
		return nil, toGRPCStatus(err)
	}

	if user.Status != 1 {
		recordLoginLog(l.svcCtx, user.Id, user.Username, in.ClientIp, in.UserAgent, "账号已被禁用", false)
		return nil, toGRPCStatus(errs.New(errs.CodeForbidden, "账号已被禁用"))
	}

	accessToken, err := jwthelper.GenerateToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.Issuer,
		l.svcCtx.Config.JWT.AccessExpire,
		user.Id,
		user.Username,
		false,
	)
	if err != nil {
		recordLoginLog(l.svcCtx, user.Id, user.Username, in.ClientIp, in.UserAgent, "生成访问令牌失败", false)
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "生成访问令牌失败", err))
	}

	refreshToken, err := jwthelper.GenerateToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.Issuer,
		l.svcCtx.Config.JWT.RefreshExpire,
		user.Id,
		user.Username,
		true,
	)
	if err != nil {
		recordLoginLog(l.svcCtx, user.Id, user.Username, in.ClientIp, in.UserAgent, "生成刷新令牌失败", false)
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "生成刷新令牌失败", err))
	}

	recordLoginLog(l.svcCtx, user.Id, user.Username, in.ClientIp, in.UserAgent, "飞书登录成功", true)

	go createUnreadNoticeNotifications(l.svcCtx, user.Id)

	return &iam.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// findOrCreateUser 按 open_id 查绑定关系；未绑定则复用 UserDomainService.CreateUser
// 建新号（不绕过用户名唯一性校验/密码加密/事务落库这套统一路径），再写绑定记录。
func (l *LoginFeishuLogic) findOrCreateUser(userInfo *feishu.UserInfo) (*iammodel.AdminUser, error) {
	bind, err := l.svcCtx.Domain.IAM.UserThirdParty.FindByOpenID(l.ctx, consts.FeishuProvider, userInfo.OpenId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询第三方账号绑定失败", err)
	}

	if bind != nil {
		user, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, bind.UserId)
		if err != nil {
			return nil, errs.Wrap(errs.CodeInternalError, "查询绑定用户失败", err)
		}
		return user, nil
	}

	username := "feishu_" + userInfo.OpenId
	nickname := userInfo.Name
	if nickname == "" {
		nickname = "飞书用户"
	}

	user, err := l.svcCtx.Domain.IAM.UserService.CreateUser(l.ctx, iamdomain.CreateUserInput{
		Username:     username,
		Nickname:     nickname,
		Password:     uuid.NewString(),
		Avatar:       userInfo.AvatarUrl,
		Status:       1,
		DepartmentId: l.resolveDefaultDepartmentID(),
	})
	if err != nil {
		// 建号失败大概率是并发首次登录竞态（两个请求同时发现未绑定、同时建号），或者
		// 上一次登录建号成功但绑定写入失败留下的孤儿账号——用户名 feishu_<open_id> 是
		// 确定性生成的，按用户名找回已存在账号自愈，而不是把这种偶发冲突当成登录失败。
		existing, findErr := l.svcCtx.Domain.IAM.User.FindByUsername(l.ctx, username)
		if findErr != nil {
			return nil, errs.Wrap(errs.CodeInternalError, "自动创建飞书用户失败", err)
		}
		user = existing
	} else {
		l.assignDefaultRole(user.Id)
	}

	if err := l.svcCtx.Domain.IAM.UserThirdParty.Create(l.ctx, &iammodel.AdminUserThirdParty{
		UserId:   user.Id,
		Provider: consts.FeishuProvider,
		OpenId:   userInfo.OpenId,
		UnionId:  userInfo.UnionId,
	}); err != nil {
		// 绑定写入失败：可能是并发请求已经写入了同一条绑定（provider+open_id 唯一键冲突），
		// 重新按 open_id 查一次绑定自愈，查不到才是真失败。
		if rebind, findErr := l.svcCtx.Domain.IAM.UserThirdParty.FindByOpenID(l.ctx, consts.FeishuProvider, userInfo.OpenId); findErr != nil || rebind == nil {
			return nil, errs.Wrap(errs.CodeInternalError, "写入第三方账号绑定失败", err)
		}
	}

	return user, nil
}

// resolveDefaultDepartmentID 查"飞书待分配"部门 ID，尽力而为：部门未初始化（尚未跑
// 迁移 SQL）时返回 0（未分配部门），不阻塞登录——用户管理列表会显示部门为空，
// 提醒管理员补跑迁移或手动分配，而不是让飞书用户完全无法登录。
func (l *LoginFeishuLogic) resolveDefaultDepartmentID() uint64 {
	dept, err := l.svcCtx.Domain.IAM.Department.FindByName(l.ctx, consts.FeishuDefaultDepartmentName)
	if err != nil {
		logx.Errorf("查询飞书待分配部门失败，新用户将不归属任何部门: err=%v", err)
		return 0
	}
	return dept.Id
}

// assignDefaultRole 给新建的飞书用户分配默认角色，尽力而为：角色未初始化（尚未跑迁移
// SQL）或分配失败都只记日志，不影响登录本身——避免角色数据缺失导致用户完全无法登录。
func (l *LoginFeishuLogic) assignDefaultRole(userID uint64) {
	role, err := l.svcCtx.Domain.IAM.Role.FindByCode(l.ctx, consts.FeishuDefaultRoleCode)
	if err != nil {
		logx.Errorf("查询飞书默认角色失败，跳过分配: userId=%d, err=%v", userID, err)
		return
	}
	if err := l.svcCtx.Domain.IAM.UserRole.UpdateUserRoles(l.ctx, userID, []uint64{role.Id}); err != nil {
		logx.Errorf("分配飞书默认角色失败: userId=%d, roleId=%d, err=%v", userID, role.Id, err)
	}
}
