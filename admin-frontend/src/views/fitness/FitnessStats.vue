<template>
  <div class="page">
    <el-card>
      <el-form :inline="true" @submit.prevent>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="range"
            type="daterange"
            value-format="YYYY-MM-DD"
            range-separator="至"
            start-placeholder="开始"
            end-placeholder="结束"
            :clearable="false"
          />
        </el-form-item>
        <el-form-item label="成员">
          <el-select
            v-model="userId"
            filterable
            remote
            clearable
            :remote-method="searchMembers"
            placeholder="全部成员"
            style="width: 200px"
          >
            <el-option
              v-for="m in memberOptions"
              :key="m.userId"
              :label="m.nickname || m.username"
              :value="m.userId"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="listLoading" @click="search">{{ t('common.search') }}</el-button>
          <el-button v-permission="'fitness_stats:export'" type="success" @click="handleExport">{{ t('common.export') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-row v-loading="overviewLoading" :gutter="12">
      <el-col :xs="24" :md="8">
        <el-card class="stat-card">
          <div class="stat-card__label">今日打卡人数（{{ overview?.date }}）</div>
          <div class="stat-card__value">{{ overview?.todayCheckinCount ?? '-' }} <small>/ {{ overview?.memberCount ?? '-' }} 人</small></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card class="stat-card">
          <div class="stat-card__label">本周平均完成率</div>
          <div class="stat-card__value">{{ overview ? formatRate(overview.weekAvgRate) : '-' }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card class="stat-card">
          <div class="stat-card__label">近 {{ overview?.dailyCounts.length || 0 }} 天每日打卡人数</div>
          <div class="daily-bars">
            <div
              v-for="d in overview?.dailyCounts || []"
              :key="d.date"
              class="daily-bars__item"
              :title="`${d.date}：${d.count} 人`"
            >
              <span class="daily-bars__bar" :style="{height: `${barHeight(d.count)}%`}"></span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="12">
      <el-col :xs="24" :lg="8">
        <el-card header="连续打卡排行">
          <el-table :data="overview?.streakRanking || []" size="small">
            <el-table-column type="index" label="#" width="44" />
            <el-table-column label="成员">
              <template #default="{row}">
                <el-link type="primary" :underline="false" @click="pickMember(row.userId, row.nickname)">{{ row.nickname }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="streakDays" label="连续(天)" width="80" />
            <el-table-column label="本周" width="70">
              <template #default="{row}">{{ formatRate(row.weekRate) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="16">
        <el-card header="打卡明细">
          <el-table
            v-loading="listLoading"
            :data="checkins"
            size="small"
            border
          >
            <el-table-column prop="date" label="日期" width="105" />
            <el-table-column label="成员" width="100">
              <template #default="{row}">
                <el-link type="primary" :underline="false" @click="pickMember(row.userId, row.nickname)">{{ row.nickname }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="templateName" label="模板" width="100" />
            <el-table-column label="训练" width="140">
              <template #default="{row}">
                {{ getTrainingTypeLabel(row.trainingType) }}
                <el-tag v-if="row.trainingStatus" size="small" :type="trainingStatusTagType(row.trainingStatus)">
                  {{ getTrainingStatusLabel(row.trainingStatus) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="动作" width="70">
              <template #default="{row}">{{ row.exerciseDone }}/{{ row.exerciseTotal }}</template>
            </el-table-column>
            <el-table-column label="饮食 计划/别的/没吃" width="140">
              <template #default="{row}">{{ row.mealOnPlan }} / {{ row.mealOther }} / {{ row.mealSkipped }}</template>
            </el-table-column>
            <el-table-column label="步数" width="110">
              <template #default="{row}">{{ row.steps }}<span v-if="row.stepGoal" class="muted"> / {{ row.stepGoal }}</span></template>
            </el-table-column>
            <el-table-column prop="waterCups" label="喝水" width="60" />
            <el-table-column label="完成度" width="80">
              <template #default="{row}">{{ row.score < 0 ? '休息' : formatRate(row.score) }}</template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="listQuery.page"
            v-model:page-size="listQuery.pageSize"
            :total="checkinTotal"
            :page-sizes="[20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            class="pagination"
            @current-change="loadCheckins"
            @size-change="loadCheckins"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-card v-if="userId" v-loading="detailLoading">
      <template #header>
        <div class="detail-head">
          <span>个人数据：{{ detail?.member.nickname || '' }}</span>
          <span v-if="detail" class="muted">
            当前模板 {{ detail.member.templateName || '-' }}｜连续 {{ detail.member.streakDays }} 天｜本周 {{ formatRate(detail.member.weekRate) }}｜区间完成率 {{ formatRate(detail.rangeRate) }}
          </span>
        </div>
      </template>
      <template v-if="detail">
        <h4 class="section-title">打卡热力图</h4>
        <FitnessHeatmap :days="detail.calendar" />
        <h4 class="section-title">体重趋势</h4>
        <FitnessWeightChart :records="detail.bodyRecords" :target="detail.member.targetWeightKg" empty-text="区间内没有身体数据" />
        <h4 class="section-title">模板切换记录</h4>
        <FitnessSwitchTable :data="detail.switches" :show-user="false" />
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref, watch} from 'vue'
import {useRoute} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import FitnessHeatmap from '@/components/fitness/FitnessHeatmap.vue'
import FitnessWeightChart from '@/components/fitness/FitnessWeightChart.vue'
import FitnessSwitchTable from '@/components/fitness/FitnessSwitchTable.vue'
import {fitnessApi} from '@/api/fitness'
import type {
  FitnessCheckinRow,
  FitnessMemberDetailResp,
  FitnessMemberItem,
  FitnessStatsOverviewResp
} from '@/api/generated/admin'
import {useDictOptions} from '@/composables/useDictOptions'
import {addDays, formatRate} from '@/utils/fitness'
import {FITNESS_ADMIN_STATS_PATH, FitnessTrainingStatus} from '@/constants/fitness'

const route = useRoute()
const {t} = useI18n()
const {getLabel: getTrainingTypeLabel} = useDictOptions('fitness_training_type')
const {getLabel: getTrainingStatusLabel} = useDictOptions('fitness_training_status')

const trainingStatusTagType = (status: number): 'success' | 'warning' | 'info' => {
  if (status === FitnessTrainingStatus.Done) {
    return 'success'
  }
  return status === FitnessTrainingStatus.Partial ? 'warning' : 'info'
}

const DEFAULT_RANGE_DAYS = 29

const overview = ref<FitnessStatsOverviewResp | null>(null)
const overviewLoading = ref(false)
const range = ref<[string, string] | null>(null)
const userId = ref<number | undefined>(route.query.userId ? Number(route.query.userId) : undefined)
const memberOptions = ref<Pick<FitnessMemberItem, 'userId' | 'nickname' | 'username'>[]>([])

const checkins = ref<FitnessCheckinRow[]>([])
const checkinTotal = ref(0)
const listLoading = ref(false)
const listQuery = reactive({page: 1, pageSize: 20})

const detail = ref<FitnessMemberDetailResp | null>(null)
const detailLoading = ref(false)

const barHeight = (count: number): number => {
  const max = Math.max(1, ...(overview.value?.dailyCounts || []).map((d) => d.count))
  return Math.max(4, Math.round((count / max) * 100))
}

const filters = () => ({
  userId: userId.value || undefined,
  startDate: range.value?.[0],
  endDate: range.value?.[1]
})

const loadOverview = async () => {
  overviewLoading.value = true
  try {
    overview.value = await fitnessApi.statsOverview({})
    if (!range.value) {
      range.value = [addDays(overview.value.date, -DEFAULT_RANGE_DAYS), overview.value.date]
    }
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '加载失败')
  } finally {
    overviewLoading.value = false
  }
}

const loadCheckins = async () => {
  listLoading.value = true
  try {
    const resp = await fitnessApi.checkinList({page: listQuery.page, pageSize: listQuery.pageSize, ...filters()})
    checkins.value = resp.list || []
    checkinTotal.value = resp.total
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '加载失败')
  } finally {
    listLoading.value = false
  }
}

const loadDetail = async () => {
  if (!userId.value) {
    detail.value = null
    return
  }
  detailLoading.value = true
  try {
    detail.value = await fitnessApi.memberDetail({
      userId: userId.value,
      startDate: range.value?.[0],
      endDate: range.value?.[1]
    })
    const m = detail.value.member
    if (!memberOptions.value.some((o) => o.userId === m.userId)) {
      memberOptions.value = [{userId: m.userId, nickname: m.nickname, username: m.username}, ...memberOptions.value]
    }
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '加载失败')
  } finally {
    detailLoading.value = false
  }
}

const search = () => {
  listQuery.page = 1
  loadCheckins()
  loadDetail()
}

const searchMembers = async (keyword: string) => {
  try {
    const resp = await fitnessApi.memberList({page: 1, pageSize: 20, keyword: keyword || undefined})
    memberOptions.value = resp.list || []
  } catch {
    memberOptions.value = []
  }
}

const pickMember = (id: number, nickname: string) => {
  if (!memberOptions.value.some((o) => o.userId === id)) {
    memberOptions.value = [{userId: id, nickname, username: ''}, ...memberOptions.value]
  }
  userId.value = id
}

watch(userId, search)

// 后台菜单页 keepAlive：从成员列表带不同 userId 跳过来时页面不会重新挂载
watch(
  () => route.query.userId,
  (v) => {
    if (v && route.path === FITNESS_ADMIN_STATS_PATH) {
      userId.value = Number(v)
    }
  }
)

const handleExport = async () => {
  try {
    await fitnessApi.checkinExport(filters())
    ElMessage.success('已创建异步导出任务，请在右下角任务列表查看进度')
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.exportFail'))
  }
}

onMounted(async () => {
  await loadOverview()
  searchMembers('')
  search()
})
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-card {
  height: 100%;
}

.stat-card__label {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.stat-card__value {
  margin-top: 8px;
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.stat-card__value small {
  font-size: 13px;
  font-weight: 400;
  color: var(--color-text-secondary);
}

.daily-bars {
  display: flex;
  align-items: flex-end;
  gap: 3px;
  height: 48px;
  margin-top: 8px;
}

.daily-bars__item {
  flex: 1;
  height: 100%;
  display: flex;
  align-items: flex-end;
}

.daily-bars__bar {
  width: 100%;
  border-radius: 2px 2px 0 0;
  background: var(--color-primary);
}

.pagination {
  margin-top: 12px;
  justify-content: flex-end;
}

.detail-head {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: baseline;
}

.muted {
  color: var(--color-text-secondary);
  font-size: 12.5px;
}

.section-title {
  margin: 16px 0 8px;
  font-size: 14px;
  color: var(--color-text-primary);
}

.section-title:first-child {
  margin-top: 0;
}
</style>
