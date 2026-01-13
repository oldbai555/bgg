<template>
  <div class="video-player-page">
    <el-card>
      <!-- 顶部操作栏 -->
      <div class="header">
        <el-input
          v-model="videoUrl"
          placeholder="请输入视频链接（支持m3u8、mp4等格式）..."
          clearable
          class="url-input"
          @keydown.enter="handlePlay"
        >
          <template #append>
            <el-button type="primary" :loading="loading" @click="handlePlay">
              <el-icon><VideoPlay /></el-icon>
              播放
            </el-button>
          </template>
        </el-input>
        <el-button
          v-if="showPlayer"
          type="info"
          plain
          @click="handleReset"
        >
          <el-icon><RefreshLeft /></el-icon>
          重新输入
        </el-button>
      </div>

      <!-- 播放器区域 -->
      <div v-show="showPlayer" class="player-container">
        <div ref="dplayerRef" class="dplayer-container"></div>
      </div>

      <!-- 提示信息 -->
      <div v-if="!showPlayer" class="tips">
        <el-empty description="请输入视频链接并点击播放" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted, onBeforeUnmount, watch, nextTick} from 'vue'
import {useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {VideoPlay, RefreshLeft} from '@element-plus/icons-vue'
import DPlayer from 'dplayer'

const route = useRoute()

const videoUrl = ref('')
const videoCover = ref('')
const showPlayer = ref(false)
const loading = ref(false)
const dplayerRef = ref<HTMLDivElement | null>(null)
let player: DPlayer | null = null

// 从路由参数获取视频URL
onMounted(() => {
  const urlParam = route.query.url as string
  if (urlParam) {
    videoUrl.value = decodeURIComponent(urlParam)
    handlePlay()
  }
})

// 监听路由变化
watch(() => route.query.url, (newUrl) => {
  if (newUrl) {
    videoUrl.value = decodeURIComponent(newUrl as string)
    handlePlay()
  }
})

// 播放视频
const handlePlay = async () => {
  if (!videoUrl.value || videoUrl.value.trim() === '') {
    ElMessage.warning('请输入视频链接')
    return
  }

  const url = videoUrl.value.trim()

  // 验证URL格式
  try {
    new URL(url)
  } catch {
    ElMessage.error('视频链接格式不正确')
    return
  }

  loading.value = true
  showPlayer.value = true

  try {
    // 销毁旧的播放器
    if (player) {
      player.destroy()
      player = null
    }

    // 等待DOM更新
    await nextTick()

    if (!dplayerRef.value) {
      ElMessage.error('播放器容器初始化失败')
      loading.value = false
      return
    }

    // 判断视频类型
    const isM3u8 = url.includes('.m3u8')
    
    // 直接播放，失败就失败
    initPlayer({
      url: url,
      type: isM3u8 ? 'hls' : 'auto'
    })

  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '未知错误'
    ElMessage.error(`播放失败: ${message}`)
    loading.value = false
    showPlayer.value = false
  }
}

// 初始化 DPlayer
const initPlayer = (options: {url: string; type: 'hls' | 'auto'}) => {
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
    player = new DPlayer({
      container: dplayerRef.value,
      video: {
        url: options.url,
        type: options.type,
        pic: videoCover.value || undefined
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
        // HLS 配置，用于处理 CORS
        xhrSetup: (xhr: XMLHttpRequest) => {
          // 设置 CORS 相关请求头
          xhr.withCredentials = false
        }
      }
    })

    // 监听播放器事件（使用原生 video 元素的事件）
    if (player.video) {
      player.video.addEventListener('loadstart', () => {
        loading.value = true
      })

      player.video.addEventListener('canplay', () => {
        loading.value = false
      })

      player.video.addEventListener('error', () => {
        loading.value = false
      })
    }
  } catch (error) {
    console.error('初始化播放器失败:', error)
    ElMessage.error('播放器初始化失败')
    loading.value = false
    throw error
  }
}

// 重置播放器
const handleReset = () => {
  if (player) {
    try {
      player.destroy()
    } catch (e) {
      console.warn('销毁播放器失败:', e)
    }
    player = null
  }
  showPlayer.value = false
  videoUrl.value = ''
  videoCover.value = ''
  loading.value = false
}

// 组件销毁时清理播放器
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
.video-player-page {
  padding: 20px;
  height: 100%;

  .header {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    align-items: center;

    .url-input {
      flex: 1;
    }
  }

  .player-container {
    width: 100%;
    min-height: 400px;

    .dplayer-container {
      width: 100%;
      height: 600px;
      max-width: 100%;
    }
  }

  .tips {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 400px;
  }
}
</style>
