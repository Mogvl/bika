<template>
  <div class="page-container">
    <div class="section-title">🎮 游戏区</div>

    <div v-if="loading" class="loading">加载中</div>
    <div v-else>
      <div class="comic-grid">
        <div v-for="game in games" :key="game._id" class="comic-card" @click="goGame(game._id)">
          <img :src="getCoverUrl(game.thumb)" :alt="game.title" class="comic-cover" loading="lazy" @error="handleImgError" />
          <div class="comic-info">
            <div class="comic-title">{{ game.title }}</div>
            <div class="comic-author">{{ game.author }}</div>
          </div>
        </div>
      </div>
      <div class="page-load-more">
        <button v-if="hasMore" class="load-more-btn" @click="loadMore">加载更多</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getGames } from '@/api'

const router = useRouter()
const games = ref<any[]>([])
const page = ref(1)
const hasMore = ref(true)
const loading = ref(true)

onMounted(() => loadGames())

async function loadGames() {
  loading.value = true
  try {
    const res = await getGames(page.value)
    const data = res.data
    const gamesData = data?.games
    const newGames = Array.isArray(gamesData?.docs) ? gamesData.docs : (Array.isArray(gamesData) ? gamesData : [])
    games.value.push(...newGames)
    hasMore.value = page.value < (gamesData?.pages || data?.pages || 1)
  } catch {} finally { loading.value = false }
}

async function loadMore() { page.value++; await loadGames() }

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = '#f0f0f0'
}

function goGame(id: string) { router.push(`/game/${id}`) }
</script>
