<template>
  <div v-loading="loading" class="blog-detail-page">
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
              <MdPreview
                :editor-id="'public-blog-detail'"
                :model-value="detail.content || ''"
                :preview-theme="'github'"
                @html-changed="handleContentRendered"
              />
            </div>

            <!-- 相邻文章导航 -->
            <div v-if="prevArticle || nextArticle" class="detail-navigation">
              <router-link
                v-if="prevArticle"
                :to="`/public/blog/${prevArticle.id}`"
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
                :to="`/public/blog/${nextArticle.id}`"
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
import {ref, computed, onMounted, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {MdPreview} from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
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

const loading = ref(false)
const detail = ref<PublicBlogArticleDetailResp | null>(null)
const prevArticle = ref<{id: number; title: string; publishTime: number} | null>(null)
const nextArticle = ref<{id: number; title: string; publishTime: number} | null>(null)

const currentTagId = computed(() => {
  // 从文章标签中获取第一个标签ID作为当前分类
  return detail.value?.tags?.[0]?.id || 0
})

// 计算字数和阅读时间
const wordCount = computed(() => {
  if (!detail.value?.content) {
return 0
}
  // 去除Markdown语法，统计纯文本字数
  const text = detail.value.content
    .replace(/#{1,6}\s+/g, '') // 标题
    .replace(/\*\*([^*]+)\*\*/g, '$1') // 加粗
    .replace(/\*([^*]+)\*/g, '$1') // 斜体
    .replace(/`([^`]+)`/g, '') // 行内代码
    .replace(/```[\s\S]*?```/g, '') // 代码块
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1') // 链接
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, '') // 图片
    .replace(/\n/g, '') // 换行
    .trim()
  return text.length
})

const readingTime = computed(() => {
  const count = wordCount.value
  if (count === 0) {
return 0
}
  // 按300字/分钟计算
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

const handleContentRendered = () => {
  // Markdown渲染完成后，为标题添加ID（BlogTOC组件会处理）
  // 这里可以触发TOC更新
}

const handleTagSelect = (tagId: number) => {
  router.push({
    path: '/public/blog',
    query: {tagId: tagId > 0 ? tagId : undefined}
  })
}

const loadDetail = async () => {
  const id = Number(route.params.id)
  if (!id) {
    ElMessage.error('文章ID不正确')
    return
  }

  loading.value = true
  try {
    const req: PublicBlogArticleDetailReq = {id}
    const resp = await blogApi.publicDetail(req)
    detail.value = resp

    // 加载相邻文章
    await Promise.all([loadPrevArticle(id), loadNextArticle(id)])
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载文章失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

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

// 监听路由参数变化，重新加载文章详情
watch(
  () => route.params.id,
  (newId, oldId) => {
    // 只有当 ID 真正变化时才重新加载
    if (newId && newId !== oldId) {
      // 重置状态
      detail.value = null
      prevArticle.value = null
      nextArticle.value = null
      // 重新加载
      loadDetail()
    }
  },
  {immediate: false} // 不在初始化时执行，因为 onMounted 会处理
)

onMounted(() => {
  loadDetail()
})
</script>

<style scoped lang="scss">
@import '@/styles/blog.scss';

.blog-detail-page {
  min-height: 100vh;
  background: #f5f5f5;
  // 页面整体可以滚动，但左侧和右侧使用 sticky 定位固定

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
}

// 移动端适配
@media (max-width: 768px) {
  .blog-detail-page {
    .blog-page-container {
      padding-top: 50px;
    }
  }
}
</style>
