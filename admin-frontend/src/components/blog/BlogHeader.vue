<template>
  <header class="blog-header">
    <div class="header-container">
      <div class="header-left">
        <div class="blog-logo" @click="goHome">
          <span class="logo-text">个人博客</span>
        </div>
        <nav class="header-nav">
          <span class="nav-item active">主页</span>
        </nav>
      </div>
      <div class="header-right">
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索"
            clearable
            @keydown.enter="handleSearch"
            @clear="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        <div v-if="socialInfoList.length > 0" class="social-icons">
          <a
            v-for="info in socialInfoList"
            :key="info.id"
            :href="normalizeUrl(info.url)"
            target="_blank"
            rel="noopener noreferrer"
            class="social-icon"
            :title="info.name"
          >
            <span>{{ getSocialIconText(info.name) }}</span>
          </a>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import {ref, onMounted, watch} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {Search} from '@element-plus/icons-vue'
import {contentApi} from '@/api/content'
import type {PublicBlogSocialInfoListResp} from '@/api/generated/admin'

const router = useRouter()
const route = useRoute()

const searchKeyword = ref('')
const socialInfoList = ref<PublicBlogSocialInfoListResp['list']>([])

// 从路由参数同步搜索关键词
watch(
  () => route.query.keyword,
  (keyword) => {
    if (typeof keyword === 'string') {
      searchKeyword.value = keyword
    } else if (!keyword) {
      searchKeyword.value = ''
    }
  },
  {immediate: true}
)

const normalizeUrl = (url: string): string => {
  if (!url) {
    return ''
  }
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  return `https://${url}`
}

const getSocialIconText = (name: string): string => {
  // 返回社交平台名称的首字母或缩写
  const nameLower = name.toLowerCase()
  if (nameLower.includes('github')) {
    return 'G'
  }
  if (nameLower.includes('gitee')) {
    return '码'
  }
  if (nameLower.includes('bilibili') || nameLower.includes('b站')) {
    return 'B'
  }
  if (nameLower.includes('wechat') || nameLower.includes('微信')) {
    return '微'
  }
  if (nameLower.includes('email') || nameLower.includes('邮件')) {
    return '邮'
  }
  return name.charAt(0).toUpperCase()
}

const goHome = () => {
  router.push('/blog')
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/blog',
      query: {keyword: searchKeyword.value.trim()}
    })
  } else {
    router.push('/blog')
  }
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
  // 初始化时从路由参数读取搜索关键词
  if (route.query.keyword && typeof route.query.keyword === 'string') {
    searchKeyword.value = route.query.keyword
  }
})
</script>

<style scoped lang="scss">
.blog-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border-bottom: 1px solid #f0f0f0;

  .header-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 32px;

    .blog-logo {
      cursor: pointer;
      .logo-text {
        font-size: 20px;
        font-weight: 600;
        color: #333;
      }
    }

    .header-nav {
      display: flex;
      gap: 24px;

      .nav-item {
        font-size: 15px;
        color: #666;
        cursor: pointer;
        padding: 4px 0;
        position: relative;
        transition: color 0.2s;

        &:hover {
          color: #409eff;
        }

        &.active {
          color: #409eff;
          font-weight: 500;

          &::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            height: 2px;
            background: #409eff;
          }
        }
      }
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;

    .search-box {
      width: 200px;

      :deep(.el-input__wrapper) {
        border-radius: 20px;
      }
    }

    .social-icons {
      display: flex;
      gap: 12px;

      .social-icon {
        width: 32px;
        height: 32px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        background: #f5f5f5;
        color: #666;
        text-decoration: none;
        transition: all 0.2s;
        font-size: 14px;

        &:hover {
          background: #409eff;
          color: #fff;
        }
      }
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .blog-header {
    .header-container {
      padding: 0 12px;
      height: 50px;
    }

    .header-left {
      gap: 16px;

      .blog-logo .logo-text {
        font-size: 18px;
      }

      .header-nav {
        display: none; // 移动端隐藏导航
      }
    }

    .header-right {
      gap: 8px;

      .search-box {
        width: 150px;
      }

      .social-icons {
        gap: 8px;

        .social-icon {
          width: 28px;
          height: 28px;
          font-size: 12px;
        }
      }
    }
  }
}
</style>
