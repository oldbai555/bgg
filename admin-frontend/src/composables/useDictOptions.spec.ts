import {describe, it, expect, beforeEach} from 'vitest'
import {createPinia, setActivePinia} from 'pinia'
import {useDictOptions, getDictOptions} from './useDictOptions'
import {useDictStore} from '@/stores/dict'

describe('useDictOptions', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('字典 store 里没有数据时回退到 defaultValue', () => {
    const {options} = useDictOptions('not_loaded_code', [{label: '默认', value: 1}])
    expect(options.value).toEqual([{label: '默认', value: 1}])
  })

  it('字典 store 里有数据时优先使用 store 数据，而不是 defaultValue', () => {
    const dictStore = useDictStore()
    dictStore.dicts.status = [
      {id: 1, typeId: 1, label: '启用', value: '1', sort: 0, status: 1, remark: '', createdAt: 0},
      {id: 2, typeId: 1, label: '禁用', value: '2', sort: 1, status: 1, remark: '', createdAt: 0}
    ]

    const {options} = useDictOptions('status', [{label: '默认', value: 1}])
    expect(options.value).toEqual([
      {label: '启用', value: '1'},
      {label: '禁用', value: '2'}
    ])
  })

  it('getLabel 未命中时回退返回原始值的字符串形式', () => {
    const {getLabel} = useDictOptions('status')
    expect(getLabel(1)).toBe('1')
  })

  it('getLabel 命中字典项时返回对应 label', () => {
    const dictStore = useDictStore()
    dictStore.dicts.status = [
      {id: 1, typeId: 1, label: '启用', value: '1', sort: 0, status: 1, remark: '', createdAt: 0}
    ]
    const {getLabel} = useDictOptions('status')
    expect(getLabel(1)).toBe('启用')
    expect(getLabel('1')).toBe('启用')
  })

  it('非响应式的 getDictOptions 辅助函数行为与 useDictOptions 的 options 一致', () => {
    expect(getDictOptions('not_loaded_code', [{label: '默认', value: 1}])).toEqual([
      {label: '默认', value: 1}
    ])
  })
})
