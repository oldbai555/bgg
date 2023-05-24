<template>
  <div>
    <!--顶部搜索栏-->
    <a-card :title="routerTitle">
      <a-row :gutter="[16,16]">
        <a-col :span="6">
          <a-input-search
              placeholder="一个万能的搜索框......"
              enter-button
              allowClear
              @search="handleSearch"
              v-model:value="searchValue"
          />
        </a-col>
        <a-col :span="4">
          <a-button type="info" @click="resetSearch">重置</a-button>
        </a-col>
      </a-row>
      <a-row justify="end">
        <a-col :span="2">
          <a-button type="primary" @click="showAddModal">新增</a-button>
        </a-col>
      </a-row>
    </a-card>

    <!--表格部分-->
    <a-table
        :columns="columns"
        :data-source="list"
        :pagination="pagination"
        :loading="loading"
        @change="handleTableChange"
        :scroll="{ x: 1500, y: 300 }"
    >
      <template #bodyCell="{ column, record }">

        <template v-if="column.key === 'avatar'">
          <a-image
              :width="80"
              :height="80"
              :src="record.avatar"
              fallback="data:image/png;base64,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"
          />
        </template>
        <template v-else-if="column.key === 'github'">
          <a :href="record?.github" target="_blank">点击跳转</a>
        </template>

        <template v-else-if="column.key === 'action'">
        <span>
          <a @click="showDetailsModal(record.id)">查看</a>
          <a-divider type="vertical"/>
          <a @click="deleteById(record.id)">删除</a>
        </span>
        </template>

      </template>
    </a-table>

    <!--添加弹窗-->
    <Add @handleComplete="handleComplete" ref="addRef"/>

    <!--详情弹窗-->
    <Details @handleComplete="handleComplete" ref="detailsRef"/>
  </div>
</template>
<script lang="ts" setup>
import Add from './add.vue';
import Details from './details.vue';
import type {TableProps} from 'ant-design-vue';
import {message} from 'ant-design-vue';
import {useRouter} from 'vue-router';
import {reactive, ref} from 'vue';
import type {Options, Paginate} from "@/plugin/api/model/lb";
import {DefaultOption, DefaultOrderBy} from "@/plugin/api/model/lb";
import lbuser from "../../plugin/api/lbuser"
import {format} from 'date-fns'
import {ModelUser_Role} from "@/plugin/api/model/lbuser";

// 表格的列定义
const columns = [

  {
    title: 'id',
    dataIndex: 'id',
    key: 'id',
    fixed: 'left',
  },

  {
    title: '头像',
    dataIndex: 'avatar',
    key: 'avatar',
  },

  {
    title: '账号',
    dataIndex: 'username',
    key: 'username',
  },

  {
    title: '昵称',
    dataIndex: 'nickname',
    key: 'nickname',
  },

  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
  },

  {
    title: 'GitHub',
    dataIndex: 'github',
    key: 'github',
  },

  {
    title: '描述',
    dataIndex: 'desc',
    key: 'desc',
  },

  {
    title: '角色',
    dataIndex: 'role',
    key: 'role',
    customRender: (val: Number) => {
      switch (val) {
        case ModelUser_Role.RoleAdmin:
          return "管理员";
        default:
          return "-";
      }
    }
  },

  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    customRender: (val: any) => {
      if (val?.text) {
        return format(val?.text * 1000, "yyyy年MM月dd日 hh:mm:ss")
      }
      return
    }
  },

  {
    title: '操作',
    key: 'action',
    fixed: 'right',
  },
];

// 默认分页大小
const defaultPage = 1
const defaultSize = 10

// 路由
const router = useRouter();

// 路由标题
const routerTitle = ref(router.currentRoute.value.meta?.title);

// 分页
const pagination = reactive({
  total: 0,
  current: defaultPage,
  pageSize: defaultSize,
});

// 接收的列表数据
const list = ref([{}]);

// 表格的loading动画
const loading = ref(false);

// 搜索框的内容
const searchValue = ref<string>("")

// 条件构造器
const opts = reactive<Options>({
  opt_list: [
    {
      key: DefaultOption.DefaultOptionOrderBy,
      value: DefaultOrderBy.DefaultOrderByCreatedAtDesc.toString(),
    }
  ],
  size: defaultSize,
  page: defaultPage,
  skip_total: false,
});

// 请求列表数据
const getList = async (listOption: Options) => {
  // 开启 Loading
  loading.value = true;
  try {
    const resp = await lbuser.getUserList({
      options: listOption,
    });

    // 列表赋值
    list.value = resp.list;

    // 更新一下分页
    changePagination(resp.paginate);
  } catch (error: any) {
    message.error(error);
  }
  // 关闭 Loading
  loading.value = false;
};

// 删除
const deleteById = async (id: number | undefined) => {
  if (id) {
    try {
      await lbuser.delUser({id: id})
    } catch (error: any) {
      message.error(error);
    }
  }
  await resetSearch()
}

// 更新分页
const changePagination = (paginate: Paginate | undefined) => {
  pagination.current = Number(paginate?.page);
  pagination.pageSize = Number(paginate?.size);
  pagination.total = Number(paginate?.total);
}

// 表格更新事件
const handleTableChange: TableProps['onChange'] = (pagination, filters, sorter, {currentDataSource}) => {
  opts.page = Number(pagination.current);
  opts.size = Number(pagination.pageSize);
  getList(opts);

  console.log("filter :", filters);
  console.log("sorter :", sorter);
  console.log("currentDataSource :", currentDataSource);
};

// 执行搜索
const handleSearch = (searchValue: string) => {
  console.log(searchValue)
}

// 重置搜索
const resetSearch = async () => {
  opts.size = defaultSize;
  opts.page = defaultPage;
  opts.skip_total = false;
  opts.opt_list = [
    {
      key: DefaultOption.DefaultOptionOrderBy,
      value: DefaultOrderBy.DefaultOrderByCreatedAtDesc.toString(),
    }
  ];
  searchValue.value = "";

  await getList(opts)
}

// 指向添加组件
const addRef = ref()

// 展示添加窗口
const showAddModal = () => {
  addRef.value.show()
}

// 指向详情组件
const detailsRef = ref()

// 展示详情窗口
const showDetailsModal = (id: number | undefined) => {
  detailsRef.value.show(id)
}

// 完成弹窗操作
const handleComplete = () => {
  resetSearch()
}

// 初始化调用一下
getList(opts)
</script>

