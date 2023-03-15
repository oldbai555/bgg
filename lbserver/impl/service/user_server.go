package service

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/client/lbuser"
	"github.com/oldbai555/bgg/lbserver/impl/cache"
	"github.com/oldbai555/bgg/lbserver/impl/constant"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

var UserServer LbuserServer

type LbuserServer struct {
	*lbuser.UnimplementedLbuserServer
}

func (a *LbuserServer) ResetPassword(ctx context.Context, req *lbuser.ResetPasswordReq) (*lbuser.ResetPasswordRsp, error) {
	var rsp lbuser.ResetPasswordRsp
	var err error
	req.NewPassword = utils.StrMd5(req.NewPassword)

	err = User.UpdateById(ctx, req.Id, map[string]interface{}{
		lbuser.FieldPassword_: req.NewPassword,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbuserServer) GetFrontUser(ctx context.Context, req *lbuser.GetFrontUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp

	user, err := User.GetByAdmin(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Avatar = user.Avatar
	rsp.Email = user.Email
	rsp.Github = user.Github
	rsp.Nickname = user.Nickname
	rsp.Desc = user.Desc
	return &rsp, nil
}

func (a *LbuserServer) Login(ctx context.Context, req *lbuser.LoginReq) (*lbuser.LoginRsp, error) {
	var rsp lbuser.LoginRsp
	var err error

	user, err := User.GetByUserName(ctx, req.Username)
	if err != nil {
		log.Errorf("err:%v", err)
		if err == gorm.ErrRecordNotFound {
			err = lbuser.ErrUserNotFound
		}
		return nil, err
	}

	if user.Password != utils.StrMd5(req.Password) {
		return nil, lbuser.ErrPasswordInvalid
	}

	genToken, err := webtool.GenToken(ctx, user.Id)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	sid := utils.GenUUID()
	err = cache.SetLoginUserToken(ctx, sid, genToken, time.Hour)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	err = cache.SetLoginUser(ctx, fmt.Sprintf("%d", user.Id), user, time.Hour)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.Sid = sid

	return &rsp, nil
}
func (a *LbuserServer) UpdateLoginUserInfo(ctx context.Context, req *lbuser.UpdateLoginUserInfoReq) (*lbuser.UpdateLoginUserInfoRsp, error) {
	var rsp lbuser.UpdateLoginUserInfoRsp
	var err error

	claims, ok := ctx.Value(constant.CtxWithClaim).(*webtool.Claims)
	if !ok {
		log.Errorf("err is : %v", lbuser.ErrUserLoginExpired)
		return nil, lbuser.ErrUserLoginExpired
	}

	var updateMap = map[string]interface{}{
		lbuser.FieldAvatar_:   req.User.Avatar,
		lbuser.FieldEmail_:    req.User.Email,
		lbuser.FieldGithub_:   req.User.Github,
		lbuser.FieldNickname_: req.User.Nickname,
		lbuser.FieldDesc:      req.User.Desc,
	}

	err = User.UpdateById(ctx, claims.UserId, updateMap)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetUserList(ctx context.Context, req *lbuser.GetUserListReq) (*lbuser.GetUserListRsp, error) {
	var rsp lbuser.GetUserListRsp
	var err error

	rsp.List, rsp.Page, err = User.FindPage(ctx, req.ListOption)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) DelUser(ctx context.Context, req *lbuser.DelUserReq) (*lbuser.DelUserRsp, error) {
	var rsp lbuser.DelUserRsp
	var err error

	err = User.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) UpdateUserNameWithRole(ctx context.Context, req *lbuser.UpdateUserNameWithRoleReq) (*lbuser.UpdateUserNameWithRoleRsp, error) {
	var rsp lbuser.UpdateUserNameWithRoleRsp

	exit, err := User.CheckUserNameExit(ctx, req.Id, req.Username)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if exit {
		return nil, lbuser.ErrUserExit
	}

	err = User.UpdateById(ctx, req.Id, map[string]interface{}{
		lbuser.FieldUsername_: req.Username,
		lbuser.FieldRole_:     req.Role,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) Logout(ctx context.Context, req *lbuser.LogoutReq) (*lbuser.LogoutRsp, error) {
	var rsp lbuser.LogoutRsp

	err := cache.DelLoginUserToken(ctx, req.Sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetLoginUser(ctx context.Context, req *lbuser.GetLoginUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp
	claims, ok := ctx.Value(constant.CtxWithClaim).(*webtool.Claims)
	if !ok {
		log.Errorf("err is : %v", constant.ErrGetLoginFail)
		return nil, constant.ErrGetLoginFail
	}

	user, err := User.GetById(ctx, claims.UserId)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	rsp.Avatar = user.Avatar
	rsp.Email = user.Email
	rsp.Github = user.Github
	rsp.Nickname = user.Nickname
	rsp.Desc = user.Desc
	return &rsp, nil
}
func (a *LbuserServer) AddUser(ctx context.Context, req *lbuser.AddUserReq) (*lbuser.AddUserRsp, error) {
	var rsp lbuser.AddUserRsp

	exit, err := User.CheckUserNameExit(ctx, 0, req.User.Username)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if exit {
		return nil, lbuser.ErrUserExit
	}

	err = User.Create(ctx, req.User)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetUser(ctx context.Context, req *lbuser.GetUserReq) (*lbuser.GetUserRsp, error) {
	var rsp lbuser.GetUserRsp
	var err error

	rsp.User, err = User.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
