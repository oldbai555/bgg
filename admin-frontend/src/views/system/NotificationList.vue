<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="消息来源">
          <el-select
            v-model="query.sourceType"
            placeholder="请选择消息来源"
            clearable
            style="width: 150px"
          >
            <el-option
              v-for="item in sourceTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="已读状态">
          <el-select
            v-model="query.readStatus"
            placeholder="请选择已读状态"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in readStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="Number(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">{{ t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作按钮 -->
    <el-card class="mb-12">
      <el-space>
        <el-button type="success" :loading="readAllLoading" @click="handleReadAll">
          全部已读
        </el-button>
        <el-button type="warning" :loading="clearReadLoading" @click="handleClearRead">
          清除已读
        </el-button>
      </el-space>
    </el-card>

    <!-- D2Table 组件 -->
    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        :drawer-columns="drawerColumns"
        :have-edit="false"
        :have-detail="true"
        delete-permission=""
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
      >
        <!-- 自定义列 -->
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'sourceType'" :type="getSourceTypeTag(row.sourceType)">
            {{ getSourceTypeLabel(row.sourceType) }}
          </el-tag>
          <el-tag v-else-if="column.prop === 'readStatus'" :type="row.readStatus === 2 ? 'success' : 'warning'">
            {{ readStatusOptions.find(opt => Number(opt.value) === row.readStatus)?.label || (row.readStatus === 2 ? '已读' : '未读') }}
          </el-tag>
          <span v-else-if="column.prop === 'readAt'">
            {{ row.readAt ? formatTime(row.readAt) : '-' }}
          </span>
          <el-tooltip
            v-else-if="column.prop === 'content'"
            :content="row.content"
            placement="top"
            :disabled="!row.content || row.content.length <= 50"
          >
            <div class="content-cell">
              {{ row.content }}
            </div>
          </el-tooltip>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {systemApi} from '@/api/system'
import type {NotificationItem} from '@/api/generated/admin'
import {useI18n} from 'vue-i18n'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {useDictOptions} from '@/composables/useDictOptions'

const {t} = useI18n()

const query = reactive({
  page: 1,
  pageSize: 10,
  sourceType: '',
  // 1 = 未读；2 = 已读；undefined 表示不筛选（与 NoticeList 行为一致）
  readStatus: undefined as number | undefined
})
const list = ref<NotificationItem[]>([])
const total = ref(0)
const loading = ref(false)
const readAllLoading = ref(false)
const clearReadLoading = ref(false)

// 消息来源选项（使用字典）
const {options: sourceTypeOptions, getLabel: getSourceTypeLabel} = useDictOptions(
  'notification_source_type',
  [
    {label: '在线聊天', value: 'chat'},
    {label: '系统公告', value: 'notice'},
    {label: '系统通知', value: 'system'}
  ]
)

// 已读状态选项（字典 read_status：1=未读，2=已读；0 由前端表示「全部」）
const {options: readStatusOptions} = useDictOptions('read_status', [
  {label: '未读', value: '1'},
  {label: '已读', value: '2'}
])

// 获取消息来源标签颜色
const getSourceTypeTag = (sourceType: string): string | undefined => {
  const map: Record<string, string> = {
    'chat': 'primary',
    'notice': 'warning',
    'system': 'info'
  }
  return map[sourceType] || undefined
}

// 格式化时间
const formatTime = (timestamp: number): string => {
  if (!timestamp) {
    return '-'
  }
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 表格列配置
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'sourceType', label: '消息来源', width: 120},
  {prop: 'title', label: '消息标题', minWidth: 200},
  {prop: 'content', label: '消息内容', minWidth: 300},
  {prop: 'readStatus', label: '已读状态', width: 100},
  {prop: 'readAt', label: '已读时间', width: 180},
  {prop: 'createdAt', label: '创建时间', width: 180, type: D2TableElemType.ConvertTime}
])

// 详情抽屉列配置
const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'sourceType', label: '消息来源', type: D2TableElemType.Tag},
  {prop: 'title', label: '消息标题'},
  {prop: 'content', label: '消息内容'},
  {prop: 'readStatus', label: '已读状态', type: D2TableElemType.Tag},
  {prop: 'readAt', label: '已读时间'},
  {prop: 'createdAt', label: '创建时间', type: D2TableElemType.ConvertTime}
])

const loadData = async () => {
  loading.value = true
  try {
    const req: Record<string, unknown> = {
      page: query.page,
      pageSize: query.pageSize
    }
    if (query.sourceType) {
      req.sourceType = query.sourceType
    }
    // readStatus：0 不传表示不筛选，其余（1=未读，2=已读）直接透传给后端（由后端映射到 DB 值）
    if ((query.readStatus ?? 0) > 0) {
      req.readStatus = query.readStatus
    }
    const resp = await systemApi.notificationList(req)
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.searchFailed')
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  query.page = 1
  query.pageSize = 10
  query.sourceType = ''
  query.readStatus = undefined
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

const handleDelete = (index: number, row: NotificationItem) => {
  ElMessageBox.confirm('确定要删除该消息通知吗？', '确认删除', {type: 'warning'})
    .then(async () => {
      await systemApi.notificationDelete({id: row.id})
      ElMessage.success('删除成功')
      loadData()
    })
    .catch(() => {})
}

const handleReadAll = async () => {
  readAllLoading.value = true
  try {
    await systemApi.notificationReadAll()
    ElMessage.success('全部已读操作成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '操作失败'
    ElMessage.error(message)
  } finally {
    readAllLoading.value = false
  }
}

const handleClearRead = async () => {
  ElMessageBox.confirm('确定要清除所有已读消息吗？此操作不可恢复。', '确认清除', {type: 'warning'})
    .then(async () => {
      clearReadLoading.value = true
      try {
        await systemApi.notificationClearRead()
        ElMessage.success('清除已读消息成功')
        loadData()
      } catch (err: unknown) {
        const message = err instanceof Error ? err.message : '操作失败'
        ElMessage.error(message)
      } finally {
        clearReadLoading.value = false
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
.mb-12 {
  margin-bottom: 12px;
}
.content-cell {
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: default;
}
</style>

