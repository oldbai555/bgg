<template>
  <div class="task-float-ball">
    <el-badge
      :value="unfinishedCount"
      :hidden="unfinishedCount === 0"
      class="task-float-ball__badge"
    >
      <el-button
        type="primary"
        circle
        class="task-float-ball__button"
        @click="togglePanel"
      >
        <el-icon><Clock /></el-icon>
      </el-button>
    </el-badge>

    <transition name="el-fade-in">
      <div v-if="panelVisible" class="task-float-ball__panel">
        <div class="task-float-ball__header">
          <span>最近任务</span>
          <el-icon class="task-float-ball__close" @click="panelVisible = false">
            <Close />
          </el-icon>
        </div>
        <div class="task-float-ball__body">
          <el-skeleton
            v-if="wsStore.recentTasksLoading"
            :rows="3"
            animated
          />
          <el-empty
            v-else-if="recentTasks.length === 0"
            description="暂无任务"
          />
          <el-scrollbar v-else class="task-float-ball__list">
            <div
              v-for="item in recentTasks"
              :key="item.id"
              class="task-float-ball__item"
              @click="handleClickTask(item)"
            >
              <div class="task-float-ball__item-main">
                <span class="task-float-ball__item-name">
                  {{ item.name }}
                </span>
                <el-tag
                  size="small"
                  :type="statusTagType(item.status)"
                  class="task-float-ball__item-status"
                >
                  {{ getStatusLabel(item.status) }}
                </el-tag>
              </div>
              <div class="task-float-ball__item-meta">
                <span class="task-float-ball__item-time">
                  {{ formatTime(item.createdAt) }}
                </span>
                <el-button
                  v-if="canDownload(item)"
                  size="small"
                  type="primary"
                  text
                  @click.stop="downloadResult(item)"
                >
                  下载
                </el-button>
              </div>
            </div>
          </el-scrollbar>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import {useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {Clock, Close} from '@element-plus/icons-vue'
import {useWebSocketStore} from '@/stores/websocket'
import {useDictOptions} from '@/composables/useDictOptions'
import type {TaskItem} from '@/api/generated/admin'

const router = useRouter()
const wsStore = useWebSocketStore()

const panelVisible = ref(false)

// 字典：任务状态
const {getLabel: getStatusLabel} = useDictOptions('task_status')

const recentTasks = computed<TaskItem[]>(() => wsStore.recentTasks || [])

// 角标逻辑：显示最近任务总数，更直观
const unfinishedCount = computed(() => recentTasks.value.length)

const togglePanel = () => {
  panelVisible.value = !panelVisible.value
  if (panelVisible.value && recentTasks.value.length === 0) {
    wsStore.refreshRecentTasks().catch((err) => {
      console.error('加载最近任务失败:', err)
    })
  }
}

const statusTagType = (status: number | undefined) => {
  // 对应字典 task_status：1=未开始，2=进行中，3=已完成，4=失败
  switch (status) {
    case 3:
      return 'success'
    case 2:
      return 'warning'
    case 4:
      return 'danger'
    default:
      return 'info'
  }
}

const formatTime = (ts: number | undefined) => {
  if (!ts) {
return ''
}
  const d = new Date(ts * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(
    d.getDate()
  )} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const getFileUrlFromResult = (item: TaskItem): string => {
  if (!item.result) {
return ''
}
  try {
    const parsed = JSON.parse(item.result as unknown as string)
    if (parsed && typeof parsed.fileUrl === 'string') {
      return parsed.fileUrl
    }
  } catch {
    // ignore
  }
  return ''
}

const canDownload = (item: TaskItem) => {
  // 仅在任务已完成且有 fileUrl 时展示下载按钮
  return item.status === 3 && !!getFileUrlFromResult(item)
}

const downloadResult = (item: TaskItem) => {
  const url = getFileUrlFromResult(item)
  if (!url) {
    ElMessage.warning('暂无可下载的文件')
    return
  }
  window.open(url, '_blank')
}

const handleClickTask = (item: {id: number}) => {
  panelVisible.value = false
  router
    .push({
      path: '/admin/system/task',
      query: {taskId: String(item.id)}
    })
    .catch((err) => {
      console.error('跳转任务列表失败:', err)
      ElMessage.error('跳转任务列表失败')
    })
}

onMounted(() => {
  // 初始化时加载一次最近任务列表
  wsStore.refreshRecentTasks().catch(() => {})
})
</script>

<style scoped lang="scss">
.task-float-ball {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 2000;

  &__button {
    width: 50px;
    height: 50px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }

  &__badge {
    .el-badge__content {
      cursor: default;
    }
  }

  &__panel {
    position: absolute;
    right: 64px;
    bottom: 0;
    width: 320px;
    max-height: 400px;
    background: var(--el-bg-color);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    border-bottom: 1px solid var(--el-border-color-lighter);
    font-size: 14px;
    font-weight: 500;
  }

  &__close {
    cursor: pointer;
  }

  &__body {
    padding: 8px 0;
  }

  &__list {
    max-height: 340px;
    padding: 0 8px 8px;
  }

  &__item {
    padding: 6px 8px;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
      background-color: var(--el-fill-color-light);
    }
  }

  &__item-main {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 2px;
  }

  &__item-name {
    font-size: 13px;
    color: var(--el-text-color-primary);
    margin-right: 8px;
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__item-status {
    flex-shrink: 0;
  }

  &__item-meta {
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
}
</style>

