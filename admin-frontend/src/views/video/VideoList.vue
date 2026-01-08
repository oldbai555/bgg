<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="搜索视频名称、描述" clearable />
        </el-form-item>
        <el-form-item label="来源类型">
          <el-select v-model="query.sourceType" placeholder="全部" clearable style="width: 150px">
            <el-option label="全部" :value="0" />
            <el-option
              v-for="option in sourceTypeOptions"
              :key="option.value"
              :label="option.label"
              :value="Number(option.value)"
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
        create-permission="video:create"
        update-permission="video:update"
        delete-permission="video:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
        @onclick-update-row="handleUpdate"
        @onclick-add-row="handleAdd"
      >
        <!-- 自定义列 -->
        <template #cell="{row, column}">
          <div v-if="column.prop === 'cover'" class="cover-cell">
            <el-image
              v-if="row.cover"
              :src="row.cover"
              :preview-src-list="[row.cover]"
              fit="cover"
              style="width: 80px; height: 60px; border-radius: 4px;"
            />
            <span v-else class="no-cover">无封面</span>
          </div>
          <el-tag
            v-else-if="column.prop === 'sourceType'"
            :type="row.sourceType === 1 ? 'primary' : 'success'"
            size="small"
          >
            {{ getSourceTypeLabel(row.sourceType) }}
          </el-tag>
          <span v-else-if="column.prop === 'duration'">{{ formatDuration(row.duration) }}</span>
        </template>
        <!-- 自定义操作列 -->
        <template #action="{row}">
          <el-button
            type="primary"
            link
            size="small"
            @click="handlePlay(row)"
          >
            播放
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
import {videoList, videoCreate, videoUpdate, videoDelete} from '@/api/generated/admin'
import type {VideoItem, VideoCreateReq, VideoUpdateReq} from '@/api/generated/admin'
import {useI18n} from 'vue-i18n'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {formatUnixTime} from '@/utils/date'
import {useDictOptions} from '@/composables/useDictOptions'

const {t} = useI18n()
const router = useRouter()

const query = reactive({
  page: 1,
  pageSize: 10,
  keyword: '',
  sourceType: 0 as number // 0=全部，1=手动添加，2=采集
})
const list = ref<VideoItem[]>([])
const total = ref(0)
const loading = ref(false)

// 视频来源类型选项（字典 video_source_type：1=手动添加，2=采集）
const {options: sourceTypeOptions, getLabel: getSourceTypeLabel} = useDictOptions('video_source_type', [
  {label: '手动添加', value: '1'},
  {label: '采集', value: '2'}
])

// 格式化时长（秒转时分秒）
const formatDuration = (seconds: number): string => {
  if (!seconds || seconds < 0) {
return '00:00'
}
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  if (h > 0) {
    return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  }
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

// 表格列配置
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'cover', label: '封面', width: 120},
  {prop: 'name', label: '视频名称', minWidth: 200},
  {prop: 'sourceType', label: '来源类型', width: 120},
  {prop: 'duration', label: '时长', width: 100},
  {prop: 'playUrl', label: '播放链接', minWidth: 300, showOverflowTooltip: true},
  {prop: 'createdAt', label: '创建时间', width: 180, type: D2TableElemType.ConvertTime}
])

// 详情/编辑抽屉列配置
const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'name', label: '视频名称', type: D2TableElemType.EditInput, required: true},
  {prop: 'cover', label: '封面URL', type: D2TableElemType.EditInput},
  {prop: 'sourceType', label: '来源类型', type: D2TableElemType.EditSelect, options: sourceTypeOptions},
  {prop: 'duration', label: '时长（秒）', type: D2TableElemType.EditInput},
  {prop: 'playUrl', label: '播放链接', type: D2TableElemType.EditInput, required: true},
  {prop: 'description', label: '描述', type: D2TableElemType.EditTextarea}
])

// 新增抽屉列配置
const drawerAddColumns = computed<DrawerColumn[]>(() => [
  {prop: 'name', label: '视频名称', type: D2TableElemType.EditInput, required: true},
  {prop: 'cover', label: '封面URL', type: D2TableElemType.EditInput},
  {prop: 'sourceType', label: '来源类型', type: D2TableElemType.EditSelect, options: sourceTypeOptions, default: 1},
  {prop: 'duration', label: '时长（秒）', type: D2TableElemType.EditInput},
  {prop: 'playUrl', label: '播放链接', type: D2TableElemType.EditInput, required: true},
  {prop: 'description', label: '描述', type: D2TableElemType.EditTextarea}
])

const loadData = async () => {
  loading.value = true
  try {
    const req: Record<string, unknown> = {
      page: query.page,
      pageSize: query.pageSize
    }
    if (query.keyword) {
      req.keyword = query.keyword
    }
    // sourceType：0 不传表示不筛选，其余（1=手动添加，2=采集）直接透传给后端
    if (query.sourceType > 0) {
      req.type = query.sourceType
    }
    const resp = await videoList(req)
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
  query.sourceType = 0
  loadData()
}

const handlePlay = (row: VideoItem) => {
  // 跳转到视频播放器页面，传递视频URL
  router.push({
    path: '/video/player',
    query: {url: row.playUrl}
  })
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

const handleUpdate = async (row: VideoItem) => {
  try {
    await videoUpdate(row as VideoUpdateReq)
    ElMessage.success('更新成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '更新失败'
    ElMessage.error(message)
  }
}

const handleAdd = async (row: Record<string, unknown>) => {
  try {
    // 确保sourceType字段正确传递（API使用type字段，默认1=手动添加）
    const createData = {
      ...row,
      type: (row.sourceType as number) || 1
    } as VideoCreateReq
    await videoCreate(createData)
    ElMessage.success('新增成功')
    loadData()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '新增失败'
    ElMessage.error(message)
  }
}

const handleDelete = (index: number, row: VideoItem) => {
  ElMessageBox.confirm(t('common.confirmDelete'), t('common.confirm'), {type: 'warning'})
    .then(async () => {
      await videoDelete({id: row.id})
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
.cover-cell {
  display: flex;
  align-items: center;
  justify-content: center;
}
.no-cover {
  color: #999;
  font-size: 12px;
}
</style>
