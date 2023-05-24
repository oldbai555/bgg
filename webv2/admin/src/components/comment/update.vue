<template>
  <div>
    <!--  弹窗  -->
    <a-modal :visible="visible" title="编辑" @ok="handleOk" , @cancel="handleCancel">
      <!--   加载中  -->
      <a-spin tip="loading..." :spinning="loading">
        <!--  表单  -->
        <a-form ref="formRef" :model="formState" layout="vertical" class="list">

          <!--    表单字段    -->
          <a-form-item name="status" label="status">
            <a-select style="width: 130px" placeholder="请进行审核" v-model:value="formState.status">
              <a-select-option v-for="item in statusList" :key="item" :value="item">
                <span v-if="item === 1">正常</span>
                <span v-else-if="item ===2">审核中</span>
                <span v-else-if="item ===3">撤下</span>
              </a-select-option>
            </a-select>
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
import type {ModelComment} from '@/plugin/api/model/lbblog';
import {ModelComment_Status} from "@/plugin/api/model/lbblog";
import lbblog from "../../plugin/api/lbblog"

// 指向表单
const formRef = ref<FormInstance>();

// 重置表单
const resetForm = () => {
  formRef.value!.resetFields();
  console.log("reset form ok");
};

// 表单字段
const formState = ref<ModelComment | undefined>({

  id: undefined,

  created_at: undefined,

  updated_at: undefined,

  deleted_at: undefined,

  article_id: undefined,

  user_id: undefined,

  user_email: undefined,

  content: undefined,

  status: undefined,

});

// 弹窗的显示
const visible = ref(false);

// 加载中
const loading = ref(false);

// 接收父组件传递过来的方法
const emit = defineEmits(['handleComplete']);

// 通知父组件
const notifyParent = () => {
  emit('handleComplete');
  resetForm();
};

const statusList = [ModelComment_Status.StatusNormal, ModelComment_Status.StatusReview, ModelComment_Status.StatusTakeDown]

// 展示弹窗
const show = async (id: number | undefined) => {
  loading.value = true;
  visible.value = true;
  if (id) {
    try {
      const resp = await lbblog.getComment({id: id});
      formState.value = resp.comment
    } catch (error: any) {
      message.error(error);
    }
  }
  loading.value = false;
};

// 确认
const handleOk = async (e: MouseEvent) => {
  visible.value = false;
  try {
    await lbblog.updateComment({
      comment: formState.value
    });
    console.log("update complete");
  } catch (error: any) {
    message.error(error);
  }
  notifyParent();
};

// 取消
const handleCancel = (e: MouseEvent) => {
  visible.value = false;
  notifyParent();
};

// 导出方法给父组件调用
defineExpose({
  show,
});
</script>
<style>
.list {
  width: 100%;
  height: 450px;
  overflow-y: auto;
}
</style>