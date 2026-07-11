// Package task 从 internal/domain/task/scheduler.go 原样搬迁而来，两处结构性改动：
//
//  1. 分布式锁改成原子操作（修复了一个真实 bug，见 docs/progress.md 本轮条目"发现 1"）：
//     原实现是 Exists + Setex 两步式，存在 TOCTOU 竞态窗口——两个调度器实例可能都读到锁不
//     存在然后都 Setex 成功，同一个任务被并发执行两次；releaseLock 的 Del 也不校验锁的
//     持有者，A 实例的锁到期后 B 实例刚获取到，A 如果这时候才执行到 defer releaseLock 会把
//     B 的锁误删。单实例运行时这两个问题都不会暴露，但 task-rpc 拆分后完全可能跑多副本，
//     必须借这次拆分改成原子 SET NX EX（SetnxExCtx）+ 持锁 token（释放前用 Lua script 校验
//     token 匹配才 DEL，标准 Redlock 安全释放模式）。
//  2. 不再直接持有 *repository.Repository/*hub.ChatHub，改成持有 task-rpc 自己的窄
//     TaskRepository 接口 + *redis.Redis；scanAsyncTasks/scanScheduledTasks 的裸 SQL 收进
//     TaskRepository.FindPendingAsync/FindPendingScheduled（原实现绕开仓储层直连 DB，是搬迁
//     顺手做的清理，不是新引入的范围）。
package task

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"postapocgame/admin-server/services/task/internal/consts"
	"postapocgame/admin-server/services/task/internal/interfaces"
	taskmodel "postapocgame/admin-server/services/task/internal/model/task"
	"postapocgame/admin-server/services/task/internal/repository"
)

// releaseLockScript 安全释放锁：只有当锁的值仍然等于本实例持有的 token 时才 DEL，
// GET+DEL 通过 Lua script 在 Redis 侧原子执行，避免释放到其他实例新获取的锁。
var releaseLockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	taskRepo      repository.TaskRepository
	redis         *redis.Redis
	ticker        *time.Ticker
	stopChan      chan struct{}
	wg            sync.WaitGroup
	notifier      *TaskNotifier
	executors     map[int]interfaces.TaskExecutor
	maxConcurrent int
	semaphore     chan struct{} // 控制并发执行的信号量
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler(taskRepo repository.TaskRepository, rds *redis.Redis, notifier *TaskNotifier, executors map[int]interfaces.TaskExecutor) *TaskScheduler {
	maxConcurrent := consts.TaskDefaultMaxConcurrent
	if maxConcurrent <= 0 {
		maxConcurrent = 10
	}

	return &TaskScheduler{
		taskRepo:      taskRepo,
		redis:         rds,
		stopChan:      make(chan struct{}),
		notifier:      notifier,
		executors:     executors,
		maxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
	}
}

// Start 启动任务调度器
func (s *TaskScheduler) Start() {
	scanInterval := time.Duration(consts.TaskDefaultScanInterval) * time.Second
	s.ticker = time.NewTicker(scanInterval)

	s.wg.Add(1)
	go s.run()

	logx.Infof("任务调度器已启动，扫描间隔：%v，最大并发：%d", scanInterval, s.maxConcurrent)
}

// Stop 停止任务调度器
func (s *TaskScheduler) Stop() {
	close(s.stopChan)
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.wg.Wait()
	logx.Infof("任务调度器已停止")
}

// run 运行调度器
func (s *TaskScheduler) run() {
	defer s.wg.Done()

	// 立即执行一次扫描
	s.scanAndExecute()

	for {
		select {
		case <-s.ticker.C:
			s.scanAndExecute()
		case <-s.stopChan:
			return
		}
	}
}

// scanAndExecute 扫描并执行任务
func (s *TaskScheduler) scanAndExecute() {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("任务调度器扫描发生 panic: %v", r)
		}
	}()

	ctx := context.Background()

	asyncTasks, err := s.taskRepo.FindPendingAsync(ctx, consts.TaskDefaultBatchSize)
	if err != nil {
		logx.Errorf("扫描异步任务失败: %v", err)
		return
	}

	scheduledTasks, err := s.taskRepo.FindPendingScheduled(ctx, consts.TaskDefaultBatchSize, time.Now().Unix())
	if err != nil {
		logx.Errorf("扫描定时任务失败: %v", err)
		return
	}

	tasks := append(asyncTasks, scheduledTasks...)

	for _, task := range tasks {
		select {
		case s.semaphore <- struct{}{}: // 获取信号量
			s.wg.Add(1)
			go s.executeTask(ctx, task)
		default:
			// 达到最大并发数，跳过本次扫描
			logx.Infof("达到最大并发数，跳过任务: taskId=%d", task.Id)
		}
	}
}

