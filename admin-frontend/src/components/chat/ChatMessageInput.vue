<template>
  <div class="message-input">
    <!-- Emoji 选择器 -->
    <div class="emoji-picker-wrapper">
      <el-popover
        placement="top-start"
        :width="300"
        trigger="click"
        popper-class="emoji-picker-popover"
      >
        <template #reference>
          <el-button
            text
            circle
            size="small"
            class="emoji-btn"
          >
            <el-icon :size="20"><ChatDotRound /></el-icon>
          </el-button>
        </template>
        <div class="emoji-picker-container">
          <!-- Emoji 分页显示 -->
          <div
            class="emoji-picker"
            :style="{ gridTemplateColumns: `repeat(${emojiColsPerRow}, 1fr)` }"
          >
            <div
              v-for="emoji in currentPageEmojis"
              :key="emoji"
              class="emoji-item"
              @click="insertEmoji(emoji)"
            >
              {{ emoji }}
            </div>
          </div>
          <!-- 分页控制器 -->
          <div v-if="totalEmojiPages > 1" class="emoji-pagination">
            <el-button
              text
              size="small"
              :disabled="currentEmojiPage === 0"
              @click="currentEmojiPage--"
            >
              上一页
            </el-button>
            <span class="emoji-page-info">{{ currentEmojiPage + 1 }} / {{ totalEmojiPages }}</span>
            <el-button
              text
              size="small"
              :disabled="currentEmojiPage >= totalEmojiPages - 1"
              @click="currentEmojiPage++"
            >
              下一页
            </el-button>
          </div>
        </div>
      </el-popover>
      <!-- 图片上传按钮 -->
      <el-upload
        :action="props.uploadUrl"
        :headers="props.uploadHeaders"
        :on-success="handleImageUploadSuccess"
        :on-error="handleImageUploadError"
        :before-upload="beforeImageUpload"
        :show-file-list="false"
        accept="image/*"
      >
        <el-button
          text
          circle
          size="small"
          class="image-btn"
        >
          <el-icon :size="20"><Picture /></el-icon>
        </el-button>
      </el-upload>
    </div>
    <el-input
      v-model="inputMessage"
      type="textarea"
      :rows="3"
      placeholder="输入消息..."
      @keydown.enter.exact.prevent="handleSend"
      @keydown.enter.shift.exact="inputMessage += '\n'"
    />
    <div class="input-actions">
      <div class="input-info">
        <span v-if="props.selectedChatName">{{ props.selectedChatName }}</span>
        <span v-else>请选择聊天</span>
      </div>
      <el-button
        type="primary"
        :disabled="(!inputMessage.trim() && !pendingImageUrl) || !props.wsConnected"
        @click="handleSend"
      >
        发送
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, computed, nextTick} from 'vue'
import {ElMessage} from 'element-plus'
import {ChatDotRound, Picture} from '@element-plus/icons-vue'
import type {FileUploadResp} from '@/api/generated/admin'
import {buildFileUrlFromResponse} from '@/utils/file'

const props = withDefaults(defineProps<{
  wsConnected: boolean
  hasSelectedChat: boolean
  selectedChatName: string
  uploadUrl: string
  uploadHeaders: Record<string, string>
  emojiColsPerRow?: number
  emojiRows?: number
  onSendText: (text: string) => Promise<void>
  onSendImage: (imageUrl: string) => Promise<void>
}>(), {
  emojiColsPerRow: 8,
  emojiRows: 3
})

const inputMessage = ref('')
const pendingImageUrl = ref('')

// Emoji 列表
const emojiList = [
  '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂', '🙂', '🙃',
  '😉', '😊', '😇', '🥰', '😍', '🤩', '😘', '😗', '😚', '😙',
  '😋', '😛', '😜', '🤪', '😝', '🤑', '🤗', '🤭', '🤫', '🤔',
  '🤐', '🤨', '😐', '😑', '😶', '😏', '😒', '🙄', '😬', '🤥',
  '😌', '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢', '🤮',
  '👍', '👎', '👌', '✌️', '🤞', '🤟', '🤘', '👏', '🙌', '👐',
  '❤️', '💛', '💚', '💙', '💜', '🖤', '🤍', '🤎', '💔', '❣️'
]

