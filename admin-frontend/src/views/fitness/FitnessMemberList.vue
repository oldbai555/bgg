<template>
  <div class="page">
    <el-card>
      <el-form :inline="true" :model="query" @submit.prevent>
        <el-form-item label="成员">
          <el-input v-model="query.keyword" placeholder="昵称/用户名" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="search">{{ t('common.search') }}</el-button>
          <el-button @click="reset">{{ t('common.reset') }}</el-button>
          <el-button @click="openSwitches(null)">全部切换记录</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        :drawer-columns="[]"
        :have-edit="false"
        :have-detail="false"
        :action-column-width="230"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      >
        <template #cell="{row, column}">
          <div v-if="column.prop === 'nickname'" class="member-cell">
            <el-avatar :size="24" :src="row.avatar">{{ (row.nickname || '').slice(0, 1) }}</el-avatar>
            <span>{{ row.nickname }}</span>
          </div>
          <span v-else-if="column.prop === 'weekRate'">{{ formatRate(row.weekRate) }}</span>
          <span v-else-if="column.prop === 'streakDays'">{{ row.streakDays ?? 0 }} 天</span>
          <span v-else-if="column.prop === 'weight'">
            {{ row.latestWeightKg || '-' }} / {{ row.targetWeightKg || '-' }}
          </span>
          <span v-else>{{ row[column.prop as keyof FitnessMemberItem] || '-' }}</span>
        </template>
        <template #action="{row}">
          <el-button
            v-permission="'fitness_member:assign'"
            size="small"
            type="primary"
            link
            @click="openAssign(row)"
          >
            指定模板
          </el-button>
          <el-button
            size="small"
            type="primary"
            link
            @click="openSwitches(row)"
          >切换记录</el-button>
          <el-button
            v-permission="'fitness_stats:list'"
            size="small"
            type="primary"
            link
            @click="goStats(row)"
          >
            打卡统计
          </el-button>
        </template>
      </D2Table>
    </el-card>

    <el-dialog v-model="assignVisible" :title="`指定模板：${assignTarget?.nickname || ''}`" width="480px">
      <el-form label-width="80px" @submit.prevent>
        <el-form-item label="当前模板">{{ assignTarget?.templateName || '-' }}</el-form-item>
        <el-form-item label="新模板" required>
          <el-select v-model="assignForm.templateId" placeholder="选择模板" style="width: 100%">
            <el-option
              v-for="tpl in templates"
              :key="tpl.id"
              :label="`${tpl.name}（${tpl.stage}）`"
              :value="tpl.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="原因">
          <el-input
            v-model="assignForm.reason"
            type="textarea"
            :rows="2"
            maxlength="200"
            placeholder="如：出差两周改忙碌周"
          />
        </el-form-item>
        <el-alert
          type="info"
          :closable="false"
          show-icon
          title="从明天开始生效，今天已有的打卡不受影响；明天之前再次指定会覆盖这次。"
        />
      </el-form>
      <template #footer>
        <el-button @click="assignVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="assigning"
          :disabled="!assignForm.templateId"
          @click="submitAssign"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="switchVisible" :title="switchUser ? `切换记录：${switchUser.nickname}` : '全部切换记录'" width="960px">
      <FitnessSwitchTable v-loading="switchLoading" :data="switches" :show-user="!switchUser" />
      <el-pagination
        v-model:current-page="switchQuery.page"
        :page-size="switchQuery.pageSize"
        :total="switchTotal"
        layout="total, prev, pager, next"
        class="switch-pagination"
        @current-change="loadSwitches"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import {useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import D2Table from '@/components/common/D2Table.vue'
import FitnessSwitchTable from '@/components/fitness/FitnessSwitchTable.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessMemberItem, FitnessSwitchItem, FitnessTemplateItem} from '@/api/generated/admin'
import {FITNESS_ADMIN_STATS_PATH, FitnessStatus} from '@/constants/fitness'
import type {TableColumn} from '@/types/table'
import {formatRate} from '@/utils/fitness'

const router = useRouter()
const {t} = useI18n()

const query = reactive({page: 1, pageSize: 20, keyword: ''})
const list = ref<FitnessMemberItem[]>([])
const total = ref(0)
const loading = ref(false)

const templates = ref<FitnessTemplateItem[]>([])
const assignVisible = ref(false)
const assigning = ref(false)
const assignTarget = ref<FitnessMemberItem | null>(null)
const assignForm = reactive({templateId: 0, reason: ''})

const switchVisible = ref(false)
const switchLoading = ref(false)
const switchUser = ref<FitnessMemberItem | null>(null)
const switches = ref<FitnessSwitchItem[]>([])
const switchTotal = ref(0)
const switchQuery = reactive({page: 1, pageSize: 10})

const columns = computed<TableColumn[]>(() => [
  {prop: 'userId', label: '用户ID', width: 80},
  {prop: 'nickname', label: '成员', width: 150},
  {prop: 'username', label: '用户名', width: 150},
  {prop: 'templateName', label: '当前模板', width: 120},
  {prop: 'weight', label: '最新/目标体重', width: 120},
  {prop: 'streakDays', label: '连续打卡', width: 90},
  {prop: 'weekRate', label: '本周完成率', width: 100},
  {prop: 'lastCheckinDate', label: '最近打卡', width: 110}
])

const loadData = async () => {
  loading.value = true
  try {
    const resp = await fitnessApi.memberList({page: query.page, pageSize: query.pageSize, keyword: query.keyword || undefined})
    list.value = resp.list || []
    total.value = resp.total
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '加载失败')
  } finally {
    loading.value = false
  }
}

const search = () => {
  query.page = 1
  loadData()
}

const reset = () => {
  query.keyword = ''
  search()
}

const handlePageChange = (page: number) => {
  query.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  query.pageSize = size
  query.page = 1
  loadData()
}

const openAssign = async (row: FitnessMemberItem) => {
  assignTarget.value = row
  assignForm.templateId = 0
  assignForm.reason = ''
  assignVisible.value = true
  if (templates.value.length === 0) {
    try {
      const resp = await fitnessApi.templateList({page: 1, pageSize: 100})
      templates.value = (resp.list || []).filter((tpl) => tpl.status === FitnessStatus.Enabled)
    } catch (err: unknown) {
      ElMessage.error(err instanceof Error ? err.message : '加载模板失败')
    }
  }
}

const submitAssign = async () => {
  if (!assignTarget.value) {
    return
  }
  assigning.value = true
  try {
    const resp = await fitnessApi.memberAssign({
      userId: assignTarget.value.userId,
      templateId: assignForm.templateId,
      reason: assignForm.reason
    })
    ElMessage.success(`已指定，${resp.effectiveDate} 起生效`)
    assignVisible.value = false
    loadData()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  } finally {
    assigning.value = false
  }
}

const loadSwitches = async () => {
  switchLoading.value = true
  try {
    const resp = await fitnessApi.switchList({
      userId: switchUser.value?.userId,
      page: switchQuery.page,
      pageSize: switchQuery.pageSize
    })
    switches.value = resp.list || []
    switchTotal.value = resp.total
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '加载失败')
  } finally {
    switchLoading.value = false
  }
}

const openSwitches = (row: FitnessMemberItem | null) => {
  switchUser.value = row
  switchQuery.page = 1
  switchVisible.value = true
  loadSwitches()
}

const goStats = (row: FitnessMemberItem) => {
  router.push({path: FITNESS_ADMIN_STATS_PATH, query: {userId: String(row.userId)}})
}

onMounted(loadData)
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.switch-pagination {
  margin-top: 12px;
  justify-content: flex-end;
}
</style>
