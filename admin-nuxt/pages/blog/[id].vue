<template>
  <div :class="['blog-detail-page', { 'is-loading': loading }]">
    <ClientOnly>
      <MetricReporter module="blog_article_detail" :biz-id="Number(route.params.id) || 0" />
    </ClientOnly>
    <ClientOnly>
      <BlogHeader />
    </ClientOnly>
    <div class="blog-page-container">
      <div class="blog-content-wrapper">
        <!-- 左侧分类导航 -->
        <aside class="blog-sidebar-left">
          <ClientOnly>
            <BlogCategoryNav :selected-tag-id="currentTagId" @select="handleTagSelect" />
          </ClientOnly>
        </aside>

        <!-- 中间文章内容 -->
        <main class="blog-main">
          <div v-if="detail" class="blog-detail-container">
            <!-- 返回按钮 -->
            <div class="back-link" @click="goBack">← 返回列表</div>
            <h1 class="detail-title">{{ detail.title }}</h1>

            <div class="detail-meta">
              <span class="meta-item">
                <span class="author">{{ detail.authorName || '匿名' }}</span>
              </span>
              <span class="meta-item">
                <span class="time">{{ formatTime(detail.publishTime || 0) }}</span>
              </span>
              <span v-if="wordCount > 0" class="meta-item">
                <span class="word-count">约{{ wordCount }}字</span>
              </span>
              <span v-if="readingTime > 0" class="meta-item">
                <span class="reading-time">大约{{ readingTime }}分钟</span>
              </span>
              <span v-if="detail.tags?.length" class="meta-item">
                <ClientOnly>
                  <el-tag
                    v-for="tag in detail.tags"
                    :key="tag.id"
                    size="small"
                    effect="plain"
                  >
                    {{ tag.name }}
                  </el-tag>
                </ClientOnly>
              </span>
            </div>

            <div v-if="detail.cover" class="detail-cover">
              <img :src="detail.cover" :alt="detail.title" @error="handleImageError" />
            </div>

            <div class="detail-content">
              <ClientOnly>
                <template #default>
                  <component
                    v-if="mdPreviewLoaded && MdPreview && detail?.content"
                    :is="MdPreview"
                    :editor-id="'public-blog-detail'"
                    :model-value="detail.content"
                    :preview-theme="'github'"
                    @html-changed="handleContentRendered"
                  />
                  <div v-else-if="!mdPreviewLoaded || !MdPreview" class="loading-placeholder">
                    加载中... (mdPreviewLoaded: {{ mdPreviewLoaded }}, MdPreview: {{ !!MdPreview }})
                  </div>
                  <div v-else-if="!detail?.content" class="empty-content">
                    暂无内容 (detail: {{ !!detail }}, content: {{ detail?.content?.substring(0, 50) }})
                  </div>
                </template>
                <template #fallback>
                  <div class="loading-placeholder">加载中...</div>
                </template>
              </ClientOnly>
            </div>

            <!-- 相邻文章导航 -->
            <div v-if="prevArticle || nextArticle" class="detail-navigation">
              <NuxtLink
                v-if="prevArticle"
                :to="`/blog/${prevArticle.id}`"
                class="nav-item nav-prev"
              >
                <div class="nav-label">← 上一页</div>
                <div class="nav-title">{{ prevArticle.title }}</div>
              </NuxtLink>
              <div v-else class="nav-item nav-prev disabled">
                <div class="nav-label">← 上一页</div>
                <div class="nav-title">没有更多了</div>
              </div>

              <NuxtLink
                v-if="nextArticle"
                :to="`/blog/${nextArticle.id}`"
                class="nav-item nav-next"
              >
                <div class="nav-label">下一页 →</div>
                <div class="nav-title">{{ nextArticle.title }}</div>
              </NuxtLink>
              <div v-else class="nav-item nav-next disabled">
                <div class="nav-label">下一页 →</div>
                <div class="nav-title">没有更多了</div>
              </div>
            </div>
          </div>

          <div v-else-if="!loading" class="empty">文章不存在或已下架</div>
        </main>

        <!-- 右侧目录 -->
        <aside class="blog-sidebar">
          <ClientOnly>
            <BlogTOC v-if="detail?.content" :content="detail.content" />
          </ClientOnly>
        </aside>
      </div>
    </div>

    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
