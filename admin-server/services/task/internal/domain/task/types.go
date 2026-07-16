//staticcheck:ignore SA5008 // "optional" is a go-zero framework extension for JSON tags
package task

// TaskParamsReq 通用任务参数结构
type TaskParamsReq struct {
	// 基础字段（所有任务类型共有）
	Module string `json:"module"` // 模块名称（如：operation_log、audit_log、login_log等）
}

// TaskResultResp 通用任务结果结构
type TaskResultResp struct {
	// 基础字段（所有任务类型共有）
	Success bool   `json:"success"`          // 是否成功
	Message string `json:"message,optional"` // 结果消息（可选）
}

// ExcelExportParams Excel导出任务参数（增量扩展）
type ExcelExportParams struct {
	TaskParamsReq
	Filters map[string]interface{} `json:"filters"`         // 筛选条件
	Fields  []string               `json:"fields,optional"` // 导出字段（可选）
}

// ExcelExportResult Excel导出任务结果（增量扩展）
type ExcelExportResult struct {
	TaskResultResp
	FileURL     string `json:"fileUrl"`     // 文件下载URL
	FileName    string `json:"fileName"`    // 文件名
	FileSize    int64  `json:"fileSize"`    // 文件大小（字节）
	RecordCount int64  `json:"recordCount"` // 记录数量
}
