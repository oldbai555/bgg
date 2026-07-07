package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/hub"
	"postapocgame/admin-server/internal/interfaces"

	"github.com/zeromicro/go-zero/core/logx"
	taskmodel "postapocgame/admin-server/internal/model/task"
	"postapocgame/admin-server/internal/repository"
	taskrepo "postapocgame/admin-server/internal/repository/task"
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	repo          *repository.Repository
	chatHub       *hub.ChatHub
	ticker        *time.Ticker
	stopChan      chan struct{}
	wg            sync.WaitGroup
	notifier      *TaskNotifier
	executors     map[int]interfaces.TaskExecutor
	maxConcurrent int
	semaphore     chan struct{} // 控制并发执行的信号量
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler(repo *repository.Repository, chatHub *hub.ChatHub, executors map[int]interfaces.TaskExecutor) *TaskScheduler {
	maxConcurrent := consts.TaskDefaultMaxConcurrent
	if maxConcurrent <= 0 {
		maxConcurrent = 10
	}

	return &TaskScheduler{
		repo:          repo,
		chatHub:       chatHub,
		stopChan:      make(chan struct{}),
		notifier:      NewTaskNotifier(repo, chatHub),
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
	taskRepo := taskrepo.NewTaskRepository(s.repo)

	// 1. 扫描待执行的异步任务（execution_type=2, status=1, scheduled_at=0）
	asyncTasks, err := s.scanAsyncTasks(ctx, taskRepo)
	if err != nil {
		logx.Errorf("扫描异步任务失败: %v", err)
		return
	}

	// 2. 扫描待执行的定时任务（execution_type=2, status=1, scheduled_at>0且<=now）
	scheduledTasks, err := s.scanScheduledTasks(ctx, taskRepo)
	if err != nil {
		logx.Errorf("扫描定时任务失败: %v", err)
		return
	}

	// 合并任务列表
	tasks := append(asyncTasks, scheduledTasks...)

	// 3. 并发执行任务（受maxConcurrent限制）
	for _, task := range tasks {
		// 检查是否达到最大并发数
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

// scanAsyncTasks 扫描异步任务
func (s *TaskScheduler) scanAsyncTasks(ctx context.Context, taskRepo taskrepo.TaskRepository) ([]taskmodel.AdminTask, error) {
	// 使用自定义SQL查询异步任务
	// 查询条件：execution_type=2（异步），status=1（未开始），scheduled_at=0（立即执行），deleted_at=0（未删除）
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"execution_type": consts.TaskExecutionTypeAsync},
		sq.Eq{"status": consts.TaskStatusPending},
		sq.Eq{"scheduled_at": 0},
	}

	sqlStr, args, err := sq.Select("*").
		From("`admin_task`").
		Where(conditions).
		OrderBy("created_at ASC").
		Limit(uint64(consts.TaskDefaultBatchSize)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("异步任务SQL生成失败: %w", err)
	}

	var tasks []taskmodel.AdminTask
	if err := s.repo.DB.QueryRowsCtx(ctx, &tasks, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("查询异步任务失败: %w", err)
	}

	return tasks, nil
}

// scanScheduledTasks 扫描定时任务
func (s *TaskScheduler) scanScheduledTasks(ctx context.Context, taskRepo taskrepo.TaskRepository) ([]taskmodel.AdminTask, error) {
	now := time.Now().Unix()

	// 使用自定义SQL查询定时任务
	// 查询条件：execution_type=2（异步），status=1（未开始），scheduled_at>0且<=now，deleted_at=0（未删除）
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"execution_type": consts.TaskExecutionTypeAsync},
		sq.Eq{"status": consts.TaskStatusPending},
		sq.Gt{"scheduled_at": 0},
		sq.LtOrEq{"scheduled_at": now},
	}

	sqlStr, args, err := sq.Select("*").
		From("`admin_task`").
		Where(conditions).
		OrderBy("scheduled_at ASC").
		Limit(uint64(consts.TaskDefaultBatchSize)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("定时任务SQL生成失败: %w", err)
	}

	var tasks []taskmodel.AdminTask
	if err := s.repo.DB.QueryRowsCtx(ctx, &tasks, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("查询定时任务失败: %w", err)
	}

	return tasks, nil
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

	// 1. 获取分布式锁（幂等性保证）
	lockKey := fmt.Sprintf("%s%d", consts.RedisTaskLockPrefix, task.Id)
	locked, err := s.acquireLock(ctx, lockKey)
	if !locked {
		if err != nil {
			logx.Errorf("获取任务锁失败: taskId=%d, error: %v", task.Id, err)
		} else {
			logx.Infof("任务正在执行中，跳过: taskId=%d", task.Id)
		}
		return
	}
	defer s.releaseLock(ctx, lockKey)

	// 2. 再次检查任务状态（双重检查）
	taskRepo := taskrepo.NewTaskRepository(s.repo)
	currentTask, err := taskRepo.FindOne(ctx, task.Id)
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
	err = taskRepo.UpdateStatus(ctx, task.Id, consts.TaskStatusRunning, now, 0)
	if err != nil {
		logx.Errorf("更新任务状态失败: taskId=%d, error: %v", task.Id, err)
		return
	}

	// 更新task对象的状态和时间
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

	// 获取任务参数JSON
	paramsJSON := ""
	if task.Params.Valid {
		paramsJSON = task.Params.String
	}

	// 执行任务
	resultJSON, err := executor.Execute(ctxWithTimeout, &task, paramsJSON)

	// 7. 处理执行结果
	if err != nil {
		s.handleTaskError(ctx, &task, err)
		return
	}

	// 8. 更新任务状态为"已完成"
	finishedAt := time.Now().Unix()
	err = taskRepo.UpdateResult(ctx, task.Id, consts.TaskStatusCompleted, resultJSON, "", finishedAt)
	if err != nil {
		logx.Errorf("更新任务结果失败: taskId=%d, error: %v", task.Id, err)
		return
	}

	// 更新task对象的状态和结果
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

	// 构建失败结果JSON
	resultJSON := fmt.Sprintf(`{"success":false,"message":"%s"}`, errorMessage)

	// 更新任务状态为"失败"
	taskRepo := taskrepo.NewTaskRepository(s.repo)
	finishedAt := time.Now().Unix()
	updateErr := taskRepo.UpdateResult(ctx, task.Id, consts.TaskStatusFailed, resultJSON, errorMessage, finishedAt)
	if updateErr != nil {
		logx.Errorf("更新任务失败状态失败: taskId=%d, error: %v", task.Id, updateErr)
		return
	}

	// 更新task对象的状态和错误信息
	task.Status = consts.TaskStatusFailed
	task.ErrorMessage = errorMessage
	task.FinishedAt = finishedAt

	// 发送任务失败通知
	s.notifier.NotifyTaskStatusChange(ctx, task)

	logx.Errorf("任务执行失败: taskId=%d, name=%s, error: %v", task.Id, task.Name, err)
}

// acquireLock 获取分布式锁
func (s *TaskScheduler) acquireLock(ctx context.Context, lockKey string) (bool, error) {
	// 检查锁是否已存在
	exists, err := s.repo.Redis.Exists(lockKey)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil // 锁已存在
	}

	// 设置锁（使用Setex，过期时间使用常量）
	err = s.repo.Redis.Setex(lockKey, "1", int(consts.TaskDefaultLockTimeout))
	if err != nil {
		return false, err
	}

	return true, nil
}

// releaseLock 释放分布式锁
func (s *TaskScheduler) releaseLock(ctx context.Context, lockKey string) {
	_, err := s.repo.Redis.Del(lockKey)
	if err != nil {
		logx.Errorf("释放任务锁失败: lockKey=%s, error: %v", lockKey, err)
	}
}
