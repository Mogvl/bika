<template>
  <div class="detail-page" v-if="!loading">
    <div class="detail-header">
      <div class="detail-cover-wrap">
        <img :src="getCoverUrl(game.thumb)" :alt="game.title" class="detail-cover" />
      </div>
      <div class="detail-info">
        <h1 class="detail-title">{{ game.title }}</h1>
        <p class="detail-author">厂商: {{ game.author }}</p>
        <p class="detail-meta">章节: {{ game.epsCount }} | 喜欢: {{ game.totalLikes }}</p>
        <div class="detail-tags">
          <span v-for="cat in game.categories" :key="cat" class="tag">{{ cat }}</span>
        </div>
        <button class="btn btn-primary" @click="startPlaying" :disabled="!eps.length">
          {{ eps.length ? '开始游戏' : '暂无章节' }}
        </button>
      </div>
    </div>
    <div class="detail-section">
      <h3 class="section-title">简介</h3>
      <p class="detail-desc">{{ game.description || '暂无简介' }}</p>
    </div>
    <div class="detail-section">
      <h3 class="section-title">章节列表 ({{ eps.length }})</h3>
      <div class="eps-list">
        <div v-for="ep in eps" :key="ep._id" class="eps-item" @click="goReader(ep._id)">
          <span class="eps-order">第 {{ ep.order }} 话</span>
          <span class="eps-title">{{ ep.title || `第 ${ep.order} 话` }}</span>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading">加载中</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getGameDetail, getGameEps } from '@/api'

const route = useRoute()
const router = useRouter()
const game = ref<any>({})
const eps = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  const id = route.params.id as string
  try {
    const [detailRes, epsRes] = await Promise.all([getGameDetail(id), getGameEps(id)])
    game.value = detailRes.data?.game || {}
    eps.value = epsRes.data?.eps || []
  } catch {} finally { loading.value = false }
})

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function startPlaying() { if (eps.value.length) goReader(eps.value[0]._id) }
function goReader(epsId: string) { router.push(`/reader/${route.params.id}/${epsId}`) }
</script>

<style scoped>
.detail-page { max-width: 800px; margin: 0 auto; }
.detail-header { display: flex; gap: 16px; padding: 16px; background: var(--bg-card); }
.detail-cover-wrap { flex-shrink: 0; width: 140px; }
.detail-cover { width: 100%; border-radius: 8px; aspect-ratio: 3/4; object-fit: cover; }
.detail-info { flex: 1; min-width: 0; }
.detail-title { font-size: 18px; font-weight: 600; margin-bottom: 8px; }
.detail-author { font-size: 13px; color: var(--text-secondary); margin-bottom: 4px; }
.detail-meta { font-size: 12px; color: var(--text-secondary); margin-bottom: 8px; }
.detail-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 12px; }
.tag { padding: 2px 10px; border-radius: 12px; font-size: 11px; background: #fce4ec; color: var(--primary); }
.detail-section { padding: 16px; background: var(--bg-card); margin-top: 8px; }
.detail-desc { font-size: 14px; line-height: 1.7; color: var(--text-secondary); }
.eps-list { display: flex; flex-direction: column; }
.eps-item { display: flex; align-items: center; padding: 12px 8px; border-bottom: 1px solid var(--border); cursor: pointer; gap: 12px; }
.eps-item:hover { background: var(--bg); }
.eps-order { font-size: 13px; color: var(--primary); font-weight: 500; min-width: 56px; }
.eps-title { flex: 1; font-size: 14px; }
</style>
