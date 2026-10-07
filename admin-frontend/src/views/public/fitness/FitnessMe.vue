<template>
  <div class="fitness-me-page public-detail-page">
    <MetricReporter module="fitness_me" />
    <PublicHeader />
    <FitnessTabBar />
    <div class="page-shell">
      <div v-loading="loading" class="page-layout">
        <section v-if="profile" class="detail-card fitness-me-head">
          <el-avatar :size="52" :src="profile.avatar">{{ profile.nickname.slice(0, 1) }}</el-avatar>
          <div class="fitness-me-head__name">{{ profile.nickname }}</div>
          <div class="fitness-me-stats">
            <div class="fitness-me-stats__item">
              <b>{{ profile.streakDays }}</b>
              <span>{{ t('fitness.streak') }}（{{ t('fitness.days') }}）</span>
            </div>
            <div class="fitness-me-stats__item">
              <b>{{ formatRate(profile.weekRate) }}</b>
              <span>{{ t('fitness.weekRate') }}</span>
            </div>
            <div class="fitness-me-stats__item">
              <b>{{ profile.latestWeightKg || '-' }}</b>
              <span>{{ t('fitness.weight') }}</span>
            </div>
          </div>
        </section>

        <FitnessSuggestionCard
          v-if="profile?.suggestion && profile.suggestion.status === FitnessSuggestionStatus.Pending"
          :suggestion="profile.suggestion"
          @decided="load"
        />

        <section v-if="profile" class="detail-card">
          <h2 class="title">{{ t('fitness.currentTemplate') }}</h2>
          <div v-if="profile.currentTemplate" class="fitness-tpl is-current">
            <div class="fitness-tpl__name">
              {{ profile.currentTemplate.name }}
              <el-tag size="small" effect="plain">{{ profile.currentTemplate.stage }}</el-tag>
            </div>
            <p class="fitness-tpl__summary">{{ profile.currentTemplate.summary }}</p>
          </div>
          <el-alert
            v-if="profile.pendingSwitch"
            type="warning"
            :closable="false"
            show-icon
            :title="t('fitness.pendingSwitch', {date: profile.pendingSwitch.effectiveDate, name: profile.pendingSwitch.templateName})"
            class="fitness-me__alert"
          />
          <h3 class="fitness-me__subtitle">{{ t('fitness.switchTemplate') }}</h3>
          <p class="fitness-me__hint">{{ t('fitness.switchHint') }}</p>
          <div class="fitness-tpl-list">
            <div
              v-for="tpl in switchable"
              :key="tpl.id"
              class="fitness-tpl"
              :class="{'is-pending': profile.pendingSwitch?.templateId === tpl.id}"
            >
              <div class="fitness-tpl__name">
                {{ tpl.name }}
                <el-tag size="small" effect="plain">{{ tpl.stage }}</el-tag>
              </div>
              <p class="fitness-tpl__summary">{{ tpl.summary }}</p>
              <el-button
                size="small"
                type="primary"
                plain
                :disabled="profile.pendingSwitch?.templateId === tpl.id"
                @click="switchTo(tpl)"
              >
                {{ t('fitness.switchTemplate') }}
              </el-button>
            </div>
          </div>
        </section>

        <section class="detail-card">
          <el-form label-position="top" class="fitness-profile-form" @submit.prevent>
            <el-form-item :label="t('fitness.height')">
              <el-input-number
                v-model="form.heightCm"
                :min="0"
                :max="250"
                :step="1"
                :precision="1"
              />
            </el-form-item>
            <el-form-item :label="t('fitness.targetWeight')">
              <el-input-number
                v-model="form.targetWeightKg"
                :min="0"
                :max="300"
                :step="0.5"
                :precision="1"
              />
            </el-form-item>
            <el-form-item :label="t('fitness.preference')">
              <el-radio-group v-model="form.trainingPreference">
                <el-radio-button
                  v-for="opt in options('fitness_training_preference')"
                  :key="opt.value"
                  :value="Number(opt.value)"
                >
                  {{ opt.label }}
                </el-radio-button>
              </el-radio-group>
            </el-form-item>
          </el-form>
          <el-button
            type="primary"
            round
            :loading="saving"
            class="fitness-me__save"
            @click="saveProfile"
          >
            {{ t('common.save') }}
          </el-button>
        </section>

        <el-button class="fitness-me__logout" @click="logout">{{ t('fitness.logout') }}</el-button>
      </div>
    </div>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import {useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import PublicHeader from '@/components/common/PublicHeader.vue'
import MetricReporter from '@/components/common/MetricReporter.vue'
import IcpFooter from '@/components/common/IcpFooter.vue'
import FitnessTabBar from '@/components/fitness/FitnessTabBar.vue'
import FitnessSuggestionCard from '@/components/fitness/FitnessSuggestionCard.vue'
import {fitnessApi} from '@/api/fitness'
import type {FitnessMyProfileResp, FitnessTemplateBrief} from '@/api/generated/admin'
import {useUserStore} from '@/stores/user'
import {useFitnessDicts} from '@/composables/useFitnessDicts'
import {FITNESS_LOGIN_PATH, FitnessSuggestionStatus, FitnessTrainingPreference} from '@/constants/fitness'
import {formatRate} from '@/utils/fitness'

const router = useRouter()
const {t} = useI18n()
const userStore = useUserStore()
const {options} = useFitnessDicts()

const profile = ref<FitnessMyProfileResp | null>(null)
const loading = ref(false)
const saving = ref(false)
const form = reactive<{heightCm: number; targetWeightKg: number; trainingPreference: number}>({
  heightCm: 0,
  targetWeightKg: 0,
  trainingPreference: FitnessTrainingPreference.Morning
})

const switchable = computed(() => (profile.value?.templates || []).filter((tpl) => tpl.id !== profile.value?.currentTemplate?.id))

const load = async () => {
  loading.value = true
  try {
    const resp = await fitnessApi.myProfile()
    profile.value = resp
    form.heightCm = resp.heightCm
    form.targetWeightKg = resp.targetWeightKg
    form.trainingPreference = resp.trainingPreference || FitnessTrainingPreference.Morning
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('fitness.loadFailed'))
  } finally {
    loading.value = false
  }
}

