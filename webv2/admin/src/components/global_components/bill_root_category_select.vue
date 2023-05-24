<template>
  <div>
    <a-select
        allowClear
        placeholder="请选择分类"
        :value="val"
        @change="onChange"
        :options="options"
        :disabled="disabled"
    />
  </div>
</template>

<script setup lang="ts">
import type {SelectProps} from 'ant-design-vue';
import {RootCategory} from "@/plugin/api/model/lbbill";
import {defineEmits, defineProps, ref, withDefaults} from "vue";


interface Props {
  val?: number,
  placeholder: string,
  disabled?: boolean,
}

const props = withDefaults(defineProps<Props>(), {
  val: undefined,
  placeholder: "请选择",
  disabled: false,
})

const emit = defineEmits(['change', 'update:val'])

const options = ref<SelectProps['options']>([
  {
    value: RootCategory.RootCategoryIncome,
    label: "收入",
  },
  {
    value: RootCategory.RootCategoryExpenditure,
    label: "支出",
  }
])

const onChange = (value: number | undefined) => {
  emit('update:val', value);
  emit('change', value);
}

</script>

<style scoped>

</style>