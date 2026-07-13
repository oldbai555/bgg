<template>
  <div class="page">
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="Key ID">
          <el-input v-model="query.sdkKeyId" :placeholder="t('sdk.keyIdPlaceholder')" />
        </el-form-item>
        <el-form-item label="API Code">
          <el-input v-model="query.apiCode" placeholder="get:/sdk/file/upload" />
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-select v-model="query.respCode" clearable :placeholder="t('common.all')">
            <el-option label="2xx" :value="200" />
            <el-option label="4xx" :value="400" />
            <el-option label="5xx" :value="500" />
          </el-select>
        </el-form-item>
        <el-form-item label="IP">
          <el-input v-model="query.ip" placeholder="127.0.0.1" />
        </el-form-item>
        <el-form-item :label="t('common.timeRange')">
          <el-date-picker
            v-model="timeRange"
            type="datetimerange"
            range-separator="-"
            :start-placeholder="t('common.startTime')"
            :end-placeholder="t('common.endTime')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSearch">{{ t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
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
        :have-edit="false"
        :have-detail="false"
        :drawer-columns="[]"
        :drawer-add-columns="[]"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      >
        <template #cell="{row, column}">
          <span v-if="column.prop === 'respCode'">
            {{ row.respCode }}
          </span>
        </template>
      </D2Table>
      <div class="mt-12">
        <el-button
          v-permission="'sdk:call_log:export'"
          type="primary"
          :loading="exporting"
          @click="handleExport"
        >
          {{ t('common.export') }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage} from 'element-plus'
import {useI18n} from 'vue-i18n'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn} from '@/types/table'
import {sdkApi} from '@/api/sdk'
import type {SdkCallLogItem, SdkCallLogExportReq} from '@/api/generated/adminComponents'

const {t} = useI18n()

const query = reactive({
  page: 1,
  pageSize: 10,
  sdkKeyId: '',
  apiCode: '',
  respCode: undefined as number | undefined,
  ip: ''
})

const timeRange = ref<[Date, Date] | null>(null)
const list = ref<SdkCallLogItem[]>([])
const total = ref(0)
const loading = ref(false)
const exporting = ref(false)

const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'sdkKeyId', label: 'Key ID', width: 90},
  {prop: 'apiCode', label: 'API Code', width: 220},
  {prop: 'path', label: t('common.path'), width: 220},
  {prop: 'method', label: t('common.method'), width: 90},
  {prop: 'ip', label: 'IP', width: 140},
  {prop: 'respCode', label: t('common.statusCode'), width: 110},
  {prop: 'durationMs', label: t('sdk.durationMs'), width: 110},
  {prop: 'createdAt', label: t('common.createdAt'), width: 180, type: D2TableElemType.ConvertTime}
])

const loadData = async () => {
  loading.value = true
  try {
    const [start, end] = timeRange.value ?? []
    const startTime = start ? Math.floor(start.getTime() / 1000) : undefined
    const endTime = end ? Math.floor(end.getTime() / 1000) : undefined

    const resp = await sdkApi.sdkCallLogList({
      page: query.page,
      pageSize: query.pageSize,
      sdkKeyId: query.sdkKeyId ? Number(query.sdkKeyId) : undefined,
      apiCode: query.apiCode || undefined,
      respCode: query.respCode,
      ip: query.ip || undefined,
      startTime,
      endTime
    })
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.searchFail')
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  query.page = 1
  void loadData()
}

const handleReset = () => {
  query.page = 1
  query.pageSize = 10
  query.sdkKeyId = ''
  query.apiCode = ''
  query.respCode = undefined
  query.ip = ''
  timeRange.value = null
  void loadData()
}

const handlePageChange = (page: number) => {
  query.page = page
  void loadData()
}

const handleSizeChange = (size: number) => {
  query.pageSize = size
  query.page = 1
  void loadData()
}

const handleExport = async () => {
  exporting.value = true
  try {
    const [start, end] = timeRange.value ?? []
    const payload: SdkCallLogExportReq = {
      sdkKeyId: query.sdkKeyId ? Number(query.sdkKeyId) : undefined,
      apiCode: query.apiCode || undefined,
      respCode: query.respCode,
      ip: query.ip || undefined,
      startTime: start ? Math.floor(start.getTime() / 1000) : undefined,
      endTime: end ? Math.floor(end.getTime() / 1000) : undefined
    }
    await sdkApi.sdkCallLogExport(payload)
    ElMessage.success('已创建异步导出任务，请在右下角任务列表查看进度')
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.exportFail')
    ElMessage.error(message)
  } finally {
    exporting.value = false
  }
}

onMounted(() => {
  void loadData()
})
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.mb-12 {
  margin-bottom: 12px;
}
.mt-12 {
  margin-top: 12px;
}
</style>