// Nuxt 3 自动导入 composables，无需手动导入 useRouter、useRoute
import {ref, computed, onMounted, watch, shallowRef, nextTick} from 'vue'
import {ElMessage} from 'element-plus'
// md-editor-v3 只能在客户端使用，使用动态导入避免 SSR 错误
const MdPreview = shallowRef<any>(null)
const mdPreviewLoaded = ref(false)

// 在客户端加载 md-editor-v3
if (process.client) {
  import('md-editor-v3').then((module) => {
    // 检查导出方式：可能是命名导出或默认导出
    MdPreview.value = module.MdPreview || module.default?.MdPreview || module.default
    mdPreviewLoaded.value = true
  }).catch((err) => {
    console.error('加载 md-editor-v3 失败:', err)
  })
  // 使用 style.css（md-editor-v3 没有单独的 preview.css）
  import('md-editor-v3/lib/style.css').catch((err) => {
    console.error('加载 md-editor-v3 样式失败:', err)
  })
}
import {blogApi} from '@/api/blog'
import type {
  PublicBlogArticleDetailReq,
  PublicBlogArticleDetailResp,
  PublicBlogArticlePrevReq,
  PublicBlogArticleNextReq
} from '@/api/generated/admin'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import BlogHeader from '@/components/blog/BlogHeader.vue'
import BlogCategoryNav from '@/components/blog/BlogCategoryNav.vue'
import BlogTOC from '@/components/blog/BlogTOC.vue'

// Nuxt 3 自动导入 useRouter 和 useRoute
const route = useRoute()
const router = useRouter()

// 定义页面元数据（Nuxt 3 规范）
definePageMeta({
  layout: false
})

// 使用 useAsyncData 支持 SSR
const { data: detailData, error: detailError, pending: loading, refresh: refreshDetail } = await useAsyncData(
  `blog-detail-${route.params.id}`,
  async () => {
    const id = Number(route.params.id)
    if (!id) {
      throw new Error('文章ID不正确')
    }
    try {
      const req: PublicBlogArticleDetailReq = {id}
      const resp = await blogApi.publicArticleDetail(req)
      return resp
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '加载文章失败'
      throw new Error(message)
    }
  },
  {
    server: true, // 允许在服务端执行
    default: () => null
  }
)

// 同步 detailData 到 detail ref（保持向后兼容）
const detail = computed(() => detailData.value)

// 调试：监听数据变化
watch(detailData, (newData) => {
  console.log('detailData changed:', newData)
  if (newData?.content) {
    console.log('Content length:', newData.content.length)
  }
}, { immediate: true })

const prevArticle = ref<{id: number; title: string; publishTime: number} | null>(null)
const nextArticle = ref<{id: number; title: string; publishTime: number} | null>(null)

const currentTagId = computed(() => {
  // 从文章标签中获取第一个标签ID作为当前分类
  return detail.value?.tags?.[0]?.id || 0
})

