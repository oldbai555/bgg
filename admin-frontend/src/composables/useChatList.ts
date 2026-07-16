/**
 * 会话列表页数据/业务逻辑 Composable
 * 从 views/chat/ChatList.vue 拆分而来：负责聊天列表/消息的加载、发送、
 * WebSocket 消息同步，不涉及任何 DOM 操作（滚动等交由调用方处理）
 */
import {ref, reactive, computed, nextTick, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {useUserStore} from '@/stores/user'
import {useWebSocketStore, MessageType} from '@/stores/websocket'
import {chatApi} from '@/api/chat'
import {systemApi} from '@/api/system'
import type {
  ChatMessageItem,
  ChatMessageListReq,
  ChatMessageSendReq,
  ChatItem
} from '@/api/generated/admin'

interface UseChatListOptions {
  // 消息列表发生变化（新增/加载完成）时回调，用于调用方做滚动到底部等 DOM 操作
  onMessagesChanged?: () => void
}

export function useChatList(options: UseChatListOptions = {}) {
  const userStore = useUserStore()
  const wsStore = useWebSocketStore()

  const currentUserId = computed(() => userStore.profile?.id || 0)
  const currentUsername = computed(() => userStore.profile?.username || '')
  const wsConnected = computed(() => wsStore.connected)

  const selectedChatId = ref<number | null>(null)
  const selectedChat = ref<ChatItem | null>(null)
  const messages = ref<ChatMessageItem[]>([])
  const chats = ref<ChatItem[]>([])

  // 聊天配置：消息数量限制（从字典获取，默认30）
  const chatMessageLimit = ref(30)
  // Emoji 分页配置（从字典获取，默认 8 列 x 3 行）
  const emojiColsPerRow = ref(8)
  const emojiRows = ref(3)

  // 查询参数
  const query = reactive<ChatMessageListReq>({
    page: 1,
    pageSize: 30, // 默认30，将从字典加载后更新
    chatId: 0
  })

  const notifyMessagesChanged = () => {
    options.onMessagesChanged?.()
  }

  // 处理来自 WebSocket 的聊天消息
  const handleChatMessage = (data: Record<string, unknown>) => {
    // 检查是否是当前选中的聊天
    if (!(selectedChatId.value && data.chatId && Number(data.chatId) === Number(selectedChatId.value))) {
      return
    }

    // 智能判断消息类型：如果 content 是完整的 URL 且看起来是图片，则可能是图片消息
    let messageType = Number(data.messageType) || 1
    if (!data.messageType && data.content) {
      const content = String(data.content)
      if ((content.startsWith('http://') || content.startsWith('https://')) &&
          (content.includes('/uploads/') || content.includes('/files/') ||
           /\.(jpg|jpeg|png|gif|webp|bmp|svg)(\?|$)/i.test(content))) {
        messageType = 2 // 图片消息
      }
    }

    // 检查是否已存在相同的消息（避免重复添加）：
    // 优先通过 messageId 匹配；否则通过 内容+发送者+时间戳（10秒内）匹配
    const messageId = data.messageId && Number(data.messageId) > 0 ? Number(data.messageId) : null
    const content = String(data.content || '')
    const fromId = Number(data.fromId || 0)
    const wsCreatedAt = Number(data.createdAt || Math.floor(Date.now() / 1000))

    const existingIndex = messages.value.findIndex(msg => {
      const msgContent = String(msg.content || '')
      const msgFromId = Number(msg.fromUserId || 0)
      const msgCreatedAt = Number(msg.createdAt || 0)

      if (msgContent !== content || msgFromId !== fromId) {
        return false
      }

      const timeDiff = Math.abs(msgCreatedAt - wsCreatedAt)
      const timeThreshold = 10

      if (timeDiff <= timeThreshold) {
        if (messageId) {
          const msgIdNum = Number(msg.id)
          if (msgIdNum && msgIdNum === messageId) {
            return true
          }
        }
        return true
      }

      return false
    })

    if (existingIndex >= 0) {
      // 如果已存在，更新消息（使用服务器返回的ID和类型），不滚动到底部（避免干扰用户）
      const finalMessageId = messageId || messages.value[existingIndex].id
      messages.value[existingIndex] = {
        ...messages.value[existingIndex],
        id: finalMessageId,
        messageType: messageType,
        createdAt: wsCreatedAt
      }
    } else {
      // 收到新消息，添加到消息列表
      const finalMessageId = messageId || Date.now()
      const newMessage: ChatMessageItem = {
        id: finalMessageId,
        chatId: Number(data.chatId) || 0,
        fromUserId: Number(data.fromId) || 0,
        fromUserName: String(data.fromName || ''),
        content: content,
        messageType: messageType,
        status: 1,
        createdAt: wsCreatedAt
      }
      messages.value.push(newMessage)
      if (messages.value.length > chatMessageLimit.value) {
        messages.value = messages.value.slice(-chatMessageLimit.value)
      }
      notifyMessagesChanged()
    }
  }

  // 监听 WebSocket 消息（使用全局 store）
  watch(
    () => wsStore.lastMessage,
    (newMessage) => {
      if (!newMessage) {
        return
      }

      // 只处理聊天相关的消息
      if (newMessage.type === MessageType.CHAT || newMessage.type === 'chat') {
        handleChatMessage(newMessage)
      } else if (newMessage.type === 'join') {
        ElMessage.info(`${newMessage.fromName} 加入了聊天室`)
      } else if (newMessage.type === 'leave') {
        ElMessage.info(`${newMessage.fromName} 离开了聊天室`)
      }
    }
  )

  // 加载聊天配置（从字典获取）
  const loadChatConfig = async () => {
    try {
      const resp = await systemApi.dictGet({code: 'chat_config'})
      if (resp && resp.items && resp.items.length > 0) {
        // 查找"聊天窗口消息数量"配置项
        const limitItem = resp.items.find(item => item.label === '聊天窗口消息数量')
        if (limitItem && limitItem.value) {
          const limit = parseInt(limitItem.value, 10)
          if (!isNaN(limit) && limit > 0) {
            chatMessageLimit.value = limit
            query.pageSize = limit
          }
        }

        const colsItem = resp.items.find(item => item.label === 'Emoji每行显示数量')
        if (colsItem && colsItem.value) {
          const cols = parseInt(colsItem.value, 10)
          if (!isNaN(cols) && cols > 0) {
            emojiColsPerRow.value = cols
          }
        }

        const rowsItem = resp.items.find(item => item.label === 'Emoji显示行数')
        if (rowsItem && rowsItem.value) {
          const rows = parseInt(rowsItem.value, 10)
          if (!isNaN(rows) && rows > 0) {
            emojiRows.value = rows
          }
        }
      }
    } catch (err: unknown) {
      console.warn('加载聊天配置失败，使用默认值:', err)
      chatMessageLimit.value = 30
      query.pageSize = 30
    }
  }

  // 加载消息列表
  const loadMessages = async () => {
    if (!selectedChatId.value) {
      messages.value = []
      return
    }

    try {
      query.page = 1
      query.pageSize = chatMessageLimit.value // 使用从字典获取的限制值
      query.chatId = selectedChatId.value

      const resp = await chatApi.chatMessageList(query)
      const allMessages = (resp.list || []).reverse() // 反转列表，最新的在底部
      // 只保留最新的N条消息（N为字典配置的值）
      messages.value = allMessages.slice(-chatMessageLimit.value)
      nextTick(() => {
        notifyMessagesChanged()
      })
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '加载消息失败'
      ElMessage.error(message)
    }
  }

  // 加载聊天列表
  const loadChats = async () => {
    try {
      const resp = await chatApi.chatList()
      chats.value = resp.list || []
      // 如果没有选中的聊天，默认选中第一个
      if (chats.value.length > 0 && !selectedChatId.value) {
        selectChat(chats.value[0])
      }
    } catch (err: unknown) {
      console.error('加载聊天列表失败:', err)
    }
  }

  // 选择聊天
  const selectChat = (chat: ChatItem) => {
    selectedChatId.value = chat.chatId
    selectedChat.value = chat
    loadMessages()
  }

  // 发送文本消息（调用方需保证已选中聊天，否则抛错交由调用方的 catch 处理）
  const sendTextMessage = async (text: string) => {
    if (!selectedChatId.value) {
      throw new Error('请先选择一个聊天')
    }
    const chatId = selectedChatId.value
    const req: ChatMessageSendReq = {
      chatId,
      content: text,
      messageType: 1 // 文本消息
    }
    await chatApi.chatMessageSend(req)

    // 检查是否已有相同的消息（WebSocket 可能已经推送了）
    const localCreatedAt = Math.floor(Date.now() / 1000)
    const existingIndex = messages.value.findIndex(msg => {
      const msgContent = String(msg.content || '')
      const msgFromId = Number(msg.fromUserId || 0)
      const msgCreatedAt = Number(msg.createdAt || 0)

      if (msgContent !== text || msgFromId !== Number(currentUserId.value)) {
        return false
      }

      // 时间戳差异在 5 秒内（WebSocket 消息可能先到达）
      const timeDiff = Math.abs(msgCreatedAt - localCreatedAt)
      return timeDiff <= 5
    })

    if (existingIndex < 0) {
      // 立即在本地添加消息（使用临时ID，WebSocket返回后会更新或合并）
      const tempId = Date.now()
      const localMessage: ChatMessageItem = {
        id: tempId,
        chatId,
        fromUserId: Number(currentUserId.value),
        fromUserName: currentUsername.value,
        content: text,
        messageType: 1,
        status: 1,
        createdAt: localCreatedAt
      }
      messages.value.push(localMessage)
    }
    if (messages.value.length > chatMessageLimit.value) {
      messages.value = messages.value.slice(-chatMessageLimit.value)
    }
    notifyMessagesChanged()
    // 注意：消息也会通过 WebSocket 推送回来，但我们已经提前显示了，避免重复
  }

  // 发送图片消息（调用方需保证已选中聊天，否则抛错交由调用方的 catch 处理）
  const sendImageMessage = async (imageUrl: string) => {
    if (!selectedChatId.value) {
      throw new Error('请先选择一个聊天')
    }
    const chatId = selectedChatId.value
    const req: ChatMessageSendReq = {
      chatId,
      content: imageUrl,
      messageType: 2 // 图片消息
    }
    await chatApi.chatMessageSend(req)

    // 立即在本地添加消息，确保图片能正确显示
    const localMessage: ChatMessageItem = {
      id: Date.now(), // 临时ID，WebSocket返回后会更新
      chatId,
      fromUserId: Number(currentUserId.value),
      fromUserName: currentUsername.value,
      content: imageUrl,
      messageType: 2,
      status: 1,
      createdAt: Math.floor(Date.now() / 1000)
    }
    messages.value.push(localMessage)
    if (messages.value.length > chatMessageLimit.value) {
      messages.value = messages.value.slice(-chatMessageLimit.value)
    }
    notifyMessagesChanged()
  }

  return {
    currentUserId,
    currentUsername,
    wsConnected,
    chats,
    selectedChatId,
    selectedChat,
    messages,
    emojiColsPerRow,
    emojiRows,
    loadChatConfig,
    loadChats,
    selectChat,
    sendTextMessage,
    sendImageMessage
  }
}
