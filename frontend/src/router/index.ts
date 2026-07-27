import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
    { path: '/register', name: 'register', component: () => import('@/views/Register.vue') },
    { path: '/home', name: 'home', component: () => import('@/views/Home.vue'), meta: { requiresAuth: true } },
    { path: '/category', name: 'category', component: () => import('@/views/Category.vue'), meta: { requiresAuth: true } },
    { path: '/search', name: 'search', component: () => import('@/views/Search.vue'), meta: { requiresAuth: true } },
    { path: '/comic/:id', name: 'comic-detail', component: () => import('@/views/ComicDetail.vue'), meta: { requiresAuth: true } },
    { path: '/reader/:id/:epsId', name: 'reader', component: () => import('@/views/Reader.vue'), meta: { requiresAuth: true } },
    { path: '/leaderboard', name: 'leaderboard', component: () => import('@/views/Leaderboard.vue'), meta: { requiresAuth: true } },
    { path: '/favourites', name: 'favourites', component: () => import('@/views/Favourites.vue'), meta: { requiresAuth: true } },
    { path: '/games', name: 'games', component: () => import('@/views/GameList.vue'), meta: { requiresAuth: true } },
    { path: '/game/:id', name: 'game-detail', component: () => import('@/views/GameDetail.vue'), meta: { requiresAuth: true } },
    { path: '/profile', name: 'profile', component: () => import('@/views/Profile.vue'), meta: { requiresAuth: true } },
    { path: '/settings', name: 'settings', component: () => import('@/views/Settings.vue'), meta: { requiresAuth: true } },
    { path: '/history', name: 'history', component: () => import('@/views/History.vue'), meta: { requiresAuth: true } },
    { path: '/downloads', name: 'downloads', component: () => import('@/views/Download.vue'), meta: { requiresAuth: true } },
  ],
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
