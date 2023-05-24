<template>
  <div>
    <!--  弹窗  -->
    <a-modal :visible="visible" title="新增" @ok="handleOk" , @cancel="handleCancel">
      <!--   加载中  -->
      <a-spin tip="loading..." :spinning="loading">
        <!--  表单  -->
        <a-form ref="formRef" :model="formState" layout="vertical" class="form_overflow">

          <!--    表单字段    -->
          <a-form-item name="category_id" label="分类">
            <bill-category-select
                style="width: 130px"
                v-model:val="formState.category_id"
                placeholder="请选择分类"
            />
          </a-form-item>

          <!--    表单字段    -->
          <a-form-item name="amount" label="金额">
            <input-number
                style="width: 130px"
                v-model:val="formState.amount"
                placeholder="请输入金额"
            />
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="root_category" label="收入/支出">
            <bill-root-category-select
                style="width: 130px"
                v-model:val="formState.root_category"
                placeholder="请选择收入/支出"
            />
          </a-form-item>

          <!--    表单字段    -->
          <a-form-item name="remark" label="备注">
            <a-input v-model:value="formState.remark" :allowClear="true" :showCount="true"/>
          </a-form-item>

        </a-form>
      </a-spin>
    </a-modal>
  </div>
</template>
<script setup lang="ts">
import {defineEmits, ref} from 'vue';
import type {FormInstance} from 'ant-design-vue';
import {message} from 'ant-design-vue';
import type {ModelBill} from '@/plugin/api/model/lbbill';
import lbbill from "../../plugin/api/lbbill"
import InputNumber from '../global_components/input_number.vue';
import BillCategorySelect from '../global_components/bill_category_select.vue';
import BillRootCategorySelect from '../global_components/bill_root_category_select.vue';


// 指向表单
const formRef = ref<FormInstance>();

// 重置表单
const resetForm = () => {
  formRef.value!.resetFields();
  console.log("reset form ok")
};

// 表单字段
const formState = ref<ModelBill | undefined>({

  id: undefined,

  created_at: undefined,

  updated_at: undefined,

  deleted_at: undefined,

  creator_uid: undefined,

  amount: undefined,

  category_id: undefined,

  date_unix: undefined,

  root_category: undefined,

  remark: undefined,

});

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
    const resp = await lbbill.addBill({
      bill: formState.value
    });
    console.log("add complete , id is ", resp.bill?.id);
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