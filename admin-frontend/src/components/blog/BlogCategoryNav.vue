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
  width: 100%; // 占满父容器 grid 单元格
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-light);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
  box-sizing: border-box;
  position: sticky;
  top: 88px; // PublicHeader 高度(64px) + .page-shell 顶部内边距(24px)

  .nav-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--color-border-light);
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
        color: var(--color-text-regular);
      }

      &:hover {
        background: var(--color-bg-secondary);

        .nav-label {
          color: var(--color-primary);
        }
      }

      &.active {
        background: color-mix(in srgb, var(--color-primary) 12%, transparent);

        .nav-label {
          color: var(--color-primary);
          font-weight: 500;
        }
      }
    }
  }
}

// 移动端适配：横向导航条，随内容自然滚动，不再 sticky
@include mobile {
  .blog-category-nav {
    position: static;
    padding: 8px;
    margin-bottom: 14px;

    .nav-title {
      display: none;
    }

    .nav-list {
      flex-direction: row;
      overflow-x: auto;
      gap: 6px;

      .nav-item {
        white-space: nowrap;
      }
    }
  }
}
</style>