// executeTask 执行任务
func (s *TaskScheduler) executeTask(ctx context.Context, task taskmodel.AdminTask) {
	defer func() {
		<-s.semaphore // 释放信号量
		s.wg.Done()
	}()

	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("任务执行发生 panic: taskId=%d, error: %v", task.Id, r)
			s.handleTaskError(ctx, &task, fmt.Errorf("任务执行发生 panic: %v", r))
		}
	}()

	// 1. 获取分布式锁（幂等性保证，原子 SET NX EX + token）
	lockKey := fmt.Sprintf("%s%d", consts.RedisTaskLockPrefix, task.Id)
	token, locked, err := s.acquireLock(ctx, lockKey)
	if !locked {
		if err != nil {
			logx.Errorf("获取任务锁失败: taskId=%d, error: %v", task.Id, err)
		} else {
			logx.Infof("任务正在执行中，跳过: taskId=%d", task.Id)
		}
		return
	}
	defer s.releaseLock(ctx, lockKey, token)

	// 2. 再次检查任务状态（双重检查）
	currentTask, err := s.taskRepo.FindOne(ctx, task.Id)
	if err != nil {
		logx.Errorf("查询任务状态失败: taskId=%d, error: %v", task.Id, err)
		return
	}

	if currentTask.Status != consts.TaskStatusPending {
		logx.Infof("任务状态已变更，跳过执行: taskId=%d, status=%d", task.Id, currentTask.Status)
		return
	}

	// 3. 更新任务状态为"进行中"
	now := time.Now().Unix()
	if err := s.taskRepo.UpdateStatus(ctx, task.Id, consts.TaskStatusRunning, now, 0); err != nil {
		logx.Errorf("更新任务状态失败: taskId=%d, error: %v", task.Id, err)
		return
	}

	task.Status = consts.TaskStatusRunning
	task.StartedAt = now

	// 4. 发送任务开始通知
	s.notifier.NotifyTaskStatusChange(ctx, &task)

	// 5. 获取任务执行器
	executor, ok := s.executors[int(task.Type)]
	if !ok {
		err := fmt.Errorf("未找到任务执行器: taskType=%d", task.Type)
		logx.Errorf("执行任务失败: taskId=%d, error: %v", task.Id, err)
		s.handleTaskError(ctx, &task, err)
		return
	}

	// 6. 执行任务（带超时控制）
	taskTimeout := time.Duration(consts.TaskDefaultTaskTimeout) * time.Second
	ctxWithTimeout, cancel := context.WithTimeout(ctx, taskTimeout)
	defer cancel()

	paramsJSON := ""
	if task.Params.Valid {
		paramsJSON = task.Params.String
	}

	resultJSON, err := executor.Execute(ctxWithTimeout, &task, paramsJSON)

	if err != nil {
		s.handleTaskError(ctx, &task, err)
		return
	}

	// 8. 更新任务状态为"已完成"
	finishedAt := time.Now().Unix()
	if err := s.taskRepo.UpdateResult(ctx, task.Id, consts.TaskStatusCompleted, resultJSON, "", finishedAt); err != nil {
		logx.Errorf("更新任务结果失败: taskId=%d, error: %v", task.Id, err)
		return
	}

	task.Status = consts.TaskStatusCompleted
	task.FinishedAt = finishedAt
	if task.Result.Valid {
		task.Result.String = resultJSON
	}

	// 9. 发送任务完成通知
	s.notifier.NotifyTaskStatusChange(ctx, &task)

	logx.Infof("任务执行完成: taskId=%d, name=%s", task.Id, task.Name)
}

// handleTaskError 处理任务执行错误
func (s *TaskScheduler) handleTaskError(ctx context.Context, task *taskmodel.AdminTask, err error) {
	errorMessage := err.Error()
	if len(errorMessage) > 1000 {
		errorMessage = errorMessage[:1000]
	}

	resultBytes, marshalErr := json.Marshal(TaskResultResp{Success: false, Message: errorMessage})
	if marshalErr != nil {
		// 理论上不会失败（纯字符串字段），兜底成一个手写但转义安全的 JSON。
		logx.Errorf("序列化任务失败结果失败: taskId=%d, error: %v", task.Id, marshalErr)
		resultBytes, _ = json.Marshal(map[string]any{"success": false, "message": "任务执行失败（结果序列化异常）"})
	}
	resultJSON := string(resultBytes)

	finishedAt := time.Now().Unix()
	if updateErr := s.taskRepo.UpdateResult(ctx, task.Id, consts.TaskStatusFailed, resultJSON, errorMessage, finishedAt); updateErr != nil {
		logx.Errorf("更新任务失败状态失败: taskId=%d, error: %v", task.Id, updateErr)
		return
	}

	task.Status = consts.TaskStatusFailed
	task.ErrorMessage = errorMessage
	task.FinishedAt = finishedAt

	s.notifier.NotifyTaskStatusChange(ctx, task)

	logx.Errorf("任务执行失败: taskId=%d, name=%s, error: %v", task.Id, task.Name, err)
}

// acquireLock 原子获取分布式锁：SET key token NX EX ttl。返回本实例持有的 token，
// 供 releaseLock 校验所有权后再删除，避免误删其他实例持有的锁。
func (s *TaskScheduler) acquireLock(ctx context.Context, lockKey string) (string, bool, error) {
	token := uuid.NewString()
	ok, err := s.redis.SetnxExCtx(ctx, lockKey, token, consts.TaskDefaultLockTimeout)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return token, true, nil
}

// releaseLock 安全释放分布式锁：只有锁的值仍是本实例的 token 时才删除。
func (s *TaskScheduler) releaseLock(ctx context.Context, lockKey, token string) {
	if _, err := s.redis.ScriptRunCtx(ctx, releaseLockScript, []string{lockKey}, token); err != nil {
		logx.Errorf("释放任务锁失败: lockKey=%s, error: %v", lockKey, err)
	}
}
