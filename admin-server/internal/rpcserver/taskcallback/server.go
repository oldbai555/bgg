// Package taskcallback 实现 pkg/taskcallback.TaskCallbackServer，供 task-rpc 回调取导出数据 /
// 登记导出文件。当前阶段单体内嵌一个 zrpc server 提前实现这份接口（admin.go 里和 REST server
// 并存），后续 iam-rpc/sdk-rpc 真正拆分时把这个实现原样搬过去，不改契约。
//
// FetchExportData 的 5 个分支是从 internal/domain/task/executors/excel_export_executor.go 的
// export{OperationLog,AuditLog,LoginLog,PerformanceLog}（查询部分）原样迁移过来的，新增了此前
// 缺失的 sdk_call_log 分支（复用已存在的 SdkAdminRepository.ExportCallLogs，修复了一个真实
// bug：consts.TaskModuleSdkCallLog 常量和导出方法都已就绪，只是 Execute 的 switch 一直没接上，
// 命中 default 直接报错"不支持的导出模块"，见 docs/progress.md 对应条目）。
//
// RegisterExportFile 是从 generateCSVFile 的"文件已存在则复用记录，否则登记 admin_file"这段
// 尾部逻辑原样迁移过来的。
package taskcallback

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/internal/consts"
	systemmodel "postapocgame/admin-server/internal/model/system"
	"postapocgame/admin-server/internal/repository"
	monitoringrepo "postapocgame/admin-server/internal/repository/monitoring"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
	systemrepo "postapocgame/admin-server/internal/repository/system"
	pb "postapocgame/admin-server/pkg/taskcallback/pb"
)

type Server struct {
	pb.UnimplementedTaskCallbackServer
	repo *repository.Repository
}

func NewServer(repo *repository.Repository) *Server {
	return &Server{repo: repo}
}

// FetchExportData 按 module 分支查询要导出的数据，返回通用的 headers/rows_json 结构。
func (s *Server) FetchExportData(ctx context.Context, req *pb.FetchExportDataRequest) (*pb.FetchExportDataResponse, error) {
	filters := map[string]interface{}{}
	if req.FiltersJson != "" {
		if err := json.Unmarshal([]byte(req.FiltersJson), &filters); err != nil {
			return nil, fmt.Errorf("解析 filters_json 失败: %w", err)
		}
	}

	var (
		headers []string
		rows    [][]string
		err     error
	)

	switch req.Module {
	case consts.TaskModuleOperationLog:
		headers, rows, err = s.fetchOperationLog(ctx, filters)
	case consts.TaskModuleAuditLog:
		headers, rows, err = s.fetchAuditLog(ctx, filters)
	case consts.TaskModuleLoginLog:
		headers, rows, err = s.fetchLoginLog(ctx, filters)
	case consts.TaskModulePerformanceLog:
		headers, rows, err = s.fetchPerformanceLog(ctx, filters)
	case consts.TaskModuleSdkCallLog:
		headers, rows, err = s.fetchSdkCallLog(ctx, filters)
	default:
		return nil, fmt.Errorf("不支持的导出模块: %s", req.Module)
	}
	if err != nil {
		return nil, err
	}

	rowsJSON := make([]string, 0, len(rows))
	for _, row := range rows {
		obj := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(row) {
				obj[h] = row[i]
			}
		}
		b, marshalErr := json.Marshal(obj)
		if marshalErr != nil {
			return nil, fmt.Errorf("序列化导出行失败: %w", marshalErr)
		}
		rowsJSON = append(rowsJSON, string(b))
	}

	return &pb.FetchExportDataResponse{
		RowsJson:   rowsJSON,
		Headers:    headers,
		TotalCount: int64(len(rows)),
	}, nil
}

func filterString(filters map[string]interface{}, key string) string {
	if v, ok := filters[key].(string); ok {
		return v
	}
	return ""
}

func filterUint64(filters map[string]interface{}, key string) uint64 {
	if v, ok := filters[key].(float64); ok {
		return uint64(v)
	}
	return 0
}

func filterInt64(filters map[string]interface{}, key string) int64 {
	if v, ok := filters[key].(float64); ok {
		return int64(v)
	}
	return 0
}

