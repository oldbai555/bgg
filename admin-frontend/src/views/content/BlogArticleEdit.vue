<template>
  <div v-loading="loading" class="page">
    <el-card>
      <div class="header">
        <div class="title">{{ isEdit ? '编辑文章' : '新增文章' }}</div>
        <div class="actions">
          <el-button native-type="button" @click.prevent="goBack">返回</el-button>
          <el-button
            v-if="canEdit"
            v-permission="isEdit ? 'blog_article:update' : 'blog_article:create'"
            type="primary"
            native-type="button"
            @click.prevent="handleSave"
          >
            {{ saveButtonText }}
          </el-button>
          <el-button
            v-if="canEdit && canSubmitAudit"
            v-permission="'blog_article:submit'"
            type="warning"
            native-type="button"
            @click.prevent="handleSubmit"
          >
            提交审核
          </el-button>
        </div>
      </div>

      <!-- 上架状态提示 -->
      <el-alert
        v-if="isEdit && articleStatus === 4"
        type="warning"
        :closable="false"
        show-icon
        class="mb-12"
      >
        <template #title>
          <span>当前文章已上架，不可编辑。如需编辑，请先下架文章。</span>
        </template>
      </el-alert>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        :disabled="!canEdit"
        label-width="90px"
        class="form"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入文章标题" clearable />
        </el-form-item>

        <el-form-item label="标签" prop="tagIds">
          <el-select
            v-model="form.tagIds"
            multiple
            filterable
            placeholder="请选择标签（至少 1 个）"
            style="width: 100%"
          >
            <el-option
              v-for="t in tagOptions"
              :key="t.value"
              :label="t.label"
              :value="t.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="封面">
          <ImageUpload v-model="form.cover" tip="上传封面图（可选）" />
        </el-form-item>

        <el-form-item label="摘要">
          <el-input
            v-model="form.summary"
            type="textarea"
            :rows="3"
            placeholder="可选：不填则公共列表用正文截断"
          />
        </el-form-item>

        <el-form-item label="正文" prop="content">
          <div class="editor">
            <div class="editor-toolbar">
              <span class="hint">Markdown 编辑区（支持预览与快捷键）</span>
              <div class="toolbar-right">
                <el-popover placement="bottom" width="260" trigger="click">
                  <template #reference>
                    <el-button size="small">上传正文图片</el-button>
                  </template>
                  <div class="inline-upload">
                    <ImageUpload v-model="contentImageUrl" tip="上传后点击“插入到正文”" />
                    <el-button
                      size="small"
                      type="primary"
                      :disabled="!contentImageUrl"
                      @click="insertContentImage"
                    >
                      插入到正文
                    </el-button>
                  </div>
                </el-popover>
              </div>
            </div>

            <div class="editor-body">
              <MdEditor
                v-model="form.content"
                :preview-theme="'github'"
                :language="'zh-CN'"
                style="min-height: 420px"
              />
            </div>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import type {FormInstance, FormRules} from 'element-plus'
import {MdEditor} from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import {contentApi} from '@/api/content'
import ImageUpload from '@/components/common/ImageUpload.vue'
import type {BlogArticleCreateReq, BlogArticleUpdateReq} from '@/api/generated/admin'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const formRef = ref<FormInstance>()

const id = computed(() => Number(route.params.id || 0))
const isEdit = computed(() => id.value > 0)

const form = reactive<{
  title: string;
  tagIds: number[];
  cover: string;
  summary: string;
  content: string;
}>({
  title: '',
  tagIds: [],
  cover: '',
  summary: '',
  content: ''
})

const rules: FormRules = {
  title: [{required: true, message: '请输入标题', trigger: 'blur'}],
  tagIds: [{type: 'array', required: true, message: '请选择至少一个标签', trigger: 'change'}],
  content: [{required: true, message: '请输入正文内容', trigger: 'blur'}]
}

const tagOptions = ref<Array<{ label: string; value: number }>>([])

// 正文图片临时 URL（用于插入 markdown）
const contentImageUrl = ref('')

// 文章状态（用于控制编辑权限）
const articleStatus = ref<number>(0) // 0=新增，1=草稿，2=待审核，3=审核通过-未上架，4=上架，5=下架
const auditStatus = ref<number>(0) // 审核状态

