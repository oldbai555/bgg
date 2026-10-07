<template>
  <div class="fitness-today-page public-detail-page">
    <MetricReporter module="fitness_today" />
    <PublicHeader />
    <FitnessTabBar />
    <div class="page-shell">
      <div v-loading="loading" class="page-layout">
        <section class="detail-card fitness-head">
          <FitnessDayBar
            v-if="date"
            :date="date"
            :today="today"
            @update:date="changeDate"
          />
          <div v-if="day" class="fitness-head__meta">
            <span>{{ day.templateName || '-' }}</span>
            <el-tag
              v-if="day.isOverride"
              size="small"
              type="warning"
              effect="plain"
            >{{ t('fitness.override') }}</el-tag>
            <el-tag
              v-if="day.checkin"
              size="small"
              :type="day.checkin.trainingStatus ? 'success' : 'info'"
              effect="plain"
            >
              {{ checkinSummary }}
            </el-tag>
          </div>
        </section>

        <FitnessSuggestionCard
          v-if="day?.suggestion && day.suggestion.status === FitnessSuggestionStatus.Pending"
          :suggestion="day.suggestion"
          @decided="load"
        />

        <template v-if="day?.hasPlan && day.content">
          <section class="detail-card">
            <h2 class="title">{{ t('fitness.training') }}</h2>
            <FitnessTrainingBlock :content="day.content" />
          </section>

          <section class="detail-card">
            <h2 class="title">{{ t('fitness.meals') }}</h2>
            <FitnessMealList :meals="day.content.meals || []" />
          </section>
        </template>
        <section v-else-if="day" class="detail-card">
          <div class="empty">{{ t('fitness.noPlan') }}</div>
        </section>

        <section v-if="day && (day.content?.tip || day.tips.length)" class="detail-card">
          <h2 class="title">{{ t('fitness.tips') }}</h2>
          <p v-if="day.content?.tip" class="fitness-tip">{{ day.content.tip }}</p>
          <el-collapse>
            <el-collapse-item
              v-for="tip in day.tips"
              :key="tip.id"
              :title="tip.title"
              :name="tip.id"
            >
              <div class="content fitness-tip__content">{{ tip.content }}</div>
            </el-collapse-item>
          </el-collapse>
        </section>

        <div v-if="day?.canCheckin" class="fitness-actions">
          <el-button
            type="primary"
            size="large"
            round
            @click="router.push('/front/fitness/checkin')"
          >
            {{ t('fitness.goCheckin') }}
          </el-button>
        </div>
      </div>
    </div>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import {useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import PublicHeader from '@/components/common/PublicHeader.vue'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import FitnessTabBar from '@/components/fitness/FitnessTabBar.vue'
import FitnessDayBar from '@/components/fitness/FitnessDayBar.vue'
import FitnessTrainingBlock from '@/components/fitness/FitnessTrainingBlock.vue'
import FitnessMealList from '@/components/fitness/FitnessMealList.vue'
import FitnessSuggestionCard from '@/components/fitness/FitnessSuggestionCard.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessMyDayResp} from '@/api/generated/admin'
import {useFitnessDicts} from '@/composables/useFitnessDicts'
import {useFitnessDate} from '@/composables/useFitnessDate'
import {FitnessSuggestionStatus} from '@/constants/fitness'

const router = useRouter()
const {t} = useI18n()
const {label} = useFitnessDicts()
const {date, today, syncToday} = useFitnessDate()

const day = ref<FitnessMyDayResp | null>(null)
const loading = ref(false)

const checkinSummary = computed(() => {
  const c = day.value?.checkin
  if (!c) {
    return ''
  }
  if (c.trainingStatus) {
    return label('fitness_training_status', c.trainingStatus)
  }
  return t('fitness.tabCheckin')
})

const load = async () => {
  loading.value = true
  try {
    const resp = await fitnessApi.myDay({date: date.value || undefined})
    if (syncToday(resp.today)) {
      return await load()
    }
    date.value = resp.date
    day.value = resp
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

onMounted(load)
</script>

<style scoped lang="scss">
@import '@/styles/public-detail.scss';

.fitness-today-page {
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

.fitness-head__meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.fitness-tip {
  margin: 0 0 8px;
  font-size: 14px;
  line-height: 1.7;
  color: var(--color-text-primary);
}

.fitness-tip__content {
  white-space: pre-line;
}

.fitness-actions {
  display: flex;
  justify-content: center;
}

@include mobile {
  // 给底部固定 tab 栏让位，IcpFooter 仍在可视区内
  .fitness-today-page {
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
