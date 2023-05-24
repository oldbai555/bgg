<template>
  <div>
    <a-row justify="end">
      <a-col style="margin-right: 10px">
        <a-dropdown>
          <a-badge :count="count">
            <a-avatar shape="square" size="large" :src="avatar"
                      :size="{ xs: 24, sm: 32, md: 40, lg: 64, xl: 80, xxl: 100 }"/>
          </a-badge>
          <template #overlay>
            <a-menu>
              <a-menu-item>
                <a>个人信息</a>
              </a-menu-item>
              <a-menu-item>
                <a @click="logout">退出登陆</a>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </a-col>
    </a-row>

  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from "vue";
import {useRouter} from "vue-router";
import {message} from "ant-design-vue"
import userApi from "../../plugin/api/lbuser";
import {clearAllCaches} from "../../plugin/utils/cache";


export default defineComponent({
  setup() {
    const router = useRouter()
    const avatar = ref("");
    const count = ref(0);

    // 获取登陆的信息
    const getLoginInfo = async () => {
      try {
        const resp = await userApi.getLoginUser({})
        avatar.value = resp.avatar
      } catch (error: any) {
        message.error(error)
        // 拿不到数据，清理一下缓存，需要重新登陆
        clearAllCaches()
        await router.push({
          name: "login",
        })
      }
    }
    getLoginInfo()

    const logout = async () => {
      try {
        const _ = await userApi.logout({});
      } catch (error: any) {
        message.error(error)
      }

      // 清理缓存
      clearAllCaches()

      // 到登陆页面
      await router.push({
        name: "login",
      })
    }

    return {
      avatar,
      count,

      logout,
    }
  },
})


</script>
