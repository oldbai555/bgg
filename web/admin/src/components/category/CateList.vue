<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="4">
          <a-button type="primary" @click="addCateVisible = true">新增分类</a-button>
        </a-col>
      </a-row>

      <a-table
        rowKey="id"
        :columns="columns"
        :pagination="pagination"
        :dataSource="Catelist"
        bordered
        @change="handleTableChange"
      >
        <template slot="action" slot-scope="data">
          <div class="actionSlot">
            <a-button type="primary" icon="edit" style="margin-right: 15px" @click="editCate(data.id)">编辑</a-button>
            <a-button type="danger" icon="delete" style="margin-right: 15px" @click="deleteCate(data.id)"
            >删除
            </a-button
            >
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 新增分类区域 -->
    <a-modal
      closable
      title="新增分类"
      :visible="addCateVisible"
      width="60%"
      @ok="addCateOk"
      @cancel="addCateCancel"
      destroyOnClose
    >
      <a-form-model :model="newCate" :rules="addCateRules" ref="addCateRef">
        <a-form-model-item label="分类名称" prop="name">
          <a-input v-model="newCate.name"></a-input>
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <!-- 编辑分类区域 -->
    <a-modal
      closable
      destroyOnClose
      title="编辑分类"
      :visible="editCateVisible"
      width="60%"
      @ok="editCateOk"
      @cancel="editCateCancel"
    >
      <a-form-model :model="CateInfo" :rules="CateRules" ref="addCateRef">
        <a-form-model-item label="分类名称" prop="name">
          <a-input v-model="CateInfo.name"></a-input>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </div>
</template>

<script>
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: '10%',
    key: 'id',
    align: 'center',
  },
  {
    title: '分类名',
    dataIndex: 'name',
    width: '20%',
    key: 'name',
    align: 'center',
  },
  {
    title: '操作',
    width: '30%',
    key: 'action',
    align: 'center',
    scopedSlots: {customRender: 'action'},
  },
]

export default {
  data() {
    return {
      pagination: {
        pageSizeOptions: ['5', '10', '20'],
        pageSize: 5,
        total: 0,
        showSizeChanger: true,
        showTotal: (total) => `共${total}条`,
      },
      Catelist: [],
      CateInfo: {
        name: '',
        id: 0,
      },
      newCate: {
        name: '',
      },
      columns,
      queryParam: {
        pagesize: 5,
        pagenum: 1,
      },
      editVisible: false,
      CateRules: {
        name: [
          {
            validator: (rule, value, callback) => {
              if (this.CateInfo.name === '') {
                callback(new Error('请输入分类名'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      addCateRules: {
        name: [
          {
            validator: (rule, value, callback) => {
              if (this.newCate.name === '') {
                callback(new Error('请输入分类名'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      editCateVisible: false,
      addCateVisible: false,
    }
  },
  created() {
    this.getCateList()
  },
  methods: {
    // 获取分类列表
    async getCateList() {
      const listOption = {
        limit: this.queryParam.pagesize,
        offset: (this.queryParam.pagenum - 1) * this.queryParam.pagesize,
        options: []
      }
      const {data: res} = await this.$http.post('blog/GetCategoryList', {
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

      this.Catelist = res.data.list
      this.pagination.total = res.data.page.total
    },

    // 更改分页
    handleTableChange(pagination, filters, sorter) {
      var pager = {...this.pagination}
      pager.current = pagination.current
      pager.pageSize = pagination.pageSize
      this.queryParam.pagesize = pagination.pageSize
      this.queryParam.pagenum = pagination.current

      if (pagination.pageSize !== this.pagination.pageSize) {
        this.queryParam.pagenum = 1
        pager.current = 1
      }
      this.pagination = pager
      this.getCateList()
    },

    // 删除分类
    deleteCate(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '确定要删除该分类吗？一旦删除，无法恢复',
        onOk: async () => {
          const {data: res} = await this.$http.post(`blog/DelCategory`, {
            id: id,
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
          await this.getCateList()
        },
        onCancel: () => {
          this.$message.info('已取消删除')
        },
      })
    },

    // 新增分类
    addCateOk() {
      this.$refs.addCateRef.validate(async (valid) => {
        if (!valid) {
          this.$message.error('参数不符合要求，请重新输入')
          return
        }

        const cate = {
          name: this.newCate.name
        }

        const {data: res} = await this.$http.post('blog/AddCategory', {
          category: cate,
        })

        if (res.code !== 200) {
          this.$message.error(res.message)
          if (res.code === 401) {
            window.sessionStorage.clear()
            await this.$router.push('/login')
          }
          return
        }

        this.$refs.addCateRef.resetFields()
        this.addCateVisible = false
        this.$message.success('添加分类成功')
        await this.getCateList()
      })
    },

    addCateCancel() {
      this.$refs.addCateRef.resetFields()
      this.addCateVisible = false
      this.$message.info('新增分类已取消')
    },

    // 编辑分类
    async editCate(id) {
      this.editCateVisible = true
      const {data: res} = await this.$http.post(`blog/GetCategory`, {
        id: id,
      })

      if (res.code !== 200) {
        this.$message.error(res.message)
        if (res.code === 401) {
          window.sessionStorage.clear()
          await this.$router.push('/login')
        }
        return
      }

      this.CateInfo = res.data.category
    },

    editCateOk() {
      this.$refs.addCateRef.validate(async (valid) => {
        if (!valid) {
          this.$message.error('参数不符合要求，请重新输入')
          return
        }

        const {data: res} = await this.$http.post(`blog/UpdateCategory`, {
          category: this.CateInfo,
        })

        if (res.code !== 200) {
          this.$message.error(res.message)
          if (res.code === 401) {
            window.sessionStorage.clear()
            await this.$router.push('/login')
          }
          return
        }

        this.editCateVisible = false
        this.$message.success('更新分类信息成功')
        await this.getCateList()
      })
    },

    editCateCancel() {
      this.$refs.addCateRef.resetFields()
      this.editCateVisible = false
      this.$message.info('编辑已取消')
    },
  },
}
</script>

<style scoped>
.actionSlot {
  display: flex;
  justify-content: center;
}
</style>
