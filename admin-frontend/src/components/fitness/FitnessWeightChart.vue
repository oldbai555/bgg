<template>
  <div class="fitness-chart">
    <div v-if="records.length === 0" class="fitness-chart__empty">{{ emptyText }}</div>
    <template v-else>
      <svg :viewBox="`0 0 ${W} ${H}`" class="fitness-chart__svg" role="img">
        <line
          v-for="g in grid"
          :key="g.y"
          :x1="PAD_L"
          :x2="W - PAD_R"
          :y1="g.y"
          :y2="g.y"
          class="fitness-chart__grid"
        />
        <text
          v-for="g in grid"
          :key="`t${g.y}`"
          :x="PAD_L - 6"
          :y="g.y + 4"
          class="fitness-chart__axis"
          text-anchor="end"
        >
          {{ g.label }}
        </text>
        <line
          v-if="targetY !== null"
          :x1="PAD_L"
          :x2="W - PAD_R"
          :y1="targetY"
          :y2="targetY"
          class="fitness-chart__target"
        />
        <g :transform="`translate(${PAD_L},${PAD_T})`">
          <path :d="weightPath" class="fitness-chart__weight" />
          <path :d="avgPath" class="fitness-chart__avg" />
          <circle
            v-for="p in points"
            :key="p.x"
            :cx="p.x"
            :cy="p.y"
            r="2.5"
            class="fitness-chart__dot"
          />
        </g>
        <text :x="PAD_L" :y="H - 4" class="fitness-chart__axis">{{ records[0].date.slice(5) }}</text>
        <text
          :x="W - PAD_R"
          :y="H - 4"
          class="fitness-chart__axis"
          text-anchor="end"
        >
          {{ records[records.length - 1].date.slice(5) }}
        </text>
      </svg>
      <div class="fitness-chart__legend">
        <span><i class="is-weight"></i>{{ weightLabel }}</span>
        <span><i class="is-avg"></i>{{ avgLabel }}</span>
        <span v-if="targetY !== null"><i class="is-target"></i>{{ targetLabel }} {{ target }}</span>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import type {FitnessBodyRecordItem} from '@/api/generated/admin'
import {buildLinePath, valueRange} from '@/utils/fitness'

const props = withDefaults(
  defineProps<{
    /** 按日期升序，缺测的日期可以不出现 */
    records: FitnessBodyRecordItem[];
    target?: number;
    emptyText?: string;
    weightLabel?: string;
    avgLabel?: string;
    targetLabel?: string;
  }>(),
  {target: 0, emptyText: '暂无数据', weightLabel: '体重', avgLabel: '7 日均值', targetLabel: '目标'}
)

const W = 640
const H = 220
const PAD_L = 40
const PAD_R = 12
const PAD_T = 12
const PAD_B = 24
const innerW = W - PAD_L - PAD_R
const innerH = H - PAD_T - PAD_B

const range = computed(() => {
  const values = props.records.flatMap((r) => [r.weightKg, r.weightAvg7])
  if (props.target > 0) {
    values.push(props.target)
  }
  return valueRange(values)
})

const toY = (v: number) => PAD_T + innerH - ((v - range.value[0]) / (range.value[1] - range.value[0] || 1)) * innerH

const grid = computed(() => {
  const [min, max] = range.value
  const steps = 4
  return Array.from({length: steps + 1}, (_, i) => {
    const v = min + ((max - min) * i) / steps
    return {y: +toY(v).toFixed(1), label: v.toFixed(1)}
  })
})

const weightPath = computed(() =>
  buildLinePath(props.records.map((r) => r.weightKg), innerW, innerH, range.value[0], range.value[1])
)
const avgPath = computed(() =>
  buildLinePath(props.records.map((r) => r.weightAvg7), innerW, innerH, range.value[0], range.value[1])
)

const points = computed(() => {
  const n = props.records.length
  const step = n > 1 ? innerW / (n - 1) : 0
  return props.records
    .map((r, i) => ({v: r.weightKg, x: +(i * step).toFixed(1)}))
    .filter((p) => p.v > 0)
    .map((p) => ({x: p.x, y: +(toY(p.v) - PAD_T).toFixed(1)}))
})

const targetY = computed(() => (props.target > 0 ? +toY(props.target).toFixed(1) : null))
</script>

<style scoped lang="scss">
.fitness-chart {
  width: 100%;

  &__empty {
    padding: 32px 0;
    text-align: center;
    color: var(--color-text-secondary);
    font-size: 13px;
  }

  &__svg {
    display: block;
    width: 100%;
    max-width: 720px;
    height: auto;
  }

  &__grid {
    stroke: var(--color-border-light);
    stroke-width: 1;
  }

  &__target {
    stroke: var(--color-success);
    stroke-dasharray: 4 4;
  }

  &__axis {
    font-size: 11px;
    fill: var(--color-text-secondary);
  }

  &__weight {
    fill: none;
    stroke: var(--color-primary);
    stroke-width: 1.5;
    opacity: 0.55;
  }

  &__avg {
    fill: none;
    stroke: var(--color-warning);
    stroke-width: 2.5;
  }

  &__dot {
    fill: var(--color-primary);
  }

  &__legend {
    display: flex;
    flex-wrap: wrap;
    gap: 14px;
    margin-top: 6px;
    font-size: 12px;
    color: var(--color-text-secondary);

    i {
      display: inline-block;
      width: 14px;
      height: 3px;
      margin-right: 4px;
      vertical-align: middle;
      border-radius: 2px;

      &.is-weight {
        background: var(--color-primary);
      }

      &.is-avg {
        background: var(--color-warning);
      }

      &.is-target {
        background: var(--color-success);
      }
    }
  }
}
</style>
