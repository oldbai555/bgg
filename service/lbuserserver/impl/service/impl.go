package service

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lbuser"
	cache2 "github.com/oldbai555/bgg/service/lbuserserver/impl/cache"
	"github.com/oldbai555/bgg/service/lbuserserver/impl/dao/impl/mysql"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

var ServerImpl LbuserServer

type LbuserServer struct {
	*lbuser.UnimplementedLbuserServer
}

func (a *LbuserServer) ResetPassword(ctx context.Context, req *lbuser.ResetPasswordReq) (*lbuser.ResetPasswordRsp, error) {
	var rsp lbuser.ResetPasswordRsp
	var err error
	req.NewPassword = utils.StrMd5(req.NewPassword)

	_, err = mysql.UserOrm.UpdateById(ctx, req.Id, map[string]interface{}{
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

	user, err := mysql.UserOrm.GetByAdmin(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.BaseUser = lbuser.NewBaseUser(user)
	return &rsp, nil
}

func (a *LbuserServer) Login(ctx context.Context, req *lbuser.LoginReq) (*lbuser.LoginRsp, error) {
	var rsp lbuser.LoginRsp
	var err error

	user, err := mysql.UserOrm.GetByUserName(ctx, req.Username)
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

	sid := utils.GenUUID()
	genToken, err := webtool.GenToken(ctx, user.Id, sid)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	err = cache2.SetLoginUserToken(ctx, sid, genToken, time.Hour)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	err = cache2.SetLoginUser(ctx, fmt.Sprintf("%d", user.Id), user, time.Hour)
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

	var updateMap = map[string]interface{}{
		lbuser.FieldAvatar_:   req.User.Avatar,
		lbuser.FieldEmail_:    req.User.Email,
		lbuser.FieldGithub_:   req.User.Github,
		lbuser.FieldNickname_: req.User.Nickname,
		lbuser.FieldDesc:      req.User.Desc,
	}

	_, err = mysql.UserOrm.UpdateById(ctx, req.User.Id, updateMap)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetUserList(ctx context.Context, req *lbuser.GetUserListReq) (*lbuser.GetUserListRsp, error) {
	var rsp lbuser.GetUserListRsp
	var err error

	rsp.List, rsp.Paginate, err = mysql.UserOrm.FindPaginate(ctx, req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) DelUser(ctx context.Context, req *lbuser.DelUserReq) (*lbuser.DelUserRsp, error) {
	var rsp lbuser.DelUserRsp
	var err error

	_, err = mysql.UserOrm.DeleteById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) UpdateUserNameWithRole(ctx context.Context, req *lbuser.UpdateUserNameWithRoleReq) (*lbuser.UpdateUserNameWithRoleRsp, error) {
	var rsp lbuser.UpdateUserNameWithRoleRsp

	exit, err := mysql.UserOrm.CheckUserNameExit(ctx, req.Id, req.Username)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if exit {
		return nil, lbuser.ErrUserExit
	}

	_, err = mysql.UserOrm.UpdateById(ctx, req.Id, map[string]interface{}{
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

	err := cache2.DelLoginUserToken(ctx, req.Sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetLoginUser(ctx context.Context, req *lbuser.GetLoginUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp

	token, err := cache2.GetLoginUserToken(ctx, req.GetSid())
	if err != nil && !cache2.IsNotFoundErr(err) {
		log.Errorf("err:%v", err)
		return nil, err
	}

	parseToken, claims, err := webtool.ParseToken(token)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if err != nil || !parseToken.Valid {
		return nil, lbuser.ErrUserLoginExpired
	}

	u, err := mysql.UserOrm.GetById(ctx, claims.GetUserId())
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	rsp.BaseUser = lbuser.NewBaseUser(u)
	return &rsp, nil
}
func (a *LbuserServer) AddUser(ctx context.Context, req *lbuser.AddUserReq) (*lbuser.AddUserRsp, error) {
	var rsp lbuser.AddUserRsp

	exit, err := mysql.UserOrm.CheckUserNameExit(ctx, 0, req.User.Username)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if exit {
		return nil, lbuser.ErrUserExit
	}

	_, err = mysql.UserOrm.Create(ctx, req.User)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetUser(ctx context.Context, req *lbuser.GetUserReq) (*lbuser.GetUserRsp, error) {
	var rsp lbuser.GetUserRsp
	var err error

	rsp.User, err = mysql.UserOrm.GetById(ctx, req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
