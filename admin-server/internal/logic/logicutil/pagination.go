package logicutil

// NormalizePage 统一分页参数的默认值与上限控制。
func NormalizePage(page, pageSize, defaultPageSize, maxPageSize int64) (int64, int64) {
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
