import Vue from 'vue'
import VueRouter from 'vue-router'
import Layout from '@/views/Layout.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import ImageLibrary from '@/views/ImageLibrary.vue'
import AudioLibrary from '@/views/AudioLibrary.vue'
import TaskManagement from '@/views/TaskManagement.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/',
    component: Layout,
    redirect: '/images',
    children: [
      {
        path: '/images',
        name: 'ImageLibrary',
        component: ImageLibrary
      },
      {
        path: '/audios',
        name: 'AudioLibrary',
        component: AudioLibrary
      },
      {
        path: '/tasks',
        name: 'TaskManagement',
        component: TaskManagement
      }
    ]
  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
});

// 添加路由守卫
router.beforeEach((to, from, next) => {
  // 白名单：无需 token 的页面
  const publicPages = ['/login', '/register'];
  if (publicPages.includes(to.path)) {
    next();
    return;
  }

  // 获取token
  const token = localStorage.getItem('emo-token');

  // 如果未登录且访问非白名单页面，则跳转到登录页
  if (!token && !publicPages.includes(to.path)) {
    next('/login');
    return;
  }

  // 如果已登录且尝试访问登录页，则重定向到首页
  if (token && to.path === '/login') {
    next('/');
    return;
  }

  // 其他情况正常跳转
  next();
});

export default router