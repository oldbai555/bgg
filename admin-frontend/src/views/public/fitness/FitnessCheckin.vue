<template>
  <div class="fitness-checkin-page public-detail-page">
    <MetricReporter module="fitness_checkin" />
    <PublicHeader />
    <FitnessTabBar />
    <div class="page-shell">
      <div v-loading="loading" class="page-layout">
        <section class="detail-card">
          <FitnessDayBar
            v-if="date && today"
            :date="date"
            :today="today"
            :max-date="today"
            :min-date="minDate"
            @update:date="changeDate"
          />
          <el-alert
            v-if="day && !day.canCheckin"
            :title="t('fitness.checkinClosed')"
            type="info"
            :closable="false"
            show-icon
            class="fitness-checkin__alert"
          />
        </section>

        <template v-if="day">
          <section class="detail-card">
            <h2 class="title">{{ t('fitness.training') }}</h2>
            <FitnessTrainingBlock
              v-if="day.content"
              :done="form.done"
              :content="day.content"
              checkable
              :disabled="!day.canCheckin"
              @update:done="onDoneChange"
            />
            <div class="fitness-field">
              <span class="fitness-field__label">{{ t('fitness.trainingStatus') }}</span>
              <el-radio-group v-model="form.trainingStatus" :disabled="!day.canCheckin">
                <el-radio-button v-for="opt in options('fitness_training_status')" :key="opt.value" :value="Number(opt.value)">
                  {{ opt.label }}
                </el-radio-button>
              </el-radio-group>
            </div>
          </section>

          <section class="detail-card">
            <h2 class="title">{{ t('fitness.meals') }}</h2>
            <ul class="fitness-meal-checks">
              <li v-for="m in mealRows" :key="m.slot" class="fitness-meal-checks__item">
                <div class="fitness-meal-checks__head">
                  <b>{{ label('fitness_meal_slot', m.slot) }}</b>
                  <span v-if="m.food" class="fitness-meal-checks__food">{{ m.food }}</span>
                </div>
                <el-radio-group v-model="form.meals[m.slot].status" size="small" :disabled="!day.canCheckin">
                  <el-radio-button v-for="opt in options('fitness_meal_status')" :key="opt.value" :value="Number(opt.value)">
                    {{ opt.label }}
                  </el-radio-button>
                </el-radio-group>
                <el-input
                  v-if="form.meals[m.slot].status === FitnessMealStatus.Other"
                  v-model="form.meals[m.slot].note"
                  :placeholder="t('fitness.mealNote')"
                  maxlength="100"
                  size="small"
                  :disabled="!day.canCheckin"
                />
              </li>
            </ul>
          </section>

          <section class="detail-card fitness-numbers">
            <div class="fitness-field">
              <span class="fitness-field__label">
                {{ t('fitness.steps') }}
                <small v-if="stepGoal">/ {{ stepGoal }}</small>
              </span>
              <el-input-number
                v-model="form.steps"
                :min="0"
                :max="200000"
                :step="500"
                :disabled="!day.canCheckin"
              />
            </div>
            <div class="fitness-field">
              <span class="fitness-field__label">{{ t('fitness.water') }}</span>
              <el-input-number
                v-model="form.waterCups"
                :min="0"
                :max="30"
                :disabled="!day.canCheckin"
              />
            </div>
          </section>

          <div v-if="day.canCheckin" class="fitness-actions">
            <el-button
              type="primary"
              size="large"
              round
              :loading="saving"
              @click="save"
            >
              {{ t('fitness.saveCheckin') }}
            </el-button>
          </div>
        </template>
      </div>
    </div>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import PublicHeader from '@/components/common/PublicHeader.vue'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import FitnessTabBar from '@/components/fitness/FitnessTabBar.vue'
import FitnessDayBar from '@/components/fitness/FitnessDayBar.vue'
import FitnessTrainingBlock from '@/components/fitness/FitnessTrainingBlock.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessMealCheck, FitnessMyDayResp} from '@/api/generated/admin'
import {useFitnessDicts} from '@/composables/useFitnessDicts'
import {useFitnessDate} from '@/composables/useFitnessDate'
import {
  FITNESS_BACKFILL_DAYS,
  FITNESS_MEAL_SLOTS,
  FitnessMealStatus,
  FitnessTrainingStatus
} from '@/constants/fitness'
import {addDays, canCheckinOn, suggestTrainingStatus} from '@/utils/fitness'

const {t} = useI18n()
const {label, options} = useFitnessDicts()
const {date, today, syncToday} = useFitnessDate()

const day = ref<FitnessMyDayResp | null>(null)
const loading = ref(false)
const saving = ref(false)

