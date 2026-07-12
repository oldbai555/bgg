package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkInterfaceUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkInterfaceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceUpdateLogic {
	return &SdkInterfaceUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkInterfaceUpdateLogic) SdkInterfaceUpdate(in *sdk.SdkInterfaceUpdateRequest) (*sdk.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "ID 不能为空"))
	}

	iface, err := l.svcCtx.Admin.FindInterface(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "接口不存在", err))
	}

	if strings.TrimSpace(in.Name) != "" {
		iface.Name = in.Name
	}
	// Path 或 Method 变更时，自动重新生成 API Code
	pathChanged := strings.TrimSpace(in.Path) != "" && in.Path != iface.Path
	methodChanged := strings.TrimSpace(in.Method) != "" && strings.ToUpper(in.Method) != iface.Method
	if pathChanged || methodChanged {
		// 根据新的 path 和 method 生成 API Code
		newPath := iface.Path
		newMethod := iface.Method
		if pathChanged {
			newPath = in.Path
		}
		if methodChanged {
			newMethod = strings.ToUpper(in.Method)
		}
		newApiCode := l.svcCtx.Public.BuildInterfaceCode(newMethod, newPath)
		// 检查新生成的 API Code 是否与其他记录冲突（排除自己）
		existing, err := l.svcCtx.Admin.FindInterfaceByCode(l.ctx, newApiCode)
		if err == nil && existing != nil && existing.Id != in.Id {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该接口路径和方法组合已存在"))
		}
		iface.Path = newPath
		iface.Method = newMethod
		iface.ApiCode = newApiCode
	}
	if in.RateLimitDefault != 0 {
		iface.RateLimitDefault = in.RateLimitDefault
	}
	if in.Status == 1 || in.Status == 2 {
		iface.Status = in.Status
	}
	if in.Remark != "" {
		iface.Remark = in.Remark
	}

	if err := l.svcCtx.Admin.UpdateInterface(l.ctx, iface); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新接口失败", err))
	}

	return &sdk.Empty{}, nil
}
