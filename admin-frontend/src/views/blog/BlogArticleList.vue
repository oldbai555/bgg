<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query" class="search-form">
        <el-form-item label="标题">
          <el-input v-model="query.title" placeholder="搜索文章标题" clearable />
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
        <el-form-item label="审核状态">
          <el-select
            v-model="query.auditStatus"
            placeholder="全部"
            clearable
            style="width: 140px"
          >
            <el-option label="全部" :value="0" />
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

    <!-- 文章列表 -->
    <el-card>
      <div class="toolbar">
        <el-button v-permission="'blog_article:create'" type="primary" @click="goCreate">
          新增文章
        </el-button>
      </div>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        :have-edit="false"
        :have-detail="false"
        :drawer-columns="[]"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
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
          <el-tag
            v-else-if="column.prop === 'isTop'"
            :type="getIsTopValue(row) === 1 ? 'warning' : 'info'"
            size="small"
          >
            {{ getIsTopValue(row) === 1 ? '置顶' : '普通' }}
          </el-tag>
          <div v-else-if="column.prop === 'tagNames'" class="tag-names">
            <el-tag
              v-for="tag in row.tagNames || []"
              :key="tag"
              size="small"
              class="mr-4"
            >
              {{ tag }}
            </el-tag>
          </div>
        </template>

        <template #action="{ row }">
          <el-button
            v-permission="'blog_article:update'"
            type="primary"
            link
            size="small"
            @click="goEdit(row.id)"
          >
            编辑
          </el-button>
          <el-button
            v-if="canSubmit(row)"
            type="primary"
            link
            size="small"
            @click="handleSubmitAudit(row)"
          >
            提交审核
          </el-button>
          <el-button
            v-if="canPublish(row)"
            type="success"
            link
            size="small"
            @click="handlePublish(row)"
          >
            上架
          </el-button>
          <el-button
            v-if="canUnpublish(row)"
            type="warning"
            link
            size="small"
            @click="handleUnpublish(row)"
          >
            下架
          </el-button>
          <el-button
            v-if="canTop(row)"
            v-permission="'blog_article:top'"
            type="warning"
            link
            size="small"
            @click="handleTop(row)"
          >
            置顶
          </el-button>
          <el-button
            v-if="canUntop(row)"
            v-permission="'blog_article:untop'"
            type="info"
            link
            size="small"
            @click="handleUntop(row)"
          >
            取消置顶
          </el-button>
          <el-button
            v-permission="'blog_article:delete'"
            type="danger"
            link
            size="small"
            @click="handleDelete(0, row)"
          >
            删除
          </el-button>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {useRouter} from 'vue-router'
import {blogApi} from '@/api/blog'
import type {
  BlogArticleItem,
  BlogArticleListReq
} from '@/api/generated/admin'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn} from '@/types/table'
import {useDictOptions} from '@/composables/useDictOptions'

const router = useRouter()

const query = reactive<BlogArticleListReq>({
  page: 1,
  pageSize: 10,
  title: '',
  status: 0, // 字典 blog_article_status，0=全部
  auditStatus: 0, // 字典 blog_article_audit_status，0=全部
  tagId: 0
})

const list = ref<BlogArticleItem[]>([])
const total = ref(0)
const loading = ref(false)

const {options: statusOptions, getLabel: getStatusLabel} = useDictOptions('blog_article_status', [])
const {options: auditStatusOptions, getLabel: getAuditStatusLabel} = useDictOptions(
  'blog_article_audit_status',
  []
)

const statusTagType = (status: number) => {
  // 1=草稿，2=待审核，3=审核通过-未上架，4=上架，5=下架
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
  // 1=未提交，2=待审核，3=审核通过，4=审核驳回
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
  {prop: 'isTop', label: '置顶', width: 80},
  {prop: 'tagNames', label: '标签', minWidth: 180},
  {prop: 'publishTime', label: '上架时间', width: 180, type: D2TableElemType.ConvertTime},
  {prop: 'createdAt', label: '创建时间', width: 180, type: D2TableElemType.ConvertTime}
])

const buildListReq = (): BlogArticleListReq => {
  const req: BlogArticleListReq = {
    page: query.page,
    pageSize: query.pageSize,
    title: query.title,
    status: query.status || 0,
    auditStatus: query.auditStatus || 0,
    tagId: query.tagId || 0
  }
  return req
}

const loadData = async () => {
  loading.value = true
  try {
    const resp = await blogApi.articleList(buildListReq())
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
  query.auditStatus = 0
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

const goCreate = () => {
  router.push('/blog/article/edit')
}

const goEdit = (id: number) => {
  router.push(`/blog/article/edit/${id}`)
}

const handleDelete = (_index: number, row: BlogArticleItem) => {
  ElMessageBox.confirm('确认删除该文章？', '提示', {type: 'warning'})
    .then(async () => {
      await blogApi.articleDelete({id: row.id})
      ElMessage.success('删除成功')
      loadData()
    })
    .catch(() => {})
}

const canSubmit = (row: BlogArticleItem) => {
  // 草稿或审核驳回可以提交审核
  return row.status === 1 || row.auditStatus === 4
}

const canPublish = (row: BlogArticleItem) => {
  // 审核通过且未上架可以上架
  return row.auditStatus === 3 && row.status !== 4
}

const canUnpublish = (row: BlogArticleItem) => {
  // 已上架可以下架
  return row.status === 4
}

// 获取 isTop 值（兼容不同数据类型）
const getIsTopValue = (row: BlogArticleItem): number => {
  const article = row as BlogArticleItem & {isTop?: number | string}
  const isTop = article.isTop
  if (isTop === 1 || isTop === '1' || String(isTop) === '1') {
    return 1
  }
  return 0
}

const canTop = (row: BlogArticleItem) => {
  // 已上架且未置顶可以置顶
  return row.status === 4 && getIsTopValue(row) !== 1
}

const canUntop = (row: BlogArticleItem) => {
  // 已置顶可以取消置顶
  return getIsTopValue(row) === 1
}

const handleSubmitAudit = async (row: BlogArticleItem) => {
  try {
    await blogApi.articleSubmit({id: row.id})
    ElMessage.success('已提交审核')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '提交失败'
    ElMessage.error(message)
  }
}

const handlePublish = async (row: BlogArticleItem) => {
  try {
    await blogApi.articlePublish({id: row.id})
    ElMessage.success('已上架')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '上架失败'
    ElMessage.error(message)
  }
}

const handleUnpublish = async (row: BlogArticleItem) => {
  try {
    await blogApi.articleUnpublish({id: row.id})
    ElMessage.success('已下架')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '下架失败'
    ElMessage.error(message)
  }
}

const handleTop = async (row: BlogArticleItem) => {
  try {
    await blogApi.articleTop({id: row.id})
    ElMessage.success('已置顶')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '置顶失败'
    ElMessage.error(message)
  }
}

const handleUntop = async (row: BlogArticleItem) => {
  try {
    await blogApi.articleUntop({id: row.id})
    ElMessage.success('已取消置顶')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '取消置顶失败'
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

.tag-names {
  display: flex;
  flex-wrap: wrap;
}

.mr-4 {
  margin-right: 4px;
  margin-bottom: 4px;
}
</style>

