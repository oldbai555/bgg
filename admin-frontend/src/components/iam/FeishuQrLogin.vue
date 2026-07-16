<template>
  <div class="feishu-login">
    <el-icon :size="64" class="feishu-login__icon"><ChatDotRound /></el-icon>
    <p class="feishu-login__desc">{{ t('auth.feishuLoginDesc') }}</p>
    <el-button type="primary" size="large" class="feishu-login__btn" :disabled="!!loadError" @click="startLogin">
      {{ t('auth.loginByFeishu') }}
    </el-button>
    <p v-if="loadError" class="feishu-login__error">{{ loadError }}</p>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {ChatDotRound} from '@element-plus/icons-vue'
import {useI18n} from 'vue-i18n'

// 标准 OAuth 授权码流程（飞书官方文档 2026-07 核实）：直接跳转到飞书授权页，
// 飞书自行根据设备展示扫码/免扫码 UI，授权后浏览器 302 带 code+state 回跳
// redirect_uri，不需要任何 JS SDK 或 postMessage 中转。
// 文档：https://open.feishu.cn/document/common-capabilities/sso/api/obtain-oauth-code
const AUTHORIZE_URL = 'https://accounts.feishu.cn/open-apis/authen/v1/authorize'
const STATE_STORAGE_KEY = 'feishu_login_state'
// 与后台已开通的权限点保持一致（用户信息/工号/手机号/邮箱）
const SCOPES = [
  'directory:employee.base.base:read',
  'contact:user.employee_id:readonly',
  'contact:user.phone:readonly',
  'contact:user.email:readonly'
]

const {t} = useI18n()
const loadError = ref('')

function randomState(): string {
  return `${Date.now().toString(36)}${Math.random().toString(36).slice(2)}`
}

function startLogin() {
  const appId = import.meta.env.VITE_FEISHU_APP_ID as string
  if (!appId) {
    loadError.value = t('auth.feishuMissingAppId')
    return
  }

  // 拼上 Vite base path（本项目部署在 /bgg/ 子路径下，见 vite.config.ts base: '/bgg/'），
  // 必须和飞书开放平台后台「安全设置 - 重定向URL」里登记的地址逐字符一致，否则报 20029。
  const redirectUri = `${window.location.origin}${import.meta.env.BASE_URL}admin/login/feishu/callback`
  const state = randomState()
  sessionStorage.setItem(STATE_STORAGE_KEY, state)

  const params = new URLSearchParams({
    client_id: appId,
    redirect_uri: redirectUri,
    response_type: 'code',
    state,
    scope: SCOPES.join(' ')
  })

  window.location.href = `${AUTHORIZE_URL}?${params.toString()}`
}
</script>

<style scoped>
.feishu-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  min-height: 220px;
  padding: 24px 0;
}
.feishu-login__icon {
  color: var(--color-primary);
}
.feishu-login__desc {
  color: var(--color-text-regular);
  font-size: 14px;
  text-align: center;
  margin: 0;
}
.feishu-login__btn {
  width: 100%;
  max-width: 260px;
}
.feishu-login__error {
  margin-top: 4px;
  color: var(--color-danger, #f56c6c);
  font-size: 14px;
  text-align: center;
}
</style>
