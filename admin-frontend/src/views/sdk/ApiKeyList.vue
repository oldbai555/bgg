<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item :label="t('common.name')">
          <el-input v-model="query.name" :placeholder="t('common.search')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">{{ t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- D2Table：API Key 管理 -->
    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        :drawer-columns="drawerColumns"
        :drawer-add-columns="drawerAddColumns"
        :have-edit="true"
        :have-detail="true"
        create-permission="sdk_key:create"
        update-permission="sdk_key:update"
        delete-permission="sdk_key:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <!-- 自定义列渲染 -->
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'status'" :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? t('status.enabled') : row.status === 2 ? t('status.disabled') : '-' }}
          </el-tag>
          <span v-else-if="column.prop === 'expireAt'">
            <!-- expireAt 为秒级时间戳，0 表示永不过期 -->
            {{ row.expireAt ? formatUnixTime(row.expireAt) : t('common.neverExpire') }}
          </span>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {useI18n} from 'vue-i18n'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {formatUnixTime} from '@/utils/date'
import {sdkApiKeyList, sdkApiKeyCreate, sdkApiKeyUpdate, sdkApiKeyDelete} from '@/api/generated/admin'
import type {SdkApiKeyItem, SdkApiKeyCreateReq, SdkApiKeyUpdateReq} from '@/api/generated/adminComponents'

const {t} = useI18n()

const query = reactive({
  page: 1,
  pageSize: 10,
  name: '',
  status: undefined as number | undefined
})
const list = ref<SdkApiKeyItem[]>([])
const total = ref(0)
const loading = ref(false)

// 表格列配置
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'name', label: t('common.name')},
  {prop: 'apiKey', label: 'API Key'},
  {prop: 'apiSecret', label: 'API Secret'},
  {prop: 'status', label: t('common.status'), width: 100},
  {prop: 'expireAt', label: t('sdk.expireAt'), width: 180},
  {prop: 'createdAt', label: t('common.createdAt'), width: 180, type: D2TableElemType.ConvertTime}
])

// 详情/编辑抽屉列配置
const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'name', label: t('common.name'), type: D2TableElemType.EditInput, required: true},
  {
    prop: 'status',
    label: t('common.status'),
    type: D2TableElemType.Select,
    options: [
      {label: t('status.enabled'), value: 1},
      {label: t('status.disabled'), value: 2}
    ]
  },
  {
    prop: 'expireAt',
    label: t('sdk.expireAt'),
    type: D2TableElemType.Datetime,
    placeholder: t('sdk.expireAtPlaceholder')
  },
  {
    prop: 'ipWhitelist',
    label: t('sdk.ipWhitelist'),
    type: D2TableElemType.EditTextarea
  },
  {
    prop: 'remark',
    label: t('common.remark'),
    type: D2TableElemType.EditTextarea
  }
])

// 新增抽屉列配置
const drawerAddColumns = computed<DrawerColumn[]>(() => [
  {prop: 'name', label: t('common.name'), required: true},
  {
    prop: 'status',
    label: t('common.status'),
    type: D2TableElemType.Select,
    options: [
      {label: t('status.enabled'), value: 1},
      {label: t('status.disabled'), value: 2}
    ]
  },
  {
    prop: 'expireAt',
    label: t('sdk.expireAt'),
    type: D2TableElemType.Datetime,
    placeholder: t('sdk.expireAtPlaceholder')
  },
  {
    prop: 'ipWhitelist',
    label: t('sdk.ipWhitelist'),
    type: D2TableElemType.EditTextarea
  },
  {
    prop: 'remark',
    label: t('common.remark'),
    type: D2TableElemType.EditTextarea
  }
])

const loadData = async () => {
  loading.value = true
  try {
    const resp = await sdkApiKeyList({
      page: query.page,
      pageSize: query.pageSize,
      name: query.name || undefined,
      status: query.status || 0 // 0表示不按状态过滤
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

const handleReset = () => {
  query.page = 1
  query.pageSize = 10
  query.name = ''
  query.status = undefined
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

const handleUpdate = async (row: SdkApiKeyItem) => {
  try {
    const expireAt = row.expireAt ? Number(row.expireAt) : 0
    const payload: SdkApiKeyUpdateReq = {
      id: row.id,
      name: row.name,
      status: row.status,
      expireAt,
      ipWhitelist: row.ipWhitelist,
      remark: row.remark
    }
    await sdkApiKeyUpdate(payload)
    ElMessage.success(t('common.updateSuccess'))
    void loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.updateFail')
    ElMessage.error(message)
  }
}

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    const payload = {
      ...(row as Record<string, unknown>),
      expireAt: row.expireAt ? Number(row.expireAt) : 0
    } as unknown as SdkApiKeyCreateReq
    await sdkApiKeyCreate(payload)
    ElMessage.success(t('common.createSuccess'))
    void loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.createFail')
    ElMessage.error(message)
  }
}

const handleDelete = (index: number, row: SdkApiKeyItem) => {
  ElMessageBox.confirm(t('common.confirmDelete'), t('common.confirm'), {type: 'warning'})
    .then(async () => {
      await sdkApiKeyDelete({id: row.id})
      ElMessage.success(t('common.deleteSuccess'))
      void loadData()
    })
    .catch(() => {})
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
</style>
