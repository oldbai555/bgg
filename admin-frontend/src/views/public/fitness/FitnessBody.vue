<template>
  <div class="fitness-body-page public-detail-page">
    <MetricReporter module="fitness_body" />
    <PublicHeader />
    <FitnessTabBar />
    <div class="page-shell">
      <div class="page-layout">
        <section class="detail-card">
          <h2 class="title">{{ t('fitness.recordBody') }}</h2>
          <el-form label-position="top" class="fitness-body-form" @submit.prevent>
            <el-form-item :label="t('fitness.date')">
              <el-date-picker
                v-model="form.date"
                type="date"
                value-format="YYYY-MM-DD"
                :clearable="false"
                :editable="false"
                :disabled-date="disabledDate"
              />
            </el-form-item>
            <el-form-item :label="t('fitness.weight')">
              <el-input-number
                v-model="form.weightKg"
                :min="0"
                :max="300"
                :step="0.1"
                :precision="1"
              />
            </el-form-item>
            <el-form-item :label="t('fitness.waist')">
              <el-input-number
                v-model="form.waistCm"
                :min="0"
                :max="200"
                :step="0.5"
                :precision="1"
              />
            </el-form-item>
          </el-form>
          <el-button
            type="primary"
            round
            :loading="saving"
            class="fitness-body-form__submit"
            @click="save"
          >
            {{ t('common.save') }}
          </el-button>
        </section>

        <section v-loading="loading" class="detail-card">
          <div class="fitness-body-head">
            <h2 class="title">{{ t('fitness.trend') }}</h2>
            <el-radio-group v-model="days" size="small" @change="load">
              <el-radio-button v-for="n in [30, 90, 180]" :key="n" :value="n">{{ n }} {{ t('fitness.days') }}</el-radio-button>
            </el-radio-group>
          </div>
          <FitnessWeightChart
            :records="records"
            :empty-text="t('fitness.noBodyData')"
            :weight-label="t('fitness.weight')"
            :avg-label="t('fitness.avg7')"
          />
          <el-table
            v-if="records.length"
            :data="recent"
            size="small"
            class="fitness-body-table"
          >
            <el-table-column prop="date" :label="t('fitness.date')" min-width="100" />
            <el-table-column :label="t('fitness.weight')" min-width="90">
              <template #default="{row}">{{ row.weightKg || '-' }}</template>
            </el-table-column>
            <el-table-column :label="t('fitness.avg7')" min-width="90">
              <template #default="{row}">{{ row.weightAvg7 || '-' }}</template>
            </el-table-column>
            <el-table-column :label="t('fitness.waist')" min-width="90">
              <template #default="{row}">{{ row.waistCm || '-' }}</template>
            </el-table-column>
          </el-table>
        </section>
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
import FitnessWeightChart from '@/components/fitness/FitnessWeightChart.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessBodyRecordItem} from '@/api/generated/admin'
import {useFitnessDate} from '@/composables/useFitnessDate'
import {FITNESS_BACKFILL_DAYS} from '@/constants/fitness'
import {addDays, localDateStr} from '@/utils/fitness'

const {t} = useI18n()
const {today} = useFitnessDate()

// 今天以后端（Asia/Shanghai）为准：访问过计划/打卡页就有；直接打开本页时用本机日期兜底，越界由后端校验
const baseToday = computed(() => today.value || localDateStr(new Date()))

const records = ref<FitnessBodyRecordItem[]>([])
const days = ref(90)
const loading = ref(false)
const saving = ref(false)
const form = reactive({date: baseToday.value, weightKg: 0, waistCm: 0})

const recent = computed(() => [...records.value].reverse().slice(0, 14))

const disabledDate = (d: Date): boolean => {
  const s = localDateStr(d)
  return s > baseToday.value || s < addDays(baseToday.value, -FITNESS_BACKFILL_DAYS)
}

const load = async () => {
  loading.value = true
  try {
    const resp = await fitnessApi.myBodyRecordList({days: days.value})
    records.value = [...(resp.list || [])].sort((a, b) => a.date.localeCompare(b.date))
    const latest = records.value[records.value.length - 1]
    if (latest && !form.weightKg) {
      form.weightKg = latest.weightKg
      form.waistCm = latest.waistCm
    }
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('fitness.loadFailed'))
  } finally {
    loading.value = false
  }
}

const save = async () => {
  saving.value = true
  try {
    await fitnessApi.myBodyRecordSave({date: form.date, weightKg: form.weightKg || 0, waistCm: form.waistCm || 0})
    ElMessage.success(t('fitness.saved'))
    await load()
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

.fitness-body-page {
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

.fitness-body-form {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 16px;

  &__submit {
    width: 100%;
  }
}

.fitness-body-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
}

.fitness-body-table {
  margin-top: 16px;
}

@include mobile {
  .fitness-body-page {
    padding-bottom: calc(60px + env(safe-area-inset-bottom));

    .detail-card {
      padding: 16px;
    }
  }

  .fitness-body-form {
    grid-template-columns: 1fr;
  }
}
</style>
