<template>
  <div id="app-root">
    <header v-if="auth.isLoggedIn && !isReaderPage" class="app-header">
      <div class="header-left">
        <button v-if="showBack" class="back-btn" @click="goBack">←</button>
        <button class="menu-btn" @click="showMenu = !showMenu">☰</button>
        <router-link to="/home" class="logo">PicACG</router-link>
      </div>
      <div class="header-right">
        <router-link to="/search" class="nav-link" aria-label="搜索">🔍</router-link>
        <router-link to="/games" class="nav-link" aria-label="游戏">🎮</router-link>
        <router-link to="/leaderboard" class="nav-link" aria-label="排行榜">🏆</router-link>
        <router-link to="/favourites" class="nav-link" aria-label="收藏">❤️</router-link>
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
        <router-link to="/local" @click="showMenu = false" class="menu-item">📚 本地库</router-link>
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

const isReaderPage = computed(() => route.name === 'reader' || route.name === 'local-reader')

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
  --primary: #ff5fa8;
  --primary-dark: #f0449a;
  --accent: #ff8fc2;
  --bg: #1a0f18;
  --bg-soft: #241722;
  --bg-card: #2c1c29;
  --bg-elev: #372333;
  --text: #ffeef6;
  --text-secondary: #d9a8c2;
  --text-muted: #a97a94;
  --border: rgba(255, 143, 194, 0.14);
  --border-strong: rgba(255, 143, 194, 0.28);
  --shadow: 0 8px 30px rgba(0, 0, 0, 0.4);
  --shadow-glow: 0 0 0 1px rgba(255, 95, 168, 0.3), 0 8px 32px rgba(255, 95, 168, 0.22);
  --radius: 14px;
  --radius-sm: 10px;
  --radius-lg: 22px;
  --safe-top: env(safe-area-inset-top, 0px);
  --font-display: 'Plus Jakarta Sans', 'Noto Sans SC', -apple-system, sans-serif;
  --font-body: 'Noto Sans SC', 'Plus Jakarta Sans', -apple-system, sans-serif;
  --success: #51d8a6;
  --danger: #ff6b7d;
  --warn: #f5b754;
  --info: #6db7ff;
}

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: var(--font-body);
  background: var(--bg);
  color: var(--text);
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
}

a { color: inherit; text-decoration: none; }

.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 58px;
  background: rgba(26, 15, 24, 0.82);
  backdrop-filter: blur(18px) saturate(140%);
  -webkit-backdrop-filter: blur(18px) saturate(140%);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 100;
}

.header-left { display: flex; align-items: center; gap: 12px; }
.header-right { display: flex; align-items: center; gap: 18px; }

.menu-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
  color: var(--text);
  transition: color 0.2s, transform 0.2s;
}
.menu-btn:hover { color: var(--primary); }

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
  color: var(--text);
  margin-right: 2px;
  transition: color 0.2s, transform 0.2s;
}
.back-btn:hover { color: var(--primary); transform: translateX(-2px); }

.logo {
  font-family: var(--font-display);
  font-size: 19px;
  font-weight: 700;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ff9acb 0%, #ff5fa8 55%, #ff3d94 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav-link {
  font-size: 17px;
  padding: 6px;
  border-radius: 10px;
  transition: transform 0.2s, background 0.2s;
  opacity: 0.85;
}
.nav-link:hover { opacity: 1; transform: translateY(-1px); background: var(--bg-elev); }

.user-name {
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  max-width: 90px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color 0.2s;
}
.user-name:hover { color: var(--primary); }

.app-main {
  padding-top: 58px;
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
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(2px);
  z-index: 200;
  animation: fadeIn 0.2s ease;
}

.side-menu {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 264px;
  background: var(--bg-card);
  z-index: 201;
  box-shadow: 4px 0 32px rgba(0,0,0,0.5);
  display: flex;
  flex-direction: column;
  animation: slideIn 0.25s cubic-bezier(0.22, 1, 0.36, 1);
}

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideIn { from { transform: translateX(-100%); } to { transform: translateX(0); } }

.menu-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 22px 18px;
  border-bottom: 1px solid var(--border);
}

.menu-header h3 {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #ff9acb, #ff5fa8);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.menu-header button {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: var(--text-secondary);
}
.menu-header button:hover { color: var(--text); }

.menu-items {
  flex: 1;
  padding: 10px 10px;
  overflow-y: auto;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 13px 14px;
  font-size: 14.5px;
  color: var(--text);
  border-radius: 12px;
  transition: background 0.2s, color 0.2s, transform 0.15s;
}

.menu-item:hover { background: var(--bg-elev); color: var(--primary); }
.menu-item.logout { color: #ff6b81; }
.menu-item.logout:hover { background: rgba(255, 107, 129, 0.12); }

.menu-items hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 10px 12px;
}

.user-menu {
  position: fixed;
  top: 62px;
  right: 12px;
  width: 190px;
  background: var(--bg-elev);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  z-index: 201;
  padding: 12px;
  animation: fadeIn 0.2s ease;
}

.user-info { margin-bottom: 8px; padding-bottom: 10px; border-bottom: 1px solid var(--border); }
.user-level { font-size: 12px; color: var(--primary); }
.user-name { font-size: 14px; font-weight: 600; color: var(--text); }

