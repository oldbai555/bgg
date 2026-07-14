import {describe, it, expect, beforeEach, afterEach, vi} from 'vitest'
import {createPinia, setActivePinia} from 'pinia'
import {useWebSocketStore, MessageType} from './websocket'
import {useUserStore} from './user'
import {taskApi} from '@/api/task'

vi.mock('@/api/task', () => ({
  taskApi: {
    taskRecent: vi.fn()
  }
}))

describe('useWebSocketStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('未登录（无 token）时 connect() 直接跳过，不进入 connecting 状态', async () => {
    const store = useWebSocketStore()
    await store.connect()

    expect(store.connecting).toBe(false)
    expect(store.connected).toBe(false)
    expect(store.ws).toBeNull()
  })

  it('已在连接中或已连接时 connect() 是幂等的，不重复发起', async () => {
    const store = useWebSocketStore()
    useUserStore().token = 'test-token'
    store.connecting = true

    await store.connect()

    expect(store.ws).toBeNull()
  })

  it('disconnect() 关闭连接并复位状态', () => {
    const store = useWebSocketStore()
    const close = vi.fn()
    store.ws = {close} as unknown as WebSocket
    store.connected = true
    store.connecting = true
    store.reconnectAttempts = 3

    store.disconnect()

    expect(close).toHaveBeenCalled()
    expect(store.ws).toBeNull()
    expect(store.connected).toBe(false)
    expect(store.connecting).toBe(false)
    expect(store.reconnectAttempts).toBe(0)
  })

  it('sendMessage() 未连接时不发送，只报错', () => {
    const store = useWebSocketStore()
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    store.sendMessage({type: 'chat', content: 'hi'})

    expect(errorSpy).toHaveBeenCalled()
    errorSpy.mockRestore()
  })

  it('sendMessage() 已连接时通过 ws.send 发送 JSON', () => {
    const store = useWebSocketStore()
    const send = vi.fn()
    store.ws = {send} as unknown as WebSocket
    store.connected = true

    store.sendMessage({type: 'chat', content: 'hi'})

    expect(send).toHaveBeenCalledWith(JSON.stringify({type: 'chat', content: 'hi'}))
  })

  it('handleMessage() 总是广播 lastMessage，供其它 store 订阅消费', () => {
    const store = useWebSocketStore()
    store.handleMessage({type: MessageType.SYSTEM, content: 'ping'})

    expect(store.lastMessage).toEqual({type: MessageType.SYSTEM, content: 'ping'})
  })

  describe('任务浮球刷新（防抖）', () => {
    beforeEach(() => {
      vi.useFakeTimers()
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('收到 TASK_PROGRESS 消息后延迟刷新最近任务列表', async () => {
      useUserStore().token = 'test-token'
      vi.mocked(taskApi.taskRecent).mockResolvedValue({list: [], total: 0} as never)
      const store = useWebSocketStore()

      store.handleMessage({type: MessageType.TASK_PROGRESS, taskId: 't1', taskName: '导出'})
      await vi.advanceTimersByTimeAsync(500)

      expect(taskApi.taskRecent).toHaveBeenCalled()
    })

    it('收到带 taskId 的 NOTIFICATION 消息同样触发刷新', async () => {
      useUserStore().token = 'test-token'
      vi.mocked(taskApi.taskRecent).mockResolvedValue({list: [], total: 0} as never)
      const store = useWebSocketStore()

      store.handleMessage({type: MessageType.NOTIFICATION, taskId: 't1', content: '完成'})
      await vi.advanceTimersByTimeAsync(500)

      expect(taskApi.taskRecent).toHaveBeenCalled()
    })

    it('普通 NOTIFICATION（无 taskId）不触发任务刷新', async () => {
      useUserStore().token = 'test-token'
      vi.mocked(taskApi.taskRecent).mockResolvedValue({list: [], total: 0} as never)
      const store = useWebSocketStore()

      store.handleMessage({type: MessageType.NOTIFICATION, content: '普通通知'})
      await vi.advanceTimersByTimeAsync(500)

      expect(taskApi.taskRecent).not.toHaveBeenCalled()
    })
  })
})
