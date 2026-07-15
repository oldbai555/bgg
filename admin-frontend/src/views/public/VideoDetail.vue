<template>
  <div class="video-detail-page public-detail-page">
    <MetricReporter ref="metricReporterRef" module="video_detail" :biz-id="Number(route.params.id) || 0" />
    <PublicHeader />
    <div class="page-shell">
      <div class="page-layout">
        <div class="detail-card">
          <div class="back-link" @click="goBack">← 返回列表</div>

          <!-- 视频播放器 -->
          <div :class="['video-container', {'is-loading': loading}]">
            <div class="video-wrapper">
              <div ref="dplayerRef" class="dplayer-container"></div>
            </div>
          </div>

          <h1 class="title">{{ video.name || '未命名视频' }}</h1>
          <div class="meta">
            <span class="video-code">{{ video.godNum || '-' }}</span>
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
import {contentApi} from '@/api/content'
import {copyToClipboard as copyText} from '@/utils/clipboard'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import PublicHeader from '@/components/common/PublicHeader.vue'

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
const metricReporterRef = ref<InstanceType<typeof MetricReporter> | null>(null)
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
          metricReporterRef.value?.report({event: 'play', bizId: video.value.id})
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
    router.push('/front/videos')
    return
  }

  const idNum = Number(id)
  if (!idNum || idNum === 0) {
    ElMessage.error('视频ID格式错误')
    router.push('/front/videos')
    return
  }

  loading.value = true
  try {
    const resp = await contentApi.publicVideoDetail({id: idNum})
    video.value = resp
    hasReportedPlay.value = false
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    ElMessage.error(message)
    router.push('/front/videos')
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
            path: '/front/videos',
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
    router.push('/front/videos')
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

.video-detail-page .page-layout {
  grid-template-columns: 1fr;
  max-width: 860px;
  margin: 0 auto;
  width: 100%;
}

.video-detail-page .detail-card {
  .video-container {
    background: var(--color-bg-secondary);
    border-radius: 10px;
    overflow: hidden;
    margin-bottom: 20px;
    position: relative;

    &.is-loading {
      min-height: 300px;

      &::before {
        content: '';
        position: absolute;
        inset: 0;
        background-color: rgba(0, 0, 0, 0.05);
        z-index: 999;
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
        border: 3px solid var(--color-primary);
        border-radius: 50%;
        border-top-color: transparent;
        animation: spin 1s linear infinite;
        z-index: 1000;
      }
    }

    @include mobile {
      :deep(.dplayer-loading),
      :deep(.dplayer-loading-icon) {
        display: none !important;
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
    padding-top: 56.25%; // 16:9
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

  .video-code {
    background: var(--color-bg-secondary);
    padding: 4px 10px;
    border-radius: 6px;
    font-weight: 500;
    font-size: 13px;
  }

  .magnet-section {
    margin-top: 8px;
  }

  .section-title {
    font-size: 17px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin: 0 0 14px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .magnet-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .magnet-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 13px;
    background: var(--color-bg-secondary);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    border: 1px solid transparent;

    &:hover {
      border-color: var(--color-primary);
    }
  }

  .magnet-icon {
    flex-shrink: 0;
    width: 32px;
    height: 32px;
    background: linear-gradient(135deg, var(--color-primary), var(--color-success));
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-weight: 600;
    font-size: 13px;
  }

  .magnet-text {
    flex: 1;
    font-family: monospace;
    font-size: 12.5px;
    color: var(--color-text-regular);
    word-break: break-all;
    line-height: 1.5;
  }

  .copy-icon {
    flex-shrink: 0;
    width: 30px;
    height: 30px;
    background: var(--color-primary);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 15px;
  }
}

@include mobile {
  .video-detail-page .detail-card {
    .video-title,
    .title {
      font-size: 19px;
    }

    .magnet-item {
      align-items: flex-start;
    }
  }
}
</style>
