// 职责边界：只管 WebSocket 连接生命周期（连接/重连/心跳）+ 原始消息广播（lastMessage）+ 任务浮球的最近任务列表。
// 未读消息列表/已读状态属于 stores/notification.ts，不要把两者揉回同一个 store。
import {defineStore} from 'pinia'
import {useUserStore} from './user'
import {getWebSocketBaseURL} from '@/composables/useAppConfig'
import {taskApi} from '@/api/task'
import type {TaskItem} from '@/api/generated/admin'

// WebSocket 消息类型
export enum MessageType {
  CHAT = 'chat', // 聊天消息
  TASK_PROGRESS = 'task_progress', // 任务进度
  NOTIFICATION = 'notification', // 通知消息
  SYSTEM = 'system' // 系统消息
}

// WebSocket 消息结构
export interface WSMessage {
  type: MessageType | string;
  fromId?: number;
  fromName?: string;
  toId?: number;
  roomId?: string;
  content?: string;
  messageId?: number;
  createdAt?: number; // 秒级时间戳
  chatId?: number; // 聊天ID（用于区分群聊和私聊）
  messageType?: number; // 消息类型（1=文本，2=图片等）
  // 任务进度相关
  taskId?: string;
  taskName?: string;
  progress?: number;
  status?: string;
  // 通知相关
  title?: string;
  level?: 'info' | 'success' | 'warning' | 'error';
}

interface WebSocketState {
  connected: boolean;
  connecting: boolean;
  reconnectAttempts: number;
  maxReconnectAttempts: number;
  reconnectDelay: number;
  ws: WebSocket | null;
  lastMessage: WSMessage | null;
  // 最近任务列表（用于浮动任务球）
  recentTasks: TaskItem[];
  recentTasksLoading: boolean;
  recentTasksUpdateTimer: number | null;
}

const RECONNECT_DELAY_BASE = 3000 // 基础重连延迟（毫秒）
const MAX_RECONNECT_ATTEMPTS = 10

