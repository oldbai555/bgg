package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"postapocgame/admin-server/internal/model"

	"github.com/pkg/errors"
)

// SdkRepository 封装 SDK 相关数据访问。
type SdkRepository struct {
	repo *Repository
}

func NewSdkRepository(repo *Repository) *SdkRepository {
	return &SdkRepository{repo: repo}
}

func (r *SdkRepository) FindKeyByApiKey(ctx context.Context, apiKey string) (*model.SdkKey, error) {
	return r.repo.SdkKeyModel.FindOneByApiKey(ctx, apiKey)
}

func (r *SdkRepository) FindInterfaceByCode(ctx context.Context, apiCode string) (*model.SdkInterface, error) {
	return r.repo.SdkInterfaceModel.FindOneByApiCode(ctx, apiCode)
}

func (r *SdkRepository) FindInterfaceByPathMethod(ctx context.Context, path, method string) (*model.SdkInterface, error) {
	// 为避免新增索引，这里简单通过接口编码唯一键做映射，约定 api_code = method:lower + path
	code := r.BuildInterfaceCode(method, path)
	return r.FindInterfaceByCode(ctx, code)
}

func (r *SdkRepository) BuildInterfaceCode(method, path string) string {
	return fmt.Sprintf("%s:%s", strings.ToLower(method), path)
}

func (r *SdkRepository) FindKeyApiBinding(ctx context.Context, sdkKeyId, sdkInterfaceId uint64) (*model.SdkKeyApi, error) {
	return r.repo.SdkKeyApiModel.FindOneBySdkKeyIdSdkInterfaceId(ctx, sdkKeyId, sdkInterfaceId)
}

// GetDefaultRateLimit 返回接口默认值或字典默认值（sdk_rate_limit_default）；若均无则 60。
func (r *SdkRepository) GetDefaultRateLimit(ctx context.Context, interfaceDefault int64) int64 {
	if interfaceDefault > 0 {
		return interfaceDefault
	}
	typeId, err := r.repo.AdminDictTypeModel.FindIdByCode(ctx, "sdk_rate_limit_default")
	if err == nil && typeId > 0 {
		items, _, _ := r.repo.AdminDictItemModel.FindPageByTypeId(ctx, typeId, 1, 1)
		if len(items) > 0 {
			if v := parseInt64(items[0].Value); v > 0 {
				return v
			}
		}
	}
	return 60
}

func parseInt64(s string) int64 {
	if s == "" {
		return 0
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	return 0
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
func (r *SdkRepository) SaveCallLog(ctx context.Context, log *model.SdkCallLog) error {
	_, err := r.repo.SdkCallLogModel.Insert(ctx, log)
	return errors.WithStack(err)
}
