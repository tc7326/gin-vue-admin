import { createRouter, createWebHashHistory } from 'vue-router'

// 页面的路由 / 重定向到 /login
const routes = [{
  path: '/',
  redirect: '/login'
},
{
  path: '/init',
  name: 'Init',
  component: () => import('@/view/init/index.vue')
},
{
  path: '/login',
  name: 'Login',
  component: () => import('@/view/login/index.vue')
},
{
  path: '/register',//新增mc注册的路由
  name: 'Register',
  component: () => import('@/plugin/register/view/index.vue')
},
{
  path: '/:catchAll(.*)',
  meta: {
    closeTab: true,
  },
  component: () => import('@/view/error/index.vue')
}
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
