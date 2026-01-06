<template>
  <div class="page">
    <!-- 搜索表单和上传按钮 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item :label="t('common.name')">
          <el-input v-model="query.name" :placeholder="t('common.search')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">{{ t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <el-upload
            v-permission="'file:create'"
            :action="uploadUrl"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            :show-file-list="false"
            style="display: inline-block; margin-left: 10px;"
          >
            <el-button type="success">上传文件</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
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
        :drawer-add-columns="drawerAddColumns"
        :have-edit="true"
        :have-detail="true"
        create-permission="file:create"
        update-permission="file:update"
        delete-permission="file:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <!-- 自定义状态列和操作列 -->
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'status'" :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? t('status.enabled') : t('status.disabled') }}
          </el-tag>
          <div v-else-if="column.prop === 'path'" class="path-cell">
            <span class="path-text">{{ row.path }}</span>
            <el-button
              type="primary"
              link
              size="small"
              class="copy-btn"
              @click="handleCopyPath(row)"
            >
              <el-icon><DocumentCopy /></el-icon>
              复制
            </el-button>
          </div>
        </template>
        <!-- 自定义操作列 -->
        <template #action="{row}">
          <el-button
            v-permission="'file:list'"
            type="primary"
            link
            size="small"
            @click="handleDownload(row)"
          >
            下载
          </el-button>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {DocumentCopy} from '@element-plus/icons-vue'
import {fileList, fileCreate, fileUpdate, fileDelete} from '@/api/generated/admin'
import type {FileItem, FileCreateReq, FileUpdateReq} from '@/api/generated/admin'
import {useI18n} from 'vue-i18n'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {useUserStore} from '@/stores/user'
import {copyToClipboard} from '@/utils/clipboard'
import {useAppConfig} from '@/composables/useAppConfig'

const {t} = useI18n()

const query = reactive({
  page: 1,
  pageSize: 10,
  name: ''
})
const list = ref<FileItem[]>([])
const total = ref(0)
const loading = ref(false)

// 应用配置
const {storageBaseURL, initConfig} = useAppConfig()

// 文件上传配置
const uploadUrl = computed(() => {
  // 开发环境：始终使用 vite 代理路径（避免 CORS）
  if (import.meta.env.DEV) {
    return '/api/v1/files/upload'
  }
  // 生产环境：使用字典配置的 baseURL
  if (storageBaseURL.value) {
    return `${storageBaseURL.value}/api/v1/files/upload`
  }
  // 生产环境默认使用网关路径
  return ''
})

const userStore = useUserStore()
const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${userStore.token}`
  }
})

// 表格列配置
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'name', label: '短链', width: 200},
  {prop: 'originalName', label: '原始文件名', minWidth: 200, showOverflowTooltip: true},
  {prop: 'path', label: '访问路径', minWidth: 250},
  {prop: 'status', label: t('common.status'), width: 100},
  {prop: 'createdAt', label: t('common.createdAt'), width: 180, type: D2TableElemType.ConvertTime}
])

// 详情/编辑抽屉列配置
const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'name', label: '短链', type: D2TableElemType.Tag},
  {prop: 'originalName', label: '原始文件名', type: D2TableElemType.Tag},
  {prop: 'path', label: '访问路径', type: D2TableElemType.Tag},
  {
    prop: 'status',
    label: t('common.status'),
    type: D2TableElemType.Select,
    options: [
      {label: t('status.enabled'), value: 1},
      {label: t('status.disabled'), value: 0}
    ]
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
      {label: t('status.disabled'), value: 0}
    ]
  }
])

const loadData = async () => {
  loading.value = true
  try {
    const resp = await fileList({...query})
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
  query.name = ''
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

const handleUpdate = async (row: FileItem) => {
  try {
    await fileUpdate(row as FileUpdateReq)
    ElMessage.success('更新成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '更新失败'
    ElMessage.error(message)
  }
}

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    await fileCreate(row as FileCreateReq)
    ElMessage.success('新增成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '新增失败'
    ElMessage.error(message)
  }
}

const handleDelete = (index: number, row: FileItem) => {
  ElMessageBox.confirm(t('common.confirmDelete'), t('common.confirm'), {type: 'warning'})
    .then(async () => {
      await fileDelete({id: row.id})
      ElMessage.success(t('common.delete'))
      loadData()
    })
    .catch(() => {})
}

// 文件上传前验证
const beforeUpload = (file: File) => {
  const isValidSize = file.size / 1024 / 1024 < 50 // 50MB
  if (!isValidSize) {
    ElMessage.error('文件大小不能超过 50MB')
    return false
  }
  return true
}

// 文件上传成功
const handleUploadSuccess = (_response: unknown) => {
  ElMessage.success('文件上传成功')
  loadData()
}

// 文件上传失败
const handleUploadError = (error: Error | unknown) => {
  const message = error instanceof Error ? error.message : '未知错误'
  ElMessage.error('文件上传失败：' + message)
}

// 复制路径（baseUrl + path）
const handleCopyPath = async (row: FileItem) => {
  const fullUrl = row.baseUrl ? `${row.baseUrl}${row.path}` : row.path
  const success = await copyToClipboard(fullUrl)
  if (success) {
    ElMessage.success('已复制到剪贴板')
  } else {
    ElMessage.error('复制失败')
  }
}

// 文件下载
const handleDownload = async (row: FileItem) => {
  try {
    // 构建下载URL（与上传URL逻辑保持一致）
    let downloadUrl = ''
    if (import.meta.env.DEV) {
      // 开发环境：使用 vite 代理路径
      downloadUrl = `/api/v1/files/download?id=${row.id}`
    } else {
      // 生产环境：使用字典配置的 baseURL
      if (storageBaseURL.value) {
        downloadUrl = `${storageBaseURL.value}/api/v1/files/download?id=${row.id}`
      } else {
        // 生产环境默认使用网关路径
        downloadUrl = ''
      }
    }

    // 使用 fetch 下载文件（携带认证 token）
    const response = await fetch(downloadUrl, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.msg || errorData.message || '下载失败')
    }

    // 获取文件 Blob
    const blob = await response.blob()

    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = row.originalName || row.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    // 释放 URL 对象
    window.URL.revokeObjectURL(url)
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '未知错误'
    ElMessage.error('下载失败：' + message)
  }
}

onMounted(async () => {
  // 初始化配置
  await initConfig()
  // 加载数据
  loadData()
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
.path-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.path-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.copy-btn {
  flex-shrink: 0;
}
</style>

