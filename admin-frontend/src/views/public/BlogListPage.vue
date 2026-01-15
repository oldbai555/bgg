<template>
  <div class="blog-list-page">
    <MetricReporter module="blog_article_list" :biz-id="0" />
    <BlogHeader />
    <div class="blog-page-container">
      <div class="blog-content-wrapper">
        <!-- 左侧分类导航 -->
        <aside class="blog-sidebar-left">
          <BlogCategoryNav :selected-tag-id="query.tagId" @select="handleTagSelect" />
        </aside>

        <!-- 中间文章列表 -->
        <main class="blog-main">
          <div v-loading="loading" class="article-list">
            <div
              v-for="item in list"
              :key="item.id"
              class="blog-article-card"
              @click="goToDetail(item.id)"
            >
              <div class="article-cover">
                <img
                  v-if="item.cover"
                  :src="item.cover"
                  :alt="item.title"
                  @error="handleImageError"
                />
                <div v-else class="cover-fallback">{{ firstChar(item.title) }}</div>
                <div v-if="item.isTop === 1" class="top-badge">置顶</div>
              </div>
              <div class="article-content">
                <div class="article-title">
                  <span v-if="item.isTop === 1" class="top-icon">📌</span>
                  {{ item.title }}
                </div>
                <div class="article-meta">
                  <span class="author">{{ item.authorName || '匿名' }}</span>
                  <span class="dot">·</span>
                  <span class="time">{{ formatTime(item.publishTime) }}</span>
                </div>
                <div class="article-summary">{{ item.summary || '暂无摘要' }}</div>
                <div v-if="item.tagNames?.length" class="article-tags">
                  <el-tag
                    v-for="tag in item.tagNames"
                    :key="tag"
                    size="small"
                    effect="plain"
                  >
                    {{ tag }}
                  </el-tag>
                </div>
              </div>
            </div>

            <div v-if="!loading && list.length === 0" class="empty-message">暂无文章</div>
          </div>

          <!-- 分页 -->
          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="query.page"
              v-model:page-size="query.size"
              :total="total"
              :page-sizes="[10, 20, 30, 50]"
              :layout="paginationLayout"
              :size="isMobile ? 'small' : 'default'"
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
            />
          </div>
        </main>

        <!-- 右侧侧边栏 -->
        <aside class="blog-sidebar">
          <BlogAuthorCard />
          <BlogSocialLinks />
        </aside>
      </div>
    </div>

    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, computed, onMounted, onUnmounted, nextTick} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {blogApi} from '@/api/blog'
import type {PublicBlogArticleListReq, PublicBlogArticleItem} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import BlogHeader from '@/components/blog/BlogHeader.vue'
import BlogCategoryNav from '@/components/blog/BlogCategoryNav.vue'
import BlogAuthorCard from '@/components/blog/BlogAuthorCard.vue'
import BlogSocialLinks from '@/components/blog/BlogSocialLinks.vue'

const router = useRouter()
const route = useRoute()

const SCROLL_STATE_KEY = 'public_blog_list_state'

const query = reactive<PublicBlogArticleListReq>({
  page: 1,
  size: 10,
  keyword: '',
  tagId: 0
})

const list = ref<PublicBlogArticleItem[]>([])
const total = ref(0)
const loading = ref(false)

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

const firstChar = (text: string): string => {
  if (!text) {
return '文'
}
  return Array.from(text)[0] ?? '文'
}

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}

const buildReq = (): PublicBlogArticleListReq => {
  const req: PublicBlogArticleListReq = {
    page: query.page,
    size: query.size
  }
  if (query.keyword) {
    req.keyword = query.keyword
  }
  if (query.tagId && query.tagId > 0) {
    req.tagId = query.tagId
  }
  return req
}

const pendingScrollTop = ref<number | null>(null)
const isMobile = ref(false)

const paginationLayout = computed(() =>
  isMobile.value ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'
)

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const handleResize = () => {
  checkMobile()
}

const restoreScrollPosition = async (scrollTop: number) => {
  await nextTick()
  await new Promise((resolve) => requestAnimationFrame(resolve))
  window.scrollTo({top: scrollTop, behavior: 'auto'})
}

const loadData = async () => {
  loading.value = true
  const shouldRestoreScroll = pendingScrollTop.value !== null
  const scrollTopToRestore = pendingScrollTop.value

  try {
    const resp = await blogApi.publicList(buildReq())
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    ElMessage.error(message)
  } finally {
    loading.value = false

    if (shouldRestoreScroll && scrollTopToRestore !== null) {
      pendingScrollTop.value = null
      await restoreScrollPosition(scrollTopToRestore)
    }
  }
}

const updateRouteQuery = () => {
  router.replace({
    path: route.path,
    query: {
      ...route.query,
      page: String(query.page),
      size: String(query.size),
      keyword: query.keyword || undefined,
      tagId: query.tagId && query.tagId > 0 ? String(query.tagId) : undefined
    }
  })
}

const handlePageChange = (page: number) => {
  query.page = page
  updateRouteQuery()
  pendingScrollTop.value = null
  loadData()
}

const handleSizeChange = (size: number) => {
  query.size = size
  query.page = 1
  updateRouteQuery()
  pendingScrollTop.value = null
  loadData()
}

const handleTagSelect = (tagId: number) => {
  query.tagId = tagId
  query.page = 1
  updateRouteQuery()
  pendingScrollTop.value = null
  loadData()
}

const goToDetail = (id: number) => {
  // 保存当前状态
  try {
    const state = {
      page: query.page,
      size: query.size,
      keyword: query.keyword,
      tagId: query.tagId,
      scrollTop: window.scrollY
    }
    sessionStorage.setItem(SCROLL_STATE_KEY, JSON.stringify(state))
  } catch {
    // 忽略存储错误
  }

  router.push(`/public/blog/${id}`)
}

// 从路由参数初始化查询条件
const initFromRoute = () => {
  const page = Number(route.query.page) || 1
  const size = Number(route.query.size) || 10
  const keyword = (route.query.keyword as string) || ''
  const tagId = Number(route.query.tagId) || 0

  query.page = page
  query.size = size
  query.keyword = keyword
  query.tagId = tagId
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', handleResize)
  initFromRoute()
  loadData()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped lang="scss">
@import '@/styles/blog.scss';

.blog-list-page {
  min-height: 100vh;
  background: #f5f5f5;
  // 页面整体可以滚动，但左侧和右侧使用 sticky 定位固定

  .article-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
    margin-bottom: 24px;
  }

  .cover-fallback {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 48px;
    font-weight: 600;
    color: #fff;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .empty-message {
    text-align: center;
    color: #999;
    padding: 40px 0;
    font-size: 14px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 20px 0;
  }

  .blog-sidebar-left {
    // 左侧导航栏样式已在组件内部定义
    // 使用 sticky 定位，在 blog.scss 中定义
  }
}

// 移动端适配
@media (max-width: 768px) {
  .blog-list-page {
    .blog-page-container {
      padding-top: 50px;
    }

    .article-list {
      gap: 16px;
    }
  }
}
</style>
