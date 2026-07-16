<template>
  <!-- 纯逻辑组件：不渲染任何内容 -->
  <span v-if="false"></span>
</template>

<script setup lang="ts">
import {onMounted, watch} from 'vue'
import {monitoringApi} from '@/api/monitoring'

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

/** 命令式触发时可覆盖 event/bizId（如"视频真正开始播放"这类一次性业务事件，不适合用 props 声明式表达） */
const report = (override?: {event?: MetricEvent; bizId?: number}) => {
  const module = (props.module || '').trim()
  if (!props.enabled || !module) {
    return
  }
  const bizId = Number(override?.bizId ?? props.bizId ?? 0)
  const event = override?.event ?? props.event
  monitoringApi
    .metricReport({
      module,
      bizId,
      event
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

defineExpose({report})
</script>

