package executors

import (
	"context"
	"database/sql"
	"encoding/csv"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/internal/consts"
	taskdomain "postapocgame/admin-server/internal/domain/task"
	monitoringmodel "postapocgame/admin-server/internal/model/monitoring"
	taskmodel "postapocgame/admin-server/internal/model/task"
	"postapocgame/admin-server/internal/repository"
)

var (
	adminFileColumns = []string{
		"id", "name", "original_name", "path", "base_url", "size",
		"mime_type", "ext", "storage_type", "status", "created_at", "updated_at", "deleted_at",
	}
	operationLogColumns = []string{
		"id", "user_id", "username", "operation_type", "operation_object", "method", "path",
		"request_params", "response_code", "response_msg", "ip_address", "user_agent",
		"duration", "created_at", "updated_at", "deleted_at",
	}
)

func newTestRepo(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
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

	return repo, sqlMock, func() {
		_ = db.Close()
		mr.Close()
		_ = os.RemoveAll(consts.UploadDir)
	}
}

func TestExecute_UnsupportedModule(t *testing.T) {
	repo, _, cleanup := newTestRepo(t)
	defer cleanup()

	executor := &ExcelExportExecutor{svcCtx: repo}
	paramsJSON := `{"module":"not_a_real_module"}`

	_, err := executor.Execute(context.Background(), &taskmodel.AdminTask{}, paramsJSON)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "不支持的导出模块")
}

func TestExportOperationLog_Success(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery("(?i)count.*admin_operation_log").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	sqlMock.ExpectQuery("(?i)from.*admin_operation_log").
		WillReturnRows(sqlmock.NewRows(operationLogColumns).
			AddRow(1, 7, "alice", "create", "user", "POST", "/api/v1/iam/user/create", nil, 200, "ok", "127.0.0.1", "ua", 12, 0, 0, 0).
			AddRow(2, 7, "alice", "update", "user", "PUT", "/api/v1/iam/user/update", nil, 200, "ok", "127.0.0.1", "ua", 8, 0, 0, 0))
	// generateCSVFile 内部：FindByName 未命中 -> 走 Create 落库
	sqlMock.ExpectQuery("(?i)from `admin_file`").WillReturnError(sql.ErrNoRows)
	sqlMock.ExpectExec("(?i)insert into `admin_file`").WillReturnResult(sqlmock.NewResult(1, 1))

	executor := &ExcelExportExecutor{svcCtx: repo}
	fileURL, fileName, fileSize, recordCount, err := executor.exportOperationLog(context.Background(), taskdomain.ExcelExportParams{
		TaskParamsReq: taskdomain.TaskParamsReq{Module: consts.TaskModuleOperationLog},
	})

	require.NoError(t, err)
	assert.EqualValues(t, 2, recordCount)
	assert.NotEmpty(t, fileURL)
	assert.NotEmpty(t, fileName) // 返回的是人类可读的原始文件名（供展示），不是磁盘上的存储文件名
	assert.Positive(t, fileSize)
	assert.NoError(t, sqlMock.ExpectationsWereMet())

	// 磁盘上实际存储的文件名是内容 MD5（见 generateCSVFile 的去重设计），从 fileURL 里取真实文件名。
	f, err := os.Open(consts.UploadDir + "/" + path.Base(fileURL))
	require.NoError(t, err)
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	require.NoError(t, err)
	require.Len(t, rows, 3)
	// encoding/csv 不会自动剥离写入的 UTF-8 BOM，第一个字段前会带上 BOM 字节。
	assert.Equal(t, "ID", strings.TrimPrefix(rows[0][0], "\ufeff"))
}

func TestGenerateCSVFile_DBWriteFailsCleansUpFile(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery("(?i)from `admin_file`").WillReturnError(sql.ErrNoRows)
	sqlMock.ExpectExec("(?i)insert into `admin_file`").WillReturnError(assert.AnError)

	executor := &ExcelExportExecutor{svcCtx: repo}
	_, _, _, _, err := executor.generateCSVFile(context.Background(), "测试模块",
		[]monitoringmodel.AdminOperationLog{},
		func(interface{}) []string { return nil },
		[]string{"ID"})

	require.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())

	// generateCSVFile 内部计算出的最终文件名是 MD5(内容).csv，这里只需确认 UploadDir 下
	// 没有残留任何 .csv 文件——补偿删除逻辑生效。
	entries, err := os.ReadDir(consts.UploadDir)
	if err == nil {
		for _, e := range entries {
			assert.NotContains(t, e.Name(), ".csv")
		}
	}
}

func TestGenerateCSVFile_FileAlreadyExists(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery("(?i)from `admin_file`").
		WillReturnRows(sqlmock.NewRows(adminFileColumns).
			AddRow(1, "existing", "existing_original.csv", "/api/v1/files/uploads/existing", "", 100,
				nil, nil, "local", 1, 0, 0, 0))

	executor := &ExcelExportExecutor{svcCtx: repo}
	_, fileName, fileSize, _, err := executor.generateCSVFile(context.Background(), "测试模块",
		[]monitoringmodel.AdminOperationLog{},
		func(interface{}) []string { return nil },
		[]string{"ID"})

	require.NoError(t, err)
	assert.Equal(t, "existing_original.csv", fileName)
	assert.EqualValues(t, 100, fileSize)
	// 命中已存在记录时不应该再执行 INSERT（没有为此设置 ExpectExec，若发生会导致 ExpectationsWereMet 失败或
	// sqlmock 报"未预期的调用"）。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
