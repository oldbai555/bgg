// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/model/sdk"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
)

type SdkApiKeyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkApiKeyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyCreateLogic {
	return &SdkApiKeyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkApiKeyCreateLogic) SdkApiKeyCreate(req *types.SdkApiKeyCreateReq) (resp *types.SdkApiKeyCreateResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errs.New(errs.CodeBadRequest, "名称不能为空")
	}

	repo := sdkrepo.NewSdkAdminRepository(l.svcCtx.Repository)

	apiKey, apiSecret, err := l.generateUniqueKeyPair(repo)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "生成 API Key 失败", err)
	}

	status := req.Status
	if status != 1 && status != 2 {
		status = 1 // 默认启用
	}

	data := &sdk.SdkKey{
		Name:        req.Name,
		ApiKey:      apiKey,
		ApiSecret:   apiSecret,
		Status:      status,
		ExpireAt:    req.ExpireAt,
		IpWhitelist: req.IpWhitelist,
		Remark:      req.Remark,
	}

	id, err := repo.CreateSdkKey(l.ctx, data)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建 API Key 失败", err)
	}

	return &types.SdkApiKeyCreateResp{
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
