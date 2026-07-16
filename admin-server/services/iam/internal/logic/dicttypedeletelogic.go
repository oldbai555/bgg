package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictTypeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeDeleteLogic {
	return &DictTypeDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictTypeDeleteLogic) DictTypeDelete(in *iam.DictTypeDeleteRequest) (*iam.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "字典类型ID不能为空"))
	}

	items, err := l.svcCtx.Domain.System.DictItem.FindByTypeID(l.ctx, in.Id)
	if err == nil && len(items) > 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该字典类型下存在字典项，无法删除"))
	}

	if err := l.svcCtx.Domain.System.DictType.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除字典类型失败", err))
	}
	return &iam.Empty{}, nil
}
