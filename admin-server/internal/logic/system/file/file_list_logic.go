// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileListLogic) FileList(req *types.FileListReq) (resp *types.FileListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.FileList(l.ctx, &iamclient.FileListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询文件列表失败", err)
	}

	items := make([]types.FileItem, 0, len(rpcResp.List))
	for _, f := range rpcResp.List {
		items = append(items, types.FileItem{
			Id:           f.Id,
			Name:         f.Name,
			OriginalName: f.OriginalName,
			Path:         f.Path,
			BaseUrl:      f.BaseUrl,
			Status:       f.Status,
			CreatedAt:    f.CreatedAt,
		})
	}

	return &types.FileListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
