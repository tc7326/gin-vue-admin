<template>
  <div id="userLayout" class="w-full h-full relative">
    <div
      class="rounded-lg flex items-center justify-evenly w-full h-full bg-white md:w-screen md:h-screen md:bg-[#7B68EE]">
      <div class="md:w-3/5 w-10/12 h-full flex items-center justify-evenly">
        <!-- 页面背景时蓝色的 这是加了一个白色的斜面 -->
        <div class="oblique h-[130%] w-3/5 bg-white transform -rotate-12 absolute -ml-52" />
        <!-- 分割斜块 -->
        <div class="z-[999] pt-12 pb-10 md:w-96 w-full  rounded-lg flex flex-col justify-between box-border">
          <div>
            <div class="flex items-center justify-center">
              <img class="w-24" :src="$GIN_VUE_ADMIN.appLogo" alt>
            </div>
            <div class="mb-9">
              <p class="text-center text-sm font-normal text-gray-500 mt-2.5">&nbsp;</p>
              <p class="text-center text-4xl font-bold">账号登录</p>
              <p class="text-center text-sm font-normal text-gray-500 mt-2.5"></p>
            </div>
            <el-form ref="loginForm" :model="loginFormData" :rules="rules" :validate-on-rule-change="false"
              @keyup.enter="submitForm">
              <!-- keyup.enter 是 按回车后执行的fun -->
              <el-form-item prop="username" class="mb-6">
                <el-input v-model="loginFormData.username" size="large" placeholder="用户名/QQ邮箱" suffix-icon="user" />
              </el-form-item>
              <el-form-item prop="password" class="mb-6">
                <el-input v-model="loginFormData.password" show-password size="large" type="password" placeholder="密码" />
              </el-form-item>
              <el-form-item v-if="loginFormData.openCaptcha" prop="captcha" class="mb-6">
                <div class="flex w-full justify-between">
                  <el-input v-model="loginFormData.captcha" placeholder="验证码" size="large" class="flex-1 mr-5" />
                  <div class="w-2/5 h-11 bg-[#c3d4f2] rounded">
                    <img v-if="picPath" class="w-full h-full" :src="picPath" alt="验证码" @click="submitVerify()">
                  </div>
                </div>
              </el-form-item>
              <el-form-item class="mb-6">
                <el-button class="shadow shadow-blue-600 h-11 w-full" type="primary" size="large"
                  @click="submitForm">登&nbsp;录</el-button>
              </el-form-item>
              <el-form-item class="mb-6">
                <el-button class="h-11 w-full" size="large" @click="GoRegister">没账号，去注册</el-button>
              </el-form-item>

              <!-- 数据库初始化 初始化完成后 屏蔽此内容 -->
              <!-- <el-form-item class="mb-6">
                <el-button class="shadow shadow-blue-600 h-11 w-full" type="primary" size="large"
                  @click="checkInit">前往初始化</el-button>
              </el-form-item> -->

            </el-form>
          </div>
        </div>
      </div>
      <div class="hidden md:block w-1/2 float-right text-center">
        
        <img class="h-full w-1/4" src="@/assets/Allay_Dancing.webp" alt="banner">
      </div>
    </div>

    <BottomInfo class="left-0 right-0 absolute bottom-3 mx-auto  w-full z-20">
      <!-- <div class="links items-center justify-center gap-2 hidden md:flex">
        <a href="http://doc.henrongyi.top/" target="_blank">
          <img src="@/assets/docs.png" class="w-8 h-8" alt="文档">
        </a>
        <a href="https://support.qq.com/product/371961" target="_blank">
          <img src="@/assets/kefu.png" class="w-8 h-8" alt="客服">
        </a>
        <a href="https://github.com/flipped-aurora/gin-vue-admin" target="_blank">
          <img src="@/assets/github.png" class="w-8 h-8" alt="github">
        </a>
        <a href="https://space.bilibili.com/322210472" target="_blank">
          <img src="@/assets/video.png" class="w-8 h-8" alt="视频站">
        </a>
      </div> -->
    </BottomInfo>
    
  </div>
</template>

<script>
export default {
  name: 'Login',
}
</script>

<script setup>
import { captcha } from '@/api/user'
import { checkDB } from '@/api/initdb'
import BottomInfo from '@/view/layout/bottomInfo/bottomInfo.vue'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
const router = useRouter()
// 验证函数
const checkUsername = (rule, value, callback) => {
  const regName = /^[a-zA-Z0-9_-]{3,16}$/;
  if (!regName.test(value)) {
    return callback(new Error('用户名仅能包含字母、数字和下划线'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入密码'))
  } else {
    callback()
  }
}
const checkVerify = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入验证码'))
  } else {
    callback()
  }
}

// 获取验证码校验
const submitVerify = () => {

  //validateField 验证单个字段 https://blog.csdn.net/Alan_ran/article/details/125336443
  loginForm.value.validateField(['username', 'password'], async (v) => {
    if (v) {

      //数据校验通过 获取新的验证码
      loginVerify()

    }
  })
}

// 获取验证码
const loginVerify = () => {
  captcha({}).then(async (ele) => {
    rules.captcha.push({
      max: ele.data.captchaLength,
      min: ele.data.captchaLength,
      message: `请输入${ele.data.captchaLength}位验证码`,
      trigger: 'blur',
    })
    picPath.value = ele.data.picPath
    loginFormData.captchaId = ele.data.captchaId
    loginFormData.openCaptcha = ele.data.openCaptcha
  })
}
//进入页面时 先获取验证码
loginVerify()

// 登录相关操作
const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: '',
  password: '',
  captcha: '',
  captchaId: '',
  openCaptcha: false,
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [{ validator: checkVerify, trigger: 'blur' },],
})

// 用户数据
const userStore = useUserStore()
// 异步登录请求？包括保存数据等操作？
const login = async () => {
  return await userStore.LoginIn(loginFormData)
}
// 登录提交
const submitForm = () => {
  loginForm.value.validate(async (v) => {
    if (v) {
      const flag = await login()
      if (!flag) {
        loginVerify()
      }
    } else {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      })
      return false
    }
  })
}

// 去注册页
const GoRegister = () => {
  //导航到注册页面
  // router.push({ name: 'Register', replace: true })//如果携带 replace 则浏览器的返回按钮会返回到选项卡
  router.push({ name: 'Register' })
}

// 跳转初始化
const checkInit = async () => {
  const res = await checkDB()
  if (res.code === 0) {
    if (res.data?.needInit) {
      userStore.NeedInit()
      router.push({ name: 'Init' })
    } else {
      ElMessage({
        type: 'info',
        message: '已配置数据库信息，无法初始化',
      })
    }
  }
}

</script>
