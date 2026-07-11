package taskcallback_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/repository"
	taskcallbacksrv "postapocgame/admin-server/internal/rpcserver/taskcallback"
	pb "postapocgame/admin-server/pkg/taskcallback/pb"
)

func newTestRepo(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	conn := sqlx.NewSqlConnFromDB(db)

	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisConf := redis.RedisConf{Host: mr.Addr(), Type: "node"}
	cacheConf := cache.CacheConf{{RedisConf: redisConf, Weight: 100}}
	rdb, err := redis.NewRedis(redisConf)
	require.NoError(t, err)

	repo, err := repository.NewRepository(conn, cacheConf, rdb)
	require.NoError(t, err)

	return repo, mock, func() {
		_ = db.Close()
		mr.Close()
	}
}

func TestFetchExportData_OperationLog(t *testing.T) {
	repo, mock, cleanup := newTestRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "username", "operation_type", "operation_object", "method", "path",
		"request_params", "response_code", "response_msg", "ip_address", "user_agent", "duration",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(1, 2, "alice", "create", "user", "POST", "/api/v1/users", nil, 200, "ok", "127.0.0.1", "curl", 15, 1700000000, 1700000000, 0)
	mock.ExpectQuery(`(?i)select count\(\*\) from .*admin_operation_log`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`(?i)select \* from .*admin_operation_log`).WillReturnRows(rows)

	srv := taskcallbacksrv.NewServer(repo)
	resp, err := srv.FetchExportData(context.Background(), &pb.FetchExportDataRequest{
		Module:      consts.TaskModuleOperationLog,
		FiltersJson: `{}`,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	require.Len(t, resp.RowsJson, 1)

	var row map[string]string
	require.NoError(t, json.Unmarshal([]byte(resp.RowsJson[0]), &row))
	assert.Equal(t, "alice", row["用户名"])
	assert.Equal(t, "POST", row["请求方法"])
}

// TestFetchExportData_SdkCallLog 验证发现 2 的 bug 修复：sdk_call_log 此前命中 default 分支
// 直接报错"不支持的导出模块"，现在应该走通 SdkAdminRepository.ExportCallLogs。
func TestFetchExportData_SdkCallLog(t *testing.T) {
	repo, mock, cleanup := newTestRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "sdk_key_id", "sdk_interface_id", "api_code", "path", "method", "ip", "user_agent",
		"req_body", "resp_body", "resp_code", "duration_ms", "deleted_at", "created_at", "updated_at",
	}).AddRow(1, 5, 1, "video.list", "/sdk/video/list", "GET", "1.2.3.4", "sdk-client", nil, nil, 200, 30, 0, 1700000000, 1700000000)
	mock.ExpectQuery(`(?i)select \* from .*sdk_call_log`).WillReturnRows(rows)

	srv := taskcallbacksrv.NewServer(repo)
	resp, err := srv.FetchExportData(context.Background(), &pb.FetchExportDataRequest{
		Module:      consts.TaskModuleSdkCallLog,
		FiltersJson: `{}`,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	require.Len(t, resp.RowsJson, 1)

	var row map[string]string
	require.NoError(t, json.Unmarshal([]byte(resp.RowsJson[0]), &row))
	assert.Equal(t, "video.list", row["接口编码"])
}

func TestFetchExportData_UnsupportedModule(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	srv := taskcallbacksrv.NewServer(repo)
	_, err := srv.FetchExportData(context.Background(), &pb.FetchExportDataRequest{
		Module:      "not_a_real_module",
		FiltersJson: `{}`,
	})
	require.Error(t, err)
}

func TestRegisterExportFile_CreatesNewRecord(t *testing.T) {
	repo, mock, cleanup := newTestRepo(t)
	defer cleanup()

	// getStorageBaseURL（RegisterExportFile 里先于 FindByName 调用）：字典类型查不到，
	// 退化成空 baseURL，不影响主流程。
	mock.ExpectQuery(`(?i)select .* from .*admin_dict_type`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description", "status", "created_at", "updated_at", "deleted_at"}))

	// FindByName：查不到已有文件
	mock.ExpectQuery(`(?i)select .* from .*admin_file.*where.*name`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "original_name", "path", "base_url", "size", "mime_type", "ext", "storage_type", "status", "created_at", "updated_at", "deleted_at"}))

	// Create：insert admin_file
	mock.ExpectExec(`(?i)insert into .*admin_file`).
		WillReturnResult(sqlmock.NewResult(42, 1))

	srv := taskcallbacksrv.NewServer(repo)
	resp, err := srv.RegisterExportFile(context.Background(), &pb.RegisterExportFileRequest{
		FileName:     "abc123.csv",
		OriginalName: "操作日志_20260711.csv",
		StoragePath:  "/api/v1/files/uploads/abc123.csv",
		FileSize:     1024,
		UploadedBy:   1,
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(42), resp.FileId)
	assert.Equal(t, "/api/v1/files/uploads/abc123.csv", resp.AccessUrl)
}
