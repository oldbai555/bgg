<template>
  <div class="video-detail-page public-detail-page">
    <ClientOnly>
      <MetricReporter module="video_detail" :biz-id="Number(route.params.id) || 0" />
    </ClientOnly>
    <div class="container">
      <!-- 返回按钮 -->
      <div class="back-link" @click="goBack">← 返回列表</div>

      <!-- 视频播放器 -->
      <div :class="['video-container', { 'is-loading': loading }]">
        <div class="video-wrapper">
          <div ref="dplayerRef" class="dplayer-container"></div>
        </div>

        <div class="video-info-section">
          <h1 class="video-title">{{ video.name || '未命名视频' }}</h1>
          <div class="video-meta">
            <span class="video-code">{{ video.godNum || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- 磁力链接 -->
      <div class="magnet-section">
        <h2 class="section-title">🧲 磁力链接</h2>
        <div v-if="video.xlzzUrls && video.xlzzUrls.length > 0" class="magnet-list">
          <div
            v-for="(url, index) in video.xlzzUrls"
            :key="index"
            class="magnet-item"
            @click="copyToClipboard(url)"
          >
            <div class="magnet-icon">{{ index + 1 }}</div>
            <div class="magnet-text">{{ url }}</div>
            <div class="copy-icon">📋</div>
          </div>
        </div>
        <div v-else class="empty-message">暂无磁力链接</div>
      </div>
    </div>

    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
// Nuxt 3 自动导入 composables，无需手动导入 useRouter、useRoute
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
// DPlayer 只能在客户端使用，使用动态导入避免 SSR 错误
import type DPlayerType from 'dplayer'
import type {PublicVideoDetailResp} from '@/api/generated/admin'
import {videoApi} from '@/api/video'
import {metricApi} from '@/api/metric'
import {copyToClipboard as copyText} from '@/utils/clipboard'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'

// Nuxt 3 自动导入 useRouter 和 useRoute
const router = useRouter()
const route = useRoute()

// 定义页面元数据（Nuxt 3 规范）
definePageMeta({
  layout: false
})

const video = ref<PublicVideoDetailResp>({
  id: 0,
  uuid: '',
  name: '',
  godNum: '',
  playUrl: '',
  xlzzUrls: [],
  cover: '',
  duration: 0,
  description: '',
  type: 0,
  createdAt: 0,
  updatedAt: 0
})
const loading = ref(false)
const dplayerRef = ref<HTMLDivElement | null>(null)
let player: DPlayerType | null = null
let DPlayer: typeof DPlayerType | null = null // 动态导入的 DPlayer 类
let isInitializing = false // 防止重复初始化的标志
const hasReportedPlay = ref(false)

// 动态加载 DPlayer（只在客户端）
const loadDPlayer = async (): Promise<typeof DPlayerType | null> => {
  if (typeof window === 'undefined') {
    return null
  }
  
  if (DPlayer) {
    return DPlayer
  }
  
  try {
    const dplayerModule = await import('dplayer')
    DPlayer = dplayerModule.default
    return DPlayer
  } catch (err) {
    console.error('加载 DPlayer 失败:', err)
    return null
  }
}

// 初始化播放器
const initPlayer = async () => {
  if (!dplayerRef.value || typeof window === 'undefined') {
    return
  }

  // 如果正在初始化或还在加载中，跳过
  if (isInitializing || loading.value) {
    return
  }

  const playUrl = video.value.playUrl
  if (!playUrl || playUrl.trim() === '') {
    return
  }

  // 验证 URL 格式
  try {
    new URL(playUrl)
  } catch {
    console.warn('无效的视频 URL:', playUrl)
    return
  }

  // 如果播放器已存在且 URL 相同，跳过
  if (player && player.video) {
    const currentSrc = player.video.src || player.video.currentSrc
    if (currentSrc && currentSrc === playUrl) {
      return
    }
  }

  isInitializing = true

  try {
    // 销毁旧的播放器
    if (player) {
      try {
        player.destroy()
      } catch (e) {
        console.warn('销毁旧播放器失败:', e)
      }
      player = null
    }

    // 等待 DOM 更新
    await nextTick()

    if (!dplayerRef.value) {
      isInitializing = false
      return
    }

    const isM3u8 = playUrl.includes('.m3u8')

    // 设置加载状态
    loading.value = true

    // 对于 m3u8 视频，使用后端代理来避免 CORS 问题
    // 后端代理支持 Range 请求（边下边播），hls.js 会自动利用此功能
    let finalUrl = playUrl
    if (isM3u8 && typeof window !== 'undefined') {
      // 使用后端代理接口（支持 Range 请求，实现边下边播）
      try {
        const runtimeConfig = useRuntimeConfig()
        const baseURL = runtimeConfig.public.apiBase || 'http://localhost:20000'
        finalUrl = `${baseURL}/api/v1/m3u8/proxy?url=${encodeURIComponent(playUrl)}`
      } catch {
        // 如果获取失败，使用默认值
        finalUrl = `http://localhost:20000/api/v1/m3u8/proxy?url=${encodeURIComponent(playUrl)}`
      }
    }

    // 直接播放，失败就失败
    await initDPlayer({
      url: finalUrl,
      type: isM3u8 ? 'hls' : 'auto'
    })
  } catch (err: unknown) {
    console.error('播放器初始化失败:', err)
    const message = err instanceof Error ? err.message : '播放器初始化失败'
    if (process.client) {
      ElMessage.error(message)
    }
    loading.value = false
    isInitializing = false
  }
}

// 初始化 DPlayer
const initDPlayer = async (options: {url: string; type: 'hls' | 'auto'}) => {
  if (!dplayerRef.value || typeof window === 'undefined') {
    return
  }

  // 动态加载 DPlayer
  const DPlayerClass = await loadDPlayer()
  if (!DPlayerClass) {
    console.error('DPlayer 加载失败')
    ElMessage.error('播放器加载失败')
    loading.value = false
    isInitializing = false
    return
  }

  // 销毁旧播放器
  if (player) {
    try {
      player.destroy()
    } catch (e) {
      console.warn('销毁旧播放器失败:', e)
    }
    player = null
  }

  try {
    const hlsConfig: Record<string, unknown> = {}

    // 由于已经使用后端代理，不需要自定义 loader
    // 后端代理会处理 CORS 问题

    player = new DPlayerClass({
      container: dplayerRef.value,
      video: {
        url: options.url,
        type: options.type,
        pic: video.value.cover || undefined
      },
      autoplay: false,
      theme: '#b7daff',
      loop: false,
      lang: 'zh-cn',
      screenshot: true,
      hotkey: true,
      preload: 'auto',
      volume: 0.7,
      mutex: true,
      playbackSpeed: [0.5, 0.75, 1, 1.25, 1.5, 2],
      hlsConfig: Object.keys(hlsConfig).length > 0 ? hlsConfig : {
        // HLS 配置
        // 后端代理已支持 Range 请求（边下边播），hls.js 会自动发送 Range 请求
        // 无需额外配置，播放器会自动利用 Range 请求实现边下边播
        xhrSetup: (xhr: XMLHttpRequest) => {
          xhr.withCredentials = false
          // Range 请求由 hls.js 自动处理，无需手动设置
        }
      }
    })

    // 播放器创建完成后，检查并关闭 loading
    await nextTick()
    
    // 立即检查一次
    if (player?.video) {
      const videoElement = player.video
      // 如果视频元素已经有数据，立即关闭 loading
      if (videoElement.readyState >= 2) {
        loading.value = false
        isInitializing = false
      } else {
        // 否则延迟检查（给播放器一些时间加载，特别是移动端）
        setTimeout(() => {
          if (videoElement.readyState >= 2) {
            loading.value = false
            isInitializing = false
          }
        }, 1000)
      }
    } else {
      // 如果播放器创建失败，延迟关闭 loading
      setTimeout(() => {
        loading.value = false
        isInitializing = false
      }, 1000)
    }

    // 监听播放器事件
    if (player.video) {
      const videoElement = player.video

      // 设置超时检测（10秒后如果还在加载，可能是加载失败）
      const loadingTimeout = setTimeout(() => {
        if (loading.value) {
          console.warn('视频加载超时，可能加载失败')
          console.log('视频元素状态:', {
            readyState: videoElement.readyState,
            networkState: videoElement.networkState,
            error: videoElement.error,
            src: videoElement.src,
            currentSrc: videoElement.currentSrc
          })
          loading.value = false
          ElMessage.warning('视频加载超时，请检查网络连接或视频链接')
        }
      }, 10000)

      // 开始加载
      videoElement.addEventListener('loadstart', () => {
        loading.value = true
      })

      // 可以播放
      videoElement.addEventListener('canplay', () => {
        console.log('视频可以播放')
        clearTimeout(loadingTimeout)
        loading.value = false
      })

      // 可以开始播放（更早的事件）
      videoElement.addEventListener('loadeddata', () => {
        console.log('视频数据已加载')
        clearTimeout(loadingTimeout)
        loading.value = false
      })

      // 元数据已加载
      videoElement.addEventListener('loadedmetadata', () => {
        console.log('视频元数据已加载，时长:', videoElement.duration)
        clearTimeout(loadingTimeout)
      })

      // 播放中
      videoElement.addEventListener('playing', () => {
        console.log('视频正在播放')
        clearTimeout(loadingTimeout)
        loading.value = false

        // 上报播放事件（每次进入详情页仅上报一次，避免频繁触发）
        if (!hasReportedPlay.value && video.value.id) {
          hasReportedPlay.value = true
          metricApi
            .report({
              module: 'video_detail',
              bizId: video.value.id,
              event: 'play'
            })
            .catch(() => {})
        }
      })

      // 错误处理
      videoElement.addEventListener('error', (_e: Event) => {
        clearTimeout(loadingTimeout)
        isInitializing = false
        const error = videoElement.error
        if (error) {
          // 忽略 "Empty src attribute" 错误，这通常是播放器切换时的临时状态
          if (error.message && error.message.includes('Empty src')) {
            return
          }
          console.error('视频错误详情:', {
            code: error.code,
            message: error.message,
            readyState: videoElement.readyState,
            networkState: videoElement.networkState,
            src: videoElement.src,
            currentSrc: videoElement.currentSrc
          })
        }
        loading.value = false
        ElMessage.error('视频加载失败，请检查视频链接')
      })

      // 等待中
      videoElement.addEventListener('waiting', () => {
        loading.value = true
      })

      // 可以继续播放
      videoElement.addEventListener('canplaythrough', () => {
        clearTimeout(loadingTimeout)
        loading.value = false
        isInitializing = false
      })

      // 移动端特殊处理：确保在视频可以播放时关闭 loading
      // 有些移动端浏览器可能不会触发某些事件，使用更早的事件来关闭 loading
      videoElement.addEventListener('loadedmetadata', () => {
        // 元数据加载完成后，延迟关闭 loading（给播放器一些时间渲染）
        setTimeout(() => {
          if (videoElement.readyState >= 2) { // HAVE_CURRENT_DATA
            loading.value = false
          }
        }, 500)
      })

      // 移动端特殊处理：定期检查并关闭 loading
      // 有些移动端浏览器可能不会触发某些事件，使用定期检查确保 loading 关闭
      const loadingCheckInterval = setInterval(() => {
        // 如果视频元素已准备好（有元数据或数据），关闭 loading
        if (videoElement.readyState >= 2) { // HAVE_CURRENT_DATA
          if (loading.value) {
            console.log('通过定期检查关闭 loading，readyState:', videoElement.readyState)
            loading.value = false
            isInitializing = false
          }
          // 如果 readyState >= 3，说明数据已足够，可以停止检查
          if (videoElement.readyState >= 3) {
            clearInterval(loadingCheckInterval)
          }
        }
      }, 500) // 每 500ms 检查一次
      
      // 5秒后强制关闭 loading 和检查（移动端可能加载较慢，但不应超过5秒）
      setTimeout(() => {
        clearInterval(loadingCheckInterval)
        if (loading.value) {
          console.log('超时强制关闭 loading')
          loading.value = false
          isInitializing = false
        }
      }, 5000)
    }
  } catch (error) {
    console.error('初始化播放器失败:', error)
    ElMessage.error('播放器初始化失败')
    loading.value = false
    throw error
  }
}

// 加载视频详情（使用 Nuxt 3 的 useAsyncData 支持 SSR）
const { data: videoData, error: videoError, pending, refresh } = await useAsyncData(
  `video-${route.params.id}`,
  async () => {
    const id = route.params.id as string

    if (!id) {
      throw new Error('视频ID不能为空')
    }

    const idNum = Number(id)
    if (!idNum || idNum === 0) {
      throw new Error('视频ID格式错误')
    }

    try {
      const resp = await videoApi.publicDetail({id: idNum})
      return resp
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '加载失败'
      throw new Error(message)
    }
  },
  {
    server: true, // 允许在服务端执行
    default: () => ({
      id: 0,
      uuid: '',
      name: '',
      godNum: '',
      playUrl: '',
      xlzzUrls: [],
      cover: '',
      duration: 0,
      description: '',
      type: 0,
      createdAt: 0,
      updatedAt: 0
    } as PublicVideoDetailResp)
  }
)

// 同步 videoData 到 video ref
watch(videoData, (newData) => {
  if (newData) {
    video.value = newData
    hasReportedPlay.value = false
  }
}, { immediate: true })

// 同步 loading 状态（只在数据加载时，播放器加载状态由播放器事件控制）
watch(pending, (isPending) => {
  // 只在数据加载时设置 loading，播放器初始化后由播放器事件控制
  if (isPending) {
    loading.value = true
  } else {
    // 数据加载完成后，如果播放器还未初始化，等待播放器初始化
    // 如果播放器已初始化，由播放器事件控制 loading 状态
    if (!player && !isInitializing) {
      // 数据已加载但播放器未初始化，等待播放器初始化
      // loading 状态由播放器初始化过程控制
    }
  }
}, { immediate: true })

// 处理错误
watch(videoError, (err) => {
  if (err) {
    console.error('加载视频详情失败:', err)
    if (process.client) {
      ElMessage.error(err.message || '加载失败')
      router.push('/videos')
    }
  }
}, { immediate: true })

// 兼容旧的 loadData 函数（用于路由参数变化时）
const loadData = async () => {
  const id = route.params.id as string

  if (!id) {
    if (process.client) {
      ElMessage.error('视频ID不能为空')
      router.push('/videos')
    }
    return
  }

  const idNum = Number(id)
  if (!idNum || idNum === 0) {
    if (process.client) {
      ElMessage.error('视频ID格式错误')
      router.push('/videos')
    }
    return
  }

  // 刷新数据
  await refresh()
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  const success = await copyText(text)
  if (success) {
    ElMessage.success('已复制到剪贴板 ✓')
  } else {
    ElMessage.error('复制失败，请手动复制')
  }
}

// 返回列表
const goBack = () => {
  // 优先走浏览器历史，保留列表分页与滚动状态
  // 如果是从列表页进入的，router.back() 会恢复列表页的状态（包括滚动位置）
  if (typeof window !== 'undefined' && window.history.length > 1) {
    router.back()
  } else {
    // 如果没有历史记录，尝试从 sessionStorage 恢复状态
    if (typeof window !== 'undefined') {
      try {
        const raw = sessionStorage.getItem('public_video_list_state')
        if (raw) {
          const parsed = JSON.parse(raw) as {
            page?: number
            size?: number
            content?: string
          }
          router.push({
            path: '/videos',
            query: {
              ...(parsed.page && {page: String(parsed.page)}),
              ...(parsed.size && {size: String(parsed.size)}),
              ...(parsed.content && {content: parsed.content})
            }
          })
          return
        }
      } catch {
        // 忽略解析错误
      }
    }
    // 兜底：直接跳转到列表页
    router.push('/videos')
  }
}

// 监听 video.playUrl 和 loading 状态，确保在合适的时机初始化播放器
watch(
    () => video.value.playUrl,
    (playUrl) => {
      // 只有当 playUrl 存在且有效时才初始化播放器
      if (playUrl && playUrl.trim() !== '' && dplayerRef.value && !loading.value && typeof window !== 'undefined') {
        // 延迟初始化，确保 DOM 已完全更新
        setTimeout(() => {
          initPlayer()
        }, 200)
      }
    },
    {immediate: false}
)

// 监听路由参数变化
watch(
  () => route.params.id,
  async (newId, oldId) => {
    if (newId && newId !== oldId && process.client) {
      // 重置状态
      video.value = {
        id: 0,
        uuid: '',
        name: '',
        godNum: '',
        playUrl: '',
        xlzzUrls: [],
        cover: '',
        duration: 0,
        description: '',
        type: 0,
        createdAt: 0,
        updatedAt: 0
      }
      hasReportedPlay.value = false
      // 销毁旧播放器
      if (player) {
        try {
          player.destroy()
        } catch (e) {
          console.warn('销毁旧播放器失败:', e)
        }
        player = null
      }
      // 重新加载
      await loadData()
    }
  },
  {immediate: false}
)

// 客户端挂载时初始化播放器
onMounted(() => {
  if (process.client && video.value.playUrl) {
    // 延迟初始化，确保 DOM 已完全更新
    setTimeout(() => {
      initPlayer()
    }, 200)
  }
})

onBeforeUnmount(() => {
  if (player) {
    try {
      player.destroy()
    } catch (e) {
      console.warn('销毁播放器失败:', e)
    }
    player = null
  }
})
</script>

<style scoped lang="scss">
@import '@/assets/styles/public-detail.scss';

// 视频详情页特定样式
.video-detail-page {

  .video-container {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;
    position: relative;

    &.is-loading {
      min-height: 400px; // 确保加载时容器有足够高度

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: rgba(255, 255, 255, 0.8);
        z-index: 999;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      &::after {
        content: '';
        position: absolute;
        top: 50%;
        left: 50%;
        width: 30px;
        height: 30px;
        margin-top: -15px;
        margin-left: -15px;
        border: 3px solid #409eff;
        border-radius: 50%;
        border-top-color: transparent;
        animation: spin 1s linear infinite;
        z-index: 1000;
      }
    }

    // 移动端：隐藏 DPlayer 自己的加载状态（如果有）
    @media (max-width: 768px) {
      // 隐藏 DPlayer 控制栏中的加载指示器
      :deep(.dplayer-loading),
      :deep(.dplayer-loading-icon) {
        display: none !important;
      }

      // 确保我们的 loading 状态在移动端正确显示和隐藏
      &.is-loading {
        &::before,
        &::after {
          display: block;
        }
      }

      // 播放器准备好后，确保 loading 隐藏
      &:not(.is-loading) {
        &::before,
        &::after {
          display: none !important;
        }
      }
    }
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .video-wrapper {
    position: relative;
    width: 100%;
    padding-top: 56.25%; /* 16:9 aspect ratio */
    background: #000;
    overflow: hidden;
  }

  .dplayer-container {
    position: absolute !important;
    top: 0 !important;
    left: 0 !important;
    width: 100% !important;
    height: 100% !important;
  }

  .video-info-section {
    padding: 25px;
  }

  .video-title {
    font-size: 24px;
    font-weight: 600;
    color: #333;
    margin-bottom: 15px;
  }

  .video-meta {
    display: flex;
    align-items: center;
    gap: 15px;
    color: #666;
    font-size: 14px;
    margin-bottom: 25px;
    padding-bottom: 20px;
    border-bottom: 1px solid #e0e0e0;
  }

  .video-code {
    background: #f5f5f5;
    padding: 6px 12px;
    border-radius: 6px;
    font-weight: 500;
  }

  .magnet-section {
    background: white;
    border-radius: 12px;
    padding: 25px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  }

  .section-title {
    font-size: 20px;
    font-weight: 600;
    color: #333;
    margin-bottom: 15px;
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .magnet-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .magnet-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 15px;
    background: #f8f9fa;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    border: 2px solid transparent;

    &:hover {
      background: #e9ecef;
      border-color: #667eea;
      transform: translateX(5px);
    }
  }

  .magnet-icon {
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: 600;
  }

  .magnet-text {
    flex: 1;
    font-family: monospace;
    font-size: 13px;
    color: #333;
    word-break: break-all;
    line-height: 1.5;
  }

  .copy-icon {
    flex-shrink: 0;
    width: 36px;
    height: 36px;
    background: #667eea;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 18px;
    transition: background 0.3s;
  }

  .magnet-item:hover .copy-icon {
    background: #5568d3;
  }

  .empty-message {
    text-align: center;
    color: #999;
    padding: 40px;
    font-size: 16px;
  }

  @media (max-width: 768px) {
    .video-info-section {
      padding: 18px 14px 20px;
    }

    .video-title {
      font-size: 20px;
      line-height: 1.3;
    }

    .video-meta {
      flex-wrap: wrap;
      gap: 8px;
      font-size: 13px;
    }

    .magnet-section {
      padding: 18px 14px 20px;
    }

    .magnet-item {
      align-items: flex-start;
      padding: 12px;
    }

    .magnet-text {
      font-size: 12px;
    }
  }
}
</style>
