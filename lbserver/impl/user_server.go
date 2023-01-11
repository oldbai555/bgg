package impl

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/bgg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"

	"time"
)

var lbuserServer LbuserServer

type LbuserServer struct {
	*lbuser.UnimplementedLbuserServer
}

func (u *LbuserServer) Login(ctx context.Context, req *lbuser.LoginReq) (*lbuser.LoginRsp, error) {
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

func (u *LbuserServer) Logout(ctx context.Context, req *lbuser.LogoutReq) (*lbuser.LogoutRsp, error) {
	var rsp lbuser.LogoutRsp

	err := lb.Rdb.Del(ctx, req.Sid).Err()
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	return &rsp, nil
}

func (u *LbuserServer) GetLoginUser(ctx context.Context, req *lbuser.GetLoginUserReq) (*lbuser.GetLoginUserRsp, error) {
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
