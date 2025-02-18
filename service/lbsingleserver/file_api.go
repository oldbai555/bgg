package lbsingleserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/uctx"
)

func (a *LbsingleServer) DelFileList(ctx context.Context, req *lbsingle.DelFileListReq) (*lbsingle.DelFileListRsp, error) {
	var rsp lbsingle.DelFileListRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	listRsp, err := a.GetFileList(ctx, &lbsingle.GetFileListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, []string{lbsingle.FieldId_, lbsingle.FieldSortUrl_, lbsingle.FieldPath_}),
	})
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	if len(listRsp.List) == 0 {
		log.Infof("list is empty")
		return &rsp, nil
	}

	// todo 清缓存 和 文件

	idList := utils.PluckUint64List(listRsp.List, lbsingle.FieldId)
	_, err = OrmFile.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	return &rsp, err
}
func (a *LbsingleServer) GetFile(ctx context.Context, req *lbsingle.GetFileReq) (*lbsingle.GetFileRsp, error) {
	var rsp lbsingle.GetFileRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	data, err := OrmFile.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	rsp.Data = data

	return &rsp, err
}
