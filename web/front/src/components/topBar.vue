<template>
  <div>
    <v-app-bar mobileBreakpoint="sm" app dark flat color="indigo darken-2">
      <v-app-bar-nav-icon dark class="hidden-md-and-up" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>
        <v-app-bar-nav-icon class="mx-15 hidden-md-and-down">
          <v-avatar size="40" color="grey">
            <img src="https://baifile-1309918034.cos.ap-guangzhou.myqcloud.com/public/link-info/assets/images/20230131-162920.webp?q-sign-algorithm=sha1&q-ak=AKIDMdBgXmJYhQRDFl7jQpHSkkhEW8SWJ45pxWAo97i99jJ2h8LTQSDsGFjEQWo2rcZA&q-sign-time=1677482696;1677486296&q-key-time=1677482696;1677486296&q-header-list=host&q-url-param-list=&q-signature=f2e9596899d7ba53538d53aa8769bf1c200b8048&x-cos-security-token=cmQy9hKYWHE2Aky9IqB24lrnPcSjSxxa1e185b23d8a0b37adf32b2a10b8bf60ftMDG9SnQ3PrH8A865yHjc-fydov9_kpZF7BPkyZeDjRnJvoJ-Wpa65XPyeNPDsDxNxXmofHTRctr4MZ85Bpp737B50JU2SVcKLA1GoYqcuwZWenmQK-ZXaJcodc7Ws7kE767ySQJcIltv4TTc_QG9OZ8XK6d6hkfBi-Jo_NIOQDhL1UhFhNE4PAHLy-EYJpJ" alt />
          </v-avatar>
        </v-app-bar-nav-icon>
      </v-toolbar-title>

      <v-tabs dark center-active centered class="hidden-sm-and-down">
        <v-tab @click="$router.push('/')">首页</v-tab>
        <v-tab
          v-for="item in cateList"
          :key="item.id"
          text
          @click="gotoCate(item.id)"
        >{{ item.name }}</v-tab>
      </v-tabs>

      <v-spacer></v-spacer>

      <v-responsive class="hidden-sm-and-down" color="white">
        <v-text-field
          dense
          flat
          hide-details
          solo-inverted
          rounded
          placeholder="请输入文章标题查找"
          dark
          append-icon="mdi-text-search"
          v-model="searchName"
          @change="searchTitle(searchName)"
        ></v-text-field>
      </v-responsive>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" color="indigo" dark app temporary>
      <v-list>
        <v-list-item-title>
          <v-btn href="/" dark text>
            <v-icon small>mdi-home</v-icon>首页
          </v-btn>
        </v-list-item-title>

        <v-list-item
          v-model="group"
          active-class="deep-purple--text text--accent-4"
          v-for="item in cateList"
          :key="item.id"
        >
          <v-list-item-title>
            <v-btn dark text @click="gotoCate(item.id)">{{ item.name }}</v-btn>
          </v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
  </div>
</template>

<script>
export default {
  data() {
    return {
      drawer: false,
      group: null,
      valid: true,
      registerformvalid: true,
      cateList: [],
      searchName: '',
      formdata: {
        username: '',
        password: ''
      },
      checkPassword: '',
      dialog: false,
      headers: {
        Authorization: '',
        username: ''
      },
      nameRules: [
        (v) => !!v || '用户名不能为空',
        (v) =>
          (v && v.length >= 4 && v.length <= 12) ||
          '用户名必须在4到12个字符之间'
      ],
      passwordRules: [
        (v) => !!v || '密码不能为空',
        (v) =>
          (v && v.length >= 6 && v.length <= 20) || '密码必须在6到20个字符之间'
      ],
      checkPasswordRules: [
        (v) => !!v || '密码不能为空',
        (v) =>
          (v && v.length >= 6 && v.length <= 20) || '密码必须在6到20个字符之间',
        (v) => v === this.formdata.password || '密码两次输入不一致，请检查'
      ]
    }
  },
  watch: {
    group() {
      this.drawer = false
    }
  },
  created() {
    this.GetCateList()
  },
  mounted() {
    this.headers = {
      Authorization: `Bearer ${window.sessionStorage.getItem('token')}`,
      username: window.sessionStorage.getItem('username')
    }
  },
  methods: {
    // 获取分类
    async GetCateList() {
      const queryParam = {
        pagesize: 10,
        pagenum: 1,
      };
      const listOption = {
        limit: queryParam.pagesize,
        offset: (queryParam.pagenum - 1) * queryParam.pagesize,
        options: []
      }

      const {data: res} = await this.$http.post('public/GetCategoryList', {
        list_option: listOption
      })

      if (res.code !== 200) {
        this.$message.error(res.message)
        return
      }

      this.cateList = res.data.list
    },

    // 查找文章标题
    searchTitle(title) {
      if (title.length == 0) return this.$message.error('你还没填入搜索内容哦')
      this.$router.push(`/search/${title}`)
    },

    gotoCate(cid) {
      this.$router.push(`/category/${cid}`).catch((err) => err)
    },
    // 登录
    async login() {
      if (!this.$refs.loginFormRef.validate())
        return this.$message.error('输入数据非法，请检查输入的用户名和密码')
      const { data: res } = await this.$http.post('loginfront', this.formdata)
      if (res.status !== 200) return this.$message.error(res.message)
      window.sessionStorage.setItem('username', res.data)
      window.sessionStorage.setItem('user_id', res.id)
      this.$message.success('登录成功')
      this.$router.go(0)
    },

    // 退出
    loginout() {
      window.sessionStorage.clear('token')
      window.sessionStorage.clear('username')
      this.$message.success('退出成功')
      this.$router.go(0)
    },

    // 注册
    async registerUser() {
      if (!this.$refs.registerformRef.validate())
        return this.$message.error('输入数据非法，请检查输入的用户名和密码')
      const { data: res } = await this.$http.post('user/add', {
        username: this.formdata.username,
        password: this.formdata.password,
        role: 2
      })
      if (res.status !== 200) return this.$message.error(res.message)
      this.$message.success('注册成功')
      this.$router.go(0)
    }
  }
}
</script>

<style></style>
