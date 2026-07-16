// Package executors 从 internal/domain/task/executors/excel_export_executor.go 改造而来。
// 原实现的 Execute 按 params.Module 分支，每个分支都直接持有 *repository.Repository 查询
// 对应业务表——task-rpc 拆分后只有 admin_task 一张表，拿不到这些数据，必须改成回调
// TaskCallback.FetchExportData（见 pkg/taskcallback/taskcallback.proto、
// internal/rpcserver/taskcallback/server.go）取数据，本地只负责生成 CSV 文件。
//
// 因此原来的 5 个 export{OperationLog,AuditLog,LoginLog,PerformanceLog,SdkCallLog}方法
// 收敛成一个通用的 Execute + generateCSVFile：不再关心 module 内部字段结构，只处理
// TaskCallback 返回的通用 headers/rows_json。generateCSVFile 的落盘+MD5+登记逻辑本身
// （原 excel_export_executor.go:334-547）原样保留，只是把"登记进 admin_file"从直连
// systemrepo.FileRepository 换成回调 RegisterExportFile。
package executors

import (
	"context"
	"crypto/md5"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"postapocgame/admin-server/pkg/taskcallback"
	pb "postapocgame/admin-server/pkg/taskcallback/pb"
	"postapocgame/admin-server/services/task/internal/consts"
	"postapocgame/admin-server/services/task/internal/domain/task"
	taskmodel "postapocgame/admin-server/services/task/internal/model/task"
)

// moduleDisplayNames 仅用于生成文件名前缀，和原实现每个 export 方法传给 generateCSVFile
// 的 moduleName 字面量一一对应。
var moduleDisplayNames = map[string]string{
	consts.TaskModuleOperationLog:   "操作日志",
	consts.TaskModuleAuditLog:       "审计日志",
	consts.TaskModuleLoginLog:       "登录日志",
	consts.TaskModulePerformanceLog: "性能监控日志",
	consts.TaskModuleSdkCallLog:     "SDK调用日志",
}

// ModuleServiceRoute 是 module -> 拥有该模块数据的服务的 TaskCallback 客户端的静态路由表。
// 当前阶段（iam-rpc/sdk-rpc 还没真正拆分成独立进程）全部 module 指向同一个单体内嵌的
// TaskCallback server，路由表退化成"只有一条真实映射"，但接口形状按未来多目标设计，
// 拆分时只改路由表的值，不改这段执行器逻辑。见 17-async-eventing.md 第 1.3 节。
type ModuleServiceRoute map[string]taskcallback.Client

// GenericExportExecutor 通用导出任务执行器
type GenericExportExecutor struct {
	moduleRoutes ModuleServiceRoute
	// fileRegistryClient 固定指向拥有 admin_file 表的服务（当前是 iam，物理上和
	// moduleRoutes 里的值是同一个单体连接，但结构上分开，避免把"数据从哪来"和
	// "文件登记去哪"两件事绑死成同一个概念——iam-rpc 拆分之后这两者会明确分开。
	fileRegistryClient taskcallback.Client
}

// NewGenericExportExecutor 创建通用导出任务执行器
func NewGenericExportExecutor(moduleRoutes ModuleServiceRoute, fileRegistryClient taskcallback.Client) *GenericExportExecutor {
	return &GenericExportExecutor{
		moduleRoutes:       moduleRoutes,
		fileRegistryClient: fileRegistryClient,
	}
}

// GetType 获取任务类型
func (e *GenericExportExecutor) GetType() int {
	return int(consts.TaskTypeExcelExport)
}

// Execute 执行导出任务：回调 FetchExportData 取数据 -> 本地生成 CSV -> 回调
// RegisterExportFile 登记文件。
func (e *GenericExportExecutor) Execute(ctx context.Context, taskModel *taskmodel.AdminTask, paramsJSON string) (string, error) {
	var params task.ExcelExportParams
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		return "", fmt.Errorf("解析任务参数失败: %w", err)
	}

	client, ok := e.moduleRoutes[params.Module]
	if !ok {
		return "", fmt.Errorf("不支持的导出模块: %s", params.Module)
	}

	filtersJSON, err := json.Marshal(params.Filters)
	if err != nil {
		return "", fmt.Errorf("序列化筛选条件失败: %w", err)
	}

	data, err := client.FetchExportData(ctx, &pb.FetchExportDataRequest{
		Module:      params.Module,
		TaskId:      taskModel.Id,
		RequestedBy: taskModel.UserId,
		FiltersJson: string(filtersJSON),
	})
	if err != nil {
		return "", fmt.Errorf("取导出数据失败: %w", err)
	}

	moduleName := moduleDisplayNames[params.Module]
	if moduleName == "" {
		moduleName = params.Module
	}

	fileURL, fileName, fileSize, err := e.generateCSVFile(ctx, moduleName, data.Headers, data.RowsJson, taskModel.UserId)
	if err != nil {
		return "", fmt.Errorf("导出失败: %w", err)
	}

	result := task.ExcelExportResult{
		TaskResultResp: task.TaskResultResp{
			Success: true,
			Message: "导出成功",
		},
		FileURL:     fileURL,
		FileName:    fileName,
		FileSize:    fileSize,
		RecordCount: data.TotalCount,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化任务结果失败: %w", err)
	}
	return string(resultJSON), nil
}

