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

        <template v-if="column.key === 'status'">
          <a-tag v-if="record.status === 1">正常</a-tag>
          <a-tag v-else-if="record.status ===2">审核中</a-tag>
          <a-tag v-else-if="record.status ===3">撤下</a-tag>
          <a-tag v-else>待审核</a-tag>
        </template>

        <template v-if="column.key === 'action'">
        <span>
          <a @click="showDetailsModal(record.id)">查看</a>
          <a-divider type="vertical"/>
          <a @click="showUpdateModal(record.id)">编辑</a>
          <a-divider type="vertical"/>
          <a @click="deleteById(record.id)">删除</a>
        </span>
        </template>

      </template>
    </a-table>

    <!--更新弹窗-->
    <Update @handleComplete="handleComplete" ref="updateRef"/>

    <!--详情弹窗-->
    <Details @handleComplete="handleComplete" ref="detailsRef"/>
  </div>
</template>
<script lang="ts" setup>
import Update from './update.vue';
import Details from './details.vue';
import type {TableProps} from 'ant-design-vue';
import {message} from 'ant-design-vue';
import {useRouter} from 'vue-router';
import {reactive, ref} from 'vue';
import type {Options, Paginate} from "@/plugin/api/model/lb";
import lbblog from "../../plugin/api/lbblog"
import {format} from "date-fns";
import {DefaultOption, DefaultOrderBy} from "@/plugin/api/model/lb";

// 表格的列定义
const columns = [

  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    fixed: 'left',
  },

  {
    title: '文章作者',
    dataIndex: 'article_id',
    key: 'article_id',
  },

  {
    title: '评论用户',
    dataIndex: 'user_id',
    key: 'user_id',
  },

  {
    title: '用户邮箱',
    dataIndex: 'user_email',
    key: 'user_email',
  },

  {
    title: '内容',
    dataIndex: 'content',
    key: 'content',
  },

  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
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
    const resp = await lbblog.getCommentList({
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
const deleteById = (id: number | undefined) => {
  if (id) {
    try {
      lbblog.delComment({id: id})
    } catch (error: any) {
      message.error(error);
    }
  }
  resetSearch()
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
const resetSearch = () => {
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

  getList(opts)
}

// 指向添加组件
const addRef = ref()

// 展示添加窗口
const showAddModal = () => {
  addRef.value.show()
}

// 指向更新组件
const updateRef = ref()

// 展示更新窗口
const showUpdateModal = (id: number | undefined) => {
  updateRef.value.show(id)
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

