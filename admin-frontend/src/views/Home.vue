<template>
  <div class="home-page">
    <main class="home">
      <span class="home__mark">BGG</span>

      <h1 class="home__path" :style="pathStyle" :aria-label="`${origin}${typed}`">
        <span class="home__path-base">{{ origin }}</span><span class="home__path-tail">{{ typed }}</span><span class="home__caret" aria-hidden="true"></span>
      </h1>
      <p class="home__lead">博客、视频、每日训练计划和管理后台，都从这里进。</p>

      <section v-for="group in groups" :key="group.prefix" class="home__group">
        <div class="home__group-label">{{ group.prefix }}</div>
        <router-link
          v-for="(entry, idx) in group.entries"
          :key="entry.path"
          :to="entry.path"
          class="entry"
          :style="{'--entry-color': entry.color, '--i': group.offset + idx}"
          @mouseenter="point(entry)"
          @mouseleave="point(null)"
          @focus="point(entry)"
          @blur="point(null)"
        >
          <span class="entry__glyph" aria-hidden="true">{{ entry.glyph }}</span>
          <span class="entry__text">
            <span class="entry__title">{{ entry.title }}</span>
            <span class="entry__desc">{{ entry.desc }}</span>
          </span>
          <span class="entry__route" aria-hidden="true">{{ entry.path }}</span>
          <span class="entry__arrow" aria-hidden="true">→</span>
        </router-link>
      </section>
    </main>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {computed, onBeforeUnmount, ref} from 'vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import {FITNESS_HOME_PATH} from '@/constants/fitness'

interface Entry {
  path: string
  title: string
  desc: string
  glyph: string
  color: string
}

const frontEntries: Entry[] = [
  {path: '/front/blog', title: '博客', desc: '技术文章，按标签分类，支持搜索', glyph: '博', color: 'var(--color-primary)'},
  {path: '/front/videos', title: '视频', desc: '视频列表，悬停预览片段', glyph: '视', color: 'var(--color-warning)'},
  {path: FITNESS_HOME_PATH, title: '今日计划', desc: '训练与饮食打卡 · 飞书登录', glyph: '练', color: 'var(--color-success)'}
]
const adminEntries: Entry[] = [
  {path: '/admin/login', title: '后台', desc: '管理系统 · 账号或飞书登录', glyph: '管', color: 'var(--color-text-secondary)'}
]
const groups = [
  {prefix: '/front', entries: frontEntries, offset: 0},
  {prefix: '/admin', entries: adminEntries, offset: frontEntries.length}
]

// 地址栏前缀跟随实际部署域名与 base（线上 oldbai.top/bgg/，本地 localhost:5173/bgg/）
const origin = `${window.location.host}${import.meta.env.BASE_URL}`
const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

const typed = ref('')
const activeColor = ref('')
const pathStyle = computed(() => (activeColor.value ? {'--caret-color': activeColor.value, '--tail-color': activeColor.value} : undefined))

let target = ''
let timer: number | undefined

const step = () => {
  const current = typed.value
  if (current === target) {
    window.clearInterval(timer)
    timer = undefined
    if (!target) {
      activeColor.value = ''
    }
    return
  }
  typed.value = target.startsWith(current) ? target.slice(0, current.length + 1) : current.slice(0, -1)
}

const point = (entry: Entry | null) => {
  target = entry ? entry.path.slice(1) : ''
  if (entry) {
    activeColor.value = entry.color
  }
  if (reduceMotion) {
    typed.value = target
    activeColor.value = entry ? entry.color : ''
    return
  }
  if (timer === undefined) {
    timer = window.setInterval(step, 16)
  }
}

onBeforeUnmount(() => window.clearInterval(timer))
</script>

<style scoped lang="scss">
$mono: ui-monospace, 'SF Mono', 'JetBrains Mono', Menlo, Consolas, monospace;

.home-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
  box-sizing: border-box;
}

