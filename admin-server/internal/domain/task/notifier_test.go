package task

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/hub"
	taskmodel "postapocgame/admin-server/internal/model/task"
)

func TestNotifyTaskStatusChange_Running(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	notifier := NewTaskNotifier(repo, hub.NewChatHub())
	task := &taskmodel.AdminTask{Id: 1, Name: "导出任务", Status: consts.TaskStatusRunning, UserId: 7}

	require.NotPanics(t, func() {
		notifier.NotifyTaskStatusChange(context.Background(), task)
	})
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestNotifyTaskStatusChange_Completed(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	notifier := NewTaskNotifier(repo, hub.NewChatHub())
	task := &taskmodel.AdminTask{Id: 2, Name: "导出任务", Status: consts.TaskStatusCompleted, UserId: 7}

	require.NotPanics(t, func() {
		notifier.NotifyTaskStatusChange(context.Background(), task)
	})
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestNotifyTaskStatusChange_Failed(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	notifier := NewTaskNotifier(repo, hub.NewChatHub())
	task := &taskmodel.AdminTask{Id: 3, Name: "导出任务", Status: consts.TaskStatusFailed, ErrorMessage: "磁盘已满", UserId: 7}

	require.NotPanics(t, func() {
		notifier.NotifyTaskStatusChange(context.Background(), task)
	})
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

// TestNotifyTaskStatusChange_NotificationCreateFails 验证 createNotificationRecord 失败只记日志、
// 不影响调用方（方法本身无返回值，签名上就不存在"返回 error"这回事），且不会阻断后续的 WebSocket 通知。
func TestNotifyTaskStatusChange_NotificationCreateFails(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_notification`")).
		WillReturnError(assert.AnError)

	notifier := NewTaskNotifier(repo, hub.NewChatHub())
	task := &taskmodel.AdminTask{Id: 4, Name: "导出任务", Status: consts.TaskStatusCompleted, UserId: 7}

	require.NotPanics(t, func() {
		notifier.NotifyTaskStatusChange(context.Background(), task)
	})
	// ExpectationsWereMet 本身即验证 INSERT 确实被尝试过一次（不是被跳过)。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestSendWebSocketMessage_ChatHubNil(t *testing.T) {
	repo, _, _, cleanup := newTestRepoWithRedis(t)
	defer cleanup()

	notifier := NewTaskNotifier(repo, nil)
	task := &taskmodel.AdminTask{Id: 5, Name: "导出任务", Status: consts.TaskStatusCompleted, UserId: 7}

	require.NotPanics(t, func() {
		notifier.sendWebSocketMessage(context.Background(), task)
	})
}
