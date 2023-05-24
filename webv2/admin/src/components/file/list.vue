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
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'file_name'">
        <span
            v-if="record.file_type === 'image/jpeg' || record.file_type === 'image/jpg' || record.file_type === 'image/png' ">
                   <a-image
                       :width="80"
                       :height="80"
                       :src="record.url"
                       fallback="data:image/png;base64,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"
                   />
        </span>
          <span v-else>
             <a :href="record.url" target="_blank">{{ record.file_name }}</a>
          </span>
        </template>

        <template v-else-if="column.key === 'action'">
        <span>
          <a @click="deleteById(record.id)">删除</a>
        </span>
        </template>

      </template>
    </a-table>
  </div>
</template>
<script lang="ts" setup>
import type {TableProps} from 'ant-design-vue';
import {message} from 'ant-design-vue';
import {useRouter} from 'vue-router';
import {reactive, ref} from 'vue';
import type {Options, Paginate} from "../../plugin/api/model/lb";
import lbstore from "../../plugin/api/lbstore"
import {format} from "date-fns";
import type {ModelFile} from "@/plugin/api/model/lbstore";
import {DefaultOption, DefaultOrderBy} from "../../plugin/api/model/lb";

// 表格的列定义
const columns = [

  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    fixed: 'left',
  },

  {
    title: '上传用户',
    dataIndex: 'creator_uid',
    key: 'creator_uid',
  },

  {
    title: '文件名称',
    dataIndex: 'file_name',
    key: 'file_name',
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
const list = ref<ModelFile[]>([]);

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
    const resp = await lbstore.getFileList({
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
  // if (id) {
  //   try {
  //     lbstore.delFile({id: id})
  //   } catch (error: any) {
  //     message.error(error);
  //   }
  // }
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

// 初始化调用一下
getList(opts)
</script>

