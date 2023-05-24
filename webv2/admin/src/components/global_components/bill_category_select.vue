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

// 分类选择器
import {defineEmits, defineProps, ref, withDefaults} from "vue";
import lbbill from "@/plugin/api/lbbill";
import {DefaultOption, DefaultOrderBy} from "@/plugin/api/model/lb";
import type {SelectProps} from 'ant-design-vue';
import {message} from "ant-design-vue";
import type {ModelBillCategory} from "@/plugin/api/model/lbbill";

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


const options = ref<SelectProps['options']>([])
const getBillCategoryList = async () => {
  options.value = []
  try {
    const resp = await lbbill.getBillCategoryList({
      options: {
        opt_list: [
          {
            key: DefaultOption.DefaultOptionOrderBy,
            value: DefaultOrderBy.DefaultOrderByCreatedAtDesc.toString(),
          }
        ],
        size: 2000,
        page: 1,
        skip_total: true,
      },
    });
    // 列表赋值
    resp.list.forEach((e: ModelBillCategory) => {
      options.value!.push({
        value: e.id,
        label: e.name,
      })
    })
  } catch (error: any) {
    message.error(error);
  }
}
getBillCategoryList()

const onChange = (value: number | undefined) => {
  emit('update:val', value);
  emit('change', value);
}

</script>

<style scoped>

</style>