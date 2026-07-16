<!-- 即时通讯 UI（会话列表+消息气泡+输入区），不是管理后台 CRUD 列表，不适用 D2Table -->
<template>
  <div class="chat-container">
    <el-card class="chat-card">
      <template #header>
        <div class="chat-header">
          <span class="chat-title">在线聊天</span>
          <div class="chat-status">
            <el-tag :type="chat.wsConnected.value ? 'success' : 'danger'" size="small">
              {{ chat.wsConnected.value ? '已连接' : '未连接' }}
            </el-tag>
            <el-button
              v-if="!chat.wsConnected.value"
              type="primary"
              size="small"
              @click="wsStore.connect()"
            >
              重新连接
            </el-button>
          </div>
        </div>
      </template>

      <div class="chat-content">
        <!-- 左侧：聊天列表 -->
        <div class="chat-sidebar">
          <div class="sidebar-header">
            <h3>聊天列表 ({{ chat.chats.value.length }})</h3>
          </div>
          <div class="user-list">
            <ChatListItem
              v-for="item in chat.chats.value"
              :key="item.chatId"
              :chat="item"
              :active="chat.selectedChatId.value === item.chatId"
              @click="chat.selectChat(item)"
            />
            <div
              v-if="chat.chats.value.length === 0"
              class="empty-users"
            >
              <el-empty description="暂无聊天" :image-size="80" />
            </div>
          </div>
        </div>

        <!-- 右侧：聊天区域 -->
        <div class="chat-main">
          <!-- 消息列表 -->
          <div ref="messageListRef" class="message-list">
            <ChatMessageBubble
              v-for="message in chat.messages.value"
              :key="message.id"
              :message="message"
              :is-self="Number(message.fromUserId) === Number(chat.currentUserId.value)"
            />
            <div v-if="chat.messages.value.length === 0" class="empty-message">
              <el-empty description="暂无消息，开始聊天吧~" />
            </div>
          </div>

          <!-- 输入区域 -->
          <ChatMessageInput
            :ws-connected="chat.wsConnected.value"
            :has-selected-chat="!!chat.selectedChatId.value"
            :selected-chat-name="chat.selectedChat.value?.name || ''"
            :upload-url="uploadUrl"
            :upload-headers="uploadHeaders"
            :emoji-cols-per-row="chat.emojiColsPerRow.value"
            :emoji-rows="chat.emojiRows.value"
            :on-send-text="chat.sendTextMessage"
            :on-send-image="chat.sendImageMessage"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {ref, computed, onMounted, onUnmounted, nextTick} from 'vue'
import {useUserStore} from '@/stores/user'
import {useWebSocketStore} from '@/stores/websocket'
import {useAppConfig} from '@/composables/useAppConfig'
import {useChatList} from '@/composables/useChatList'
import ChatListItem from '@/components/chat/ChatListItem.vue'
import ChatMessageBubble from '@/components/chat/ChatMessageBubble.vue'
import ChatMessageInput from '@/components/chat/ChatMessageInput.vue'

const userStore = useUserStore()
const wsStore = useWebSocketStore()
const messageListRef = ref<HTMLElement>()

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}

const chat = useChatList({onMessagesChanged: scrollToBottom})

// 应用配置
const {storageBaseURL, initConfig} = useAppConfig()

// 文件上传配置
const uploadUrl = computed(() => {
  // 开发环境：始终使用 vite 代理路径（避免 CORS）
  if (import.meta.env.DEV) {
    return '/api/v1/files/upload'
  }
  // 生产环境：使用字典配置的 baseURL
  if (storageBaseURL.value) {
    return `${storageBaseURL.value}/api/v1/files/upload`
  }
  // 生产环境默认使用网关路径
  return ''
})
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

onMounted(async () => {
  // 初始化应用配置（从字典获取）
  await initConfig()

  // 确保用户信息已加载（刷新页面后可能需要重新获取）
  if (userStore.token && !userStore.profile) {
    await userStore.fetchProfile(true)
  }

  // 先加载聊天配置，再加载聊天列表
  await chat.loadChatConfig()
  await chat.loadChats()
  // WebSocket 连接由全局 store 管理，这里不需要手动连接
  // 但确保连接已建立
  if (!chat.wsConnected.value && userStore.token) {
    wsStore.connect()
  }
})

onUnmounted(() => {
  // 不断开 WebSocket，因为可能在其他页面还需要使用
})
</script>

<style scoped lang="scss">
.chat-container {
  height: calc(100vh - 120px);
  padding: 20px;

  .chat-card {
    height: 100%;

    :deep(.el-card__body) {
      height: calc(100% - 60px);
      padding: 0;
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }
  }
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .chat-title {
    font-size: 18px;
    font-weight: 600;
  }

  .chat-status {
    display: flex;
    align-items: center;
    gap: 10px;
  }
}

.chat-content {
  display: flex;
  height: 100%;
  border-top: 1px solid var(--el-border-color-light);
}

.chat-sidebar {
  width: 250px;
  border-right: 1px solid var(--el-border-color-light);
  display: flex;
  flex-direction: column;

  .sidebar-header {
    padding: 15px;
    border-bottom: 1px solid var(--el-border-color-light);

    h3 {
      margin: 0;
      font-size: 14px;
      font-weight: 600;
    }
  }

  .user-list {
    flex: 1;
    overflow-y: auto;
    padding: 10px;

    .empty-users {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 200px;
    }
  }
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: var(--el-bg-color-page);

  .empty-message {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
  }
}
</style>
