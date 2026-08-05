<template>
  <div id="app-root">
    <header v-if="auth.isLoggedIn && !isReaderPage" class="app-header">
      <div class="header-left">
        <button v-if="showBack" class="back-btn" @click="goBack">←</button>
        <button class="menu-btn" @click="showMenu = !showMenu">☰</button>
        <router-link to="/home" class="logo">PicACG</router-link>
      </div>
      <div class="header-right">
        <router-link to="/search" class="nav-link">🔍</router-link>
        <router-link to="/games" class="nav-link">🎮</router-link>
        <router-link to="/leaderboard" class="nav-link">🏆</router-link>
        <router-link to="/favourites" class="nav-link">❤️</router-link>
        <span class="user-name" @click="showUserMenu = !showUserMenu">
          {{ auth.user?.name || '用户' }}
        </span>
      </div>
    </header>

    <!-- 侧边菜单 -->
    <div v-if="showMenu" class="menu-overlay" @click="showMenu = false"></div>
    <aside v-if="showMenu" class="side-menu">
      <div class="menu-header">
        <h3>哔咔漫画</h3>
        <button @click="showMenu = false">✕</button>
      </div>
      <nav class="menu-items">
        <router-link to="/home" @click="showMenu = false" class="menu-item">🏠 首页</router-link>
        <router-link to="/category" @click="showMenu = false" class="menu-item">📂 分类浏览</router-link>
        <router-link to="/search" @click="showMenu = false" class="menu-item">🔍 搜索</router-link>
        <router-link to="/leaderboard" @click="showMenu = false" class="menu-item">🏆 排行榜</router-link>
        <router-link to="/favourites" @click="showMenu = false" class="menu-item">❤️ 我的收藏</router-link>
        <router-link to="/games" @click="showMenu = false" class="menu-item">🎮 游戏区</router-link>
        <router-link to="/chat" @click="showMenu = false" class="menu-item">💬 聊天室</router-link>
        <router-link to="/fried" @click="showMenu = false" class="menu-item">📱 锅贴广场</router-link>
        <router-link to="/downloads" @click="showMenu = false" class="menu-item">📥 下载管理</router-link>
        <router-link to="/history" @click="showMenu = false" class="menu-item">📖 阅读历史</router-link>
        <router-link to="/my-comments" @click="showMenu = false" class="menu-item">💬 我的评论</router-link>
        <router-link to="/profile" @click="showMenu = false" class="menu-item">👤 个人资料</router-link>
        <router-link to="/settings" @click="showMenu = false" class="menu-item">⚙️ 设置</router-link>
        <router-link to="/help" @click="showMenu = false" class="menu-item">❓ 帮助</router-link>
        <hr />
        <a href="#" @click.prevent="logout" class="menu-item logout">🚪 退出登录</a>
      </nav>
    </aside>

    <main class="app-main" :class="{ 'reader-mode': isReaderPage }">
      <router-view />
    </main>

    <!-- 用户菜单下拉 -->
    <div v-if="showUserMenu" class="user-menu-overlay" @click="showUserMenu = false"></div>
    <div v-if="showUserMenu" class="user-menu">
      <div class="user-info">
        <p class="user-level">Lv.{{ auth.user?.level || '?' }}</p>
        <p class="user-name">{{ auth.user?.name || '用户' }}</p>
      </div>
      <router-link to="/profile" @click="showUserMenu = false">个人资料</router-link>
      <router-link to="/settings" @click="showUserMenu = false">设置</router-link>
      <a href="#" @click.prevent="logout">退出登录</a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
auth.init()
const route = useRoute()
const router = useRouter()
const showMenu = ref(false)
const showUserMenu = ref(false)

const isReaderPage = computed(() => route.name === 'reader')

// 二级页面显示返回按钮（详情/列表/设置等）
const backRoutes = new Set(['comic-detail', 'game-detail', 'search', 'category-comics', 'profile', 'settings', 'history', 'my-comments', 'downloads', 'chat', 'fried', 'help', 'leaderboard', 'favourites', 'games'])
const showBack = computed(() => backRoutes.has(route.name as string))

function goBack() {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/home')
  }
}

function logout() {
  auth.logout()
  showMenu.value = false
  showUserMenu.value = false
  window.location.hash = '/login'
}
</script>

<style>
:root {
  --primary: #e74c8b;
  --primary-dark: #c0397a;
  --bg: #f5f5f5;
  --bg-card: #ffffff;
  --bg-header: #ffffff;
  --text: #333333;
  --text-secondary: #666666;
  --border: #eeeeee;
  --shadow: 0 2px 8px rgba(0,0,0,0.08);
  --radius: 12px;
  --safe-top: env(safe-area-inset-top, 0px);
}

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  background: var(--bg);
  color: var(--text);
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
}

