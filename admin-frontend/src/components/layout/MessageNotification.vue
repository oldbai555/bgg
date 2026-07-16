<template>
  <el-popover
    placement="bottom-end"
    :width="350"
    trigger="click"
    popper-class="message-notification-popover"
    @show="handlePopoverShow"
  >
    <template #reference>
      <el-badge :value="unreadCount" :hidden="unreadCount === 0" :max="99">
        <el-button text circle class="app-header__action-btn">
          <el-icon :size="18">
            <Bell />
          </el-icon>
        </el-button>
      </el-badge>
    </template>

    <div class="message-notification">
      <div class="message-notification__header">
        <span class="message-notification__title">消息通知</span>
        <div class="message-notification__actions">
          <el-button
            text
            size="small"
            :loading="readAllLoading"
            @click="handleMarkAllAsRead"
          >全部已读</el-button>
          <el-button
            text
            size="small"
            :loading="clearReadLoading"
            @click="handleClearRead"
          >清除已读</el-button>
        </div>
      </div>

      <div class="message-notification__list">
        <div
          v-for="message in displayMessages"
          :key="message.id"
          class="message-notification__item"
          :class="{ unread: message.readStatus === 1 }"
          @click="handleMessageClick(message)"
        >
          <div class="message-notification__item-icon">
            <el-icon v-if="message.sourceType === 'chat'" :size="20">
              <ChatDotRound />
            </el-icon>
            <el-icon v-else-if="message.sourceType === 'notice'" :size="20">
              <Bell />
            </el-icon>
            <el-icon v-else :size="20">
              <Bell />
            </el-icon>
          </div>
          <div class="message-notification__item-content">
            <div class="message-notification__item-title">{{ message.title }}</div>
            <div class="message-notification__item-text">{{ message.content }}</div>
            <div class="message-notification__item-time">{{ formatTime(message.createdAt) }}</div>
          </div>
          <div v-if="message.readStatus === 1" class="message-notification__item-dot"></div>
        </div>

        <el-empty
          v-if="displayMessages.length === 0"
          description="暂无消息"
          :image-size="80"
        />
      </div>

      <div v-if="displayMessages.length > 0" class="message-notification__footer">
        <el-button text size="small" @click="handleViewAll">查看全部</el-button>
      </div>
    </div>
  </el-popover>

  <!-- 公告阅读框 -->
  <NoticeReader
    v-model:visible="noticeReaderVisible"
    :notices="noticeNotifications"
    @read="handleNoticeRead"
  />
</template>

<script setup lang="ts">
import {computed, ref, onMounted, onUnmounted, watch} from 'vue'
import {useRouter} from 'vue-router'
import {Bell, ChatDotRound} from '@element-plus/icons-vue'
import {systemApi} from '@/api/system'
import type {NotificationItem} from '@/api/generated/admin'
import {ElMessage} from 'element-plus'
import NoticeReader from '@/components/common/NoticeReader.vue'
import {MessageType} from '@/stores/websocket'
import {useNotificationStore} from '@/stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

const notifications = ref<NotificationItem[]>([])
const loading = ref(false)
const readAllLoading = ref(false)
const clearReadLoading = ref(false)
const noticeReaderVisible = ref(false)
const noticeNotifications = ref<NotificationItem[]>([])

// 合并API通知和WebSocket消息
const allUnreadMessages = computed(() => {
  const apiNotifications = notifications.value.filter(n => n.readStatus === 1)
  const wsMessages = notificationStore.unreadMessages.filter(m => !m.read && m.type === MessageType.CHAT)

  // 将WebSocket消息转换为通知格式
  const wsNotifications: NotificationItem[] = wsMessages.map((msg, index) => ({
    id: Number(msg.id.replace(/\D/g, '')) || Date.now() + index, // 从ID中提取数字
    userId: 0, // WebSocket消息没有userId
    sourceType: 'chat',
    sourceId: 0,
    title: msg.title,
    content: msg.content,
    readStatus: msg.read ? 2 : 1, // 字典 read_status：1=未读，2=已读
    readAt: 0,
    createdAt: Math.floor(msg.timestamp / 1000), // 转换为秒级时间戳
    updatedAt: Math.floor(msg.timestamp / 1000)
  }))

  // 合并并去重（根据sourceType和content）
  const merged = [...apiNotifications, ...wsNotifications]
  const unique = merged.filter((item, index, self) =>
    index === self.findIndex(t =>
      t.sourceType === item.sourceType &&
      t.content === item.content &&
      Math.abs(t.createdAt - item.createdAt) < 5 // 5秒内的相同消息视为重复
    )
  )

  // 按时间倒序排序
  return unique.sort((a, b) => b.createdAt - a.createdAt)
})

const unreadCount = computed(() => allUnreadMessages.value.length)
const displayMessages = computed(() => allUnreadMessages.value.slice(0, 10)) // 最多显示 10 条

// 聊天页面路径（从字典读取）
const chatPath = ref('/admin/chatroom/chat') // 默认值