// 编辑权限：上架状态（4）不可编辑，其他状态可编辑
const canEdit = computed(() => {
  if (!isEdit.value) {
    return true
  } // 新增模式始终可编辑
  return articleStatus.value !== 4 // 上架状态不可编辑
})

// 保存按钮文案：根据原状态显示不同文案
const saveButtonText = computed(() => {
  if (!isEdit.value) {
    return '保存草稿'
  }
  if (articleStatus.value === 1) {
    return '保存草稿'
  } // 草稿状态
  return '保存并重新提交审核' // 其他状态（待审核、审核通过-未上架、下架）
})

// 是否可以提交审核：仅草稿状态可以提交审核
const canSubmitAudit = computed(() => {
  if (!isEdit.value) {
    return false
  } // 新增模式需要先保存
  return articleStatus.value === 1 // 仅草稿状态可以提交审核
})

const loadTagOptions = async () => {
  // 使用专用 tags/options 接口（仅返回启用标签，更轻量）
  const resp = await contentApi.tagOptions({limit: 1000})
  tagOptions.value = (resp.list || []).map((t) => ({
    label: t.name,
    value: t.id
  }))
}

const loadDetail = async () => {
  if (!isEdit.value || id.value <= 0) {
    return
  }
  try {
    const resp = await contentApi.articleDetail({id: id.value})
    form.title = resp.title
    form.cover = resp.cover || ''
    form.summary = resp.summary || ''
    form.content = resp.content || ''
    form.tagIds = (resp.tags || []).map((t) => Number(t.id))
    // 记录文章状态和审核状态
    articleStatus.value = resp.status
    auditStatus.value = resp.auditStatus
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载文章详情失败'
    ElMessage.error(message)
  }
}

const goBack = () => {
  router.push('/admin/blog/article')
}

const handleSave = async () => {
  if (!canEdit.value) {
    ElMessage.warning('当前状态不允许编辑')
    return
  }

  const ok = await formRef.value?.validate().catch(() => false)
  if (!ok) {
    return
  }

  loading.value = true
  try {
    if (isEdit.value) {
      const req: BlogArticleUpdateReq = {
        id: id.value,
        title: form.title,
        content: form.content,
        tagIds: form.tagIds,
        cover: form.cover,
        summary: form.summary
      }
      await contentApi.articleUpdate(req)
      // 根据原状态显示不同的成功提示
      if (articleStatus.value === 1) {
        ElMessage.success('保存成功')
      } else {
        ElMessage.success('保存成功，文章已重新提交审核')
        // 更新状态为待审核
        articleStatus.value = 2
        auditStatus.value = 2
        // 重新提交审核后回到文章列表
        goBack()
      }
    } else {
      const req: BlogArticleCreateReq = {
        title: form.title,
        content: form.content,
        tagIds: form.tagIds,
        cover: form.cover,
        summary: form.summary
      }
      await contentApi.articleCreate(req)
      ElMessage.success('创建成功')
      goBack()
    }
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '保存失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!isEdit.value) {
    ElMessage.warning('请先保存文章，再提交审核')
    return
  }
  if (!canSubmitAudit.value) {
    ElMessage.warning('当前状态不允许提交审核')
    return
  }
  loading.value = true
  try {
    await contentApi.articleSubmit({id: id.value})
    ElMessage.success('已提交审核')
    // 更新状态为待审核
    articleStatus.value = 2
    auditStatus.value = 2
    goBack()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '提交失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const insertContentImage = () => {
  if (!contentImageUrl.value) {
    return
  }
  const md = `\n![](${contentImageUrl.value})\n`
  form.content = (form.content || '') + md
  ElMessage.success('已插入到正文末尾')
  contentImageUrl.value = ''
}

onMounted(async () => {
  loading.value = true
  try {
    await loadTagOptions()
    await loadDetail()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped lang="scss">
.page {
  padding: 16px 24px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;

  .title {
    font-size: 18px;
    font-weight: 600;
  }

  .actions {
    display: flex;
    gap: 8px;
  }
}

.editor {
  width: 100%;

  .editor-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;

    .hint {
      color: var(--color-text-regular);
      font-size: 12px;
    }
  }

  .editor-body {
    // 单列布局，编辑器占满整行
    width: 100%;
  }
}

.inline-upload {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>

