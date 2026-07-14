<template>
  <div :class="['blog-detail-page', 'public-detail-page', {'is-loading': loading}]">
    <MetricReporter module="blog_article_detail" :biz-id="Number(route.params.id) || 0" />
    <PublicHeader />
    <div class="page-shell">
      <div class="page-layout">
        <!-- 左侧分类导航 -->
        <BlogCategoryNav :selected-tag-id="currentTagId" @select="handleTagSelect" />

        <!-- 中间文章内容 -->
        <article v-if="detail" class="detail-card">
          <div class="back-link" @click="goBack">← 返回列表</div>
          <h1 class="title">{{ detail.title }}</h1>

          <div class="meta">
            <span>{{ detail.authorName || '匿名' }}</span>
            <span class="dot">·</span>
            <span>{{ formatTime(detail.publishTime || 0) }}</span>
            <template v-if="wordCount > 0">
              <span class="dot">·</span>
              <span>约{{ wordCount }}字</span>
            </template>
            <template v-if="readingTime > 0">
              <span class="dot">·</span>
              <span>大约{{ readingTime }}分钟</span>
            </template>
            <span v-if="detail.tags?.length" class="tags">
              <el-tag
                v-for="tag in detail.tags"
                :key="tag.id"
                size="small"
                effect="plain"
              >
                {{ tag.name }}
              </el-tag>
            </span>
          </div>

          <div v-if="detail.cover" class="cover">
            <img :src="detail.cover" :alt="detail.title" @error="handleImageError" />
          </div>

          <div class="content">
            <component
              v-if="mdPreviewLoaded && MdPreview && detail?.content"
              :is="MdPreview"
              :editor-id="'public-blog-detail'"
              :model-value="detail.content"
              :preview-theme="'github'"
              @html-changed="handleContentRendered"
            />
            <div v-else-if="!mdPreviewLoaded || !MdPreview" class="loading-placeholder">
              加载中...
            </div>
            <div v-else-if="!detail?.content" class="empty-content">
              暂无内容
            </div>
          </div>

          <!-- 相邻文章导航 -->
          <div v-if="prevArticle || nextArticle" class="detail-navigation">
            <router-link
              v-if="prevArticle"
              :to="`/front/blog/${prevArticle.id}`"
              class="nav-item nav-prev"
            >
              <div class="nav-label">← 上一页</div>
              <div class="nav-title">{{ prevArticle.title }}</div>
            </router-link>
            <div v-else class="nav-item nav-prev disabled">
              <div class="nav-label">← 上一页</div>
              <div class="nav-title">没有更多了</div>
            </div>

            <router-link
              v-if="nextArticle"
              :to="`/front/blog/${nextArticle.id}`"
              class="nav-item nav-next"
            >
              <div class="nav-label">下一页 →</div>
              <div class="nav-title">{{ nextArticle.title }}</div>
            </router-link>
            <div v-else class="nav-item nav-next disabled">
              <div class="nav-label">下一页 →</div>
              <div class="nav-title">没有更多了</div>
            </div>
          </div>
        </article>

        <div v-else-if="!loading" class="detail-card">
          <div class="empty">文章不存在或已下架</div>
        </div>

        <!-- 右侧目录 -->
        <BlogTOC v-if="detail?.content" :content="detail.content" />
      </div>
    </div>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {ref, computed, onMounted, watch, shallowRef, nextTick} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {contentApi} from '@/api/content'
import type {
  PublicBlogArticleDetailReq,
  PublicBlogArticleDetailResp,
  PublicBlogArticlePrevReq,
  PublicBlogArticleNextReq
} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import PublicHeader from '@/components/common/PublicHeader.vue'
import BlogCategoryNav from '@/components/blog/BlogCategoryNav.vue'
import BlogTOC from '@/components/blog/BlogTOC.vue'

const route = useRoute()
const router = useRouter()

// md-editor-v3 动态导入
const MdPreview = shallowRef<any>(null)
const mdPreviewLoaded = ref(false)

// 在客户端加载 md-editor-v3
if (typeof window !== 'undefined') {
  import('md-editor-v3').then((module) => {
    // md-editor-v3 的类型声明只有具名导出，但不同打包场景下运行时可能落在 default 上，这里做防御性兜底
    const mod = module as unknown as {MdPreview?: unknown; default?: {MdPreview?: unknown} | unknown}
    MdPreview.value = mod.MdPreview || (mod.default as {MdPreview?: unknown} | undefined)?.MdPreview || mod.default
    mdPreviewLoaded.value = true
  }).catch((err) => {
    console.error('加载 md-editor-v3 失败:', err)
  })
  import('md-editor-v3/lib/style.css').catch((err) => {
    console.error('加载 md-editor-v3 样式失败:', err)
  })
}

const detail = ref<PublicBlogArticleDetailResp | null>(null)
const loading = ref(false)
const prevArticle = ref<{id: number; title: string; publishTime: number} | null>(null)
const nextArticle = ref<{id: number; title: string; publishTime: number} | null>(null)

