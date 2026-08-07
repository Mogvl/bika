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
/* ==================== 柔雾粉紫 × MIUIX 设计令牌 ==================== */
:root {
  /* 主色：粉→紫渐变家族（保留 MIUIX 前的粉紫基底） */
  --primary: #a475d0;        /* 紫罗兰主色 */
  --primary-dark: #8b5fc0;
  --primary-soft: #efe5fb;   /* 浅紫（毛玻璃底） */
  --accent: #e6a8d0;         /* 柔和粉 */
  --accent-soft: #fbeef6;

  /* 底色：带柔雾渐变的冷白，MIUIX 玻璃感来自半透明层叠 */
  --bg: linear-gradient(160deg, #f6f0fb 0%, #fbf7fc 40%, #f3eef9 100%);
  --bg-solid: #f8f4fb;
  --bg-soft: rgba(164, 117, 208, 0.07);
  --bg-card: rgba(255, 255, 255, 0.62);
  --bg-glass: rgba(255, 255, 255, 0.5);
  --bg-elev: rgba(255, 255, 255, 0.8);

  /* 文字：暖灰紫，MIUIX 强对比标题 */
  --text: #332b4d;
  --text-strong: #221c33;
  --text-secondary: #83789b;
  --text-muted: #aba1bf;

  /* 边框：MIUIX 风格的极细高光玻璃边 */
  --border: rgba(164, 117, 208, 0.1);
  --border-strong: rgba(164, 117, 208, 0.2);
  --glass-highlight: inset 0 1px 0 rgba(255, 255, 255, 0.9);

  /* 阴影：MIUIX 多层级错落散射，柔和无刺眼描边 */
  --shadow-sm: 0 2px 8px rgba(90, 60, 120, 0.05), 0 6px 20px rgba(90, 60, 120, 0.05);
  --shadow: 0 6px 18px rgba(90, 60, 120, 0.06), 0 18px 48px rgba(164, 117, 208, 0.1);
  --shadow-glow: 0 0 0 1px rgba(164, 117, 208, 0.05), 0 10px 34px rgba(164, 117, 208, 0.14);

  /* 圆角：MIUIX 大圆润胶囊感 */
  --radius: 20px;
  --radius-sm: 16px;
  --radius-lg: 30px;
  --radius-pill: 999px;

  --safe-top: env(safe-area-inset-top, 0px);
  --font-display: 'Plus Jakarta Sans', 'Noto Sans SC', -apple-system, sans-serif;
  --font-body: 'Noto Sans SC', 'Plus Jakarta Sans', -apple-system, sans-serif;

  /* 状态色（低饱和） */
  --success: #58c98f;
  --danger: #ec8fa0;
  --warn: #e7b877;
  --info: #8fb8ef;

  /* MIUIX 毛玻璃统一滤镜 */
  --glass-blur: blur(28px) saturate(150%);
}

html, body { background: var(--bg-bg, #f6f0fb); }
html, body, #app-root { background: transparent; }

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: var(--font-body);
  background: linear-gradient(160deg, #f6f0fb 0%, #fbf7fc 40%, #f3eef9 100%);
  background-attachment: fixed;
  color: var(--text);
  line-height: 1.7;
  -webkit-font-smoothing: antialiased;
}

a { color: inherit; text-decoration: none; }

/* ==================== 顶栏 ==================== */
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(250, 247, 252, 0.6);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border-bottom: 1px solid rgba(164, 117, 208, 0.08);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9), 0 6px 24px rgba(90, 60, 120, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  z-index: 100;
}

.header-left { display: flex; align-items: center; gap: 14px; }
.header-right { display: flex; align-items: center; gap: 22px; }

.menu-btn, .back-btn {
  background: none;
  border: none;
  font-size: 19px;
  cursor: pointer;
  padding: 6px;
  color: var(--text-secondary);
  transition: color 0.25s, transform 0.25s;
}
.menu-btn, .back-btn {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-pill);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.55);
  box-shadow: var(--shadow-sm);
}
.menu-btn:hover, .back-btn:hover { color: var(--primary); }
.back-btn:hover { transform: translateX(-2px); }

.logo {
  font-family: var(--font-display);
  font-size: 19px;
  font-weight: 700;
  letter-spacing: -0.02em;
  background: linear-gradient(120deg, #b98ce8 0%, #a475d0 50%, #e08bb8 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav-link {
  font-size: 17px;
  padding: 6px;
  border-radius: var(--radius-pill);
  transition: color 0.25s, transform 0.25s;
  color: var(--text-secondary);
}
.nav-link {
  width: 40px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.55);
  box-shadow: var(--shadow-sm);
}
.nav-link:hover { color: var(--primary); transform: translateY(-1px); }

.user-name {
  font-size: 13px;
  color: var(--text-muted);
  cursor: pointer;
  max-width: 90px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color 0.25s;
}
.user-name:hover { color: var(--primary); }

