// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

import (
	"context"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"postapocgame/admin-server/internal/model/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkInterfaceCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkInterfaceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceCreateLogic {
	return &SdkInterfaceCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkInterfaceCreateLogic) SdkInterfaceCreate(req *types.SdkInterfaceCreateReq) error {
	if req == nil {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Path) == "" || strings.TrimSpace(req.Method) == "" {
		return errs.New(errs.CodeBadRequest, "名称/path/method 不能为空")
	}

	// 根据 path 和 method 自动生成 apiCode（与中间件校验逻辑保持一致）
	apiCode := l.svcCtx.Domain.SDK.Public.BuildInterfaceCode(req.Method, req.Path)

	if _, err := l.svcCtx.Domain.SDK.Admin.FindInterfaceByCode(l.ctx, apiCode); err == nil {
		return errs.New(errs.CodeBadRequest, "该接口路径和方法组合已存在")
	}

	data := &sdk.SdkInterface{
		Name:             req.Name,
		ApiCode:          apiCode, // 后端自动生成
		Path:             req.Path,
		Method:           strings.ToUpper(req.Method),
		RateLimitDefault: req.RateLimitDefault,
		Status:           req.Status,
		Remark:           req.Remark,
	}

	if _, err := l.svcCtx.Domain.SDK.Admin.CreateInterface(l.ctx, data); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建接口失败", err)
	}

	return nil
}
