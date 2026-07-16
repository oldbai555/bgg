import {describe, it, expect, beforeEach, vi} from 'vitest'
import {nextTick} from 'vue'
import {createPinia, setActivePinia} from 'pinia'
import {useChatList} from './useChatList'
import {useUserStore} from '@/stores/user'
import {useWebSocketStore, MessageType} from '@/stores/websocket'
import {chatApi} from '@/api/chat'
import {systemApi} from '@/api/system'

vi.mock('@/api/chat', () => ({
  chatApi: {
    chatList: vi.fn(),
    chatMessageList: vi.fn(),
    chatMessageSend: vi.fn()
  }
}))

vi.mock('@/api/system', () => ({
  systemApi: {
    dictGet: vi.fn()
  }
}))

vi.mock('element-plus', () => ({
  ElMessage: Object.assign(vi.fn(), {
    info: vi.fn(),
    success: vi.fn(),
    warning: vi.fn(),
    error: vi.fn()
  })
}))

describe('useChatList', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    useUserStore().profile = {id: 1, username: 'alice'} as never
  })

  describe('loadChatConfig', () => {
    it('从 chat_config 字典读取消息数量/emoji 行列配置', async () => {
      vi.mocked(systemApi.dictGet).mockResolvedValue({
        items: [
          {label: '聊天窗口消息数量', value: '50'},
          {label: 'Emoji每行显示数量', value: '6'},
          {label: 'Emoji显示行数', value: '4'}
        ]
      } as never)

      const {loadChatConfig, emojiColsPerRow, emojiRows} = useChatList()
      await loadChatConfig()

      expect(emojiColsPerRow.value).toBe(6)
      expect(emojiRows.value).toBe(4)
    })

    it('字典请求失败时回退默认值 30，不向上抛错', async () => {
      vi.mocked(systemApi.dictGet).mockRejectedValue(new Error('network error'))

      const {loadChatConfig} = useChatList()
      await expect(loadChatConfig()).resolves.toBeUndefined()
    })
  })

  describe('loadChats / selectChat', () => {
    it('加载聊天列表后，未选中任何会话时自动选中第一个并加载其消息', async () => {
      vi.mocked(chatApi.chatList).mockResolvedValue({
        list: [{chatId: 10, name: '群聊A'}, {chatId: 20, name: '群聊B'}]
      } as never)
      vi.mocked(chatApi.chatMessageList).mockResolvedValue({list: []} as never)

      const {loadChats, selectedChatId, chats} = useChatList()
      await loadChats()
      await nextTick()

      expect(chats.value).toHaveLength(2)
      expect(selectedChatId.value).toBe(10)
      expect(chatApi.chatMessageList).toHaveBeenCalledWith(
        expect.objectContaining({chatId: 10})
      )
    })

    it('selectChat 切换会话后按新 chatId 重新加载消息，最新消息在底部', async () => {
      vi.mocked(chatApi.chatMessageList).mockResolvedValue({
        list: [
          {id: 3, chatId: 20, content: 'c3', createdAt: 300},
          {id: 2, chatId: 20, content: 'c2', createdAt: 200},
          {id: 1, chatId: 20, content: 'c1', createdAt: 100}
        ]
      } as never)

      const {selectChat, messages, selectedChatId} = useChatList()
      selectChat({chatId: 20, name: '群聊B'} as never)
      await nextTick()

      expect(selectedChatId.value).toBe(20)
      expect(messages.value.map((m) => m.id)).toEqual([1, 2, 3])
    })
  })

  describe('sendTextMessage / sendImageMessage', () => {
    it('未选中会话时发送文本消息抛错，不调用接口', async () => {
      const {sendTextMessage} = useChatList()
      await expect(sendTextMessage('hi')).rejects.toThrow('请先选择一个聊天')
      expect(chatApi.chatMessageSend).not.toHaveBeenCalled()
    })

    it('发送文本消息成功后立即在本地追加一条消息', async () => {
      vi.mocked(chatApi.chatMessageList).mockResolvedValue({list: []} as never)
      vi.mocked(chatApi.chatMessageSend).mockResolvedValue({} as never)

      const {selectChat, sendTextMessage, messages} = useChatList()
      selectChat({chatId: 20, name: '群聊B'} as never)
      await nextTick()

      await sendTextMessage('你好')

      expect(chatApi.chatMessageSend).toHaveBeenCalledWith(
        expect.objectContaining({chatId: 20, content: '你好', messageType: 1})
      )
      expect(messages.value).toHaveLength(1)
      expect(messages.value[0].content).toBe('你好')
      expect(messages.value[0].fromUserId).toBe(1)
    })

    it('未选中会话时发送图片消息抛错，不调用接口', async () => {
      const {sendImageMessage} = useChatList()
      await expect(sendImageMessage('http://x/a.png')).rejects.toThrow('请先选择一个聊天')
      expect(chatApi.chatMessageSend).not.toHaveBeenCalled()
    })

    it('发送图片消息成功后本地追加一条 messageType=2 的消息', async () => {
      vi.mocked(chatApi.chatMessageList).mockResolvedValue({list: []} as never)
      vi.mocked(chatApi.chatMessageSend).mockResolvedValue({} as never)

      const {selectChat, sendImageMessage, messages} = useChatList()
      selectChat({chatId: 20, name: '群聊B'} as never)
      await nextTick()

      await sendImageMessage('http://x/a.png')

      expect(messages.value).toHaveLength(1)
      expect(messages.value[0].messageType).toBe(2)
      expect(messages.value[0].content).toBe('http://x/a.png')
    })
  })

  describe('WebSocket 消息同步（订阅 wsStore.lastMessage）', () => {
    it('收到当前选中会话的 CHAT 消息会追加到 messages', async () => {
      vi.mocked(chatApi.chatMessageList).mockResolvedValue({list: []} as never)
      const {selectChat, messages} = useChatList()
      selectChat({chatId: 20, name: '群聊B'} as never)
      await nextTick()

      const wsStore = useWebSocketStore()
      wsStore.handleMessage({
        type: MessageType.CHAT,
        chatId: 20,
        fromId: 2,
        fromName: '张三',
        content: '你好',
        messageId: 999,
        createdAt: Math.floor(Date.now() / 1000)
      })
      await nextTick()

      expect(messages.value).toHaveLength(1)
      expect(messages.value[0].id).toBe(999)
    })

    it('收到非当前选中会话的 CHAT 消息会被忽略', async () => {
      vi.mocked(chatApi.chatMessageList).mockResolvedValue({list: []} as never)
      const {selectChat, messages} = useChatList()
      selectChat({chatId: 20, name: '群聊B'} as never)
      await nextTick()

      const wsStore = useWebSocketStore()
      wsStore.handleMessage({type: MessageType.CHAT, chatId: 999, fromId: 2, content: '别的群'})
      await nextTick()

      expect(messages.value).toHaveLength(0)
    })
  })
})
