// 职责边界：WebSocket 推送的未读消息列表（聊天/任务/系统通知）+ 已读状态。
// 连接生命周期属于 stores/websocket.ts；本 store 只订阅其 lastMessage（单向依赖，不反向 import）。
import {defineStore} from 'pinia'
import {ref, computed, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {useWebSocketStore, MessageType} from './websocket'
import {useUserStore} from './user'
import type {WSMessage} from './websocket'

// 未读消息
export interface UnreadMessage {
  id: string;
  type: MessageType | string;
  title: string;
  content: string;
  timestamp: number;
  read: boolean;
}

export const useNotificationStore = defineStore('notification', () => {
  const unreadMessages = ref<UnreadMessage[]>([])

  const unreadCount = computed(() => unreadMessages.value.filter((m) => !m.read).length)
  const hasUnreadChat = computed(() =>
    unreadMessages.value.some((m) => !m.read && m.type === MessageType.CHAT)
  )

  // 添加未读消息
  function addUnreadMessage(message: UnreadMessage) {
    unreadMessages.value.unshift(message)
    // 限制未读消息数量，最多保留 50 条
    if (unreadMessages.value.length > 50) {
      unreadMessages.value = unreadMessages.value.slice(0, 50)
    }
  }

  // 标记消息为已读
  function markAsRead(messageId: string) {
    const message = unreadMessages.value.find((m) => m.id === messageId)
    if (message) {
      message.read = true
    }
  }

  // 标记所有消息为已读
  function markAllAsRead() {
    unreadMessages.value.forEach((m) => {
      m.read = true
    })
  }

  // 清除已读消息
  function clearReadMessages() {
    unreadMessages.value = unreadMessages.value.filter((m) => !m.read)
  }

  // 处理聊天消息
  function handleChatMessage(data: WSMessage) {
    const userStore = useUserStore()
    const currentUserId = userStore.profile?.id || 0

    // 只有不是自己发的消息才需要显示未读通知
    const isMyMessage = data.fromId && Number(data.fromId) === Number(currentUserId)
    if (isMyMessage) {
      return
    }

    // 检查当前是否在聊天页面（支持多种路径格式）
    // 注意：使用 hash 路由时，路径在 hash 中
    const currentPath = window.location.pathname
    const currentHash = window.location.hash
    // 检查是否在 ChatList.vue 页面（/chatroom/chat）
    // 注意：/chatroom/chat 是 admin_menu.path（URL 路由），Phase 1 域目录重组只改了 component 字段
    // （chatroom/ChatList → chat/ChatList），path 字段未变，这里沿用旧路径是正确的，不是遗漏
    const isInChatListPage = currentHash === '#/chatroom/chat' ||
                             currentHash.startsWith('#/chatroom/chat?') ||
                             currentPath.includes('/chatroom/chat')
    const isInOtherChatPage = currentHash.includes('/temp/chat') ||
                              currentHash.includes('/chat') ||
                              currentPath.includes('/temp/chat') ||
                              currentPath.includes('/chat')
    const isInChatPage = isInChatListPage || isInOtherChatPage

    // 如果不在聊天页面，添加到未读消息并显示提示
    if (!isInChatPage) {
      const chatId = data.chatId || 0
      const isGroupChat = chatId > 0
      const title = isGroupChat
        ? `群聊消息：来自 ${data.fromName || '未知用户'}`
        : `来自 ${data.fromName || '未知用户'}`

      const content = data.content || ''
      const displayContent = data.messageType === 2 ? '[图片]' : content

      addUnreadMessage({
        id: `chat_${data.messageId || Date.now()}_${chatId}`,
        type: MessageType.CHAT,
        title: title,
        content: displayContent,
        timestamp: Date.now(),
        read: false
      })

      ElMessage.info({
        message: `${title}: ${displayContent}`,
        duration: 3000,
        showClose: true
      })
    }
  }

  // 处理任务进度
  function handleTaskProgress(data: WSMessage) {
    if (data.taskName) {
      addUnreadMessage({
        id: `task_${data.taskId || Date.now()}`,
        type: MessageType.TASK_PROGRESS,
        title: `任务进度: ${data.taskName}`,
        content: `进度: ${data.progress || 0}% - ${data.status || ''}`,
        timestamp: Date.now(),
        read: false
      })
    }
  }

  // 处理通知消息
  function handleNotification(data: WSMessage) {
    const level = data.level || 'info'
    if (data.taskId) {
      // 任务相关通知：提示 + 记录未读，任务浮球的刷新由 websocket store 自己处理
      ElMessage[level](data.content || data.title || '任务状态已更新')
      addUnreadMessage({
        id: `task_notify_${data.taskId || Date.now()}`,
        type: MessageType.NOTIFICATION,
        title: data.title || '任务通知',
        content: data.content || '',
        timestamp: Date.now(),
        read: false
      })
      return
    }

    ElMessage[level](data.content || data.title || '新通知')
    addUnreadMessage({
      id: `notify_${Date.now()}`,
      type: MessageType.NOTIFICATION,
      title: data.title || '通知',
      content: data.content || '',
      timestamp: Date.now(),
      read: false
    })
  }

  // 订阅 websocket store 广播的原始消息（单向依赖：连接 store 广播，本 store 消费）
  const wsStore = useWebSocketStore()
  watch(
    () => wsStore.lastMessage,
    (data) => {
      if (!data) {
        return
      }
      switch (data.type) {
        case MessageType.CHAT:
          handleChatMessage(data)
          break
        case MessageType.TASK_PROGRESS:
          handleTaskProgress(data)
          break
        case MessageType.NOTIFICATION:
          handleNotification(data)
          break
        case MessageType.SYSTEM:
          // 系统消息通常不需要添加到未读消息列表
          break
        default:
          break
      }
    }
  )

  return {
    unreadMessages,
    unreadCount,
    hasUnreadChat,
    addUnreadMessage,
    markAsRead,
    markAllAsRead,
    clearReadMessages
  }
})
