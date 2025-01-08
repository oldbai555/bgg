package lbsingleserver

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/bgg/service/lbsingleserver/wsmgr"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
	"golang.org/x/sync/singleflight"
	"strings"
)

var OnceSvrImpl = &LbsingleServer{}

type LbsingleServer struct {
	lbsingle.UnimplementedLbsingleServer
	singleGroup singleflight.Group
}

func (a *LbsingleServer) GetFileList(ctx context.Context, req *lbsingle.GetFileListReq) (*lbsingle.GetFileListRsp, error) {
	var rsp lbsingle.GetFileListRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmFile.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbsingle.ModelFile](req.ListOption, db)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = core.NewOptionsProcessor(req.ListOption).
		AddString(lbsingle.GetFileListReq_ListOptionLikeName, func(val string) error {
			db.WhereLike(lbsingle.FieldName_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) Login(ctx context.Context, req *lbsingle.LoginReq) (*lbsingle.LoginRsp, error) {
	var rsp lbsingle.LoginRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	var user *lbsingle.ModelUser
	user, err = OrmUser.NewBaseScope().Where(map[string]interface{}{
		lbsingle.FieldUsername_: req.Username,
		lbsingle.FieldPassword_: utils.StrMd5(req.Password),
	}).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	sid := utils.GenUUID()
	err = cache.SetLoginInfo(sid, user.ToBaseUser())
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Token = sid
	rsp.User = user.ToBaseUser()

	return &rsp, err
}

func (a *LbsingleServer) Logout(c context.Context, req *lbsingle.LogoutReq) (*lbsingle.LogoutRsp, error) {
	var rsp lbsingle.LogoutRsp
	var err error
	uCtx, err := uctx.ToUCtx(c)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	baseUser := uCtx.ExtInfo().(*lbbase.BaseUser)

	err = cache.DelLoginInfo(uCtx.Sid())
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	conn := wsmgr.GetConnByUid(baseUser.Id)
	if conn != nil {
		newCtx := bctx.NewCtx(context.Background())
		newCtx.SetExtInfo(conn)
		_, err = handleWebsocketDataTypeLogout(newCtx, wsmgr.PacketWebsocketDataByLogout())
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}

	return &rsp, err
}

func (a *LbsingleServer) GetLoginUser(ctx context.Context, req *lbsingle.GetLoginUserReq) (*lbsingle.GetLoginUserRsp, error) {
	var rsp lbsingle.GetLoginUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	rsp.User = uCtx.ExtInfo().(*lbbase.BaseUser)

	return &rsp, err
}

func (a *LbsingleServer) UpdateLoginUserInfo(ctx context.Context, req *lbsingle.UpdateLoginUserInfoReq) (*lbsingle.UpdateLoginUserInfoRsp, error) {
	var rsp lbsingle.UpdateLoginUserInfoRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	baseUser := uCtx.ExtInfo().(*lbbase.BaseUser)

	_, err = OrmUser.NewBaseScope().Where(lbsingle.FieldId_, baseUser.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmUser.NewBaseScope().Where(lbsingle.FieldId_, baseUser.Id).Update(uCtx, utils.OrmStruct2Map(req.User, lbsingle.FieldRole))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = cache.SetLoginInfo(uCtx.Sid(), req.User)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) ResetPassword(ctx context.Context, req *lbsingle.ResetPasswordReq) (*lbsingle.ResetPasswordRsp, error) {
	var rsp lbsingle.ResetPasswordRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	baseUser := uCtx.ExtInfo().(*lbbase.BaseUser)

	oldPsw := utils.Md5(req.OldPassword)
	newPsw := utils.Md5(req.NewPassword)

	_, err = OrmUser.NewBaseScope().Where(map[string]interface{}{
		lbsingle.FieldId_:       baseUser.Id,
		lbsingle.FieldPassword_: oldPsw,
	}).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, lbsingle.ErrOldPasswordNotEqual
	}

	if oldPsw == newPsw {
		return nil, lbsingle.ErrOldPwdEqualNewPwd
	}

	_, err = OrmUser.NewBaseScope().Where(map[string]interface{}{
		lbsingle.FieldId_: baseUser.Id,
	}).Update(uCtx, map[string]interface{}{
		lbsingle.FieldPassword_: newPsw,
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = cache.DelLoginInfo(uCtx.Sid())
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) SyncFile(ctx context.Context, _ *lbsingle.SyncFileReq) (*lbsingle.SyncFileRsp, error) {
	var rsp lbsingle.SyncFileRsp
	var err error

	// 单体环境 直接用单机的 singleflight 防并发
	_, err, _ = a.singleGroup.Do("syncfile", func() (interface{}, error) {
		err := SyncFileIndex(ctx, true)
		if err != nil {
			return nil, lberr.Wrap(err)
		}
		return nil, nil
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) AddUser(ctx context.Context, req *lbsingle.AddUserReq) (*lbsingle.AddUserRsp, error) {
	var rsp lbsingle.AddUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	// check user
	user := req.Data
	user.Password = strings.TrimSpace(user.Password)
	if len(user.Password) < 6 || len(user.Password) > 16 {
		return nil, lbsingle.ErrPasswordLength
	}
	user.Password = utils.Md5(user.Password)

	err = OrmUser.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data.ToBaseUser()

	return &rsp, err
}

func (a *LbsingleServer) DelUserList(ctx context.Context, req *lbsingle.DelUserListReq) (*lbsingle.DelUserListRsp, error) {
	var rsp lbsingle.DelUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetUserList(ctx, &lbsingle.GetUserListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, lbsingle.FieldId_),
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbsingle.FieldId)
	_, err = OrmUser.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) UpdateUser(ctx context.Context, req *lbsingle.UpdateUserReq) (*lbsingle.UpdateUserRsp, error) {
	var rsp lbsingle.UpdateUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmUser.NewBaseScope().Where(lbsingle.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmUser.NewBaseScope().Where(lbsingle.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map4Update(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) GetUser(ctx context.Context, req *lbsingle.GetUserReq) (*lbsingle.GetUserRsp, error) {
	var rsp lbsingle.GetUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmUser.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data.ToBaseUser()

	return &rsp, err
}

func (a *LbsingleServer) GetUserList(ctx context.Context, req *lbsingle.GetUserListReq) (*lbsingle.GetUserListRsp, error) {
	var rsp lbsingle.GetUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmUser.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbsingle.ModelUser](req.ListOption, db)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = core.NewOptionsProcessor(req.ListOption).
		AddString(lbsingle.GetUserListReq_ListOptionLikeNickname, func(val string) error {
			db.WhereLike(lbsingle.FieldNickname_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		AddString(lbsingle.GetUserListReq_ListOptionLikeUsername, func(val string) error {
			db.WhereLike(lbsingle.FieldUsername_, fmt.Sprintf("%%%s%%", val))
			return nil
		}).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}

func (a *LbsingleServer) CheckAuthSys(_ context.Context, req *lbsingle.CheckAuthSysReq) (*lbsingle.CheckAuthSysRsp, error) {
	var rsp lbsingle.CheckAuthSysRsp
	var err error

	uCtx := uctx.NewBaseUCtx()
	uCtx.SetSid(req.Sid)
	auth, err := CheckAuth(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.User = auth
	return &rsp, err
}
func (a *LbsingleServer) AddDailyShortSentences(ctx context.Context, req *lbsingle.AddDailyShortSentencesReq) (*lbsingle.AddDailyShortSentencesRsp, error) {
	var rsp lbsingle.AddDailyShortSentencesRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = OrmDailyShortSentences.NewBaseScope().Create(uCtx, req.Data)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = req.Data

	return &rsp, err
}
func (a *LbsingleServer) DelDailyShortSentencesList(ctx context.Context, req *lbsingle.DelDailyShortSentencesListReq) (*lbsingle.DelDailyShortSentencesListRsp, error) {
	var rsp lbsingle.DelDailyShortSentencesListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetDailyShortSentencesList(ctx, &lbsingle.GetDailyShortSentencesListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, lbsingle.FieldId_),
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbsingle.FieldId)
	_, err = OrmDailyShortSentences.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) UpdateDailyShortSentences(ctx context.Context, req *lbsingle.UpdateDailyShortSentencesReq) (*lbsingle.UpdateDailyShortSentencesRsp, error) {
	var rsp lbsingle.UpdateDailyShortSentencesRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmDailyShortSentences.NewBaseScope().Where(lbsingle.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	_, err = OrmDailyShortSentences.NewBaseScope().Where(lbsingle.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) GetDailyShortSentences(ctx context.Context, req *lbsingle.GetDailyShortSentencesReq) (*lbsingle.GetDailyShortSentencesRsp, error) {
	var rsp lbsingle.GetDailyShortSentencesRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmDailyShortSentences.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbsingleServer) GetDailyShortSentencesList(ctx context.Context, req *lbsingle.GetDailyShortSentencesListReq) (*lbsingle.GetDailyShortSentencesListRsp, error) {
	var rsp lbsingle.GetDailyShortSentencesListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmDailyShortSentences.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbsingle.ModelDailyShortSentences](req.ListOption, db)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
