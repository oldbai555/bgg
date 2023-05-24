<template>
  <div class="footer d-flex align-items-center">
    <a target="_blank" href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=45130202000252"
       style="display:inline-block;text-decoration:none;height:20px;line-height:20px;">
      <img
          src="@/assets/image/beian.png"
          style="float:left;"/>
      <p
          style="float:left;height:20px;line-height:20px;margin: 0px 0px 0px 5px; color:#939393;">桂公网安备
        45130202000252号</p></a>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";

export default defineComponent({})
</script>

<style scoped>
.footer {
  background-color: #eee;
  text-align: center;
  height: 100%;
}

.footer span {
  height: 100%;
  font-size: 18px;
}
</style>