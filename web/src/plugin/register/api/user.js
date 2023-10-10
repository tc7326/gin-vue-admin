import service from '@/utils/request'

// @Summary 用户邮箱验证注册
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /register/user [post]
export const userEmailRegister = (data) => {
  return service({
    url: '/register/user',
    method: 'post',
    data: data,
  })
}
