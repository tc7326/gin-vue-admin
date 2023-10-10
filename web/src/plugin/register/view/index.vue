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
              <p class="text-center text-4xl font-bold">账号注册</p>
              <p class="text-center text-sm font-normal text-gray-500 mt-2.5"></p>
            </div>

            <!-- 注册的表单 -->
            <el-form ref="loginForm" :model="loginFormData" :rules="rules" :validate-on-rule-change="false"
              @keyup.enter="submitReg">
              <!-- 账号 keyup.enter 是 按回车后执行的fun -->
              <el-form-item prop="username" class="mb-6">
                <el-input v-model="loginFormData.username" size="large" placeholder="用户名" suffix-icon="user" />
              </el-form-item>
              <!-- 密码 -->
              <el-form-item prop="password" class="mb-6">
                <el-input v-model="loginFormData.password" show-password size="large" type="password" placeholder="密码" />
              </el-form-item>
              <!-- 加入邮箱输入框 -->
              <el-form-item prop="email" class="mb-6">
                <el-input v-model="loginFormData.email" size="large" placeholder="QQ邮箱" type="email" />
              </el-form-item>
              <!-- 验证码 -->
              <el-form-item v-if="loginFormData.openCaptcha" prop="captcha" class="mb-6">
                <div class="flex w-full justify-between">
                  <!-- 验证码输入框 -->
                  <el-input v-model="loginFormData.captcha" placeholder="验证码" size="large" class="flex-1 mr-5" />

                  <!-- 获取验证码按钮 -->
                  <el-button class="w-2/5 h-11 rounded" type="primary" :disabled="time > 0" @click="submitEmailCaptcha">{{
                    time > 0 ? `(${time}s)后重新获取` : '获取验证码' }}</el-button>
                </div>
              </el-form-item>
              <el-form-item class="mb-6">
                <!-- 注册按钮 -->
                <el-button class="shadow shadow-blue-600 h-11 w-full" type="primary" size="large"
                  @click="submitReg">注&nbsp;册</el-button>
              </el-form-item>
              <el-form-item class="mb-6">
                <el-button class="h-11 w-full" size="large" @click="goLoginPage">有账号，去登录</el-button>
              </el-form-item>
            </el-form>

          </div>
        </div>
      </div>
      <div class="hidden md:block w-1/2 float-right text-center">
        <img class="h-full w-1/2" src="@/assets/ic_piglin_loading.gif" alt="banner">
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
import { sendEmailCaptcha } from '@/api/user'
import BottomInfo from '@/view/layout/bottomInfo/bottomInfo.vue'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
const router = useRouter()

// 验证函数
const checkUsername = (rule, value, callback) => {
  //用户名
  const regName = /^[a-zA-Z0-9_-]{3,16}$/;
  if (!regName.test(value)) {
    return callback(new Error('用户名仅能包含字母、数字和下划线'))
  } else {
    callback()
  }
}

const checkPassword = (rule, value, callback) => {
  //密码规则匹配 密码只做长度校验 反正要hash
  if (value.length < 6) {
    return callback(new Error('密码最少输入6位'))
  } else {
    callback()
  }
}

//校验邮箱格式 这里由于MC玩家 大多交流平台是qq群 所以这里限制QQ邮箱 也能在一定程度上防熊
const checkEmail = (rule, value, callback) => {
  const regEmail = /^[a-zA-Z0-9_-]{6,20}@qq\.com$/;
  if (!regEmail.test(value)) {
    return callback(new Error('请输入正确的QQ邮箱'))
  } else {
    callback()
  }
}
//校验验证码 填写了验证码才能提交注册
const checkEmailCaptcha = (rule, value, callback) => {
  const regCaptcha = /^[0-9]{6}$/;
  if (!regCaptcha.test(value)) {
    return callback(new Error('验证码格式不正确'))
  } else {
    callback()
  }
}

// 定义登录表单
const loginForm = ref(null)

//定义 注册 表单结构体
const loginFormData = reactive({
  username: '',
  password: '',
  email: '',//邮箱
  captcha: '',//用户输入的验证码
  captchaId: '',
  openCaptcha: true,//默认显示验证码
})

//注册按钮提交规则
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  email: [{ validator: checkEmail, trigger: 'blur' }],
  captcha: [{ validator: checkEmailCaptcha, trigger: 'blur' }],
})

// 用户数据
const userStore = useUserStore()

// 去登录页
const goLoginPage = () => {
  router.push({ name: 'Login', replace: true })
}


// 玩家注册异步请求
const register = async () => {

  //调用具体的注册请求API
  return await userStore.EmailRegister(loginFormData)
  
}


// 玩家注册提交
const submitReg = () => {
  loginForm.value.validate(async (v) => {
    if (v) {
      //注册请求
      const flag = await register()


    } else {

      ElMessage({
        type: 'error',
        message: '请正确填写注册信息',
        showClose: true,
      })

      return false
    }
  })
}

// 发送验证码
const submitEmailCaptcha = () => {

  //validateField 验证单个字段 https://blog.csdn.net/Alan_ran/article/details/125336443
  loginForm.value.validateField(['username', 'password', 'email'], async (v) => {
    if (v) {

      //先按钮倒计时 不管发没发成功
      getCode()

      //数据校验通过 发送验证码
      const flag = await sendEmail()

    } else {

      //数据校验失败 提示
      ElMessage({
        type: 'error',
        message: '请正确填写注册信息',
        showClose: true,
      })

      return false
    }
  })
}

// 玩家发送邮箱验证码异步请求
const sendEmail = async () => {

  const res = await sendEmailCaptcha(loginFormData)
  if (res.code === 0) {

    ElMessage.success('发送成功,请查收')

    //填充后台返回的验证码缓存的key数据 注册接口后台要该用值获取缓存校验
    loginFormData.captchaId = res.data.captchaId

    return true

  } else {
    ElMessage.error('发送失败,请稍后重试')

  }
}

//倒计时时间
const time = ref(0)

//验证码倒计时
const getCode = async () => {
  time.value = 120
  let timer = setInterval(() => {
    time.value--
    if (time.value <= 0) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

</script>
