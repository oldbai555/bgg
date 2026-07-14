import {describe, it, expect} from 'vitest'
import {generateUniqueRouteName} from './index'

describe('generateUniqueRouteName', () => {
  it('把路径转换成下划线分隔的路由名', () => {
    const used = new Set<string>()
    expect(generateUniqueRouteName('/iam/user', used)).toBe('iam_user')
  })

  it('根路径 / 转换成 root', () => {
    const used = new Set<string>()
    expect(generateUniqueRouteName('/', used)).toBe('root')
  })

  it('重复路径通过递增后缀去重', () => {
    const used = new Set<string>()
    const first = generateUniqueRouteName('/system/config', used)
    const second = generateUniqueRouteName('/system/config', used)
    const third = generateUniqueRouteName('/system/config', used)

    expect(first).toBe('system_config')
    expect(second).toBe('system_config_1')
    expect(third).toBe('system_config_2')
  })

  it('每次生成的名字都会被记录进 usedNames，避免跨调用重复', () => {
    const used = new Set<string>()
    const name = generateUniqueRouteName('/chat/chat', used)
    expect(used.has(name)).toBe(true)
  })
})