// generateCSVFile 生成CSV文件（通用方法），从原 excel_export_executor.go 迁移，
// 落盘/BOM/MD5 去重这部分逻辑不变，"登记进 admin_file" 从直连 repository 换成回调 RPC。
func (e *GenericExportExecutor) generateCSVFile(ctx context.Context, moduleName string, headers []string, rowsJSON []string, uploadedBy uint64) (string, string, int64, error) {
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.csv", moduleName, timestamp)
	fileSystemPath := filepath.Join(consts.UploadDir, fileName)

	if err := os.MkdirAll(consts.UploadDir, 0755); err != nil {
		return "", "", 0, fmt.Errorf("创建上传目录失败: %w", err)
	}

	file, err := os.Create(fileSystemPath)
	if err != nil {
		return "", "", 0, fmt.Errorf("创建文件失败: %w", err)
	}

	if _, err := file.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		_ = file.Close()
		return "", "", 0, fmt.Errorf("写入BOM失败: %w", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.Write(headers); err != nil {
		writer.Flush()
		_ = file.Close()
		return "", "", 0, fmt.Errorf("写入CSV表头失败: %w", err)
	}

	for _, rowJSON := range rowsJSON {
		var rowObj map[string]string
		if err := json.Unmarshal([]byte(rowJSON), &rowObj); err != nil {
			writer.Flush()
			_ = file.Close()
			return "", "", 0, fmt.Errorf("解析导出行失败: %w", err)
		}
		row := make([]string, len(headers))
		for i, h := range headers {
			row[i] = rowObj[h]
		}
		if err := writer.Write(row); err != nil {
			writer.Flush()
			_ = file.Close()
			return "", "", 0, fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = file.Close()
		return "", "", 0, fmt.Errorf("刷新CSV缓冲区失败: %w", err)
	}
	if err := file.Close(); err != nil {
		return "", "", 0, fmt.Errorf("关闭文件失败: %w", err)
	}

	fileInfo, err := os.Stat(fileSystemPath)
	if err != nil {
		return "", "", 0, fmt.Errorf("获取文件信息失败: %w", err)
	}
	fileSize := fileInfo.Size()

	md5Hash, err := calculateFileMD5(fileSystemPath)
	if err != nil {
		return "", "", 0, fmt.Errorf("计算文件MD5失败: %w", err)
	}

	finalFileName := md5Hash + ".csv"
	finalFilePath := filepath.Join(consts.UploadDir, finalFileName)
	if fileName != finalFileName {
		if err := os.Rename(fileSystemPath, finalFilePath); err != nil {
			return "", "", 0, fmt.Errorf("重命名文件失败: %w", err)
		}
		fileName = finalFileName
	}

	originalName := fmt.Sprintf("%s_%s.csv", moduleName, timestamp)
	storagePath := fmt.Sprintf("%s/%s", consts.PathFileUploads, fileName)

	resp, err := e.fileRegistryClient.RegisterExportFile(ctx, &pb.RegisterExportFileRequest{
		FileName:     fileName,
		OriginalName: originalName,
		StoragePath:  storagePath,
		FileSize:     uint64(fileSize),
		UploadedBy:   uploadedBy,
	})
	if err != nil {
		// 登记失败时删除已落盘的文件，避免留下无记录的孤儿文件
		_ = os.Remove(finalFilePath)
		return "", "", 0, fmt.Errorf("登记导出文件失败: %w", err)
	}

	return resp.AccessUrl, originalName, fileSize, nil
}

func calculateFileMD5(filePath string) (string, error) {
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