.app-main {
  padding-top: 64px;
  min-height: 100vh;
}
.app-main.reader-mode { padding-top: 0; }

/* ==================== 菜单 ==================== */
.menu-overlay, .user-menu-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(60, 40, 80, 0.28);
  backdrop-filter: blur(3px);
  z-index: 200;
  animation: fadeIn 0.25s ease;
}

.side-menu {
  position: fixed;
  top: 0; left: 0; bottom: 0;
  width: 268px;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  z-index: 201;
  box-shadow: var(--glass-highlight), 16px 0 56px rgba(90, 60, 120, 0.12);
  display: flex;
  flex-direction: column;
  animation: slideIn 0.32s cubic-bezier(0.22, 1, 0.36, 1);
}

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideIn { from { transform: translateX(-100%); } to { transform: translateX(0); } }

.menu-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28px 24px 20px;
  border-bottom: 1px solid var(--border);
}
.menu-header h3 {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(120deg, #b98ce8, #e08bb8);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.menu-header button {
  background: none;
  border: none;
  font-size: 17px;
  cursor: pointer;
  color: var(--text-muted);
}
.menu-header button:hover { color: var(--text); }

.menu-items {
  flex: 1;
  padding: 16px 14px;
  overflow-y: auto;
}
.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 13px 16px;
  font-size: 14.5px;
  color: var(--text);
  border-radius: var(--radius-sm);
  transition: background 0.2s, color 0.2s;
}
.menu-item:hover { background: rgba(164, 117, 208, 0.08); color: var(--primary-dark); }
.menu-item.logout { color: var(--danger); }
.menu-item.logout:hover { background: rgba(236, 143, 160, 0.12); }

.menu-items hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 14px 16px;
}

.user-menu {
  position: fixed;
  top: 70px;
  right: 18px;
  width: 200px;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid var(--border);
  border-radius: 22px;
  box-shadow: var(--glass-highlight), var(--shadow);
  z-index: 201;
  padding: 14px;
  animation: fadeIn 0.25s ease;
}
.user-info { margin-bottom: 10px; padding-bottom: 12px; border-bottom: 1px solid var(--border); }
.user-level { font-size: 12px; color: var(--primary); }
.user-name { font-size: 14px; font-weight: 600; color: var(--text); }
.user-menu a {
  display: block;
  padding: 10px 6px;
  font-size: 14px;
  color: var(--danger);
  transition: color 0.2s;
}
.user-menu a:hover { color: var(--primary); }

/* ==================== 通用布局 ==================== */
.page-container {
  max-width: 760px;
  margin: 0 auto;
  padding: 28px 20px 48px;
  animation: fadeUp 0.4s ease;
}
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.section-title {
  font-family: var(--font-display);
  font-size: 26px;
  font-weight: 800;
  letter-spacing: -0.03em;
  margin-bottom: 20px;
  padding: 0 4px;
  color: var(--text-strong);
}

.comic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 20px 16px;
}

.comic-card {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: var(--radius-sm);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: transform 0.3s cubic-bezier(0.22, 1, 0.36, 1), box-shadow 0.3s;
  cursor: pointer;
  position: relative;
}
.comic-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow);
}

.comic-cover {
  width: 100%;
  aspect-ratio: 3/4;
  object-fit: cover;
  background: rgba(164, 117, 208, 0.08);
  display: block;
}

.comic-info {
  padding: 10px 12px 12px;
}
.comic-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}
.comic-author {
  font-size: 11px;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ==================== 加载 / 空 / 错误 ==================== */
.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 64px;
  color: var(--text-muted);
  font-size: 14px;
}
.loading::after {
  content: '';
  width: 20px;
  height: 20px;
  border: 2px solid var(--border-strong);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-left: 10px;
}
@keyframes spin { to { transform: rotate(360deg); } }

.error-msg {
  text-align: center;
  padding: 64px 20px;
  color: var(--danger);
}
.empty-state {
  text-align: center;
  padding: 96px 20px;
  color: var(--text-muted);
}

/* ==================== 按钮 ==================== */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 28px;
  border: none;
  border-radius: var(--radius-pill);
  font-size: 14.5px;
  font-weight: 600;
  font-family: var(--font-body);
  color: var(--text);
  cursor: pointer;
  background: rgba(255, 255, 255, 0.6);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  transition: background 0.25s, color 0.25s, transform 0.15s, box-shadow 0.25s;
}
.btn:hover { background: var(--primary-soft); color: var(--primary-dark); }
.btn:active { transform: scale(0.97); }

