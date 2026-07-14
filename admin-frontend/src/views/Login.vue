<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-card__header">
        <div class="logo">Admin System</div>
        <div class="subtitle">{{ t('common.welcome') }}</div>
      </div>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        class="login-form"
      >
        <el-form-item :label="t('auth.username')" prop="username">
          <el-input
            v-model="form.username"
            size="large"
            placeholder="admin"
            autocomplete="username"
            clearable
          />
        </el-form-item>
        <el-form-item :label="t('auth.password')" prop="password">
          <el-input
            v-model="form.password"
            size="large"
            type="password"
            placeholder="••••••"
            show-password
            autocomplete="current-password"
          />
        </el-form-item>
        <el-form-item class="login-actions">
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="full-btn"
            @click="handleSubmit"
          >
            {{ t('common.login') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <IcpFooter class="login-page__footer" />
  </div>
</template>

<script setup lang="ts">
import {reactive, ref} from 'vue'
import {ElForm, ElMessage} from 'element-plus'
import {useRouter} from 'vue-router'
import {useUserStore} from '@/stores/user'
import {useI18n} from 'vue-i18n'
import IcpFooter from '@/components/common/IcpFooter.vue'

const router = useRouter()
const userStore = useUserStore()
const {t} = useI18n()

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{required: true, message: t('auth.username'), trigger: 'blur'}],
  password: [{required: true, message: t('auth.password'), trigger: 'blur'}]
}

const formRef = ref<InstanceType<typeof ElForm>>()
const loading = ref(false)

const handleSubmit = () => {
  formRef.value?.validate(async (valid) => {
    if (!valid) {
return
}
    loading.value = true
    try {
      await userStore.login(form)
      ElMessage.success(t('auth.loginSuccess'))
      // 登录成功后跳转到后台管理首页（Dashboard）
      router.push('/dashboard')
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : t('auth.loginFail')
      ElMessage.error(message)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: grid;
  grid-template-rows: 1fr auto;
  justify-items: center;
  align-items: center;
  background: var(--gradient-login-bg);
  padding: 24px;
  box-sizing: border-box;
  overflow: hidden; /* 防止出现滚动条 */
}

.login-card {
  width: min(420px, calc(100vw - 48px));
  padding: 32px 32px 24px;
  background: var(--color-card);
  border-radius: 16px;
  box-shadow: 0 10px 28px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.04);
}
.login-card__header {
  text-align: center;
  margin-bottom: 24px;
}
.logo {
  font-size: 22px;
  font-weight: 700;
  color: var(--color-primary);
}
.subtitle {
  margin-top: 6px;
  color: var(--color-text-regular);
  font-size: 14px;
}
.login-form :deep(.el-form-item__label) {
  font-weight: 600;
  color: var(--color-text-primary);
}
.login-actions {
  margin-top: 8px;
}
.full-btn {
  width: 100%;
}

.login-page__footer {
  margin: 0;
}
</style>

