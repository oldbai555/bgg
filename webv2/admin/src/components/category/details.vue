<template>
  <div>
    <!--  弹窗  -->
    <a-modal :visible="visible" title="分类详情" @ok="handleOk" , @cancel="handleCancel">
      <!--   加载中  -->
      <a-spin tip="loading..." :spinning="loading">
        <!--  表单  -->
        <a-form ref="formRef" :model="formState" layout="vertical">
          <!--    表单字段    -->
          <a-form-item name="name" label="分类名称">
            <a-input v-model:value="formState.name" disabled/>
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
import type {ModelCategory} from '../../plugin/api/model/lbblog'
import lbblog from "../../plugin/api/lbblog";

// 指向表单
const formRef = ref<FormInstance>();

// 重置表单
const resetForm = () => {
  formRef.value!.resetFields();
  console.log("reset form ok")
};

// 表单字段
const formState = ref<ModelCategory | undefined>({
  id: 0,
  created_at: 0,
  updated_at: 0,
  deleted_at: 0,
  name: "",
});

// 弹窗的显示
const visible = ref(false)

// 加载中
const loading = ref(false)

// 接收父组件传递过来的方法
const emit = defineEmits(['handleComplete'])

// 通知父组件
const notifyParent = () => {
  emit('handleComplete')
  resetForm()
}

// 展示弹窗
const show = async (id: number | undefined) => {
  loading.value = true
  visible.value = true
  if (id) {
    try {
      const resp = await lbblog.getCategory({id: id})
      formState.value = resp.category
    } catch (error: any) {
      message.error(error);
    }
  }
  loading.value = false
}

// 确认
const handleOk = async (e: MouseEvent) => {
  visible.value = false
  notifyParent()
};

// 取消
const handleCancel = (e: MouseEvent) => {
  visible.value = false
  notifyParent()
}

// 导出方法给父组件调用
defineExpose({
  show,
})
</script>

