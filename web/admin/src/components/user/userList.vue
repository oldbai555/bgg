<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="6">
          <a-input-search
            v-model="queryParam.username"
            placeholder="输入用户名查找"
            enter-button
            allowClear
            @search="searchUser"
          />
        </a-col>
        <a-col :span="4">
          <a-button type="primary" @click="addUserVisible = true">新增</a-button>
        </a-col>
      </a-row>

      <a-table
        :rowKey="(record,index)=>{return index}"
        :columns="columns"
        :pagination="pagination"
        :dataSource="userlist"
        bordered
        @change="handleTableChange"
      >

        <span slot="role" slot-scope="data">{{ data === 1 ? '管理员' : '订阅者' }}</span>

        <template slot="action" slot-scope="data">
          <div class="actionSlot">
            <a-button
              type="primary"
              icon="edit"
              style="margin-right: 15px"
              @click="editUser(data.id)"
            >编辑
            </a-button>
            <a-button
              type="danger"
              icon="delete"
              style="margin-right: 15px"
              @click="deleteUser(data.id)"
            >删除
            </a-button>
            <a-button type="info" icon="info" @click="ChangePassword(data.id)">修改密码</a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 新增用户区域 -->
    <a-modal
      closable
      title="新增用户"
      :visible="addUserVisible"
      width="60%"
      @ok="addUserOk"
      @cancel="addUserCancel"
      destroyOnClose
    >
      <a-form-model :model="newUser" :rules="addUserRules" ref="addUserRef">
        <a-form-model-item label="用户名" prop="username">
          <a-input v-model="newUser.username"></a-input>
        </a-form-model-item>
        <a-form-model-item has-feedback label="密码" prop="password">
          <a-input-password v-model="newUser.password"></a-input-password>
        </a-form-model-item>
        <a-form-model-item has-feedback label="确认密码" prop="checkpass">
          <a-input-password v-model="newUser.checkpass"></a-input-password>
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <!-- 编辑用户区域 -->
    <a-modal
      closable
      destroyOnClose
      title="编辑用户"
      :visible="editUserVisible"
      width="60%"
      @ok="editUserOk"
      @cancel="editUserCancel"
    >
      <a-form-model :model="userInfo" :rules="userRules" ref="addUserRef">
        <a-form-model-item label="用户名" prop="username">
          <a-input v-model="userInfo.username"></a-input>
        </a-form-model-item>
        <a-form-model-item label="是否为管理员">
          <a-switch
            :checked="IsAdmin"
            checked-children="是"
            un-checked-children="否"
            @change="adminChange"
          />
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <!-- 修改密码 -->
    <a-modal
      closable
      title="修改密码"
      :visible="changePasswordVisible"
      width="60%"
      @ok="changePasswordOk"
      @cancel="changePasswordCancel"
      destroyOnClose
    >
      <a-form-model :model="changePassword" :rules="changePasswordRules" ref="changePasswordRef">
        <a-form-model-item has-feedback label="密码" prop="password">
          <a-input-password v-model="changePassword.password"></a-input-password>
        </a-form-model-item>
        <a-form-model-item has-feedback label="确认密码" prop="checkpass">
          <a-input-password v-model="changePassword.checkpass"></a-input-password>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </div>
</template>

