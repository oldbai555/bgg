<template>
  <div class="fitness-login-page public-detail-page">
    <MetricReporter module="fitness_login" />
    <PublicHeader />
    <div class="page-shell">
      <div class="page-layout">
        <section class="detail-card fitness-login">
          <h1 class="title">{{ t('fitness.loginTitle') }}</h1>
          <p class="fitness-login__desc">{{ t('fitness.loginDesc') }}</p>
          <div v-if="working" class="fitness-login__status">
            <el-icon class="fitness-login__spinner" :size="28"><Loading /></el-icon>
            <span>{{ t('fitness.loggingIn') }}</span>
          </div>
          <template v-else>
            <p v-if="errorMessage" class="fitness-login__error">{{ errorMessage }}</p>
            <el-button
              type="primary"
              size="large"
              round
              class="fitness-login__btn"
              @click="startOAuth"
            >
              {{ t('fitness.loginByFeishu') }}
            </el-button>
          </template>
        </section>
      </div>
    </div>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {Loading} from '@element-plus/icons-vue'
import PublicHeader from '@/components/common/PublicHeader.vue'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import {fitnessApi} from '@/api/fitness'
import {useUserStore} from '@/stores/user'
import {FITNESS_HOME_PATH, FITNESS_LOGIN_PATH} from '@/constants/fitness'

// 飞书客户端内：H5 JSSDK tt.requestAccess 免登（拿不到再退到 requestAuthCode，都不行就走 OAuth 跳转）；
// 普通浏览器：OAuth 授权码跳转，回跳到本页带 code+state。两条路后端都用 /authen/v2/oauth/token 换 token，
// 只有 OAuth 需要 redirect_uri（后端按 mode 区分）。只需要登录用户基础信息，不申请额外 scope。
// 文档：https://open.feishu.cn/document/client-docs/h5/development-guide/step-3
const AUTHORIZE_URL = 'https://accounts.feishu.cn/open-apis/authen/v1/authorize'
const JSSDK_URL = 'https://lf1-cdn-tos.bytegoofy.com/goofy/lark/op/h5-js-sdk-1.5.26.js'
const STATE_KEY = 'fitness_login_state'
const REDIRECT_KEY = 'fitness_login_redirect'
const JSSDK_ERRNO_UNSUPPORTED = 103

interface FeishuAuthResult {
  code: string;
}
interface FeishuAuthError {
  errno?: number;
  errString?: string;
}
interface FeishuTT {
  requestAccess?: (opts: {
    appID: string;
    scopeList: string[];
    state?: string;
    success: (res: FeishuAuthResult) => void;
    fail: (err: FeishuAuthError) => void;
  }) => void;
  requestAuthCode?: (opts: {
    appId: string;
    success: (res: FeishuAuthResult) => void;
    fail: (err: FeishuAuthError) => void;
  }) => void;
}
interface FeishuWindow {
  tt?: FeishuTT;
  h5sdk?: {ready: (cb: () => void) => void; error: (cb: (err: unknown) => void) => void};
}

const route = useRoute()
const router = useRouter()
const {t} = useI18n()
const userStore = useUserStore()

const working = ref(true)
const errorMessage = ref('')
let appId = ''

const randomState = () => `${Date.now().toString(36)}${Math.random().toString(36).slice(2)}`
const inFeishuClient = () => /Lark|Feishu/i.test(navigator.userAgent)

const redirectTarget = (): string => {
  const fromQuery = route.query.redirect as string | undefined
  const target = fromQuery || sessionStorage.getItem(REDIRECT_KEY) || FITNESS_HOME_PATH
  sessionStorage.removeItem(REDIRECT_KEY)
  // 只允许站内身材管理路径，避免 redirect 参数被拿来做开放跳转
  return target.startsWith(FITNESS_HOME_PATH) && !target.startsWith(FITNESS_LOGIN_PATH) ? target : FITNESS_HOME_PATH
}

const finish = async (code: string, state: string, mode: 'h5' | 'oauth') => {
  await userStore.loginByFeishuMobile(code, state, mode)
  router.replace(redirectTarget())
}

const fail = (err: unknown) => {
  working.value = false
  errorMessage.value = err instanceof Error && err.message ? err.message : t('fitness.loginFailed')
}

