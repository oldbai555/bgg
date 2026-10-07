<template>
  <el-table :data="data" size="small" border>
    <el-table-column
      v-if="showUser"
      prop="nickname"
      label="成员"
      width="110"
    />
    <el-table-column label="切换" min-width="180">
      <template #default="{row}">{{ row.fromTemplateName || '（无）' }} → <b>{{ row.toTemplateName }}</b></template>
    </el-table-column>
    <el-table-column prop="effectiveDate" label="生效日期" width="110" />
    <el-table-column label="来源" width="110">
      <template #default="{row}">
        <el-tag size="small" effect="plain">{{ getLabel(row.source) }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column
      prop="reason"
      label="原因"
      min-width="200"
      show-overflow-tooltip
    />
    <el-table-column label="操作人" width="100">
      <template #default="{row}">{{ row.operatorName || '-' }}</template>
    </el-table-column>
    <el-table-column label="记录时间" width="170">
      <template #default="{row}">{{ formatUnixTime(row.createdAt) }}</template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type {FitnessSwitchItem} from '@/api/generated/admin'
import {useDictOptions} from '@/composables/useDictOptions'
import {formatUnixTime} from '@/utils/date'

withDefaults(defineProps<{data: FitnessSwitchItem[]; showUser?: boolean}>(), {showUser: true})

const {getLabel} = useDictOptions('fitness_switch_source')
</script>
