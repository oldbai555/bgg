<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="Method">
          <el-input v-model="query.method" placeholder="HTTP Method" clearable />
        </el-form-item>
        <el-form-item label="Path">
          <el-input v-model="query.path" placeholder="API Path" clearable />
        </el-form-item>
        <el-form-item label="Status">
          <el-input v-model.number="query.statusCode" placeholder="HTTP Status Code" clearable />
        </el-form-item>
        <el-form-item label="Slow Flag">
          <el-select
            v-model="query.isSlow"
            placeholder="Slow or not"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in slowStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="Number(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Start Time">
          <el-date-picker
            v-model="query.startTime"
            type="datetime"
            placeholder="Start Time"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            clearable
          />
        </el-form-item>
        <el-form-item label="End Time">
          <el-date-picker
            v-model="query.endTime"
            type="datetime"
            placeholder="End Time"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">Search</el-button>
          <el-button @click="handleReset">Reset</el-button>
          <el-button type="success" @click="handleExport">Export</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- D2Table 组件（只读列表） -->
    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        :drawer-columns="drawerColumns"
        :drawer-add-columns="drawerAddColumns"
        :have-edit="false"
        :have-detail="false"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      >
        <!-- 自定义慢接口标记列 -->
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'isSlow'" :type="row.isSlow === 1 ? 'danger' : 'info'">
            {{ slowStatusOptions.find(opt => Number(opt.value) === row.isSlow)?.label || (row.isSlow === 1 ? 'Slow' : 'Normal') }}
          </el-tag>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage} from 'element-plus'
import {monitoringApi} from '@/api/monitoring'
import type {
  PerformanceLogItem,
  PerformanceLogListReq,
  PerformanceLogExportReq
} from '@/api/generated/admin'
import {useI18n} from 'vue-i18n'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {useDictOptions} from '@/composables/useDictOptions'

const {t} = useI18n()

const query = reactive<PerformanceLogListReq & { page: number; pageSize: number }>({
  page: 1,
  pageSize: 10,
  method: '',
  path: '',
  isSlow: undefined,
  statusCode: undefined,
  startTime: '',
  endTime: ''
})
const list = ref<PerformanceLogItem[]>([])
const total = ref(0)
const loading = ref(false)

// 慢查询状态选项（字典 performance_log_slow_status：1=Slow，2=Normal；0 由前端表示「全部」）
const {options: slowStatusOptions} = useDictOptions('performance_log_slow_status', [
  {label: 'Slow', value: '1'},
  {label: 'Normal', value: '2'}
])

// 表格列配置（只读性能日志字段）
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'method', label: 'Method', width: 90},
  {prop: 'path', label: 'Path', minWidth: 220, showOverflowTooltip: true},
  {prop: 'statusCode', label: 'Status', width: 90},
  {prop: 'duration', label: 'Duration (ms)', width: 120},
  {prop: 'isSlow', label: 'Slow Flag', width: 100},
  {prop: 'username', label: t('common.username'), width: 140},
  {prop: 'ipAddress', label: 'IP', width: 140},
  {prop: 'createdAt', label: t('common.createdAt'), width: 180, type: D2TableElemType.ConvertTime}
])

// 占位：只读模式但 D2Table 要求必传
const drawerColumns: DrawerColumn[] = []
const drawerAddColumns: DrawerColumn[] = []

const loadData = async () => {
  loading.value = true
  try {
    const resp = await monitoringApi.performanceLogList({...query})
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.search')
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  query.page = 1
  query.pageSize = 10
  query.method = ''
  query.path = ''
  query.isSlow = undefined
  query.statusCode = undefined
  query.startTime = ''
  query.endTime = ''
  loadData()
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

const handleExport = async () => {
  try {
    const req: PerformanceLogExportReq = {
      method: query.method || undefined,
      path: query.path || undefined,
      isSlow: query.isSlow,
      statusCode: query.statusCode,
      startTime: query.startTime || undefined,
      endTime: query.endTime || undefined
    }
    await monitoringApi.performanceLogExport(req)
    ElMessage.success('已创建异步导出任务，请在右下角任务列表查看进度')
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.exportFail')
    ElMessage.error(message)
  }
}

onMounted(loadData)
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
</style>

