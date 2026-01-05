<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="搜索内容或作者" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select
            v-model="query.type"
            placeholder="请选择类型"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in sentenceTypeOptions"
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
        create-permission="daily_short_sentence:create"
        update-permission="daily_short_sentence:update"
        delete-permission="daily_short_sentence:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <!-- 自定义类型列 -->
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'type'" :type="row.type === 2 ? 'success' : 'info'">
            {{ sentenceTypeOptions.find(opt => Number(opt.value) === row.type)?.label || (row.type === 1 ? '普通' : '文学') }}
          </el-tag>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {dailyShortSentenceList, dailyShortSentenceCreate, dailyShortSentenceUpdate, dailyShortSentenceDelete} from '@/api/generated/admin'
import type {DailyShortSentenceItem, DailyShortSentenceCreateReq, DailyShortSentenceUpdateReq} from '@/api/generated/admin'
import {useI18n} from 'vue-i18n'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {useDictOptions} from '@/composables/useDictOptions'

const {t} = useI18n()

const query = reactive({
  page: 1,
  pageSize: 10,
  keyword: '',
  type: undefined as number | undefined
})
const list = ref<DailyShortSentenceItem[]>([])
const total = ref(0)
const loading = ref(false)

// 短句类型选项
const {options: sentenceTypeOptions} = useDictOptions(
  'daily_short_sentence_type',
  [
    {label: '普通', value: '1'},
    {label: '文学', value: '2'}
  ]
)

// 表格列配置
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'content', label: '短句内容', minWidth: 200},
  {prop: 'type', label: '类型', width: 100},
  {prop: 'literatureAuthor', label: '作者', width: 120},
  {prop: 'createdAt', label: t('common.createdAt'), width: 180}
])

// 详情/编辑抽屉列配置
const drawerColumns = computed((): DrawerColumn[] => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {
    prop: 'type',
    label: '类型',
    type: D2TableElemType.Select,
    options: sentenceTypeOptions.value.map(opt => ({label: opt.label, value: Number(opt.value)}))
  },
  {prop: 'content', label: '短句内容', type: D2TableElemType.EditTextarea, required: true},
  {prop: 'literatureAuthor', label: '作者', type: D2TableElemType.EditInput},
  {prop: 'img', label: '图片URL', type: D2TableElemType.EditInput},
  {prop: 'convertImg', label: '转换图片URL', type: D2TableElemType.EditInput}
])

// 新增抽屉列配置
const drawerAddColumns = computed<DrawerColumn[]>(() => [
  {
    prop: 'type',
    label: '类型',
    type: D2TableElemType.Select,
    options: sentenceTypeOptions.value.map(opt => ({label: opt.label, value: Number(opt.value)}))
  },
  {prop: 'content', label: '短句内容', type: D2TableElemType.EditTextarea, required: true},
  {prop: 'literatureAuthor', label: '作者', type: D2TableElemType.EditInput},
  {prop: 'img', label: '图片URL', type: D2TableElemType.EditInput},
  {prop: 'convertImg', label: '转换图片URL', type: D2TableElemType.EditInput}
])

const loadData = async () => {
  loading.value = true
  try {
    const resp = await dailyShortSentenceList({
      page: query.page,
      pageSize: query.pageSize,
      keyword: query.keyword || undefined,
      type: query.type
    })
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
  query.keyword = ''
  query.type = undefined
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

const handleUpdate = async (row: DailyShortSentenceItem) => {
  try {
    await dailyShortSentenceUpdate({
      id: row.id,
      type: row.type,
      content: row.content,
      literatureAuthor: row.literatureAuthor,
      img: row.img,
      convertImg: row.convertImg
    } as DailyShortSentenceUpdateReq)
    ElMessage.success('更新成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '更新失败'
    ElMessage.error(message)
  }
}

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    await dailyShortSentenceCreate({
      type: (row.type as number) || 1,
      content: row.content as string,
      literatureAuthor: row.literatureAuthor as string,
      img: row.img as string,
      convertImg: row.convertImg as string
    } as DailyShortSentenceCreateReq)
    ElMessage.success('新增成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '新增失败'
    ElMessage.error(message)
  }
}

const handleDelete = (index: number, row: DailyShortSentenceItem) => {
  ElMessageBox.confirm(t('common.confirmDelete'), t('common.confirm'), {type: 'warning'})
    .then(async () => {
      await dailyShortSentenceDelete({id: row.id})
      ElMessage.success(t('common.delete'))
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
