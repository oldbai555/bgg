<template>
  <div>
    <!--  弹窗  -->
    <a-modal :visible="visible" title="详情" @ok="handleOk" , @cancel="handleCancel">
      <!--   加载中  -->
      <a-spin tip="loading..." :spinning="loading">
        <!--  表单  -->
        <a-form ref="formRef" :model="formState" layout="vertical" class="list">

          <!--    表单字段    -->
          <a-form-item name="ID" label="id">
            <a-input v-model:value="formState.id" :allowClear="true" :showCount="true" disabled/>
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="创建时间" label="created_at">
            <a-input v-model:value="formState.created_at" :allowClear="true" :showCount="true" disabled/>
          </a-form-item>

          <!--    表单字段    -->
          <a-form-item name="文章作者" label="article_id">
            <a-input v-model:value="formState.article_id" :allowClear="true" :showCount="true" disabled/>
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="评论用户" label="user_id">
            <a-input v-model:value="formState.user_id" :allowClear="true" :showCount="true" disabled/>
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="用户邮箱" label="user_email">
            <a-input v-model:value="formState.user_email" :allowClear="true" :showCount="true" disabled/>
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="内容" label="content">
            <a-input v-model:value="formState.content" :allowClear="true" :showCount="true" disabled/>
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="状态" label="status">
            <a-tag v-if="formState.status === 1">正常</a-tag>
            <a-tag v-else-if="formState.status ===2">审核中</a-tag>
            <a-tag v-else-if="formState.status ===3">撤下</a-tag>
            <a-tag v-else>待审核</a-tag>
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
import type {ModelComment} from '@/plugin/api/model/lbblog';
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

// 展示弹窗
const show = async (id: number | undefined) => {
  loading.value = true;
  visible.value = true;
  if (id) {
    try {
      const resp = await lbblog.getComment({id: id});
      formState.value = resp.comment;
    } catch (error: any) {
      message.error(error);
    }
  }
  loading.value = false;
};

// 确认
const handleOk = async (e: MouseEvent) => {
  visible.value = false;
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