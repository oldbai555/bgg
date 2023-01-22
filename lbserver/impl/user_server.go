package impl

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/lbconst"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/bgg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"

	"time"
)

var lbuserServer LbuserServer

type LbuserServer struct {
	*lbuser.UnimplementedLbuserServer
}

func (a *LbuserServer) Login(ctx context.Context, req *lbuser.LoginReq) (*lbuser.LoginRsp, error) {
	var rsp lbuser.LoginRsp
	var user lbuser.ModelUser
	err := UserOrm.NewScope().Eq(lbuser.FieldUsername_, req.Username).First(ctx, &user)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	if user.Password != utils.StrMd5(req.Password) {
		return nil, ErrPasswordInvalid
	}

	genToken, err := webtool.GenToken(ctx, user.Id)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	sid := utils.GenUUID()
	err = lb.Rdb.Set(ctx, sid, genToken, time.Hour).Err()
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	err = lb.Rdb.SetJson(ctx, fmt.Sprintf("%d", user.Id), user, time.Hour)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.Sid = sid
	return &rsp, nil
}

func (a *LbuserServer) Logout(ctx context.Context, req *lbuser.LogoutReq) (*lbuser.LogoutRsp, error) {
	var rsp lbuser.LogoutRsp

	err := lb.Rdb.Del(ctx, req.Sid).Err()
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (a *LbuserServer) GetLoginUser(ctx context.Context, req *lbuser.GetLoginUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp
	var user lbuser.ModelUser

	claims, ok := ctx.Value(CtxWithClaim).(*webtool.Claims)
	if !ok {
		log.Errorf("err is : %v", ErrGetLoginFail)
		return nil, ErrGetLoginFail
	}

	err := UserOrm.NewScope().Eq(lbuser.FieldId_, claims.UserId).First(ctx, &user)
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

func (a *LbuserServer) UpdateLoginUserInfo(ctx context.Context, req *lbuser.UpdateLoginUserInfoReq) (*lbuser.UpdateLoginUserInfoRsp, error) {
	var rsp lbuser.UpdateLoginUserInfoRsp

	claims, ok := ctx.Value(CtxWithClaim).(*webtool.Claims)
	if !ok {
		log.Errorf("err is : %v", ErrGetLoginFail)
		return nil, ErrGetLoginFail
	}

	var updateMap = map[string]interface{}{
		lbuser.FieldAvatar_:   req.User.Avatar,
		lbuser.FieldEmail_:    req.User.Email,
		lbuser.FieldGithub_:   req.User.Github,
		lbuser.FieldNickname_: req.User.Nickname,
		lbuser.FieldDesc:      req.User.Desc,
	}
	err := UserOrm.NewScope().Eq(lbuser.FieldId_, claims.UserId).Update(ctx, updateMap)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) AddUser(ctx context.Context, req *lbuser.AddUserReq) (*lbuser.AddUserRsp, error) {
	var rsp lbuser.AddUserRsp

	isEmpty, err := UserOrm.NewScope().Eq(lbuser.FieldUsername_, req.User.Username).IsEmpty(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	if !isEmpty {
		return nil, lberr.NewErr(int32(lbuser.ErrCode_ErrUserExit), "用户名已存在")
	}

	req.User.Password = utils.StrMd5(req.User.Password)
	err = UserOrm.NewScope().Create(ctx, req.User)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetUserList(ctx context.Context, req *lbuser.GetUserListReq) (*lbuser.GetUserListRsp, error) {
	var rsp lbuser.GetUserListRsp

	db := UserOrm.NewList(req.ListOption).OrderByDesc(lbuser.FieldId_)
	err := lbconst.NewListOptionProcessor(req.ListOption).
		AddString(lbuser.GetUserListReq_ListOptionLikeUsername, func(val string) error {
			db.Like(lbuser.FieldUsername_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		Process()
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	rsp.Page, err = db.FindPage(ctx, &rsp.List)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	return &rsp, nil
}
func (a *LbuserServer) DelUser(ctx context.Context, req *lbuser.DelUserReq) (*lbuser.DelUserRsp, error) {
	var rsp lbuser.DelUserRsp

	err := UserOrm.NewScope().Eq(lbuser.FieldId_, req.Id).Delete(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetUser(ctx context.Context, req *lbuser.GetUserReq) (*lbuser.GetUserRsp, error) {
	var rsp lbuser.GetUserRsp

	err := UserOrm.NewScope().Eq(lbuser.FieldId_, req.Id).First(ctx, &rsp.User)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) UpdateUserNameWithRole(ctx context.Context, req *lbuser.UpdateUserNameWithRoleReq) (*lbuser.UpdateUserNameWithRoleRsp, error) {
	var rsp lbuser.UpdateUserNameWithRoleRsp

	isEmpty, err := UserOrm.NewScope().NotEq(lbuser.FieldId_, req.Id).Eq(lbuser.FieldUsername_, req.Username).IsEmpty(ctx)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	if !isEmpty {
		return nil, lberr.NewErr(int32(lbuser.ErrCode_ErrUserExit), "用户名已存在")
	}

	UserOrm.NewScope().Eq(lbuser.FieldId_, req.Id).Update(ctx, map[string]interface{}{
		lbuser.FieldUsername_: req.Username,
		lbuser.FieldRole_:     req.Role,
	})
	return &rsp, nil
}
func (a *LbuserServer) ResetPassword(ctx context.Context, req *lbuser.ResetPasswordReq) (*lbuser.ResetPasswordRsp, error) {
	var rsp lbuser.ResetPasswordRsp

	req.NewPassword = utils.StrMd5(req.NewPassword)

	err := UserOrm.NewScope().Eq(lbuser.FieldId_, req.Id).Eq(lbuser.FieldPassword_, req.OldPassword).Update(ctx, map[string]interface{}{
		lbuser.FieldPassword_: req.NewPassword,
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &rsp, nil
}
func (a *LbuserServer) GetFrontUser(ctx context.Context, req *lbuser.GetFrontUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp

	var user lbuser.ModelUser
	err := UserOrm.NewScope().Eq(lbuser.FieldId_, 4).First(ctx, &user)
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
