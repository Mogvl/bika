<template>
  <div class="page-container">
    <div class="section-title">📚 本地库</div>

    <div v-if="loading" class="loading">加载中</div>

    <div v-else-if="comics.length === 0" class="empty-state">
      <p>还没有下载任何漫画</p>
      <p class="empty-tip">在漫画详情页点击"下载"后，下载完成的漫画会出现在这里</p>
    </div>

    <div v-else>
      <!-- 本地漫画列表 -->
      <div class="comic-grid">
        <div v-for="comic in comics" :key="comic.path" class="comic-card" @click="toggleExpand(comic)">
          <img :src="getCoverUrl(comic)" :alt="comic.title" class="comic-cover" loading="lazy" @error="handleImgError" />
          <div class="comic-info">
            <div class="comic-title">{{ comic.title }}</div>
            <div class="comic-author">{{ comic.eps.length }} 个章节</div>
          </div>
        </div>
      </div>

      <!-- 展开的章节列表 -->
      <div v-if="expanded" class="local-detail">
        <div class="detail-header">
          <button class="btn-back" @click="expanded = null">← 返回列表</button>
          <h3 class="detail-title">{{ expanded.title }}</h3>
        </div>
        <div class="eps-list">
          <div
            v-for="ep in expanded.eps"
            :key="ep.path"
            class="eps-item"
            @click="openReader(expanded, ep)"
          >
            <span class="eps-order">📖</span>
            <span class="eps-title">{{ ep.title }}</span>
            <span class="eps-count">{{ ep.pages.length }} 页</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getLocalLibrary, getLocalImageUrl } from '@/api'

interface LocalEps {
  title: string
  path: string
  pages: string[]
}

interface LocalComic {
  title: string
  path: string
  cover: string
  eps: LocalEps[]
}

const router = useRouter()
const comics = ref<LocalComic[]>([])
const loading = ref(true)
const expanded = ref<LocalComic | null>(null)

onMounted(() => loadLibrary())

async function loadLibrary() {
  loading.value = true
  try {
    const res = await getLocalLibrary()
    comics.value = res.data?.comics || []
  } catch {} finally { loading.value = false }
}

function toggleExpand(comic: LocalComic) {
  if (expanded.value?.path === comic.path) {
    expanded.value = null
  } else {
    expanded.value = comic
  }
}

function getCoverUrl(comic: LocalComic): string {
  if (!comic.cover) return ''
  return getLocalImageUrl(comic.path + '/' + comic.cover)
}

function openReader(comic: LocalComic, ep: LocalEps) {
  // 使用本地库阅读器路由
  router.push({
    path: '/local-reader',
    query: {
      comic: comic.path,
      eps: ep.path,
      title: comic.title,
    },
  })
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = '#f0f0f0'
}
</script>

<style scoped>
.comic-card { position: relative; }
.empty-tip { font-size: 13px; color: #999; margin-top: 8px; }

.local-detail {
  margin-top: 16px;
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 16px;
  box-shadow: var(--shadow);
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.btn-back {
  background: none;
  border: none;
  font-size: 15px;
  color: var(--primary);
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 16px;
}

.btn-back:hover { background: #fce4ec; }

.detail-title { font-size: 18px; font-weight: 600; flex: 1; }

.eps-list { display: flex; flex-direction: column; }

.eps-item {
  display: flex;
  align-items: center;
  padding: 12px 8px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  gap: 12px;
}

.eps-item:hover { background: var(--bg); }

.eps-order { font-size: 16px; }
.eps-title { flex: 1; font-size: 14px; }
.eps-count { font-size: 12px; color: var(--text-secondary); }
</style>