a { color: inherit; text-decoration: none; }

.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: var(--bg-header);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 100;
  box-shadow: var(--shadow);
}

.header-left { display: flex; align-items: center; gap: 12px; }
.header-right { display: flex; align-items: center; gap: 16px; }

.menu-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
  color: var(--text);
}

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
  color: var(--text);
  margin-right: 4px;
}

.back-btn:hover {
  color: var(--primary);
}

.logo {
  font-size: 18px;
  font-weight: 700;
  color: var(--primary);
}

.nav-link {
  font-size: 18px;
  padding: 4px;
}

.user-name {
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
}

.app-main {
  padding-top: 56px;
  min-height: 100vh;
}

.app-main.reader-mode {
  padding-top: 0;
}

.menu-overlay, .user-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.4);
  z-index: 200;
}

.side-menu {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 260px;
  background: var(--bg-card);
  z-index: 201;
  box-shadow: 2px 0 12px rgba(0,0,0,0.15);
  display: flex;
  flex-direction: column;
}

.menu-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 16px;
  border-bottom: 1px solid var(--border);
}

.menu-header h3 { color: var(--primary); font-size: 18px; }
.menu-header button {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: var(--text-secondary);
}

.menu-items {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
}

.menu-item {
  display: block;
  padding: 14px 20px;
  font-size: 15px;
  color: var(--text);
  transition: background 0.2s;
}

.menu-item:hover { background: var(--bg); }
.menu-item.logout { color: #e74c3c; }

.menu-items hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 8px 16px;
}

.user-menu {
  position: fixed;
  top: 56px;
  right: 12px;
  width: 180px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
  z-index: 201;
  padding: 12px;
}

.user-info { margin-bottom: 8px; padding-bottom: 8px; border-bottom: 1px solid var(--border); }
.user-level { font-size: 12px; color: var(--primary); }
.user-name { font-size: 14px; font-weight: 600; }

.user-menu a {
  display: block;
  padding: 8px 4px;
  font-size: 14px;
  color: #e74c3c;
}

/* 通用组件样式 */
.page-container {
  max-width: 100%;
  padding: 12px;
  margin: 0 auto;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
  padding: 0 4px;
}

.comic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}

.comic-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: var(--shadow);
  transition: transform 0.2s;
  cursor: pointer;
}

.comic-card:hover { transform: translateY(-2px); }

.comic-cover {
  width: 100%;
  aspect-ratio: 3/4;
  object-fit: cover;
  background: var(--bg);
  display: block;
}

.comic-info {
  padding: 8px 10px;
}

.comic-title {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.comic-author {
  font-size: 11px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
  color: var(--text-secondary);
  font-size: 14px;
}

.loading::after {
  content: '';
  width: 20px;
  height: 20px;
  border: 2px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-left: 8px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-msg {
  text-align: center;
  padding: 40px 20px;
  color: #e74c3c;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 10px 24px;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: var(--primary);
  color: white;
}

.btn-primary:hover { background: var(--primary-dark); }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

.input-field {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 15px;
  outline: none;
  transition: border-color 0.2s;
}

.input-field:focus {
  border-color: var(--primary);
}

.page-load-more {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.load-more-btn {
  padding: 8px 32px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  color: var(--text-secondary);
}

.load-more-btn:hover { border-color: var(--primary); color: var(--primary); }

/* 分类标签 */
.category-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px;
  overflow-x: auto;
}

.category-tag {
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 13px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.category-tag:hover { border-color: var(--primary); color: var(--primary); }
.category-tag.active { background: var(--primary); color: white; border-color: var(--primary); }

/* 搜索框 */
.search-box {
  position: relative;
  padding: 12px;
}

.search-box input {
  width: 100%;
  padding: 12px 16px 12px 40px;
  border: none;
  border-radius: 24px;
  background: var(--bg-card);
  font-size: 15px;
  outline: none;
  box-shadow: var(--shadow);
}

.search-box .search-icon {
  position: absolute;
  left: 24px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
}

/* 标签(tabs) */
.tabs {
  display: flex;
  gap: 0;
  border-bottom: 1px solid var(--border);
  background: var(--bg-card);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.tab-item {
  flex: 1;
  min-width: 60px;
  padding: 12px 16px;
  text-align: center;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
  white-space: nowrap;
}

.tab-item.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
  font-weight: 500;
}
</style>
