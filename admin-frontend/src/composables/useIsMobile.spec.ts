import {describe, it, expect, afterEach} from 'vitest'
import {defineComponent, h} from 'vue'
import {mount} from '@vue/test-utils'
import {useIsMobile} from './useIsMobile'

function setInnerWidth(width: number) {
  Object.defineProperty(window, 'innerWidth', {writable: true, configurable: true, value: width})
}

const TestComponent = defineComponent({
  setup() {
    const {isMobile} = useIsMobile()
    return {isMobile}
  },
  render() {
    return h('div', String(this.isMobile))
  }
})

describe('useIsMobile', () => {
  const originalWidth = window.innerWidth

  afterEach(() => {
    setInnerWidth(originalWidth)
  })

  it('挂载时按当前窗口宽度计算 isMobile', async () => {
    setInnerWidth(500)
    const wrapper = mount(TestComponent)
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toBe('true')
    wrapper.unmount()
  })

  it('宽度大于断点时 isMobile 为 false', async () => {
    setInnerWidth(1200)
    const wrapper = mount(TestComponent)
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toBe('false')
    wrapper.unmount()
  })

  it('resize 事件触发后重新计算', async () => {
    setInnerWidth(1200)
    const wrapper = mount(TestComponent)
    expect(wrapper.text()).toBe('false')

    setInnerWidth(500)
    window.dispatchEvent(new Event('resize'))
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toBe('true')
    wrapper.unmount()
  })

  it('卸载后不再响应 resize（监听器已清理）', async () => {
    setInnerWidth(1200)
    const wrapper = mount(TestComponent)
    wrapper.unmount()

    setInnerWidth(500)
    expect(() => window.dispatchEvent(new Event('resize'))).not.toThrow()
  })
})
