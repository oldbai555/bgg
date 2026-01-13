<template>
  <div class="video-list-page">
    <div class="container">
      <h1>🎬 视频列表</h1>

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

      <!-- 视频网格 -->
      <div v-loading="loading" class="video-grid">
        <div
          v-for="video in list"
          :key="video.id"
          class="video-card"
        >
          <div
            class="video-thumbnail"
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
          <div class="video-info">
            <div
              class="video-title"
              @click.stop="goToDetail(video.id)"
            >
              {{ video.name || '未命名视频' }}
            </div>
            <div class="video-code">{{ video.godNum || '-' }}</div>
          </div>
        </div>
        <div v-if="!loading && list.length === 0" class="empty-message">
          暂无视频数据
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <div class="page-info">
          共 {{ total }} 条，第 {{ query.page }} / {{ totalPages }} 页
        </div>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, computed, onMounted, onUnmounted, nextTick} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {publicVideoList} from '@/api/generated/admin'
import type {PublicVideoListReq, PublicVideoItem} from '@/api/generated/admin'

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

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))

// 响应式分页 layout：移动端使用简化布局
const paginationLayout = computed(() => {
  return isMobile.value ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'
})

// 检测屏幕尺寸
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

// 监听窗口大小变化
const handleResize = () => {
  checkMobile()
}

// 获取封面URL
const getCoverUrl = (godNum: string): string => {
  if (!godNum) {
return ''
}
  return `https://fourhoi.com/${godNum}/cover-t.jpg`
}

// 获取预览视频URL
const getPreviewUrl = (godNum: string): string => {
  if (!godNum) {
return ''
}
  return `https://fourhoi.com/${godNum}/preview.mp4`
}

// 图片加载失败处理
const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjgwIiBoZWlnaHQ9IjE1OCIgdmlld0JveD0iMCAwIDI4MCAxNTgiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyODAiIGhlaWdodD0iMTU4IiBmaWxsPSIjRjVGNUY1Ii8+CjxwYXRoIGQ9Ik0xNDAgNzkuNUwxMzAgODkuNUwxMzAgNjkuNUwxNDAgNzkuNVoiIGZpbGw9IiNDQ0NDQ0MiLz4KPHBhdGggZD0iTTE0MCA3OS41TDE1MCA4OS41TDE1MCA2OS41TDE0MCA3OS41WiIgZmlsbD0iI0NDQ0NDQyIvPgo8L3N2Zz4='
}

// 鼠标悬停视频卡片
const handleThumbnailHover = (video: PublicVideoItem) => {
  hoveringVideoId.value = video.id
}

// 设置视频引用
const setVideoRef = (videoId: number, el: unknown) => {
  if (el && el instanceof HTMLVideoElement) {
    videoRefs.value.set(videoId, el)
  } else {
    videoRefs.value.delete(videoId)
  }
}

// 鼠标离开视频卡片
const handleThumbnailLeave = (video: PublicVideoItem) => {
  hoveringVideoId.value = null
  // 停止预览视频播放
  const videoEl = videoRefs.value.get(video.id)
  if (videoEl) {
    videoEl.pause()
    videoEl.currentTime = 0
  }
}

// 恢复滚动位置（在 DOM 完全渲染后）
const restoreScrollPosition = async (scrollTop: number) => {
  // 等待 Vue 完成 DOM 更新
  await nextTick()
  // 等待浏览器完成渲染
  await new Promise(resolve => requestAnimationFrame(resolve))
  await new Promise(resolve => requestAnimationFrame(resolve))
  // 额外等待，确保图片等资源加载完成
  await new Promise(resolve => setTimeout(resolve, 100))

  // 恢复滚动位置
  window.scrollTo({top: scrollTop, behavior: 'auto'})

  // 如果滚动位置仍然不对，可能是内容高度变化，尝试再次恢复
  setTimeout(() => {
    const currentScroll = window.scrollY
    const diff = Math.abs(currentScroll - scrollTop)
    if (diff > 50) {
      // 如果差异较大，再次尝试恢复
      window.scrollTo({top: scrollTop, behavior: 'auto'})
    }
  }, 200)
}

// 加载数据
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

    // 响应直接返回数据（无 code/msg 包装），拦截器会直接返回原始数据
    const resp = await publicVideoList(req)
    list.value = resp.list
    total.value = resp.total
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    ElMessage.error(message)
  } finally {
    loading.value = false

    // 如果需要恢复滚动位置，等待 DOM 完全渲染后再恢复
    if (shouldRestoreScroll && scrollTopToRestore !== null) {
      pendingScrollTop.value = null
      await restoreScrollPosition(scrollTopToRestore)
    }
  }
}

