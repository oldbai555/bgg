import {describe, it, expect, beforeEach, afterEach, vi} from 'vitest'
import {nextTick} from 'vue'
import {createPinia, setActivePinia} from 'pinia'
import {useNotificationStore} from './notification'
import {useWebSocketStore, MessageType} from './websocket'
import {useUserStore} from './user'

vi.mock('element-plus', () => ({
  ElMessage: Object.assign(vi.fn(), {
    info: vi.fn(),
    success: vi.fn(),
    warning: vi.fn(),
    error: vi.fn()
  })
}))

describe('useNotificationStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    history.pushState(null, '', '/bgg/admin/dashboard')
  })

  afterEach(() => {
    history.pushState(null, '', '/bgg/admin/dashboard')
  })

  it('addUnreadMessage 新消息插到最前面，超过 50 条按最新 50 条截断', () => {
    const store = useNotificationStore()
    for (let i = 0; i < 51; i++) {
      store.addUnreadMessage({
        id: `m${i}`,
        type: MessageType.SYSTEM,
        title: 't',
        content: 'c',
        timestamp: i,
        read: false
      })
    }

    expect(store.unreadMessages).toHaveLength(50)
    expect(store.unreadMessages[0].id).toBe('m50')
  })

  it('markAsRead / markAllAsRead / clearReadMessages 按预期变更已读状态', () => {
    const store = useNotificationStore()
    store.addUnreadMessage({id: 'a', type: MessageType.SYSTEM, title: 't', content: 'c', timestamp: 1, read: false})
    store.addUnreadMessage({id: 'b', type: MessageType.SYSTEM, title: 't', content: 'c', timestamp: 2, read: false})

    store.markAsRead('a')
    expect(store.unreadMessages.find((m) => m.id === 'a')?.read).toBe(true)
    expect(store.unreadCount).toBe(1)

    store.markAllAsRead()
    expect(store.unreadCount).toBe(0)

    store.clearReadMessages()
    expect(store.unreadMessages).toHaveLength(0)
  })

  it('unreadCount / hasUnreadChat 只统计未读消息', () => {
    const store = useNotificationStore()
    store.addUnreadMessage({id: 'chat1', type: MessageType.CHAT, title: 't', content: 'c', timestamp: 1, read: false})
    store.addUnreadMessage({id: 'sys1', type: MessageType.SYSTEM, title: 't', content: 'c', timestamp: 2, read: true})

    expect(store.unreadCount).toBe(1)
    expect(store.hasUnreadChat).toBe(true)

    store.markAsRead('chat1')
    expect(store.hasUnreadChat).toBe(false)
  })

  it('订阅 websocket store 的 lastMessage：不在聊天页时收到他人 CHAT 消息会计入未读', async () => {
    const store = useNotificationStore()
    const userStore = useUserStore()
    userStore.profile = {id: 1} as never
    const wsStore = useWebSocketStore()

    wsStore.handleMessage({
      type: MessageType.CHAT,
      fromId: 2,
      fromName: '张三',
      content: '你好',
      messageId: 100,
      chatId: 0
    })
    await nextTick()

    expect(store.unreadCount).toBe(1)
    expect(store.unreadMessages[0].type).toBe(MessageType.CHAT)
  })

  it('自己发出的 CHAT 消息不计入未读', async () => {
    const store = useNotificationStore()
    const userStore = useUserStore()
    userStore.profile = {id: 1} as never
    const wsStore = useWebSocketStore()

    wsStore.handleMessage({type: MessageType.CHAT, fromId: 1, content: '我发的'})
    await nextTick()

    expect(store.unreadCount).toBe(0)
  })

  it('在聊天页面时收到 CHAT 消息不计入未读', async () => {
    const store = useNotificationStore()
    const userStore = useUserStore()
    userStore.profile = {id: 1} as never
    history.pushState(null, '', '/bgg/admin/chatroom/chat')
    const wsStore = useWebSocketStore()

    wsStore.handleMessage({type: MessageType.CHAT, fromId: 2, fromName: '张三', content: '你好'})
    await nextTick()

    expect(store.unreadCount).toBe(0)
  })

  it('在聊天记录管理/群组管理页面（同前缀但不是聊天页）收到 CHAT 消息仍计入未读', async () => {
    const store = useNotificationStore()
    const userStore = useUserStore()
    userStore.profile = {id: 1} as never
    history.pushState(null, '', '/bgg/admin/chatroom/chat-message')
    const wsStore = useWebSocketStore()

    wsStore.handleMessage({type: MessageType.CHAT, fromId: 2, fromName: '张三', content: '你好'})
    await nextTick()

    expect(store.unreadCount).toBe(1)
  })

  it('带 taskName 的 TASK_PROGRESS 消息计入未读，不带 taskName 的不计入', async () => {
    const store = useNotificationStore()
    const wsStore = useWebSocketStore()

    wsStore.handleMessage({type: MessageType.TASK_PROGRESS, taskId: 't1', progress: 50})
    await nextTick()
    expect(store.unreadCount).toBe(0)

    wsStore.handleMessage({type: MessageType.TASK_PROGRESS, taskId: 't1', taskName: '导出报表', progress: 80})
    await nextTick()
    expect(store.unreadCount).toBe(1)
  })

  it('NOTIFICATION 消息计入未读', async () => {
    const store = useNotificationStore()
    const wsStore = useWebSocketStore()

    wsStore.handleMessage({type: MessageType.NOTIFICATION, title: '系统通知', content: '维护中', level: 'warning'})
    await nextTick()

    expect(store.unreadCount).toBe(1)
    expect(store.unreadMessages[0].title).toBe('系统通知')
  })

  it('SYSTEM 消息不计入未读列表', async () => {
    const store = useNotificationStore()
    const wsStore = useWebSocketStore()

    wsStore.handleMessage({type: MessageType.SYSTEM, content: '心跳'})
    await nextTick()

    expect(store.unreadCount).toBe(0)
  })
})
