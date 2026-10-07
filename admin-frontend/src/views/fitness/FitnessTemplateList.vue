<template>
  <div class="page">
    <el-card>
      <el-form :inline="true" :model="query" @submit.prevent>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="模板名称/编码" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="search">{{ t('common.search') }}</el-button>
          <el-button @click="reset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
      <el-alert
        type="info"
        :closable="false"
        show-icon
        title="7 套种子模板的编码（starter/standard/advanced/evening/busy/plateau/maintain）被切换建议规则引用，可以改名称和内容，不要改编码或删除。"
      />
    </el-card>

    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        :drawer-columns="drawerColumns"
        :drawer-add-columns="drawerAddColumns"
        :have-detail="false"
        hav-custom-btn
        hav-custom-str="周计划"
        :action-column-width="200"
        create-permission="fitness_template:create"
        update-permission="fitness_template:update"
        delete-permission="fitness_template:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
        @onclick-btn-custom="openWeek"
      >
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'status'" :type="row.status === FitnessStatus.Enabled ? 'success' : 'info'">
            {{ row.status === FitnessStatus.Enabled ? t('status.enabled') : t('status.disabled') }}
          </el-tag>
          <span v-else-if="column.prop === 'updatedAt'">{{ formatUnixTime(row.updatedAt) }}</span>
          <span v-else>{{ row[column.prop as keyof FitnessTemplateItem] }}</span>
        </template>
      </D2Table>
    </el-card>

    <FitnessWeekEditor
      v-if="weekTarget"
      v-model="weekVisible"
      :template-id="weekTarget.id"
      :template-name="weekTarget.name"
    />
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import D2Table from '@/components/common/D2Table.vue'
import FitnessWeekEditor from '@/components/fitness/FitnessWeekEditor.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessTemplateItem} from '@/api/generated/admin'
import {FitnessStatus} from '@/constants/fitness'
import {D2TableElemType, type DrawerColumn, type TableColumn} from '@/types/table'
import {formatUnixTime} from '@/utils/date'

const {t} = useI18n()

const query = reactive({page: 1, pageSize: 20, keyword: ''})
const list = ref<FitnessTemplateItem[]>([])
const total = ref(0)
const loading = ref(false)
const weekVisible = ref(false)
const weekTarget = ref<FitnessTemplateItem | null>(null)

const statusOptions = computed(() => [
  {label: t('status.enabled'), value: FitnessStatus.Enabled},
  {label: t('status.disabled'), value: FitnessStatus.Disabled}
])

const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 70},
  {prop: 'code', label: '编码', width: 110},
  {prop: 'name', label: '名称', width: 130},
  {prop: 'level', label: '强度', width: 70},
  {prop: 'stage', label: '阶段', width: 110},
  {prop: 'summary', label: '说明'},
  {prop: 'dailyKcal', label: '日热量', width: 90},
  {prop: 'dailyProtein', label: '蛋白质g', width: 90},
  {prop: 'sort', label: t('common.order'), width: 70},
  {prop: 'status', label: t('common.status'), width: 80},
  {prop: 'updatedAt', label: '更新时间', width: 170}
])

const editableColumns = computed<DrawerColumn[]>(() => [
  {prop: 'name', label: '名称', type: D2TableElemType.EditInput, required: true},
  {prop: 'level', label: '强度(0~9)', type: D2TableElemType.Number},
  {prop: 'stage', label: '阶段', type: D2TableElemType.EditInput},
  {prop: 'summary', label: '说明', type: D2TableElemType.EditTextarea},
  {prop: 'dailyKcal', label: '日热量 kcal', type: D2TableElemType.Number},
  {prop: 'dailyProtein', label: '蛋白质 g', type: D2TableElemType.Number},
  {prop: 'sort', label: t('common.order'), type: D2TableElemType.Number},
  {prop: 'status', label: t('common.status'), type: D2TableElemType.Select, options: statusOptions.value}
])

const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'code', label: '编码', type: D2TableElemType.Tag},
  ...editableColumns.value
])

const drawerAddColumns = computed<DrawerColumn[]>(() => [
  {prop: 'code', label: '编码（英文，创建后不可改）', required: true},
  ...editableColumns.value
])

const loadData = async () => {
  loading.value = true
  try {
    const resp = await fitnessApi.templateList({page: query.page, pageSize: query.pageSize, keyword: query.keyword || undefined})
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

const toPayload = (row: Record<string, unknown>) => ({
  name: String(row.name || '').trim(),
  level: Number(row.level || 0),
  stage: String(row.stage || ''),
  summary: String(row.summary || ''),
  dailyKcal: Number(row.dailyKcal || 0),
  dailyProtein: Number(row.dailyProtein || 0),
  sort: Number(row.sort || 0),
  status: row.status === undefined ? FitnessStatus.Enabled : Number(row.status)
})

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    await fitnessApi.templateCreate({code: String(row.code || '').trim(), ...toPayload(row)})
    ElMessage.success(t('common.createSuccess'))
    loadData()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  }
}

const handleUpdate = async (row: FitnessTemplateItem) => {
  try {
    await fitnessApi.templateUpdate({id: row.id, ...toPayload(row as unknown as Record<string, unknown>)})
    ElMessage.success(t('common.updateSuccess'))
    loadData()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  }
}

const handleDelete = (_index: number, row: FitnessTemplateItem) => {
  ElMessageBox.confirm(`删除模板「${row.name}」？正在使用该模板的成员需先改派。`, t('common.confirm'), {type: 'warning'})
    .then(async () => {
      try {
        await fitnessApi.templateDelete({id: row.id})
        ElMessage.success(t('common.deleteSuccess'))
        loadData()
      } catch (err: unknown) {
        ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
      }
    })
    .catch(() => {})
}

const openWeek = (_index: number, row: FitnessTemplateItem) => {
  weekTarget.value = row
  weekVisible.value = true
}

onMounted(loadData)
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