const saveProfile = async () => {
  saving.value = true
  try {
    await fitnessApi.myProfileUpdate({...form})
    ElMessage.success(t('fitness.saved'))
    await load()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  } finally {
    saving.value = false
  }
}

const switchTo = async (tpl: FitnessTemplateBrief) => {
  try {
    await ElMessageBox.confirm(`${t('fitness.switchHint')}`, `${t('fitness.switchTemplate')}：${tpl.name}`, {
      type: 'warning'
    })
  } catch {
    return
  }
  try {
    const resp = await fitnessApi.myTemplateSwitch({templateId: tpl.id})
    ElMessage.success(t('fitness.pendingSwitch', {date: resp.effectiveDate, name: tpl.name}))
    await load()
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : t('common.submitFailed'))
  }
}

const logout = async () => {
  await userStore.logout()
  router.replace(FITNESS_LOGIN_PATH)
}

onMounted(load)
</script>

<style scoped lang="scss">
@import '@/styles/public-detail.scss';

.fitness-me-page {
  .page-shell {
    max-width: 760px;
  }

  .page-layout {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .detail-card {
    padding: 20px 24px;

    .title {
      font-size: 17px;
      margin-bottom: 12px;
    }
  }
}

.fitness-me-head {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;

  &__name {
    font-size: 17px;
    font-weight: 700;
    color: var(--color-text-primary);
  }
}

.fitness-me-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  width: 100%;
  margin-top: 8px;

  &__item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;

    b {
      font-size: 22px;
      color: var(--color-primary);
    }

    span {
      font-size: 12px;
      color: var(--color-text-secondary);
    }
  }
}

.fitness-me__alert {
  margin-top: 12px;
}

.fitness-me__subtitle {
  margin: 18px 0 4px;
  font-size: 15px;
  color: var(--color-text-primary);
}

.fitness-me__hint {
  margin: 0 0 12px;
  font-size: 12.5px;
  color: var(--color-text-secondary);
}

.fitness-tpl-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.fitness-tpl {
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid var(--color-border-light);
  background: var(--color-bg-secondary);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;

  &.is-current {
    border-color: var(--color-primary);
    background: var(--color-bg-card);
  }

  &.is-pending {
    border-color: var(--color-warning);
  }

  &__name {
    display: flex;
    align-items: center;
    gap: 6px;
    font-weight: 600;
    color: var(--color-text-primary);
  }

  &__summary {
    margin: 0;
    font-size: 12.5px;
    line-height: 1.6;
    color: var(--color-text-regular);
  }
}

.fitness-profile-form {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0 16px;
}

.fitness-me__save,
.fitness-me__logout {
  width: 100%;
}

@include mobile {
  .fitness-me-page {
    padding-bottom: calc(60px + env(safe-area-inset-bottom));

    .detail-card {
      padding: 16px;
    }
  }

  .fitness-tpl-list,
  .fitness-profile-form {
    grid-template-columns: 1fr;
  }
}
</style>
