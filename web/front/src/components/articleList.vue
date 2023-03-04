<template>
  <v-col>
    <v-card
        class="ma-3 d-flex flex-no-wrap justify-space-between align-center"
        v-for="item in artList"
        :key="item.id"
        link
        @click="$router.push(`article/detail/${item.id}`)"
    >
      <v-avatar class="ma-3 hidden-sm-and-down" size="125" tile>
        <v-img :src="item.img"></v-img>
      </v-avatar>

      <v-col>
        <v-card-title>
          <v-chip color="pink" outlined label class="mr-3 white--text">
            {{
              CateMap[item.category_id].name
            }}
          </v-chip>
          <div>{{ item.title }}</div>
        </v-card-title>
        <v-card-subtitle class="mt-1" v-text="item.desc"></v-card-subtitle>
        <v-divider class="mx-4"></v-divider>
        <v-card-text class="d-flex align-center">
          <div class="d-flex align-center">
            <v-icon class="mr-1" small>{{ 'mdi-calendar-month' }}</v-icon>
            <span>{{ item.created_at*1000 | dateformat('YYYY-MM-DD') }}</span>
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
    </v-card>
    <div class="text-center">
      <v-pagination
          color="indigo"
          total-visible="7"
          v-model="queryParam.pagenum"
          :length="Math.ceil(total / queryParam.pagesize)"
          @input="getArtList()"
      ></v-pagination>
    </div>
  </v-col>
</template>
<script>

export default {
  data() {
    return {
      artList: [],
      queryParam: {
        pagesize: 5,
        pagenum: 1,
        title: "",
      },
      CateMap: {},
      count: 0,
      total: 0
    }
  },
  created() {
  },
  mounted() {
    this.getArtList()
  },
  filters: {},
  methods: {
    // 获取文章列表
    async getArtList() {
      const listOption = {
        limit: this.queryParam.pagesize,
        offset: (this.queryParam.pagenum - 1) * this.queryParam.pagesize,
        options: [
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
    }
  }
}
</script>
<style lang=""></style>
