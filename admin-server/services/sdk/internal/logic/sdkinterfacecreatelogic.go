package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkInterfaceCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkInterfaceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceCreateLogic {
	return &SdkInterfaceCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkInterfaceCreateLogic) SdkInterfaceCreate(in *sdk.SdkInterfaceCreateRequest) (*sdk.Empty, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Path) == "" || strings.TrimSpace(in.Method) == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "名称/path/method 不能为空"))
	}

	// 根据 path 和 method 自动生成 apiCode（与 VerifyApiKey 校验逻辑保持一致）
	apiCode := l.svcCtx.Public.BuildInterfaceCode(in.Method, in.Path)

	if _, err := l.svcCtx.Admin.FindInterfaceByCode(l.ctx, apiCode); err == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该接口路径和方法组合已存在"))
	}

	data := &sdkmodel.SdkInterface{
		Name:             in.Name,
		ApiCode:          apiCode, // 后端自动生成
		Path:             in.Path,
		Method:           strings.ToUpper(in.Method),
		RateLimitDefault: in.RateLimitDefault,
		Status:           in.Status,
		Remark:           in.Remark,
	}

	if _, err := l.svcCtx.Admin.CreateInterface(l.ctx, data); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建接口失败", err))
	}

	return &sdk.Empty{}, nil
}
