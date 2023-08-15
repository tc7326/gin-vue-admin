import service from '@/utils/request'
// @Summary 用户注册
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /register [post]
export const userRegister = (data) => {
  return service({
    url: '/register/register',
    method: 'post',
    data: data,
  })
}
