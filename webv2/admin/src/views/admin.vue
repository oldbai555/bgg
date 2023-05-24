<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider v-model:collapsed="collapsed" collapsible breakpoint="lg">
      <Sider/>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0">
        <Header/>
      </a-layout-header>
      <a-layout-content style="margin: 0 16px">
        <RouterView :key="$route.fullPath"/>
      </a-layout-content>
      <a-layout-footer style="text-align: center">
        <Footer/>
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue';
import Sider from '../components/admin/sider.vue'
import Footer from '../components/admin/footer.vue'
import Header from '../components/admin/header.vue'


export default defineComponent({
  components: {
    Sider,
    Footer,
    Header,
  },
  data() {
    return {
      collapsed: ref<boolean>(false),
    };
  },
});
</script>
<style>
#components-layout-demo-side .logo {
  height: 32px;
  margin: 16px;
  background: rgba(255, 255, 255, 0.3);
}

.site-layout .site-layout-background {
  background: #fff;
}

[data-theme='dark'] .site-layout .site-layout-background {
  background: #141414;
}
</style>