export const useWebSocketStore = defineStore('websocket', {
  state: (): WebSocketState => ({
    connected: false,
    connecting: false,
    reconnectAttempts: 0,
    maxReconnectAttempts: MAX_RECONNECT_ATTEMPTS,
    reconnectDelay: RECONNECT_DELAY_BASE,
    ws: null,
    lastMessage: null,
    recentTasks: [],
    recentTasksLoading: false,
    recentTasksUpdateTimer: null
  }),

  actions: {
    // 连接 WebSocket
    async connect() {
      const userStore = useUserStore()

      // 在线聊天无需权限，只要登录就可以使用
      // 移除权限检查

      if (this.connecting || this.connected) {
        return
      }

      const token = userStore.token
      if (!token) {
        console.log('未登录，跳过 WebSocket 连接')
        return
      }

      this.connecting = true

      // 开发环境强制使用本地服务，忽略字典配置
      let wsBaseURL: string
      if (import.meta.env.DEV) {
        // 开发环境：直接连接本地后端服务
        wsBaseURL = 'localhost:20000'
      } else {
        // 生产环境：从字典获取配置，如果没有则使用默认值
        wsBaseURL = await getWebSocketBaseURL()
        if (!wsBaseURL) {
          wsBaseURL = 'oldbai.top'
        }
      }

      // 构建 WebSocket URL
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsPath = '/api/v1/chats/ws'
      let wsUrl: string

      // 处理 baseURL：
      // 开发环境：直接连接后端，路径为 ws://localhost:20000/api/v1/chats/ws
      // 生产环境：
      //   1. 如果 baseURL 包含 /ws，说明已经配置了完整的路径（如 oldbai.top/ws）
      //      那么直接拼接路径：ws://oldbai.top/ws/api/v1/chats/ws
      //   2. 如果 baseURL 不包含 /ws，说明只配置了域名（如 oldbai.top）
      //      那么直接连接后端：ws://oldbai.top/api/v1/chats/ws
      if (import.meta.env.DEV) {
        // 开发环境：直接连接本地后端，不添加 /ws 前缀
        wsUrl = `${protocol}//${wsBaseURL}${wsPath}?token=${encodeURIComponent(token)}&roomId=default`
      } else if (wsBaseURL.includes('/ws')) {
        // 生产环境：baseURL 已经包含 /ws，直接拼接
        wsUrl = `${protocol}//${wsBaseURL}${wsPath}?token=${encodeURIComponent(token)}&roomId=default`
      } else {
        // 生产环境：baseURL 不包含 /ws，直接连接后端（不添加 /ws 前缀）
        wsUrl = `${protocol}//${wsBaseURL}${wsPath}?token=${encodeURIComponent(token)}&roomId=default`
      }

      try {
        const ws = new WebSocket(wsUrl)

        ws.onopen = () => {
          this.connected = true
          this.connecting = false
          this.reconnectAttempts = 0
          this.reconnectDelay = RECONNECT_DELAY_BASE
        }

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data) as WSMessage
            this.handleMessage(data)
          } catch (err) {
            console.error('解析 WebSocket 消息失败:', err)
          }
        }

        ws.onerror = (error) => {
          console.error('WebSocket 错误:', error)
          console.error('WebSocket URL:', wsUrl.replace(/token=[^&]+/, 'token=***'))
          // 尝试输出更详细的错误信息
          if (ws.readyState === WebSocket.CLOSED) {
            console.error('WebSocket 连接已关闭，可能的原因：')
            console.error('1. 后端服务未运行')
            console.error('2. Nginx 配置错误（如果通过代理）')
            console.error('3. 路径不匹配')
            console.error('4. 防火墙或网络问题')
          }
          this.connecting = false
        }

        ws.onclose = () => {
          this.connected = false
          this.connecting = false
          this.ws = null
          console.log('WebSocket 连接已断开')

          // 自动重连
          if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++
            const delay = Math.min(
              this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1),
              30000
            )
            console.log(`将在 ${delay}ms 后尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
            setTimeout(() => {
              this.connect()
            }, delay)
          } else {
            console.error('达到最大重连次数，停止重连')
          }
        }

        this.ws = ws
      } catch (err) {
        console.error('创建 WebSocket 连接失败:', err)
        this.connecting = false
      }
    },

    // 断开连接
    disconnect() {
      if (this.ws) {
        this.ws.close()
        this.ws = null
      }
      this.connected = false
      this.connecting = false
      this.reconnectAttempts = 0
    },

    // 广播原始消息给外部订阅者（stores/notification.ts、composables/useChatList.ts 等通过 watch lastMessage 消费），
    // 顺带刷新任务浮球——这是本 store 自己持有的状态，不下放给消费方各自处理。
    handleMessage(data: WSMessage) {
      this.lastMessage = data

      if (data.type === MessageType.TASK_PROGRESS) {
        this.scheduleRefreshRecentTasks()
      } else if (data.type === MessageType.NOTIFICATION && data.taskId) {
        this.scheduleRefreshRecentTasks()
      }
    },

    // 调用后端接口刷新最近任务列表（带防抖）
    scheduleRefreshRecentTasks() {
      if (this.recentTasksUpdateTimer) {
        clearTimeout(this.recentTasksUpdateTimer)
      }
      this.recentTasksUpdateTimer = window.setTimeout(() => {
        this.refreshRecentTasks().catch((err) => {
          console.error('刷新最近任务列表失败:', err)
        })
      }, 500)
    },

    async refreshRecentTasks() {
      const userStore = useUserStore()
      if (!userStore.token) {
        return
      }
      this.recentTasksLoading = true
      try {
        const resp = await taskApi.taskRecent({})
        this.recentTasks = resp.list || []
      } finally {
        this.recentTasksLoading = false
      }
    },

    // 发送消息（如果需要）
    sendMessage(message: unknown) {
      if (this.ws && this.connected) {
        this.ws.send(JSON.stringify(message))
      } else {
        console.error('WebSocket 未连接，无法发送消息')
      }
    }
  }
})
