package task

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/interfaces"
	interfacesmocks "postapocgame/admin-server/internal/mocks/interfaces"
	taskmodel "postapocgame/admin-server/internal/model/task"
	"postapocgame/admin-server/internal/repository"
)

var adminTaskColumns = []string{
	"id", "name", "type", "execution_type", "status", "params", "result",
	"error_message", "user_id", "scheduled_at", "started_at", "finished_at",
	"created_at", "updated_at", "deleted_at",
}

func newTestRepoWithRedis(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, *miniredis.Miniredis, func()) {
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

	return repo, sqlMock, mr, func() {
		_ = db.Close()
		mr.Close()
	}
}

func TestScanAsyncTasks(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	sqlMock.ExpectQuery("(?i)from `admin_task`").
		WithArgs(0, consts.TaskExecutionTypeAsync, consts.TaskStatusPending, 0).
		WillReturnRows(sqlmock.NewRows(adminTaskColumns))

	scheduler := NewTaskScheduler(repo, nil, nil)
	tasks, err := scheduler.scanAsyncTasks(context.Background())

	require.NoError(t, err)
	assert.Empty(t, tasks)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestScanScheduledTasks(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	sqlMock.ExpectQuery("(?i)from `admin_task`").
		WithArgs(0, consts.TaskExecutionTypeAsync, consts.TaskStatusPending, 0, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(adminTaskColumns))

	scheduler := NewTaskScheduler(repo, nil, nil)
	tasks, err := scheduler.scanScheduledTasks(context.Background())

	require.NoError(t, err)
	assert.Empty(t, tasks)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

// callExecuteTask 模拟 scanAndExecute 调用 executeTask 前的前置动作（占用信号量 + wg.Add），
// 直接调用 executeTask 而不启动 goroutine，因为 executeTask 的 defer 里有 <-s.semaphore 和
// s.wg.Done()，不预先占位会导致测试阻塞在信号量接收上或 wg 计数为负而 panic。
func callExecuteTask(s *TaskScheduler, ctx context.Context, task taskmodel.AdminTask) {
	s.wg.Add(1)
	s.semaphore <- struct{}{}
	s.executeTask(ctx, task)
}

func TestExecuteTask_Success(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	task := taskmodel.AdminTask{Id: 1, Name: "导出任务", Type: 1, Status: consts.TaskStatusPending, UserId: 7}

	// 1. FindOne 双重检查
	sqlMock.ExpectQuery("(?i)from `admin_task`").
		WillReturnRows(sqlmock.NewRows(adminTaskColumns).
			AddRow(1, "导出任务", 1, consts.TaskExecutionTypeAsync, consts.TaskStatusPending, nil, nil, "", 7, 0, 0, 0, 0, 0, 0))
	// 2. UpdateStatus -> 进行中
	sqlMock.ExpectExec("(?i)update `admin_task`").WillReturnResult(sqlmock.NewResult(0, 1))
	// 3. 通知：进行中（写 admin_notification）
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).WillReturnResult(sqlmock.NewResult(1, 1))
	// 4. UpdateResult -> 已完成
	sqlMock.ExpectExec("(?i)update `admin_task`").WillReturnResult(sqlmock.NewResult(0, 1))
	// 5. 通知：已完成
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).WillReturnResult(sqlmock.NewResult(2, 1))

	executor := interfacesmocks.NewTaskExecutor(t)
	executor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return(`{"success":true}`, nil)

	scheduler := NewTaskScheduler(repo, nil, map[int]interfaces.TaskExecutor{1: executor})
	callExecuteTask(scheduler, context.Background(), task)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	executor.AssertExpectations(t)
}

func TestExecuteTask_LockHeld(t *testing.T) {
	repo, sqlMock, mr, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	task := taskmodel.AdminTask{Id: 2, Name: "导出任务", Type: 1, Status: consts.TaskStatusPending}
	lockKey := consts.RedisTaskLockPrefix + "2"
	require.NoError(t, mr.Set(lockKey, "1"))

	executor := interfacesmocks.NewTaskExecutor(t) // 不应该被调用

	scheduler := NewTaskScheduler(repo, nil, map[int]interfaces.TaskExecutor{1: executor})
	callExecuteTask(scheduler, context.Background(), task)

	// 锁已被持有，应直接跳过：没有任何 SQL 被执行，executor 也没有被调用。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
	executor.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything, mock.Anything)
}

func TestExecuteTask_ExecutorNotFound(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	task := taskmodel.AdminTask{Id: 3, Name: "未知任务", Type: 99, Status: consts.TaskStatusPending}

	sqlMock.ExpectQuery("(?i)from `admin_task`").
		WillReturnRows(sqlmock.NewRows(adminTaskColumns).
			AddRow(3, "未知任务", 99, consts.TaskExecutionTypeAsync, consts.TaskStatusPending, nil, nil, "", 0, 0, 0, 0, 0, 0, 0))
	sqlMock.ExpectExec("(?i)update `admin_task`").WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).WillReturnResult(sqlmock.NewResult(1, 1))
	// handleTaskError -> UpdateResult(失败) + 通知
	sqlMock.ExpectExec("(?i)update `admin_task`").WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).WillReturnResult(sqlmock.NewResult(2, 1))

	scheduler := NewTaskScheduler(repo, nil, map[int]interfaces.TaskExecutor{}) // 没有注册任何执行器
	callExecuteTask(scheduler, context.Background(), task)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestExecuteTask_ExecutorPanic(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	task := taskmodel.AdminTask{Id: 4, Name: "会panic的任务", Type: 1, Status: consts.TaskStatusPending}

	sqlMock.ExpectQuery("(?i)from `admin_task`").
		WillReturnRows(sqlmock.NewRows(adminTaskColumns).
			AddRow(4, "会panic的任务", 1, consts.TaskExecutionTypeAsync, consts.TaskStatusPending, nil, nil, "", 0, 0, 0, 0, 0, 0, 0))
	sqlMock.ExpectExec("(?i)update `admin_task`").WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).WillReturnResult(sqlmock.NewResult(1, 1))
	// executeTask 自身的 recover() 捕获 panic 后调用 handleTaskError -> UpdateResult(失败) + 通知
	sqlMock.ExpectExec("(?i)update `admin_task`").WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).WillReturnResult(sqlmock.NewResult(2, 1))

	executor := interfacesmocks.NewTaskExecutor(t)
	executor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Run(func(mock.Arguments) {
		panic("boom")
	}).Return("", nil)

	scheduler := NewTaskScheduler(repo, nil, map[int]interfaces.TaskExecutor{1: executor})

	require.NotPanics(t, func() {
		callExecuteTask(scheduler, context.Background(), task)
	})

	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
