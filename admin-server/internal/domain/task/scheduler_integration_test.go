//go:build integration

package task

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/interfaces"
	taskmodel "postapocgame/admin-server/internal/model/task"
	"postapocgame/admin-server/internal/repository"
	taskrepo "postapocgame/admin-server/internal/repository/task"
)

// noopExecutor 是仅供集成测试使用的最小任务执行器，不代表任何真实业务类型。
type noopExecutor struct{ taskType int }

func (e *noopExecutor) GetType() int { return e.taskType }

func (e *noopExecutor) Execute(ctx context.Context, task *taskmodel.AdminTask, paramsJSON string) (string, error) {
	return `{"ok":true}`, nil
}

// 08-testing-strategy.md §5 场景 5：Task 调度器跑一个完整周期——提交任务 → 调度器拾取 →
// 执行 → 状态更新为完成。直接调用包内未导出的 scanAndExecute（本文件与 scheduler.go 同包），
// 避免等待 consts.TaskDefaultScanInterval 的真实 ticker。
func TestIntegration_TaskScheduler_FullCycle(t *testing.T) {
	dsn := os.Getenv("TEST_MYSQL_DSN")
	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	if dsn == "" || redisAddr == "" {
		t.Skip("TEST_MYSQL_DSN/TEST_REDIS_ADDR 未设置，跳过集成测试")
	}

	conn := sqlx.NewMysql(dsn)
	rawDB, err := conn.RawDB()
	require.NoError(t, err)
	require.NoError(t, rawDB.Ping())

	rdb, err := redis.NewRedis(redis.RedisConf{Host: redisAddr, Type: "node"})
	require.NoError(t, err)
	cacheConf := cache.CacheConf{{RedisConf: redis.RedisConf{Host: redisAddr, Type: "node"}, Weight: 100}}
	repo, err := repository.NewRepository(conn, cacheConf, rdb)
	require.NoError(t, err)

	ctx := context.Background()
	tr := taskrepo.NewTaskRepository(repo)

	const testTaskType = 999999 // 不与任何真实字典 task_type 冲突的测试专用类型
	taskID, err := tr.Create(ctx, &taskmodel.AdminTask{
		Name:          "集成测试任务",
		Type:          testTaskType,
		ExecutionType: consts.TaskExecutionTypeAsync,
		Status:        consts.TaskStatusPending,
		UserId:        0,
		ScheduledAt:   0, // 立即执行
	})
	require.NoError(t, err)
	require.NotZero(t, taskID)

	executors := map[int]interfaces.TaskExecutor{
		testTaskType: &noopExecutor{taskType: testTaskType},
	}
	scheduler := NewTaskScheduler(repo, nil, executors)
	scheduler.scanAndExecute()

	deadline := time.Now().Add(3 * time.Second)
	var finalStatus int64
	for time.Now().Before(deadline) {
		got, err := tr.FindOne(ctx, taskID)
		require.NoError(t, err)
		finalStatus = got.Status
		if finalStatus == consts.TaskStatusCompleted || finalStatus == consts.TaskStatusFailed {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	assert.Equal(t, consts.TaskStatusCompleted, finalStatus, "调度器应该拾取任务并在一个周期内执行完成")
}
