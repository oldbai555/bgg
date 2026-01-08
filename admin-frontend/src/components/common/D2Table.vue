<template>
  <div class="d2-table">
    <!-- 新增按钮 -->
    <div
      v-if="drawerAddColumns && drawerAddColumns.length > 0 && canCreate"
      class="d2-table__toolbar"
    >
      <el-button type="primary" @click="showAdd">
        {{ t('common.create') }}
      </el-button>
    </div>

    <!-- 表格 -->
    <el-table
      :data="displayedData"
      :max-height="maxHeight"
      border
      style="width: 100%"
    >
      <el-table-column
        v-for="(column, index) in columns"
        :key="index"
        :prop="column.prop"
        :label="column.label"
        :width="column.width || undefined"
        :fixed="column.fixed"
      >
        <template #default="scope">
          <slot
            name="cell"
            :row="scope.row"
            :column="column"
            :index="index"
          >
            <!-- 时间戳转换 -->
            <el-tag v-if="column.type === D2TableElemType.ConvertTime">
              {{ formatUnixTime(scope.row[column.prop]) }}
            </el-tag>
            <!-- 标签显示 -->
            <el-tag v-else-if="column.type === D2TableElemType.Tag">
              {{ scope.row[column.prop] }}
            </el-tag>
            <!-- 枚举转描述 -->
            <el-tag v-else-if="column.type === D2TableElemType.EnumToDesc">
              {{ column.enum2StrMap?.[scope.row[column.prop]] || scope.row[column.prop] }}
            </el-tag>
            <!-- 下载链接（带 baseUrl） -->
            <el-link
              v-else-if="column.type === D2TableElemType.DownloadWithSortUrl"
              type="primary"
              :href="`${baseUrl}/${scope.row[column.prop]}`"
              target="_blank"
            >
              {{ t('common.download') }}
            </el-link>
            <!-- 复制链接 -->
            <el-button
              v-else-if="column.type === D2TableElemType.CopyUrl"
              type="primary"
              link
              @click="handleCopyUrl(`${baseUrl}/${scope.row[column.prop]}`)"
            >
              {{ t('common.copy') }}
            </el-button>
            <!-- 跳转链接 -->
            <el-link
              v-else-if="column.type === D2TableElemType.LinkJump"
              type="primary"
              :href="scope.row[column.prop]"
              target="_blank"
            >
              {{ t('common.view') }}
            </el-link>
            <!-- 图片（带 baseUrl） -->
            <div v-else-if="column.type === D2TableElemType.ImageWithSortUrl" class="d2-table__image">
              <el-image
                style="width: 100px; height: 100px"
                :src="`${baseUrl}/${scope.row[column.prop]}`"
                fit="cover"
              >
                <template #error>
                  <div class="image-slot">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <!-- 图片 -->
            <div v-else-if="column.type === D2TableElemType.Image" class="d2-table__image">
              <el-image
                style="width: 100px; height: 100px"
                :src="scope.row[column.prop]"
                fit="cover"
              >
                <template #error>
                  <div class="image-slot">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <!-- 默认文本 -->
            <span v-else>
              {{ scope.row[column.prop] }}
            </span>
          </slot>
        </template>
      </el-table-column>

      <!-- 操作列 -->
      <el-table-column fixed="right" :label="t('common.actions')" :width="actionColumnWidth">
        <template #default="scope">
          <slot name="action" :row="scope.row" :index="scope.$index">
            <el-button
              v-if="haveDetail && canDetail"
              size="small"
              type="primary"
              link
              @click="handleEdit(scope.$index, scope.row, false)"
            >
              {{ t('common.view') }}
            </el-button>
            <el-button
              v-if="havCustomBtn && canCustom"
              size="small"
              type="primary"
              link
              @click="handleBtnCustom(scope.$index, scope.row)"
            >
              {{ havCustomStr }}
            </el-button>
            <el-button
              v-if="haveEdit && canUpdate"
              size="small"
              type="warning"
              link
              @click="handleEdit(scope.$index, scope.row, true)"
            >
              {{ t('common.edit') }}
            </el-button>
            <el-button
              v-if="canDelete"
              size="small"
              type="danger"
              link
              @click="handleDelete(scope.$index, scope.row)"
            >
              {{ t('common.delete') }}
            </el-button>
          </slot>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="d2-table__pagination">
      <el-pagination
        v-model:current-page="currentPageModel"
        v-model:page-size="pageSizeModel"
        :page-sizes="pageSizes"
        :total="total"
        :layout="paginationLayout"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 详情/编辑抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      :title="isEdit ? t('common.edit') : t('common.detail')"
      direction="rtl"
      :size="drawerWidth"
    >
      <el-form :model="drawerRow" label-width="120px">
        <el-form-item
          v-for="(column, index) in drawerColumns"
          :key="index"
          :label="column.label"
          :required="column.required"
          v-show="!isEdit || !column.disabled"
        >
          <!-- 下载链接 -->
          <el-link
            v-if="column.type === D2TableElemType.DownloadWithSortUrl"
            type="primary"
            :href="`${baseUrl}/${drawerRow[column.prop]}`"
            target="_blank"
          >
            {{ t('common.download') }}
          </el-link>
          <!-- 复制链接 -->
          <el-button
            v-else-if="column.type === D2TableElemType.CopyUrl"
            type="primary"
            link
            @click="handleCopyUrl(`${baseUrl}/${drawerRow[column.prop]}`)"
          >
            {{ t('common.copy') }}
          </el-button>
          <!-- 图片（带 baseUrl） -->
          <div v-else-if="column.type === D2TableElemType.ImageWithSortUrl" class="d2-table__image">
            <el-image
              style="width: 100px; height: 100px"
              :src="`${baseUrl}/${drawerRow[column.prop]}`"
              fit="cover"
            >
              <template #error>
                <div class="image-slot">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </div>
          <!-- 下拉选择 -->
          <el-select
            v-else-if="column.type === D2TableElemType.Select"
            v-model="drawerRow[column.prop]"
            :disabled="!isEdit || column.disabled"
            style="width: 360px"
          >
            <el-option
              v-for="item in column.options"
              :key="String(item.value)"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <!-- 可编辑输入框 -->
          <el-input
            v-else-if="column.type === D2TableElemType.EditInput"
            v-model="drawerRow[column.prop]"
            :disabled="!isEdit || column.disabled"
            style="width: 360px"
          />
          <!-- 可编辑文本域（详情模式下显示为只读文本域，支持滚动） -->
          <el-input
            v-else-if="column.type === D2TableElemType.EditTextarea"
            v-model="drawerRow[column.prop]"
            type="textarea"
            :rows="!isEdit ? 8 : 4"
            :disabled="!isEdit"
            style="width: 360px"
          />
          <el-input-number
            v-else-if="column.type === D2TableElemType.Number"
            v-model="drawerRow[column.prop]"
            :min="0"
            :disabled="!isEdit"
            style="width: 360px"
          />
          <el-date-picker
            v-else-if="column.type === D2TableElemType.Datetime"
            v-model="drawerRow[column.prop]"
            type="datetime"
            value-format="X"
            :placeholder="column.placeholder"
            :disabled="!isEdit"
            style="width: 360px"
          />
          <!-- 只读文本域（用于详情显示） -->
          <el-input
            v-else-if="column.type === D2TableElemType.Textarea"
            v-model="drawerRow[column.prop]"
            type="textarea"
            :rows="4"
            disabled
            style="width: 360px"
          />
          <!-- 字节转MB -->
          <el-tag v-else-if="column.type === D2TableElemType.Byte2MB">
            {{ formatBytes(drawerRow[column.prop]) }}
          </el-tag>
          <!-- 枚举转描述 -->
          <el-tag v-else-if="column.type === D2TableElemType.EnumToDesc">
            {{ column.enum2StrMap?.[drawerRow[column.prop]] || drawerRow[column.prop] }}
          </el-tag>
          <!-- 时间戳转换 -->
          <el-tag v-else-if="column.type === D2TableElemType.ConvertTime">
            {{ formatUnixTime(drawerRow[column.prop]) }}
          </el-tag>
          <!-- 图片 -->
          <div v-else-if="column.type === D2TableElemType.Image" class="d2-table__image">
            <ImageUpload
              v-if="isEdit && drawerRow[column.prop] !== undefined"
              v-model="drawerRow[column.prop]"
            />
            <el-image
              v-else-if="!isEdit && drawerRow[column.prop]"
              style="width: 200px; height: 200px; max-width: 100%"
              :src="drawerRow[column.prop]"
              fit="contain"
              :preview-src-list="[drawerRow[column.prop]].filter(Boolean)"
            >
              <template #error>
                <div class="image-slot">
                  <el-icon><Picture /></el-icon>
                  <span style="font-size: 12px; margin-top: 4px">加载失败</span>
                </div>
              </template>
            </el-image>
            <span v-else-if="!isEdit && !drawerRow[column.prop]" style="color: #909399">暂无图片</span>
          </div>
          <!-- 默认标签 -->
          <el-tag v-else>
            {{ drawerRow[column.prop] }}
          </el-tag>
        </el-form-item>

        <el-form-item v-if="haveEdit && isEdit && canUpdate">
          <el-button type="primary" @click="updateItem">
            {{ t('common.save') }}
          </el-button>
          <el-button @click="cancelEdit">
            {{ t('common.cancel') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- 新增抽屉 -->
    <el-drawer
      v-model="drawerVisibleAdd"
      :title="t('common.create')"
      direction="rtl"
      :size="drawerWidth"
    >
      <el-form :model="drawerAddRow" label-width="120px">
        <el-form-item
          v-for="(column, index) in drawerAddColumns"
          :key="index"
          :label="column.label"
          :required="column.required"
        >
          <!-- 下拉选择 -->
          <el-select
            v-if="column.type === D2TableElemType.Select"
            v-model="drawerAddRow[column.prop]"
            style="width: 360px"
          >
            <el-option
              v-for="item in column.options"
              :key="String(item.value)"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <!-- 图片 -->
          <div v-else-if="column.type === D2TableElemType.Image" class="d2-table__image">
            <ImageUpload v-model="drawerAddRow[column.prop]" />
          </div>
          <!-- 可编辑文本域 -->
          <el-input
            v-else-if="column.type === D2TableElemType.EditTextarea"
            v-model="drawerAddRow[column.prop]"
            type="textarea"
            :rows="4"
            style="width: 360px"
          />
          <el-input-number
            v-else-if="column.type === D2TableElemType.Number"
            v-model="drawerAddRow[column.prop]"
            :min="0"
            style="width: 360px"
          />
          <el-date-picker
            v-else-if="column.type === D2TableElemType.Datetime"
            v-model="drawerAddRow[column.prop]"
            type="datetime"
            value-format="X"
            :placeholder="column.placeholder"
            style="width: 360px"
          />
          <!-- 默认输入框 -->
          <el-input
            v-else
            v-model="drawerAddRow[column.prop]"
            style="width: 360px"
          />
        </el-form-item>

        <el-form-item>
          <el-button v-if="canCreate" type="primary" @click="handleAdd">
            {{ t('common.create') }}
          </el-button>
          <el-button @click="cancelAdd">
            {{ t('common.cancel') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import {computed, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {Picture} from '@element-plus/icons-vue'
import {useI18n} from 'vue-i18n'
import {usePermission} from '@/hooks/usePermission'
import {D2TableElemType, type TableColumn, type DrawerColumn} from '@/types/table'
import {formatUnixTime} from '@/utils/date'
import {copyToClipboard} from '@/utils/clipboard'
import ImageUpload from './ImageUpload.vue'

const {t} = useI18n()
const {hasPermission} = usePermission()

// Props
interface Props {
  /** 表格列配置 */
  columns: TableColumn[];
  /** 表格数据 */
  data: Record<string, unknown>[]
  /** 总条数 */
  total: number;
  /** 每页显示条数 */
  pageSize?: number;
  /** 当前页码 */
  currentPage?: number;
  /** 基础URL（用于文件下载、图片显示等） */
  baseUrl?: string;
  /** 是否显示编辑按钮 */
  haveEdit?: boolean;
  /** 是否显示查看按钮 */
  haveDetail?: boolean;
  /** 详情/编辑抽屉列配置 */
  drawerColumns: DrawerColumn[];
  /** 新增抽屉列配置 */
  drawerAddColumns?: DrawerColumn[];
  /** 是否显示自定义按钮 */
  havCustomBtn?: boolean;
  /** 自定义按钮文本 */
  havCustomStr?: string;
  /** 表格最大高度 */
  maxHeight?: string | number;
  /** 操作列宽度 */
  actionColumnWidth?: number;
  /** 抽屉宽度 */
  drawerWidth?: string | number;
  /** 分页每页条数选项 */
  pageSizes?: number[];
  /** 分页布局 */
  paginationLayout?: string;
  /** 新增权限编码（可选） */
  createPermission?: string;
  /** 编辑权限编码（可选） */
  updatePermission?: string;
  /** 删除权限编码（可选） */
  deletePermission?: string;
  /** 查看详情权限编码（可选） */
  detailPermission?: string;
  /** 自定义按钮权限编码（可选） */
  customPermission?: string;
}

const props = withDefaults(defineProps<Props>(), {
  pageSize: 10,
  currentPage: 1,
  baseUrl: '',
  haveEdit: true,
  haveDetail: true,
  havCustomBtn: false,
  havCustomStr: '自定义按钮',
  maxHeight: 600,
  actionColumnWidth: 220,
  drawerWidth: '50%',
  pageSizes: () => [10, 20, 50, 100],
  paginationLayout: 'total, sizes, prev, pager, next, jumper',
  createPermission: '',
  updatePermission: '',
  deletePermission: '',
  detailPermission: '',
  customPermission: ''
})

// Emits
const emit = defineEmits<{
  'size-change': [size: number]
  'current-change': [page: number]
  'onclick-delete': [index: number, row: Record<string, unknown>]
  'onclick-updateRow': [row: Record<string, unknown>]
  'onclick-addRow': [row: Record<string, unknown>]
  'onclick-btnCustom': [index: number, row: Record<string, unknown>]
}>()

// 内部状态
const drawerVisible = ref(false)
const drawerVisibleAdd = ref(false)
const drawerRow = ref<Record<string, unknown>>({})
const drawerAddRow = ref<Record<string, unknown>>({})
const isEdit = ref(false)

// 分页模型（支持 v-model）
const currentPageModel = computed({
  get: () => props.currentPage,
  set: (val) => {
    // 通过 emit 通知父组件更新，而不是直接修改内部状态
    emit('current-change', val)
  }
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (val) => {
    // 通过 emit 通知父组件更新，而不是直接修改内部状态
    emit('size-change', val)
  }
})

// 计算属性
const displayedData = computed(() => props.data)

// 权限相关计算属性（未传权限编码时默认允许）
const canCreate = computed(
  () => !props.createPermission || hasPermission(props.createPermission)
)
const canUpdate = computed(
  () => !props.updatePermission || hasPermission(props.updatePermission)
)
const canDelete = computed(
  () => !props.deletePermission || hasPermission(props.deletePermission)
)
const canDetail = computed(
  () => !props.detailPermission || hasPermission(props.detailPermission)
)
const canCustom = computed(
  () => !props.customPermission || hasPermission(props.customPermission)
)

// 方法
const formatBytes = (bytes: number): string => {
  if (!bytes) {
    return '0MB'
  }
  return `${(bytes / (1024 * 1024)).toFixed(2)}MB`
}

const handleSizeChange = (size: number) => {
  emit('size-change', size)
}

const handleCurrentChange = (page: number) => {
  emit('current-change', page)
}

const handleEdit = (index: number, row: Record<string, unknown>, edit: boolean) => {
  cancelAdd()
  isEdit.value = edit
  drawerRow.value = {...row}
  drawerVisible.value = true
}

const handleDelete = (index: number, row: Record<string, unknown>) => {
  emit('onclick-delete', index, row)
}

const updateItem = () => {
  emit('onclick-updateRow', drawerRow.value)
  cancelEdit()
}

const cancelEdit = () => {
  drawerRow.value = {}
  drawerVisible.value = false
  isEdit.value = false
}

const showAdd = () => {
  cancelEdit()
  drawerAddRow.value = {}
  drawerVisibleAdd.value = true
}

const handleAdd = () => {
  emit('onclick-addRow', drawerAddRow.value)
  cancelAdd()
}

const cancelAdd = () => {
  drawerAddRow.value = {}
  drawerVisibleAdd.value = false
}

const handleCopyUrl = async (url: string) => {
  const success = await copyToClipboard(url)
  if (success) {
    ElMessage.success(t('common.copySuccess') || '链接已复制到剪贴板')
  } else {
    ElMessage.error(t('common.copyFail') || '复制失败，请手动复制')
  }
}

const handleBtnCustom = (index: number, row: Record<string, unknown>) => {
  emit('onclick-btnCustom', index, row)
}
</script>

<style scoped lang="scss">
.d2-table {
  &__toolbar {
    margin-bottom: 16px;
  }

  &__pagination {
    display: flex;
    justify-content: center;
    margin-top: 16px;
  }

  &__image {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  &__upload {
    margin-top: 8px;
  }
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-placeholder);
  font-size: 20px;
}
</style>

