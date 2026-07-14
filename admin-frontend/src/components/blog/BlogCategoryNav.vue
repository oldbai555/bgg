<template>
  <aside class="blog-category-nav">
    <div class="nav-title">分类</div>
    <div class="nav-list">
      <div
        class="nav-item"
        :class="{active: selectedTagId === 0}"
        @click="handleSelect(0)"
      >
        <span class="nav-label">全部</span>
      </div>
      <div
        v-for="tag in tagList"
        :key="tag.id"
        class="nav-item"
        :class="{active: selectedTagId === tag.id}"
        @click="handleSelect(tag.id)"
      >
        <span class="nav-label">{{ tag.name }}</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import {ref, onMounted} from 'vue'
import {contentApi} from '@/api/content'
import type {PublicBlogTagListResp} from '@/api/generated/admin'

defineProps<{
  selectedTagId?: number
}>()

const emit = defineEmits<{
  (e: 'select', tagId: number): void
}>()

const tagList = ref<PublicBlogTagListResp['list']>([])

const handleSelect = (tagId: number) => {
  emit('select', tagId)
}

const loadTags = async () => {
  try {
    const resp = await contentApi.publicTagList()
    tagList.value = resp.list || []
  } catch (err) {
    console.error('加载标签列表失败:', err)
  }
}

onMounted(() => {
  loadTags()
})
</script>

<style scoped lang="scss">
.blog-category-nav {
  width: 100%; // 占满父容器（blog-sidebar-left）
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  position: relative; // 确保定位上下文
  z-index: 10; // 提高层级，确保在中间内容之上
  box-sizing: border-box; // 确保 padding 包含在宽度内

  .nav-title {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f0f0f0;
  }

  .nav-list {
    display: flex;
    flex-direction: column;
    gap: 4px;

    .nav-item {
      padding: 8px 12px;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.2s;

      .nav-label {
        font-size: 14px;
        color: #666;
      }

      &:hover {
        background: #f5f5f5;

        .nav-label {
          color: #409eff;
        }
      }

      &.active {
        background: #e6f4ff;

        .nav-label {
          color: #409eff;
          font-weight: 500;
        }
      }
    }
  }
}

// 移动端适配
@include mobile {
  .blog-category-nav {
    width: 100%;
    margin-bottom: 16px;
  }
}
</style>
