<template>
  <ul class="fitness-meals">
    <li v-for="meal in meals" :key="meal.slot" class="fitness-meals__item">
      <div class="fitness-meals__head">
        <span class="fitness-meals__slot">{{ dictStore.getDictLabel('fitness_meal_slot', meal.slot) }}</span>
        <span v-if="meal.time" class="fitness-meals__time">{{ meal.time }}</span>
        <span v-if="meal.place" class="fitness-meals__place">{{ meal.place }}</span>
      </div>
      <div class="fitness-meals__food">
        {{ meal.food }}
        <span v-if="meal.portion" class="fitness-meals__portion">{{ meal.portion }}</span>
      </div>
      <div v-if="meal.kcal || meal.protein" class="fitness-meals__nutri">
        <el-tag v-if="meal.kcal" size="small" effect="plain">{{ meal.kcal }} kcal</el-tag>
        <el-tag
          v-if="meal.protein"
          size="small"
          type="success"
          effect="plain"
        >蛋白质 {{ meal.protein }}g</el-tag>
      </div>
      <p v-if="meal.howToOrder" class="fitness-meals__line">
        <b>{{ t('fitness.howToOrder') }}</b>{{ meal.howToOrder }}
      </p>
      <p v-if="meal.alternatives?.length" class="fitness-meals__line">
        <b>{{ t('fitness.alternatives') }}</b>{{ meal.alternatives.join('；') }}
      </p>
      <slot name="extra" :meal="meal"></slot>
    </li>
  </ul>
</template>

<script setup lang="ts">
import {useI18n} from 'vue-i18n'
import type {FitnessMeal} from '@/api/generated/admin'
import {useDictStore} from '@/stores/dict'

defineProps<{meals: FitnessMeal[]}>()

const {t} = useI18n()
const dictStore = useDictStore()
</script>

<style scoped lang="scss">
.fitness-meals {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 12px;

  &__item {
    padding: 12px 14px;
    border-radius: 10px;
    border: 1px solid var(--color-border-light);
    background: var(--color-bg-card);
    display: flex;
    flex-direction: column;
    gap: 6px;
    font-size: 14px;
    color: var(--color-text-regular);
  }

  &__head {
    display: flex;
    align-items: baseline;
    gap: 8px;
  }

  &__slot {
    font-weight: 700;
    color: var(--color-text-primary);
  }

  &__time {
    color: var(--color-primary);
    font-weight: 600;
  }

  &__place {
    font-size: 12px;
    color: var(--color-text-secondary);
  }

  &__food {
    color: var(--color-text-primary);
    line-height: 1.6;
  }

  &__portion {
    margin-left: 6px;
    font-size: 12.5px;
    color: var(--color-text-secondary);
  }

  &__nutri {
    display: flex;
    gap: 6px;
  }

  &__line {
    margin: 0;
    font-size: 13px;
    line-height: 1.6;

    b {
      margin-right: 6px;
      color: var(--color-text-primary);
    }
  }
}
</style>
