<template>
  <div>
    <div class="logo"/>
    <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleClick">
      <a-menu-item key="index">
        <pie-chart-outlined/>
        <span>仪表盘</span>
      </a-menu-item>
      <a-sub-menu>
        <template #title>
            <span>
              <read-outlined/>
              <span>博客</span>
            </span>
        </template>
        <a-menu-item key="article_list">文章管理</a-menu-item>
        <a-menu-item key="category_list">分类管理</a-menu-item>
        <a-menu-item key="comment_list">评论管理</a-menu-item>
      </a-sub-menu>

      <a-sub-menu>
        <template #title>
            <span>
              <account-book-outlined/>
              <span>记账</span>
            </span>
        </template>
        <a-menu-item key="bill_list">清单管理</a-menu-item>
        <a-menu-item key="bill_category_list">分类管理</a-menu-item>
      </a-sub-menu>

      <a-sub-menu>
        <template #title>
            <span>
              <team-outlined/>
              <span>团队</span>
            </span>
        </template>
        <a-menu-item key="user_list">成员管理</a-menu-item>
      </a-sub-menu>

      <a-sub-menu>
        <template #title>
            <span>
              <file-outlined/>
              <span>存储</span>
            </span>
        </template>
        <a-menu-item key="file_list">文件管理</a-menu-item>
      </a-sub-menu>

      <a-sub-menu>
        <template #title>
            <span>
              <desktop-outlined/>
        <span>网站</span>
            </span>
        </template>
        <a-menu-item key="profile">个人中心</a-menu-item>
        <a-menu-item key="game">测试ws</a-menu-item>
      </a-sub-menu>


    </a-menu>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from "vue";
import {
  AccountBookOutlined,
  DesktopOutlined,
  FileOutlined,
  PieChartOutlined,
  ReadOutlined,
  TeamOutlined,
  UserOutlined
} from "@ant-design/icons-vue";
import {useRouter} from "vue-router";
import type {MenuProps} from 'ant-design-vue';

const defaultSelectedKeys = ['index']

export default defineComponent({
  components: {
    PieChartOutlined,
    ReadOutlined,
    AccountBookOutlined,
    FileOutlined,
    DesktopOutlined,
    UserOutlined,
    TeamOutlined,
  },
  setup() {
    const router = useRouter()
    const selectedKeys = ref(defaultSelectedKeys)
    const handleClick: MenuProps['onClick'] = async menuInfo => {
      await router.push({
        name: String(menuInfo.key),
      })
    };

    return {
      selectedKeys,
      handleClick
    }
  },
});
</script>


<style>
.logo {
  height: 32px;
  margin: 16px;
  background: rgba(255, 255, 255, 0.3);
}
</style>