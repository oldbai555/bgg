<template>
  <div class="fitness-training">
    <div class="fitness-training__head">
      <el-tag :type="typeTag" effect="dark" round>{{ dictLabel('fitness_training_type', content.trainingType) }}</el-tag>
      <span v-if="content.trainingTime" class="fitness-training__time">{{ content.trainingTime }}</span>
      <span v-if="content.stepGoal" class="fitness-training__steps">
        {{ t('fitness.stepGoal') }} {{ content.stepGoal }}
      </span>
    </div>

    <p v-if="content.warmup" class="fitness-training__line">
      <b>{{ t('fitness.warmup') }}</b>{{ content.warmup }}
    </p>

    <ol v-if="exercises.length" class="fitness-training__list">
      <li v-for="(ex, i) in exercises" :key="i" class="fitness-training__item">
        <el-checkbox
          v-if="checkable"
          :model-value="done.includes(i)"
          :disabled="disabled"
          @update:model-value="(v: boolean | string | number) => toggle(i, !!v)"
        />
        <div class="fitness-training__item-body">
          <div class="fitness-training__item-name">
            {{ ex.name }}
            <span class="fitness-training__dose">{{ dose(ex) }}</span>
          </div>
          <div v-if="ex.note" class="fitness-training__note">{{ ex.note }}</div>
        </div>
      </li>
    </ol>

    <p v-if="content.stretch" class="fitness-training__line">
      <b>{{ t('fitness.stretch') }}</b>{{ content.stretch }}
    </p>
    <p v-if="content.trainingNote" class="fitness-training__note">{{ content.trainingNote }}</p>
  </div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import type {FitnessDayContent, FitnessExercise} from '@/api/generated/admin'
import {useDictStore} from '@/stores/dict'
import {FitnessTrainingType} from '@/constants/fitness'

const props = withDefaults(
  defineProps<{
    content: FitnessDayContent;
    /** 打卡页可勾选动作 */
    checkable?: boolean;
    done?: number[];
    disabled?: boolean;
  }>(),
  {checkable: false, done: () => [], disabled: false}
)

const emit = defineEmits<{(e: 'update:done', value: number[]): void}>()
const {t} = useI18n()
const dictStore = useDictStore()

const dictLabel = (code: string, value: number) => dictStore.getDictLabel(code, value)
const exercises = computed(() => props.content.exercises || [])

const typeTag = computed(() => {
  switch (props.content.trainingType) {
    case FitnessTrainingType.Strength:
      return 'danger'
    case FitnessTrainingType.Cardio:
      return 'warning'
    default:
      return 'info'
  }
})

const dose = (ex: FitnessExercise): string => {
  if (ex.sets && ex.reps) {
    return `${ex.sets}×${ex.reps}`
  }
  return ex.duration || ex.reps || ''
}

const toggle = (index: number, checked: boolean) => {
  const set = new Set(props.done)
  if (checked) {
    set.add(index)
  } else {
    set.delete(index)
  }
  emit('update:done', [...set].sort((a, b) => a - b))
}
</script>

<style scoped lang="scss">
.fitness-training {
  display: flex;
  flex-direction: column;
  gap: 10px;
  font-size: 14px;
  color: var(--color-text-regular);

  &__head {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
  }

  &__time {
    font-weight: 600;
    color: var(--color-text-primary);
  }

  &__steps {
    color: var(--color-text-secondary);
    font-size: 13px;
  }

  &__line {
    margin: 0;
    line-height: 1.7;

    b {
      color: var(--color-text-primary);
      margin-right: 8px;
    }
  }

  &__list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  &__item {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    padding: 10px 12px;
    border-radius: 8px;
    background: var(--color-bg-secondary);
  }

  &__item-body {
    flex: 1;
    min-width: 0;
  }

  &__item-name {
    font-weight: 600;
    color: var(--color-text-primary);
  }

  &__dose {
    margin-left: 8px;
    font-weight: 500;
    color: var(--color-primary);
  }

  &__note {
    margin: 2px 0 0;
    font-size: 12.5px;
    color: var(--color-text-secondary);
    line-height: 1.6;
  }
}
</style>
