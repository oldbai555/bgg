package lbsingleserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/uctx"
	"os"
)

func (a *LbsingleServer) DelFileList(ctx context.Context, req *lbsingle.DelFileListReq) (*lbsingle.DelFileListRsp, error) {
	var rsp lbsingle.DelFileListRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	listRsp, err := a.GetFileList(ctx, &lbsingle.GetFileListReq{
		ListOption: req.ListOption.
			SetSkipTotal().
			AddOpt(core.DefaultListOption_DefaultListOptionSelect, []string{lbsingle.FieldId_, lbsingle.FieldSortUrl_, lbsingle.FieldPath_}),
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

	idList := utils.PluckUint64List(listRsp.List, lbsingle.FieldId)
	_, err = OrmFile.NewBaseScope().WhereIn(lbsingle.FieldId_, idList).Delete(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, err
}

func (a *LbsingleServer) GetFile(ctx context.Context, req *lbsingle.GetFileReq) (*lbsingle.GetFileRsp, error) {
	var rsp lbsingle.GetFileRsp
	var err error
	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	data, err := OrmFile.NewBaseScope().Where(lbsingle.FieldId_, req.Id).First(uCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	rsp.Data = data

	return &rsp, err
}
