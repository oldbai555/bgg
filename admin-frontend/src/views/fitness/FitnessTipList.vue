<template>
  <div class="page">
    <el-card>
      <el-form :inline="true" :model="query" @submit.prevent>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="标题/内容" clearable />
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
        title="通用提示（如运动衣物护理、外卖怎么点）在手机端每天的计划页底部统一展示，只在这里维护一份；某一天专属的提示写在模板周计划的「当天提示」里。"
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
        create-permission="fitness_tip:create"
        update-permission="fitness_tip:update"
        delete-permission="fitness_tip:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'status'" :type="row.status === FitnessStatus.Enabled ? 'success' : 'info'">
            {{ row.status === FitnessStatus.Enabled ? t('status.enabled') : t('status.disabled') }}
          </el-tag>
          <div v-else-if="column.prop === 'content'" class="content-cell">{{ row.content }}</div>
          <span v-else-if="column.prop === 'updatedAt'">{{ formatUnixTime(row.updatedAt) }}</span>
          <span v-else>{{ row[column.prop as keyof FitnessTipItem] }}</span>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import D2Table from '@/components/common/D2Table.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessTipItem} from '@/api/generated/admin'
import {FitnessStatus} from '@/constants/fitness'
import {D2TableElemType, type DrawerColumn, type TableColumn} from '@/types/table'
import {formatUnixTime} from '@/utils/date'

const {t} = useI18n()

const query = reactive({page: 1, pageSize: 20, keyword: ''})
const list = ref<FitnessTipItem[]>([])
const total = ref(0)
const loading = ref(false)

const statusOptions = computed(() => [
  {label: t('status.enabled'), value: FitnessStatus.Enabled},
  {label: t('status.disabled'), value: FitnessStatus.Disabled}
])

const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 70},
  {prop: 'code', label: '编码', width: 110},
  {prop: 'title', label: '标题', width: 160},
  {prop: 'content', label: '内容'},
  {prop: 'sort', label: t('common.order'), width: 70},
  {prop: 'status', label: t('common.status'), width: 80},
  {prop: 'updatedAt', label: '更新时间', width: 170}
])

const editableColumns = computed<DrawerColumn[]>(() => [
  {prop: 'title', label: '标题', type: D2TableElemType.EditInput, required: true},
  {prop: 'content', label: '内容', type: D2TableElemType.EditTextarea, required: true},
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
    const resp = await fitnessApi.tipList({page: query.page, pageSize: query.pageSize, keyword: query.keyword || undefined})
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
  title: String(row.title || ''),
  content: String(row.content || ''),
  sort: Number(row.sort || 0),
  status: row.status === undefined ? FitnessStatus.Enabled : Number(row.status)
})

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    await fitnessApi.tipCreate({code: String(row.code || '').trim(), ...toPayload(row)})
    ElMessage.success(t('common.createSuccess'))
    loadData()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  }
}

const handleUpdate = async (row: FitnessTipItem) => {
  try {
    await fitnessApi.tipUpdate({id: row.id, code: row.code, ...toPayload(row as unknown as Record<string, unknown>)})
    ElMessage.success(t('common.updateSuccess'))
    loadData()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  }
}

const handleDelete = (_index: number, row: FitnessTipItem) => {
  ElMessageBox.confirm(t('common.confirmDelete'), t('common.confirm'), {type: 'warning'})
    .then(async () => {
      try {
        await fitnessApi.tipDelete({id: row.id})
        ElMessage.success(t('common.deleteSuccess'))
        loadData()
      } catch (err: unknown) {
        ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
      }
    })
    .catch(() => {})
}

onMounted(loadData)
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.content-cell {
  white-space: pre-line;
  line-height: 1.6;
  max-height: 96px;
  overflow: hidden;
}
</style>
