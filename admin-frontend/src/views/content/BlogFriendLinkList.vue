<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="名称/备注">
          <el-input v-model="query.keyword" placeholder="搜索名称或备注" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="query.status"
            placeholder="全部"
            clearable
            style="width: 140px"
          >
            <el-option label="全部" :value="0" />
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="Number(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- D2Table 组件 -->
    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.size"
        :current-page="query.page"
        :drawer-columns="drawerColumns"
        :drawer-add-columns="drawerAddColumns"
        :have-edit="true"
        :have-detail="true"
        create-permission="blog_friend_link:create"
        update-permission="blog_friend_link:update"
        delete-permission="blog_friend_link:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <!-- 自定义状态列 -->
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'status'" :type="row.status === 1 ? 'success' : 'info'">
            {{ getStatusLabel(row.status as number) }}
          </el-tag>
          <el-link
            v-else-if="column.prop === 'url'"
            :href="row.url"
            target="_blank"
            type="primary"
            :underline="false"
          >
            {{ row.url }}
          </el-link>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {contentApi} from '@/api/content'
import type {
  BlogFriendLinkItem,
  BlogFriendLinkCreateReq,
  BlogFriendLinkUpdateReq
} from '@/api/generated/admin'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {useDictOptions} from '@/composables/useDictOptions'

// 状态选项（从字典获取）
const {options: statusOptions, getLabel: getStatusLabel} = useDictOptions('blog_friend_link_status', [
  {label: '启用', value: 1},
  {label: '禁用', value: 2}
])

const query = reactive({
  page: 1,
  size: 10,
  keyword: '',
  status: 0
})
const list = ref<BlogFriendLinkItem[]>([])
const total = ref(0)
const loading = ref(false)

// 表格列配置
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'name', label: '名称', width: 150},
  {prop: 'url', label: '链接地址'},
  {prop: 'remark', label: '备注'},
  {prop: 'orderNum', label: '排序', width: 100},
  {prop: 'status', label: '状态', width: 100},
  {prop: 'createdAt', label: '创建时间', width: 180}
])

// 详情/编辑抽屉列配置
const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'name', label: '名称', type: D2TableElemType.EditInput, required: true},
  {prop: 'url', label: '链接地址', type: D2TableElemType.EditInput, required: true},
  {prop: 'remark', label: '备注', type: D2TableElemType.EditInput},
  {
    prop: 'orderNum',
    label: '排序值',
    type: D2TableElemType.Number,
    required: false
  },
  {
    prop: 'status',
    label: '状态',
    type: D2TableElemType.Select,
    required: true,
    options: statusOptions.value
  }
])

// 新增抽屉列配置
const drawerAddColumns = computed<DrawerColumn[]>(() => [
  {prop: 'name', label: '名称', type: D2TableElemType.EditInput, required: true},
  {prop: 'url', label: '链接地址', type: D2TableElemType.EditInput, required: true},
  {prop: 'remark', label: '备注', type: D2TableElemType.EditInput},
  {
    prop: 'orderNum',
    label: '排序值',
    type: D2TableElemType.Number,
    required: false
  },
  {
    prop: 'status',
    label: '状态',
    type: D2TableElemType.Select,
    required: true,
    options: statusOptions.value
  }
])

const loadData = async () => {
  loading.value = true
  try {
    const resp = await contentApi.friendLinkList({...query})
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '查询失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  query.page = 1
  query.size = 10
  query.keyword = ''
  query.status = 0
  loadData()
}

const handlePageChange = (page: number) => {
  query.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  query.size = size
  query.page = 1
  loadData()
}

const handleUpdate = async (row: BlogFriendLinkItem) => {
  try {
    await contentApi.friendLinkUpdate(row as BlogFriendLinkUpdateReq)
    ElMessage.success('更新成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '更新失败'
    ElMessage.error(message)
  }
}

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    await contentApi.friendLinkCreate(row as unknown as BlogFriendLinkCreateReq)
    ElMessage.success('新增成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '新增失败'
    ElMessage.error(message)
  }
}

const handleDelete = (index: number, row: BlogFriendLinkItem) => {
  ElMessageBox.confirm('确定要删除这条友情链接吗？', '确认删除', {type: 'warning'})
    .then(async () => {
      await contentApi.friendLinkDelete({id: row.id})
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
}
.mb-12 {
  margin-bottom: 12px;
}
</style>