const currentTagId = computed(() => {
  return detail.value?.tags?.[0]?.id || 0
})

// 计算字数和阅读时间
const wordCount = computed(() => {
  if (!detail.value?.content) {
    return 0
  }
  const text = detail.value.content
    .replace(/#{1,6}\s+/g, '')
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/`([^`]+)`/g, '')
    .replace(/```[\s\S]*?```/g, '')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, '')
    .replace(/\n/g, '')
    .trim()
  return text.length
})

const readingTime = computed(() => {
  const count = wordCount.value
  if (count === 0) {
    return 0
  }
  return Math.max(1, Math.ceil(count / 300))
})

const formatTime = (ts: number): string => {
  if (!ts) {
    return '-'
  }
  const d = new Date(ts * 1000)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}

const handleContentRendered = (html: string) => {
  // Markdown渲染完成后，为标题添加ID
  if (typeof window !== 'undefined') {
    nextTick(() => {
      const toolbarWrapper = document.querySelector('.md-editor-toolbar-wrapper')
      const toolbar = document.querySelector('.md-editor-toolbar')
      if (toolbarWrapper) {
        (toolbarWrapper as HTMLElement).style.display = 'none'
      }
      if (toolbar) {
        (toolbar as HTMLElement).style.display = 'none'
      }

      const inputWrapper = document.querySelector('.md-editor-input-wrapper')
      if (inputWrapper) {
        (inputWrapper as HTMLElement).style.display = 'none'
      }

      const footer = document.querySelector('.md-editor-footer')
      if (footer) {
        (footer as HTMLElement).style.display = 'none'
      }

      const catalog = document.querySelector('.md-editor-catalog')
      if (catalog) {
        (catalog as HTMLElement).style.display = 'none'
      }
    })
  }
}

const handleTagSelect = (tagId: number) => {
  router.push({
    path: '/front/blog',
    query: {tagId: tagId > 0 ? tagId : undefined}
  })
}

// 返回列表
const goBack = () => {
  if (typeof window !== 'undefined' && window.history.length > 1) {
    router.back()
  } else {
    if (typeof window !== 'undefined') {
      try {
        const raw = sessionStorage.getItem('public_blog_list_state')
        if (raw) {
          const parsed = JSON.parse(raw) as {
            page?: number
            size?: number
            keyword?: string
            tagId?: number
            scrollTop?: number
          }
          router.push({
            path: '/front/blog',
            query: {
              ...(parsed.page && {page: String(parsed.page)}),
              ...(parsed.size && {size: String(parsed.size)}),
              ...(parsed.keyword && {keyword: parsed.keyword}),
              ...(parsed.tagId && parsed.tagId > 0 && {tagId: String(parsed.tagId)})
            }
          })
          return
        }
      } catch {
        // 忽略解析错误
      }
    }
    router.push('/front/blog')
  }
}

// 加载文章详情
const loadDetail = async () => {
  const id = Number(route.params.id)
  if (!id) {
    ElMessage.error('文章ID不正确')
    return
  }

  loading.value = true
  try {
    const req: PublicBlogArticleDetailReq = {id}
    const resp = await contentApi.publicArticleDetail(req)
    detail.value = resp
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载文章失败'
    ElMessage.error(message)
    detail.value = null
  } finally {
    loading.value = false
  }
}

// 加载相邻文章
const loadPrevArticle = async (currentId: number) => {
  try {
    const req: PublicBlogArticlePrevReq = {id: currentId}
    const resp = await contentApi.publicArticlePrev(req)
    if (resp.id) {
      prevArticle.value = {
        id: resp.id,
        title: resp.title,
        publishTime: resp.publishTime
      }
    } else {
      prevArticle.value = null
    }
  } catch (err) {
    console.error('加载上一篇文章失败:', err)
    prevArticle.value = null
  }
}

const loadNextArticle = async (currentId: number) => {
  try {
    const req: PublicBlogArticleNextReq = {id: currentId}
    const resp = await contentApi.publicArticleNext(req)
    if (resp.id) {
      nextArticle.value = {
        id: resp.id,
        title: resp.title,
        publishTime: resp.publishTime
      }
    } else {
      nextArticle.value = null
    }
  } catch (err) {
    console.error('加载下一篇文章失败:', err)
    nextArticle.value = null
  }
}

// 监听详情数据变化，加载相邻文章
watch(() => detail.value?.id, async (newId) => {
  if (newId) {
    await Promise.all([loadPrevArticle(newId), loadNextArticle(newId)])
  }
}, {immediate: true})

// 监听路由参数变化
watch(
  () => route.params.id,
  (newId, oldId) => {
    if (newId && newId !== oldId) {
      prevArticle.value = null
      nextArticle.value = null
      loadDetail()
    }
  },
  {immediate: false}
)

onMounted(() => {
  loadDetail()
})
</script>

<style scoped lang="scss">
@import '@/styles/public-detail.scss';

.blog-detail-page .page-layout {
  grid-template-columns: 200px 1fr 240px;
}

.blog-detail-page .detail-card {
  .detail-navigation {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid var(--color-border-light);
    display: flex;
    justify-content: space-between;
    gap: 16px;

    .nav-item {
      flex: 1;
      padding: 12px 16px;
      border-radius: 6px;
      background: var(--color-bg-secondary);
      text-decoration: none;
      color: var(--color-text-regular);
      transition: all 0.2s;

      &:hover {
        background: color-mix(in srgb, var(--color-primary) 12%, transparent);
        color: var(--color-primary);
      }

      &.disabled {
        cursor: not-allowed;
        opacity: 0.5;
      }

      .nav-label {
        font-size: 12px;
        color: var(--color-text-secondary);
        margin-bottom: 4px;
      }

      .nav-title {
        font-size: 14px;
        font-weight: 500;
      }
    }
  }

  // 全局隐藏 rn-wrapper 行号（md-editor-v3 使用的行号包装器）
  // 在 .content 作用域内的所有层级都要隐藏
  .content {
    :deep([rn-wrapper]),
    :deep(span[rn-wrapper]),
    :deep(*[rn-wrapper]) {
      display: none !important;
      visibility: hidden !important;
      width: 0 !important;
      height: 0 !important;
      padding: 0 !important;
      margin: 0 !important;
      opacity: 0 !important;
      position: absolute !important;
      left: -9999px !important;
    }
    min-height: 200px;

    :deep(.md-editor) {
      border: none !important;
      box-shadow: none !important;
      background: transparent !important;
      display: block !important;
      width: 100% !important;
      min-width: 0 !important;
      min-height: 200px !important;
      height: auto !important;
      box-sizing: border-box !important;
      overflow: visible !important;
    }

    :deep(.md-editor-toolbar-wrapper),
    :deep(.md-editor-toolbar) {
      display: none !important;
      visibility: hidden !important;
    }

    :deep(.md-editor-input-wrapper) {
      display: none !important;
      visibility: hidden !important;
    }

    :deep(.md-editor-content) {
      display: block !important;
      width: 100% !important;
      min-width: 0 !important;
      height: auto !important;
    }

    :deep(.md-editor-footer) {
      display: none !important;
      visibility: hidden !important;
    }

    :deep(.md-editor-catalog) {
      display: none !important;
      visibility: hidden !important;
    }

    :deep(.md-editor-preview-wrapper) {
      display: block !important;
      width: 100% !important;
      min-width: 0 !important;
      max-width: 100% !important;
      padding: 0 !important;
      background: transparent !important;
      visibility: visible !important;
      opacity: 1 !important;
      min-height: 200px !important;
      height: auto !important;
      overflow: visible !important;
    }

    :deep(.md-editor-preview) {
      display: block !important;
      width: 100% !important;
      min-width: 0 !important;
      max-width: 100% !important;
      padding: 0 !important;
      background: transparent !important;
      visibility: visible !important;
      opacity: 1 !important;
      min-height: 200px !important;
      height: auto !important;
      overflow: visible !important;

      // 隐藏 rn-wrapper 行号（md-editor-v3 使用的行号包装器）
      // 在预览区域内的所有层级都要隐藏
      [rn-wrapper],
      span[rn-wrapper],
      *[rn-wrapper] {
        display: none !important;
        visibility: hidden !important;
        width: 0 !important;
        height: 0 !important;
        padding: 0 !important;
        margin: 0 !important;
        opacity: 0 !important;
      }

      // 隐藏代码块行号（针对 highlight.js 生成的行号）
      pre {
        .hljs-ln-numbers,
        .hljs-ln-n {
          display: none !important;
        }

        .hljs-ln-code,
        .hljs-ln-line {
          padding-left: 0 !important;
          margin-left: 0 !important;
        }

        .hljs-ln {
          .hljs-ln-numbers {
            display: none !important;
          }
        }

        code {
          padding-left: 0 !important;
          margin-left: 0 !important;
        }

        [class*='line-number'],
        [class*='ln-numbers'],
        [class*='ln-n'] {
          display: none !important;
        }

        [rn-wrapper],
        span[rn-wrapper] {
          display: none !important;
        }
      }

      // 隐藏 CodeMirror 编辑器相关的行号（如果存在）
      .cm-gutters,
      .cm-gutter,
      .cm-lineNumbers,
      .cm-gutterElement,
      .cm-gutter.cm-lineNumbers {
        display: none !important;
      }

      .cm-content {
        padding-left: 0 !important;
      }

      .line-number,
      .code-line-number,
      .line-numbers {
        display: none !important;
      }
    }

    .loading-placeholder,
    .empty-content {
      text-align: center;
      padding: 40px 0;
      color: var(--color-text-secondary);
    }
  }
}

@include mobile {
  .blog-detail-page .page-layout {
    grid-template-columns: 1fr;
  }
}
</style>
