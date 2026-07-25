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
    },
    {
      path: '/category',
      name: 'category',
      component: () => import('@/views/Category.vue'),
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('@/views/Search.vue'),
    },
    {
      path: '/comic/:id',
      name: 'comic-detail',
      component: () => import('@/views/ComicDetail.vue'),
    },
    {
      path: '/reader/:id/:epsId',
      name: 'reader',
      component: () => import('@/views/Reader.vue'),
    },
    {
      path: '/leaderboard',
      name: 'leaderboard',
      component: () => import('@/views/Leaderboard.vue'),
    },
    {
      path: '/favourites',
      name: 'favourites',
      component: () => import('@/views/Favourites.vue'),
    },
  ],
})

export default router
