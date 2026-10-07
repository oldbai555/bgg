<template>
  <nav class="fitness-tabbar">
    <router-link
      v-for="tab in tabs"
      :key="tab.path"
      :to="tab.path"
      class="fitness-tabbar__item"
      :class="{'is-active': route.path === tab.path}"
    >
      <el-icon :size="20"><component :is="tab.icon" /></el-icon>
      <span>{{ t(tab.label) }}</span>
    </router-link>
  </nav>
</template>

<script setup lang="ts">
import {useRoute} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {Calendar, CircleCheck, TrendCharts, User} from '@element-plus/icons-vue'

const route = useRoute()
const {t} = useI18n()

const tabs = [
  {path: '/front/fitness', label: 'fitness.tabToday', icon: Calendar},
  {path: '/front/fitness/checkin', label: 'fitness.tabCheckin', icon: CircleCheck},
  {path: '/front/fitness/body', label: 'fitness.tabBody', icon: TrendCharts},
  {path: '/front/fitness/me', label: 'fitness.tabMe', icon: User}
]
</script>

<style scoped lang="scss">
// 桌面端作为页面顶部的二级导航，移动端固定在底部（页面需给 .page-shell 预留底部空间）
.fitness-tabbar {
  display: flex;
  gap: 8px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px 0;

  &__item {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 8px 14px;
    border-radius: 999px;
    color: var(--color-text-regular);
    background: var(--color-bg-card);
    border: 1px solid var(--color-border-light);
    text-decoration: none;
    font-size: 14px;

    &.is-active {
      color: #fff;
      background: var(--color-primary);
      border-color: var(--color-primary);
    }
  }
}

@include mobile {
  .fitness-tabbar {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 100;
    gap: 0;
    max-width: none;
    padding: 6px 0 calc(6px + env(safe-area-inset-bottom));
    background: var(--color-bg-primary);
    border-top: 1px solid var(--color-border);

    &__item {
      flex: 1;
      flex-direction: column;
      gap: 2px;
      padding: 2px 0;
      border: none;
      border-radius: 0;
      background: transparent;
      font-size: 11px;

      &.is-active {
        color: var(--color-primary);
        background: transparent;
      }
    }
  }
}
</style>