const ensureAppId = async () => {
  if (!appId) {
    appId = (await fitnessApi.loginConfig()).appId
  }
  return appId
}

const startOAuth = async () => {
  try {
    working.value = true
    const id = await ensureAppId()
    const state = randomState()
    sessionStorage.setItem(STATE_KEY, state)
    if (route.query.redirect) {
      sessionStorage.setItem(REDIRECT_KEY, route.query.redirect as string)
    }
    // 必须和后端 Feishu.MobileRedirectUri、飞书开放平台「重定向 URL」白名单逐字符一致
    const redirectUri = `${window.location.origin}${import.meta.env.BASE_URL}${FITNESS_LOGIN_PATH.slice(1)}`
    const params = new URLSearchParams({client_id: id, redirect_uri: redirectUri, response_type: 'code', state})
    window.location.href = `${AUTHORIZE_URL}?${params.toString()}`
  } catch (err) {
    fail(err)
  }
}

const loadJSSDK = (): Promise<void> =>
  new Promise((resolve, reject) => {
    if ((window as unknown as FeishuWindow).h5sdk) {
      resolve()
      return
    }
    const script = document.createElement('script')
    script.src = JSSDK_URL
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('JSSDK load failed'))
    document.head.appendChild(script)
  })

const requestH5Code = async (id: string): Promise<string> => {
  await loadJSSDK()
  const w = window as unknown as FeishuWindow
  return new Promise((resolve, reject) => {
    if (!w.h5sdk || !w.tt) {
      reject(new Error('JSSDK unavailable'))
      return
    }
    const tt = w.tt
    const byAuthCode = () => {
      if (!tt.requestAuthCode) {
        reject(new Error('requestAuthCode unavailable'))
        return
      }
      tt.requestAuthCode({appId: id, success: (res) => resolve(res.code), fail: (e) => reject(new Error(e.errString || 'requestAuthCode failed'))})
    }
    w.h5sdk.error((e) => reject(e instanceof Error ? e : new Error('h5sdk error')))
    w.h5sdk.ready(() => {
      if (!tt.requestAccess) {
        byAuthCode()
        return
      }
      tt.requestAccess({
        appID: id,
        scopeList: [],
        success: (res) => resolve(res.code),
        fail: (e) => (e.errno === JSSDK_ERRNO_UNSUPPORTED ? byAuthCode() : reject(new Error(e.errString || 'requestAccess failed')))
      })
    })
  })
}

onMounted(async () => {
  const code = route.query.code as string | undefined
  if (code) {
    const state = (route.query.state as string | undefined) || ''
    const saved = sessionStorage.getItem(STATE_KEY) || ''
    sessionStorage.removeItem(STATE_KEY)
    if (!saved || saved !== state) {
      fail(new Error(t('fitness.stateMismatch')))
      return
    }
    try {
      await finish(code, state, 'oauth')
    } catch (err) {
      fail(err)
    }
    return
  }

  if (userStore.token) {
    router.replace(redirectTarget())
    return
  }

  if (!inFeishuClient()) {
    working.value = false
    return
  }
  let h5Code = ''
  try {
    h5Code = await requestH5Code(await ensureAppId())
  } catch {
    // 客户端内拿不到免登码（JSSDK 加载失败/版本过低/未配置 H5 可信域名）时退到 OAuth 跳转，飞书内会自动授权
    await startOAuth()
    return
  }
  try {
    await finish(h5Code, '', 'h5')
  } catch (err) {
    // 后端拒绝（如非本企业飞书账号）直接展示，不再重试
    fail(err)
  }
})
</script>

<style scoped lang="scss">
@import '@/styles/public-detail.scss';

.fitness-login-page {
  .page-shell {
    max-width: 480px;
  }

  .page-layout {
    grid-template-columns: 1fr;
    padding-top: 8vh;
  }
}

.fitness-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;

  &__desc {
    margin: 0;
    color: var(--color-text-secondary);
    font-size: 14px;
  }

  &__status {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--color-text-regular);
  }

  &__spinner {
    animation: fitness-login-spin 1s linear infinite;
    color: var(--color-primary);
  }

  &__error {
    margin: 0;
    color: var(--color-danger);
    font-size: 13px;
  }

  &__btn {
    width: 100%;
    margin-top: 8px;
  }
}

@keyframes fitness-login-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
