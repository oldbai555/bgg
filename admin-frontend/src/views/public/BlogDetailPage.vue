<template>
  <div v-loading="loading" class="blog-detail-page public-detail-page">
    <MetricReporter module="blog_article_detail" :biz-id="Number(route.params.id) || 0" />
    <div v-if="detail" class="container">
      <div class="back-link" @click="goBack">← 返回列表</div>

      <h1 class="title">{{ detail.title }}</h1>

      <div class="meta">
        <span class="author">{{ detail.authorName || '匿名' }}</span>
        <span class="dot">·</span>
        <span class="time">{{ formatTime(detail.publishTime || detail.createdAt) }}</span>
        <span v-if="detail.tags?.length" class="dot">·</span>
        <span v-if="detail.tags?.length" class="tags">
          <el-tag
            v-for="tag in detail.tags"
            :key="tag.id"
            size="small"
            effect="plain"
          >
            {{ tag.name }}
          </el-tag>
        </span>
      </div>

      <div v-if="detail.cover" class="cover">
        <img :src="detail.cover" :alt="detail.title" @error="handleImageError" />
      </div>

      <div class="content">
        <!-- 使用 MdPreview 渲染 Markdown，确保正文图片/链接可见 -->
        <MdPreview :editor-id="'public-blog-detail'" :model-value="detail.content || ''" :preview-theme="'github'" />
      </div>
    </div>

    <div v-else-if="!loading" class="empty">文章不存在或已下架</div>

    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {MdPreview} from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import {blogApi} from '@/api/blog'
import type {PublicBlogArticleDetailReq, PublicBlogArticleDetailResp} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<PublicBlogArticleDetailResp | null>(null)

const formatTime = (ts: number): string => {
  if (!ts) {
return '-'
}
  const d = new Date(ts * 1000)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}

const loadDetail = async () => {
  const id = Number(route.params.id)
  if (!id) {
    ElMessage.error('文章ID不正确')
    return
  }

  loading.value = true
  try {
    const req: PublicBlogArticleDetailReq = {id}
    const resp = await blogApi.publicDetail(req)
    detail.value = resp
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '加载文章失败'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  // 优先走浏览器历史，保留列表页的查询条件与滚动位置
  if (window.history.length > 1) {
    router.back()
    return
  }

  // 没有历史记录时，尝试从 sessionStorage 中恢复列表状态
  try {
    const raw = sessionStorage.getItem('public_blog_list_state')
    if (raw) {
      const parsed = JSON.parse(raw) as {
        page?: number;
        size?: number;
        keyword?: string;
      }
      router.push({
        path: '/public/blog',
        query: {
          ...(parsed.page && {page: String(parsed.page)}),
          ...(parsed.size && {size: String(parsed.size)}),
          ...(parsed.keyword && {keyword: parsed.keyword})
        }
      })
      return
    }
  } catch {
    // 忽略解析错误
  }

  // 兜底：直接返回列表页
  router.push('/public/blog')
}

onMounted(loadDetail)
</script>

<style scoped lang="scss">
@import '@/styles/public-detail.scss';

// 博客详情页特定样式（如需覆盖通用样式，可在此添加）
.blog-detail-page {
  // 保留原有类名以兼容可能的其他引用
}
</style>

