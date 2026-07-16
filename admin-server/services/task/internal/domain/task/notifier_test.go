package task

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"postapocgame/admin-server/services/task/internal/consts"
	taskmodel "postapocgame/admin-server/services/task/internal/model/task"
)

func TestNotifyTaskStatusChange_PublishesForRunningCompletedFailed(t *testing.T) {
	rds, cleanup := newTestRedis(t)
	defer cleanup()

	notifier := NewTaskNotifier(rds)
	ctx := context.Background()

	for _, status := range []int64{consts.TaskStatusRunning, consts.TaskStatusCompleted, consts.TaskStatusFailed} {
		notifier.NotifyTaskStatusChange(ctx, &taskmodel.AdminTask{Id: 1, Name: "t", UserId: 2, Status: status})
	}

	info, err := rds.XInfoStreamCtx(ctx, StreamTaskNotification)
	require.NoError(t, err)
	assert.EqualValues(t, 3, info.Length)
}

func TestNotifyTaskStatusChange_SkipsPending(t *testing.T) {
	rds, cleanup := newTestRedis(t)
	defer cleanup()

	notifier := NewTaskNotifier(rds)
	notifier.NotifyTaskStatusChange(context.Background(), &taskmodel.AdminTask{Id: 1, Status: consts.TaskStatusPending})

	// Pending 不发布事件，流从未被创建过，XINFO STREAM 应该报"不存在"。
	_, err := rds.XInfoStreamCtx(context.Background(), StreamTaskNotification)
	assert.Error(t, err)
}