<script>
import {formatDate} from '../../plugin/time'

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: '10%',
    key: 'id',
    align: 'center',
  },
  {
    title: '用户名',
    dataIndex: 'username',
    width: '20%',
    key: 'username',
    align: 'center',
  },
  {
    title: '注册时间',
    dataIndex: 'created_at',
    width: '20%',
    key: 'created_at',
    align: 'center',
    customRender: (val) => {
      const date = new Date(val * 1000);
      return val ? formatDate(date, 'yyyy年MM月dd日 hh:mm:ss') : '暂无'
    },
  },
  {
    title: '角色',
    dataIndex: 'role',
    width: '20%',
    key: 'role',
    align: 'center',
    scopedSlots: {customRender: 'role'},
  },
  {
    title: '昵称',
    dataIndex: 'nickname',
    width: '20%',
    key: 'nickname',
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
      userlist: [],
      userInfo: {
        username: '',
        password: '',
        role: 2,
        checkPass: '',
      },
      newUser: {
        username: '',
        password: '',
        role: 2,
        checkPass: '',
      },
      changePassword: {
        id: 0,
        password: '',
        oldPassword: '',
        checkPass: '',
      },
      columns,
      queryParam: {
        username: '',
        pagesize: 5,
        pagenum: 1,
      },
      editVisible: false,
      userRules: {
        username: [
          {
            validator: (rule, value, callback) => {
              if (this.userInfo.username === '') {
                callback(new Error('请输入用户名'))
              }
              if ([...this.userInfo.username].length < 4 || [...this.userInfo.username].length > 12) {
                callback(new Error('用户名应当在4到12个字符之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        password: [
          {
            validator: (rule, value, callback) => {
              if (this.userInfo.password === '') {
                callback(new Error('请输入密码'))
              }
              if ([...this.userInfo.password].length < 6 || [...this.userInfo.password].length > 20) {
                callback(new Error('密码应当在6到20位之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        checkpass: [
          {
            validator: (rule, value, callback) => {
              if (this.userInfo.checkpass === '') {
                callback(new Error('请输入密码'))
              }
              if (this.userInfo.password !== this.userInfo.checkpass) {
                callback(new Error('密码不一致，请重新输入'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      addUserRules: {
        username: [
          {
            validator: (rule, value, callback) => {
              if (this.newUser.username === '') {
                callback(new Error('请输入用户名'))
              }
              if ([...this.newUser.username].length < 4 || [...this.newUser.username].length > 12) {
                callback(new Error('用户名应当在4到12个字符之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        password: [
          {
            validator: (rule, value, callback) => {
              if (this.newUser.password === '') {
                callback(new Error('请输入密码'))
              }
              if ([...this.newUser.password].length < 6 || [...this.newUser.password].length > 20) {
                callback(new Error('密码应当在6到20位之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        checkpass: [
          {
            validator: (rule, value, callback) => {
              if (this.newUser.checkpass === '') {
                callback(new Error('请输入密码'))
              }
              if (this.newUser.password !== this.newUser.checkpass) {
                callback(new Error('密码不一致，请重新输入'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      changePasswordRules: {
        password: [
          {
            validator: (rule, value, callback) => {
              if (this.changePassword.password === '') {
                callback(new Error('请输入密码'))
              }
              if ([...this.changePassword.password].length < 6 || [...this.changePassword.password].length > 20) {
                callback(new Error('密码应当在6到20位之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        checkpass: [
          {
            validator: (rule, value, callback) => {
              if (this.changePassword.checkpass === '') {
                callback(new Error('请输入密码'))
              }
              if (this.changePassword.password !== this.changePassword.checkpass) {
                callback(new Error('密码不一致，请重新输入'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      editUserVisible: false,
      addUserVisible: false,
      changePasswordVisible: false,
    }
  },
  created() {
    this.getUserList()
  },
  computed: {
    IsAdmin: function () {
      return this.userInfo.role === 1;
    },
  },
  methods: {
    // 获取用户列表
    async getUserList() {
      const listOption = {
        limit: this.queryParam.pagesize,
        offset: (this.queryParam.pagenum - 1) * this.queryParam.pagesize,
        options: []
      }
      const {data: res} = await this.$http.post('user/GetUserList', {
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

      this.userlist = res.data.list
      this.pagination.total = res.data.page.total
    },

    // 搜索用户
    async searchUser() {
      const listOption = {
        limit: this.queryParam.pagesize,
        offset: (this.queryParam.pagenum - 1) * this.queryParam.pagesize,
        options: [
          {
            type: 1,
            value: this.queryParam.username + '',
          },
        ]
      }

      const {data: res} = await this.$http.post('user/GetUserList', {
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

      this.userlist = res.data.list
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
      this.getUserList()
    },
    // 删除用户
    deleteUser(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '确定要删除该用户吗？一旦删除，无法恢复',
        onOk: async () => {
          const {data: res} = await this.$http.post(`user/DelUser`, {
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
          await this.getUserList()
        },
        onCancel: () => {
          this.$message.info('已取消删除')
        },
      })
    },
    // 新增用户
    addUserOk() {
      this.$refs.addUserRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const {data: res} = await this.$http.post('user/AddUser', {
          user: {
            username: this.newUser.username,
            password: this.newUser.password,
            role: this.newUser.role,
          }
        })

        if (res.code !== 200) {
          this.$message.error(res.message)
          if (res.code === 401) {
            window.sessionStorage.clear()
            await this.$router.push('/login')
          }
          return
        }

        this.$refs.addUserRef.resetFields()
        this.addUserVisible = false
        this.$message.success('添加用户成功')
        await this.getUserList()
      })
    },
    addUserCancel() {
      this.$refs.addUserRef.resetFields()
      this.addUserVisible = false
      this.$message.info('新增用户已取消')
    },
    adminChange(checked) {
      if (checked) {
        this.userInfo.role = 1
      } else {
        this.userInfo.role = 2
      }
    },
    // 编辑用户
    async editUser(id) {
      this.editUserVisible = true
      const {data: res} = await this.$http.post(`user/GetUser`, {
        id: id,
      })
      this.userInfo = res.data.user
    },
    editUserOk() {
      this.$refs.addUserRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const {data: res} = await this.$http.post(`user/UpdateUserNameWithRole`, {
          id: this.userInfo.id,
          username: this.userInfo.username,
          role: this.userInfo.role,
        })

        if (res.code !== 200) {
          this.$message.error(res.message)
          if (res.code === 401) {
            window.sessionStorage.clear()
            await this.$router.push('/login')
          }
          return
        }

        this.editUserVisible = false
        this.$message.success('更新用户信息成功')
        this.$refs.addUserRef.resetFields()
        await this.getUserList()
      })
    },
    editUserCancel() {
      this.$refs.addUserRef.resetFields()
      this.editUserVisible = false
      this.$message.info('编辑已取消')
    },

    // 修改密码
    async ChangePassword(id) {
      this.changePasswordVisible = true
      const {data: res} = await this.$http.post(`user/GetUser`, {
        id: id,
      })
      this.changePassword.id = res.data.user.id
      this.changePassword.oldPassword = res.data.user.password
    },
    changePasswordOk() {
      this.$refs.changePasswordRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const {data: res} = await this.$http.post(`user/ResetPassword`, {
          new_password: this.changePassword.password,
          old_password: this.changePassword.oldPassword,
          id: this.changePassword.id,
        })

        if (res.code !== 200) {
          this.$message.error(res.message)
          if (res.code === 401) {
            window.sessionStorage.clear()
            await this.$router.push('/login')
          }
          return
        }

        this.changePasswordVisible = false
        this.$message.success('修改密码成功')
        this.$refs.changePasswordRef.resetFields()
        await this.getUserList()
      })
    },
    changePasswordCancel() {
      this.$refs.changePasswordRef.resetFields()
      this.changePasswordVisible = false
      this.$message.info('已取消')
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
