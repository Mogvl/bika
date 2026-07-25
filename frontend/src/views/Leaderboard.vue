<template>
  <div class="page-container">
    <div class="section-title">🏆 排行榜</div>

    <div class="tabs">
      <div v-for="tab in tabs" :key="tab.value" class="tab-item" :class="{ active: currentTT === tab.value }" @click="switchTab(tab.value)">
        {{ tab.label }}
      </div>
    </div>

    <div v-if="loading" class="loading">加载中</div>
    <div v-else class="comic-grid" style="padding-top: 12px;">
      <div v-for="(comic, idx) in comics" :key="comic._id" class="comic-card" @click="goComic(comic._id)">
        <div class="rank-badge" :class="'rank-' + (idx + 1)">{{ idx + 1 }}</div>
        <img :src="getCoverUrl(comic.thumb)" :alt="comic.title" class="comic-cover" loading="lazy" @error="handleImgError" />
        <div class="comic-info">
          <div class="comic-title">{{ comic.title }}</div>
          <div class="comic-author">❤️ {{ comic.totalLikes || comic.likesCount }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getLeaderboard } from '@/api'
import type { Comic, LeaderboardTT } from '@/types'

const router = useRouter()
const tabs = [
  { label: '24小时', value: 'H24' as LeaderboardTT },
  { label: '本周', value: 'D7' as LeaderboardTT },
  { label: '本月', value: 'D30' as LeaderboardTT },
]
const currentTT = ref<LeaderboardTT>('H24')
const comics = ref<Comic[]>([])
const loading = ref(true)

onMounted(() => loadLeaderboard())

async function switchTab(tt: LeaderboardTT) {
  currentTT.value = tt
  await loadLeaderboard()
}

async function loadLeaderboard() {
  loading.value = true
  try {
    const res = await getLeaderboard(currentTT.value)
    comics.value = res.data?.comics || []
  } catch {} finally { loading.value = false }
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
</style>