const formatTime = (timestamp: number) => {
  if (!timestamp) {
    return '-'
  }
  try {
    const date = new Date(timestamp * 1000) // 秒级时间戳转换为毫秒
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)

    if (minutes < 1) {
      return '刚刚'
    } else if (minutes < 60) {
      return `${minutes}分钟前`
    } else {
      return date.toLocaleString('zh-CN', {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }
  } catch {
    return '-'
  }
}

const loadNotifications = async () => {
  loading.value = true
  try {
    // 只查询未读消息（字典 read_status：1=未读）
    const resp = await systemApi.notificationList({
      page: 1,
      pageSize: 100,
      readStatus: 1
    })
    notifications.value = resp.list || []
  } catch (err: unknown) {
    console.error('加载消息通知失败:', err)
    notifications.value = []
  } finally {
    loading.value = false
  }
}

const handlePopoverShow = () => {
  loadNotifications()
}

// 加载聊天页面路径配置
const loadChatPath = async () => {
  try {
    const resp = await systemApi.dictGet({code: 'chat_config'})
    if (resp && resp.items && resp.items.length > 0) {
      // 查找"在线聊天页面路径"配置项
      const pathItem = resp.items.find(item => item.label === '在线聊天页面路径')
      if (pathItem && pathItem.value) {
        chatPath.value = pathItem.value
      }
    }
  } catch (err: unknown) {
    console.warn('加载聊天页面路径配置失败，使用默认值:', err)
    chatPath.value = '/admin/chatroom/chat'
  }
}

const handleMessageClick = async (message: NotificationItem) => {
  // 根据消息来源类型处理
  if (message.sourceType === 'chat') {
    // 如果是WebSocket消息（userId为0），标记为已读
    if (message.userId === 0) {
      const wsMessage = notificationStore.unreadMessages.find(m =>
        m.type === MessageType.CHAT &&
        m.title === message.title &&
        Math.abs(m.timestamp - message.createdAt * 1000) < 5000
      )
      if (wsMessage) {
        notificationStore.markAsRead(wsMessage.id)
      }
    } else {
      // API通知，标记为已读
      try {
        await systemApi.notificationRead({id: message.id})
        await loadNotifications()
      } catch (err: unknown) {
        console.error('标记通知已读失败:', err)
      }
    }

    // 聊天消息：跳转到在线聊天页面
    router.push(chatPath.value)
  } else if (message.sourceType === 'notice') {
    // 公告消息：打开公告阅读框
    // 获取所有未读的公告通知
    const noticeNotifs = notifications.value.filter(n => n.readStatus === 1 && n.sourceType === 'notice')
    if (noticeNotifs.length > 0) {
      noticeNotifications.value = noticeNotifs
      noticeReaderVisible.value = true
    }
  }
}

const handleNoticeRead = async (notificationId: number) => {
  try {
    // 标记单个通知为已读
    await systemApi.notificationRead({id: notificationId})
    // 重新加载通知列表
    await loadNotifications()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '标记已读失败'
    ElMessage.error(message)
  }
}

const handleMarkAllAsRead = async () => {
  readAllLoading.value = true
  try {
    // 标记后端API通知为已读
    await systemApi.notificationReadAll()
    // 标记WebSocket消息为已读
    notificationStore.markAllAsRead()
    ElMessage.success('全部已读成功')
    await loadNotifications()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '操作失败'
    ElMessage.error(message)
  } finally {
    readAllLoading.value = false
  }
}

const handleClearRead = async () => {
  clearReadLoading.value = true
  try {
    await systemApi.notificationClearRead()
    notificationStore.clearReadMessages()
    ElMessage.success('清除已读消息成功')
    await loadNotifications()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '操作失败'
    ElMessage.error(message)
  } finally {
    clearReadLoading.value = false
  }
}

const handleViewAll = () => {
  router.push('/admin/system/notification')
}

// 监听公告阅读框关闭
watch(noticeReaderVisible, (newVal) => {
  if (!newVal) {
    // 关闭时重新加载通知列表
    loadNotifications()
  }
})

let refreshTimer: number | undefined

onMounted(() => {
  loadChatPath()
  loadNotifications()

  // 定期刷新通知列表（每30秒）
  refreshTimer = window.setInterval(() => {
    loadNotifications()
  }, 30000)
})

onUnmounted(() => {
  if (refreshTimer !== undefined) {
    window.clearInterval(refreshTimer)
  }
})
</script>

<style scoped lang="scss">
.message-notification {
  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid var(--el-border-color-light);
  }

  &__title {
    font-weight: 600;
    font-size: 16px;
  }

  &__actions {
    display: flex;
    gap: 8px;
  }

  &__list {
    max-height: 400px;
    overflow-y: auto;
  }

  &__item {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px 16px;
    cursor: pointer;
    transition: background-color 0.2s;
    position: relative;

    &:hover {
      background-color: var(--el-fill-color-light);
    }

    &.unread {
      background-color: var(--el-color-primary-light-9);
    }

    &-icon {
      flex-shrink: 0;
      color: var(--el-color-primary);
    }

    &-content {
      flex: 1;
      min-width: 0;
    }

    &-title {
      font-weight: 500;
      font-size: 14px;
      margin-bottom: 4px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    &-text {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-bottom: 4px;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }

    &-time {
      font-size: 11px;
      color: var(--el-text-color-placeholder);
    }

    &-dot {
      position: absolute;
      top: 16px;
      right: 16px;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background-color: var(--el-color-primary);
    }
  }

  &__footer {
    padding: 8px 16px;
    border-top: 1px solid var(--el-border-color-light);
    text-align: center;
  }
}
</style>
