package executors

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"postapocgame/admin-server/internal/task"
	"strings"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/interfaces"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/repository"

	"github.com/zeromicro/go-zero/core/logx"
)

// ExcelExportExecutor Excel导出任务执行器
type ExcelExportExecutor struct {
	svcCtx *repository.Repository
}

// NewExcelExportExecutor 创建Excel导出执行器
func NewExcelExportExecutor(svcCtx *repository.Repository) interfaces.TaskExecutor {
	return &ExcelExportExecutor{svcCtx: svcCtx}
}

// GetType 获取任务类型
func (e *ExcelExportExecutor) GetType() int {
	return 1 // 对应字典 task_type value=1（异步导出Excel）
}

// Execute 执行Excel导出任务
func (e *ExcelExportExecutor) Execute(ctx context.Context, taskModel *model.AdminTask, paramsJSON string) (string, error) {
	// 1. 解析参数JSON为ExcelExportParams
	var params task.ExcelExportParams
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		return "", fmt.Errorf("解析任务参数失败: %w", err)
	}

	// 2. 根据module类型执行不同的导出逻辑
	var fileURL, fileName string
	var fileSize, recordCount int64
	var err error

	switch params.Module {
	case consts.TaskModuleOperationLog:
		fileURL, fileName, fileSize, recordCount, err = e.exportOperationLog(ctx, params)
	case consts.TaskModuleAuditLog:
		fileURL, fileName, fileSize, recordCount, err = e.exportAuditLog(ctx, params)
	case consts.TaskModuleLoginLog:
		fileURL, fileName, fileSize, recordCount, err = e.exportLoginLog(ctx, params)
	case consts.TaskModulePerformanceLog:
		fileURL, fileName, fileSize, recordCount, err = e.exportPerformanceLog(ctx, params)
	default:
		return "", fmt.Errorf("不支持的导出模块: %s", params.Module)
	}

	if err != nil {
		return "", fmt.Errorf("导出失败: %w", err)
	}

	// 3. 构建结果JSON
	result := task.ExcelExportResult{
		TaskResultResp: task.TaskResultResp{
			Success: true,
			Message: "导出成功",
		},
		FileURL:     fileURL,
		FileName:    fileName,
		FileSize:    fileSize,
		RecordCount: recordCount,
	}

	// 4. 序列化为JSON字符串返回
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %w", err)
	}

	return string(resultJSON), nil
}

// exportOperationLog 导出操作日志
func (e *ExcelExportExecutor) exportOperationLog(ctx context.Context, params task.ExcelExportParams) (string, string, int64, int64, error) {
	// 解析筛选条件
	userId := uint64(0)
	username := ""
	operationType := ""
	operationObject := ""
	method := ""
	startTime := ""
	endTime := ""

	if params.Filters != nil {
		if v, ok := params.Filters[consts.TaskFilterUserId].(float64); ok {
			userId = uint64(v)
		}
		if v, ok := params.Filters[consts.TaskFilterUsername].(string); ok {
			username = v
		}
		if v, ok := params.Filters[consts.TaskFilterOperationType].(string); ok {
			operationType = v
		}
		if v, ok := params.Filters[consts.TaskFilterOperationObj].(string); ok {
			operationObject = v
		}
		if v, ok := params.Filters[consts.TaskFilterMethod].(string); ok {
			method = v
		}
		if v, ok := params.Filters[consts.TaskFilterStartTime].(string); ok {
			startTime = v
		}
		if v, ok := params.Filters[consts.TaskFilterEndTime].(string); ok {
			endTime = v
		}
	}

	// 查询数据
	operationLogRepo := repository.NewOperationLogRepository(e.svcCtx)
	list, _, err := operationLogRepo.FindPage(ctx, 1, 10000, userId, username, operationType, operationObject, method, startTime, endTime)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("查询操作日志失败: %w", err)
	}

	// 生成CSV文件
	return e.generateCSVFile(ctx, "操作日志", list, func(log interface{}) []string {
		opLog := log.(*model.AdminOperationLog)
		requestParams := ""
		if opLog.RequestParams.Valid {
			requestParams = opLog.RequestParams.String
		}
		return []string{
			fmt.Sprintf("%d", opLog.Id),
			fmt.Sprintf("%d", opLog.UserId),
			opLog.Username,
			opLog.OperationType,
			opLog.OperationObject,
			opLog.Method,
			opLog.Path,
			requestParams,
			fmt.Sprintf("%d", opLog.ResponseCode),
			opLog.ResponseMsg,
			opLog.IpAddress,
			opLog.UserAgent,
			fmt.Sprintf("%d", opLog.Duration),
			time.Unix(opLog.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		}
	}, []string{"ID", "用户ID", "用户名", "操作类型", "操作对象", "请求方法", "请求路径", "请求参数", "响应状态码", "响应消息", "IP地址", "用户代理", "耗时(ms)", "创建时间"})
}

