package server

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
	"golang.org/x/sync/singleflight"
)

var OnceSvrImpl = &LbsingleServer{}

type LbsingleServer struct {
	client.UnimplementedLbsingleServer
	singleGroup singleflight.Group
}

func (a *LbsingleServer) DelFileList(ctx context.Context, req *client.DelFileListReq) (*client.DelFileListRsp, error) {
	var rsp client.DelFileListRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetFileList(ctx, &client.GetFileListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, []string{client.FieldId_, client.FieldSortUrl_}),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	// 清缓存
	for _, val := range listRsp.List {
		err = cache.DelFileBySortUrl(val.SortUrl)
		if err != nil {
			log.Errorf("del %s cache failed err:%v", val.SortUrl, err)
		}
	}

	idList := utils.PluckUint64List(listRsp.List, client.FieldId)
	_, err = OrmFile.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetFile(ctx context.Context, req *client.GetFileReq) (*client.GetFileRsp, error) {
	var rsp client.GetFileRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmFile.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetFileList(ctx context.Context, req *client.GetFileListReq) (*client.GetFileListRsp, error) {
	var rsp client.GetFileListRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmFile.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelFile](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		AddString(client.GetFileListReq_ListOptionLikeName, func(val string) error {
			db.WhereLike(client.FieldName_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) Login(ctx context.Context, req *client.LoginReq) (*client.LoginRsp, error) {
	var rsp client.LoginRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	var user *client.ModelUser
	user, err = OrmUser.NewBaseScope().Where(map[string]interface{}{
		client.FieldUsername_: req.Username,
		client.FieldPassword_: utils.StrMd5(req.Password),
	}).First(uCtx)
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

	_, err = OrmUser.NewBaseScope().Where(client.FieldId_, baseUser.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmUser.NewBaseScope().Where(client.FieldId_, baseUser.Id).Update(uCtx, utils.OrmStruct2Map(req.User, client.FieldRole))
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

	_, err = OrmUser.NewBaseScope().Where(map[string]interface{}{
		client.FieldId_:       baseUser.Id,
		client.FieldPassword_: oldPsw,
	}).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, client.ErrOldPasswordNotEqual
	}

	if oldPsw == newPsw {
		return nil, client.ErrOldPwdEqualNewPwd
	}

	_, err = OrmUser.NewBaseScope().Where(map[string]interface{}{
		client.FieldId_: baseUser.Id,
	}).Update(uCtx, map[string]interface{}{
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
func (a *LbsingleServer) SyncFile(ctx context.Context, _ *client.SyncFileReq) (*client.SyncFileRsp, error) {
	var rsp client.SyncFileRsp
	var err error

	// 单体环境 直接用单机的 singleflight 防并发
	_, err, _ = a.singleGroup.Do("syncfile", func() (interface{}, error) {
		err := SyncFileIndex(ctx, true)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
