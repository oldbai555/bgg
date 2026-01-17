<template>
  <div class="video-list-page public-list-page">
    <MetricReporter module="video_list" :biz-id="0" />
    <div class="container">
      <div class="hero">
        <div class="hero-title">🎬 视频列表</div>
        <div class="hero-subtitle">发现精彩视频内容</div>

        <!-- 搜索栏 -->
        <div class="search-bar">
          <el-input
            v-model="query.content"
            placeholder="搜索视频..."
            clearable
            @keydown.enter="handleSearch"
          >
            <template #append>
              <el-button type="primary" :loading="loading" @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
      </div>

      <!-- 视频网格 -->
      <div v-loading="loading" class="list-grid">
        <div
          v-for="video in list"
          :key="video.id"
          class="list-card video-card"
        >
          <div
            class="cover video-thumbnail"
            @mouseenter="handleThumbnailHover(video)"
            @mouseleave="handleThumbnailLeave(video)"
          >
            <img
              v-if="!hoveringVideoId || hoveringVideoId !== video.id"
              :src="getCoverUrl(video.godNum)"
              :alt="video.name"
              @error="handleImageError"
            />
            <video
              v-else
              :ref="(el) => setVideoRef(video.id, el)"
              class="video-player-inline"
              :src="getPreviewUrl(video.godNum)"
              autoplay
              loop
              muted
              playsinline
            ></video>
            <div class="play-overlay">
              <div class="play-icon"></div>
            </div>
          </div>
          <div class="card-content video-info">
            <div
              class="card-title video-title"
              @click.stop="goToDetail(video.id)"
            >
              {{ video.name || '未命名视频' }}
            </div>
            <div class="card-meta video-code">{{ video.godNum || '-' }}</div>
          </div>
        </div>
        <div v-if="!loading && list.length === 0" class="empty-message">
          暂无视频数据
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          :page-sizes="[10, 20, 30, 50, 100]"
          :layout="paginationLayout"
          :size="isMobile ? 'small' : 'default'"
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
import {videoApi} from '@/api/video'
import type {PublicVideoListReq, PublicVideoItem} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'

const router = useRouter()
const route = useRoute()

const SCROLL_STATE_KEY = 'public_video_list_state'

const query = reactive({
  page: 1,
  size: 10,
  content: ''
})
const list = ref<PublicVideoItem[]>([])
const total = ref(0)
const loading = ref(false)
const hoveringVideoId = ref<number | null>(null)
const videoRefs = ref<Map<number, HTMLVideoElement>>(new Map())
const pendingScrollTop = ref<number | null>(null)
const isMobile = ref(false)

const paginationLayout = computed(() => {
  return isMobile.value ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'
})

const checkMobile = () => {
  if (typeof window !== 'undefined') {
    isMobile.value = window.innerWidth <= 768
  }
}

const handleResize = () => {
  checkMobile()
}

const getCoverUrl = (godNum: string): string => {
  if (!godNum) {
    return ''
  }
  return `https://fourhoi.com/${godNum}/cover-t.jpg`
}

const getPreviewUrl = (godNum: string): string => {
  if (!godNum) {
    return ''
  }
  return `https://fourhoi.com/${godNum}/preview.mp4`
}

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjgwIiBoZWlnaHQ9IjE1OCIgdmlld0JveD0iMCAwIDI4MCAxNTgiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyODAiIGhlaWdodD0iMTU4IiBmaWxsPSIjRjVGNUY1Ii8+CjxwYXRoIGQ9Ik0xNDAgNzkuNUwxMzAgODkuNUwxMzAgNjkuNUwxNDAgNzkuNVoiIGZpbGw9IiNDQ0NDQ0MiLz4KPHBhdGggZD0iTTE0MCA3OS41TDE1MCA4OS41TDE1MCA2OS41TDE0MCA3OS41WiIgZmlsbD0iI0NDQ0NDQyIvPgo8L3N2Zz4='
}

const handleThumbnailHover = (video: PublicVideoItem) => {
  hoveringVideoId.value = video.id
}

const setVideoRef = (videoId: number, el: unknown) => {
  if (el && el instanceof HTMLVideoElement) {
    videoRefs.value.set(videoId, el)
  } else {
    videoRefs.value.delete(videoId)
  }
}

const handleThumbnailLeave = (video: PublicVideoItem) => {
  hoveringVideoId.value = null
  const videoEl = videoRefs.value.get(video.id)
  if (videoEl) {
    videoEl.pause()
    videoEl.currentTime = 0
  }
}

const restoreScrollPosition = async (scrollTop: number) => {
  if (typeof window === 'undefined') {
    return
  }
  await nextTick()
  await new Promise(resolve => requestAnimationFrame(resolve))
  await new Promise(resolve => requestAnimationFrame(resolve))
  await new Promise(resolve => setTimeout(resolve, 100))
  window.scrollTo({top: scrollTop, behavior: 'auto'})
  setTimeout(() => {
    const currentScroll = window.scrollY
    const diff = Math.abs(currentScroll - scrollTop)
    if (diff > 50) {
      window.scrollTo({top: scrollTop, behavior: 'auto'})
    }
  }, 200)
}

