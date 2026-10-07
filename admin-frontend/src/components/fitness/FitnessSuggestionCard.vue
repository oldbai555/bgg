<template>
  <div class="fitness-suggestion">
    <div class="fitness-suggestion__title">
      <el-icon><Opportunity /></el-icon>
      {{ t('fitness.suggestion') }}：{{ suggestion.fromTemplateName }} → <b>{{ suggestion.toTemplateName }}</b>
    </div>
    <p class="fitness-suggestion__reason">{{ suggestion.reason }}</p>
    <p class="fitness-suggestion__hint">{{ t('fitness.switchHint') }}</p>
    <div class="fitness-suggestion__actions">
      <el-button size="small" :loading="submitting" @click="decide(false)">{{ t('fitness.reject') }}</el-button>
      <el-button
        size="small"
        type="primary"
        :loading="submitting"
        @click="decide(true)"
      >{{ t('fitness.accept') }}</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {Opportunity} from '@element-plus/icons-vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessSuggestionItem} from '@/api/generated/admin'

const props = defineProps<{suggestion: FitnessSuggestionItem}>()
const emit = defineEmits<{(e: 'decided', accepted: boolean): void}>()
const {t} = useI18n()
const submitting = ref(false)

const decide = async (accept: boolean) => {
  submitting.value = true
  try {
    await fitnessApi.mySuggestionDecide({id: props.suggestion.id, accept})
    ElMessage.success(t('fitness.saved'))
    emit('decided', accept)
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped lang="scss">
.fitness-suggestion {
  padding: 14px 16px;
  border-radius: 10px;
  border: 1px solid var(--color-primary);
  background: var(--color-bg-card);
  display: flex;
  flex-direction: column;
  gap: 6px;

  &__title {
    display: flex;
    align-items: center;
    gap: 6px;
    font-weight: 600;
    color: var(--color-text-primary);
  }

  &__reason {
    margin: 0;
    font-size: 13.5px;
    color: var(--color-text-regular);
    line-height: 1.6;
  }

  &__hint {
    margin: 0;
    font-size: 12px;
    color: var(--color-text-secondary);
  }

  &__actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}
</style>
