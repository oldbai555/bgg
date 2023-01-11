<template>
  <a-card>
    <a-form-model labelAlign="left" :label-col="{ span: 2 }" :wrapper-col="{ span: 12 }">
      <a-form-model-item label="作者名称">
        <a-input style="width: 300px" v-model="profileInfo.nickname"></a-input>
      </a-form-model-item>

      <a-form-model-item label="个人简介">
        <a-input type="textarea" v-model="profileInfo.desc"></a-input>
      </a-form-model-item>

      <a-form-model-item label="github">
        <a-input style="width: 300px" v-model="profileInfo.github"></a-input>
      </a-form-model-item>

      <a-form-model-item label="Email">
        <a-input style="width: 300px" v-model="profileInfo.email"></a-input>
      </a-form-model-item>

      <a-form-model-item label="头像">
        <a-upload
          listType="picture"
          name="file"
          :action="upUrl"
          :headers="headers"
          @change="avatarChange"
        >
          <a-button style="margin-right:10px">
            <a-icon type="upload"/>
            点击上传
          </a-button>

          <template v-if="profileInfo.avatar">
            <img :src="profileInfo.avatar" style="width: 120px; height: 100px"/>
          </template>
        </a-upload>
      </a-form-model-item>

      <a-form-model-item>
        <a-button type="danger" style="margin-right: 15px" @click="updateProfile">更新</a-button>
      </a-form-model-item>
    </a-form-model>
  </a-card>
</template>
<script>
import {Url} from '../../plugin/http'

export default {
  data() {
    return {
      profileInfo: {
        id: 1,
        nickname: '',
        desc: '',
        github: '',
        email: '',
        img: '',
        avatar: '',
      },
      upUrl: Url + 'upload',
      headers: {},
    }
  },
  created() {
    this.getProfileInfo()
    this.headers = {Authorization: `${window.sessionStorage.getItem('sid')}`}
  },
  methods: {
    // 获取个人设置
    async getProfileInfo() {
      const sid = window.sessionStorage.getItem("sid");
      if (sid === null) {
        this.$message.error("登录信息失效")
        window.sessionStorage.clear("sid")
        await this.$router.push('/login')
      }
      const {data: res} = await this.$http.get('/user/GetLoginUser')
      if (res.code !== 200) {
        this.$message.error(res.message)
        if (res.code === 401) {
          window.sessionStorage.clear()
          await this.$router.push('/login')
        }
        return
      }
      this.profileInfo = res.data
      this.profileInfo.id = sid
    },

    // 上传头像
    avatarChange(info) {
      if (info.file.status !== 'uploading') {
      }
      if (info.file.status === 'done') {
        this.$message.success(`图片上传成功`)
        const imgUrl = info.file.response.url
        this.profileInfo.avatar = imgUrl
      } else if (info.file.status === 'error') {
        this.$message.error(`图片上传失败`)
      }
    },

    // 更新
    async updateProfile() {
      const {data: res} = await this.$http.put(`profile/${this.profileInfo.id}`, this.profileInfo)
      if (res.status !== 200) return this.$message.error(res.message)
      this.$message.success(`个人信息更新成功`)
      await this.$router.push('/index')
    },
  },
}
</script>

<style scoped>
</style>
