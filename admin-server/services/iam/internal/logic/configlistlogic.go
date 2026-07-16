package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigListLogic {
	return &ConfigListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigListLogic) ConfigList(in *iam.ConfigListRequest) (*iam.ConfigListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	list, total, err := l.svcCtx.Domain.System.Config.FindPage(l.ctx, in.Page, in.PageSize, in.Group, in.Key)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询配置列表失败", err))
	}

	items := make([]*iam.ConfigItem, 0, len(list))
	for _, c := range list {
		value := ""
		if c.Value.Valid {
			value = c.Value.String
		}
		description := ""
		if c.Description.Valid {
			description = c.Description.String
		}
		items = append(items, &iam.ConfigItem{
			Id:          c.Id,
			Group:       c.Group,
			Key:         c.Key,
			Value:       value,
			ConfigType:  c.Type,
			Description: description,
			CreatedAt:   c.CreatedAt,
		})
	}

	return &iam.ConfigListResponse{
		Total: total,
		List:  items,
	}, nil
}
