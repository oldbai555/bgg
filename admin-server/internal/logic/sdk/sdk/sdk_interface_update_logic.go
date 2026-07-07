// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

import (
	"context"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
)

type SdkInterfaceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkInterfaceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceUpdateLogic {
	return &SdkInterfaceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkInterfaceUpdateLogic) SdkInterfaceUpdate(req *types.SdkInterfaceUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "ID 不能为空")
	}

	repo := sdkrepo.NewSdkAdminRepository(l.svcCtx.Repository)
	iface, err := repo.FindInterface(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeBadRequest, "接口不存在", err)
	}

	if strings.TrimSpace(req.Name) != "" {
		iface.Name = req.Name
	}
	// Path 或 Method 变更时，自动重新生成 API Code
	pathChanged := strings.TrimSpace(req.Path) != "" && req.Path != iface.Path
	methodChanged := strings.TrimSpace(req.Method) != "" && strings.ToUpper(req.Method) != iface.Method
	if pathChanged || methodChanged {
		// 根据新的 path 和 method 生成 API Code
		sdkRepo := sdkrepo.NewSdkRepository(l.svcCtx.Repository)
		newPath := iface.Path
		newMethod := iface.Method
		if pathChanged {
			newPath = req.Path
		}
		if methodChanged {
			newMethod = strings.ToUpper(req.Method)
		}
		newApiCode := sdkRepo.BuildInterfaceCode(newMethod, newPath)
		// 检查新生成的 API Code 是否与其他记录冲突（排除自己）
		existing, err := repo.FindInterfaceByCode(l.ctx, newApiCode)
		if err == nil && existing != nil && existing.Id != req.Id {
			return errs.New(errs.CodeBadRequest, "该接口路径和方法组合已存在")
		}
		iface.Path = newPath
		iface.Method = newMethod
		iface.ApiCode = newApiCode
	}
	if req.RateLimitDefault != 0 {
		iface.RateLimitDefault = req.RateLimitDefault
	}
	if req.Status == 1 || req.Status == 2 {
		iface.Status = req.Status
	}
	if req.Remark != "" {
		iface.Remark = req.Remark
	}

	if err := repo.UpdateInterface(l.ctx, iface); err != nil {
		return errs.Wrap(errs.CodeInternalError, "更新接口失败", err)
	}

	return nil
}
