<template>
  <el-drawer
    v-model="visible"
    :title="`周计划：${templateName}`"
    size="80%"
    destroy-on-close
  >
    <div v-loading="loading" class="fitness-week-editor">
      <el-tabs v-model="activeTab">
        <el-tab-pane
          v-for="w in 7"
          :key="w"
          :label="weekdayLabel(w)"
          :name="`w${w}`"
        >
          <template v-if="drafts[w]">
            <FitnessDayEditor v-model="drafts[w]" />
            <div class="fitness-week-editor__actions">
              <el-button
                v-permission="'fitness_template:update'"
                type="primary"
                :loading="saving"
                @click="saveWeekday(w)"
              >
                保存{{ weekdayLabel(w) }}
              </el-button>
            </div>
          </template>
        </el-tab-pane>

        <el-tab-pane label="按日期调整" name="override">
          <el-alert
            type="info"
            :closable="false"
            show-icon
            title="某一天的计划 = 按日期调整（如有）> 模板对应星期几。适合节假日、出差、临时加练等单日改动。"
          />
          <div class="fitness-week-editor__toolbar">
            <el-date-picker
              v-model="newDate"
              type="date"
              value-format="YYYY-MM-DD"
              placeholder="选择日期"
            />
            <el-button v-permission="'fitness_template:update'" :disabled="!newDate" @click="startOverride(newDate)">
              新增调整（以当天星期几的内容为底稿）
            </el-button>
          </div>
          <el-table :data="overrides" size="small" border>
            <el-table-column prop="planDate" label="日期" width="120" />
            <el-table-column label="星期" width="80">
              <template #default="{row}">{{ row.planDate ? weekdayLabel(weekdayOf(row.planDate)) : '' }}</template>
            </el-table-column>
            <el-table-column label="训练" width="90">
              <template #default="{row}">{{ getLabel(row.content?.trainingType) }}</template>
            </el-table-column>
            <el-table-column label="动作">
              <template #default="{row}">{{ (row.content?.exercises || []).map((e: FitnessExercise) => e.name).join('、') || '-' }}</template>
            </el-table-column>
            <el-table-column :label="t('common.actions')" width="140">
              <template #default="{row}">
                <el-button
                  v-permission="'fitness_template:update'"
                  link
                  type="primary"
                  @click="startOverride(row.planDate)"
                >
                  {{ t('common.edit') }}
                </el-button>
                <el-button
                  v-permission="'fitness_template:update'"
                  link
                  type="danger"
                  @click="removeOverride(row)"
                >
                  {{ t('common.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-card v-if="overrideDraft" shadow="never" class="fitness-week-editor__override">
            <template #header>编辑 {{ overrideDate }}（{{ weekdayLabel(weekdayOf(overrideDate)) }}）</template>
            <FitnessDayEditor v-model="overrideDraft" />
            <div class="fitness-week-editor__actions">
              <el-button @click="overrideDraft = null">{{ t('common.cancel') }}</el-button>
              <el-button
                v-permission="'fitness_template:update'"
                type="primary"
                :loading="saving"
                @click="saveOverride"
              >
                {{ t('common.save') }}
              </el-button>
            </div>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import {ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import FitnessDayEditor from '@/components/fitness/FitnessDayEditor.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessExercise, FitnessTemplateDayItem} from '@/api/generated/admin'
import {useDictOptions} from '@/composables/useDictOptions'
import {fromDayDraft, toDayDraft, weekdayLabel, weekdayOf, type FitnessDayDraft} from '@/utils/fitness'

const props = defineProps<{templateId: number; templateName: string}>()
const visible = defineModel<boolean>({required: true})

const {t} = useI18n()
const {getLabel} = useDictOptions('fitness_training_type')

const loading = ref(false)
const saving = ref(false)
const activeTab = ref('w1')
const drafts = ref<Record<number, FitnessDayDraft>>({})
const weekly = ref<FitnessTemplateDayItem[]>([])
const overrides = ref<FitnessTemplateDayItem[]>([])
const newDate = ref('')
const overrideDate = ref('')
const overrideDraft = ref<FitnessDayDraft | null>(null)

const load = async () => {
  loading.value = true
  try {
    const resp = await fitnessApi.templateDays({templateId: props.templateId})
    weekly.value = resp.weekly || []
    overrides.value = resp.overrides || []
    const next: Record<number, FitnessDayDraft> = {}
    for (let w = 1; w <= 7; w++) {
      next[w] = toDayDraft(weekly.value.find((d) => d.weekday === w)?.content)
    }
    drafts.value = next
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '加载失败')
  } finally {
    loading.value = false
  }
}

watch(
  () => [visible.value, props.templateId] as const,
  ([open]) => {
    if (open) {
      activeTab.value = 'w1'
      overrideDraft.value = null
      load()
    }
  },
  {immediate: true}
)

const save = async (weekday: number, planDate: string, draft: FitnessDayDraft) => {
  saving.value = true
  try {
    await fitnessApi.templateDaySave({templateId: props.templateId, weekday, planDate, content: fromDayDraft(draft)})
    ElMessage.success(t('common.updateSuccess'))
    return true
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
    return false
  } finally {
    saving.value = false
  }
}

const saveWeekday = async (weekday: number) => {
  if (await save(weekday, '', drafts.value[weekday])) {
    await load()
  }
}

const startOverride = (date: string) => {
  const existing = overrides.value.find((o) => o.planDate === date)
  overrideDate.value = date
  overrideDraft.value = toDayDraft(existing ? existing.content : weekly.value.find((d) => d.weekday === weekdayOf(date))?.content)
}

const saveOverride = async () => {
  if (!overrideDraft.value) {
    return
  }
  if (await save(0, overrideDate.value, overrideDraft.value)) {
    overrideDraft.value = null
    newDate.value = ''
    await load()
  }
}

const removeOverride = async (row: FitnessTemplateDayItem) => {
  try {
    await ElMessageBox.confirm(`删除 ${row.planDate} 的调整后，这一天回到模板${weekdayLabel(weekdayOf(row.planDate))}的内容`, t('common.confirm'), {
      type: 'warning'
    })
  } catch {
    return
  }
  try {
    await fitnessApi.templateDayDelete({id: row.id})
    ElMessage.success(t('common.deleteSuccess'))
    await load()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  }
}
</script>

<style scoped>
.fitness-week-editor__actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}

.fitness-week-editor__toolbar {
  display: flex;
  gap: 8px;
  margin: 12px 0;
}

.fitness-week-editor__override {
  margin-top: 16px;
}
</style>
