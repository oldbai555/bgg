import {describe, it, expect} from 'vitest'
import {isEnvelope} from './envelope'

describe('isEnvelope', () => {
  it('识别标准 Envelope 结构（code 为 number）', () => {
    expect(isEnvelope({code: 0, msg: 'ok', data: null})).toBe(true)
    expect(isEnvelope({code: 10003, msg: 'token expired', data: null})).toBe(true)
  })

  it('code 不是 number 时判定为非 Envelope', () => {
    expect(isEnvelope({code: '0', msg: 'ok', data: null})).toBe(false)
  })

  it('缺少 code 字段时判定为非 Envelope', () => {
    expect(isEnvelope({msg: 'ok', data: null})).toBe(false)
  })

  it('非对象或 null 时判定为非 Envelope', () => {
    expect(isEnvelope(null)).toBe(false)
    expect(isEnvelope(undefined)).toBe(false)
    expect(isEnvelope('string')).toBe(false)
    expect(isEnvelope(123)).toBe(false)
    expect(isEnvelope([])).toBe(false)
  })
})
