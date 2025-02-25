package lbossserver

import (
	"context"
	"github.com/oldbai555/bgg/service/lboss"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/micro/core"
	"github.com/oldbai555/micro/gormx"
	"github.com/oldbai555/micro/uctx"
)

var OnceSvrImpl = &LbossServer{}

type LbossServer struct {
	lboss.UnimplementedLbossServer
}

func (a *LbossServer) GetFileList(ctx context.Context, req *lboss.GetFileListReq) (*lboss.GetFileListRsp, error) {
	var rsp lboss.GetFileListRsp
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	db := OrmFile.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*lboss.ModelFile](req.ListOption, db)
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
