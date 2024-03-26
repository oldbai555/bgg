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
          <article-category-select
              style="width: 130px"
              placeholder="请选择分类"
              @change="categoryChange"
          />
        </a-col>
        <a-col :span="2">
          <a-button type="info" @click="resetSearch">重置</a-button>
        </a-col>
      </a-row>
      <a-row justify="end">
        <a-col :span="2">
          <a-button type="primary" @click="gotoEdit(undefined)">新增</a-button>
        </a-col>
      </a-row>
    </a-card>
    <!--表格部分-->
    <a-table
        :columns="columns"
        :data-source="data"
        :pagination="pagination"
        :loading="loading"
        @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">

        <template v-if="column.key === 'img'">
          <a-image
              :width="80"
              :height="80"
              :src="record.img"
              fallback="data:image/png;base64,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"
          />
        </template>

        <template v-else-if="column.key === 'category_id'">
          <a-tag v-if="(categoryMap[record.category_id]?.name)">{{ categoryMap[record.category_id].name }}</a-tag>
          <a-tag v-else>暂无分类</a-tag>
        </template>

        <template v-else-if="column.key === 'action'">
        <span>
          <a @click="gotoEdit(record.id)">编辑</a>
          <a-divider type="vertical"/>
          <a @click="deleteArticle(record.id)">删除</a>
        </span>
        </template>

      </template>
    </a-table>
  </div>
</template>

<script lang="ts" setup>
import {reactive, ref} from 'vue';
import type {TableProps} from 'ant-design-vue';
import {message} from 'ant-design-vue';
import blogApi from "@/plugin/api/lbblog";
import router from "@/router";
import {replaceOps} from "@/plugin/utils/option_opt";
import {format} from "date-fns";
import ArticleCategorySelect from "../global_components/article_category_select.vue"

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '封面',
    dataIndex: 'img',
    key: 'img',
  },
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
  },
  {
    title: '描述',
    dataIndex: 'desc',
    key: 'desc',
  },
  {
    title: '分类',
    key: 'category_id',
    dataIndex: 'category_id',
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
  },
];

// 默认分页大小
const defaultPage = 1
const defaultSize = 10

// 分页
const pagination = reactive({
  total: 0,
  current: defaultPage,
  pageSize: defaultSize,
});

// 条件构造器
const opts = reactive<lb.ListOption>({
  Options: [
    {
      key: lb.DefaultListOption.DefaultListOptionOrderBy,
      value: lb.DefaultOrderBy.DefaultOrderByCreatedAtDesc.toString(),
    }
  ],
  size: defaultSize,
  page: defaultPage,
  skip_total: false,
});

const routerTitle = "文章管理"
const data = ref([{}])
const loading = ref(false)
const categoryMap = ref({})
// 搜索框的内容
const searchValue = ref<string>("")

// 更新分页
const changePagination = (paginate: lb.Paginate | undefined) => {
  pagination.current = Number(paginate?.page);
  pagination.pageSize = Number(paginate?.size);
  pagination.total = Number(paginate?.total);
}

const getArticleList = async (listOption: lb.ListOption) => {
  loading.value = true
  try {
    const resp = await blogApi.getArticleList({
      options: listOption,
    });
    // 更新一下分页
    changePagination(resp.paginate);
    categoryMap.value = resp.category_map
    data.value = resp.list
  } catch (error: any) {
    message.error(error)
  }
  loading.value = false
};
getArticleList(opts)

const handleTableChange: TableProps['onChange'] = (
    pagination, filters, sorter, {currentDataSource}
) => {
  opts.page = Number(pagination.current);
  opts.size = Number(pagination.pageSize);
  getArticleList(opts)

  console.log("filter :", filters)
  console.log("sorter :", sorter)
  console.log("currentDataSource :", currentDataSource)
};

const gotoEdit = async (id: number | undefined) => {
  await router.push({
    name: "article_edit",
    query: {
      id: id,
    }
  })
}

const deleteArticle = async (id: number | undefined) => {
  try {
    const resp = await blogApi.delArticle({
      id: id,
    })
    console.log(resp)
    resetSearch()
  } catch (error: any) {
    message.error(error)
  }
}

// 执行搜索
const handleSearch = (searchValue: string) => {
  opts.Options = replaceOps(opts.Options, {
    key: lbblog.GetArticleListReq_ListOption.ListOptionLikeTitle,
    value: searchValue
  })
  getArticleList(opts)
  console.log(searchValue)
}

// 重置搜索
const resetSearch = () => {
  opts.size = defaultSize;
  opts.page = defaultPage;
  opts.skip_total = false;
  opts.Options = [
    {
      key: lb.DefaultListOption.DefaultListOptionOrderBy,
      value: lb.DefaultOrderBy.DefaultOrderByCreatedAtDesc.toString(),
    }
  ];

  searchValue.value = "";
  getArticleList(opts)
}

const categoryChange = (value: number | undefined) => {
  if (value) {
    opts.Options = replaceOps(opts.Options, {
      key: lbblog.GetArticleListReq_ListOption.ListOptionCategoryId,
      value: String(value)
    })
  } else {
    opts.Options = replaceOps(opts.Options, {
      key: lbblog.GetArticleListReq_ListOption.ListOptionCategoryId,
      value: '',
    })
  }

  getArticleList(opts)
}

</script>