func filterInt(filters map[string]interface{}, key string) int {
	return int(filterInt64(filters, key))
}

func (s *Server) fetchOperationLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
	repo := monitoringrepo.NewOperationLogRepository(s.repo)
	list, _, err := repo.FindPage(ctx, 1, 10000,
		filterUint64(filters, consts.TaskFilterUserId),
		filterString(filters, consts.TaskFilterUsername),
		filterString(filters, consts.TaskFilterOperationType),
		filterString(filters, consts.TaskFilterOperationObj),
		filterString(filters, consts.TaskFilterMethod),
		filterString(filters, consts.TaskFilterStartTime),
		filterString(filters, consts.TaskFilterEndTime),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("查询操作日志失败: %w", err)
	}

	headers := []string{"ID", "用户ID", "用户名", "操作类型", "操作对象", "请求方法", "请求路径", "请求参数", "响应状态码", "响应消息", "IP地址", "用户代理", "耗时(ms)", "创建时间"}
	rows := make([][]string, 0, len(list))
	for i := range list {
		opLog := &list[i]
		requestParams := ""
		if opLog.RequestParams.Valid {
			requestParams = opLog.RequestParams.String
		}
		rows = append(rows, []string{
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
		})
	}
	return headers, rows, nil
}

func (s *Server) fetchAuditLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
	repo := monitoringrepo.NewAuditLogRepository(s.repo)
	list, _, err := repo.FindPage(ctx, 1, 10000,
		filterUint64(filters, consts.TaskFilterUserId),
		filterString(filters, consts.TaskFilterUsername),
		filterString(filters, consts.TaskFilterAuditType),
		filterString(filters, consts.TaskFilterAuditObject),
		filterString(filters, consts.TaskFilterStartTime),
		filterString(filters, consts.TaskFilterEndTime),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("查询审计日志失败: %w", err)
	}

	headers := []string{"ID", "用户ID", "用户名", "审计类型", "审计对象", "审计详情", "IP地址", "用户代理", "创建时间"}
	rows := make([][]string, 0, len(list))
	for i := range list {
		auditLog := &list[i]
		auditDetail := ""
		if auditLog.AuditDetail.Valid {
			auditDetail = auditLog.AuditDetail.String
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", auditLog.Id),
			fmt.Sprintf("%d", auditLog.UserId),
			auditLog.Username,
			auditLog.AuditType,
			auditLog.AuditObject,
			auditDetail,
			auditLog.IpAddress,
			auditLog.UserAgent,
			time.Unix(auditLog.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return headers, rows, nil
}

func (s *Server) fetchLoginLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
	repo := monitoringrepo.NewLoginLogRepository(s.repo)
	list, _, err := repo.FindPage(ctx, 1, 10000,
		filterUint64(filters, consts.TaskFilterUserId),
		filterString(filters, consts.TaskFilterUsername),
		filterInt(filters, consts.TaskFilterStatus),
		filterString(filters, consts.TaskFilterStartTime),
		filterString(filters, consts.TaskFilterEndTime),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("查询登录日志失败: %w", err)
	}

	headers := []string{"ID", "用户ID", "用户名", "IP地址", "登录地点", "浏览器", "操作系统", "用户代理", "登录状态", "登录消息", "登录时间", "登出时间", "创建时间"}
	rows := make([][]string, 0, len(list))
	for i := range list {
		loginLog := &list[i]
		rows = append(rows, []string{
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
		})
	}
	return headers, rows, nil
}

func (s *Server) fetchPerformanceLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
	repo := monitoringrepo.NewPerformanceLogRepository(s.repo)
	list, _, err := repo.FindPage(ctx, 1, 10000,
		filterString(filters, consts.TaskFilterMethod),
		filterString(filters, consts.TaskFilterPath),
		filterInt64(filters, consts.TaskFilterIsSlow),
		filterInt64(filters, consts.TaskFilterStatusCode),
		filterString(filters, consts.TaskFilterStartTime),
		filterString(filters, consts.TaskFilterEndTime),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("查询性能日志失败: %w", err)
	}

	headers := []string{"ID", "用户ID", "用户名", "请求方法", "请求路径", "状态码", "耗时(ms)", "是否慢接口", "慢接口阈值(ms)", "IP地址", "User-Agent", "错误信息", "创建时间"}
	rows := make([][]string, 0, len(list))
	for i := range list {
		perf := &list[i]
		rows = append(rows, []string{
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
		})
	}
	return headers, rows, nil
}

