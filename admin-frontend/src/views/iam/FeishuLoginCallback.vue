<template>
  <div class="feishu-callback">
    <el-icon v-if="!errorMessage" class="feishu-callback__spinner" :size="32">
      <Loading />
    </el-icon>
    <p class="feishu-callback__text">{{ errorMessage || t('auth.feishuLoggingIn') }}</p>
    <el-button v-if="errorMessage" type="primary" @click="router.replace('/admin/login')">
      {{ t('common.backToLogin') }}
    </el-button>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {Loading} from '@element-plus/icons-vue'
import {useUserStore} from '@/stores/user'
import {useI18n} from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const {t} = useI18n()

const errorMessage = ref('')

const STATE_STORAGE_KEY = 'feishu_login_state'

onMounted(async () => {
  const code = route.query.code as string | undefined
  const state = (route.query.state as string | undefined) || ''
  const savedState = sessionStorage.getItem(STATE_STORAGE_KEY) || ''
  sessionStorage.removeItem(STATE_STORAGE_KEY)

  if (!code) {
    errorMessage.value = t('auth.feishuMissingCode')
    return
  }
  if (!savedState || state !== savedState) {
    errorMessage.value = t('auth.feishuStateMismatch')
    return
  }

  try {
    await userStore.loginByFeishu(code, state)
    router.push('/admin/dashboard')
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : t('auth.loginFail')
  }
})
</script>

<style scoped>
.feishu-callback {
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
}
.feishu-callback__spinner {
  animation: feishu-callback-spin 1s linear infinite;
  color: var(--color-primary);
}
.feishu-callback__text {
  color: var(--color-text-regular);
  font-size: 14px;
}
@keyframes feishu-callback-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