// exportAuditLog 导出审计日志
func (e *ExcelExportExecutor) exportAuditLog(ctx context.Context, params task.ExcelExportParams) (string, string, int64, int64, error) {
	// 解析筛选条件
	userId := uint64(0)
	username := ""
	auditType := ""
	auditObject := ""
	startTime := ""
	endTime := ""

	if params.Filters != nil {
		if v, ok := params.Filters[consts.TaskFilterUserId].(float64); ok {
			userId = uint64(v)
		}
		if v, ok := params.Filters[consts.TaskFilterUsername].(string); ok {
			username = v
		}
		if v, ok := params.Filters[consts.TaskFilterAuditType].(string); ok {
			auditType = v
		}
		if v, ok := params.Filters[consts.TaskFilterAuditObject].(string); ok {
			auditObject = v
		}
		if v, ok := params.Filters[consts.TaskFilterStartTime].(string); ok {
			startTime = v
		}
		if v, ok := params.Filters[consts.TaskFilterEndTime].(string); ok {
			endTime = v
		}
	}

	// 查询数据
	auditLogRepo := repository.NewAuditLogRepository(e.svcCtx)
	list, _, err := auditLogRepo.FindPage(ctx, 1, 10000, userId, username, auditType, auditObject, startTime, endTime)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("查询审计日志失败: %w", err)
	}

	// 生成CSV文件
	return e.generateCSVFile(ctx, "审计日志", list, func(log interface{}) []string {
		auditLog := log.(*model.AuditLog)
		auditDetail := ""
		if auditLog.AuditDetail.Valid {
			auditDetail = auditLog.AuditDetail.String
		}
		return []string{
			fmt.Sprintf("%d", auditLog.Id),
			fmt.Sprintf("%d", auditLog.UserId),
			auditLog.Username,
			auditLog.AuditType,
			auditLog.AuditObject,
			auditDetail,
			auditLog.IpAddress,
			auditLog.UserAgent,
			time.Unix(auditLog.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		}
	}, []string{"ID", "用户ID", "用户名", "审计类型", "审计对象", "审计详情", "IP地址", "用户代理", "创建时间"})
}

// exportLoginLog 导出登录日志
func (e *ExcelExportExecutor) exportLoginLog(ctx context.Context, params task.ExcelExportParams) (string, string, int64, int64, error) {
	// 解析筛选条件
	userId := uint64(0)
	username := ""
	status := 0
	startTime := ""
	endTime := ""

	if params.Filters != nil {
		if v, ok := params.Filters[consts.TaskFilterUserId].(float64); ok {
			userId = uint64(v)
		}
		if v, ok := params.Filters[consts.TaskFilterUsername].(string); ok {
			username = v
		}
		if v, ok := params.Filters[consts.TaskFilterStatus].(float64); ok {
			status = int(v)
		}
		if v, ok := params.Filters[consts.TaskFilterStartTime].(string); ok {
			startTime = v
		}
		if v, ok := params.Filters[consts.TaskFilterEndTime].(string); ok {
			endTime = v
		}
	}

	// 查询数据
	loginLogRepo := repository.NewLoginLogRepository(e.svcCtx)
	list, _, err := loginLogRepo.FindPage(ctx, 1, 10000, userId, username, status, startTime, endTime)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("查询登录日志失败: %w", err)
	}

	// 生成CSV文件
	return e.generateCSVFile(ctx, "登录日志", list, func(log interface{}) []string {
		loginLog := log.(*model.AdminLoginLog)
		return []string{
			fmt.Sprintf("%d", loginLog.Id),
			fmt.Sprintf("%d", loginLog.UserId),
			loginLog.Username,
			loginLog.IpAddress,
			loginLog.Location,
			loginLog.Browser,
			loginLog.Os,
			loginLog.UserAgent,
			fmt.Sprintf("%d", loginLog.Status),
			loginLog.Message,
			time.Unix(loginLog.LoginAt, 0).Format("2006-01-02 15:04:05"),
			time.Unix(loginLog.LogoutAt, 0).Format("2006-01-02 15:04:05"),
			time.Unix(loginLog.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		}
	}, []string{"ID", "用户ID", "用户名", "IP地址", "登录地点", "浏览器", "操作系统", "用户代理", "登录状态", "登录消息", "登录时间", "登出时间", "创建时间"})
}