const currentEmojiPage = ref(0) // 当前页码

// 计算每页显示的emoji数量（emojiColsPerRow/emojiRows 由父组件从 chat_config 字典加载后传入）
const emojisPerPage = computed(() => props.emojiColsPerRow * props.emojiRows)

// 计算总页数
const totalEmojiPages = computed(() => Math.ceil(emojiList.length / emojisPerPage.value))

// 当前页显示的emoji列表
const currentPageEmojis = computed(() => {
  const start = currentEmojiPage.value * emojisPerPage.value
  const end = start + emojisPerPage.value
  return emojiList.slice(start, end)
})

// 插入 Emoji
const insertEmoji = (emoji: string) => {
  const textarea = document.querySelector('.message-input textarea') as HTMLTextAreaElement
  if (textarea) {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    inputMessage.value = inputMessage.value.substring(0, start) + emoji + inputMessage.value.substring(end)
    // 设置光标位置
    nextTick(() => {
      textarea.focus()
      textarea.setSelectionRange(start + emoji.length, start + emoji.length)
    })
  } else {
    inputMessage.value += emoji
  }
}

// 图片上传前验证
const beforeImageUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB！')
    return false
  }
  return true
}

// 图片上传成功
const handleImageUploadSuccess = async (response: FileUploadResp) => {
  try {
    const fullUrl = buildFileUrlFromResponse(response)
    pendingImageUrl.value = fullUrl
    ElMessage.success('图片上传成功，点击发送按钮发送')
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '未知错误'
    ElMessage.error('图片上传失败：' + message)
  }
}

// 图片上传失败
const handleImageUploadError = (error: Error | unknown) => {
  const message = error instanceof Error ? error.message : '未知错误'
  ElMessage.error('图片上传失败：' + message)
}

// 发送消息（如果文本和图片都有，各自作为一条消息依次发送，与原实现保持一致）
const handleSend = async () => {
  const hasText = inputMessage.value.trim()
  const hasImage = pendingImageUrl.value

  if (!hasText && !hasImage) {
    return
  }

  if (!props.wsConnected) {
    ElMessage.warning('WebSocket 未连接，请先连接')
    return
  }

  if (!props.hasSelectedChat) {
    ElMessage.warning('请先选择一个聊天')
    return
  }

  try {
    if (hasImage) {
      const imageUrl = pendingImageUrl.value
      await props.onSendImage(imageUrl)
      pendingImageUrl.value = ''
    }

    if (hasText) {
      const messageContent = inputMessage.value.trim()
      inputMessage.value = ''
      await props.onSendText(messageContent)
    }
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '发送消息失败'
    ElMessage.error(message)
  }
}
</script>

<style scoped lang="scss">
.message-input {
  padding: 15px;
  border-top: 1px solid var(--el-border-color-light);
  background-color: var(--el-bg-color);
}

.emoji-picker-wrapper {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;

  .emoji-btn,
  .image-btn {
    color: var(--el-text-color-regular);

    &:hover {
      color: var(--el-color-primary);
    }
  }
}

.emoji-picker-container {
  display: flex;
  flex-direction: column;
}

.emoji-picker {
  display: grid;
  gap: 4px;
  padding: 8px;
  min-width: 200px;

  .emoji-item {
    font-size: 20px;
    padding: 4px;
    cursor: pointer;
    text-align: center;
    border-radius: 4px;
    transition: background-color 0.2s;

    &:hover {
      background-color: var(--el-fill-color-light);
    }
  }
}

.emoji-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-top: 1px solid var(--el-border-color-lighter);
  font-size: 12px;
}

.emoji-page-info {
  color: var(--el-text-color-secondary);
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;

  .input-info {
    display: flex;
    gap: 15px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
}
</style>
