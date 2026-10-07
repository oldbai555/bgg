<template>
  <div class="fitness-heatmap">
    <div class="fitness-heatmap__grid" :style="{gridTemplateColumns: `24px repeat(${weeks.length}, 14px)`}">
      <template v-for="w in 7" :key="`l${w}`">
        <span class="fitness-heatmap__label" :style="{gridRow: w, gridColumn: 1}">{{ w % 2 === 1 ? weekdayLabel(w).slice(1) : '' }}</span>
      </template>
      <template v-for="(week, wi) in weeks" :key="wi">
        <span
          v-for="cell in week"
          :key="cell.date"
          class="fitness-heatmap__cell"
          :class="cellClass(cell)"
          :style="{gridRow: weekdayOf(cell.date), gridColumn: wi + 2}"
          :title="cellTitle(cell)"
        ></span>
      </template>
    </div>
    <div class="fitness-heatmap__legend">
      <span><i class="lv-none"></i>未打卡</span>
      <span><i class="lv-rest"></i>休息日</span>
      <span><i class="lv-1"></i>部分完成</span>
      <span><i class="lv-3"></i>完成</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import type {FitnessCalendarDay} from '@/api/generated/admin'
import {useDictOptions} from '@/composables/useDictOptions'
import {formatRate, weekdayLabel, weekdayOf} from '@/utils/fitness'

const props = defineProps<{days: FitnessCalendarDay[]}>()
const {getLabel} = useDictOptions('fitness_training_type')

/** 按周切列：遇到周一另起一列 */
const weeks = computed(() => {
  const out: FitnessCalendarDay[][] = []
  props.days.forEach((d) => {
    if (out.length === 0 || weekdayOf(d.date) === 1) {
      out.push([])
    }
    out[out.length - 1].push(d)
  })
  return out
})

const cellClass = (d: FitnessCalendarDay): string => {
  if (!d.counted) {
    return d.checked ? 'lv-rest is-checked' : 'lv-rest'
  }
  if (!d.checked || d.score <= 0) {
    return 'lv-none'
  }
  if (d.score >= 1) {
    return 'lv-3'
  }
  return d.score >= 0.5 ? 'lv-2' : 'lv-1'
}

const cellTitle = (d: FitnessCalendarDay): string => {
  const type = d.trainingType ? getLabel(d.trainingType) : '无计划'
  const state = d.counted ? `完成度 ${formatRate(d.score)}` : '不计入完成率'
  return `${d.date} ${type}｜${d.checked ? '已打卡' : '未打卡'}｜${state}`
}
</script>

<style scoped>
.fitness-heatmap {
  overflow-x: auto;
}

.fitness-heatmap__grid {
  display: grid;
  grid-template-rows: repeat(7, 14px);
  gap: 3px;
}

.fitness-heatmap__label {
  font-size: 10px;
  line-height: 14px;
  color: var(--color-text-secondary);
}

.fitness-heatmap__cell {
  width: 14px;
  height: 14px;
  border-radius: 3px;
}

.fitness-heatmap__legend {
  display: flex;
  gap: 14px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.fitness-heatmap__legend i {
  display: inline-block;
  width: 10px;
  height: 10px;
  margin-right: 4px;
  border-radius: 2px;
  vertical-align: middle;
}

.lv-none {
  background: var(--color-border-light);
}

.lv-rest {
  background: var(--color-bg-secondary);
  outline: 1px dashed var(--color-border);
  outline-offset: -1px;
}

.lv-rest.is-checked {
  outline-color: var(--color-success);
}

.lv-1 {
  background: color-mix(in srgb, var(--color-success) 35%, transparent);
}

.lv-2 {
  background: color-mix(in srgb, var(--color-success) 65%, transparent);
}

.lv-3 {
  background: var(--color-success);
}
</style>
