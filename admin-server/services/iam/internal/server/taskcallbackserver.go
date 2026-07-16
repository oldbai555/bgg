// Package server 实现 pkg/taskcallback.TaskCallbackServer，供 task-rpc 回调取导出数据 /
// 登记导出文件。原来是单体内嵌的一个 zrpc server（和 REST server 并存），iam-rpc 拆分后
// 整体原样搬到这里，和 iam-rpc 的其他两个 gRPC service（Iam、IamCallback）同一个进程、
// 同一个端口注册（见 services/iam/iam.go），契约不变。
//
// FetchExportData 的 4 个分支（operation_log/audit_log/login_log/performance_log）是从
// internal/domain/task/executors/excel_export_executor.go 的
// export{OperationLog,AuditLog,LoginLog,PerformanceLog}（查询部分）原样迁移过来的。
// sdk_call_log 分支原来直连 SdkAdminRepository.ExportCallLogs，sdk-rpc 拆分后 sdk_call_log
// 表物理上已经不在这个进程里，改成回调 sdk-rpc 的 SdkCallLogExport（见
// services/sdk/rpc/sdk.proto、services/sdk/internal/logic/sdkcallLogexportlogic.go）。
//
// RegisterExportFile 是从 generateCSVFile 的"文件已存在则复用记录，否则登记 admin_file"这段
// 尾部逻辑原样迁移过来的，admin_file 物理上属于 iam，继续留在这里。
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/services/iam/internal/consts"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/repository"
	monitoringrepo "postapocgame/admin-server/services/iam/internal/repository/monitoring"
	systemrepo "postapocgame/admin-server/services/iam/internal/repository/system"
	pb "postapocgame/admin-server/pkg/taskcallback/pb"
	"postapocgame/admin-server/services/sdk/sdkclient"
)

type TaskCallbackServer struct {
	pb.UnimplementedTaskCallbackServer
	repo   *repository.Repository
	sdkRPC sdkclient.Sdk
}

func NewTaskCallbackServer(repo *repository.Repository, sdkRPC sdkclient.Sdk) *TaskCallbackServer {
	return &TaskCallbackServer{repo: repo, sdkRPC: sdkRPC}
}

// FetchExportData 按 module 分支查询要导出的数据，返回通用的 headers/rows_json 结构。
func (s *TaskCallbackServer) FetchExportData(ctx context.Context, req *pb.FetchExportDataRequest) (*pb.FetchExportDataResponse, error) {
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

func (s *TaskCallbackServer) fetchOperationLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
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

func (s *TaskCallbackServer) fetchAuditLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
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

func (s *TaskCallbackServer) fetchLoginLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
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

func (s *TaskCallbackServer) fetchPerformanceLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
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

// fetchSdkCallLog 回调 sdk-rpc 的 SdkCallLogExport——sdk_call_log 表物理上属于 sdk-rpc，
// 这个进程拿不到直连数据的能力了。
func (s *TaskCallbackServer) fetchSdkCallLog(ctx context.Context, filters map[string]interface{}) ([]string, [][]string, error) {
	resp, err := s.sdkRPC.SdkCallLogExport(ctx, &sdkclient.SdkCallLogExportRequest{
		MaxRows:   2000,
		SdkKeyId:  filterUint64(filters, consts.TaskFilterSdkKeyId),
		ApiCode:   filterString(filters, consts.TaskFilterApiCode),
		RespCode:  filterInt64(filters, consts.TaskFilterRespCode),
		Ip:        filterString(filters, consts.TaskFilterIP),
		StartTime: filterInt64(filters, consts.TaskFilterStartTime),
		EndTime:   filterInt64(filters, consts.TaskFilterEndTime),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("查询SDK调用日志失败: %w", err)
	}

	headers := []string{"ID", "SDK Key ID", "接口编码", "响应状态码", "IP地址", "耗时(ms)", "创建时间"}
	rows := make([][]string, 0, len(resp.List))
	for _, row := range resp.List {
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
func (s *TaskCallbackServer) RegisterExportFile(ctx context.Context, req *pb.RegisterExportFileRequest) (*pb.RegisterExportFileResponse, error) {
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
		Status:       consts.Open,
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

func (s *TaskCallbackServer) getStorageBaseURL(ctx context.Context) string {
	dictTypeRepo := systemrepo.NewDictTypeRepository(s.repo)
	dictType, err := dictTypeRepo.FindByCode(ctx, consts.DictCodeStorageBaseURL)
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
