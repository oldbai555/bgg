// Package consts 复制自 internal/consts/consts.go 里 sdk 用到的常量。按
// 16-rpc-conventions.md 第 6 节的既定策略：直接复制到各服务自己的 internal/consts，
// 不做成共享包（量很小，维护成本可忽略），后续两边如需变更需要各自同步改。
package consts

// Open 通用启用状态值（sdk_key.status / sdk_interface.status）。
const Open = 1
