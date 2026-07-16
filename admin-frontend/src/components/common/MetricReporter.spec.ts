import {describe, it, expect, vi, beforeEach} from 'vitest'
import {mount} from '@vue/test-utils'
import MetricReporter from './MetricReporter.vue'

vi.mock('@/api/monitoring', () => ({
  monitoringApi: {
    metricReport: vi.fn(() => Promise.resolve())
  }
}))

import {monitoringApi} from '@/api/monitoring'

describe('MetricReporter', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('挂载时按 props 声明式上报一次（默认 event=view）', () => {
    mount(MetricReporter, {props: {module: 'video_list', bizId: 0}})
    expect(monitoringApi.metricReport).toHaveBeenCalledWith({module: 'video_list', bizId: 0, event: 'view'})
  })

  it('命令式 report() 不传覆盖参数时，行为与声明式上报一致', () => {
    const wrapper = mount(MetricReporter, {props: {module: 'video_detail', bizId: 5, event: 'view'}})
    vi.clearAllMocks()

    wrapper.vm.report()
    expect(monitoringApi.metricReport).toHaveBeenCalledWith({module: 'video_detail', bizId: 5, event: 'view'})
  })

  it('命令式 report() 可覆盖 event/bizId，用于"播放"这类一次性业务事件', () => {
    const wrapper = mount(MetricReporter, {props: {module: 'video_detail', bizId: 5}})
    vi.clearAllMocks()

    wrapper.vm.report({event: 'play', bizId: 5})
    expect(monitoringApi.metricReport).toHaveBeenCalledWith({module: 'video_detail', bizId: 5, event: 'play'})
  })

  it('enabled=false 时命令式 report() 也不上报', () => {
    const wrapper = mount(MetricReporter, {props: {module: 'video_detail', bizId: 5, enabled: false}})
    vi.clearAllMocks()

    wrapper.vm.report({event: 'play'})
    expect(monitoringApi.metricReport).not.toHaveBeenCalled()
  })
})
