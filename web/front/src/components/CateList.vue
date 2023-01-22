<template>
  <v-container>
    <div v-if="total == 0 && isLoad" class="d-flex justify-center align-center">
      <div>
        <v-alert class="ma-5" dense outlined type="error"
          >抱歉，暂无数据！</v-alert
        >
      </div>
    </div>
    <v-sheet>
      <v-card
        class="ma-3"
        v-for="item in artList"
        :key="item.id"
        link
        @click="$router.push(`/article/detail/${item.id}`)"
      >
        <v-row no-gutters class="d-flex align-center">
          <v-avatar class="ma-3 hidden-sm-and-down" size="125" tile>
            <v-img :src="item.img"></v-img>
          </v-avatar>
          <v-col>
            <v-card-title>
              <v-chip color="purple" outlined label class="mr-3 white--text">{{
                  CateMap[item.category_id].name
              }}</v-chip>
              <div>{{ item.title }}</div>
            </v-card-title>
            <v-card-subtitle class="mt-1" v-text="item.desc"></v-card-subtitle>
            <v-divider class="mx-4"></v-divider>
            <v-card-text class="d-flex align-center">
              <div class="d-flex align-center">
                <v-icon class="mr-1" small>{{ 'mdi-calendar-month' }}</v-icon>
                <span>{{
                    item.created_at*1000 | dateformat('YYYY-MM-DD HH:MM')
                }}</span>
              </div>
              <div class="mx-4 d-flex align-center">
                <v-icon class="mr-1" small>{{ 'mdi-comment' }}</v-icon>
                <span>{{ item.comment_count }}</span>
              </div>
              <div class="mx-1 d-flex align-center">
                <v-icon class="mr-1" small>{{ 'mdi-eye' }}</v-icon>
                <span>{{ item.read_count }}</span>
              </div>
            </v-card-text>
          </v-col>
        </v-row>
      </v-card>
      <v-col>
        <div class="text-center">
          <v-pagination
            total-visible="7"
            v-model="queryParam.pagenum"
            :length="Math.ceil(total / queryParam.pagesize)"
            @input="getArtList()"
          ></v-pagination>
        </div>
      </v-col>
    </v-sheet>
  </v-container>
</template>
<script>
export default {
  props: ['cid'],
  data() {
    return {
      artList: [],
      queryParam: {
        pagesize: 5,
        pagenum: 1
      },
      CateMap: {},
      total: 0,
      isLoad: false
    }
  },
  mounted() {
    this.getArtList()
  },
  methods: {
    // 获取文章列表
    async getArtList() {
      const listOption = {
        limit: this.queryParam.pagesize,
        offset: (this.queryParam.pagenum - 1) * this.queryParam.pagesize,
        options: [
          {
            type: 1,
            value: this.cid,
          },
          {
            type: 2,
            value: this.queryParam.title,
          }
        ]
      }
      const {data: res} = await this.$http.post('public/GetArticleList', {
        list_option: listOption
      })
      if (res.code !== 200) {
        this.$message.error(res.message)
        return
      }
      this.artList = res.data.list || []
      this.CateMap = res.data.category_map || {}
      this.total = res.data.page.total
      this.isLoad = true
    }
  }
}
</script>
<style scoped>
.nodate {
  width: 100%;
  height: 100%;
}
</style>
