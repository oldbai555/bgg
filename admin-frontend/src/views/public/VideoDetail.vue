<template>
  <div class="video-detail-page public-detail-page">
    <MetricReporter module="video_detail" :biz-id="Number(route.params.id) || 0" />
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
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import type DPlayerType from 'dplayer'
import type {PublicVideoDetailResp} from '@/api/generated/admin'
import {videoApi} from '@/api/video'
import {metricApi} from '@/api/metric'
import {copyToClipboard as copyText} from '@/utils/clipboard'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'

const router = useRouter()
const route = useRoute()

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
let DPlayer: typeof DPlayerType | null = null
let isInitializing = false
const hasReportedPlay = ref(false)

// 获取 API base URL
const getApiBase = (): string => {
  if (import.meta.env.PROD) {
    return '/gateway'
  }
  return 'http://localhost:20000'
}

// 动态加载 DPlayer
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

  if (isInitializing || loading.value) {
    return
  }

  const playUrl = video.value.playUrl
  if (!playUrl || playUrl.trim() === '') {
    return
  }

  try {
    new URL(playUrl)
  } catch {
    console.warn('无效的视频 URL:', playUrl)
    return
  }

  if (player && player.video) {
    const currentSrc = player.video.src || player.video.currentSrc
    if (currentSrc && currentSrc === playUrl) {
      return
    }
  }

  isInitializing = true

  try {
    if (player) {
      try {
        player.destroy()
      } catch (e) {
        console.warn('销毁旧播放器失败:', e)
      }
      player = null
    }

    await nextTick()

    if (!dplayerRef.value) {
      isInitializing = false
      return
    }

    const isM3u8 = playUrl.includes('.m3u8')
    loading.value = true

    let finalUrl = playUrl
    if (isM3u8 && typeof window !== 'undefined') {
      try {
        const baseURL = getApiBase()
        finalUrl = `${baseURL}/api/v1/m3u8/proxy?url=${encodeURIComponent(playUrl)}`
      } catch {
        finalUrl = `http://localhost:20000/api/v1/m3u8/proxy?url=${encodeURIComponent(playUrl)}`
      }
    }

    await initDPlayer({
      url: finalUrl,
      type: isM3u8 ? 'hls' : 'auto'
    })
  } catch (err: unknown) {
    console.error('播放器初始化失败:', err)
    const message = err instanceof Error ? err.message : '播放器初始化失败'
    if (typeof window !== 'undefined') {
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

  const DPlayerClass = await loadDPlayer()
  if (!DPlayerClass) {
    console.error('DPlayer 加载失败')
    ElMessage.error('播放器加载失败')
    loading.value = false
    isInitializing = false
    return
  }

  if (player) {
    try {
      player.destroy()
    } catch (e) {
      console.warn('销毁旧播放器失败:', e)
    }
    player = null
  }

  try {
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
      hlsConfig: {
        xhrSetup: (xhr: XMLHttpRequest) => {
          xhr.withCredentials = false
        }
      }
    })

    await nextTick()
    
    if (player?.video) {
      const videoElement = player.video
      if (videoElement.readyState >= 2) {
        loading.value = false
        isInitializing = false
      } else {
        setTimeout(() => {
          if (videoElement.readyState >= 2) {
            loading.value = false
            isInitializing = false
          }
        }, 1000)
      }
    } else {
      setTimeout(() => {
        loading.value = false
        isInitializing = false
      }, 1000)
    }

    if (player.video) {
      const videoElement = player.video
      const loadingTimeout = setTimeout(() => {
        if (loading.value) {
          console.warn('视频加载超时')
          loading.value = false
          ElMessage.warning('视频加载超时，请检查网络连接或视频链接')
        }
      }, 10000)

      videoElement.addEventListener('loadstart', () => {
        loading.value = true
      })

      videoElement.addEventListener('canplay', () => {
        clearTimeout(loadingTimeout)
        loading.value = false
      })

      videoElement.addEventListener('loadeddata', () => {
        clearTimeout(loadingTimeout)
        loading.value = false
      })

      videoElement.addEventListener('loadedmetadata', () => {
        clearTimeout(loadingTimeout)
      })

      videoElement.addEventListener('playing', () => {
        clearTimeout(loadingTimeout)
        loading.value = false

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

      videoElement.addEventListener('error', () => {
        clearTimeout(loadingTimeout)
        isInitializing = false
        loading.value = false
        ElMessage.error('视频加载失败，请检查视频链接')
      })

      videoElement.addEventListener('waiting', () => {
        loading.value = true
      })

      videoElement.addEventListener('canplaythrough', () => {
        clearTimeout(loadingTimeout)
        loading.value = false
        isInitializing = false
      })
    }
  } catch (error) {
    console.error('初始化播放器失败:', error)
    ElMessage.error('播放器初始化失败')
    loading.value = false
    throw error
  }
}

// 加载视频详情
const loadData = async () => {
  const id = route.params.id as string

  if (!id) {
    ElMessage.error('视频ID不能为空')
    router.push('/videos')
    return
  }

  const idNum = Number(id)
  if (!idNum || idNum === 0) {
    ElMessage.error('视频ID格式错误')
    router.push('/videos')
    return
  }

  loading.value = true
  try {
    const resp = await videoApi.publicDetail({id: idNum})
    video.value = resp
    hasReportedPlay.value = false
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    ElMessage.error(message)
    router.push('/videos')
  } finally {
    loading.value = false
  }
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
  if (typeof window !== 'undefined' && window.history.length > 1) {
    router.back()
  } else {
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
    router.push('/videos')
  }
}

// 监听 video.playUrl 变化
watch(
  () => video.value.playUrl,
  (playUrl) => {
    if (playUrl && playUrl.trim() !== '' && dplayerRef.value && !loading.value && typeof window !== 'undefined') {
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
    if (newId && newId !== oldId && typeof window !== 'undefined') {
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
      if (player) {
        try {
          player.destroy()
        } catch (e) {
          console.warn('销毁旧播放器失败:', e)
        }
        player = null
      }
      await loadData()
    }
  },
  {immediate: false}
)

onMounted(() => {
  if (typeof window !== 'undefined') {
    loadData()
    if (video.value.playUrl) {
      setTimeout(() => {
        initPlayer()
      }, 200)
    }
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
@import '@/styles/public-detail.scss';

.video-detail-page {
  .video-container {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;
    position: relative;

    &.is-loading {
      min-height: 400px;

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

    @media (max-width: 768px) {
      :deep(.dplayer-loading),
      :deep(.dplayer-loading-icon) {
        display: none !important;
      }

      &.is-loading {
        &::before,
        &::after {
          display: block;
        }
      }

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
