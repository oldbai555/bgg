package lbwxmpserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lbwxmp"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
)

func (a *LbwxmpServer) MPUserInfo(ctx context.Context, req *lbwxmp.MPUserInfoReq) (*lbwxmp.MPUserInfoRsp, error) {
	var rsp lbwxmp.MPUserInfoRsp
	var err error

	return &rsp, err
}
func (a *LbwxmpServer) MPUserMineService(ctx context.Context, req *lbwxmp.MPUserMineServiceReq) (*lbwxmp.MPUserMineServiceRsp, error) {
	var rsp lbwxmp.MPUserMineServiceRsp
	var err error

	return &rsp, err
}
func (a *LbwxmpServer) MPUserSaveInfo(ctx context.Context, req *lbwxmp.MPUserSaveInfoReq) (*lbwxmp.MPUserSaveInfoRsp, error) {
	var rsp lbwxmp.MPUserSaveInfoRsp
	var err error

	return &rsp, err
}
func (a *LbwxmpServer) AddMpMemberUser(ctx context.Context, req *lbwxmp.AddMpMemberUserReq) (*lbwxmp.AddMpMemberUserRsp, error) {
	var rsp lbwxmp.AddMpMemberUserRsp
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
func (a *LbwxmpServer) DelMpMemberUserList(ctx context.Context, req *lbwxmp.DelMpMemberUserListReq) (*lbwxmp.DelMpMemberUserListRsp, error) {
	var rsp lbwxmp.DelMpMemberUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetMpMemberUserList(ctx, &lbwxmp.GetMpMemberUserListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, lbwxmp.FieldId_),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	idList := utils.PluckUint64List(listRsp.List, lbwxmp.FieldId)
	_, err = OrmMpMemberUser.NewBaseScope().WhereIn(lbwxmp.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) UpdateMpMemberUser(ctx context.Context, req *lbwxmp.UpdateMpMemberUserReq) (*lbwxmp.UpdateMpMemberUserRsp, error) {
	var rsp lbwxmp.UpdateMpMemberUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMemberUser.NewBaseScope().Where(lbwxmp.FieldId_, req.Data.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	_, err = OrmMpMemberUser.NewBaseScope().Where(lbwxmp.FieldId_, data.Id).Update(uCtx, utils.OrmStruct2Map(req.Data))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}
func (a *LbwxmpServer) GetMpMemberUser(ctx context.Context, req *lbwxmp.GetMpMemberUserReq) (*lbwxmp.GetMpMemberUserRsp, error) {
	var rsp lbwxmp.GetMpMemberUserRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmMpMemberUser.NewBaseScope().Where(lbwxmp.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}
func (a *LbwxmpServer) GetMpMemberUserList(ctx context.Context, req *lbwxmp.GetMpMemberUserListReq) (*lbwxmp.GetMpMemberUserListRsp, error) {
	var rsp lbwxmp.GetMpMemberUserListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	db := OrmMpMemberUser.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lbwxmp.ModelMpMemberUser](req.ListOption, db)
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
