package task

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"postapocgame/admin-server/services/task/internal/consts"
)

func newTestRedis(t *testing.T) (*redis.Redis, func()) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)

	rds, err := redis.NewRedis(redis.RedisConf{Host: mr.Addr(), Type: "node"})
	require.NoError(t, err)

	return rds, mr.Close
}

// TestAcquireLock_MutualExclusion 验证发现 1 的修复：原实现是 Exists+Setex 两步式，
// 有 TOCTOU 竞态；改成 SetnxExCtx 原子操作后，同一把锁只能被一个调用者拿到。
func TestAcquireLock_MutualExclusion(t *testing.T) {
	rds, cleanup := newTestRedis(t)
	defer cleanup()

	s := &TaskScheduler{redis: rds}
	lockKey := "task:lock:1"

	token1, ok1, err1 := s.acquireLock(context.Background(), lockKey)
	require.NoError(t, err1)
	require.True(t, ok1)
	require.NotEmpty(t, token1)

	// 第二次获取同一把锁必须失败（锁已存在）。
	token2, ok2, err2 := s.acquireLock(context.Background(), lockKey)
	require.NoError(t, err2)
	assert.False(t, ok2)
	assert.Empty(t, token2)
}

// TestReleaseLock_OnlyOwnerCanRelease 验证 releaseLock 只有在持有正确 token 时才真正
// 删除锁——模拟"A 实例的锁已经被 B 实例新获取"的场景，A 迟到的 release 不应该误删 B 的锁。
func TestReleaseLock_OnlyOwnerCanRelease(t *testing.T) {
	rds, cleanup := newTestRedis(t)
	defer cleanup()

	s := &TaskScheduler{redis: rds}
	lockKey := "task:lock:1"

	token, ok, err := s.acquireLock(context.Background(), lockKey)
	require.NoError(t, err)
	require.True(t, ok)

	// 模拟另一个实例的 token 覆盖了这把锁（比如原 token 过期后被 B 实例重新 SetnxEx）。
	require.NoError(t, rds.SetCtx(context.Background(), lockKey, "other-instance-token"))

	// 用 A 实例过期前拿到的 token 释放，不应该删除 B 实例持有的锁。
	s.releaseLock(context.Background(), lockKey, token)

	val, err := rds.GetCtx(context.Background(), lockKey)
	require.NoError(t, err)
	assert.Equal(t, "other-instance-token", val, "releaseLock 不应该删除其他实例持有的锁")

	// 用正确的 token 释放，应该成功删除。
	s.releaseLock(context.Background(), lockKey, "other-instance-token")
	exists, err := rds.ExistsCtx(context.Background(), lockKey)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestNewTaskScheduler_DefaultsMaxConcurrent(t *testing.T) {
	rds, cleanup := newTestRedis(t)
	defer cleanup()

	notifier := NewTaskNotifier(rds)
	s := NewTaskScheduler(nil, rds, notifier, nil)
	assert.Equal(t, consts.TaskDefaultMaxConcurrent, s.maxConcurrent)
	assert.NotNil(t, s.semaphore)
}
