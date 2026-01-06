<template>
  <div class="page">
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item :label="t('common.name')">
          <el-input v-model="query.name" :placeholder="t('common.search')" />
        </el-form-item>
        <el-form-item label="API Code">
          <el-input v-model="query.apiCode" placeholder="get:/sdk/file/upload" />
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-select v-model="query.status" clearable :placeholder="t('common.all')">
            <el-option :label="t('status.enabled')" :value="1" />
            <el-option :label="t('status.disabled')" :value="0" />
          </el-select>
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
        :drawer-columns="drawerColumns"
        :drawer-add-columns="drawerAddColumns"
        :have-edit="true"
        :have-detail="true"
        create-permission="sdk:interface:create"
        update-permission="sdk:interface:update"
        delete-permission="sdk:interface:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'status'" :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? t('status.enabled') : t('status.disabled') }}
          </el-tag>
          <span v-else-if="column.prop === 'createdAt'">
            {{ formatDateTime(row.createdAt) }}
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
import {formatDateTime} from '@/utils/date'
import {
  sdkInterfaceList,
  sdkInterfaceCreate,
  sdkInterfaceUpdate,
  sdkInterfaceDelete
} from '@/api/generated/admin'
import type {
  SdkInterfaceItem,
  SdkInterfaceCreateReq,
  SdkInterfaceUpdateReq
} from '@/api/generated/adminComponents'

const {t} = useI18n()

const query = reactive({
  page: 1,
  pageSize: 10,
  name: '',
  apiCode: '',
  status: undefined as number | undefined
})

const list = ref<SdkInterfaceItem[]>([])
const total = ref(0)
const loading = ref(false)

const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'name', label: t('common.name')},
  {prop: 'apiCode', label: 'API Code', width: 220},
  {prop: 'path', label: t('common.path'), width: 220},
  {prop: 'method', label: t('common.method'), width: 90},
  {prop: 'rateLimitDefault', label: t('sdk.rateLimitDefault'), width: 140},
  {prop: 'status', label: t('common.status'), width: 100},
  {prop: 'createdAt', label: t('common.createdAt'), width: 180}
])

const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'name', label: t('common.name'), type: D2TableElemType.EditInput, required: true},
  {prop: 'apiCode', label: 'API Code', type: D2TableElemType.EditInput, required: true},
  {prop: 'path', label: t('common.path'), type: D2TableElemType.EditInput, required: true},
  {
    prop: 'method',
    label: t('common.method'),
    type: D2TableElemType.EditInput,
    placeholder: 'GET/POST/PUT/DELETE'
  },
  {
    prop: 'rateLimitDefault',
    label: t('sdk.rateLimitDefault'),
    type: D2TableElemType.EditInput,
    placeholder: '每分钟默认请求上限'
  },
  {
    prop: 'status',
    label: t('common.status'),
    type: D2TableElemType.Select,
    options: [
      {label: t('status.enabled'), value: 1},
      {label: t('status.disabled'), value: 0}
    ]
  },
  {
    prop: 'remark',
    label: t('common.remark'),
    type: D2TableElemType.EditTextarea
  }
])

const drawerAddColumns = computed<DrawerColumn[]>(() => [
  {prop: 'name', label: t('common.name'), required: true},
  {prop: 'apiCode', label: 'API Code', required: true},
  {prop: 'path', label: t('common.path'), required: true},
  {
    prop: 'method',
    label: t('common.method'),
    type: D2TableElemType.EditInput,
    placeholder: 'GET/POST/PUT/DELETE'
  },
  {
    prop: 'rateLimitDefault',
    label: t('sdk.rateLimitDefault'),
    type: D2TableElemType.EditInput,
    placeholder: '每分钟默认请求上限'
  },
  {
    prop: 'status',
    label: t('common.status'),
    type: D2TableElemType.Select,
    options: [
      {label: t('status.enabled'), value: 1},
      {label: t('status.disabled'), value: 0}
    ]
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
    const resp = await sdkInterfaceList({
      page: query.page,
      pageSize: query.pageSize,
      name: query.name || undefined,
      apiCode: query.apiCode || undefined,
      status: query.status
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
  query.name = ''
  query.apiCode = ''
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

const handleUpdate = async (row: SdkInterfaceItem) => {
  try {
    const payload: SdkInterfaceUpdateReq = {
      id: row.id,
      name: row.name,
      apiCode: row.apiCode,
      path: row.path,
      method: row.method,
      rateLimitDefault: row.rateLimitDefault,
      status: row.status,
      remark: row.remark
    }
    await sdkInterfaceUpdate(payload)
    ElMessage.success(t('common.updateSuccess'))
    void loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.updateFail')
    ElMessage.error(message)
  }
}

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    const payload = row as unknown as SdkInterfaceCreateReq
    await sdkInterfaceCreate(payload)
    ElMessage.success(t('common.createSuccess'))
    void loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : t('common.createFail')
    ElMessage.error(message)
  }
}

const handleDelete = (index: number, row: SdkInterfaceItem) => {
  ElMessageBox.confirm(t('common.confirmDelete'), t('common.confirm'), {type: 'warning'})
    .then(async () => {
      await sdkInterfaceDelete({id: row.id})
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

