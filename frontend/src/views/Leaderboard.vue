<template>
  <div class="page-container">
    <div class="section-title">🏆 排行榜</div>

    <div class="tabs">
      <div v-for="tab in tabs" :key="tab.value" class="tab-item" :class="{ active: currentTT === tab.value }" @click="switchTab(tab.value)">
        {{ tab.label }}
      </div>
    </div>

    <!-- 漫画榜 -->
    <div v-if="currentTT !== 'knight'" class="comic-grid" style="padding-top: 12px;">
      <div v-for="(comic, idx) in comics" :key="comic._id" class="comic-card" @click="goComic(comic._id)">
        <div class="rank-badge" :class="'rank-' + (idx + 1)">{{ idx + 1 }}</div>
        <img :src="getCoverUrl(comic.thumb)" :alt="comic.title" class="comic-cover" loading="lazy" @error="handleImgError" />
        <div class="comic-info">
          <div class="comic-title">{{ comic.title }}</div>
          <div class="comic-author">❤️ {{ comic.totalLikes || comic.likesCount }}</div>
        </div>
      </div>
    </div>

    <!-- 骑士榜（用户排行） -->
    <div v-else>
      <div v-if="loading" class="loading">加载中</div>
      <div v-else class="knight-list" style="padding-top: 12px;">
        <div v-for="(user, idx) in knights" :key="user._id" class="knight-item">
          <div class="rank-badge knight-rank" :class="'rank-' + (idx + 1)">{{ idx + 1 }}</div>
          <div class="knight-avatar">
            <img v-if="user.avatar" :src="getUserAvatar(user)" alt="" />
            <span v-else>{{ user.name?.[0] || '?' }}</span>
          </div>
          <div class="knight-info">
            <div class="knight-name">
              {{ user.name || '匿名' }}
              <span v-if="user.level" class="knight-level">Lv.{{ user.level }}</span>
            </div>
            <div v-if="user.title" class="knight-title">{{ user.title }}</div>
          </div>
          <div class="knight-stat">{{ user.slogan || '' }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getLeaderboard, getKnightLeaderboard } from '@/api'
import type { Comic, LeaderboardTT } from '@/types'

const router = useRouter()
const tabs = [
  { label: '24小时', value: 'H24' as const },
  { label: '本周', value: 'D7' as const },
  { label: '本月', value: 'D30' as const },
  { label: '骑士榜', value: 'knight' as const },
]
type TabValue = LeaderboardTT | 'knight'
const currentTT = ref<TabValue>('H24')
const comics = ref<Comic[]>([])
const knights = ref<any[]>([])
const loading = ref(true)

onMounted(() => load())

async function switchTab(tt: TabValue) {
  currentTT.value = tt
  await load()
}

async function load() {
  loading.value = true
  try {
    if (currentTT.value === 'knight') {
      const res = await getKnightLeaderboard()
      const data = res.data
      const users = data?.users || data?.comics?.users || []
      knights.value = Array.isArray(users) ? users : (Array.isArray(users?.docs) ? users.docs : [])
      comics.value = []
    } else {
      const res = await getLeaderboard(currentTT.value as LeaderboardTT)
      const comicsData = res.data?.comics
      comics.value = Array.isArray(comicsData) ? comicsData : (Array.isArray(comicsData?.docs) ? comicsData.docs : [])
      knights.value = []
    }
  } catch {} finally { loading.value = false }
}

function getUserAvatar(user: any): string {
  if (user?.avatar?.fileServer && user?.avatar?.path) {
    return `/api/image/proxy?fileServer=${encodeURIComponent(user.avatar.fileServer)}&path=${encodeURIComponent(user.avatar.path)}`
  }
  if (typeof user?.avatar === 'string' && user.avatar) {
    return `/api/image/proxy?url=${encodeURIComponent(user.avatar)}`
  }
  return ''
}

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = '#f0f0f0'
}

function goComic(id: string) { router.push(`/comic/${id}`) }
</script>

<style scoped>
.rank-badge { position: absolute; top: 8px; left: 8px; width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 700; color: white; z-index: 10; }
.rank-1 { background: #ff4757; } .rank-2 { background: #ff6b81; } .rank-3 { background: #ffa502; }
.rank-badge:not(.rank-1):not(.rank-2):not(.rank-3) { background: rgba(0,0,0,0.5); }
.comic-card { position: relative; }

/* 骑士榜 */
.knight-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.knight-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 12px 12px 48px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.knight-rank {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
}

.knight-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  font-weight: 600;
}

.knight-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.knight-info {
  flex: 1;
  min-width: 0;
}

.knight-name {
  font-size: 15px;
  font-weight: 500;
}

.knight-level {
  font-size: 11px;
  color: var(--primary);
  background: #fce4ec;
  padding: 1px 6px;
  border-radius: 8px;
  margin-left: 6px;
}

.knight-title {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.knight-stat {
  font-size: 12px;
  color: var(--text-secondary);
  max-width: 40%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
