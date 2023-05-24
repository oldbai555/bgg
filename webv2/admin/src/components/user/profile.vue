<template>
  <a-card>
    <a-form :model="formState" labelAlign="left" >
      <a-form-item label="作者名称" name="nickname">
        <a-input style="width: 300px" v-model:value="formState.nickname"/>
      </a-form-item>

      <a-form-item label="个人简介" name="desc">
        <a-input type="textarea" v-model:value="formState.desc"/>
      </a-form-item>

      <a-form-item label="github" name="github">
        <a-input style="width: 300px" v-model:value="formState.github"/>
      </a-form-item>

      <a-form-item label="Email" name="email">
        <a-input style="width: 300px" v-model:value="formState.email"/>
      </a-form-item>

      <a-form-item label="头像" name="avatar">
        <file-upload
            accept="image/jpeg,image/jpg,image/png"
            @change="uploadChange"
        />

        <div v-if="formState.avatar">
          <img :src="formState.avatar" style="width: 120px; height: 100px"/>
        </div>

      </a-form-item>

      <!--    表单字段    -->
      <a-form-item label="角色" name="role">
        <root-select
            v-model:val="formState.role"
            :disabled="true"
        />
      </a-form-item>

      <a-form-item>
        <a-button type="danger" style="margin-right: 15px" @click="updateProfile">更新</a-button>
      </a-form-item>

    </a-form>
  </a-card>
</template>
<script setup lang="ts">
import FileUpload from "../global_components/file_upload.vue"
import {ref} from "vue";
import type {ModelUser} from "@/plugin/api/model/lbuser";
import RootSelect from '../global_components/role_select.vue'
import lbuser from "@/plugin/api/lbuser";
import {message} from 'ant-design-vue';

// 表单字段
const formState = ref<ModelUser | undefined>({

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

const uploadChange = async (data: string) => {
  formState.value!.avatar = data
}

const updateProfile = async () => {
  try {
    await lbuser.updateLoginUserInfo({
      user: formState.value
    })
    await getLoginUser()
  } catch (e: any) {
    message.error(e)
  }
}

const getLoginUser = async () => {
  try {
    const resp = await lbuser.getLoginUser({});
    formState.value!.github = resp.github
    formState.value!.email = resp.email
    formState.value!.nickname = resp.nickname
    formState.value!.avatar = resp.avatar
    formState.value!.desc = resp.desc
  } catch (e: any) {
    message.error(e)
  }
}
getLoginUser()

</script>

<style scoped>
</style>
