<template>
  <div class="blog-social-links">
    <div v-if="socialInfoList.length > 0" class="social-section">
      <div class="section-title">社交信息</div>
      <div class="social-list">
        <a
          v-for="info in socialInfoList"
          :key="info.id"
          :href="normalizeUrl(info.url)"
          target="_blank"
          rel="noopener noreferrer"
          class="social-item"
        >
          <span class="social-name">{{ info.name }}</span>
        </a>
      </div>
    </div>
    <div v-if="friendLinkList.length > 0" class="friend-link-section">
      <div class="section-title">友情链接</div>
      <div class="friend-link-list">
        <a
          v-for="link in friendLinkList"
          :key="link.id"
          :href="normalizeUrl(link.url)"
          target="_blank"
          rel="noopener noreferrer"
          class="friend-link-item"
        >
          {{ link.name }}
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted} from 'vue'
import {blogApi} from '@/api/blog'
import type {PublicBlogSocialInfoListResp, PublicBlogFriendLinkListResp} from '@/api/generated/admin'

const socialInfoList = ref<PublicBlogSocialInfoListResp['list']>([])
const friendLinkList = ref<PublicBlogFriendLinkListResp['list']>([])

const normalizeUrl = (url: string): string => {
  if (!url) {
return ''
}
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  return `https://${url}`
}

const loadSocialInfo = async () => {
  try {
    const resp = await blogApi.publicSocialInfoList()
    socialInfoList.value = resp.list || []
  } catch (err) {
    console.error('加载社交信息失败:', err)
  }
}

const loadFriendLinks = async () => {
  try {
    const resp = await blogApi.publicFriendLinkList()
    friendLinkList.value = resp.list || []
  } catch (err) {
    console.error('加载友情链接失败:', err)
  }
}

onMounted(() => {
  loadSocialInfo()
  loadFriendLinks()
})
</script>

<style scoped lang="scss">
.blog-social-links {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .social-section,
  .friend-link-section {
    background: #fff;
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);

    .section-title {
      font-size: 16px;
      font-weight: 600;
      color: #333;
      margin-bottom: 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid #f0f0f0;
    }
  }

  .social-list,
  .friend-link-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .social-item,
  .friend-link-item {
    font-size: 14px;
    color: #666;
    text-decoration: none;
    padding: 6px 0;
    transition: color 0.2s;

    &:hover {
      color: #409eff;
    }
  }
}
</style>
