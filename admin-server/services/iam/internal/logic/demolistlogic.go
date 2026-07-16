package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDemoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoListLogic {
	return &DemoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DemoListLogic) DemoList(in *iam.DemoListRequest) (*iam.DemoListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.Misc.Demo.FindPage(l.ctx, in.Page, in.PageSize, in.Name)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询演示功能列表失败", err))
	}

	items := make([]*iam.DemoItem, 0, len(list))
	for _, d := range list {
		items = append(items, &iam.DemoItem{
			Id:        d.Id,
			Name:      d.Name,
			Status:    d.Status,
			CreatedAt: d.CreatedAt,
		})
	}

	return &iam.DemoListResponse{Total: total, List: items}, nil
}
