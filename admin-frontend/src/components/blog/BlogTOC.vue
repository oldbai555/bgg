<template>
  <aside class="blog-toc">
    <div class="toc-title">此页内容</div>
    <div v-if="tocItems.length > 0" class="toc-list">
      <div
        v-for="item in tocItems"
        :key="item.id"
        class="toc-item"
        :class="{
          [`level-${item.level}`]: true,
          active: item.id === activeId
        }"
        @click="scrollToAnchor(item.id)"
      >
        <span class="toc-text">{{ item.text }}</span>
      </div>
    </div>
    <div v-else class="toc-empty">暂无目录</div>
  </aside>
</template>

<script setup lang="ts">
import {ref, onMounted, onUnmounted, watch} from 'vue'

interface TOCItem {
  id: string
  text: string
  level: number
}

const props = defineProps<{
  content: string
}>()

const tocItems = ref<TOCItem[]>([])
const activeId = ref<string>('')

// 从Markdown内容提取标题
const extractTOC = (content: string): TOCItem[] => {
  const items: TOCItem[] = []
  const lines = content.split('\n')
  let idCounter = 0

  lines.forEach((line) => {
    const match = line.match(/^(#{1,6})\s+(.+)$/)
    if (match) {
      const level = match[1].length
      const text = match[2].trim()
      const id = `toc-${idCounter++}`
      items.push({id, text, level})
    }
  })

  return items
}

// 滚动到锚点
const scrollToAnchor = (id: string) => {
  // 确保在客户端环境
  if (typeof window === 'undefined' || typeof document === 'undefined') {
    return
  }
  const element = document.getElementById(id)
  if (element) {
    element.scrollIntoView({behavior: 'smooth', block: 'start'})
    activeId.value = id
  }
}

// 监听滚动，高亮当前阅读位置
const handleScroll = () => {
  // 确保在客户端环境
  if (typeof window === 'undefined' || typeof document === 'undefined') {
    return
  }
  const headings = tocItems.value.map((item) => ({
    id: item.id,
    element: document.getElementById(item.id)
  })).filter((h) => h.element !== null)

  if (headings.length === 0) {
    return
  }

  // 找到当前视口中最接近顶部的标题
  let currentId = ''

  for (let i = headings.length - 1; i >= 0; i--) {
    const heading = headings[i]
    if (heading.element) {
      const rect = heading.element.getBoundingClientRect()
      if (rect.top <= 150) {
        currentId = heading.id
        break
      }
    }
  }

  if (currentId) {
    activeId.value = currentId
  } else if (headings.length > 0) {
    // 如果滚动到顶部，高亮第一个
    activeId.value = headings[0].id
  }
}

// 为标题添加ID（需要在Markdown渲染后调用）
const addHeadingIds = () => {
  // 确保在客户端环境
  if (typeof window === 'undefined' || typeof document === 'undefined') {
    return
  }
  const headings = document.querySelectorAll('.md-editor-preview h1, .md-editor-preview h2, .md-editor-preview h3, .md-editor-preview h4, .md-editor-preview h5, .md-editor-preview h6')
  let index = 0
  headings.forEach((heading) => {
    const text = heading.textContent?.trim() || ''
    // 找到匹配的TOC项
    const matchedItem = tocItems.value.find((item) => item.text === text)
    if (matchedItem && !heading.id) {
      heading.id = matchedItem.id
    } else if (!heading.id) {
      // 如果没有匹配的TOC项，使用索引创建ID
      heading.id = `toc-heading-${index++}`
    }
  })
}

watch(() => props.content, (newContent) => {
  if (!newContent || typeof newContent !== 'string') {
    tocItems.value = []
    return
  }
  
  const extracted = extractTOC(newContent)
  tocItems.value = extracted
  
  // 延迟执行，等待Markdown渲染完成
  if (typeof window !== 'undefined') {
    setTimeout(() => {
      addHeadingIds()
      handleScroll()
    }, 500)
  }
}, {immediate: true})

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('scroll', handleScroll, {passive: true})
    setTimeout(() => {
      addHeadingIds()
      handleScroll()
    }, 500)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('scroll', handleScroll)
  }
})
</script>

<style scoped lang="scss">
.blog-toc {
  position: relative; // 改为 relative，因为父容器 blog-sidebar 已经使用了 sticky
  // 移除 top: 80px，由父容器控制定位
  width: 100%; // 改为 100%，占满父容器 blog-sidebar 的宽度
  max-width: 240px; // 最大宽度限制为 240px
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  max-height: calc(100vh - 100px);
  overflow-y: auto;
  overflow-x: visible; // 确保横向不裁剪
  box-sizing: border-box; // 确保 padding 包含在宽度内

  .toc-title {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f0f0f0;
  }

  .toc-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .toc-item {
    padding: 6px 8px;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;

    .toc-text {
      font-size: 14px;
      color: #666;
      line-height: 1.5;
    }

    &:hover {
      background: #f5f5f5;

      .toc-text {
        color: #409eff;
      }
    }

    &.active {
      background: #e6f4ff;

      .toc-text {
        color: #409eff;
        font-weight: 500;
      }
    }

    // 不同层级的缩进
    &.level-1 {
      padding-left: 8px;
    }
    &.level-2 {
      padding-left: 16px;
    }
    &.level-3 {
      padding-left: 24px;
    }
    &.level-4 {
      padding-left: 32px;
    }
    &.level-5 {
      padding-left: 40px;
    }
    &.level-6 {
      padding-left: 48px;
    }
  }

  .toc-empty {
    font-size: 14px;
    color: #999;
    text-align: center;
    padding: 20px 0;
  }
}

// 移动端适配
@include mobile {
  .blog-toc {
    display: none; // 移动端隐藏，通过浮层显示
  }
}
</style>
