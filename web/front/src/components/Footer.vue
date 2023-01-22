<template>
  <v-footer padless color="indigo darken-2">
    <v-row justify="center" no-gutters>
      <v-col class="py-4 text-center white--text" cols="12">{{ new Date().getFullYear() }} - 老白的博客</v-col>

      <div class="text-center white--text">
        <a class="text-center white--text" target="_blank"
           href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=45130202000252"
           style="display:inline-block;text-decoration:none;height:20px;line-height:20px;"><img
            src="../assets/备案图标.png"
            style="float:left;"/>
          <p
              style="float:left;height:20px;line-height:20px;margin: 0px 0px 0px 5px; color:#939393;">桂公网安备
            45130202000252号</p></a>
      </div>
    </v-row>
  </v-footer>
</template>

<script>
export default {
  data() {
    return {
      icp_record: ''
    }
  },
  mounted() {
    this.$root.$on('msg', (msg) => {
      this.icp_record = msg
    })
  }
}
</script>

<style>
</style>