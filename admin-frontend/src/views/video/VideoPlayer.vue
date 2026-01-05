<template>
  <div class="video-player-page">
    <el-card>
      <!-- 顶部操作栏 -->
      <div class="header">
        <el-input
          v-model="videoUrl"
          placeholder="请输入视频链接（支持m3u8、mp4等格式）..."
          clearable
          @keydown.enter="handlePlay"
          class="url-input"
        >
          <template #append>
            <el-button type="primary" @click="handlePlay" :loading="loading">
              <el-icon><VideoPlay /></el-icon>
              播放
            </el-button>
          </template>
        </el-input>
        <el-button
          v-if="showPlayer"
          @click="handleReset"
          type="info"
          plain
        >
          <el-icon><RefreshLeft /></el-icon>
          重新输入
        </el-button>
      </div>

      <!-- 播放器区域 -->
      <div class="player-container" v-show="showPlayer">
        <video
          ref="videoPlayerRef"
          class="video-js vjs-default-skin vjs-big-play-centered"
          controls
          preload="auto"
          :poster="videoCover"
        ></video>
      </div>

      <!-- 提示信息 -->
      <div v-if="!showPlayer" class="tips">
        <el-empty description="请输入视频链接并点击播放" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted, onBeforeUnmount, watch} from 'vue';
import {useRoute} from 'vue-router';
import {ElMessage} from 'element-plus';
import {VideoPlay, RefreshLeft} from '@element-plus/icons-vue';
import videojs from 'video.js';
import 'video.js/dist/video-js.css';
import {videoProxy} from '@/api/generated/admin';

const route = useRoute();

const videoUrl = ref('');
const videoCover = ref('');
const showPlayer = ref(false);
const loading = ref(false);
const videoPlayerRef = ref<HTMLVideoElement | null>(null);
let player: videojs.Player | null = null;

// 从路由参数获取视频URL
onMounted(() => {
  const urlParam = route.query.url as string;
  if (urlParam) {
    videoUrl.value = decodeURIComponent(urlParam);
    handlePlay();
  }
});

// 监听路由变化
watch(() => route.query.url, (newUrl) => {
  if (newUrl) {
    videoUrl.value = decodeURIComponent(newUrl as string);
    handlePlay();
  }
});

// 播放视频
const handlePlay = async () => {
  if (!videoUrl.value || videoUrl.value.trim() === '') {
    ElMessage.warning('请输入视频链接');
    return;
  }

  const url = videoUrl.value.trim();

  // 验证URL格式
  try {
    new URL(url);
  } catch (e) {
    ElMessage.error('视频链接格式不正确');
    return;
  }

  loading.value = true;
  showPlayer.value = true;

  try {
    // 销毁旧的播放器
    if (player) {
      player.dispose();
      player = null;
    }

    // 等待DOM更新
    await new Promise(resolve => setTimeout(resolve, 100));

    if (!videoPlayerRef.value) {
      ElMessage.error('播放器初始化失败');
      return;
    }

    // 初始化播放器
    player = videojs(videoPlayerRef.value, {
      autoplay: false,
      controls: true,
      preload: 'auto',
      fluid: true,
      responsive: true,
      playbackRates: [0.5, 1, 1.25, 1.5, 2],
      html5: {
        vhs: {
          overrideNative: true
        },
        nativeVideoTracks: false,
        nativeAudioTracks: false,
        nativeTextTracks: false
      }
    });

    // 尝试直接播放
    let playSource = url;
    let useProxy = false;

    // 如果是m3u8格式，先尝试直接播放，失败则使用代理
    if (url.includes('.m3u8')) {
      try {
        player.src({
          src: url,
          type: 'application/x-mpegURL'
        });
        
        // 监听错误事件
        player.one('error', async () => {
          const error = player?.error();
          if (error && error.code === 4) {
            // 媒体资源无法加载，尝试使用代理
            ElMessage.info('直接播放失败，尝试使用代理播放...');
            await tryProxyPlay(url);
          } else {
            ElMessage.error(`播放失败: ${error?.message || '未知错误'}`);
            loading.value = false;
          }
        });

        // 监听加载成功
        player.one('loadedmetadata', () => {
          ElMessage.success('视频加载成功');
          loading.value = false;
        });

        // 设置超时
        setTimeout(() => {
          if (loading.value && player?.error()) {
            // 如果还在加载且出现错误，尝试代理
            tryProxyPlay(url);
          }
        }, 5000);

      } catch (e) {
        // 直接播放失败，尝试代理
        await tryProxyPlay(url);
      }
    } else {
      // 非m3u8格式，直接播放
      player.src(url);
      player.one('loadedmetadata', () => {
        ElMessage.success('视频加载成功');
        loading.value = false;
      });
      player.one('error', () => {
        const error = player?.error();
        ElMessage.error(`播放失败: ${error?.message || '未知错误'}`);
        loading.value = false;
      });
    }

  } catch (error: any) {
    ElMessage.error(`播放失败: ${error.message || '未知错误'}`);
    loading.value = false;
    showPlayer.value = false;
  }
};

// 尝试使用代理播放
const tryProxyPlay = async (url: string) => {
  try {
    // 构建代理URL
    const proxyUrl = `/api/v1/videos/proxy?url=${encodeURIComponent(url)}`;
    
    // 设置播放源为代理URL
    if (player) {
      if (url.includes('.m3u8')) {
        player.src({
          src: proxyUrl,
          type: 'application/x-mpegURL'
        });
      } else {
        player.src(proxyUrl);
      }

      player.one('loadedmetadata', () => {
        ElMessage.success('通过代理播放成功');
        loading.value = false;
      });

      player.one('error', () => {
        const error = player?.error();
        ElMessage.error(`代理播放也失败: ${error?.message || '视频链接无法播放'}`);
        loading.value = false;
      });
    }
  } catch (error: any) {
    ElMessage.error(`代理播放失败: ${error.message || '未知错误'}`);
    loading.value = false;
  }
};

// 重置播放器
const handleReset = () => {
  if (player) {
    player.pause();
    player.dispose();
    player = null;
  }
  showPlayer.value = false;
  videoUrl.value = '';
  videoCover.value = '';
};

// 组件销毁时清理播放器
onBeforeUnmount(() => {
  if (player) {
    player.dispose();
    player = null;
  }
});
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

    :deep(.video-js) {
      width: 100%;
      height: 600px;
      max-width: 100%;

      .vjs-big-play-button {
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
      }
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

