import {defineStore} from 'pinia';
import {useUserStore} from './user';
import {usePermission} from '@/hooks/usePermission';
import {ElMessage} from 'element-plus';
import {getWebSocketBaseURL} from '@/composables/useAppConfig';

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

// 未读消息
export interface UnreadMessage {
  id: string;
  type: MessageType | string;
  title: string;
  content: string;
  timestamp: number;
  read: boolean;
}

interface WebSocketState {
  connected: boolean;
  connecting: boolean;
  reconnectAttempts: number;
  maxReconnectAttempts: number;
  reconnectDelay: number;
  ws: WebSocket | null;
  unreadMessages: UnreadMessage[];
  lastMessage: WSMessage | null;
}

const RECONNECT_DELAY_BASE = 3000; // 基础重连延迟（毫秒）
const MAX_RECONNECT_ATTEMPTS = 10;

export const useWebSocketStore = defineStore('websocket', {
  state: (): WebSocketState => ({
    connected: false,
    connecting: false,
    reconnectAttempts: 0,
    maxReconnectAttempts: MAX_RECONNECT_ATTEMPTS,
    reconnectDelay: RECONNECT_DELAY_BASE,
    ws: null,
    unreadMessages: [],
    lastMessage: null
  }),

  getters: {
    unreadCount: (state) => state.unreadMessages.filter((m) => !m.read).length,
    hasUnreadChat: (state) =>
      state.unreadMessages.some((m) => !m.read && m.type === MessageType.CHAT)
  },

  actions: {
    // 连接 WebSocket
    async connect() {
      const userStore = useUserStore();
      const {hasPermission} = usePermission();

      // 在线聊天无需权限，只要登录就可以使用
      // 移除权限检查

      if (this.connecting || this.connected) {
        return;
      }

      const token = userStore.token;
      if (!token) {
        console.log('未登录，跳过 WebSocket 连接');
        return;
      }

      this.connecting = true;

      // 开发环境强制使用本地服务，忽略字典配置
      let wsBaseURL: string;
      if (import.meta.env.DEV) {
        // 开发环境：直接连接本地后端服务
        wsBaseURL = 'localhost:20000';
      } else {
        // 生产环境：从字典获取配置，如果没有则使用默认值
        wsBaseURL = await getWebSocketBaseURL();
        if (!wsBaseURL) {
          wsBaseURL = 'oldbai.top';
        }
      }

      // 构建 WebSocket URL
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsPath = '/api/v1/chats/ws';
      let wsUrl: string;
      
      // 处理 baseURL：
      // 开发环境：直接连接后端，路径为 ws://localhost:20000/api/v1/chats/ws
      // 生产环境：
      //   1. 如果 baseURL 包含 /ws，说明已经配置了完整的路径（如 oldbai.top/ws）
      //      那么直接拼接路径：ws://oldbai.top/ws/api/v1/chats/ws
      //   2. 如果 baseURL 不包含 /ws，说明只配置了域名（如 oldbai.top）
      //      那么直接连接后端：ws://oldbai.top/api/v1/chats/ws
      if (import.meta.env.DEV) {
        // 开发环境：直接连接本地后端，不添加 /ws 前缀
        wsUrl = `${protocol}//${wsBaseURL}${wsPath}?token=${encodeURIComponent(token)}&roomId=default`;
      } else if (wsBaseURL.includes('/ws')) {
        // 生产环境：baseURL 已经包含 /ws，直接拼接
        wsUrl = `${protocol}//${wsBaseURL}${wsPath}?token=${encodeURIComponent(token)}&roomId=default`;
      } else {
        // 生产环境：baseURL 不包含 /ws，直接连接后端（不添加 /ws 前缀）
        wsUrl = `${protocol}//${wsBaseURL}${wsPath}?token=${encodeURIComponent(token)}&roomId=default`;
      }

      try {
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
          this.connected = true;
          this.connecting = false;
          this.reconnectAttempts = 0;
          this.reconnectDelay = RECONNECT_DELAY_BASE;
        };

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data) as WSMessage;
            this.handleMessage(data);
          } catch (err) {
            console.error('解析 WebSocket 消息失败:', err);
          }
        };

        ws.onerror = (error) => {
          console.error('WebSocket 错误:', error);
          console.error('WebSocket URL:', wsUrl.replace(/token=[^&]+/, 'token=***'));
          // 尝试输出更详细的错误信息
          if (ws.readyState === WebSocket.CLOSED) {
            console.error('WebSocket 连接已关闭，可能的原因：');
            console.error('1. 后端服务未运行');
            console.error('2. Nginx 配置错误（如果通过代理）');
            console.error('3. 路径不匹配');
            console.error('4. 防火墙或网络问题');
          }
          this.connecting = false;
        };

        ws.onclose = () => {
          this.connected = false;
          this.connecting = false;
          this.ws = null;
          console.log('WebSocket 连接已断开');

          // 自动重连
          if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = Math.min(
              this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1),
              30000
            );
            console.log(`将在 ${delay}ms 后尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
            setTimeout(() => {
              this.connect();
            }, delay);
          } else {
            console.error('达到最大重连次数，停止重连');
          }
        };

        this.ws = ws;
      } catch (err) {
        console.error('创建 WebSocket 连接失败:', err);
        this.connecting = false;
      }
    },

    // 断开连接
    disconnect() {
      if (this.ws) {
        this.ws.close();
        this.ws = null;
      }
      this.connected = false;
      this.connecting = false;
      this.reconnectAttempts = 0;
    },

    // 处理接收到的消息
    handleMessage(data: WSMessage) {
      this.lastMessage = data;

      // 根据消息类型处理
      switch (data.type) {
        case MessageType.CHAT:
          this.handleChatMessage(data);
          break;
        case MessageType.TASK_PROGRESS:
          this.handleTaskProgress(data);
          break;
        case MessageType.NOTIFICATION:
          this.handleNotification(data);
          break;
        case MessageType.SYSTEM:
          this.handleSystemMessage(data);
          break;
        default:
          // 未知消息类型，忽略
          break;
      }
    },

    // 处理聊天消息
    handleChatMessage(data: WSMessage) {
      // 获取当前用户ID
      const userStore = useUserStore();
      const currentUserId = userStore.profile?.id || 0;
      
      // 只有不是自己发的消息才需要显示未读通知
      const isMyMessage = data.fromId && Number(data.fromId) === Number(currentUserId);
      if (isMyMessage) {
        // 自己发的消息不需要添加到未读消息
        return;
      }

      // 检查当前是否在聊天页面（支持多种路径格式）
      // 注意：使用 hash 路由时，路径在 hash 中
      const currentPath = window.location.pathname;
      const currentHash = window.location.hash;
      // 检查是否在 ChatList.vue 页面（/chatroom/chat）
      const isInChatListPage = currentHash === '#/chatroom/chat' || 
                               currentHash.startsWith('#/chatroom/chat?') ||
                               currentPath.includes('/chatroom/chat');
      // 检查是否在其他聊天相关页面
      const isInOtherChatPage = currentHash.includes('/temp/chat') || 
                                currentHash.includes('/chat') ||
                                currentPath.includes('/temp/chat') || 
                                currentPath.includes('/chat');
      const isInChatPage = isInChatListPage || isInOtherChatPage;

      // 如果不在聊天页面，添加到未读消息并显示提示
      if (!isInChatPage) {
        // 构建消息标题和内容
        const chatId = data.chatId || 0;
        const isGroupChat = chatId > 0; // 如果有chatId且大于0，可能是群聊
        const title = isGroupChat 
          ? `群聊消息：来自 ${data.fromName || '未知用户'}`
          : `来自 ${data.fromName || '未知用户'}`;
        
        const content = data.content || '';
        // 如果是图片消息，显示特殊提示
        const displayContent = data.messageType === 2 ? '[图片]' : content;

        // 添加到未读消息
        this.addUnreadMessage({
          id: `chat_${data.messageId || Date.now()}_${chatId}`,
          type: MessageType.CHAT,
          title: title,
          content: displayContent,
          timestamp: Date.now(),
          read: false
        });

        // 显示提示消息
        ElMessage.info({
          message: `${title}: ${displayContent}`,
          duration: 3000,
          showClose: true
        });
      }
    },

    // 处理任务进度
    handleTaskProgress(data: WSMessage) {
      // 可以在这里处理任务进度更新
      console.log('任务进度更新:', data);
      // 如果需要，也可以添加到未读消息
      if (data.taskName) {
        this.addUnreadMessage({
          id: `task_${data.taskId || Date.now()}`,
          type: MessageType.TASK_PROGRESS,
          title: `任务进度: ${data.taskName}`,
          content: `进度: ${data.progress || 0}% - ${data.status || ''}`,
          timestamp: Date.now(),
          read: false
        });
      }
    },

    // 处理通知消息
    handleNotification(data: WSMessage) {
      const level = data.level || 'info';
      ElMessage[level](data.content || data.title || '新通知');

      this.addUnreadMessage({
        id: `notify_${Date.now()}`,
        type: MessageType.NOTIFICATION,
        title: data.title || '通知',
        content: data.content || '',
        timestamp: Date.now(),
        read: false
      });
    },

    // 处理系统消息
    handleSystemMessage(data: WSMessage) {
      console.log('系统消息:', data);
      // 系统消息通常不需要添加到未读消息列表
    },

    // 添加未读消息
    addUnreadMessage(message: UnreadMessage) {
      this.unreadMessages.unshift(message);
      // 限制未读消息数量，最多保留 50 条
      if (this.unreadMessages.length > 50) {
        this.unreadMessages = this.unreadMessages.slice(0, 50);
      }
    },

    // 标记消息为已读
    markAsRead(messageId: string) {
      const message = this.unreadMessages.find((m) => m.id === messageId);
      if (message) {
        message.read = true;
      }
    },

    // 标记所有消息为已读
    markAllAsRead() {
      this.unreadMessages.forEach((m) => {
        m.read = true;
      });
    },

    // 清除已读消息
    clearReadMessages() {
      this.unreadMessages = this.unreadMessages.filter((m) => !m.read);
    },

    // 发送消息（如果需要）
    sendMessage(message: any) {
      if (this.ws && this.connected) {
        this.ws.send(JSON.stringify(message));
      } else {
        console.error('WebSocket 未连接，无法发送消息');
      }
    }
  }
});

