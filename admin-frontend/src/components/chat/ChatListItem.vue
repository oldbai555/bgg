<template>
  <div
    class="user-item"
    :class="{ active: props.active }"
    @click="emit('click')"
  >
    <el-avatar :size="32" :src="props.chat.avatar || ''">
      {{ props.chat.name?.charAt(0).toUpperCase() || 'C' }}
    </el-avatar>
    <div class="user-info">
      <div class="user-name">
        {{ props.chat.name }}
        <el-tag
          v-if="props.chat.type === 2"
          size="small"
          type="info"
          style="margin-left: 4px"
        >群组</el-tag>
      </div>
      <div v-if="props.chat.type === 1" class="user-desc">{{ formatChatDesc(props.chat) }}</div>
      <div v-else-if="props.chat.description" class="user-desc">{{ props.chat.description }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {ChatItem} from '@/api/generated/admin'

const props = defineProps<{
  chat: ChatItem
  active: boolean
}>()

const emit = defineEmits<{
  click: []
}>()

// 格式化聊天描述：部门-角色-用户昵称（仅私聊）
const formatChatDesc = (chat: ChatItem): string => {
  if (chat.type !== 1) {
    return chat.description || ''
  }
  const parts: string[] = []
  if (chat.departmentName) {
    parts.push(chat.departmentName)
  }
  if (chat.roleNames && chat.roleNames.length > 0) {
    parts.push(chat.roleNames.join('、'))
  }
  if (chat.nickname) {
    parts.push(chat.nickname)
  }
  return parts.join('-') || chat.username || ''
}
</script>

<style scoped lang="scss">
.user-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: var(--el-fill-color-light);
  }

  &.active {
    background-color: var(--el-color-primary-light-9);
  }
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.user-desc {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
