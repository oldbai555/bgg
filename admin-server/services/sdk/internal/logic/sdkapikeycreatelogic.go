package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
	sdkrepo "postapocgame/admin-server/services/sdk/internal/repository/sdk"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkApiKeyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyCreateLogic {
	return &SdkApiKeyCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkApiKeyCreateLogic) SdkApiKeyCreate(in *sdk.SdkApiKeyCreateRequest) (*sdk.SdkApiKeyCreateResponse, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "名称不能为空"))
	}

	apiKey, apiSecret, err := l.generateUniqueKeyPair(l.svcCtx.Admin)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "生成 API Key 失败", err))
	}

	status := in.Status
	if status != 1 && status != 2 {
		status = 1 // 默认启用
	}

	data := &sdkmodel.SdkKey{
		Name:        in.Name,
		ApiKey:      apiKey,
		ApiSecret:   apiSecret,
		Status:      status,
		ExpireAt:    in.ExpireAt,
		IpWhitelist: in.IpWhitelist,
		Remark:      in.Remark,
	}

	id, err := l.svcCtx.Admin.CreateSdkKey(l.ctx, data)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建 API Key 失败", err))
	}

	return &sdk.SdkApiKeyCreateResponse{
		Id:        id,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	}, nil
}

func (l *SdkApiKeyCreateLogic) generateUniqueKeyPair(repo *sdkrepo.SdkAdminRepository) (string, string, error) {
	for i := 0; i < 5; i++ {
		key := randomHex(24)
		secret := randomHex(32)
		_, err := repo.FindSdkKeyByApiKey(l.ctx, key)
		if err == nil {
			continue
		}
		// allow not found
		if _, err := repo.FindSdkKeyByApiSecret(l.ctx, secret); err == nil {
			continue
		}
		return key, secret, nil
	}
	return "", "", errs.New(errs.CodeInternalError, "生成唯一 API Key 失败")
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return strings.ToLower(hex.EncodeToString(b))
}
