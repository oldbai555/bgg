<template>
  <div class="video-detail-page">
    <div class="container">
      <!-- 返回按钮 -->
      <el-button class="back-btn" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>

      <!-- 视频播放器 -->
      <div v-loading="loading" class="video-container">
        <div class="video-wrapper">
          <video
            :key="`video-${video.id || 'new'}`"
            ref="videoPlayerRef"
            class="video-js vjs-default-skin"
            controls
            preload="auto"
          ></video>
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
import {ref, onMounted, onBeforeUnmount, nextTick, watch} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {ArrowLeft} from '@element-plus/icons-vue'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import {publicVideoDetail} from '@/api/generated/admin'
import type {PublicVideoDetailResp} from '@/api/generated/admin'

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
const videoPlayerRef = ref<HTMLVideoElement | null>(null)
let player: ReturnType<typeof videojs> | null = null

// 等待元素在 DOM 中
const waitForElementInDOM = async (element: HTMLElement, maxAttempts = 20): Promise<boolean> => {
  for (let i = 0; i < maxAttempts; i++) {
    // 检查元素是否在 DOM 中（使用多种方法）
    if (document.body.contains(element) ||
        element.isConnected ||
        (element.offsetParent !== null || element.style.display !== 'none')) {
      return true
    }
    await new Promise(resolve => setTimeout(resolve, 50))
  }
  return false
}

// 初始化播放器
const initPlayer = async () => {
  if (!videoPlayerRef.value || !video.value.playUrl) {
    return
  }

  // 如果还在加载中，等待加载完成
  if (loading.value) {
    return
  }

  try {
    // 销毁旧的播放器
    if (player) {
      player.dispose()
      player = null
    }

    // 等待 Vue 完成 DOM 更新
    await nextTick()

    // 使用 requestAnimationFrame 确保浏览器已完成渲染
    await new Promise(resolve => requestAnimationFrame(resolve))
    await new Promise(resolve => requestAnimationFrame(resolve))

    // 额外等待，确保 loading 状态已更新
    await new Promise(resolve => setTimeout(resolve, 100))

    if (!videoPlayerRef.value) {
      return
    }

    // 等待元素在 DOM 中
    const isInDOM = await waitForElementInDOM(videoPlayerRef.value)

    if (!isInDOM) {
      return
    }

    // 再次确认元素可见（不是被 v-loading 隐藏）
    const isHidden = videoPlayerRef.value.offsetParent === null && videoPlayerRef.value.style.display === 'none'

    if (isHidden) {
      await new Promise(resolve => setTimeout(resolve, 200))
      const stillHidden = videoPlayerRef.value.offsetParent === null && videoPlayerRef.value.style.display === 'none'
      if (stillHidden) {
        return
      }
    }

    // 检查 video 元素是否已经被 video.js 初始化
    // 如果已经初始化，先清理
    try {
      // 先通过 ID 查找可能的旧实例
      const videoElement = videoPlayerRef.value
      if (videoElement && videoElement.id) {
        const existingPlayerById = videojs.getPlayer(videoElement.id)
        if (existingPlayerById) {
          existingPlayerById.dispose()
        }
      }

      // 再通过元素查找
      if (videoElement) {
        try {
          const existingPlayer = videojs.getPlayer(videoElement)
          if (existingPlayer) {
            existingPlayer.dispose()
          }
        } catch (_e) {
          // 如果没有已存在的播放器，忽略错误
        }
      }
    } catch (_e) {
      // 忽略检查错误
    }

    // dispose 后等待 Vue 重新渲染（如果使用了 key，元素会被重新创建）
    await nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    // 重新获取元素引用（因为 dispose 可能已经移除了旧元素）
    if (!videoPlayerRef.value) {
      // 等待 Vue 重新渲染元素
      let retryCount = 0
      const maxRetries = 20
      while (retryCount < maxRetries && !videoPlayerRef.value) {
        await nextTick()
        await new Promise(resolve => setTimeout(resolve, 50))
        retryCount++
      }

      if (!videoPlayerRef.value) {
        return
      }
    }

    const playUrl = video.value.playUrl
    const isM3u8 = playUrl.includes('.m3u8')

    // 如果元素不在 DOM 中，等待它出现
    if (!document.body.contains(videoPlayerRef.value) || !videoPlayerRef.value.isConnected) {
      let retryCount = 0
      const maxRetries = 20
      while (retryCount < maxRetries && (!document.body.contains(videoPlayerRef.value) || !videoPlayerRef.value.isConnected)) {
        await nextTick()
        await new Promise(resolve => setTimeout(resolve, 50))
        retryCount++
      }

      if (!videoPlayerRef.value || (!document.body.contains(videoPlayerRef.value) && !videoPlayerRef.value.isConnected)) {
        return
      }
    }

    // 初始化video.js播放器
    player = videojs(videoPlayerRef.value, {
      autoplay: false, // 不自动播放，让用户手动点击
      controls: true,
      preload: 'auto',
      fluid: false, // 禁用 fluid，使用固定宽高比
      responsive: false, // 禁用 responsive，使用固定宽高比
      playbackRates: [0.5, 1, 1.25, 1.5, 2],
      html5: {
        vhs: {
          overrideNative: true
        },
        nativeVideoTracks: false,
        nativeAudioTracks: false,
        nativeTextTracks: false
      }
    })

    if (!player) {
      return
    }

    // 等待播放器准备就绪
    const currentPlayer = player
    if (!currentPlayer) {
      return
    }

    currentPlayer.ready(() => {
      // 对于 m3u8 格式，先尝试直接播放
      if (isM3u8) {
        // 先尝试直接播放原始 URL
        currentPlayer.src({
          src: playUrl,
          type: 'application/x-mpegURL'
        })

        // 监听错误，如果直接播放失败，切换到代理地址
        const errorHandler = () => {
          const playerInstance = player
          if (!playerInstance) {
return
}

          // 切换到代理地址
          const proxyUrl = `/api/v1/m3u8/proxy?url=${encodeURIComponent(playUrl)}`
          playerInstance.src({
            src: proxyUrl,
            type: 'application/x-mpegURL'
          })

          // 移除错误监听，避免循环
          playerInstance.off('error', errorHandler)
        }

        currentPlayer.one('error', errorHandler)
      } else {
        // 非 m3u8 格式，直接播放
        currentPlayer.src(playUrl)
      }

      // 等待播放器控件完全初始化后再监听按钮（延迟设置）
      setTimeout(() => {
        // 监听大播放按钮点击
        const bigPlayButton = currentPlayer.getChild('bigPlayButton')
        if (bigPlayButton) {
          bigPlayButton.on('click', () => {
            const playerInstance = player
            if (playerInstance) {
              // 检查是否有错误
              const error = playerInstance.error()
              if (error) {
                ElMessage.error(`播放失败: ${error.message || '未知错误'}`)
                return
              }
              // 尝试播放
              const playPromise = playerInstance.play()
              if (playPromise !== undefined) {
                playPromise
                  .catch((err) => {
                    ElMessage.error(`播放失败: ${err.message || '未知错误'}`)
                  })
              }
            }
          })
        }
      }, 500)

      // 监听最终错误（如果代理也失败）
      currentPlayer.on('error', () => {
        const playerInstance = player
        if (!playerInstance) {
return
}
        const error = playerInstance.error()
        if (error) {
          ElMessage.error(`播放失败: ${error.message || '未知错误'}`)
        }
      })
    })
  } catch (err: unknown) {
    console.error('播放器初始化失败:', err)
    const message = err instanceof Error ? err.message : '播放器初始化失败'
    ElMessage.error(message)
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
    ElMessage.error('复制失败，请手动复制')
  }

  document.body.removeChild(textarea)
}

