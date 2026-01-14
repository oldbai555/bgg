<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query">
        <el-form-item label="任务名称">
          <el-input v-model="query.name" placeholder="任务名称" clearable />
        </el-form-item>
        <el-form-item label="任务类型">
          <el-select
            v-model="query.type"
            placeholder="任务类型"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in taskTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="Number(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="执行类型">
          <el-select
            v-model="query.executionType"
            placeholder="执行类型"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in executionTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="Number(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="任务状态">
          <el-select
            v-model="query.status"
            placeholder="任务状态"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="Number(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="创建用户ID">
          <el-input
            v-model.number="query.userId"
            placeholder="创建用户ID"
            clearable
          />
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            v-model="createdTimeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="x"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">
            查询
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 任务列表 -->
    <el-card>
      <D2Table
        :columns="columns"
        :data="list"
        :total="total"
        :page-size="query.pageSize || 20"
        :current-page="query.page || 1"
        :drawer-columns="drawerColumns"
        :have-edit="false"
        :have-detail="true"
        detail-permission="task:detail"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import D2Table from '@/components/common/D2Table.vue'
import {
  D2TableElemType,
  type DrawerColumn,
  type TableColumn
} from '@/types/table'
import {
  taskList,
  type TaskItem,
  type TaskListReq
} from '@/api/generated/admin'
import {useDictOptions} from '@/composables/useDictOptions'

const query = reactive<TaskListReq>({
  page: 1,
  pageSize: 20,
  name: '',
  type: undefined,
  executionType: undefined,
  status: undefined,
  userId: undefined,
  startTime: undefined,
  endTime: undefined
})

const list = ref<TaskItem[]>([])
const total = ref(0)
const loading = ref(false)

// 创建时间范围（毫秒时间戳，value-format="x"）
const createdTimeRange = ref<[string, string] | []>([])

// 字典选项
const {options: taskTypeOptions} = useDictOptions('task_type')
const {options: executionTypeOptions} = useDictOptions('task_execution_type')
const {options: statusOptions} = useDictOptions('task_status')

// 枚举转描述映射
const taskTypeMap = computed<Record<number, string>>(() => {
  const m: Record<number, string> = {}
  taskTypeOptions.value.forEach((item) => {
    const val = Number(item.value)
    if (!Number.isNaN(val)) {
      m[val] = item.label
    }
  })
  return m
})

const executionTypeMap = computed<Record<number, string>>(() => {
  const m: Record<number, string> = {}
  executionTypeOptions.value.forEach((item) => {
    const val = Number(item.value)
    if (!Number.isNaN(val)) {
      m[val] = item.label
    }
  })
  return m
})

const statusMap = computed<Record<number, string>>(() => {
  const m: Record<number, string> = {}
  statusOptions.value.forEach((item) => {
    const val = Number(item.value)
    if (!Number.isNaN(val)) {
      m[val] = item.label
    }
  })
  return m
})

// 表格列
const columns = computed<TableColumn[]>(() => [
  {prop: 'id', label: 'ID', width: 80},
  {prop: 'name', label: '任务名称', width: 200},
  {
    prop: 'type',
    label: '任务类型',
    width: 120,
    type: D2TableElemType.EnumToDesc,
    enum2StrMap: taskTypeMap.value
  },
  {
    prop: 'executionType',
    label: '执行类型',
    width: 120,
    type: D2TableElemType.EnumToDesc,
    enum2StrMap: executionTypeMap.value
  },
  {
    prop: 'status',
    label: '任务状态',
    width: 120,
    type: D2TableElemType.EnumToDesc,
    enum2StrMap: statusMap.value
  },
  {prop: 'userId', label: '创建用户ID', width: 120},
  {
    prop: 'scheduledAt',
    label: '计划执行时间',
    width: 180,
    type: D2TableElemType.ConvertTime
  },
  {
    prop: 'startedAt',
    label: '开始时间',
    width: 180,
    type: D2TableElemType.ConvertTime
  },
  {
    prop: 'finishedAt',
    label: '完成时间',
    width: 180,
    type: D2TableElemType.ConvertTime
  },
  {
    prop: 'createdAt',
    label: '创建时间',
    width: 180,
    type: D2TableElemType.ConvertTime
  }
])