.home {
  flex: 1;
  width: 100%;
  max-width: 760px;
  margin: 0 auto;
  padding: 72px 24px 32px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: center;

  &__mark {
    font-family: $mono;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.24em;
    color: var(--color-primary);
  }

  &__path {
    margin: 20px 0 0;
    font-family: $mono;
    font-size: clamp(22px, 3.6vw, 36px);
    font-weight: 500;
    line-height: 1.2;
    letter-spacing: -0.02em;
    overflow-wrap: anywhere;
    min-height: 1.2em;
  }

  &__path-base {
    color: var(--color-text-secondary);
  }

  &__path-tail {
    color: var(--tail-color, var(--color-text-primary));
  }

  &__caret {
    display: inline-block;
    width: 0.5em;
    height: 1em;
    margin-left: 2px;
    vertical-align: -0.12em;
    background: var(--caret-color, var(--color-primary));
    transition: background-color 0.18s;
    animation: caret-blink 1.1s steps(1) infinite;
  }

  &__lead {
    margin: 14px 0 48px;
    font-size: 15px;
    line-height: 1.6;
    color: var(--color-text-regular);
  }

  &__group + &__group {
    margin-top: 28px;
  }

  &__group-label {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 10px;
    font-family: $mono;
    font-size: 12px;
    color: var(--color-text-secondary);

    &::after {
      content: '';
      flex: 1;
      height: 1px;
      background: var(--color-border);
    }
  }
}

.entry {
  display: grid;
  grid-template-columns: 44px 1fr auto 20px;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  margin-bottom: 8px;
  border: 1px solid var(--color-border-light);
  border-radius: 10px;
  background: var(--color-bg-card);
  color: inherit;
  text-decoration: none;
  transition: background-color 0.18s, border-color 0.18s;
  animation: entry-in 0.4s cubic-bezier(0.2, 0.7, 0.2, 1) both;
  animation-delay: calc(var(--i) * 60ms + 120ms);

  &__glyph {
    width: 44px;
    height: 44px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    font-weight: 600;
    color: var(--entry-color);
    background: color-mix(in srgb, var(--entry-color) 12%, transparent);
    transition: background-color 0.18s, color 0.18s;
  }

  &__text {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  &__title {
    font-size: 17px;
    font-weight: 600;
  }

  &__desc {
    font-size: 13px;
    color: var(--color-text-secondary);
  }

  &__route {
    font-family: $mono;
    font-size: 13px;
    color: var(--color-text-secondary);
    transition: color 0.18s;
  }

  &__arrow {
    font-size: 16px;
    color: var(--color-text-secondary);
    transition: transform 0.18s, color 0.18s;
  }

  &:hover,
  &:focus-visible {
    border-color: color-mix(in srgb, var(--entry-color) 45%, transparent);
    background: color-mix(in srgb, var(--entry-color) 6%, var(--color-bg-card));

    .entry__glyph {
      color: #fff;
      background: var(--entry-color);
    }

    .entry__route,
    .entry__arrow {
      color: var(--entry-color);
    }

    .entry__arrow {
      transform: translateX(4px);
    }
  }

  &:focus-visible {
    outline: 2px solid var(--entry-color);
    outline-offset: 2px;
  }
}

@keyframes caret-blink {
  50% {
    opacity: 0;
  }
}

@keyframes entry-in {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
}

@media (prefers-reduced-motion: reduce) {
  .home__caret,
  .entry {
    animation: none;
  }
}

@include mobile {
  .home {
    padding: 40px 16px 24px;
    justify-content: flex-start;

    &__lead {
      margin-bottom: 32px;
      font-size: 14px;
    }
  }

  .entry {
    grid-template-columns: 40px 1fr 16px;
    gap: 12px;
    padding: 12px;

    &__glyph {
      width: 40px;
      height: 40px;
      font-size: 16px;
    }

    &__route {
      display: none;
    }
  }
}
</style>