const loadData = async () => {
  loading.value = true
  const shouldRestoreScroll = pendingScrollTop.value !== null
  const scrollTopToRestore = pendingScrollTop.value

  try {
    const req: PublicVideoListReq = {
      page: query.page,
      size: query.size
    }
    if (query.content) {
      req.content = query.content
    }

    const resp = await videoApi.publicList(req)
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

const handleSearch = () => {
  query.page = 1
  updateRouteQuery()
  pendingScrollTop.value = null
  if (typeof window !== 'undefined') {
    try {
      sessionStorage.removeItem(SCROLL_STATE_KEY)
    } catch {
      // 忽略清除失败
    }
  }
  loadData()
}

const handlePageChange = (page: number) => {
  query.page = page
  updateRouteQuery()
  pendingScrollTop.value = null
  loadData()
  if (typeof window !== 'undefined') {
    window.scrollTo({top: 0, behavior: 'smooth'})
  }
}

const handleSizeChange = (size: number) => {
  query.size = size
  query.page = 1
  updateRouteQuery()
  pendingScrollTop.value = null
  loadData()
}

const goToDetail = (id: number) => {
  if (typeof window !== 'undefined') {
    try {
      const state = {
        page: query.page,
        size: query.size,
        content: query.content,
        scrollTop: window.scrollY,
        ts: Date.now()
      }
      sessionStorage.setItem(SCROLL_STATE_KEY, JSON.stringify(state))
    } catch {
      // 忽略存储失败
    }
  }

  router.push({
    path: `/videos/${id}`,
    query: {
      page: String(query.page),
      size: String(query.size),
      content: query.content || undefined
    }
  })
}

const updateRouteQuery = () => {
  router.replace({
    path: route.path,
    query: {
      ...route.query,
      page: String(query.page),
      size: String(query.size),
      content: query.content || undefined
    }
  })
}

onMounted(() => {
  const page = route.query.page
  const size = route.query.size
  const content = route.query.content

  if (page) {
    query.page = Number(page)
  }
  if (size) {
    query.size = Number(size)
  }
  if (content) {
    query.content = String(content)
  }

  if (typeof window !== 'undefined') {
    try {
      const raw = sessionStorage.getItem(SCROLL_STATE_KEY)
      if (raw) {
        const parsed = JSON.parse(raw) as {
          page?: number
          size?: number
          content?: string
          scrollTop?: number
          ts?: number
        }
        const now = Date.now()
        if (!parsed.ts || now - parsed.ts < 60 * 60 * 1000) {
          if (!page && parsed.page && parsed.page > 0) {
            query.page = parsed.page
          }
          if (!size && parsed.size && parsed.size > 0) {
            query.size = parsed.size
          }
          if (!content && typeof parsed.content === 'string') {
            query.content = parsed.content
          }
          if (page || size) {
            if (typeof parsed.scrollTop === 'number' && parsed.scrollTop > 0) {
              pendingScrollTop.value = parsed.scrollTop
            }
          }
        } else {
          sessionStorage.removeItem(SCROLL_STATE_KEY)
        }
      }
    } catch {
      // 忽略解析错误
    }
  }

  updateRouteQuery()
  loadData()

  checkMobile()
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleResize)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleResize)
  }
})
</script>

<style scoped lang="scss">
@import '@/styles/public-list.scss';

.video-list-page {
  background: linear-gradient(135deg, #fff7e6 0%, #ffe9d9 45%, #ffd1a4 100%);

  .video-card {
    transition: transform 0.25s, box-shadow 0.25s;

    &:hover {
      transform: translateY(-6px);
      box-shadow: 0 12px 48px rgba(0, 0, 0, 0.12);
    }
  }

  .video-thumbnail {
    position: relative;
    padding-top: 56.25%; /* 16:9 */
    background: #f0f0f0;
    overflow: hidden;
    border-radius: 10px 0 0 10px;

    img,
    video {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.3s;
    }

    img {
      transform: scale(1);
    }

    .video-card:hover & img {
      transform: scale(1.06);
    }

    .video-player-inline {
      display: block;
      background: #000;
      border: 0;
    }
  }

  .play-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.28);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.2s;
  }

  .video-card:hover .play-overlay {
    opacity: 1;
  }

  .play-icon {
    width: 60px;
    height: 60px;
    background: rgba(255, 255, 255, 0.92);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;

    &::after {
      content: '';
      border-left: 20px solid #667eea;
      border-top: 12px solid transparent;
      border-bottom: 12px solid transparent;
      margin-left: 4px;
    }
  }

  .video-info {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .video-title {
    cursor: pointer;
    position: relative;
    z-index: 1;
    flex-shrink: 0;

    &:hover {
      color: #667eea;
    }
  }

  .video-code {
    position: relative;
    z-index: 0;
    flex-shrink: 0;
  }

  @media (max-width: 768px) {
    .video-thumbnail {
      border-radius: 10px 0 0 10px;
    }
  }
}
</style>
