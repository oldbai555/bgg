<template>
  <div class="message-item" :class="{ 'message-self': props.isSelf }">
    <div class="message-avatar">
      <el-avatar :size="36">
        {{ props.message.fromUserName?.charAt(0).toUpperCase() || 'U' }}
      </el-avatar>
    </div>
    <div class="message-content">
      <div class="message-header">
        <span class="message-username">{{ props.message.fromUserName }}</span>
        <span class="message-time">{{ formatTime(props.message.createdAt) }}</span>
      </div>
      <!-- 消息内容：根据消息类型显示 -->
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-if="props.message.messageType === 1" class="message-text" v-html="formatMessageContent(props.message.content)"></div>
      <div v-else-if="props.message.messageType === 2" class="message-image">
        <el-image
          :src="props.message.content"
          fit="cover"
          style="max-width: 300px; max-height: 300px; border-radius: 4px;"
          :preview-src-list="[props.message.content]"
          preview-teleported
        >
          <template #error>
            <div class="image-error">图片加载失败</div>
          </template>
        </el-image>
      </div>
      <div v-else class="message-text">{{ props.message.content }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {ChatMessageItem} from '@/api/generated/admin'

const props = defineProps<{
  message: ChatMessageItem
  isSelf: boolean
}>()

// 格式化消息内容（用于显示 Emoji，转义 HTML 特殊字符）
const formatMessageContent = (content: string) => {
  return content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
    .replace(/\n/g, '<br>')
}

// 格式化时间（接受秒级时间戳）
const formatTime = (timestamp: number) => {
  if (!timestamp) {
    return ''
  }
  const date = new Date(timestamp * 1000) // 秒级时间戳转换为毫秒
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)

  if (minutes < 1) {
    return '刚刚'
  } else if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', {hour: '2-digit', minute: '2-digit'})
  } else {
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}
</script>

<style scoped lang="scss">
.message-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 20px;

  &.message-self {
    flex-direction: row-reverse;

    .message-content {
      align-items: flex-end;

      .message-header {
        flex-direction: row-reverse;
      }

      .message-text {
        background-color: var(--el-color-primary);
        color: white;
      }
    }
  }
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-width: 60%;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.message-username {
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.message-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.message-text {
  padding: 10px 14px;
  border-radius: 8px;
  background-color: var(--el-bg-color);
  font-size: 14px;
  line-height: 1.5;
  word-wrap: break-word;
  white-space: pre-wrap;
}

.message-image {
  max-width: 300px;
  margin-top: 4px;

  .image-error {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 200px;
    height: 200px;
    background-color: var(--el-fill-color-light);
    color: var(--el-text-color-placeholder);
    border-radius: 4px;
  }
}
</style>
