<template>
  <div>
    <span style="margin-right:24px">
      <a-badge :count="0">
        <a-avatar shape="square" size="large" :src="avatarValue"/>
      </a-badge>
    </span>
    <span>
      <a-button type="danger" @click="loginOut">退出</a-button>
    </span>
  </div>
</template>

<script>
export default {
  data() {
    return {
      avatarValue: "",
    };
  },
  methods: {
    async loginOut() {
      const {data: res} = await this.$http.get('/user/Logout')
      if (res.code !== 200) {
        return this.$message.error(res.message)
      }

      this.$message.info("退出登录");
      window.sessionStorage.clear("sid");
      await this.$router.push('/login');
    },
    async InitUser() {
      const {data: res} = await this.$http.get('/user/GetLoginUser')
      if (res.code !== 200) {
        this.$message.error(res.message)
        if (res.code === 401) {
          window.sessionStorage.clear()
          await this.$router.push('/login')
        }
        return
      }
      this.avatarValue = res.data.avatar;
    }
  },
  created() {
    this.InitUser()
  }
}
</script>
