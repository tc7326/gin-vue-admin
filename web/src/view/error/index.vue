<template>
  <div>
    <div class="w-full h-screen bg-gray-50 flex items-center justify-center">
      <div class="flex flex-col items-center text-2xl gap-4">
        <img src="../../assets/Charged_Creeper_Animated.webp">
        <p>404 NOT FOUND</p>
        <p style="font-size:18px;line-height:40px;">页面被苦力怕炸没了</p>
        <p>↓</p>
        <el-button @click="toDashboard">返回首页</el-button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Error'
}
</script>

<script setup>
import {useUserStore} from "@/pinia/modules/user";
import {useRouter} from "vue-router";

const userStore = useUserStore()
const router = useRouter()
const toDashboard = () =>{
  router.push({name:userStore.userInfo.authority.defaultRouter})
}
</script>
