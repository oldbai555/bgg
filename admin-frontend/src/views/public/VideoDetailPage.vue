<template>
  <div class="video-detail-page">
    <div class="container">
      <!-- 返回按钮 -->
      <el-button class="back-btn" @click="goBack">
        <el-icon>
          <ArrowLeft />
        </el-icon>
        返回列表
      </el-button>

      <!-- 视频播放器 -->
      <div v-loading="loading" class="video-container">
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
  </div>
</template>

<script setup lang="ts">
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {ArrowLeft} from '@element-plus/icons-vue'
import DPlayer from 'dplayer'
import type {PublicVideoDetailResp} from '@/api/generated/admin'
import {publicVideoDetail} from '@/api/generated/admin'

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
let player: import('dplayer').default | null = null
let isInitializing = false // 防止重复初始化的标志

// 初始化播放器
const initPlayer = async () => {
  if (!dplayerRef.value) {
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
    if (isM3u8) {
      // 使用后端代理接口（支持 Range 请求，实现边下边播）
      const baseURL = import.meta.env.PROD ? '/gateway/api' : '/api'
      finalUrl = `${baseURL}/v1/m3u8/proxy?url=${encodeURIComponent(playUrl)}`
    }

    // 直接播放，失败就失败
    initDPlayer({
      url: finalUrl,
      type: isM3u8 ? 'hls' : 'auto'
    })
  } catch (err: unknown) {
    console.error('播放器初始化失败:', err)
    const message = err instanceof Error ? err.message : '播放器初始化失败'
    ElMessage.error(message)
    loading.value = false
    isInitializing = false
  }
}

// 初始化 DPlayer
const initDPlayer = (options: {url: string; type: 'hls' | 'auto'}) => {
  if (!dplayerRef.value) {
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
    const hlsConfig: any = {}

    // 由于已经使用后端代理，不需要自定义 loader
    // 后端代理会处理 CORS 问题

    player = new DPlayer({
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
      })

      // 错误处理
      videoElement.addEventListener('error', (e) => {
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
    router.push('/public/videos')
    return
  }

  const idNum = Number(id)
  if (!idNum || idNum === 0) {
    ElMessage.error('视频ID格式错误')
    router.push('/public/videos')
    return
  }

  loading.value = true

  try {
    // 响应直接返回数据（无 code/msg 包装），拦截器会直接返回原始数据
    const resp = await publicVideoDetail({id: idNum})
    video.value = resp
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载失败'
    console.error('加载视频详情失败:', err)
    ElMessage.error(message)
    router.push('/public/videos')
  } finally {
    loading.value = false
    // watch 会自动监听 video.value.playUrl 和 loading.value 的变化并初始化播放器
  }
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text)
      ElMessage.success('已复制到剪贴板 ✓')
    } else {
      // Fallback方法
      fallbackCopy(text)
    }
  } catch (_err) {
    console.log('copyToClipboard', _err)
    fallbackCopy(text)
  }
}

// Fallback复制方法
const fallbackCopy = (text: string) => {
  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.select()

  try {
    document.execCommand('copy')
    ElMessage.success('已复制到剪贴板 ✓')
  } catch (_err) {
    console.log('fallbackCopy', _err)
    ElMessage.error('复制失败，请手动复制')
  }

  document.body.removeChild(textarea)
}

// 返回列表
const goBack = () => {
  // 优先走浏览器历史，保留列表分页与滚动状态
  // 如果是从列表页进入的，router.back() 会恢复列表页的状态（包括滚动位置）
  if (window.history.length > 1) {
    router.back()
  } else {
    // 如果没有历史记录，尝试从 sessionStorage 恢复状态
    try {
      const raw = sessionStorage.getItem('public_video_list_state')
      if (raw) {
        const parsed = JSON.parse(raw) as {
          page?: number
          size?: number
          content?: string
        }
        router.push({
          path: '/public/videos',
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
    // 兜底：直接跳转到列表页
    router.push('/public/videos')
  }
}

// 监听 video.playUrl 和 loading 状态，确保在合适的时机初始化播放器
watch(
    () => video.value.playUrl,
    (playUrl) => {
      // 只有当 playUrl 存在且有效时才初始化播放器
      if (playUrl && playUrl.trim() !== '' && dplayerRef.value && !loading.value) {
        // 延迟初始化，确保 DOM 已完全更新
        setTimeout(() => {
          initPlayer()
        }, 200)
      }
    },
    {immediate: false}
)

onMounted(() => {
  loadData()
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
.video-detail-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .back-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 14px;
    margin-bottom: 20px;
    transition: background 0.3s;

    &:hover {
      background: rgba(255, 255, 255, 0.3);
    }
  }

  .video-container {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;
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
    padding: 12px;
  }
}
</style>