// 返回列表
const goBack = () => {
  router.push('/public/videos')
}

// 监听 video.playUrl 和 loading 状态，确保在合适的时机初始化播放器
watch(
  () => [video.value.playUrl, loading.value],
  ([playUrl, isLoading]) => {
    if (playUrl && !isLoading && videoPlayerRef.value) {
      // 延迟初始化，确保 DOM 已完全更新
      setTimeout(() => {
        initPlayer()
      }, 200)
    }
  },
  {immediate: false}
)

onMounted(() => {
  // 检查是否有旧的 video.js 实例需要清理
  if (videoPlayerRef.value) {
    try {
      const existingPlayer = videojs.getPlayer(videoPlayerRef.value)
      if (existingPlayer) {
        existingPlayer.dispose()
      }
    } catch (_e) {
      // 如果没有已存在的播放器，忽略错误
    }
  }

  loadData()
})

onBeforeUnmount(() => {
  if (player) {
    player.dispose()
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

  :deep(.video-js) {
    position: absolute !important;
    top: 0 !important;
    left: 0 !important;
    width: 100% !important;
    height: 100% !important;
    padding: 0 !important;
  }

  :deep(.video-js .vjs-tech) {
    position: absolute !important;
    top: 0 !important;
    left: 0 !important;
    width: 100% !important;
    height: 100% !important;
    object-fit: contain;
  }

  :deep(.video-js .vjs-poster) {
    position: absolute !important;
    top: 0 !important;
    left: 0 !important;
    width: 100% !important;
    height: 100% !important;
  }

  :deep(.video-js .vjs-big-play-button) {
    position: absolute !important;
    top: 50% !important;
    left: 50% !important;
    transform: translate(-50%, -50%) !important;
    margin: 0 !important;
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

