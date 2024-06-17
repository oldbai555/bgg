package server

import (
	"context"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/bgg/singlesrv/server/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/uctx"
)

var OnceSvrImpl = &LbsingleServer{}

type LbsingleServer struct {
	client.UnimplementedLbsingleServer
}

func (a *LbsingleServer) DelFileList(ctx context.Context, req *client.DelFileListReq) (*client.DelFileListRsp, error) {
	var rsp client.DelFileListRsp
	var err error

	listRsp, err := a.GetFileList(ctx, &client.GetFileListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(lb.DefaultListOption_DefaultListOptionSelect, []string{client.FieldId_, client.FieldSortUrl_}),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, client.FieldId)
	_, err = mysql.File.NewScope(ctx).In(client.FieldId_, idList).Delete(&client.ModelFile{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 清缓存
	for _, val := range listRsp.List {
		err = cache.DelFileBySortUrl(val.SortUrl)
		if err != nil {
			log.Errorf("del %s cache failed err:%v", val.SortUrl, err)
		}
	}

	return &rsp, err
}

func (a *LbsingleServer) GetFile(ctx context.Context, req *client.GetFileReq) (*client.GetFileRsp, error) {
	var rsp client.GetFileRsp
	var err error

	var data client.ModelFile
	err = mysql.File.NewScope(ctx).Where(client.FieldId_, req.Id).First(&data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = &data

	return &rsp, err
}

func (a *LbsingleServer) GetFileList(ctx context.Context, req *client.GetFileListReq) (*client.GetFileListRsp, error) {
	var rsp client.GetFileListRsp
	var err error

	db := mysql.File.NewList(ctx, req.ListOption)
	err = lb.ProcessDefaultOptions(req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = lb.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(&rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) Login(ctx context.Context, req *client.LoginReq) (*client.LoginRsp, error) {
	var rsp client.LoginRsp
	var err error

	var user client.ModelUser
	err = mysql.User.NewScope(ctx).AndMap(map[string]interface{}{
		client.FieldUsername_: req.Username,
		client.FieldPassword_: utils.StrMd5(req.Password),
	}).First(&user)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	sid := utils.GenUUID()
	err = cache.SetLoginInfo(sid, user.ToBaseUser())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Token = sid
	rsp.User = user.ToBaseUser()

	return &rsp, err
}

func (a *LbsingleServer) Logout(ctx context.Context, req *client.LogoutReq) (*client.LogoutRsp, error) {
	var rsp client.LogoutRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = cache.DelLoginInfo(uCtx.Sid())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetLoginUser(ctx context.Context, req *client.GetLoginUserReq) (*client.GetLoginUserRsp, error) {
	var rsp client.GetLoginUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.User = uCtx.ExtInfo().(*client.BaseUser)

	return &rsp, err
}

func (a *LbsingleServer) UpdateLoginUserInfo(ctx context.Context, req *client.UpdateLoginUserInfoReq) (*client.UpdateLoginUserInfoRsp, error) {
	var rsp client.UpdateLoginUserInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	baseUser := uCtx.ExtInfo().(*client.BaseUser)

	err = mysql.User.NewScope(ctx).Eq(client.FieldId_, baseUser.Id).First(&client.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = mysql.User.NewScope(ctx).Eq(client.FieldId_, baseUser.Id).Update(utils.OrmStruct2Map(req.User, client.FieldRole))
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

func (a *LbsingleServer) ResetPassword(ctx context.Context, req *client.ResetPasswordReq) (*client.ResetPasswordRsp, error) {
	var rsp client.ResetPasswordRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	baseUser := uCtx.ExtInfo().(*client.BaseUser)

	oldPsw := utils.Md5(req.OldPassword)
	newPsw := utils.Md5(req.NewPassword)

	err = mysql.User.NewScope(ctx).AndMap(map[string]interface{}{
		client.FieldId_:       baseUser.Id,
		client.FieldPassword_: oldPsw,
	}).First(&client.ModelUser{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, client.ErrOldPasswordNotEqual
	}

	if oldPsw == newPsw {
		return nil, client.ErrOldPwdEqualNewPwd
	}

	_, err = mysql.User.NewScope(ctx).AndMap(map[string]interface{}{
		client.FieldId_: baseUser.Id,
	}).Update(map[string]interface{}{
		client.FieldPassword_: newPsw,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = cache.DelLoginInfo(uCtx.Sid())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
