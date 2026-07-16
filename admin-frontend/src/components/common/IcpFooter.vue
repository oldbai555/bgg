<template>
  <footer class="icp-footer">
    <a
      class="icp-footer__link"
      :href="href"
      target="_blank"
      rel="noopener noreferrer"
    >
      {{ text }}
    </a>
  </footer>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    /**
     * 展示文本，例如：icp: 桂ICP备2023001017号-1
     */
    text?: string
    /**
     * 备案查询链接（默认工信部备案系统）
     */
    href?: string
  }>(),
  {
    text: 'ICP: 桂ICP备2023001017号-1',
    href: 'https://beian.miit.gov.cn/'
  }
)
</script>

<style scoped lang="scss">
.icp-footer {
  width: 100%;
  padding: 12px 16px;
  text-align: center;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.2;
}

.icp-footer__link {
  color: inherit;
  text-decoration: none;
}

.icp-footer__link:hover {
  text-decoration: underline;
}
</style>