// 搜索
const handleSearch = () => {
  query.page = 1
  updateRouteQuery()
  // 清除待恢复的滚动位置（因为用户主动搜索了）
  pendingScrollTop.value = null
  // 清除 sessionStorage 中的旧状态
  try {
    sessionStorage.removeItem(SCROLL_STATE_KEY)
  } catch {
    // 忽略清除失败
  }
  loadData()
}

// 分页变化
const handlePageChange = (page: number) => {
  query.page = page
  updateRouteQuery()
  // 清除待恢复的滚动位置（因为用户主动切换了页面）
  pendingScrollTop.value = null
  loadData()
  // 滚动到顶部
  window.scrollTo({top: 0, behavior: 'smooth'})
}

// 每页数量变化
const handleSizeChange = (size: number) => {
  query.size = size
  query.page = 1
  updateRouteQuery()
  // 清除待恢复的滚动位置（因为用户主动改变了每页数量）
  pendingScrollTop.value = null
  loadData()
}

// 跳转到详情页
const goToDetail = (id: number) => {
  // 记录当前列表状态与滚动位置，便于返回时恢复
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

  router.push({
    path: `/public/videos/${id}`,
    query: {
      page: String(query.page),
      size: String(query.size),
      content: query.content || undefined
    }
  })
}

// 同步路由 query，便于刷新与跨页面返回
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

// 初始化：从路由参数获取查询条件
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

  // 尝试从 sessionStorage 中恢复状态（包括从详情页返回的情况）
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
      // 简单过期控制：1 小时内的记录才恢复
      if (!parsed.ts || now - parsed.ts < 60 * 60 * 1000) {
        // 如果路由参数中没有分页信息，使用 sessionStorage 中的
        if (!page && parsed.page && parsed.page > 0) {
          query.page = parsed.page
        }
        if (!size && parsed.size && parsed.size > 0) {
          query.size = parsed.size
        }
        if (!content && typeof parsed.content === 'string') {
          query.content = parsed.content
        }
        // 如果是从详情页返回（有路由 query），恢复滚动位置
        if (page || size) {
          if (typeof parsed.scrollTop === 'number' && parsed.scrollTop > 0) {
            pendingScrollTop.value = parsed.scrollTop
          }
        }
      } else {
        // 过期了，清除旧状态
        sessionStorage.removeItem(SCROLL_STATE_KEY)
      }
    }
  } catch {
    // 忽略解析错误
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
.video-list-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
  }

  h1 {
    color: #fff;
    text-align: center;
    margin-bottom: 30px;
    font-size: 2.4rem;
  }

  .search-bar {
    background: #fff;
    border-radius: 12px;
    padding: 18px;
    margin-bottom: 24px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  }

  .video-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 24px;
    margin-bottom: 32px;
  }

  .video-card {
    background: #fff;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
    transition: transform 0.25s, box-shadow 0.25s;
    cursor: pointer;

    &:hover {
      transform: translateY(-6px);
      box-shadow: 0 12px 48px rgba(0, 0, 0, 0.12);
    }
  }

  .video-thumbnail {
    position: relative;
    width: 100%;
    padding-top: 56.25%; /* 16:9 */
    background: #f0f0f0;
    overflow: hidden;
    border-top-left-radius: 12px;
    border-top-right-radius: 12px;

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
    padding: 14px;
  }

  .video-title {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin-bottom: 8px;
    line-height: 1.2;
    cursor: pointer;

    &:hover {
      color: #667eea;
    }
  }

  .video-code {
    font-size: 14px;
    color: #666;
  }

  .empty-message {
    grid-column: 1 / -1;
    text-align: center;
    color: #fff;
    padding: 30px;
    font-size: 16px;
  }

  .pagination {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    background: #fff;
    padding: 16px;
    border-radius: 10px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);

    .page-info {
      color: #666;
      font-size: 14px;
      white-space: nowrap;
    }

    // 确保分页组件在小屏下有足够空间
    :deep(.el-pagination) {
      flex-wrap: wrap;
      justify-content: center;
    }

    // 移动端隐藏部分元素，避免挤压
    :deep(.el-pagination__total),
    :deep(.el-pagination__sizes),
    :deep(.el-pagination__jump) {
      @media (max-width: 768px) {
        display: none;
      }
    }
  }

  @media (max-width: 768px) {
    padding: 12px;

    h1 {
      font-size: 1.6rem;
      margin-bottom: 20px;
    }

    .search-bar {
      padding: 12px;
      margin-bottom: 16px;
    }

    .video-grid {
      grid-template-columns: 1fr;
      gap: 14px;
      margin-bottom: 20px;
    }

    .video-title {
      font-size: 1.4rem;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .video-code {
      font-size: 12px;
    }

    .pagination {
      padding: 12px;
      gap: 8px;

      .page-info {
        font-size: 12px;
      }
    }
  }
}
</style>

