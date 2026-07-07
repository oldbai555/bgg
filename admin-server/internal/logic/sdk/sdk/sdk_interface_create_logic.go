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
	"postapocgame/admin-server/internal/model/sdk"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
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
	sdkRepo := sdkrepo.NewSdkRepository(l.svcCtx.Repository)
	apiCode := sdkRepo.BuildInterfaceCode(req.Method, req.Path)

	repo := sdkrepo.NewSdkAdminRepository(l.svcCtx.Repository)
	if _, err := repo.FindInterfaceByCode(l.ctx, apiCode); err == nil {
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

	if _, err := repo.CreateInterface(l.ctx, data); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建接口失败", err)
	}

	return nil
}
