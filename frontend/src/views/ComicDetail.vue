<template>
  <div class="detail-page" v-if="!loading">
    <!-- 头部信息 -->
    <div class="detail-header">
      <div class="detail-cover-wrap">
        <img :src="getCoverUrl(comic.thumb)" :alt="comic.title" class="detail-cover" />
      </div>
      <div class="detail-info">
        <h1 class="detail-title">{{ comic.title }}</h1>
        <p class="detail-author">作者: {{ comic.author }}</p>
        <p class="detail-meta">
          章节: {{ comic.epsCount }} | 喜欢: {{ comic.totalLikes || comic.likesCount }}
        </p>
        <div class="detail-tags">
          <span v-for="cat in comic.categories" :key="cat" class="tag">{{ cat }}</span>
        </div>
        <button
          class="btn btn-primary"
          @click="startReading"
          :disabled="!eps.length"
        >
          {{ eps.length ? '开始阅读' : '暂无章节' }}
        </button>
      </div>
    </div>

    <!-- 简介 -->
    <div class="detail-section">
      <h3 class="section-title">简介</h3>
      <p class="detail-desc">{{ comic.description || '暂无简介' }}</p>
    </div>

    <!-- 章节列表 -->
    <div class="detail-section">
      <h3 class="section-title">章节列表 ({{ eps.length }})</h3>
      <div class="eps-list">
        <div
          v-for="ep in eps"
          :key="ep._id"
          class="eps-item"
          @click="goReader(ep._id)"
        >
          <span class="eps-order">第 {{ ep.order }} 话</span>
          <span class="eps-title">{{ ep.title || `第 ${ep.order} 话` }}</span>
          <span class="eps-pages">{{ ep.pagesCount }}页</span>
        </div>
        <div v-if="eps.length === 0" class="empty-state">暂无章节</div>
      </div>
      <div v-if="epsHasMore" class="page-load-more">
        <button class="load-more-btn" @click="loadMoreEps">加载更多章节</button>
      </div>
    </div>
  </div>

  <div v-else class="loading">加载中</div>

  <div v-if="error" class="error-msg">{{ error }}</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComicDetail, getComicEps } from '@/api'
import type { Comic, EP } from '@/types'

const route = useRoute()
const router = useRouter()
const comic = ref<Comic>({} as Comic)
const eps = ref<EP[]>([])
const loading = ref(true)
const error = ref('')
const epsPage = ref(1)
const epsHasMore = ref(false)

onMounted(async () => {
  const id = route.params.id as string
  if (!id) {
    error.value = '参数错误'
    loading.value = false
    return
  }
  await loadDetail(id)
})

async function loadDetail(id: string) {
  try {
    const [detailRes, epsRes] = await Promise.all([
      getComicDetail(id),
      getComicEps(id, 1),
    ])
    comic.value = detailRes.data?.comic || {}
    const epsData = epsRes.data
    eps.value = epsData?.eps || []
    epsHasMore.value = epsPage.value < (epsData?.pages || 1)
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

async function loadMoreEps() {
  epsPage.value++
  try {
    const res = await getComicEps(route.params.id as string, epsPage.value)
    const data = res.data
    eps.value.push(...(data?.eps || []))
    epsHasMore.value = epsPage.value < (data?.pages || 1)
  } catch {}
}

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function startReading() {
  if (eps.value.length > 0) {
    goReader(eps.value[0]._id)
  }
}

function goReader(epsId: string) {
  router.push(`/reader/${route.params.id}/${epsId}`)
}
</script>

<style scoped>
.detail-page {
  max-width: 800px;
  margin: 0 auto;
}

.detail-header {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: var(--bg-card);
}

.detail-cover-wrap {
  flex-shrink: 0;
  width: 140px;
}

.detail-cover {
  width: 100%;
  border-radius: 8px;
  aspect-ratio: 3/4;
  object-fit: cover;
}

.detail-info {
  flex: 1;
  min-width: 0;
}

.detail-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  line-height: 1.3;
}

.detail-author {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.detail-meta {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.detail-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.tag {
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 11px;
  background: #fce4ec;
  color: var(--primary);
}

.detail-section {
  padding: 16px;
  background: var(--bg-card);
  margin-top: 8px;
}

.detail-desc {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-secondary);
}

.eps-list {
  display: flex;
  flex-direction: column;
}

.eps-item {
  display: flex;
  align-items: center;
  padding: 12px 8px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background 0.2s;
  gap: 12px;
}

.eps-item:hover { background: var(--bg); }
.eps-item:active { background: #f0f0f0; }

.eps-order {
  font-size: 13px;
  color: var(--primary);
  font-weight: 500;
  min-width: 56px;
}

.eps-title {
  flex: 1;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.eps-pages {
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