.user-menu a {
  display: block;
  padding: 9px 4px;
  font-size: 14px;
  color: #ff6b81;
  transition: color 0.2s;
}
.user-menu a:hover { color: #ff8fc2; }

/* ==================== 通用组件样式 ==================== */
.page-container {
  max-width: 720px;
  margin: 0 auto;
  padding: 14px;
  animation: fadeUp 0.3s ease;
}

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.section-title {
  font-family: var(--font-display);
  font-size: 19px;
  font-weight: 600;
  margin-bottom: 14px;
  padding: 0 4px;
  color: var(--text);
  letter-spacing: -0.01em;
}

.comic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(132px, 1fr));
  gap: 14px;
}

.comic-card {
  background: var(--bg-card);
  border-radius: var(--radius-sm);
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.25);
  transition: transform 0.22s cubic-bezier(0.22, 1, 0.36, 1), box-shadow 0.22s;
  cursor: pointer;
  position: relative;
}

.comic-card:hover {
  transform: translateY(-4px) scale(1.01);
  box-shadow: 0 12px 28px rgba(0,0,0,0.45), 0 0 0 1px rgba(255, 95, 168, 0.2);
}

.comic-cover {
  width: 100%;
  aspect-ratio: 3/4;
  object-fit: cover;
  background: var(--bg-soft);
  display: block;
}

.comic-info {
  padding: 8px 10px 10px;
}

.comic-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 3px;
}

.comic-author {
  font-size: 11px;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px;
  color: var(--text-secondary);
  font-size: 14px;
}

.loading::after {
  content: '';
  width: 22px;
  height: 22px;
  border: 2px solid var(--border-strong);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-left: 10px;
}

@keyframes spin { to { transform: rotate(360deg); } }

.error-msg {
  text-align: center;
  padding: 48px 20px;
  color: #ff6b81;
}

.empty-state {
  text-align: center;
  padding: 72px 20px;
  color: var(--text-secondary);
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 22px;
  border: none;
  border-radius: 12px;
  font-size: 14.5px;
  font-weight: 500;
  font-family: var(--font-body);
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.2s, background 0.2s, opacity 0.2s;
}

.btn:active { transform: scale(0.97); }

.btn-primary {
  background: linear-gradient(135deg, #ff5fa8 0%, #f0449a 100%);
  color: #fff;
  box-shadow: 0 4px 16px rgba(255, 95, 168, 0.28);
}
.btn-primary:hover { box-shadow: 0 6px 22px rgba(255, 95, 168, 0.4); transform: translateY(-1px); }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; transform: none; }

.input-field {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border);
  border-radius: 12px;
  font-size: 15px;
  font-family: var(--font-body);
  outline: none;
  background: var(--bg-soft);
  color: var(--text);
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
}

.input-field::placeholder { color: var(--text-muted); }

.input-field:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(255, 95, 168, 0.15);
}

select.input-field option { background: var(--bg-card); color: var(--text); }

.page-load-more {
  display: flex;
  justify-content: center;
  padding: 24px;
}

.load-more-btn {
  padding: 9px 32px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s;
}
.load-more-btn:hover { border-color: var(--primary); color: var(--primary); background: var(--bg-elev); }

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
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}
.category-tag:hover { border-color: var(--primary); color: var(--primary); }
.category-tag.active {
  background: linear-gradient(135deg, #ff5fa8, #f0449a);
  color: #fff;
  border-color: transparent;
  box-shadow: 0 4px 12px rgba(255, 95, 168, 0.3);
}

/* 搜索框 */
.search-box {
  position: relative;
  padding: 12px;
}

.search-box input {
  width: 100%;
  padding: 13px 16px 13px 44px;
  border: 1px solid var(--border);
  border-radius: 26px;
  background: var(--bg-card);
  color: var(--text);
  font-size: 15px;
  font-family: var(--font-body);
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.search-box input:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(255, 95, 168, 0.15); }
.search-box input::placeholder { color: var(--text-muted); }

.search-box .search-icon {
  position: absolute;
  left: 27px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
  font-size: 15px;
}

/* 标签(tabs) */
.tabs {
  display: flex;
  gap: 6px;
  padding: 8px 12px 0;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border-bottom: 1px solid var(--border);
  background: rgba(13, 13, 18, 0.6);
}

.tab-item {
  flex: 1;
  min-width: 60px;
  padding: 12px 16px;
  text-align: center;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 10px 10px 0 0;
  transition: all 0.2s;
  white-space: nowrap;
  position: relative;
}

.tab-item:hover { color: var(--text); }

.tab-item.active {
  color: var(--primary);
  font-weight: 600;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: 0;
  transform: translateX(-50%);
  width: 28px;
  height: 2px;
  border-radius: 2px;
  background: linear-gradient(90deg, #ff9acb, #ff5fa8);
}

/* 模态框 */
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 300;
  animation: fadeIn 0.2s ease;
}
.modal {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 26px;
  width: 90%; max-width: 400px;
  display: flex; flex-direction: column; gap: 14px;
  box-shadow: var(--shadow);
  animation: popIn 0.25s cubic-bezier(0.22, 1, 0.36, 1);
}
@keyframes popIn {
  from { opacity: 0; transform: scale(0.94) translateY(8px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}
.modal h3 { margin-bottom: 6px; font-family: var(--font-display); font-size: 18px; }
.modal .error-tip { font-size: 13px; }

/* 空/错误状态 */
.empty-state p:first-child { font-size: 15px; margin-bottom: 6px; }
.empty-tip { font-size: 13px; color: var(--text-muted); margin-top: 8px; }
</style>
