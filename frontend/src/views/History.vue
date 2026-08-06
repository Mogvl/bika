<template>
  <div class="page-container">
    <div class="history-header">
      <div class="section-title">📖 阅读历史</div>
      <button v-if="history.length > 0" class="clear-btn" @click="handleClear">清空</button>
    </div>

    <div v-if="history.length === 0" class="empty-state">
      <p>暂无阅读历史</p>
      <p class="empty-tip">阅读漫画后会自动记录</p>
    </div>

    <div v-else class="comic-grid">
      <div
        v-for="item in history"
        :key="item.bookId"
        class="comic-card"
        @click="goComic(item.bookId)"
      >
        <img :src="getCoverUrl(item.thumb)" :alt="item.title" class="comic-cover" loading="lazy" @error="handleImgError" />
        <div class="comic-info">
          <div class="comic-title">{{ item.title }}</div>
          <div class="comic-author">{{ item.epsTitle || '继续阅读' }}</div>
          <div class="comic-time">{{ formatTime(item.lastRead) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getHistory, clearHistory } from '@/utils/history'

const router = useRouter()
const history = ref<any[]>([])

onMounted(() => {
  loadHistory()
})

function loadHistory() {
  history.value = getHistory()
}

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = 'var(--bg-soft)'
}

function goComic(id: string) {
  router.push(`/comic/${id}`)
}

function handleClear() {
  if (confirm('确定清空所有阅读历史？')) {
    clearHistory()
    history.value = []
  }
}

function formatTime(timestamp: number): string {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / 86400000)
  if (days > 0) return `${days}天前`
  const hours = Math.floor(diff / 3600000)
  if (hours > 0) return `${hours}小时前`
  const minutes = Math.floor(diff / 60000)
  if (minutes > 0) return `${minutes}分钟前`
  return '刚刚'
}
</script>

<style scoped>
.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.clear-btn {
  padding: 6px 16px;
  background: var(--danger);
  color: white;
  border: none;
  border-radius: 16px;
  font-size: 13px;
  cursor: pointer;
}

.empty-tip {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 8px;
}

.comic-time {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
