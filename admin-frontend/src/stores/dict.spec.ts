import {describe, it, expect, beforeEach, vi} from 'vitest'
import {createPinia, setActivePinia} from 'pinia'
import {useDictStore} from './dict'
import {systemApi} from '@/api/system'

vi.mock('@/api/system', () => ({
  systemApi: {
    dictBatchGet: vi.fn()
  }
}))

const item = (label: string, value: string) => ({
  id: 1,
  typeId: 1,
  label,
  value,
  sort: 0,
  status: 1,
  remark: '',
  createdAt: 0
})

describe('useDictStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loadDicts 成功后写入 dicts 并标记 loaded', async () => {
    vi.mocked(systemApi.dictBatchGet).mockResolvedValue({
      dicts: {
        notice_type: {items: [item('普通公告', '1')]}
      }
    } as never)

    const store = useDictStore()
    await store.loadDicts(['notice_type'])

    expect(store.loaded).toBe(true)
    expect(store.getDictItems('notice_type')).toHaveLength(1)
    expect(store.getDictLabel('notice_type', '1')).toBe('普通公告')
  })

  it('getDictLabel 未命中时回退返回值本身的字符串形式', () => {
    const store = useDictStore()
    expect(store.getDictLabel('unknown_code', 5)).toBe('5')
  })

  it('getDictOptions 把字典项映射成 {label, value} 选项', () => {
    const store = useDictStore()
    store.dicts.status = [item('启用', '1'), item('禁用', '2')]
    expect(store.getDictOptions('status')).toEqual([
      {label: '启用', value: '1'},
      {label: '禁用', value: '2'}
    ])
  })

  it('loadDicts 失败时向上抛出错误，不吞掉异常', async () => {
    vi.mocked(systemApi.dictBatchGet).mockRejectedValue(new Error('network error'))
    const store = useDictStore()
    await expect(store.loadDicts(['notice_type'])).rejects.toThrow('network error')
  })

  it('clearDicts 清空字典数据和 loaded 标记', () => {
    const store = useDictStore()
    store.dicts.status = [item('启用', '1')]
    store.loaded = true
    store.loadAt = Date.now()

    store.clearDicts()

    expect(store.dicts).toEqual({})
    expect(store.loaded).toBe(false)
    expect(store.loadAt).toBe(0)
  })
})
