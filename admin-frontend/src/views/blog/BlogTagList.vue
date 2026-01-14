<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="标签名称">
          <el-input v-model="query.name" placeholder="搜索标签" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="query.status"
            placeholder="全部"
            clearable
            style="width: 140px"
          >
            <el-option label="全部" :value="0" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 标签列表 -->
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
        create-permission="blog_tag:create"
        update-permission="blog_tag:update"
        delete-permission="blog_tag:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <template #cell="{ row, column }">
          <el-tag v-if="column.prop === 'status'" :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {blogApi} from '@/api/blog'
import type {
  BlogTagItem,
  BlogTagListReq,
  BlogTagCreateReq,
  BlogTagUpdateReq
} from '@/api/generated/admin'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'

const query = reactive<BlogTagListReq>({
  page: 1,
  pageSize: 10,
  name: '',
  status: 0 // 0=全部，1=启用，2=禁用
})

const list = ref<BlogTagItem[]>([])
const total = ref(0)
const loading = ref(false)

const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'name', label: '标签名称'},
  {prop: 'status', label: '状态', width: 100},
  {prop: 'remark', label: '备注', showOverflowTooltip: true},
  {prop: 'createdAt', label: '创建时间', width: 180, type: D2TableElemType.ConvertTime}
])

const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'name', label: '标签名称', type: D2TableElemType.EditInput, required: true},
  {
    prop: 'status',
    label: '状态',
    type: D2TableElemType.EditSelect,
    options: [
      {label: '启用', value: 1},
      {label: '禁用', value: 2}
    ],
    default: 1
  },
  {prop: 'remark', label: '备注', type: D2TableElemType.EditTextarea}
])

// 新增只需要填写标签名称（其余字段后端默认/可后续编辑）
const drawerAddColumns = computed<DrawerColumn[]>(() => [
  {prop: 'name', label: '标签名称', type: D2TableElemType.EditInput, required: true}
])

const loadData = async () => {
  loading.value = true
  try {
    const resp = await blogApi.tagList({...query})
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  query.page = 1
  query.pageSize = 10
  query.name = ''
  query.status = 0
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

const handleUpdate = async (row: BlogTagItem) => {
  try {
    const req: BlogTagUpdateReq = {
      id: row.id,
      name: row.name,
      status: row.status,
      remark: row.remark || ''
    }
    await blogApi.tagUpdate(req)
    ElMessage.success('更新成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '更新失败'
    ElMessage.error(message)
  }
}

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    const req = {
      name: row.name as string
    } as BlogTagCreateReq
    await blogApi.tagCreate(req)
    ElMessage.success('新增成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '新增失败'
    ElMessage.error(message)
  }
}

const handleDelete = (_index: number, row: BlogTagItem) => {
  ElMessageBox.confirm('确认删除该标签？', '提示', {type: 'warning'})
    .then(async () => {
      await blogApi.tagDelete({id: row.id})
      ElMessage.success('删除成功')
      loadData()
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
  padding: 16px 24px;
}

.mb-12 {
  margin-bottom: 12px;
}
</style>

