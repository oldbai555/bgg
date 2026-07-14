<template>
  <div class="blog-list-page public-list-page">
    <MetricReporter module="blog_article_list" :biz-id="0" />
    <PublicHeader />
    <div class="page-shell">
      <div class="page-intro">
        <h1 class="page-intro__title">技术博客</h1>
        <p class="page-intro__desc">记录开发过程中的思考与实践</p>
        <div class="page-intro__search">
          <el-input
            v-model="query.keyword"
            placeholder="搜索文章..."
            clearable
            @keydown.enter="handleSearch"
            @clear="handleSearch"
          >
            <template #append>
              <el-button type="primary" :loading="loading" @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
      </div>

      <div class="page-layout">
        <!-- 左侧分类导航 -->
        <BlogCategoryNav :selected-tag-id="query.tagId" @select="handleTagSelect" />

        <!-- 中间文章列表 -->
        <div>
          <div v-loading="loading" class="card-grid">
            <div
              v-for="item in list"
              :key="item.id"
              class="list-card"
              @click="goToDetail(item.id)"
            >
              <div class="cover">
                <img
                  v-if="item.cover"
                  :src="item.cover"
                  :alt="item.title"
                  @error="handleImageError"
                />
                <div v-else class="cover-fallback">{{ firstChar(item.title) }}</div>
                <div v-if="item.isTop === 1" class="top-badge">置顶</div>
              </div>
              <div class="card-content">
                <div class="card-title">
                  <span v-if="item.isTop === 1" class="top-icon">📌</span>
                  {{ item.title }}
                </div>
                <div class="card-summary">{{ item.summary || '暂无摘要' }}</div>
                <div class="card-meta">
                  <span>{{ item.authorName || '匿名' }}</span>
                  <span class="dot">·</span>
                  <span>{{ formatTime(item.publishTime) }}</span>
                </div>
                <div v-if="item.tagNames?.length" class="card-tags">
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

          <div class="pagination-bar">
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
        </div>

        <!-- 右侧侧边栏 -->
        <div class="blog-sidebar-stack">
          <BlogAuthorCard />
          <BlogSocialLinks />
        </div>
      </div>
    </div>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, computed, onMounted, onUnmounted, nextTick, watch} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {contentApi} from '@/api/content'
import type {PublicBlogArticleListReq, PublicBlogArticleItem} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import PublicHeader from '@/components/common/PublicHeader.vue'
import BlogCategoryNav from '@/components/blog/BlogCategoryNav.vue'
import BlogAuthorCard from '@/components/blog/BlogAuthorCard.vue'
import BlogSocialLinks from '@/components/blog/BlogSocialLinks.vue'
import {MOBILE_BREAKPOINT} from '@/constants/breakpoints'

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
  if (typeof window !== 'undefined') {
    isMobile.value = window.innerWidth <= MOBILE_BREAKPOINT
  }
}

const handleResize = () => {
  checkMobile()
}

const restoreScrollPosition = async (scrollTop: number) => {
  if (typeof window === 'undefined') {
    return
  }
  await nextTick()
  await new Promise((resolve) => requestAnimationFrame(resolve))
  window.scrollTo({top: scrollTop, behavior: 'auto'})
}

const loadData = async () => {
  loading.value = true
  const shouldRestoreScroll = pendingScrollTop.value !== null
  const scrollTopToRestore = pendingScrollTop.value

  try {
    const resp = await contentApi.publicArticleList(buildReq())
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    ElMessage.error(message)
  } finally {
    loading.value = false

    if (shouldRestoreScroll && scrollTopToRestore !== null && typeof window !== 'undefined') {
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

const handleSearch = () => {
  query.page = 1
  updateRouteQuery()
  pendingScrollTop.value = null
  loadData()
}

const goToDetail = (id: number) => {
  // 保存当前状态
  if (typeof window !== 'undefined') {
    try {
      const state = {
        page: query.page,
        size: query.size,
        keyword: query.keyword,
        tagId: query.tagId,
        scrollTop: window.scrollY,
        ts: Date.now()
      }
      sessionStorage.setItem(SCROLL_STATE_KEY, JSON.stringify(state))
    } catch {
      // 忽略存储错误
    }
  }

  router.push(`/front/blog/${id}`)
}

// 从详情页返回时恢复滚动位置：page/size/keyword/tagId 已经通过 route.query 恢复
// （updateRouteQuery 在导航前已写入 URL，router.back() 会带回同一个 URL），这里只需要
// 把 sessionStorage 里记录的 scrollTop 接回 pendingScrollTop，交给 loadData() 完成后恢复
const restorePendingScroll = () => {
  if (typeof window === 'undefined') {
    return
  }
  try {
    const raw = sessionStorage.getItem(SCROLL_STATE_KEY)
    if (!raw) {
      return
    }
    const parsed = JSON.parse(raw) as {scrollTop?: number; ts?: number}
    const now = Date.now()
    if (parsed.ts && now - parsed.ts >= 60 * 60 * 1000) {
      sessionStorage.removeItem(SCROLL_STATE_KEY)
      return
    }
    if (typeof parsed.scrollTop === 'number' && parsed.scrollTop > 0) {
      pendingScrollTop.value = parsed.scrollTop
    }
  } catch {
    // 忽略解析错误
  }
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

// 监听路由参数变化（特别是 keyword 和 tagId），当变化时重新加载数据
watch(
  () => route.query,
  () => {
    initFromRoute()
    pendingScrollTop.value = null
    loadData()
  },
  {deep: true}
)

onMounted(() => {
  checkMobile()
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleResize)
  }
  initFromRoute()
  restorePendingScroll()
  loadData()
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleResize)
  }
})
</script>

<style scoped lang="scss">
@import '@/styles/public-list.scss';

.blog-list-page .page-layout {
  grid-template-columns: 200px 1fr 240px;
}

.blog-sidebar-stack {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

@include mobile {
  .blog-list-page .page-layout {
    grid-template-columns: 1fr;
  }

  .blog-sidebar-stack {
    order: 2;
  }
}
</style>
