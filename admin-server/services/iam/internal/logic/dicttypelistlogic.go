package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeListLogic {
	return &DictTypeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictTypeListLogic) DictTypeList(in *iam.DictTypeListRequest) (*iam.DictTypeListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.System.DictType.FindPage(l.ctx, in.Page, in.PageSize, in.Name, in.Code)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询字典类型列表失败", err))
	}

	items := make([]*iam.DictTypeItem, 0, len(list))
	for _, dt := range list {
		description := ""
		if dt.Description.Valid {
			description = dt.Description.String
		}
		items = append(items, &iam.DictTypeItem{
			Id:          dt.Id,
			Name:        dt.Name,
			Code:        dt.Code,
			Description: description,
			Status:      dt.Status,
			CreatedAt:   dt.CreatedAt,
		})
	}

	return &iam.DictTypeListResponse{
		Total: total,
		List:  items,
	}, nil
}
