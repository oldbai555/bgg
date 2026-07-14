<template>
  <div class="blog-author-card">
    <div class="author-avatar">
      <img v-if="authorInfo.avatar" :src="authorInfo.avatar" :alt="authorInfo.nickname" />
      <div v-else class="avatar-placeholder">{{ authorInfo.nickname?.charAt(0) || 'A' }}</div>
    </div>
    <div class="author-name">{{ authorInfo.nickname || '管理员' }}</div>
    <div v-if="authorInfo.signature" class="author-signature">{{ authorInfo.signature }}</div>
    <div v-if="articleStats.totalArticles > 0" class="author-stats">
      <span class="stat-item">{{ articleStats.totalArticles }} 文章</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted} from 'vue'
import {contentApi} from '@/api/content'
import type {PublicBlogAuthorInfoResp, PublicBlogArticleStatsResp} from '@/api/generated/admin'

const authorInfo = ref<PublicBlogAuthorInfoResp>({
  id: 1,
  nickname: '',
  avatar: '',
  signature: ''
})

const articleStats = ref<PublicBlogArticleStatsResp>({
  totalArticles: 0
})

const loadAuthorInfo = async () => {
  try {
    const resp = await contentApi.publicAuthorInfo()
    authorInfo.value = resp
  } catch (err) {
    console.error('加载作者信息失败:', err)
  }
}

const loadArticleStats = async () => {
  try {
    const resp = await contentApi.publicArticleStats()
    articleStats.value = resp
  } catch (err) {
    console.error('加载文章统计失败:', err)
  }
}

onMounted(() => {
  loadAuthorInfo()
  loadArticleStats()
})
</script>

<style scoped lang="scss">
.blog-author-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-light);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
  text-align: center;

  .author-avatar {
    width: 80px;
    height: 80px;
    margin: 0 auto 12px;
    border-radius: 50%;
    overflow: hidden;
    background: var(--color-bg-secondary);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .avatar-placeholder {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 32px;
      font-weight: 600;
      color: #fff;
      background: linear-gradient(135deg, var(--color-primary), var(--color-success));
    }
  }

  .author-name {
    font-size: 18px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 8px;
  }

  .author-signature {
    font-size: 14px;
    color: var(--color-text-secondary);
    margin-bottom: 12px;
    line-height: 1.5;
  }

  .author-stats {
    padding-top: 12px;
    border-top: 1px solid var(--color-border-light);

    .stat-item {
      font-size: 14px;
      color: var(--color-text-regular);
    }
  }
}
</style>
