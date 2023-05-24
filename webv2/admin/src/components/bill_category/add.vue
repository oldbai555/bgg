<template>
  <div>
    <!--  弹窗  -->
    <a-modal :visible="visible" title="新增" @ok="handleOk" , @cancel="handleCancel">
      <!--   加载中  -->
      <a-spin tip="loading..." :spinning="loading">
        <!--  表单  -->
        <a-form ref="formRef" :model="formState" layout="vertical">

          <!--    表单字段    -->
          <a-form-item name="name" label="分类名称">
            <a-input v-model:value="formState.name" :allowClear="true" :showCount="true"/>
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="root_category" label="收入/支出">
            <a-select style="width: 130px" placeholder="请选择" v-model:value="formState.root_category">
              <a-select-option v-for="item in selectList" :key="item" :value="item">
                <span v-if="item === RootCategory.RootCategoryIncome">收入</span>
                <span v-else-if="item ===RootCategory.RootCategoryExpenditure">支出</span>
              </a-select-option>
            </a-select>
          </a-form-item>


        </a-form>
      </a-spin>
    </a-modal>
  </div>
</template>
<script setup lang="ts">
import {defineEmits, reactive, ref} from 'vue';
import type {FormInstance} from 'ant-design-vue';
import {message} from 'ant-design-vue';
import type {ModelBillCategory} from '@/plugin/api/model/lbbill';
import {RootCategory} from "@/plugin/api/model/lbbill";
import lbbill from "../../plugin/api/lbbill"

// 指向表单
const formRef = ref<FormInstance>();

// 重置表单
const resetForm = () => {
  formRef.value!.resetFields();
  console.log("reset form ok")
};

// 表单字段
const formState = ref<ModelBillCategory | undefined>({

  id: undefined,

  created_at: undefined,

  updated_at: undefined,

  deleted_at: undefined,

  name: undefined,

  root_category: undefined,

});

const selectList = [RootCategory.RootCategoryIncome, RootCategory.RootCategoryExpenditure]

// 弹窗的显示
const visible = ref(false);

// 加载中
const loading = ref(false);

// 接收父组件传递过来的方法
const emit = defineEmits(['handleComplete'])

// 通知父组件
const notifyParent = () => {
  emit('handleComplete');
  resetForm();
};

// 展示弹窗
const show = () => {
  loading.value = true;
  visible.value = true;
  loading.value = false;
};

// 确认
const handleOk = async (e: MouseEvent) => {
  visible.value = false;
  try {
    await lbbill.addBillCategory({
      category: formState.value
    });
  } catch (error: any) {
    message.error(error);
  }
  await notifyParent();
};

// 取消
const handleCancel = async (e: MouseEvent) => {
  visible.value = false;
  await notifyParent();
};

// 导出方法给父组件调用
defineExpose({
  show,
});
</script>