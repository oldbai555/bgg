package service

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/bgg/service/lbuserserver/impl/cache"
	"github.com/oldbai555/bgg/service/lbuserserver/impl/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/uctx"
)

var OnceSvrImpl = &LbuserServer{}

type LbuserServer struct {
	lbuser.UnimplementedLbuserServer
}

func (a *LbuserServer) Login(ctx context.Context, req *lbuser.LoginReq) (*lbuser.LoginRsp, error) {
	var rsp lbuser.LoginRsp
	var err error

	var user lbuser.ModelUser
	err = mysql.User.NewScope(ctx).Eq(lbuser.FieldUsername_, req.Username).Eq(lbuser.FieldPassword_, utils.StrMd5(req.Password)).First(&user)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	sid := utils.GenRandomStr()
	err = cache.SetLoginInfo(sid, user.ToBaseUser())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Sid = sid

	return &rsp, err
}
func (a *LbuserServer) Logout(ctx context.Context, req *lbuser.LogoutReq) (*lbuser.LogoutRsp, error) {
	var rsp lbuser.LogoutRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	sid := uCtx.Sid()
	err = cache.DelLoginInfo(sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbuserServer) GetLoginUser(ctx context.Context, req *lbuser.GetLoginUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp
	var err error

	sid := req.Sid
	baseUser, err := cache.GetLoginInfo(sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.BaseUser = baseUser

	return &rsp, err
}
func (a *LbuserServer) UpdateLoginUserInfo(ctx context.Context, req *lbuser.UpdateLoginUserInfoReq) (*lbuser.UpdateLoginUserInfoRsp, error) {
	var rsp lbuser.UpdateLoginUserInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	baseUser := uCtx.ExtInfo().(*lbuser.BaseUser)

	if req.User.Id != baseUser.Id {
		return nil, lbuser.ErrLoginUserNotEqualChangeUser
	}

	_, err = mysql.User.NewScope(ctx).Where(lbuser.FieldId_, baseUser.Id).Update(utils.OrmStruct2Map(req.User))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = cache.SetLoginInfo(uCtx.Sid(), req.User)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbuserServer) ResetPassword(ctx context.Context, req *lbuser.ResetPasswordReq) (*lbuser.ResetPasswordRsp, error) {
	var rsp lbuser.ResetPasswordRsp
	var err error

	if !checkPassword(req.NewPassword) {
		return nil, lbuser.ErrPasswordInvalid
	}

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	baseUser := uCtx.ExtInfo().(*lbuser.BaseUser)

	_, err = mysql.User.NewScope(ctx).Eq(lbuser.FieldId_, baseUser.Id).Eq(lbuser.FieldPassword_, utils.StrMd5(req.OldPassword)).Update(map[string]interface{}{
		lbuser.FieldPassword_: utils.StrMd5(req.NewPassword),
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbuserServer) GetFrontUser(ctx context.Context, req *lbuser.GetFrontUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp
	var err error

	var user lbuser.ModelUser
	err = mysql.User.NewScope(ctx).Gt(lbuser.FieldId_, 0).First(&user)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.BaseUser = user.ToBaseUser()

	return &rsp, err
}
func (a *LbuserServer) AddUser(ctx context.Context, req *lbuser.AddUserReq) (*lbuser.AddUserRsp, error) {
	var rsp lbuser.AddUserRsp
	var err error

	if !checkPassword(req.Data.Password) {
		return nil, lbuser.ErrPasswordInvalid
	}

	if !CheckStr(req.Data.Username) {
		return nil, lbuser.ErrUserNameInvalid
	}

	var user lbuser.ModelUser
	err = mysql.User.NewScope(ctx).Eq(lbuser.FieldUsername_, req.Data.Username).First(&user)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, lbuser.ErrUserExit
	}

	req.Data.Password = utils.StrMd5(req.Data.Password)
	_, err = mysql.User.NewScope(ctx).Create(req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbuserServer) DelUserList(ctx context.Context, req *lbuser.DelUserListReq) (*lbuser.DelUserListRsp, error) {
	var rsp lbuser.DelUserListRsp
	var err error

	listRsp, err := a.GetUserList(ctx, &lbuser.GetUserListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(lb.DefaultListOption_DefaultListOptionSelect, lbuser.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbuser.FieldId)
	_, err = mysql.User.NewScope(ctx).In(lbuser.FieldId, idList).Delete(&lbuser.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbuserServer) UpdateUser(ctx context.Context, req *lbuser.UpdateUserReq) (*lbuser.UpdateUserRsp, error) {
	var rsp lbuser.UpdateUserRsp
	var err error

	var data lbuser.ModelUser
	err = mysql.User.NewScope(ctx).Where(lbuser.FieldId_, req.Data.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = mysql.User.NewScope(ctx).Where(lbuser.FieldId_, data.Id).Update(utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbuserServer) GetUser(ctx context.Context, req *lbuser.GetUserReq) (*lbuser.GetUserRsp, error) {
	var rsp lbuser.GetUserRsp
	var err error

	var data lbuser.ModelUser
	err = mysql.User.NewScope(ctx).Where(lbuser.FieldId_, req.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = &data

	return &rsp, err
}
func (a *LbuserServer) GetUserList(ctx context.Context, req *lbuser.GetUserListReq) (*lbuser.GetUserListRsp, error) {
	var rsp lbuser.GetUserListRsp
	var err error

	db := mysql.User.NewList(ctx, req.ListOption)
	err = lb.ProcessDefaultOptions(req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = lb.NewOptionsProcessor(req.ListOption).
		AddString(lbuser.GetUserListReq_ListOptionLikeUsername, func(val string) error {
			db.Like(lbuser.FieldUsername_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Paginate, err = db.FindPaginate(&rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
