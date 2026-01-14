<template>
  <div class="blog-list-page public-list-page">
    <MetricReporter module="blog_article_list" :biz-id="0" />
    <div class="container">
      <div class="hero">
        <div class="hero-title">博客文章</div>
        <div class="hero-subtitle">发现有趣内容 · 持续更新</div>

        <!-- 搜索栏 -->
        <div class="search-bar">
          <el-input
            v-model="query.keyword"
            placeholder="搜索文章标题或内容..."
            clearable
            @keydown.enter="handleSearch"
          >
            <template #append>
              <el-button type="primary" :loading="loading" @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
      </div>

      <!-- 文章列表 -->
      <div v-loading="loading" class="list-grid">
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
          </div>
          <div class="card-content">
            <div class="card-title">{{ item.title }}</div>
            <div class="card-meta">
              <span class="author">{{ item.authorName || '匿名' }}</span>
              <span class="dot">·</span>
              <span class="time">{{ formatTime(item.publishTime || item.createdAt) }}</span>
            </div>
            <div class="card-summary">{{ item.summary || '暂无摘要' }}</div>
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

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          :page-sizes="[10, 20, 30, 50]"
          :layout="paginationLayout"
          :small="isMobile"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>

      <IcpFooter />
    </div>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, computed, onMounted, onUnmounted, nextTick} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {blogApi} from '@/api/blog'
import type {
  PublicBlogArticleListReq,
  PublicBlogArticleItem
} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'

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
      keyword: query.keyword || undefined
    }
  })
}

const handleSearch = () => {
  query.page = 1
  updateRouteQuery()
  pendingScrollTop.value = null
  loadData()
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

const goToDetail = (id: number) => {
  // 进入详情前记录当前查询条件与滚动位置，方便返回时恢复
  try {
    const state = {
      page: query.page,
      size: query.size,
      keyword: query.keyword,
      scrollTop: window.scrollY,
      ts: Date.now()
    }
    sessionStorage.setItem(SCROLL_STATE_KEY, JSON.stringify(state))
  } catch {
    // 忽略存储失败
  }

  router.push({
    path: `/public/blog/${id}`
  })
}

onMounted(() => {
  // 优先从路由 query 恢复查询参数
  const {page, size, keyword} = route.query
  if (page) {
    const p = Number(page)
    if (!Number.isNaN(p) && p > 0) {
query.page = p
}
  }
  if (size) {
    const s = Number(size)
    if (!Number.isNaN(s) && s > 0) {
query.size = s
}
  }
  if (typeof keyword === 'string') {
    query.keyword = keyword
  }

  // 再尝试从 sessionStorage 恢复（例如从详情返回）
  try {
    const raw = sessionStorage.getItem(SCROLL_STATE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw) as {
        page?: number;
        size?: number;
        keyword?: string;
        scrollTop?: number;
        ts?: number;
      }
      const now = Date.now()
      if (!parsed.ts || now - parsed.ts < 60 * 60 * 1000) {
        if (!page && parsed.page && parsed.page > 0) {
          query.page = parsed.page
        }
        if (!size && parsed.size && parsed.size > 0) {
          query.size = parsed.size
        }
        if (!keyword && typeof parsed.keyword === 'string') {
          query.keyword = parsed.keyword
        }
        if (typeof parsed.scrollTop === 'number' && parsed.scrollTop > 0) {
          pendingScrollTop.value = parsed.scrollTop
        }
      } else {
        sessionStorage.removeItem(SCROLL_STATE_KEY)
      }
    }
  } catch {
    // 忽略解析失败
  }

  updateRouteQuery()
  loadData()

  // 初始化移动端检测
  checkMobile()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped lang="scss">
@import '@/styles/public-list.scss';

// 博客列表页特定样式（如需覆盖通用样式，可在此添加）
.blog-list-page {
  // 保留原有类名以兼容可能的其他引用
}
</style>

