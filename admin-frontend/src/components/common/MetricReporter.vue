<template>
  <!-- 纯逻辑组件：不渲染任何内容 -->
  <span v-if="false"></span>
</template>

<script setup lang="ts">
import {onMounted, watch} from 'vue'
import {metricApi} from '@/api/metric'

type MetricEvent = 'view' | 'play' | string;

const props = withDefaults(
  defineProps<{
    /** 业务模块标识，如 blog_article_list/blog_article_detail/video_list/video_detail */
    module: string;
    /** 业务ID（列表页为0，详情页为具体ID） */
    bizId?: number;
    /** 事件类型：默认 view */
    event?: MetricEvent;
    /** 是否启用（方便调试/灰度） */
    enabled?: boolean;
    /** bizId 变化时是否再次上报（详情页路由切换场景有用） */
    watchBizId?: boolean;
  }>(),
  {
    bizId: 0,
    event: 'view',
    enabled: true,
    watchBizId: true
  }
)

const report = () => {
  const module = (props.module || '').trim()
  if (!props.enabled || !module) {
return
}
  const bizId = Number(props.bizId || 0)
  metricApi
    .report({
      module,
      bizId,
      event: props.event
    })
    .catch(() => {})
}

onMounted(() => {
  report()
})

watch(
  () => props.bizId,
  () => {
    if (!props.watchBizId) {
return
}
    report()
  }
)
</script>

