package logic

import (
	"strconv"

	"postapocgame/admin-server/pkg/errs"
)

// validateLength 复制自 internal/dict.ValidateLength：按 rune 计数校验长度（支持中文等多
// 字节字符），content-rpc 拆分后不再读字典决定上限，maxLength 来自
// services/content/etc/content.yaml 的静态 Limits 配置，见 18-service-extraction-runbook.md
// 2.4 节。和 16-rpc-conventions.md 第 6 节"直接复制不共享"策略一致。
func validateLength(text string, maxLength int64, fieldName string) error {
	runeCount := 0
	for range text {
		runeCount++
	}
	if int64(runeCount) > maxLength {
		return errs.New(errs.CodeBadRequest, fieldName+"长度不能超过"+strconv.FormatInt(maxLength, 10)+"个字符")
	}
	return nil
}
