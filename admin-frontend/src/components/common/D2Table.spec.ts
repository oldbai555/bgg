import {describe, it, expect, beforeEach} from 'vitest'
import {mount} from '@vue/test-utils'
import {createPinia, setActivePinia} from 'pinia'
import ElementPlus from 'element-plus'
import i18n from '@/i18n'
import D2Table from './D2Table.vue'
import type {TableColumn, DrawerColumn} from '@/types/table'

// el-table 在 jsdom 下的行内容渲染依赖真实布局/ResizeObserver，body 行会渲染成空占位（已知的 Element Plus + jsdom
// 生态限制，不是本项目代码问题），所以这里不测"点击表格行内按钮"，只测 D2Table 自己拥有的、不依赖 el-table 内部渲染的
// 逻辑：分页事件上抛（含下面这个真实 bug 的回归保护）+ 权限 prop 对操作列的条件渲染。
interface Row {
  id: number
  name: string
}

const columns: TableColumn[] = [
  {prop: 'id', label: 'ID'},
  {prop: 'name', label: '名称'}
]
const drawerColumns: DrawerColumn[] = [
  {prop: 'id', label: 'ID'},
  {prop: 'name', label: '名称'}
]
const data: Row[] = [{id: 1, name: '张三'}]

function mountTable(props: Record<string, unknown> = {}) {
  return mount(D2Table, {
    props: {
      columns,
      drawerColumns,
      data,
      total: 1,
      pageSize: 10,
      currentPage: 1,
      ...props
    },
    global: {
      plugins: [ElementPlus, i18n]
    }
  })
}

describe('D2Table', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('分页：el-pagination 同时触发 update:current-page 和 current-change（真实组件行为）时，只上抛一次 current-change', async () => {
    // 回归用例：el-pagination 一次翻页会同时 emit update:current-page（走 v-model）和 current-change（原来重复绑了
    // @current-change handler），此前 D2Table 会把 current-change 上抛两次，导致父页面重复请求列表接口。
    const wrapper = mountTable()
    const pagination = wrapper.findComponent({name: 'ElPagination'})

    pagination.vm.$emit('update:current-page', 2)
    pagination.vm.$emit('current-change', 2)
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted('current-change')).toHaveLength(1)
    expect(wrapper.emitted('current-change')?.[0]).toEqual([2])
  })

  it('分页：page-size 变化同理只上抛一次 size-change', async () => {
    const wrapper = mountTable()
    const pagination = wrapper.findComponent({name: 'ElPagination'})

    pagination.vm.$emit('update:page-size', 20)
    pagination.vm.$emit('size-change', 20)
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted('size-change')).toHaveLength(1)
    expect(wrapper.emitted('size-change')?.[0]).toEqual([20])
  })

  it('haveDetail 为 false 时不渲染查看按钮', () => {
    const wrapper = mountTable({haveDetail: false, haveEdit: false})
    const viewBtn = wrapper.findAll('button').find((btn) => btn.text() === '查看')
    expect(viewBtn).toBeUndefined()
  })
})