.btn-primary {
  background: linear-gradient(120deg, #b98ce8 0%, #a475d0 60%, #e08bb8 120%);
  color: #fff;
  box-shadow: 0 8px 24px rgba(164, 117, 208, 0.28), inset 0 1px 0 rgba(255, 255, 255, 0.3);
}
.btn-primary:hover {
  background: linear-gradient(120deg, #a475d0 0%, #b98ce8 60%, #e08bb8 120%);
  color: #fff;
  box-shadow: 0 10px 30px rgba(164, 117, 208, 0.36), inset 0 1px 0 rgba(255, 255, 255, 0.3);
}
.btn-primary:disabled { opacity: 0.45; cursor: not-allowed; transform: none; }

/* ==================== 输入 ==================== */
.input-field {
  width: 100%;
  padding: 14px 20px;
  border: 1px solid var(--border);
  border-radius: var(--radius-pill);
  font-size: 15px;
  font-family: var(--font-body);
  outline: none;
  background: rgba(255, 255, 255, 0.6);
  color: var(--text);
  transition: border-color 0.25s, box-shadow 0.25s, background 0.25s;
}
.input-field::placeholder { color: var(--text-muted); }
.input-field:focus {
  border-color: var(--primary);
  background: rgba(255, 255, 255, 0.85);
  box-shadow: 0 0 0 4px rgba(164, 117, 208, 0.1);
}
select.input-field option { background: var(--bg-solid); color: var(--text); }

/* ==================== 加载更多 ==================== */
.page-load-more {
  display: flex;
  justify-content: center;
  padding: 32px;
}
.load-more-btn {
  padding: 11px 40px;
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid var(--border);
  border-radius: var(--radius-pill);
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(10px);
  transition: all 0.25s;
}
.load-more-btn:hover { border-color: var(--primary); color: var(--primary); background: var(--primary-soft); }

/* ==================== 分类标签 ==================== */
.category-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 12px 4px;
}
.category-tag {
  padding: 9px 20px;
  border-radius: var(--radius-pill);
  font-size: 13px;
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid var(--border);
  color: var(--text-secondary);
  cursor: pointer;
  backdrop-filter: blur(10px);
  transition: all 0.25s;
  white-space: nowrap;
}
.category-tag:hover { border-color: var(--primary); color: var(--primary); }
.category-tag.active {
  background: linear-gradient(120deg, #b98ce8, #e08bb8);
  color: #fff;
  border-color: transparent;
  box-shadow: 0 4px 14px rgba(164, 117, 208, 0.2);
}

/* ==================== 搜索 ==================== */
.search-box {
  position: relative;
  padding: 4px 0 16px;
}
.search-box input {
  width: 100%;
  padding: 14px 20px 14px 48px;
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: var(--radius-pill);
  background: rgba(255, 255, 255, 0.6);
  color: var(--text);
  font-size: 15px;
  font-family: var(--font-body);
  outline: none;
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  transition: border-color 0.25s, box-shadow 0.25s;
}
.search-box input:focus { border-color: var(--primary); box-shadow: 0 0 0 4px rgba(164, 117, 208, 0.1); }
.search-box input::placeholder { color: var(--text-muted); }
.search-box .search-icon {
  position: absolute;
  left: 30px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  font-size: 15px;
}

/* ==================== Tabs ==================== */
.tabs {
  display: flex;
  gap: 8px;
  padding: 10px 20px 0;
  overflow-x: auto;
  border-bottom: 1px solid var(--border);
}
.tab-item {
  flex: 1;
  min-width: 60px;
  padding: 14px 16px;
  text-align: center;
  font-size: 14px;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.25s;
  white-space: nowrap;
  position: relative;
}
.tab-item:hover { color: var(--text); }
.tab-item.active { color: var(--primary-dark); font-weight: 700; }
.tab-item.active::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: -1px;
  transform: translateX(-50%);
  width: 40px;
  height: 3px;
  border-radius: var(--radius-pill);
  background: linear-gradient(90deg, #b98ce8, #e08bb8);
}

/* ==================== 模态框 ==================== */
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(60, 40, 80, 0.3);
  backdrop-filter: blur(6px);
  display: flex; align-items: center; justify-content: center;
  z-index: 300;
  animation: fadeIn 0.25s ease;
}
.modal {
  background: rgba(255, 255, 255, 0.74);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: var(--radius-lg);
  padding: 32px 30px;
  width: 90%; max-width: 400px;
  display: flex; flex-direction: column; gap: 16px;
  box-shadow: var(--glass-highlight), var(--shadow);
  animation: popIn 0.3s cubic-bezier(0.22, 1, 0.36, 1);
}
@keyframes popIn {
  from { opacity: 0; transform: scale(0.96) translateY(6px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}
.modal h3 { margin-bottom: 4px; font-family: var(--font-display); font-size: 17px; }
.modal .error-tip { font-size: 13px; }

/* ==================== 空/错误状态 ==================== */
.empty-state p:first-child { font-size: 15px; margin-bottom: 8px; }
.empty-tip { font-size: 13px; color: var(--text-muted); margin-top: 8px; }
</style>
