<template>
  <div :class="['blog-detail-page', { 'is-loading': loading }]">
    <MetricReporter module="blog_article_detail" :biz-id="Number(route.params.id) || 0" />
    <BlogHeader />
    <div class="blog-page-container">
      <div class="blog-content-wrapper">
        <!-- 左侧分类导航 -->
        <aside class="blog-sidebar-left">
          <BlogCategoryNav :selected-tag-id="currentTagId" @select="handleTagSelect" />
        </aside>

        <!-- 中间文章内容 -->
        <main class="blog-main">
          <div v-if="detail" class="blog-detail-container">
            <!-- 返回按钮 -->
            <div class="back-link" @click="goBack">← 返回列表</div>
            <h1 class="detail-title">{{ detail.title }}</h1>

            <div class="detail-meta">
              <span class="meta-item">
                <span class="author">{{ detail.authorName || '匿名' }}</span>
              </span>
              <span class="meta-item">
                <span class="time">{{ formatTime(detail.publishTime || 0) }}</span>
              </span>
              <span v-if="wordCount > 0" class="meta-item">
                <span class="word-count">约{{ wordCount }}字</span>
              </span>
              <span v-if="readingTime > 0" class="meta-item">
                <span class="reading-time">大约{{ readingTime }}分钟</span>
              </span>
              <span v-if="detail.tags?.length" class="meta-item">
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

            <div v-if="detail.cover" class="detail-cover">
              <img :src="detail.cover" :alt="detail.title" @error="handleImageError" />
            </div>

            <div class="detail-content">
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
                :to="`/blog/${prevArticle.id}`"
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
                :to="`/blog/${nextArticle.id}`"
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
          </div>

          <div v-else-if="!loading" class="empty">文章不存在或已下架</div>
        </main>

        <!-- 右侧目录 -->
        <aside class="blog-sidebar">
          <BlogTOC v-if="detail?.content" :content="detail.content" />
        </aside>
      </div>
    </div>

    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {ref, computed, onMounted, watch, shallowRef, nextTick} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {blogApi} from '@/api/blog'
import type {
  PublicBlogArticleDetailReq,
  PublicBlogArticleDetailResp,
  PublicBlogArticlePrevReq,
  PublicBlogArticleNextReq
} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import BlogHeader from '@/components/blog/BlogHeader.vue'
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
    MdPreview.value = module.MdPreview || module.default?.MdPreview || module.default
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
        ;(toolbarWrapper as HTMLElement).style.display = 'none'
      }
      if (toolbar) {
        ;(toolbar as HTMLElement).style.display = 'none'
      }
      
      const inputWrapper = document.querySelector('.md-editor-input-wrapper')
      if (inputWrapper) {
        ;(inputWrapper as HTMLElement).style.display = 'none'
      }
      
      const footer = document.querySelector('.md-editor-footer')
      if (footer) {
        ;(footer as HTMLElement).style.display = 'none'
      }
      
      const catalog = document.querySelector('.md-editor-catalog')
      if (catalog) {
        ;(catalog as HTMLElement).style.display = 'none'
      }
    })
  }
}

const handleTagSelect = (tagId: number) => {
  router.push({
    path: '/blog',
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
            path: '/blog',
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
    router.push('/blog')
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
    const resp = await blogApi.publicArticleDetail(req)
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
    const resp = await blogApi.publicArticlePrev(req)
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
    const resp = await blogApi.publicArticleNext(req)
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
@import '@/styles/blog.scss';

.blog-detail-page {
  height: 100vh;
  overflow: hidden;
  background: #f5f5f5;
  position: relative;

  .back-link {
    cursor: pointer;
    color: #666;
    margin-bottom: 12px;
    font-size: 14px;
    display: inline-block;
    transition: color 0.3s;
    flex-shrink: 0;

    &:hover {
      color: #409eff;
    }
  }

  .detail-navigation {
    .nav-item.disabled {
      cursor: not-allowed;
      opacity: 0.5;

      &:hover {
        background: #f5f5f5;
        color: #666;
      }
    }
  }

  .detail-content {
    min-height: 200px;
    width: 100% !important;
    height: auto !important;
    font-size: 16px;
    line-height: 1.8;
    color: #333;
    box-sizing: border-box;
    overflow: visible !important;

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
      color: #333 !important;
      font-size: 16px !important;
      line-height: 1.8 !important;
      overflow: visible !important;
    }

    .loading-placeholder,
    .empty-content {
      text-align: center;
      padding: 40px 0;
      color: #999;
    }
  }
}

@media (max-width: 768px) {
  .blog-detail-page {
    .blog-page-container {
      padding-top: 50px;
    }
  }
}
</style>
