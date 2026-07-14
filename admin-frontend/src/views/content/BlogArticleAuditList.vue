<template>
  <div class="page">
    <!-- 搜索表单（只关心待审核/审核状态） -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query" class="search-form">
        <el-form-item label="标题">
          <el-input v-model="query.title" placeholder="搜索文章标题" clearable />
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select
            v-model="query.auditStatus"
            placeholder="待审核"
            clearable
            style="width: 140px"
          >
            <el-option
              v-for="item in auditStatusOptions"
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

    <!-- 待审核文章列表 -->
    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        :drawer-columns="[]"
        :have-edit="false"
        :have-detail="false"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      >
        <template #cell="{ row, column }">
          <el-tag
            v-if="column.prop === 'status'"
            :type="statusTagType(row.status)"
            size="small"
          >
            {{ getStatusLabel(row.status) }}
          </el-tag>
          <el-tag
            v-else-if="column.prop === 'auditStatus'"
            :type="auditStatusTagType(row.auditStatus)"
            size="small"
          >
            {{ getAuditStatusLabel(row.auditStatus) }}
          </el-tag>
        </template>

        <template #action="{ row }">
          <el-button
            v-if="row.auditStatus === 2"
            type="success"
            link
            size="small"
            @click="handleAudit(row, true)"
          >
            审核通过
          </el-button>
          <el-button
            v-if="row.auditStatus === 2"
            type="danger"
            link
            size="small"
            @click="handleAudit(row, false)"
          >
            驳回
          </el-button>
          <el-button
            v-if="row.status === 4"
            type="warning"
            link
            size="small"
            @click="handleAuditUnpublish(row)"
          >
            审核下架
          </el-button>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {contentApi} from '@/api/content'
import type {BlogArticleItem, BlogArticleListReq} from '@/api/generated/admin'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn} from '@/types/table'
import {useDictOptions} from '@/composables/useDictOptions'

const query = reactive<BlogArticleListReq>({
  page: 1,
  pageSize: 10,
  title: '',
  status: 0,
  auditStatus: 2, // 默认只看待审核
  tagId: 0
})

const list = ref<BlogArticleItem[]>([])
const total = ref(0)
const loading = ref(false)

const {getLabel: getStatusLabel} = useDictOptions('blog_article_status', [])
const {options: auditStatusOptions, getLabel: getAuditStatusLabel} = useDictOptions(
  'blog_article_audit_status',
  []
)

const statusTagType = (status: number) => {
  switch (status) {
    case 1:
      return 'info'
    case 2:
      return 'warning'
    case 3:
      return 'success'
    case 4:
      return 'success'
    case 5:
      return 'danger'
    default:
      return 'info'
  }
}

const auditStatusTagType = (auditStatus: number) => {
  switch (auditStatus) {
    case 1:
      return 'info'
    case 2:
      return 'warning'
    case 3:
      return 'success'
    case 4:
      return 'danger'
    default:
      return 'info'
  }
}

const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'title', label: '标题', minWidth: 220, showOverflowTooltip: true},
  {prop: 'authorName', label: '作者', width: 120},
  {prop: 'status', label: '状态', width: 110},
  {prop: 'auditStatus', label: '审核状态', width: 110},
  {prop: 'createdAt', label: '创建时间', width: 180, type: D2TableElemType.ConvertTime}
])

const buildListReq = (): BlogArticleListReq => {
  return {
    page: query.page,
    pageSize: query.pageSize,
    title: query.title,
    status: query.status || 0,
    auditStatus: query.auditStatus || 0,
    tagId: query.tagId || 0
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const resp = await contentApi.articleList(buildListReq())
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
  query.title = ''
  query.status = 0
  query.auditStatus = 2
  query.tagId = 0
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

const handleAudit = async (row: BlogArticleItem, pass: boolean) => {
  const title = pass ? '确认审核通过该文章？' : '确认驳回该文章？'
  // 审核结果：3=通过，4=驳回（字典：blog_article_audit_status）
  const result = pass ? 3 : 4
  ElMessageBox.prompt(title, '审核', {
    inputPlaceholder: '请输入审核意见（可选）',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
    .then(async ({value}) => {
      await contentApi.articleAudit({
        id: row.id,
        result,
        remark: value || ''
      })
      ElMessage.success(pass ? '审核通过' : '已驳回')
      loadData()
    })
    .catch(() => {})
}

const handleAuditUnpublish = async (row: BlogArticleItem) => {
  ElMessageBox.prompt('确认下架该文章？', '审核下架', {
    inputPlaceholder: '请输入下架原因（可选）',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
    .then(async ({value}) => {
      await contentApi.articleAuditUnpublish({
        id: row.id,
        remark: value || ''
      })
      ElMessage.success('已下架')
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

.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
}
</style>

