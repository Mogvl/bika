import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: '/home',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/home',
      name: 'home',
      component: () => import('@/views/Home.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/category',
      name: 'category',
      component: () => import('@/views/Category.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('@/views/Search.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/comic/:id',
      name: 'comic-detail',
      component: () => import('@/views/ComicDetail.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/reader/:id/:epsId',
      name: 'reader',
      component: () => import('@/views/Reader.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/leaderboard',
      name: 'leaderboard',
      component: () => import('@/views/Leaderboard.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/favourites',
      name: 'favourites',
      component: () => import('@/views/Favourites.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// 路由守卫 - 未登录跳转登录页
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