// fetchSdkCallLog 是此前缺失的分支（发现 2 的 bug 修复）：consts.TaskModuleSdkCallLog 常量、
// SdkAdminRepository.ExportCallLogs 方法都已经存在，只是从未被 Execute 的 switch 接上过。
func (s *Server) fetchSdkCallLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
	repo := sdkrepo.NewSdkAdminRepository(s.repo)
	list, err := repo.ExportCallLogs(ctx, 2000,
		filterUint64(filters, consts.TaskFilterSdkKeyId),
		filterString(filters, consts.TaskFilterApiCode),
		filterInt64(filters, consts.TaskFilterRespCode),
		filterString(filters, consts.TaskFilterIP),
		filterInt64(filters, consts.TaskFilterStartTime),
		filterInt64(filters, consts.TaskFilterEndTime),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("查询SDK调用日志失败: %w", err)
	}

	headers := []string{"ID", "SDK Key ID", "接口编码", "响应状态码", "IP地址", "耗时(ms)", "创建时间"}
	rows := make([][]string, 0, len(list))
	for i := range list {
		row := &list[i]
		rows = append(rows, []string{
			fmt.Sprintf("%d", row.Id),
			fmt.Sprintf("%d", row.SdkKeyId),
			row.ApiCode,
			fmt.Sprintf("%d", row.RespCode),
			row.Ip,
			fmt.Sprintf("%d", row.DurationMs),
			time.Unix(row.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return headers, rows, nil
}

// RegisterExportFile 迁移自 generateCSVFile 的文件登记尾部逻辑：按 name（MD5）查重，
// 已存在则复用记录，否则新建 admin_file 记录。
func (s *Server) RegisterExportFile(ctx context.Context, req *pb.RegisterExportFileRequest) (*pb.RegisterExportFileResponse, error) {
	fileRepo := systemrepo.NewFileRepository(s.repo)
	baseURL := s.getStorageBaseURL(ctx)

	buildURL := func(storagePath string) string {
		if baseURL != "" && (strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://")) {
			return baseURL + storagePath
		}
		return storagePath
	}

	if existing, err := fileRepo.FindByName(ctx, req.FileName); err == nil && existing != nil {
		return &pb.RegisterExportFileResponse{
			FileId:    existing.Id,
			AccessUrl: buildURL(req.StoragePath),
		}, nil
	}

	now := time.Now().Unix()
	fileModel := systemmodel.AdminFile{
		Name:         req.FileName,
		OriginalName: req.OriginalName,
		Path:         req.StoragePath,
		BaseUrl:      baseURL,
		Size:         req.FileSize,
		StorageType:  "local",
		Status:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    0,
	}
	if err := fileRepo.Create(ctx, &fileModel); err != nil {
		return nil, fmt.Errorf("登记导出文件失败: %w", err)
	}

	logx.Infof("登记导出文件: fileId=%d, name=%s, uploadedBy=%d", fileModel.Id, req.FileName, req.UploadedBy)

	return &pb.RegisterExportFileResponse{
		FileId:    fileModel.Id,
		AccessUrl: buildURL(req.StoragePath),
	}, nil
}

func (s *Server) getStorageBaseURL(ctx context.Context) string {
	dictTypeRepo := systemrepo.NewDictTypeRepository(s.repo)
	dictType, err := dictTypeRepo.FindByCode(ctx, "storage_base_url")
	if err != nil {
		logx.Errorf("获取存储配置字典类型失败: %v", err)
		return ""
	}

	dictItemRepo := systemrepo.NewDictItemRepository(s.repo)
	items, err := dictItemRepo.FindByTypeID(ctx, dictType.Id)
	if err != nil || len(items) == 0 {
		logx.Errorf("获取存储配置字典项失败: %v", err)
		return ""
	}

	baseURL := items[0].Value
	if baseURL == "" {
		return ""
	}
	return strings.TrimSuffix(baseURL, "/")
}
