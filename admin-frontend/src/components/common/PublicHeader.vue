<template>
  <header class="public-header">
    <div class="public-header__inner">
      <div class="public-header__brand" @click="goHome">
        <span class="public-header__logo">BGG</span>
      </div>
      <nav class="public-header__nav">
        <router-link to="/front/blog" class="public-header__nav-item" :class="{'is-active': activeTab === 'blog'}">
          博客
        </router-link>
        <router-link to="/front/videos" class="public-header__nav-item" :class="{'is-active': activeTab === 'video'}">
          视频
        </router-link>
      </nav>
      <div v-if="socialInfoList.length > 0" class="public-header__social">
        <a
          v-for="info in socialInfoList"
          :key="info.id"
          :href="normalizeUrl(info.url)"
          target="_blank"
          rel="noopener noreferrer"
          class="public-header__social-icon"
          :title="info.name"
        >
          {{ getSocialIconText(info.name) }}
        </a>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import {ref, computed, onMounted} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {contentApi} from '@/api/content'
import type {PublicBlogSocialInfoListResp} from '@/api/generated/admin'

const router = useRouter()
const route = useRoute()

const socialInfoList = ref<PublicBlogSocialInfoListResp['list']>([])

const activeTab = computed<'blog' | 'video' | ''>(() => {
  if (route.path.startsWith('/front/blog')) return 'blog'
  if (route.path.startsWith('/front/videos')) return 'video'
  return ''
})

const normalizeUrl = (url: string): string => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  return `https://${url}`
}

const getSocialIconText = (name: string): string => {
  const nameLower = name.toLowerCase()
  if (nameLower.includes('github')) return 'G'
  if (nameLower.includes('gitee')) return '码'
  if (nameLower.includes('bilibili') || nameLower.includes('b站')) return 'B'
  if (nameLower.includes('wechat') || nameLower.includes('微信')) return '微'
  if (nameLower.includes('email') || nameLower.includes('邮件')) return '邮'
  return name.charAt(0).toUpperCase()
}

const goHome = () => {
  router.push('/front/blog')
}

const loadSocialInfo = async () => {
  try {
    const resp = await contentApi.publicSocialInfoList()
    socialInfoList.value = resp.list || []
  } catch (err) {
    console.error('加载社交信息失败:', err)
  }
}

onMounted(() => {
  loadSocialInfo()
})
</script>

<style scoped lang="scss">
// 页面级 sticky 侧栏需要预留与本组件高度一致的偏移量，改动高度时同步更新
// public-list.scss / public-detail.scss 里的 $public-header-height 注释
.public-header {
  position: sticky;
  top: 0;
  z-index: 100;
  height: 64px;
  background: var(--color-bg-primary);
  border-bottom: 1px solid var(--color-border);

  &__inner {
    max-width: 1200px;
    height: 100%;
    margin: 0 auto;
    padding: 0 24px;
    display: flex;
    align-items: center;
    gap: 32px;
  }

  &__brand {
    cursor: pointer;
    flex-shrink: 0;
  }

  &__logo {
    font-size: 18px;
    font-weight: 700;
    color: var(--color-primary);
  }

  &__nav {
    display: flex;
    gap: 24px;
    flex: 1;
  }

  &__nav-item {
    font-size: 14px;
    color: var(--color-text-regular);
    text-decoration: none;
    padding: 8px 0;
    border-bottom: 2px solid transparent;
    transition: color 0.2s;

    &:hover {
      color: var(--color-primary);
    }

    &.is-active {
      color: var(--color-primary);
      font-weight: 500;
      border-color: var(--color-primary);
    }
  }

  &__social {
    display: flex;
    gap: 10px;
  }

  &__social-icon {
    width: 30px;
    height: 30px;
    border-radius: 50%;
    background: var(--color-bg-secondary);
    color: var(--color-text-regular);
    display: flex;
    align-items: center;
    justify-content: center;
    text-decoration: none;
    font-size: 12px;
    transition: all 0.2s;

    &:hover {
      background: var(--color-primary);
      color: #fff;
    }
  }
}

@include mobile {
  .public-header {
    height: 52px;

    &__inner {
      padding: 0 14px;
      gap: 16px;
    }

    &__logo {
      font-size: 16px;
    }

    &__nav {
      gap: 16px;
    }

    &__social {
      gap: 6px;
    }

    &__social-icon {
      width: 26px;
      height: 26px;
      font-size: 11px;
    }
  }
}
</style>