// exportPerformanceLog 导出性能监控日志
func (e *ExcelExportExecutor) exportPerformanceLog(ctx context.Context, params task.ExcelExportParams) (string, string, int64, int64, error) {
	method := ""
	path := ""
	isSlow := int64(0)
	statusCode := int64(0)
	startTime := ""
	endTime := ""

	if params.Filters != nil {
		if v, ok := params.Filters[consts.TaskFilterMethod].(string); ok {
			method = v
		}
		if v, ok := params.Filters[consts.TaskFilterPath].(string); ok {
			path = v
		}
		if v, ok := params.Filters[consts.TaskFilterIsSlow].(float64); ok {
			isSlow = int64(v)
		}
		if v, ok := params.Filters[consts.TaskFilterStatusCode].(float64); ok {
			statusCode = int64(v)
		}
		if v, ok := params.Filters[consts.TaskFilterStartTime].(string); ok {
			startTime = v
		}
		if v, ok := params.Filters[consts.TaskFilterEndTime].(string); ok {
			endTime = v
		}
	}

	perfRepo := repository.NewPerformanceLogRepository(e.svcCtx)
	list, _, err := perfRepo.FindPage(ctx, 1, 10000, method, path, isSlow, statusCode, startTime, endTime)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("查询性能日志失败: %w", err)
	}

	return e.generateCSVFile(ctx, "性能监控日志", list, func(log interface{}) []string {
		perf := log.(*model.AdminPerformanceLog)
		return []string{
			fmt.Sprintf("%d", perf.Id),
			fmt.Sprintf("%d", perf.UserId),
			perf.Username,
			perf.Method,
			perf.Path,
			fmt.Sprintf("%d", perf.StatusCode),
			fmt.Sprintf("%d", perf.Duration),
			fmt.Sprintf("%d", perf.IsSlow),
			fmt.Sprintf("%d", perf.SlowThreshold),
			perf.IpAddress,
			perf.UserAgent,
			perf.ErrorMsg,
			time.Unix(perf.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		}
	}, []string{"ID", "用户ID", "用户名", "请求方法", "请求路径", "状态码", "耗时(ms)", "是否慢接口", "慢接口阈值(ms)", "IP地址", "User-Agent", "错误信息", "创建时间"})
}

