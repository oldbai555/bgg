<template>
  <a-layout class="container">
    <Nav></Nav>

    <a-layout>
      <a-layout-header class="headerBtn">
        <Header></Header>
      </a-layout-header>
      <a-layout-content>
        <router-view :key="$route.fullPath"></router-view>
      </a-layout-content>
      <a-layout-footer class="footer">
        <Footer></Footer>
      </a-layout-footer>
    </a-layout>
  </a-layout>

</template>

<script>
import Nav from '../components/admin/nav'
import Footer from '../components/admin/footer'
import Header from '../components/admin/header'
export default {
  components: { Nav, Footer, Header },
}
</script>

<style scoped>
.container {
  height: 100%;
}
.headerBtn {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
</style>
