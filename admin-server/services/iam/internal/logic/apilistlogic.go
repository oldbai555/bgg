package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiListLogic {
	return &ApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiListLogic) ApiList(in *iam.ApiListRequest) (*iam.ApiListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.IAM.Api.FindPage(l.ctx, in.Page, in.PageSize, in.Name)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询接口列表失败", err))
	}

	items := make([]*iam.ApiItem, 0, len(list))
	for _, a := range list {
		description := ""
		if a.Description.Valid {
			description = a.Description.String
		}
		items = append(items, &iam.ApiItem{
			Id:          a.Id,
			Name:        a.Name,
			Method:      a.Method,
			Path:        a.Path,
			Description: description,
			Status:      a.Status,
			CreatedAt:   int64(a.CreatedAt),
		})
	}

	return &iam.ApiListResponse{
		Total: total,
		List:  items,
	}, nil
}
