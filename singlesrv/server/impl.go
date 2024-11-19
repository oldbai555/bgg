package server

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/bgg/singlesrv/server/ctx"
	"github.com/oldbai555/bgg/singlesrv/server/wsmgr"
	"github.com/oldbai555/bgg/singlesrv/server/wxminiprogram"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
	"golang.org/x/sync/singleflight"
	"os"
	"strings"
	"time"
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
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, []string{client.FieldId_, client.FieldSortUrl_, client.FieldPath_}),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	// 清缓存 和 文件
	for _, val := range listRsp.List {
		err = cache.DelFileBySortUrl(val.SortUrl)
		if err != nil {
			log.Errorf("del %s cache failed err:%v", val.SortUrl, err)
		}
		err := os.Remove(val.Path)
		if err != nil {
			log.Errorf("remove file %s failed err:%v", val.Path, err)
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

func (a *LbsingleServer) Logout(c context.Context, req *client.LogoutReq) (*client.LogoutRsp, error) {
	var rsp client.LogoutRsp
	var err error
	uCtx, err := uctx.ToUCtx(c)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	baseUser := uCtx.ExtInfo().(*client.BaseUser)

	err = cache.DelLoginInfo(uCtx.Sid())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	conn := wsmgr.GetConnByUid(baseUser.Id)
	if conn != nil {
		newCtx := ctx.NewCtx(context.Background())
		newCtx.SetExtInfo(conn)
		_, err = handleWebsocketDataTypeLogout(newCtx, wsmgr.PacketWebsocketDataByLogout())
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
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

func (a *LbsingleServer) AddUser(ctx context.Context, req *client.AddUserReq) (*client.AddUserRsp, error) {
	var rsp client.AddUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// check user
	user := req.Data
	user.Password = strings.TrimSpace(user.Password)
	if len(user.Password) < 6 || len(user.Password) > 16 {
		return nil, client.ErrPasswordLength
	}
	user.Password = utils.Md5(user.Password)

	err = OrmUser.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data.ToBaseUser()

	return &rsp, err
}

func (a *LbsingleServer) DelUserList(ctx context.Context, req *client.DelUserListReq) (*client.DelUserListRsp, error) {
	var rsp client.DelUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetUserList(ctx, &client.GetUserListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmUser.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateUser(ctx context.Context, req *client.UpdateUserReq) (*client.UpdateUserRsp, error) {
	var rsp client.UpdateUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmUser.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmUser.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map4Update(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetUser(ctx context.Context, req *client.GetUserReq) (*client.GetUserRsp, error) {
	var rsp client.GetUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmUser.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data.ToBaseUser()

	return &rsp, err
}

func (a *LbsingleServer) GetUserList(ctx context.Context, req *client.GetUserListReq) (*client.GetUserListRsp, error) {
	var rsp client.GetUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmUser.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelUser](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		AddString(client.GetUserListReq_ListOptionLikeNickname, func(val string) error {
			db.WhereLike(client.FieldNickname_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		AddString(client.GetUserListReq_ListOptionLikeUsername, func(val string) error {
			db.WhereLike(client.FieldUsername_, fmt.Sprintf("%%%s%%", val))
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

func (a *LbsingleServer) AddMpMerchantDetails(ctx context.Context, req *client.AddMpMerchantDetailsReq) (*client.AddMpMerchantDetailsRsp, error) {
	var rsp client.AddMpMerchantDetailsRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpMerchantDetails.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpMerchantDetailsList(ctx context.Context, req *client.DelMpMerchantDetailsListReq) (*client.DelMpMerchantDetailsListRsp, error) {
	var rsp client.DelMpMerchantDetailsListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpMerchantDetailsList(ctx, &client.GetMpMerchantDetailsListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpMerchantDetails.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpMerchantDetails(ctx context.Context, req *client.UpdateMpMerchantDetailsReq) (*client.UpdateMpMerchantDetailsRsp, error) {
	var rsp client.UpdateMpMerchantDetailsRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMerchantDetails.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpMerchantDetails.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpMerchantDetails(ctx context.Context, req *client.GetMpMerchantDetailsReq) (*client.GetMpMerchantDetailsRsp, error) {
	var rsp client.GetMpMerchantDetailsRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMerchantDetails.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpMerchantDetailsList(ctx context.Context, req *client.GetMpMerchantDetailsListReq) (*client.GetMpMerchantDetailsListRsp, error) {
	var rsp client.GetMpMerchantDetailsListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpMerchantDetails.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpMerchantDetails](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpMemberUser(ctx context.Context, req *client.AddMpMemberUserReq) (*client.AddMpMemberUserRsp, error) {
	var rsp client.AddMpMemberUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpMemberUser.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpMemberUserList(ctx context.Context, req *client.DelMpMemberUserListReq) (*client.DelMpMemberUserListRsp, error) {
	var rsp client.DelMpMemberUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpMemberUserList(ctx, &client.GetMpMemberUserListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpMemberUser.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpMemberUser(ctx context.Context, req *client.UpdateMpMemberUserReq) (*client.UpdateMpMemberUserRsp, error) {
	var rsp client.UpdateMpMemberUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMemberUser.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpMemberUser.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpMemberUser(ctx context.Context, req *client.GetMpMemberUserReq) (*client.GetMpMemberUserRsp, error) {
	var rsp client.GetMpMemberUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMemberUser.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpMemberUserList(ctx context.Context, req *client.GetMpMemberUserListReq) (*client.GetMpMemberUserListRsp, error) {
	var rsp client.GetMpMemberUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpMemberUser.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpMemberUser](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpUserAddress(ctx context.Context, req *client.AddMpUserAddressReq) (*client.AddMpUserAddressRsp, error) {
	var rsp client.AddMpUserAddressRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpUserAddress.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpUserAddressList(ctx context.Context, req *client.DelMpUserAddressListReq) (*client.DelMpUserAddressListRsp, error) {
	var rsp client.DelMpUserAddressListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpUserAddressList(ctx, &client.GetMpUserAddressListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpUserAddress.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpUserAddress(ctx context.Context, req *client.UpdateMpUserAddressReq) (*client.UpdateMpUserAddressRsp, error) {
	var rsp client.UpdateMpUserAddressRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpUserAddress.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpUserAddress.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpUserAddress(ctx context.Context, req *client.GetMpUserAddressReq) (*client.GetMpUserAddressRsp, error) {
	var rsp client.GetMpUserAddressRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpUserAddress.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpUserAddressList(ctx context.Context, req *client.GetMpUserAddressListReq) (*client.GetMpUserAddressListRsp, error) {
	var rsp client.GetMpUserAddressListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpUserAddress.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpUserAddress](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpUserBill(ctx context.Context, req *client.AddMpUserBillReq) (*client.AddMpUserBillRsp, error) {
	var rsp client.AddMpUserBillRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpUserBill.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpUserBillList(ctx context.Context, req *client.DelMpUserBillListReq) (*client.DelMpUserBillListRsp, error) {
	var rsp client.DelMpUserBillListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpUserBillList(ctx, &client.GetMpUserBillListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpUserBill.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpUserBill(ctx context.Context, req *client.UpdateMpUserBillReq) (*client.UpdateMpUserBillRsp, error) {
	var rsp client.UpdateMpUserBillRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpUserBill.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpUserBill.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpUserBill(ctx context.Context, req *client.GetMpUserBillReq) (*client.GetMpUserBillRsp, error) {
	var rsp client.GetMpUserBillRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpUserBill.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpUserBillList(ctx context.Context, req *client.GetMpUserBillListReq) (*client.GetMpUserBillListRsp, error) {
	var rsp client.GetMpUserBillListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpUserBill.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpUserBill](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpProductCategory(ctx context.Context, req *client.AddMpProductCategoryReq) (*client.AddMpProductCategoryRsp, error) {
	var rsp client.AddMpProductCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpProductCategory.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpProductCategoryList(ctx context.Context, req *client.DelMpProductCategoryListReq) (*client.DelMpProductCategoryListRsp, error) {
	var rsp client.DelMpProductCategoryListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpProductCategoryList(ctx, &client.GetMpProductCategoryListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpProductCategory.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpProductCategory(ctx context.Context, req *client.UpdateMpProductCategoryReq) (*client.UpdateMpProductCategoryRsp, error) {
	var rsp client.UpdateMpProductCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpProductCategory.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpProductCategory.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpProductCategory(ctx context.Context, req *client.GetMpProductCategoryReq) (*client.GetMpProductCategoryRsp, error) {
	var rsp client.GetMpProductCategoryRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpProductCategory.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpProductCategoryList(ctx context.Context, req *client.GetMpProductCategoryListReq) (*client.GetMpProductCategoryListRsp, error) {
	var rsp client.GetMpProductCategoryListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpProductCategory.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpProductCategory](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreProduct(ctx context.Context, req *client.AddMpStoreProductReq) (*client.AddMpStoreProductRsp, error) {
	var rsp client.AddMpStoreProductRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreProduct.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreProductList(ctx context.Context, req *client.DelMpStoreProductListReq) (*client.DelMpStoreProductListRsp, error) {
	var rsp client.DelMpStoreProductListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreProductList(ctx, &client.GetMpStoreProductListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreProduct.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreProduct(ctx context.Context, req *client.UpdateMpStoreProductReq) (*client.UpdateMpStoreProductRsp, error) {
	var rsp client.UpdateMpStoreProductRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProduct.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreProduct.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProduct(ctx context.Context, req *client.GetMpStoreProductReq) (*client.GetMpStoreProductRsp, error) {
	var rsp client.GetMpStoreProductRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProduct.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductList(ctx context.Context, req *client.GetMpStoreProductListReq) (*client.GetMpStoreProductListRsp, error) {
	var rsp client.GetMpStoreProductListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreProduct.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreProduct](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreProductAttr(ctx context.Context, req *client.AddMpStoreProductAttrReq) (*client.AddMpStoreProductAttrRsp, error) {
	var rsp client.AddMpStoreProductAttrRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreProductAttr.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreProductAttrList(ctx context.Context, req *client.DelMpStoreProductAttrListReq) (*client.DelMpStoreProductAttrListRsp, error) {
	var rsp client.DelMpStoreProductAttrListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreProductAttrList(ctx, &client.GetMpStoreProductAttrListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreProductAttr.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreProductAttr(ctx context.Context, req *client.UpdateMpStoreProductAttrReq) (*client.UpdateMpStoreProductAttrRsp, error) {
	var rsp client.UpdateMpStoreProductAttrRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductAttr.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreProductAttr.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductAttr(ctx context.Context, req *client.GetMpStoreProductAttrReq) (*client.GetMpStoreProductAttrRsp, error) {
	var rsp client.GetMpStoreProductAttrRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductAttr.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductAttrList(ctx context.Context, req *client.GetMpStoreProductAttrListReq) (*client.GetMpStoreProductAttrListRsp, error) {
	var rsp client.GetMpStoreProductAttrListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreProductAttr.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreProductAttr](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreProductAttrResult(ctx context.Context, req *client.AddMpStoreProductAttrResultReq) (*client.AddMpStoreProductAttrResultRsp, error) {
	var rsp client.AddMpStoreProductAttrResultRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreProductAttrResult.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreProductAttrResultList(ctx context.Context, req *client.DelMpStoreProductAttrResultListReq) (*client.DelMpStoreProductAttrResultListRsp, error) {
	var rsp client.DelMpStoreProductAttrResultListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreProductAttrResultList(ctx, &client.GetMpStoreProductAttrResultListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreProductAttrResult.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreProductAttrResult(ctx context.Context, req *client.UpdateMpStoreProductAttrResultReq) (*client.UpdateMpStoreProductAttrResultRsp, error) {
	var rsp client.UpdateMpStoreProductAttrResultRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductAttrResult.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreProductAttrResult.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductAttrResult(ctx context.Context, req *client.GetMpStoreProductAttrResultReq) (*client.GetMpStoreProductAttrResultRsp, error) {
	var rsp client.GetMpStoreProductAttrResultRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductAttrResult.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductAttrResultList(ctx context.Context, req *client.GetMpStoreProductAttrResultListReq) (*client.GetMpStoreProductAttrResultListRsp, error) {
	var rsp client.GetMpStoreProductAttrResultListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreProductAttrResult.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreProductAttrResult](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreProductAttrValue(ctx context.Context, req *client.AddMpStoreProductAttrValueReq) (*client.AddMpStoreProductAttrValueRsp, error) {
	var rsp client.AddMpStoreProductAttrValueRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreProductAttrValue.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreProductAttrValueList(ctx context.Context, req *client.DelMpStoreProductAttrValueListReq) (*client.DelMpStoreProductAttrValueListRsp, error) {
	var rsp client.DelMpStoreProductAttrValueListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreProductAttrValueList(ctx, &client.GetMpStoreProductAttrValueListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreProductAttrValue.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreProductAttrValue(ctx context.Context, req *client.UpdateMpStoreProductAttrValueReq) (*client.UpdateMpStoreProductAttrValueRsp, error) {
	var rsp client.UpdateMpStoreProductAttrValueRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductAttrValue.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreProductAttrValue.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductAttrValue(ctx context.Context, req *client.GetMpStoreProductAttrValueReq) (*client.GetMpStoreProductAttrValueRsp, error) {
	var rsp client.GetMpStoreProductAttrValueRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductAttrValue.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductAttrValueList(ctx context.Context, req *client.GetMpStoreProductAttrValueListReq) (*client.GetMpStoreProductAttrValueListRsp, error) {
	var rsp client.GetMpStoreProductAttrValueListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreProductAttrValue.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreProductAttrValue](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreProductReply(ctx context.Context, req *client.AddMpStoreProductReplyReq) (*client.AddMpStoreProductReplyRsp, error) {
	var rsp client.AddMpStoreProductReplyRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreProductReply.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreProductReplyList(ctx context.Context, req *client.DelMpStoreProductReplyListReq) (*client.DelMpStoreProductReplyListRsp, error) {
	var rsp client.DelMpStoreProductReplyListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreProductReplyList(ctx, &client.GetMpStoreProductReplyListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreProductReply.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreProductReply(ctx context.Context, req *client.UpdateMpStoreProductReplyReq) (*client.UpdateMpStoreProductReplyRsp, error) {
	var rsp client.UpdateMpStoreProductReplyRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductReply.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreProductReply.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductReply(ctx context.Context, req *client.GetMpStoreProductReplyReq) (*client.GetMpStoreProductReplyRsp, error) {
	var rsp client.GetMpStoreProductReplyRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductReply.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductReplyList(ctx context.Context, req *client.GetMpStoreProductReplyListReq) (*client.GetMpStoreProductReplyListRsp, error) {
	var rsp client.GetMpStoreProductReplyListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreProductReply.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreProductReply](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreProductRule(ctx context.Context, req *client.AddMpStoreProductRuleReq) (*client.AddMpStoreProductRuleRsp, error) {
	var rsp client.AddMpStoreProductRuleRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreProductRule.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreProductRuleList(ctx context.Context, req *client.DelMpStoreProductRuleListReq) (*client.DelMpStoreProductRuleListRsp, error) {
	var rsp client.DelMpStoreProductRuleListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreProductRuleList(ctx, &client.GetMpStoreProductRuleListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreProductRule.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreProductRule(ctx context.Context, req *client.UpdateMpStoreProductRuleReq) (*client.UpdateMpStoreProductRuleRsp, error) {
	var rsp client.UpdateMpStoreProductRuleRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductRule.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreProductRule.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductRule(ctx context.Context, req *client.GetMpStoreProductRuleReq) (*client.GetMpStoreProductRuleRsp, error) {
	var rsp client.GetMpStoreProductRuleRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreProductRule.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreProductRuleList(ctx context.Context, req *client.GetMpStoreProductRuleListReq) (*client.GetMpStoreProductRuleListRsp, error) {
	var rsp client.GetMpStoreProductRuleListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreProductRule.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreProductRule](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreShop(ctx context.Context, req *client.AddMpStoreShopReq) (*client.AddMpStoreShopRsp, error) {
	var rsp client.AddMpStoreShopRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreShop.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreShopList(ctx context.Context, req *client.DelMpStoreShopListReq) (*client.DelMpStoreShopListRsp, error) {
	var rsp client.DelMpStoreShopListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreShopList(ctx, &client.GetMpStoreShopListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreShop.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreShop(ctx context.Context, req *client.UpdateMpStoreShopReq) (*client.UpdateMpStoreShopRsp, error) {
	var rsp client.UpdateMpStoreShopRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreShop.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreShop.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreShop(ctx context.Context, req *client.GetMpStoreShopReq) (*client.GetMpStoreShopRsp, error) {
	var rsp client.GetMpStoreShopRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreShop.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreShopList(ctx context.Context, req *client.GetMpStoreShopListReq) (*client.GetMpStoreShopListRsp, error) {
	var rsp client.GetMpStoreShopListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreShop.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreShop](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpCoupon(ctx context.Context, req *client.AddMpCouponReq) (*client.AddMpCouponRsp, error) {
	var rsp client.AddMpCouponRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpCoupon.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpCouponList(ctx context.Context, req *client.DelMpCouponListReq) (*client.DelMpCouponListRsp, error) {
	var rsp client.DelMpCouponListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpCouponList(ctx, &client.GetMpCouponListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpCoupon.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpCoupon(ctx context.Context, req *client.UpdateMpCouponReq) (*client.UpdateMpCouponRsp, error) {
	var rsp client.UpdateMpCouponRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpCoupon.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpCoupon.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpCoupon(ctx context.Context, req *client.GetMpCouponReq) (*client.GetMpCouponRsp, error) {
	var rsp client.GetMpCouponRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpCoupon.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpCouponList(ctx context.Context, req *client.GetMpCouponListReq) (*client.GetMpCouponListRsp, error) {
	var rsp client.GetMpCouponListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpCoupon.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpCoupon](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpCouponUser(ctx context.Context, req *client.AddMpCouponUserReq) (*client.AddMpCouponUserRsp, error) {
	var rsp client.AddMpCouponUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpCouponUser.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpCouponUserList(ctx context.Context, req *client.DelMpCouponUserListReq) (*client.DelMpCouponUserListRsp, error) {
	var rsp client.DelMpCouponUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpCouponUserList(ctx, &client.GetMpCouponUserListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpCouponUser.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpCouponUser(ctx context.Context, req *client.UpdateMpCouponUserReq) (*client.UpdateMpCouponUserRsp, error) {
	var rsp client.UpdateMpCouponUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpCouponUser.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpCouponUser.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpCouponUser(ctx context.Context, req *client.GetMpCouponUserReq) (*client.GetMpCouponUserRsp, error) {
	var rsp client.GetMpCouponUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpCouponUser.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpCouponUserList(ctx context.Context, req *client.GetMpCouponUserListReq) (*client.GetMpCouponUserListRsp, error) {
	var rsp client.GetMpCouponUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpCouponUser.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpCouponUser](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpOrderNumber(ctx context.Context, req *client.AddMpOrderNumberReq) (*client.AddMpOrderNumberRsp, error) {
	var rsp client.AddMpOrderNumberRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpOrderNumber.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpOrderNumberList(ctx context.Context, req *client.DelMpOrderNumberListReq) (*client.DelMpOrderNumberListRsp, error) {
	var rsp client.DelMpOrderNumberListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpOrderNumberList(ctx, &client.GetMpOrderNumberListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpOrderNumber.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpOrderNumber(ctx context.Context, req *client.UpdateMpOrderNumberReq) (*client.UpdateMpOrderNumberRsp, error) {
	var rsp client.UpdateMpOrderNumberRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpOrderNumber.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpOrderNumber.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpOrderNumber(ctx context.Context, req *client.GetMpOrderNumberReq) (*client.GetMpOrderNumberRsp, error) {
	var rsp client.GetMpOrderNumberRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpOrderNumber.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpOrderNumberList(ctx context.Context, req *client.GetMpOrderNumberListReq) (*client.GetMpOrderNumberListRsp, error) {
	var rsp client.GetMpOrderNumberListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpOrderNumber.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpOrderNumber](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreOrder(ctx context.Context, req *client.AddMpStoreOrderReq) (*client.AddMpStoreOrderRsp, error) {
	var rsp client.AddMpStoreOrderRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreOrder.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreOrderList(ctx context.Context, req *client.DelMpStoreOrderListReq) (*client.DelMpStoreOrderListRsp, error) {
	var rsp client.DelMpStoreOrderListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreOrderList(ctx, &client.GetMpStoreOrderListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreOrder.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreOrder(ctx context.Context, req *client.UpdateMpStoreOrderReq) (*client.UpdateMpStoreOrderRsp, error) {
	var rsp client.UpdateMpStoreOrderRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrder.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreOrder.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreOrder(ctx context.Context, req *client.GetMpStoreOrderReq) (*client.GetMpStoreOrderRsp, error) {
	var rsp client.GetMpStoreOrderRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrder.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreOrderList(ctx context.Context, req *client.GetMpStoreOrderListReq) (*client.GetMpStoreOrderListRsp, error) {
	var rsp client.GetMpStoreOrderListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreOrder.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreOrder](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreOrderCartInfo(ctx context.Context, req *client.AddMpStoreOrderCartInfoReq) (*client.AddMpStoreOrderCartInfoRsp, error) {
	var rsp client.AddMpStoreOrderCartInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreOrderCartInfo.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreOrderCartInfoList(ctx context.Context, req *client.DelMpStoreOrderCartInfoListReq) (*client.DelMpStoreOrderCartInfoListRsp, error) {
	var rsp client.DelMpStoreOrderCartInfoListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreOrderCartInfoList(ctx, &client.GetMpStoreOrderCartInfoListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreOrderCartInfo.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreOrderCartInfo(ctx context.Context, req *client.UpdateMpStoreOrderCartInfoReq) (*client.UpdateMpStoreOrderCartInfoRsp, error) {
	var rsp client.UpdateMpStoreOrderCartInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrderCartInfo.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreOrderCartInfo.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreOrderCartInfo(ctx context.Context, req *client.GetMpStoreOrderCartInfoReq) (*client.GetMpStoreOrderCartInfoRsp, error) {
	var rsp client.GetMpStoreOrderCartInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrderCartInfo.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreOrderCartInfoList(ctx context.Context, req *client.GetMpStoreOrderCartInfoListReq) (*client.GetMpStoreOrderCartInfoListRsp, error) {
	var rsp client.GetMpStoreOrderCartInfoListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreOrderCartInfo.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreOrderCartInfo](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpStoreOrderStatus(ctx context.Context, req *client.AddMpStoreOrderStatusReq) (*client.AddMpStoreOrderStatusRsp, error) {
	var rsp client.AddMpStoreOrderStatusRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpStoreOrderStatus.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpStoreOrderStatusList(ctx context.Context, req *client.DelMpStoreOrderStatusListReq) (*client.DelMpStoreOrderStatusListRsp, error) {
	var rsp client.DelMpStoreOrderStatusListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpStoreOrderStatusList(ctx, &client.GetMpStoreOrderStatusListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpStoreOrderStatus.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpStoreOrderStatus(ctx context.Context, req *client.UpdateMpStoreOrderStatusReq) (*client.UpdateMpStoreOrderStatusRsp, error) {
	var rsp client.UpdateMpStoreOrderStatusRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrderStatus.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpStoreOrderStatus.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreOrderStatus(ctx context.Context, req *client.GetMpStoreOrderStatusReq) (*client.GetMpStoreOrderStatusRsp, error) {
	var rsp client.GetMpStoreOrderStatusRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpStoreOrderStatus.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpStoreOrderStatusList(ctx context.Context, req *client.GetMpStoreOrderStatusListReq) (*client.GetMpStoreOrderStatusListRsp, error) {
	var rsp client.GetMpStoreOrderStatusListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreOrderStatus.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpStoreOrderStatus](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpMaterial(ctx context.Context, req *client.AddMpMaterialReq) (*client.AddMpMaterialRsp, error) {
	var rsp client.AddMpMaterialRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpMaterial.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpMaterialList(ctx context.Context, req *client.DelMpMaterialListReq) (*client.DelMpMaterialListRsp, error) {
	var rsp client.DelMpMaterialListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpMaterialList(ctx, &client.GetMpMaterialListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpMaterial.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpMaterial(ctx context.Context, req *client.UpdateMpMaterialReq) (*client.UpdateMpMaterialRsp, error) {
	var rsp client.UpdateMpMaterialRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMaterial.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpMaterial.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpMaterial(ctx context.Context, req *client.GetMpMaterialReq) (*client.GetMpMaterialRsp, error) {
	var rsp client.GetMpMaterialRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMaterial.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpMaterialList(ctx context.Context, req *client.GetMpMaterialListReq) (*client.GetMpMaterialListRsp, error) {
	var rsp client.GetMpMaterialListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpMaterial.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpMaterial](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpMaterialGroup(ctx context.Context, req *client.AddMpMaterialGroupReq) (*client.AddMpMaterialGroupRsp, error) {
	var rsp client.AddMpMaterialGroupRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpMaterialGroup.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpMaterialGroupList(ctx context.Context, req *client.DelMpMaterialGroupListReq) (*client.DelMpMaterialGroupListRsp, error) {
	var rsp client.DelMpMaterialGroupListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpMaterialGroupList(ctx, &client.GetMpMaterialGroupListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpMaterialGroup.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpMaterialGroup(ctx context.Context, req *client.UpdateMpMaterialGroupReq) (*client.UpdateMpMaterialGroupRsp, error) {
	var rsp client.UpdateMpMaterialGroupRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMaterialGroup.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpMaterialGroup.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpMaterialGroup(ctx context.Context, req *client.GetMpMaterialGroupReq) (*client.GetMpMaterialGroupRsp, error) {
	var rsp client.GetMpMaterialGroupRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMaterialGroup.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpMaterialGroupList(ctx context.Context, req *client.GetMpMaterialGroupListReq) (*client.GetMpMaterialGroupListRsp, error) {
	var rsp client.GetMpMaterialGroupListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpMaterialGroup.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpMaterialGroup](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpService(ctx context.Context, req *client.AddMpServiceReq) (*client.AddMpServiceRsp, error) {
	var rsp client.AddMpServiceRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpService.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpServiceList(ctx context.Context, req *client.DelMpServiceListReq) (*client.DelMpServiceListRsp, error) {
	var rsp client.DelMpServiceListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpServiceList(ctx, &client.GetMpServiceListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpService.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpService(ctx context.Context, req *client.UpdateMpServiceReq) (*client.UpdateMpServiceRsp, error) {
	var rsp client.UpdateMpServiceRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpService.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpService.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpService(ctx context.Context, req *client.GetMpServiceReq) (*client.GetMpServiceRsp, error) {
	var rsp client.GetMpServiceRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpService.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpServiceList(ctx context.Context, req *client.GetMpServiceListReq) (*client.GetMpServiceListRsp, error) {
	var rsp client.GetMpServiceListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpService.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpService](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) AddMpShopAds(ctx context.Context, req *client.AddMpShopAdsReq) (*client.AddMpShopAdsRsp, error) {
	var rsp client.AddMpShopAdsRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = OrmMpShopAds.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = req.Data

	return &rsp, err
}

func (a *LbsingleServer) DelMpShopAdsList(ctx context.Context, req *client.DelMpShopAdsListReq) (*client.DelMpShopAdsListRsp, error) {
	var rsp client.DelMpShopAdsListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpShopAdsList(ctx, &client.GetMpShopAdsListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, client.FieldId_),
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
	_, err = OrmMpShopAds.NewBaseScope().WhereIn(client.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateMpShopAds(ctx context.Context, req *client.UpdateMpShopAdsReq) (*client.UpdateMpShopAdsRsp, error) {
	var rsp client.UpdateMpShopAdsRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpShopAds.NewBaseScope().Where(client.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpShopAds.NewBaseScope().Where(client.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpShopAds(ctx context.Context, req *client.GetMpShopAdsReq) (*client.GetMpShopAdsRsp, error) {
	var rsp client.GetMpShopAdsRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpShopAds.NewBaseScope().Where(client.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}

func (a *LbsingleServer) GetMpShopAdsList(ctx context.Context, req *client.GetMpShopAdsListReq) (*client.GetMpShopAdsListRsp, error) {
	var rsp client.GetMpShopAdsListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpShopAds.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpShopAds](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetMpShopAdsListPublic(ctx context.Context, req *client.GetMpShopAdsListPublicReq) (*client.GetMpShopAdsListPublicRsp, error) {
	var rsp client.GetMpShopAdsListPublicRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpShopAds.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*client.ModelMpShopAds](req.ListOption, db)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	var list []*client.ModelMpShopAds
	_, err = db.FindPaginate(uCtx, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if len(rsp.List) == 0 {
		rsp.List = append(rsp.List, "https://oldbai.top/oss/download/BUOZ74", "https://oldbai.top/oss/download/BUOZ74")
	}
	rsp.IsActive = true
	return &rsp, err
}

func (a *LbsingleServer) WxMiniProgramAuthSession(c context.Context, req *client.WxMiniProgramAuthSessionReq) (*client.WxMiniProgramAuthSessionRsp, error) {
	var rsp client.WxMiniProgramAuthSessionRsp
	var err error

	conf, err := syscfg.GetWxMiniProgramConf()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	session, err := wxminiprogram.Code2Session(&client.JsCodeToSessionReq{
		JsCode: req.Code,
		Appid:  conf.AppId,
		Secret: conf.Secret,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Openid = session.Openid
	rsp.UserInfo, err = OrmMpMemberUser.NewBaseScope().Where(client.FieldMpOpenid_, session.Openid).First(ctx.NewCtx(c))
	if err != nil {
		log.Errorf("err:%v", err)
		err = nil
	}
	if rsp.UserInfo != nil {
		sid := utils.GenUUID()
		err = cache.SetLoginInfo(sid, rsp.UserInfo.ToBaseUser())
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		rsp.AccessToken = sid
	}

	return &rsp, err
}

func (a *LbsingleServer) WxMPAuthSendSmsCode(ctx context.Context, req *client.WxMPAuthSendSmsCodeReq) (*client.WxMPAuthSendSmsCodeRsp, error) {
	var rsp client.WxMPAuthSendSmsCodeRsp
	var err error

	code, err := cache.GetMpSmsCode(fmt.Sprintf("%s", req.Mobile))
	if err == nil && code != "" {
		return nil, client.ErrMpSmsCodeNoExpired
	}

	rsp.Code = utils.GetRandomString(6, utils.RandomStringModNumber)
	err = cache.SetMpSmsCode(fmt.Sprintf("%s", req.Mobile), rsp.Code)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) WxMPAuthSendSmsLogin(c context.Context, req *client.WxMPAuthSendSmsLoginReq) (*client.WxMPAuthSendSmsLoginRsp, error) {
	var rsp client.WxMPAuthSendSmsLoginRsp
	var err error

	uCtx, err := uctx.ToUCtx(c)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	ctx2, ok := uCtx.(*ctx.Ctx)
	if !ok {
		return nil, client.ErrMpCtxConvertFailed
	}

	code, err := cache.GetMpSmsCode(fmt.Sprintf("%s", req.Mobile))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if code != req.Code {
		return nil, client.ErrMpSmsCodeNotEqual
	}
	cache.DelMpSmsCode(fmt.Sprintf("%s", req.Mobile))

	mpUser := &client.ModelMpMemberUser{
		Nickname:    "微信用户_" + utils.GenUUID(),
		Mobile:      req.Mobile,
		RegisterIp:  ctx2.ClientIp,
		LastLoginAt: utils.TimeNow(),
		LastLoginIp: ctx2.ClientIp,
		LoginType:   req.From,
		MpOpenid:    req.Openid,
	}

	isEmpty, err := OrmMpMemberUser.NewBaseScope().Where(client.FieldMobile_, mpUser.Mobile).IsEmpty(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if isEmpty {
		err = OrmMpMemberUser.NewBaseScope().Create(uCtx, mpUser)
	}
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.UserInfo = mpUser
	sid := utils.GenUUID()
	err = cache.SetLoginInfo(sid, rsp.UserInfo.ToBaseUser())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.AccessToken = sid

	return &rsp, err
}

func (a *LbsingleServer) MPShopNearBy(ctx context.Context, req *client.MPShopNearByReq) (*client.MPShopNearByRsp, error) {
	var rsp client.MPShopNearByRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	db := OrmMpStoreShop.NewBaseScope()
	if req.ShopId != 0 {
		db.Where(client.FieldId_, req.ShopId)
	}
	if req.Kw != "" {
		db.WhereLike(client.FieldName_, fmt.Sprintf("%%%s%%", req.Kw))
	}
	rsp.Shop, err = db.Find(uCtx)
	return &rsp, err
}

func (a *LbsingleServer) MPShopProduct(ctx context.Context, req *client.MPShopProductReq) (*client.MPShopProductRsp, error) {
	var rsp client.MPShopProductRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreShop.NewBaseScope()
	if req.ShopId != 0 {
		db.Where(client.FieldId_, req.ShopId)
	}
	storeShop, err := db.First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.List, err = OrmMpProductCategory.NewBaseScope().Where(client.FieldMpStoreShopId_, storeShop.Id).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if len(rsp.List) == 0 {
		return &rsp, nil
	}

	rsp.ProductMap = make(map[uint64]*client.MPShopProductRsp_AppStoreProduct)
	for _, category := range rsp.List {
		var product client.MPShopProductRsp_AppStoreProduct
		find, err := OrmMpStoreProduct.NewBaseScope().Where(client.FieldMpStoreShopId_, storeShop.Id).Where(client.FieldCateId_, category.Id).Find(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}

		product.List = find
		product.ProductAttr = make(map[uint64]*client.MPShopProductRsp_AppStoreProductAttr)
		product.ProductAttrValue = make(map[uint64]*client.MPShopProductRsp_StoreProductAttrValueDO)
		for _, storeProduct := range find {
			find2, err := OrmMpStoreProductAttr.NewBaseScope().Where(client.FieldProductId_, storeProduct.Id).Find(uCtx)
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}
			find3, err := OrmMpStoreProductAttrValue.NewBaseScope().Where(client.FieldProductId_, storeProduct.Id).Find(uCtx)
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}
			product.ProductAttr[storeProduct.Id] = &client.MPShopProductRsp_AppStoreProductAttr{
				List: find2,
			}
			product.ProductAttrValue[storeProduct.Id] = &client.MPShopProductRsp_StoreProductAttrValueDO{
				ProductValue: utils.Slice2MapKeyByStructField(find3, client.FieldSku).(map[string]*client.ModelMpStoreProductAttrValue),
			}
		}
		rsp.ProductMap[category.Id] = &product
	}

	return &rsp, err
}

func (a *LbsingleServer) MPShopCouponCount(ctx context.Context, req *client.MPShopCouponCountReq) (*client.MPShopCouponCountRsp, error) {
	var rsp client.MPShopCouponCountRsp
	var err error
	return &rsp, err
}

func (a *LbsingleServer) MPShopCouponList(ctx context.Context, req *client.MPShopCouponListReq) (*client.MPShopCouponListRsp, error) {
	var rsp client.MPShopCouponListRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPShopCouponMine(ctx context.Context, req *client.MPShopCouponMineReq) (*client.MPShopCouponMineRsp, error) {
	var rsp client.MPShopCouponMineRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPShopCouponReceive(ctx context.Context, req *client.MPShopCouponReceiveReq) (*client.MPShopCouponReceiveRsp, error) {
	var rsp client.MPShopCouponReceiveRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPShopOrderCreate(ctx context.Context, req *client.MPShopOrderCreateReq) (*client.MPShopOrderCreateRsp, error) {
	var rsp client.MPShopOrderCreateRsp
	var err error

	var (
		sumPrice     int64
		couponPrice  int64
		postagePrice int64
		totalNum     uint32
	)

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	user, ok := uCtx.ExtInfo().(*client.BaseUser)
	if !ok {
		return nil, client.ErrUserNotFound
	}

	storeShop, err := OrmMpStoreShop.NewBaseScope().Where(client.FieldId_, req.ShopId).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if storeShop.DeliveryPrice > 0 {
		postagePrice = storeShop.DeliveryPrice
	}

	// todo 加锁

	specs := req.Spec
	productIds := req.ProductId
	numbers := req.Number
	couponId := req.CouponId

	// 计算商品价格
	for i, productId := range productIds {
		newSku := strings.Replace(specs[i], "|", ",", -1)
		storeProduct, err := OrmMpStoreProductAttrValue.NewBaseScope().Where(map[string]interface{}{
			client.FieldProductId_: productId,
			client.FieldSku_:       newSku,
		}).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		num := numbers[i]
		totalNum += num
		sumPrice += int64(num) * storeProduct.Price
	}

	// todo 校验优惠卷
	if couponId != 0 {
		_, err := OrmMpCouponUser.NewBaseScope().Where(client.FieldCouponId_, couponId).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		coupon, err := OrmMpCoupon.NewBaseScope().Where(client.FieldId_, couponId).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		couponPrice = coupon.Value
	}

	// 计算最终价格
	var payPrice = sumPrice - couponPrice
	if req.OrderType == "takeout" {
		payPrice += postagePrice
	}

	// 奖励积分
	var gainIntegral int64
	for _, productId := range productIds {
		product, err := OrmMpStoreProduct.NewBaseScope().Where(client.FieldId_, productId).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		gainIntegral += product.GiveIntegral
	}

	// 生成一个订单号
	orderSn := GenerateOrderID()

	// 取餐表
	orderNumber := &client.ModelMpOrderNumber{
		OrderSn: orderSn,
	}
	err = OrmMpOrderNumber.NewBaseScope().Create(uCtx, orderNumber)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	order := &client.ModelMpStoreOrder{
		OrderSn:       orderSn,
		MpUid:         user.Id,
		NumberId:      int64(orderNumber.Id),
		MpStoreShopId: storeShop.Id,
		GetAt:         uint32(time.Now().Add(time.Duration(req.GetGetTime()) * time.Minute).Unix()),
		TotalNum:      totalNum,
		CartId:        "",
		TotalPrice:    sumPrice,
		TotalPostage:  postagePrice,
		CouponId:      req.CouponId,
		CouponPrice:   couponPrice,
		PayPrice:      payPrice,
		PayPostage:    postagePrice,
		PayType:       req.PayType,
		GainIntegral:  gainIntegral,
		Mark:          req.Remark,
		ShippingType:  2, // 默认门店自提
		OrderType:     req.OrderType,
	}

	if req.OrderType == "takeout" {
		userAddress, err := OrmMpUserAddress.NewBaseScope().Where(client.FieldMpUid_, user.Id).Where(client.FieldId_, req.AddressId).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		order.UserPhone = userAddress.Phone
		order.RealName = userAddress.RealName
		order.UserAddress = userAddress.GetAddress() + " " + userAddress.GetDetail()
	}

	err = OrmMpStoreOrder.NewBaseScope().Create(uCtx, order)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 减库存
	for i, productId := range productIds {
		newSku := strings.Replace(specs[i], "|", ",", -1)
		skuVal, err := OrmMpStoreProductAttrValue.NewBaseScope().Where(map[string]interface{}{
			client.FieldProductId_: productId,
			client.FieldSku_:       newSku,
		}).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		res, err := OrmMpStoreProductAttrValue.NewBaseScope().Where(client.FieldId_, skuVal.Id).Where(fmt.Sprintf("%s > %d", client.FieldStock_, numbers[i])).Update(uCtx, map[string]interface{}{
			client.FieldStock_: skuVal.Stock - int32(numbers[i]),
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		if res.RowsAffected == 0 {
			return nil, lberr.NewInvalidArg("decrease stock failed %d %d", skuVal.Id, numbers[i])
		}

		product, err := OrmMpStoreProduct.NewBaseScope().Where(client.FieldId_, productId).First(uCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		// 添加进入购物车
		var shopCar = &client.ModelMpStoreOrderCartInfo{
			MpOrderId:    order.Id,
			OrderSn:      orderSn,
			CartId:       0,
			ProductId:    productId,
			CartInfo:     "",
			Unique:       utils.GenUUID(),
			IsAfterSales: 1,
			Title:        product.Name,
			Image:        product.Image,
			Number:       numbers[i],
			Spec:         specs[i],
			Price:        skuVal.Price,
		}
		err = OrmMpStoreOrderCartInfo.NewBaseScope().Create(uCtx, shopCar)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}

	// 订单状态
	var orderStatus = &client.ModelMpStoreOrderStatus{
		Oid:           order.Id,
		ChangeType:    "yshop_create_order",
		ChangeMessage: "订单生成",
	}
	err = OrmMpStoreOrderStatus.NewBaseScope().Create(uCtx, orderStatus)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// todo 如果超时需要取消订单

	return &rsp, err
}

func (a *LbsingleServer) MPShopOrderList(ctx context.Context, req *client.MPShopOrderListReq) (*client.MPShopOrderListRsp, error) {
	var rsp client.MPShopOrderListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpStoreOrder.NewList(&core.ListOption{
		Offset:    uint32((req.Page - 1) * req.Limit),
		Limit:     uint32(req.Limit),
		SkipTotal: true,
	})
	if req.Type >= 0 {
		db.Where(client.FieldStatus_, req.Type)
	}
	list, err := db.Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.List = list

	mpStoreShopIds := utils.PluckUint64List(rsp.List, client.FieldMpStoreShopId)
	mpStoreShopList, err := OrmMpStoreShop.WhereIn(client.FieldId_, mpStoreShopIds).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.ShopMap = utils.Slice2MapKeyByStructField(mpStoreShopList, client.FieldId).(map[uint64]*client.ModelMpStoreShop)

	orderSnList := utils.PluckStringList(rsp.List, client.FieldOrderSn)
	cartInfoList, err := OrmMpStoreOrderCartInfo.WhereIn(client.FieldOrderSn_, orderSnList).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.CartMap = make(map[string]*client.MPShopOrderListRsp_CartInfo)
	for _, info := range cartInfoList {
		cartInfo, ok := rsp.CartMap[info.OrderSn]
		if !ok {
			rsp.CartMap[info.OrderSn] = &client.MPShopOrderListRsp_CartInfo{}
			cartInfo = rsp.CartMap[info.OrderSn]
		}
		cartInfo.List = append(cartInfo.List, info)
	}

	return &rsp, err
}

func (a *LbsingleServer) MPShopOrderDetail(ctx context.Context, req *client.MPShopOrderDetailReq) (*client.MPShopOrderDetailRsp, error) {
	var rsp client.MPShopOrderDetailRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	order, err := OrmMpStoreOrder.Where(client.FieldOrderSn_, req.OrderSn).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	mpStoreShop, err := OrmMpStoreShop.Where(client.FieldId_, order.MpStoreShopId).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	cartInfoList, err := OrmMpStoreOrderCartInfo.Where(client.FieldOrderSn_, order.OrderSn).Find(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Order = order
	rsp.ShopInfo = mpStoreShop
	rsp.CartList = cartInfoList
	return &rsp, err
}

func (a *LbsingleServer) MPShopOrderTake(ctx context.Context, req *client.MPShopOrderTakeReq) (*client.MPShopOrderTakeRsp, error) {
	var rsp client.MPShopOrderTakeRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPShopOrderRefund(ctx context.Context, req *client.MPShopOrderRefundReq) (*client.MPShopOrderRefundRsp, error) {
	var rsp client.MPShopOrderRefundRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPShopOrderPay(ctx context.Context, req *client.MPShopOrderPayReq) (*client.MPShopOrderPayRsp, error) {
	var rsp client.MPShopOrderPayRsp
	var err error
	rsp.Status = "ok"
	return &rsp, err
}

func (a *LbsingleServer) MPWechatConfig(ctx context.Context, req *client.MPWechatConfigReq) (*client.MPWechatConfigRsp, error) {
	var rsp client.MPWechatConfigRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserAddressList(ctx context.Context, req *client.MPUserAddressListReq) (*client.MPUserAddressListRsp, error) {
	var rsp client.MPUserAddressListRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserAddressDefault(ctx context.Context, req *client.MPUserAddressDefaultReq) (*client.MPUserAddressDefaultRsp, error) {
	var rsp client.MPUserAddressDefaultRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserAddressAddAndEdit(ctx context.Context, req *client.MPUserAddressAddAndEditReq) (*client.MPUserAddressAddAndEditRsp, error) {
	var rsp client.MPUserAddressAddAndEditRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserAddressDel(ctx context.Context, req *client.MPUserAddressDelReq) (*client.MPUserAddressDelRsp, error) {
	var rsp client.MPUserAddressDelRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserInfo(ctx context.Context, req *client.MPUserInfoReq) (*client.MPUserInfoRsp, error) {
	var rsp client.MPUserInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	user, ok := uCtx.ExtInfo().(*client.BaseUser)
	if !ok {
		return nil, client.ErrUserNotFound
	}
	rsp.UserInfo = user

	return &rsp, err
}

func (a *LbsingleServer) MPUserMineService(ctx context.Context, req *client.MPUserMineServiceReq) (*client.MPUserMineServiceRsp, error) {
	var rsp client.MPUserMineServiceRsp
	var err error
	rsp.List = append(rsp.List, &client.ModelMpService{
		Id:     1,
		Status: 1,
		Type:   "call",
		Phone:  "14788777898",
		Name:   "电话",
		Image:  "https://oldbai.top/oss/download/BUOZ74",
	})
	return &rsp, err
}

func (a *LbsingleServer) MPUserServiceContent(ctx context.Context, req *client.MPUserServiceContentReq) (*client.MPUserServiceContentRsp, error) {
	var rsp client.MPUserServiceContentRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserSaveInfo(ctx context.Context, req *client.MPUserSaveInfoReq) (*client.MPUserSaveInfoRsp, error) {
	var rsp client.MPUserSaveInfoRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserBillList(ctx context.Context, req *client.MPUserBillListReq) (*client.MPUserBillListRsp, error) {
	var rsp client.MPUserBillListRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserMoneyList(ctx context.Context, req *client.MPUserMoneyListReq) (*client.MPUserMoneyListRsp, error) {
	var rsp client.MPUserMoneyListRsp
	var err error

	return &rsp, err
}

func (a *LbsingleServer) MPUserRecharge(ctx context.Context, req *client.MPUserRechargeReq) (*client.MPUserRechargeRsp, error) {
	var rsp client.MPUserRechargeRsp
	var err error

	return &rsp, err
}
