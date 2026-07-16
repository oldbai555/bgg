package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictItemListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictItemListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemListLogic {
	return &DictItemListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictItemListLogic) DictItemList(in *iam.DictItemListRequest) (*iam.DictItemListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.System.DictItem.FindPage(l.ctx, in.Page, in.PageSize, in.TypeId, in.Label)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询字典项列表失败", err))
	}

	items := make([]*iam.DictItemItem, 0, len(list))
	for _, di := range list {
		remark := ""
		if di.Remark.Valid {
			remark = di.Remark.String
		}
		items = append(items, &iam.DictItemItem{
			Id:        di.Id,
			TypeId:    di.TypeId,
			Label:     di.Label,
			Value:     di.Value,
			Sort:      di.Sort,
			Status:    di.Status,
			Remark:    remark,
			CreatedAt: di.CreatedAt,
		})
	}

	return &iam.DictItemListResponse{
		Total: total,
		List:  items,
	}, nil
}
