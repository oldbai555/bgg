package dict

import (
	"context"
	"strconv"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetIntValue 从字典中读取整数值配置
// code: 字典类型编码
// defaultValue: 如果字典不存在或解析失败时的默认值
// 返回: 字典项的第一个 value 解析为整数，如果失败则返回 defaultValue
func GetIntValue(ctx context.Context, repo *repository.Repository, code string, defaultValue int) int {
	dictTypeRepo := repository.NewDictTypeRepository(repo)
	dictType, err := dictTypeRepo.FindByCode(ctx, code)
	if err != nil {
		logx.WithContext(ctx).Errorf("查询字典类型失败: code=%s, error=%v, 使用默认值=%d", code, err, defaultValue)
		return defaultValue
	}

	dictItemRepo := repository.NewDictItemRepository(repo)
	items, err := dictItemRepo.FindByTypeID(ctx, dictType.Id)
	if err != nil {
		logx.WithContext(ctx).Errorf("查询字典项失败: code=%s, error=%v, 使用默认值=%d", code, err, defaultValue)
		return defaultValue
	}

	if len(items) == 0 {
		logx.WithContext(ctx).Errorf("字典项为空: code=%s, 使用默认值=%d", code, defaultValue)
		return defaultValue
	}

	// 取第一个字典项的 value 解析为整数
	valueStr := items[0].Value
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logx.WithContext(ctx).Errorf("字典值解析失败: code=%s, value=%s, error=%v, 使用默认值=%d", code, valueStr, err, defaultValue)
		return defaultValue
	}

	return value
}

// GetStringValue 从字典中读取字符串值配置
// code: 字典类型编码
// defaultValue: 如果字典不存在时的默认值
// 返回: 字典项的第一个 value，如果失败则返回 defaultValue
func GetStringValue(ctx context.Context, repo *repository.Repository, code string, defaultValue string) string {
	dictTypeRepo := repository.NewDictTypeRepository(repo)
	dictType, err := dictTypeRepo.FindByCode(ctx, code)
	if err != nil {
		logx.WithContext(ctx).Errorf("查询字典类型失败: code=%s, error=%v, 使用默认值=%s", code, err, defaultValue)
		return defaultValue
	}

	dictItemRepo := repository.NewDictItemRepository(repo)
	items, err := dictItemRepo.FindByTypeID(ctx, dictType.Id)
	if err != nil {
		logx.WithContext(ctx).Errorf("查询字典项失败: code=%s, error=%v, 使用默认值=%s", code, err, defaultValue)
		return defaultValue
	}

	if len(items) == 0 {
		logx.WithContext(ctx).Errorf("字典项为空: code=%s, 使用默认值=%s", code, defaultValue)
		return defaultValue
	}

	return items[0].Value
}

// ValidateLength 校验字符串长度（按 rune 计数，支持中文）
// text: 待校验文本
// maxLength: 最大长度
// fieldName: 字段名称（用于错误提示）
func ValidateLength(text string, maxLength int, fieldName string) error {
	if maxLength <= 0 {
		return errs.New(errs.CodeBadRequest, fieldName+"长度限制配置无效")
	}

	// 使用 rune 计数，支持中文等多字节字符
	runeCount := 0
	for range text {
		runeCount++
	}
	if runeCount > maxLength {
		return errs.New(errs.CodeBadRequest, fieldName+"长度不能超过"+strconv.Itoa(maxLength)+"个字符")
	}

	return nil
}
