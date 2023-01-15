<template>
  <div>
    <a-card>
      <a-table
        :rowKey="(record,index)=>{return index}"
        :columns="columns"
        :pagination="pagination"
        :dataSource="commentList"
        bordered
        @change="handleTableChange"
      >
        <span slot="article_id" slot-scope="data">{{ articleMap[data].title }}</span>

        <span slot="user_id" slot-scope="data"
              v-if="(userMap[data]&&userMap[data].nickname)">{{ userMap[data].nickname }}</span>
        <span slot="user_id" v-else>陌生人</span>


        <span slot="status" slot-scope="data">{{ data == 1 ? '审核通过' : '未审核' }}</span>
        <template slot="action" slot-scope="data">
          <div class="actionSlot">
            <a-button
              type="primary"
              icon="edit"
              style="margin-right: 15px"
              @click="commentCheck(data.id)"
            >通过审核
            </a-button>
            <a-button
              type="primary"
              icon="info"
              style="margin-right: 15px"
              @click="commentUncheck(data.id)"
            >撤下评论
            </a-button>
            <a-button
              type="danger"
              icon="delete"
              style="margin-right: 15px"
              @click="deleteComment(data.id)"
            >删除
            </a-button>
          </div>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script>
import {formatDate} from '../../plugin/time'

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: '2%',
    key: 'id',
    align: 'center',
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    width: '10%',
    key: 'created_at',
    align: 'center',
    customRender: (val) => {
      const date = new Date(val * 1000);
      return val ? formatDate(date, 'yyyy年MM月dd日 hh:mm:ss') : '暂无'
    },
  },
  {
    title: '评论文章',
    dataIndex: 'article_id',
    width: '7%',
    key: 'article_id',
    align: 'center',
    scopedSlots: {customRender: 'article_id'},
  },
  {
    title: '评论者',
    dataIndex: 'user_id',
    width: '7%',
    key: 'user_id',
    align: 'center',
    scopedSlots: {customRender: 'user_id'},
  },
  {
    title: '评论内容',
    dataIndex: 'content',
    width: '20%',
    key: 'content',
    align: 'center',
  },
  {
    title: '评论状态',
    dataIndex: 'status',
    width: '7%',
    key: 'status',
    align: 'center',
    scopedSlots: {customRender: 'status'},
  },
  {
    title: '操作',
    width: '20%',
    key: 'action',
    align: 'center',
    scopedSlots: {customRender: 'action'},
  },
]
export default {
  data() {
    return {
      commentList: [],
      commentInfo: {
        status: 1,
      },
      pagination: {
        pageSizeOptions: ['5', '10', '20'],
        pageSize: 5,
        total: 0,
        showSizeChanger: true,
        showTotal: (total) => `共${total}条`,
      },
      columns,
      articleMap: {},
      userMap: {},
      queryParam: {
        pagesize: 10,
        pagenum: 1,
      },
    }
  },
  created() {
    this.getCommentList()
  },
  methods: {
    // 获取评论列表
    async getCommentList() {
      const listOption = {
        limit: this.queryParam.pagesize,
        offset: (this.queryParam.pagenum - 1) * this.queryParam.pagesize,
        options: []
      }
      const {data: res} = await this.$http.post('blog/GetCommentList', {
        list_option: listOption
      })

      if (res.code !== 200) {
        this.$message.error(res.message)
        if (res.code === 401) {
          window.sessionStorage.clear()
          await this.$router.push('/login')
        }
        return
      }

      this.commentList = res.data.list || []
      this.userMap = res.data.user_map || {}
      this.articleMap = res.data.article_map || {}
      this.pagination.total = res.data.page.total
    },

    // 更改分页
    handleTableChange(pagination, filters, sorter) {
      const pager = {...this.pagination};
      pager.current = pagination.current
      pager.pageSize = pagination.pageSize
      this.queryParam.pagesize = pagination.pageSize
      this.queryParam.pagenum = pagination.current

      if (pagination.pageSize !== this.pagination.pageSize) {
        this.queryParam.pagenum = 1
        pager.current = 1
      }
      this.pagination = pager
      this.getCommentList()
    },

    // 通过审核
    commentCheck(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '要通过审核吗？',
        onOk: async () => {
          const {data: res} = await this.$http.post(`blog/GetComment`, {
            id: Number(id)
          })

          if (res.data.comment.status === 1) {
            return this.$message.error('该评论已处于显示状态，无需审核')
          }

          res.data.comment.status = 1
          const {data: uRes} = await this.$http.post(`blog/UpdateComment`, {
            comment: res.data.comment,
          })

          if (uRes.code !== 200) {
            this.$message.error(uRes.message)
            if (uRes.code === 401) {
              window.sessionStorage.clear()
              await this.$router.push('/login')
            }
            return
          }

          this.$message.success('审核成功')
          await this.getCommentList()
        },
        onCancel: () => {
          this.$message.info('已取消')
        },
      })
    },

    // 撤下评论
    commentUncheck(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '要撤下该评论吗？',
        onOk: async () => {
          const {data: res} = await this.$http.post(`blog/GetComment`, {
            id: Number(id)
          })

          if (res.data.comment.status === 2) {
            return this.$message.error('该评论已处于未审核状态，无需撤下')
          }

          if (res.data.status === 2) {
            return this.$message.error('该评论已处于未审核状态，无需撤下')
          }

          res.data.comment.status = 2
          const {data: uRes} = await this.$http.post(`blog/UpdateComment`, {
            comment: res.data.comment,
          })

          if (uRes.code !== 200) {
            this.$message.error(uRes.message)
            if (uRes.code === 401) {
              window.sessionStorage.clear()
              await this.$router.push('/login')
            }
            return
          }
          this.$message.success('评论已撤下')
          await this.getCommentList()
        },
        onCancel: () => {
          this.$message.info('已取消')
        },
      })
    },

    // 删除评论
    deleteComment(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '要删除吗？',
        onOk: async () => {
          const {data: res} = await this.$http.post(`blog/DelComment`, {
            id: Number(id)
          })

          if (res.code !== 200) {
            this.$message.error(res.message)
            if (res.code === 401) {
              window.sessionStorage.clear()
              await this.$router.push('/login')
            }
            return
          }

          this.$message.success('删除成功')
          await this.getCommentList()
        },
        onCancel: () => {
          this.$message.info('已取消')
        },
      })
    },
  },
}
</script>
<style lang="">
</style>
