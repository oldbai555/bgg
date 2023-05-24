<template>
  <div>
    <!--  弹窗  -->
    <a-modal :visible="visible" title="新增" @ok="handleOk" , @cancel="handleCancel">
      <!--   加载中  -->
      <a-spin tip="loading..." :spinning="loading">
        <!--  表单  -->
        <a-form ref="formRef" :model="formState" layout="vertical">

          <!--    表单字段    -->
          <a-form-item name="username" label="账号">
            <a-input v-model:value="formState.username" :allowClear="true" :showCount="true"/>
          </a-form-item>


          <!--    表单字段    -->
          <a-form-item name="password" label="密码">
            <a-input v-model:value="formState.password" :allowClear="true" :showCount="true" type="password"/>
          </a-form-item>

          <!--    表单字段    -->
          <a-form-item name="nickname" label="昵称">
            <a-input v-model:value="formState.nickname" :allowClear="true" :showCount="true"/>
          </a-form-item>

          <!--    表单字段    -->
          <a-form-item name="desc" label="描述">
            <a-input v-model:value="formState.desc" :allowClear="true" :showCount="true"/>
          </a-form-item>

          <!--    表单字段    -->
          <a-form-item name="role" label="角色">
            <root-select
                v-model:val="formState.role"
            />
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
import type {ModelUser} from '@/plugin/api/model/lbuser';
import lbuser from "../../plugin/api/lbuser"
import RootSelect from '../global_components/role_select.vue'

// 指向表单
const formRef = ref<FormInstance>();

// 重置表单
const resetForm = () => {
  formRef.value!.resetFields();
  console.log("reset form ok")
};

// 表单字段
const formState = reactive<ModelUser>({

  id: undefined,

  created_at: undefined,

  updated_at: undefined,

  deleted_at: undefined,

  username: undefined,

  password: undefined,

  avatar: undefined,

  nickname: undefined,

  email: undefined,

  github: undefined,

  desc: undefined,

  role: undefined,

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
    await lbuser.addUser({
      user: formState
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