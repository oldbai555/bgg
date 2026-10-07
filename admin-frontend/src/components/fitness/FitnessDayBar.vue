<template>
  <div class="fitness-daybar">
    <el-button
      :icon="ArrowLeft"
      circle
      :aria-label="t('fitness.prevDay')"
      @click="shift(-1)"
    />
    <div class="fitness-daybar__center">
      <el-date-picker
        :model-value="date"
        type="date"
        value-format="YYYY-MM-DD"
        format="MM月DD日"
        :clearable="false"
        :editable="false"
        :disabled-date="disabledDate"
        class="fitness-daybar__picker"
        @update:model-value="onPick"
      />
      <span class="fitness-daybar__weekday">{{ weekdayLabel(weekdayOf(date)) }}</span>
      <el-tag
        v-if="date === today"
        size="small"
        type="success"
        effect="plain"
      >{{ t('fitness.today') }}</el-tag>
      <el-link
        v-else-if="today"
        type="primary"
        :underline="false"
        @click="emit('update:date', today)"
      >
        {{ t('fitness.today') }}
      </el-link>
    </div>
    <el-button
      :icon="ArrowRight"
      circle
      :aria-label="t('fitness.nextDay')"
      :disabled="!!maxDate && date >= maxDate"
      @click="shift(1)"
    />
  </div>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n'
import {ArrowLeft, ArrowRight} from '@element-plus/icons-vue'
import {addDays, localDateStr, weekdayLabel, weekdayOf} from '@/utils/fitness'

const props = defineProps<{
  date: string;
  today: string;
  /** 不传表示可以往后翻（看未来计划）；打卡页传 today */
  maxDate?: string;
  minDate?: string;
}>()

const emit = defineEmits<{(e: 'update:date', value: string): void}>()
const {t} = useI18n()

const disabledDate = (d: Date): boolean => {
  const s = localDateStr(d)
  return (!!props.maxDate && s > props.maxDate) || (!!props.minDate && s < props.minDate)
}

const shift = (n: number) => {
  const next = addDays(props.date, n)
  if ((props.maxDate && next > props.maxDate) || (props.minDate && next < props.minDate)) {
    return
  }
  emit('update:date', next)
}

const onPick = (value: string | null) => {
  if (value) {
    emit('update:date', value)
  }
}
</script>

<style scoped lang="scss">
.fitness-daybar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;

  &__center {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
  }

  &__picker {
    width: 132px;
  }

  &__weekday {
    font-weight: 600;
    color: var(--color-text-primary);
    white-space: nowrap;
  }
}

@include mobile {
  .fitness-daybar__picker {
    width: 112px;
  }
}
</style>