const emptyMeals = () =>
  Object.fromEntries(FITNESS_MEAL_SLOTS.map((slot) => [slot, {status: 0, note: ''}])) as Record<
    string,
    {status: number; note: string}
  >

const form = reactive({
  done: [] as number[],
  trainingStatus: 0,
  meals: emptyMeals(),
  steps: 0,
  waterCups: 0
})

const minDate = computed(() => (today.value ? addDays(today.value, -FITNESS_BACKFILL_DAYS) : ''))
const exerciseTotal = computed(() => day.value?.content?.exercises?.length || 0)
const stepGoal = computed(() => day.value?.content?.stepGoal || day.value?.checkin?.stepGoal || 0)

const mealRows = computed(() => {
  const planned = day.value?.content?.meals || []
  if (planned.length === 0) {
    return FITNESS_MEAL_SLOTS.map((slot) => ({slot, food: ''}))
  }
  return planned.map((m) => ({slot: m.slot, food: m.food || ''}))
})

// 勾动作时自动带出训练状态；用户主动选了「跳过」就不再覆盖
const onDoneChange = (done: number[]) => {
  form.done = done
  if (form.trainingStatus !== FitnessTrainingStatus.Skipped) {
    form.trainingStatus = suggestTrainingStatus(done.length, exerciseTotal.value)
  }
}

const fillForm = (resp: FitnessMyDayResp) => {
  const c = resp.checkin
  const meals = emptyMeals()
  c?.meals?.forEach((m) => {
    meals[m.slot] = {status: m.status, note: m.note || ''}
  })
  form.meals = meals
  form.done = c?.exerciseDone ? [...c.exerciseDone] : []
  form.trainingStatus = c?.trainingStatus || 0
  form.steps = c?.steps || 0
  form.waterCups = c?.waterCups || 0
}

const load = async () => {
  loading.value = true
  try {
    const resp = await fitnessApi.myDay({date: date.value || undefined})
    if (syncToday(resp.today)) {
      return await load()
    }
    // 从计划页带过来的未来日期/超出补卡窗口的日期，打卡页回到今天
    if (!canCheckinOn(resp.date, resp.today, FITNESS_BACKFILL_DAYS)) {
      date.value = resp.today
      return await load()
    }
    date.value = resp.date
    day.value = resp
    fillForm(resp)
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('fitness.loadFailed'))
  } finally {
    loading.value = false
  }
}

const changeDate = (value: string) => {
  date.value = value
  load()
}

const save = async () => {
  saving.value = true
  try {
    const meals: FitnessMealCheck[] = Object.entries(form.meals)
      .filter(([, m]) => m.status > 0)
      .map(([slot, m]) => ({slot, status: m.status, note: m.status === FitnessMealStatus.Other ? m.note : ''}))
    const checkin = await fitnessApi.myCheckinSave({
      date: date.value,
      trainingStatus: form.trainingStatus,
      exerciseDone: form.done,
      meals,
      steps: form.steps,
      waterCups: form.waterCups
    })
    if (day.value) {
      day.value.checkin = checkin
    }
    form.trainingStatus = checkin.trainingStatus
    ElMessage.success(t('fitness.saved'))
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped lang="scss">
@import '@/styles/public-detail.scss';

.fitness-checkin-page {
  .page-shell {
    max-width: 760px;
  }

  .page-layout {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .detail-card {
    padding: 20px 24px;

    .title {
      font-size: 17px;
      margin-bottom: 12px;
    }
  }
}

.fitness-checkin__alert {
  margin-top: 12px;
}

.fitness-field {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;

  &__label {
    font-size: 14px;
    font-weight: 600;
    color: var(--color-text-primary);

    small {
      font-weight: 400;
      color: var(--color-text-secondary);
    }
  }
}

.fitness-numbers .fitness-field:first-child {
  margin-top: 0;
}

.fitness-meal-checks {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 14px;

  &__item {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding-bottom: 14px;
    border-bottom: 1px solid var(--color-border-light);

    &:last-child {
      border-bottom: none;
      padding-bottom: 0;
    }
  }

  &__head {
    display: flex;
    align-items: baseline;
    gap: 8px;
    font-size: 14px;
    color: var(--color-text-primary);
  }

  &__food {
    font-size: 12.5px;
    color: var(--color-text-secondary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.fitness-actions {
  display: flex;
  justify-content: center;
}

@include mobile {
  .fitness-checkin-page {
    padding-bottom: calc(60px + env(safe-area-inset-bottom));

    .detail-card {
      padding: 16px;
    }
  }

  .fitness-actions .el-button {
    width: 100%;
  }
}
</style>
