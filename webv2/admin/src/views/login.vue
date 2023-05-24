<template>
  <div class="container">
    <a-row type="flex" justify="center" align="middle" style="min-height:90vh;">
      <a-col :span="8">
        <a-card :bordered="false">
          <a-form
              :model="formState"
              name="basic"
              :label-col="{ span: 6 }"
              :wrapper-col="{ span: 16 }"
              autocomplete="off"
              :rules="rules"
              ref="formRef"
          >
            <!--账号-->
            <a-form-item
                label="用户名"
                name="username"
                has-feedback
                style="padding-top: 20px"
            >
              <a-input v-model:value="formState.username" autocomplete="off" placeholder="请输入用户名">
                <template #prefix>
                  <user-outlined style="color:rgba(0,0,0,.25)"/>
                </template>
                <template #suffix>
                  <a-tooltip title="手机号或邮箱">
                    <info-circle-outlined style="color: rgba(0, 0, 0, 0.45)"/>
                  </a-tooltip>
                </template>
              </a-input>
            </a-form-item>

            <!--密码-->
            <a-form-item
                label="密码"
                name="password"
                has-feedback
            >
              <a-input-password v-model:value="formState.password" autocomplete="off" placeholder="请输入密码">
                <template #prefix>
                  <lock-outlined style="color:rgba(0,0,0,.25)"/>
                </template>
              </a-input-password>
            </a-form-item>

            <!--登陆按钮-->
            <a-form-item>
              <a-row type="flex" justify="space-around" align="middle">
                <a-col :span="4" :offset="8">
                  <a-button type="primary" html-type="submit" @click="login">登录</a-button>
                </a-col>
                <a-col :span="4" :offset="8">
                  <a-button type="info" html-type="submit" @click="resetForm">取消</a-button>
                </a-col>
              </a-row>

            </a-form-item>
          </a-form>
        </a-card>

      </a-col>
    </a-row>

  </div>
</template>

<script lang="ts" setup>
import userApi from '../plugin/api/lbuser';
import type {LoginReq} from '../plugin/api/model/lbuser'
import type {FormInstance} from 'ant-design-vue';
import {message} from 'ant-design-vue';
import type {Rule} from 'ant-design-vue/es/form';
import {reactive, ref} from 'vue';
import {InfoCircleOutlined, LockOutlined, UserOutlined} from '@ant-design/icons-vue';
import {useRouter} from 'vue-router'
import {setToken} from "../plugin/utils/cache";

// 路由声明
const router = useRouter()

// 表单输入内容
const formState = reactive<LoginReq>({
  username: '',
  password: '',
});

// 指向表单的 ref
const formRef = ref<FormInstance>();

// 重置表单
const resetForm = () => {
  formRef.value!.resetFields();
};

// 表单规则
const rules: Record<string, Rule[]> = {
  username: [{required: true, message: '请输入用户名', trigger: 'blur'}, {
    min: 6,
    max: 12,
    message: '用户名必须在6到12个字符之间',
    trigger: 'blur',
  }],
  password: [{required: true, message: '请输入密码', trigger: 'blur'}, {
    min: 6,
    max: 20,
    message: '密码必须在6到20个字符之间',
    trigger: 'blur',
  }],
};

// 提交登录请求
const login = async () => {
  try {
    let resp = await userApi.login(formState)
    message.success("登录成功")
    resetForm()
    setToken(resp.sid)
    // 跳转路由
    await router.push({
      name: "admin",
    })
  } catch (error: any) {
    message.error(error)
    localStorage.clear()
  }
};
</script>


<style scoped>
.container {
  height: 100%;
  background-image: url("src/assets/image/bg_login.jpg");
}
</style>
