<template>
  <div>
    <a-select
        allowClear
        :placeholder="placeholder"
        :value="val"
        @change="onChange"
        :options="options"
        :disabled="disabled"
    />
  </div>
</template>

<script setup lang="ts">
import type {SelectProps} from 'ant-design-vue';
import {ModelUser_Role} from "@/plugin/api/model/lbuser";
import {defineEmits, defineProps, ref, withDefaults} from "vue";


interface Props {
  val?: number,
  placeholder?: string,
  disabled?: boolean,
}

const props = withDefaults(defineProps<Props>(), {
  val: undefined,
  placeholder: "请选择角色",
  disabled: false,
})

const emit = defineEmits(['change', 'update:val'])

const options = ref<SelectProps['options']>([
  {
    value: ModelUser_Role.RoleAdmin,
    label: "管理员",
  }
])

const onChange = (value: number | undefined) => {
  emit('update:val', value);
  emit('change', value);
}

</script>

<style scoped>

</style>