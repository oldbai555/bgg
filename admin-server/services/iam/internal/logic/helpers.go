package logic

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func isErrNotFound(err error) bool {
	return errors.Is(err, sqlx.ErrNotFound)
}

// normalizePage 与 gateway internal/logic/logicutil.NormalizePage 同构；
// iam-rpc 不反向依赖 gateway internal 包，因此在服务内维护一份副本。
func normalizePage(page, pageSize, defaultPageSize, maxPageSize int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if maxPageSize > 0 && pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
