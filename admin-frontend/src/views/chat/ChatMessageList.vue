<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="发送用户">
          <el-input v-model="query.fromUserName" placeholder="请输入发送用户名" clearable />
        </el-form-item>
        <el-form-item label="聊天ID">
          <el-input-number v-model="query.chatId" :min="0" placeholder="请输入聊天ID" clearable />
        </el-form-item>
        <el-form-item label="消息类型">
          <el-select
            v-model="query.messageType"
            placeholder="请选择消息类型"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in messageTypeOptions"
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
        :page-size="query.pageSize"
        :current-page="query.page"
        :drawer-columns="drawerColumns"
        :drawer-add-columns="[]"
        :have-edit="false"
        :have-detail="true"
        delete-permission="chat_message:delete"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        @onclick-delete="handleDelete"
      >
        <!-- 自定义消息类型列 -->
        <template #cell="{row, column}">
          <el-tag v-if="column.prop === 'messageType'" :type="getMessageTypeTagType(row.messageType as number)">
            {{ getMessageTypeLabel(row.messageType as number) }}
          </el-tag>
          <!-- 消息内容：如果是图片，显示图片预览 -->
          <div v-else-if="column.prop === 'content' && row.messageType === 2" class="message-content-image">
            <el-image
              :src="row.content"
              fit="cover"
              style="width: 100px; height: 100px"
              :preview-src-list="[row.content]"
              preview-teleported
            >
              <template #error>
                <div class="image-slot">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </div>
          <!-- 消息内容：文本消息，显示前50个字符 -->
          <span v-else-if="column.prop === 'content' && row.messageType === 1" class="message-content-text">
            {{ row.content.length > 50 ? row.content.substring(0, 50) + '...' : row.content }}
          </span>
          <!-- 消息内容：文件消息 -->
          <el-link
            v-else-if="column.prop === 'content' && row.messageType === 3"
            :href="row.content"
            target="_blank"
            type="primary"
          >
            查看文件
          </el-link>
          <!-- 发送时间：格式化显示 -->
          <span v-else-if="column.prop === 'createdAt'">
            {{ formatUnixTime(row.createdAt) }}
          </span>
        </template>
      </D2Table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, onMounted, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {Picture} from '@element-plus/icons-vue'
import {chatApi} from '@/api/chat'
import type {ChatMessageItem, ChatMessageListReq} from '@/api/generated/admin'
import D2Table from '@/components/common/D2Table.vue'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {buildFileUrlFromResponse} from '@/utils/file'
import {formatUnixTime} from '@/utils/date'
import {useDictOptions} from '@/composables/useDictOptions'

const query = reactive<ChatMessageListReq & {fromUserName?: string; messageType?: number}>({
  page: 1,
  pageSize: 10,
  chatId: undefined,
  fromUserName: '',
  messageType: undefined
})
const list = ref<ChatMessageItem[]>([])
const total = ref(0)
const loading = ref(false)

// 消息类型选项
const {options: messageTypeOptions, getLabel: getMessageTypeLabelFromDict} = useDictOptions(
  'chat_message_type',
  [
    {label: '文本', value: '1'},
    {label: '图片', value: '2'},
    {label: '文件', value: '3'}
  ]
)

// 表格列配置
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'fromUserName', label: '发送用户', width: 120},
  {prop: 'chatId', label: '聊天ID', width: 100},
  {prop: 'content', label: '消息内容', minWidth: 200},
  {prop: 'messageType', label: '消息类型', width: 100},
  {prop: 'createdAt', label: '发送时间', width: 180}
])

// 详情抽屉列配置（只读）
const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag},
  {prop: 'fromUserName', label: '发送用户', type: D2TableElemType.Tag},
  {prop: 'chatId', label: '聊天ID', type: D2TableElemType.Tag},
  {prop: 'content', label: '消息内容'},
  {prop: 'messageType', label: '消息类型'},
  {prop: 'createdAt', label: '发送时间', type: D2TableElemType.ConvertTime}
])

// 获取消息类型标签（使用字典）
const getMessageTypeLabel = (type: number) => {
  return getMessageTypeLabelFromDict(type) || '未知'
}

// 获取消息类型标签颜色
const getMessageTypeTagType = (type: number): 'success' | 'warning' | 'info' => {
  const typeMap: Record<number, 'success' | 'warning' | 'info'> = {
    1: 'success',
    2: 'warning',
    3: 'info'
  }
  return typeMap[type] || 'info'
}

const loadData = async () => {
  loading.value = true
  try {
    // 后端只支持按 chatId 筛选，发送用户名/消息类型在前端过滤
    const req: ChatMessageListReq = {
      page: query.page,
      pageSize: query.pageSize,
      chatId: query.chatId || undefined
    }

    const resp = await chatApi.chatMessageListAdmin(req)
    let filteredList = resp.list || []

    // 前端过滤：根据发送用户名过滤
    if (query.fromUserName) {
      filteredList = filteredList.filter(item =>
        item.fromUserName?.toLowerCase().includes(query.fromUserName!.toLowerCase())
      )
    }
    if (query.messageType !== undefined && query.messageType !== null) {
      filteredList = filteredList.filter(item => item.messageType === query.messageType)
    }

    // 处理图片消息的 URL（如果是相对路径，需要拼接 baseUrl）
    filteredList = filteredList.map(item => {
      if (item.messageType === 2 && item.content && !item.content.startsWith('http')) {
        // 图片消息，如果是相对路径，使用工具函数拼接完整 URL
        // 注意：这里假设 content 存储的是文件路径，实际可能需要从文件表查询
        item.content = buildFileUrlFromResponse({path: item.content})
      }
      return item
    })

    list.value = filteredList
    // 用后端返回的 total（按 chatId 服务端筛选后的总数），而不是前端二次过滤后的页内数量，
    // 否则分页会算错；发送用户名/消息类型是纯前端展示层过滤，不影响分页语义
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
  query.pageSize = 10
  query.chatId = undefined
  query.fromUserName = ''
  query.messageType = undefined
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

const handleDelete = (index: number, row: ChatMessageItem) => {
  ElMessageBox.confirm('确定要删除这条聊天记录吗？删除后无法恢复。', '确认删除', {type: 'warning'})
    .then(async () => {
      try {
        await chatApi.chatMessageDelete({id: row.id})
        ElMessage.success('删除成功')
        loadData()
      } catch (err: unknown) {
        const message = err instanceof Error ? err.message : '删除失败'
        ElMessage.error(message)
      }
    })
    .catch(() => {})
}

onMounted(loadData)
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.mb-12 {
  margin-bottom: 12px;
}

.message-content-image {
  display: inline-block;
}

.message-content-text {
  word-break: break-word;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-placeholder);
}
</style>

