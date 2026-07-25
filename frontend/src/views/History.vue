<template>
  <div class="page-container">
    <div class="section-title">📖 阅读历史</div>
    <div v-if="history.length === 0" class="empty-state">暂无阅读历史</div>
    <div v-else class="comic-grid">
      <div v-for="item in history" :key="item.bookId" class="comic-card" @click="goComic(item.bookId)">
        <img :src="getCoverUrl(item.thumb)" :alt="item.title" class="comic-cover" loading="lazy" @error="handleImgError" />
        <div class="comic-info">
          <div class="comic-title">{{ item.title }}</div>
          <div class="comic-author">{{ item.epsTitle || '继续阅读' }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const history = ref<any[]>([])

onMounted(() => {
  const saved = localStorage.getItem('read_history')
  if (saved) {
    try { history.value = JSON.parse(saved) } catch {}
  }
})

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
