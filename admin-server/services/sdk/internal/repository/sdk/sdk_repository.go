package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
	"postapocgame/admin-server/services/sdk/internal/repository"
)

// SdkRepository 封装 SDK 相关数据访问。
type SdkRepository struct {
	store *repository.Store
}

func NewSdkRepository(store *repository.Store) *SdkRepository {
	return &SdkRepository{store: store}
}

func (r *SdkRepository) FindKeyByApiKey(ctx context.Context, apiKey string) (*sdkmodel.SdkKey, error) {
	return r.store.SdkKeyModel.FindOneByApiKey(ctx, apiKey)
}

func (r *SdkRepository) FindInterfaceByCode(ctx context.Context, apiCode string) (*sdkmodel.SdkInterface, error) {
	return r.store.SdkInterfaceModel.FindOneByApiCode(ctx, apiCode)
}

func (r *SdkRepository) FindInterfaceByPathMethod(ctx context.Context, path, method string) (*sdkmodel.SdkInterface, error) {
	// 为避免新增索引，这里简单通过接口编码唯一键做映射，约定 api_code = method:lower + path
	code := r.BuildInterfaceCode(method, path)
	return r.FindInterfaceByCode(ctx, code)
}

func (r *SdkRepository) BuildInterfaceCode(method, path string) string {
	return fmt.Sprintf("%s:%s", strings.ToLower(method), path)
}

func (r *SdkRepository) FindKeyApiBinding(ctx context.Context, sdkKeyId, sdkInterfaceId uint64) (*sdkmodel.SdkKeyApi, error) {
	return r.store.SdkKeyApiModel.FindOneBySdkKeyIdSdkInterfaceId(ctx, sdkKeyId, sdkInterfaceId)
}

// GetDefaultRateLimit 返回接口自身默认值；若接口未设置（<=0），退回 staticDefault。
// staticDefault 原来读字典 sdk_rate_limit_default（物理属于 iam 域），拆分后 sdk-rpc
// 拿不到那张表，改成 services/sdk/etc/sdk.yaml 里的静态配置 RateLimitDefault，和
// task-rpc 的 task_recent_logic.go 处理 system 域字典依赖是同一个模式（18-service-
// extraction-runbook.md 2.1 节）。字典种子数据里的默认值是 60，静态配置沿用同一个值。
func (r *SdkRepository) GetDefaultRateLimit(interfaceDefault, staticDefault int64) int64 {
	if interfaceDefault > 0 {
		return interfaceDefault
	}
	if staticDefault > 0 {
		return staticDefault
	}
	return 60
}

// IsIPAllowed 校验 IP 白名单（空表示不限制）
func (r *SdkRepository) IsIPAllowed(ipWhitelist string, clientIP string) bool {
	if strings.TrimSpace(ipWhitelist) == "" {
		return true
	}
	parts := strings.FieldsFunc(ipWhitelist, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == ';'
	})
	clientIP = strings.TrimSpace(clientIP)
	for _, p := range parts {
		if strings.TrimSpace(p) == clientIP {
			return true
		}
	}
	return false
}

// SaveCallLog 保存调用日志（已截断/脱敏的 req/resp 由调用方处理）
func (r *SdkRepository) SaveCallLog(ctx context.Context, log *sdkmodel.SdkCallLog) error {
	_, err := r.store.SdkCallLogModel.Insert(ctx, log)
	return errors.WithStack(err)
}