// generateCSVFile 生成CSV文件（通用方法）
func (e *ExcelExportExecutor) generateCSVFile(ctx context.Context, moduleName string, data interface{}, rowConverter func(interface{}) []string, headers []string) (string, string, int64, int64, error) {
	// 创建临时文件
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.csv", moduleName, timestamp)
	fileSystemPath := filepath.Join(consts.UploadDir, fileName)

	// 确保上传目录存在
	if err := os.MkdirAll(consts.UploadDir, 0755); err != nil {
		return "", "", 0, 0, fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 创建文件
	file, err := os.Create(fileSystemPath)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("创建文件失败: %w", err)
	}

	// 写入BOM，确保Excel正确识别UTF-8
	if _, err := file.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		_ = file.Close()
		return "", "", 0, 0, fmt.Errorf("写入BOM失败: %w", err)
	}

	// 创建CSV writer
	writer := csv.NewWriter(file)

	// 写入表头
	if err := writer.Write(headers); err != nil {
		writer.Flush()
		_ = file.Close()
		return "", "", 0, 0, fmt.Errorf("写入CSV表头失败: %w", err)
	}

	// 写入数据
	recordCount := int64(0)
	switch v := data.(type) {
	case []model.AdminOperationLog:
		for i := range v {
			row := rowConverter(&v[i])
			if err := writer.Write(row); err != nil {
				writer.Flush()
				_ = file.Close()
				return "", "", 0, 0, fmt.Errorf("写入CSV数据失败: %w", err)
			}
			recordCount++
		}
	case []model.AuditLog:
		for i := range v {
			row := rowConverter(&v[i])
			if err := writer.Write(row); err != nil {
				writer.Flush()
				_ = file.Close()
				return "", "", 0, 0, fmt.Errorf("写入CSV数据失败: %w", err)
			}
			recordCount++
		}
	case []model.AdminLoginLog:
		for i := range v {
			row := rowConverter(&v[i])
			if err := writer.Write(row); err != nil {
				writer.Flush()
				_ = file.Close()
				return "", "", 0, 0, fmt.Errorf("写入CSV数据失败: %w", err)
			}
			recordCount++
		}
	case []model.AdminPerformanceLog:
		for i := range v {
			row := rowConverter(&v[i])
			if err := writer.Write(row); err != nil {
				writer.Flush()
				_ = file.Close()
				return "", "", 0, 0, fmt.Errorf("写入CSV数据失败: %w", err)
			}
			recordCount++
		}
	default:
		writer.Flush()
		_ = file.Close()
		return "", "", 0, 0, fmt.Errorf("不支持的数据类型: %T", data)
	}

	// 确保缓冲区写入磁盘
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = file.Close()
		return "", "", 0, 0, fmt.Errorf("刷新CSV缓冲区失败: %w", err)
	}

	// 关闭文件句柄（Windows 下如果不关闭会导致后续重命名失败）
	if err := file.Close(); err != nil {
		return "", "", 0, 0, fmt.Errorf("关闭文件失败: %w", err)
	}

	// 获取文件大小
	fileInfo, err := os.Stat(fileSystemPath)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("获取文件信息失败: %w", err)
	}
	fileSize := fileInfo.Size()

	// 获取基础URL（从字典中读取）
	baseURL := e.getStorageBaseURL()

	// 计算文件的MD5哈希值（用于文件名去重）
	md5Hash, err := e.calculateFileMD5(fileSystemPath)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("计算文件MD5失败: %w", err)
	}

	// 使用MD5+扩展名作为最终文件名
	finalFileName := md5Hash + ".csv"
	finalFilePath := filepath.Join(consts.UploadDir, finalFileName)

	// 如果文件名不同，重命名文件
	if fileName != finalFileName {
		if err := os.Rename(fileSystemPath, finalFilePath); err != nil {
			return "", "", 0, 0, fmt.Errorf("重命名文件失败: %w", err)
		}
		fileSystemPath = finalFilePath
		fileName = finalFileName
	}

	// 检查文件是否已存在（根据MD5）
	fileRepo := repository.NewFileRepository(e.svcCtx)
	existingFile, err := fileRepo.FindByName(ctx, fileName)
	if err == nil && existingFile != nil {
		// 文件已存在，返回已有记录
		proxyPath := fmt.Sprintf("%s/%s", consts.PathFileUploads, fileName)
		fullURL := proxyPath
		if baseURL != "" {
			if strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://") {
				fullURL = fmt.Sprintf("%s%s", baseURL, proxyPath)
			}
		}
		return fullURL, existingFile.OriginalName, int64(existingFile.Size), recordCount, nil
	}

	// 保存文件记录到数据库
	proxyPath := fmt.Sprintf("%s/%s", consts.PathFileUploads, fileName)
	now := time.Now().Unix()
	fileModel := model.AdminFile{
		Name:         fileName,
		OriginalName: fmt.Sprintf("%s_%s.csv", moduleName, timestamp),
		Path:         proxyPath,
		BaseUrl:      baseURL,
		Size:         uint64(fileSize),
		MimeType:     sql.NullString{String: "text/csv; charset=utf-8", Valid: true},
		Ext:          sql.NullString{String: "csv", Valid: true},
		StorageType:  "local",
		Status:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    0,
	}

	if err := fileRepo.Create(ctx, &fileModel); err != nil {
		// 如果数据库保存失败，删除已创建的文件
		os.Remove(fileSystemPath)
		return "", "", 0, 0, fmt.Errorf("保存文件记录失败: %w", err)
	}

	// 构建文件URL
	fullURL := proxyPath
	if baseURL != "" {
		if strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://") {
			fullURL = fmt.Sprintf("%s%s", baseURL, proxyPath)
		}
	}

	return fullURL, fileModel.OriginalName, fileSize, recordCount, nil
}

// calculateFileMD5 计算文件MD5哈希值
func (e *ExcelExportExecutor) calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// getStorageBaseURL 从字典中获取存储baseURL
func (e *ExcelExportExecutor) getStorageBaseURL() string {
	ctx := context.Background()
	dictTypeRepo := repository.NewDictTypeRepository(e.svcCtx)
	dictType, err := dictTypeRepo.FindByCode(ctx, "storage_base_url")
	if err != nil {
		logx.Errorf("获取存储配置字典类型失败: %v", err)
		return ""
	}

	dictItemRepo := repository.NewDictItemRepository(e.svcCtx)
	items, err := dictItemRepo.FindByTypeID(ctx, dictType.Id)
	if err != nil || len(items) == 0 {
		logx.Errorf("获取存储配置字典项失败: %v", err)
		return ""
	}

	baseURL := items[0].Value
	if baseURL == "" {
		logx.Errorf("字典中的baseURL为空")
		return ""
	}

	return strings.TrimSuffix(baseURL, "/")
}
