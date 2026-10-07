<template>
  <el-form label-width="96px" class="fitness-day-editor" @submit.prevent>
    <el-form-item label="训练类型" required>
      <el-radio-group v-model="model.trainingType">
        <el-radio-button v-for="opt in trainingTypeOptions" :key="opt.value" :value="Number(opt.value)">
          {{ opt.label }}
        </el-radio-button>
      </el-radio-group>
    </el-form-item>
    <el-form-item label="训练时间">
      <el-input
        v-model="model.trainingTime"
        placeholder="如 08:00-09:00"
        maxlength="32"
        style="width: 220px"
      />
    </el-form-item>
    <el-form-item label="步数目标">
      <el-input-number
        v-model="model.stepGoal"
        :min="0"
        :max="50000"
        :step="500"
      />
    </el-form-item>
    <el-form-item label="热身">
      <el-input
        v-model="model.warmup"
        type="textarea"
        :rows="2"
        maxlength="500"
      />
    </el-form-item>

    <el-form-item label="动作">
      <div class="fitness-day-editor__rows">
        <div v-for="(ex, i) in model.exercises" :key="i" class="fitness-day-editor__row">
          <el-input v-model="ex.name" placeholder="动作名" style="width: 160px" />
          <el-input-number
            v-model="ex.sets"
            :min="0"
            :max="20"
            controls-position="right"
            style="width: 90px"
          />
          <el-input v-model="ex.reps" placeholder="次数 如 8-12" style="width: 110px" />
          <el-input v-model="ex.duration" placeholder="时长 如 30分钟" style="width: 120px" />
          <el-input v-model="ex.note" placeholder="要点" style="flex: 1; min-width: 160px" />
          <el-button
            :icon="Delete"
            circle
            size="small"
            @click="model.exercises.splice(i, 1)"
          />
        </div>
        <el-button size="small" :icon="Plus" @click="model.exercises.push({name: '', sets: 3, reps: '', duration: '', note: ''})">
          添加动作
        </el-button>
      </div>
    </el-form-item>

    <el-form-item label="拉伸">
      <el-input
        v-model="model.stretch"
        type="textarea"
        :rows="2"
        maxlength="500"
      />
    </el-form-item>
    <el-form-item label="训练备注">
      <el-input
        v-model="model.trainingNote"
        type="textarea"
        :rows="2"
        maxlength="500"
      />
    </el-form-item>

    <el-form-item label="饮食">
      <div class="fitness-day-editor__rows">
        <el-card
          v-for="(meal, i) in model.meals"
          :key="i"
          shadow="never"
          class="fitness-day-editor__meal"
        >
          <div class="fitness-day-editor__row">
            <el-select v-model="meal.slot" placeholder="餐次" style="width: 110px">
              <el-option
                v-for="opt in mealSlotOptions"
                :key="opt.value"
                :label="opt.label"
                :value="String(opt.value)"
                :disabled="usedSlots.has(String(opt.value)) && meal.slot !== opt.value"
              />
            </el-select>
            <el-input v-model="meal.time" placeholder="时间 如 12:30" style="width: 110px" />
            <el-input v-model="meal.place" placeholder="地点 如 外卖/食堂" style="width: 150px" />
            <el-input-number
              v-model="meal.kcal"
              :min="0"
              :max="3000"
              controls-position="right"
              placeholder="kcal"
              style="width: 110px"
            />
            <el-input-number
              v-model="meal.protein"
              :min="0"
              :max="300"
              controls-position="right"
              placeholder="蛋白质g"
              style="width: 110px"
            />
            <el-button
              :icon="Delete"
              circle
              size="small"
              @click="model.meals.splice(i, 1)"
            />
          </div>
          <div class="fitness-day-editor__row">
            <el-input v-model="meal.food" placeholder="吃什么" style="flex: 2; min-width: 200px" />
            <el-input v-model="meal.portion" placeholder="份量" style="flex: 1; min-width: 120px" />
          </div>
          <el-input v-model="meal.howToOrder" placeholder="怎么点" />
          <el-input
            :model-value="(meal.alternatives || []).join('\n')"
            type="textarea"
            :rows="2"
            placeholder="替代选择，一行一个（2~3 个）"
            @update:model-value="(v: string) => (meal.alternatives = v.split('\n').map((s) => s.trim()).filter(Boolean))"
          />
        </el-card>
        <el-button
          size="small"
          :icon="Plus"
          :disabled="model.meals.length >= mealSlotOptions.length"
          @click="addMeal"
        >
          添加餐次
        </el-button>
      </div>
    </el-form-item>

    <el-form-item label="当天提示">
      <el-input
        v-model="model.tip"
        type="textarea"
        :rows="2"
        maxlength="500"
        placeholder="和通用提示一起展示，可为空"
      />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {Delete, Plus} from '@element-plus/icons-vue'
import {useDictOptions} from '@/composables/useDictOptions'
import type {FitnessDayDraft} from '@/utils/fitness'

const model = defineModel<FitnessDayDraft>({required: true})

const {options: trainingTypeOptions} = useDictOptions('fitness_training_type')
const {options: mealSlotOptions} = useDictOptions('fitness_meal_slot')

const usedSlots = computed(() => new Set(model.value.meals.map((m) => m.slot)))

const addMeal = () => {
  const free = mealSlotOptions.value.find((opt) => !usedSlots.value.has(String(opt.value)))
  model.value.meals.push({slot: free ? String(free.value) : '', time: '', place: '', food: '', portion: '', howToOrder: '', alternatives: []})
}
</script>

<style scoped>
.fitness-day-editor__rows {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.fitness-day-editor__row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.fitness-day-editor__meal :deep(.el-card__body) {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
}
</style>