// 详情抽屉列
const drawerColumns = computed<DrawerColumn[]>(() => [
  {prop: 'id', label: 'ID', type: D2TableElemType.Tag, disabled: true},
  {prop: 'name', label: '任务名称', type: D2TableElemType.Tag, disabled: true},
  {
    prop: 'type',
    label: '任务类型',
    type: D2TableElemType.EnumToDesc,
    enum2StrMap: taskTypeMap.value,
    disabled: true
  },
  {
    prop: 'executionType',
    label: '执行类型',
    type: D2TableElemType.EnumToDesc,
    enum2StrMap: executionTypeMap.value,
    disabled: true
  },
  {
    prop: 'status',
    label: '任务状态',
    type: D2TableElemType.EnumToDesc,
    enum2StrMap: statusMap.value,
    disabled: true
  },
  {prop: 'userId', label: '创建用户ID', type: D2TableElemType.Tag, disabled: true},
  {
    prop: 'scheduledAt',
    label: '计划执行时间',
    type: D2TableElemType.ConvertTime,
    disabled: true
  },
  {
    prop: 'startedAt',
    label: '开始时间',
    type: D2TableElemType.ConvertTime,
    disabled: true
  },
  {
    prop: 'finishedAt',
    label: '完成时间',
    type: D2TableElemType.ConvertTime,
    disabled: true
  },
  {
    prop: 'createdAt',
    label: '创建时间',
    type: D2TableElemType.ConvertTime,
    disabled: true
  },
  {
    prop: 'params',
    label: '任务参数(JSON)',
    type: D2TableElemType.Textarea,
    disabled: true
  },
  {
    prop: 'result',
    label: '任务结果(JSON)',
    type: D2TableElemType.Textarea,
    disabled: true
  },
  {
    prop: 'fileUrl',
    label: '导出文件',
    type: D2TableElemType.DownloadLink
  },
  {
    prop: 'errorMessage',
    label: '错误信息',
    type: D2TableElemType.Textarea,
    disabled: true
  }
])

const loadData = async () => {
  loading.value = true
  try {
    const req: TaskListReq = {
      page: query.page || 1,
      pageSize: query.pageSize || 20,
      name: query.name || undefined,
      type: query.type && query.type > 0 ? query.type : undefined,
      executionType:
        query.executionType && query.executionType > 0
          ? query.executionType
          : undefined,
      status: query.status && query.status > 0 ? query.status : undefined,
      userId: query.userId || undefined,
      startTime: query.startTime,
      endTime: query.endTime
    }
    const resp = await taskList(req)
    const rawList = resp.list || []
    // 解析任务结果中的 fileUrl，便于在详情里直接提供下载按钮
    list.value = rawList.map((item) => {
      let fileUrl = ''
      if (item.result) {
        try {
          const parsed = JSON.parse(item.result as unknown as string)
          if (parsed && typeof parsed.fileUrl === 'string') {
            fileUrl = parsed.fileUrl
          }
        } catch {
          // ignore JSON parse error
        }
      }
      return {
        ...item,
        fileUrl
      }
    }) as unknown as TaskItem[]
    total.value = resp.total || 0
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '查询失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  query.page = 1
  query.pageSize = 20
  query.name = ''
  query.type = undefined
  query.executionType = undefined
  query.status = undefined
  query.userId = undefined
  query.startTime = undefined
  query.endTime = undefined
  createdTimeRange.value = []
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

// 监听时间范围变化，同步到查询参数（秒级时间戳）
watch(
  createdTimeRange,
  (val) => {
    if (Array.isArray(val) && val.length === 2) {
      const [start, end] = val
      query.startTime = start ? Math.floor(Number(start) / 1000) : undefined
      query.endTime = end ? Math.floor(Number(end) / 1000) : undefined
    } else {
      query.startTime = undefined
      query.endTime = undefined
    }
  },
  {immediate: true}
)

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page {
  padding: 16px;
}

.mb-12 {
  margin-bottom: 12px;
}
</style>