// 计算字数和阅读时间
const wordCount = computed(() => {
  if (!detail.value?.content) {
    return 0
  }
  // 去除Markdown语法，统计纯文本字数
  const text = detail.value.content
    .replace(/#{1,6}\s+/g, '') // 标题
    .replace(/\*\*([^*]+)\*\*/g, '$1') // 加粗
    .replace(/\*([^*]+)\*/g, '$1') // 斜体
    .replace(/`([^`]+)`/g, '') // 行内代码
    .replace(/```[\s\S]*?```/g, '') // 代码块
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1') // 链接
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, '') // 图片
    .replace(/\n/g, '') // 换行
    .trim()
  return text.length
})

const readingTime = computed(() => {
  const count = wordCount.value
  if (count === 0) {
    return 0
  }
  // 按300字/分钟计算
  return Math.max(1, Math.ceil(count / 300))
})

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

const handleContentRendered = (html: string) => {
  // Markdown渲染完成后，为标题添加ID（BlogTOC组件会处理）
  // 这里可以触发TOC更新
  console.log('Content rendered, HTML length:', html?.length || 0)
  console.log('HTML preview:', html?.substring(0, 200) || '')
  
  // 强制隐藏工具栏和编辑器内容（CSS可能被覆盖，使用JS确保隐藏）
  if (process.client) {
    nextTick(() => {
      // 隐藏工具栏
      const toolbarWrapper = document.querySelector('.md-editor-toolbar-wrapper')
      const toolbar = document.querySelector('.md-editor-toolbar')
      if (toolbarWrapper) {
        ;(toolbarWrapper as HTMLElement).style.display = 'none'
        console.log('Toolbar wrapper hidden')
      }
      if (toolbar) {
        ;(toolbar as HTMLElement).style.display = 'none'
        console.log('Toolbar hidden')
      }
      
      // 只隐藏编辑器输入区域，不隐藏 md-editor-content（因为 preview-wrapper 是它的子元素）
      const editorContent = document.querySelector('.md-editor-content')
      const inputWrapper = document.querySelector('.md-editor-input-wrapper')
      // 不要隐藏 editorContent，因为它包含 preview-wrapper
      // if (editorContent) {
      //   ;(editorContent as HTMLElement).style.display = 'none'
      //   console.log('Editor content hidden')
      // }
      if (inputWrapper) {
        ;(inputWrapper as HTMLElement).style.display = 'none'
        console.log('Input wrapper hidden')
      }
      
      // 确保 md-editor-content 可见，这样 preview-wrapper 才能显示
      if (editorContent) {
        const editorContentEl = editorContent as HTMLElement
        editorContentEl.style.display = 'flex'
        editorContentEl.style.flexDirection = 'row'
        editorContentEl.style.width = '100%'
        editorContentEl.style.minWidth = '0'
        editorContentEl.style.boxSizing = 'border-box'
        console.log('Editor content shown (for preview)')
      }
      
      // 隐藏底部工具栏
      const footer = document.querySelector('.md-editor-footer')
      if (footer) {
        ;(footer as HTMLElement).style.display = 'none'
        console.log('Footer hidden')
      }
      
      // 隐藏目录
      const catalog = document.querySelector('.md-editor-catalog')
      if (catalog) {
        ;(catalog as HTMLElement).style.display = 'none'
        console.log('Catalog hidden')
      }
      
      // 确保预览区域可见
      const previewWrapper = document.querySelector('.md-editor-preview-wrapper')
      const preview = document.querySelector('.md-editor-preview')
      const mdEditor = document.querySelector('.md-editor')
      // 查找可能的预览元素 ID（md-editor-v3 可能会根据 editor-id 生成）
      const previewById = document.querySelector('#public-blog-detail-preview')
      
      if (mdEditor) {
        const editorStyles = window.getComputedStyle(mdEditor)
        const editorRect = mdEditor.getBoundingClientRect()
        console.log('md-editor styles:', {
          display: editorStyles.display,
          width: editorStyles.width,
          height: editorStyles.height,
          minHeight: editorStyles.minHeight,
          visibility: editorStyles.visibility,
          opacity: editorStyles.opacity,
          flexDirection: editorStyles.flexDirection,
          rectWidth: editorRect.width
        })
        
        // 获取父容器的宽度
        const parentEl = mdEditor.parentElement
        const parentWidth = parentEl ? parentEl.getBoundingClientRect().width : 0
        
        // 确保 md-editor 容器可见且有宽度
        const mdEditorEl = mdEditor as HTMLElement
        if (editorRect.width === 0 && parentWidth > 0) {
          mdEditorEl.style.width = `${parentWidth}px`
          console.log('Force set md-editor width to:', parentWidth)
        }
        
        // 设置样式（保留已有样式，追加新样式）
        const existingCssText = mdEditorEl.style.cssText || ''
        mdEditorEl.style.cssText = existingCssText + `
          display: block !important;
          visibility: visible !important;
          opacity: 1 !important;
          width: 100% !important;
          min-width: 0 !important;
          box-sizing: border-box !important;
          overflow: visible !important;
          overflow-y: visible !important;
          overflow-x: visible !important;
        `
      }
      
      if (previewWrapper) {
        const wrapperStyles = window.getComputedStyle(previewWrapper)
        const wrapperRect = previewWrapper.getBoundingClientRect()
        console.log('preview-wrapper styles before:', {
          display: wrapperStyles.display,
          height: wrapperStyles.height,
          minHeight: wrapperStyles.minHeight,
          visibility: wrapperStyles.visibility,
          opacity: wrapperStyles.opacity,
          width: wrapperStyles.width,
          rectWidth: wrapperRect.width
        })
        
        // 获取父容器的宽度
        const parentEl = previewWrapper.parentElement
        const parentWidth = parentEl ? parentEl.getBoundingClientRect().width : 0
        
        const previewWrapperEl = previewWrapper as HTMLElement
        
        // 如果宽度为 0，强制设置
        if (wrapperRect.width === 0 && parentWidth > 0) {
          previewWrapperEl.style.width = `${parentWidth}px`
          previewWrapperEl.style.flexBasis = `${parentWidth}px`
          console.log('Force set preview-wrapper width to:', parentWidth)
        }
        
        previewWrapperEl.style.cssText = `
          display: block !important;
          visibility: visible !important;
          opacity: 1 !important;
          width: 100% !important;
          min-width: 0 !important;
          max-width: 100% !important;
          height: auto !important;
          min-height: 200px !important;
          padding: 0 !important;
          background: transparent !important;
          position: relative !important;
          z-index: 1 !important;
          box-sizing: border-box !important;
          flex: 1 1 100% !important;
          flex-grow: 1 !important;
          flex-shrink: 1 !important;
          flex-basis: 100% !important;
          overflow: visible !important;
          overflow-y: visible !important;
          overflow-x: visible !important;
        `
        console.log('Preview wrapper shown')
      }
      if (preview) {
        const previewEl = preview as HTMLElement
        const previewStyles = window.getComputedStyle(previewEl)
        const previewRect = previewEl.getBoundingClientRect()
        console.log('preview styles before:', {
          display: previewStyles.display,
          height: previewStyles.height,
          minHeight: previewStyles.minHeight,
          visibility: previewStyles.visibility,
          opacity: previewStyles.opacity,
          width: previewStyles.width,
          color: previewStyles.color,
          fontSize: previewStyles.fontSize,
          rectWidth: previewRect.width
        })
        
        // 获取父容器的宽度
        const parentEl = previewEl.parentElement
        const parentWidth = parentEl ? parentEl.getBoundingClientRect().width : 0
        
        // 如果宽度为 0，强制设置
        if (previewRect.width === 0 && parentWidth > 0) {
          previewEl.style.width = `${parentWidth}px`
          previewEl.style.maxWidth = `${parentWidth}px`
          console.log('Force set preview width to:', parentWidth)
        }
        
        previewEl.style.cssText = `
          display: block !important;
          visibility: visible !important;
          opacity: 1 !important;
          width: 100% !important;
          min-width: 0 !important;
          max-width: 100% !important;
          height: auto !important;
          min-height: 200px !important;
          padding: 0 !important;
          background: transparent !important;
          color: #333 !important;
          font-size: 16px !important;
          line-height: 1.8 !important;
          position: relative !important;
          z-index: 1 !important;
          box-sizing: border-box !important;
          flex: 1 1 auto !important;
          flex-grow: 1 !important;
          flex-shrink: 1 !important;
          flex-basis: auto !important;
          overflow: visible !important;
          overflow-y: visible !important;
          overflow-x: visible !important;
        `
        
        // 移除预览元素的滚动条（如果有 ID 为 public-blog-detail-preview 的元素）
        if (previewById) {
          const previewByIdEl = previewById as HTMLElement
          previewByIdEl.style.overflow = 'visible'
          previewByIdEl.style.overflowY = 'visible'
          previewByIdEl.style.overflowX = 'visible'
          console.log('Removed scrollbar from preview element by ID')
        }
        
        console.log('Preview shown, innerHTML length:', preview.innerHTML.length)
        console.log('Preview innerHTML preview:', preview.innerHTML.substring(0, 200))
        
        // 确保所有子元素可见 - 关键：必须使用 setProperty 设置 !important
        const allChildren = preview.querySelectorAll('*')
        allChildren.forEach((child) => {
          const el = child as HTMLElement
          // 关键：使用 setProperty 设置 !important，覆盖所有 CSS 规则
          if (el.tagName === 'P' || el.tagName === 'DIV' || el.tagName === 'SPAN') {
            el.style.setProperty('display', 'block', 'important')
          } else if (el.tagName === 'H1' || el.tagName === 'H2' || el.tagName === 'H3' || el.tagName === 'H4' || el.tagName === 'H5' || el.tagName === 'H6') {
            el.style.setProperty('display', 'block', 'important')
          } else if (el.tagName === 'UL' || el.tagName === 'OL') {
            el.style.setProperty('display', 'block', 'important')
          } else if (el.tagName === 'LI') {
            el.style.setProperty('display', 'list-item', 'important')
          } else if (el.tagName === 'IMG') {
            el.style.setProperty('display', 'block', 'important')
          } else if (el.tagName === 'BR') {
            el.style.setProperty('display', 'inline', 'important')
          } else {
            el.style.setProperty('display', 'block', 'important')
          }
          el.style.setProperty('visibility', 'visible', 'important')
          el.style.setProperty('opacity', '1', 'important')
          if (el.tagName === 'H1' || el.tagName === 'H2' || el.tagName === 'H3' || el.tagName === 'H4' || el.tagName === 'H5' || el.tagName === 'H6') {
            el.style.setProperty('color', '#333', 'important')
            el.style.setProperty('font-size', el.tagName === 'H1' ? '28px' : el.tagName === 'H2' ? '24px' : '20px', 'important')
            el.style.setProperty('font-weight', '600', 'important')
            el.style.setProperty('margin-top', '24px', 'important')
            el.style.setProperty('margin-bottom', '12px', 'important')
          }
          if (el.tagName === 'P') {
            el.style.setProperty('color', '#333', 'important')
            el.style.setProperty('margin-bottom', '16px', 'important')
            el.style.setProperty('font-size', '16px', 'important')
            el.style.setProperty('line-height', '1.8', 'important')
          }
        })
        console.log('Made', allChildren.length, 'children visible with display set (important)')
        
        // 强制触发布局计算，确保宽度正确
        if (preview) {
          const previewEl = preview as HTMLElement
          // 强制重新计算布局
          void previewEl.offsetWidth
          
          // 检查所有父容器，确保它们都有宽度
          let currentEl: HTMLElement | null = previewEl.parentElement
          while (currentEl && currentEl !== document.body) {
            const rect = currentEl.getBoundingClientRect()
            const styles = window.getComputedStyle(currentEl)
            console.log('Parent element:', currentEl.className, {
              width: rect.width,
              display: styles.display,
              flexDirection: styles.flexDirection
            })
            
            // 如果父容器有宽度但预览区域没有，强制设置
            if (rect.width > 0) {
              const previewRect = previewEl.getBoundingClientRect()
              if (previewRect.width === 0) {
                // 尝试从父容器获取宽度
                previewEl.style.width = `${rect.width}px`
                previewEl.style.maxWidth = `${rect.width}px`
                console.log('Force set preview width to:', rect.width)
                
                // 也设置预览包装器的宽度
                if (previewWrapper) {
                  const wrapperEl = previewWrapper as HTMLElement
                  const widthValue = `${Number(rect.width)}px`
                  wrapperEl.style.width = widthValue
                  wrapperEl.style.maxWidth = widthValue
                }
                break
              }
            }
            currentEl = currentEl.parentElement
          }
        }
        
        // 检查第一个子元素和预览区域的实际尺寸
        const firstChild = preview?.firstElementChild as HTMLElement
        if (firstChild && preview) {
          const childStyles = window.getComputedStyle(firstChild)
          const previewRect = preview.getBoundingClientRect()
          const firstChildRect = firstChild.getBoundingClientRect()
          console.log('First child styles after:', {
            display: childStyles.display,
            visibility: childStyles.visibility,
            opacity: childStyles.opacity,
            color: childStyles.color,
            fontSize: childStyles.fontSize,
            height: childStyles.height,
            width: childStyles.width
          })
          console.log('Preview element rect:', {
            top: previewRect.top,
            left: previewRect.left,
            width: previewRect.width,
            height: previewRect.height
          })
          console.log('First child rect:', {
            top: firstChildRect.top,
            left: firstChildRect.left,
            width: firstChildRect.width,
            height: firstChildRect.height
          })
        }
      }
    })
  }
}

const handleTagSelect = (tagId: number) => {
  router.push({
    path: '/blog',
    query: {tagId: tagId > 0 ? tagId : undefined}
  })
}

// 返回列表
const goBack = () => {
  // 优先走浏览器历史，保留列表分页与滚动状态
  // 如果是从列表页进入的，router.back() 会恢复列表页的状态（包括滚动位置）
  if (typeof window !== 'undefined' && window.history.length > 1) {
    router.back()
  } else {
    // 如果没有历史记录，尝试从 sessionStorage 恢复状态
    if (typeof window !== 'undefined') {
      try {
        const raw = sessionStorage.getItem('public_blog_list_state')
        if (raw) {
          const parsed = JSON.parse(raw) as {
            page?: number
            size?: number
            keyword?: string
            tagId?: number
            scrollTop?: number
          }
          router.push({
            path: '/blog',
            query: {
              ...(parsed.page && {page: String(parsed.page)}),
              ...(parsed.size && {size: String(parsed.size)}),
              ...(parsed.keyword && {keyword: parsed.keyword}),
              ...(parsed.tagId && parsed.tagId > 0 && {tagId: String(parsed.tagId)})
            }
          })
          return
        }
      } catch {
        // 忽略解析错误
      }
    }
    // 兜底：直接跳转到列表页
    router.push('/blog')
  }
}

// 先定义加载相邻文章的函数（必须在 watch 之前定义）
const loadPrevArticle = async (currentId: number) => {
  try {
    const req: PublicBlogArticlePrevReq = {id: currentId}
    const resp = await blogApi.publicArticlePrev(req)
    if (resp.id) {
      prevArticle.value = {
        id: resp.id,
        title: resp.title,
        publishTime: resp.publishTime
      }
    } else {
      prevArticle.value = null
    }
  } catch (err) {
    console.error('加载上一篇文章失败:', err)
    prevArticle.value = null
  }
}

const loadNextArticle = async (currentId: number) => {
  try {
    const req: PublicBlogArticleNextReq = {id: currentId}
    const resp = await blogApi.publicArticleNext(req)
    if (resp.id) {
      nextArticle.value = {
        id: resp.id,
        title: resp.title,
        publishTime: resp.publishTime
      }
    } else {
      nextArticle.value = null
    }
  } catch (err) {
    console.error('加载下一篇文章失败:', err)
    nextArticle.value = null
  }
}

// 加载相邻文章（在详情加载成功后）
watch(detailData, async (newDetail) => {
  if (newDetail?.id) {
    await Promise.all([loadPrevArticle(newDetail.id), loadNextArticle(newDetail.id)])
  }
}, { immediate: true })

// 处理错误
watch(detailError, (err) => {
  if (err) {
    console.error('加载文章详情失败:', err)
    if (process.client) {
      ElMessage.error(err.message || '加载文章失败')
    }
  }
}, { immediate: true })

// 监听路由参数变化，重新加载文章详情
watch(
  () => route.params.id,
  (newId, oldId) => {
    // 只有当 ID 真正变化时才重新加载
    if (newId && newId !== oldId) {
      // 重置状态
      prevArticle.value = null
      nextArticle.value = null
      // 重新加载
      refreshDetail()
    }
  },
  {immediate: false}
)
</script>

<style scoped lang="scss">
@import '@/assets/styles/blog.scss';

.blog-detail-page {
  height: 100vh; // 固定高度，不使用 min-height，避免出现滚动条
  overflow: hidden; // 隐藏最外层滚动条
  background: #f5f5f5;
  position: relative; // 确保定位上下文

  // 返回链接
  .back-link {
    cursor: pointer;
    color: #666;
    margin-bottom: 12px;
    font-size: 14px;
    display: inline-block;
    transition: color 0.3s;
    flex-shrink: 0; // 不参与滚动，固定在顶部

    &:hover {
      color: #409eff;
    }
  }

  .detail-navigation {
    .nav-item.disabled {
      cursor: not-allowed;
      opacity: 0.5;

      &:hover {
        background: #f5f5f5;
        color: #666;
      }
    }
  }

  // 确保 MdPreview 只显示预览，隐藏编辑器工具栏和源码区
  .detail-content {
    min-height: 200px; // 确保有最小高度
    width: 100% !important; // 确保有宽度
    height: auto !important; // 高度自适应，让内容自然撑开
    font-size: 16px;
    line-height: 1.8;
    color: #333;
    box-sizing: border-box; // 确保宽度计算正确
    // 布局相关样式已在 blog.scss 中定义，这里不再重复
    // 不使用滚动，由父容器 blog-main 处理滚动
    overflow: visible !important; // 让内容自然显示

    // 确保 md-editor 容器本身可见且有宽度
    :deep(.md-editor) {
      border: none !important;
      box-shadow: none !important;
      background: transparent !important;
      display: block !important; // 改为 block，不使用 flex
      width: 100% !important; // 确保有宽度
      min-width: 0 !important; // 防止 flex 子元素溢出
      min-height: 200px !important;
      height: auto !important; // 高度自适应
      box-sizing: border-box !important; // 确保宽度计算正确
      flex-direction: column !important; // 如果使用 flex，确保纵向排列
      overflow: visible !important; // 移除滚动条
      overflow-y: visible !important; // 移除垂直滚动条
      overflow-x: visible !important; // 移除水平滚动条
    }
    
    // 如果 md-editor 内部使用了 flex 布局，确保预览区域能正确计算宽度
    :deep(.md-editor[style*="display: flex"]) {
      display: block !important; // 强制改为 block
    }

    // 隐藏工具栏 - 使用更具体的选择器
    :deep(.md-editor-toolbar-wrapper),
    :deep(.md-editor-toolbar-wrapper *),
    :deep(.md-editor-toolbar),
    :deep(.md-editor-toolbar *) {
      display: none !important;
      visibility: hidden !important;
    }

    // 只隐藏编辑器输入区域，不隐藏 md-editor-content（因为 preview-wrapper 是它的子元素）
    :deep(.md-editor-input-wrapper),
    :deep(.md-editor-input-wrapper *) {
      display: none !important;
      visibility: hidden !important;
      width: 0 !important;
      min-width: 0 !important;
      flex: none !important;
    }
    
    // 确保 md-editor-content 可见，这样 preview-wrapper 才能显示
    :deep(.md-editor-content) {
      display: block !important; // 改为 block，不使用 flex，让内容自然流动
      width: 100% !important;
      min-width: 0 !important;
      height: auto !important; // 高度自适应
    }

    // 隐藏底部工具栏（字数统计等）
    :deep(.md-editor-footer),
    :deep(.md-editor-footer *) {
      display: none !important;
      visibility: hidden !important;
    }

    // 隐藏目录（右侧）
    :deep(.md-editor-catalog),
    :deep(.md-editor-catalog *) {
      display: none !important;
      visibility: hidden !important;
    }

    // 确保预览区域显示且有宽度
    :deep(.md-editor-preview-wrapper) {
      display: block !important;
      width: 100% !important;
      min-width: 0 !important; // 防止 flex 子元素溢出
      max-width: 100% !important; // 确保不超出容器
      padding: 0 !important;
      background: transparent !important;
      border: none !important;
      box-shadow: none !important;
      visibility: visible !important;
      opacity: 1 !important;
      min-height: 200px !important;
      height: auto !important; // 高度自适应，让内容自然撑开
      position: relative !important;
      z-index: 1 !important;
      box-sizing: border-box !important; // 确保宽度计算正确
      overflow: visible !important; // 移除滚动条
      overflow-y: visible !important; // 移除垂直滚动条
      overflow-x: visible !important; // 移除水平滚动条
    }

    // 预览内容样式
    :deep(.md-editor-preview),
    :deep(#public-blog-detail-preview) {
      display: block !important;
      width: 100% !important;
      min-width: 0 !important; // 防止 flex 子元素溢出
      max-width: 100% !important; // 确保不超出容器
      padding: 0 !important;
      background: transparent !important;
      border: none !important;
      box-shadow: none !important;
      visibility: visible !important;
      opacity: 1 !important;
      min-height: 200px !important;
      height: auto !important; // 高度自适应，让内容自然撑开
      position: relative !important;
      z-index: 1 !important;
      color: #333 !important;
      box-sizing: border-box !important; // 确保宽度计算正确
      overflow: visible !important; // 移除滚动条
      overflow-y: visible !important; // 移除垂直滚动条
      overflow-x: visible !important; // 移除水平滚动条
      
      // 确保所有子元素可见 - 关键：必须设置 display
      * {
        visibility: visible !important;
        opacity: 1 !important;
        color: inherit !important;
      }
      
      // 确保段落和标题显示 - 覆盖可能的 display: none
      p, div, span, h1, h2, h3, h4, h5, h6, ul, ol, li, img {
        display: block !important;
      }
      
      li {
        display: list-item !important;
      }
      
      br {
        display: inline !important;
      }
      
      h1, h2, h3, h4, h5, h6 {
        margin-top: 24px;
        margin-bottom: 12px;
        font-weight: 600;
        color: #333;
      }

      h1 {
        font-size: 28px;
      }
      h2 {
        font-size: 24px;
      }
      h3 {
        font-size: 20px;
      }

      p {
        margin-bottom: 16px;
      }

      img {
        max-width: 100%;
        height: auto;
        border-radius: 8px;
        margin: 16px 0;
        display: block;
      }

      code {
        background: #f5f5f5;
        padding: 2px 6px;
        border-radius: 4px;
        font-size: 14px;
        color: #e83e8c;
      }

      pre {
        background: #f5f5f5;
        padding: 16px;
        border-radius: 8px;
        overflow-x: auto;
        margin: 16px 0;

        code {
          background: transparent;
          padding: 0;
          color: #333;
        }
      }

      blockquote {
        border-left: 4px solid #409eff;
        padding-left: 16px;
        margin: 16px 0;
        color: #666;
      }

      ul, ol {
        margin: 16px 0;
        padding-left: 24px;
      }

      li {
        margin: 8px 0;
      }

      table {
        width: 100%;
        border-collapse: collapse;
        margin: 16px 0;

        th, td {
          border: 1px solid #e0e0e0;
          padding: 8px 12px;
          text-align: left;
        }

        th {
          background: #f5f5f5;
          font-weight: 600;
        }
      }
    }

    .loading-placeholder,
    .empty-content {
      text-align: center;
      padding: 40px 0;
      color: #999;
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .blog-detail-page {
    .blog-page-container {
      padding-top: 50px;
    }
  }
}
</style>
