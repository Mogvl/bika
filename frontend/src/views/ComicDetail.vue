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
        <div class="detail-actions">
          <button class="btn btn-primary" @click="startReading" :disabled="!eps.length">
            {{ eps.length ? '开始阅读' : '暂无章节' }}
          </button>
          <button class="btn btn-download" @click="downloadComic">📥 下载</button>
        </div>
      </div>
    </div>

    <!-- 简介 -->
    <div class="detail-section">
      <h3 class="section-title">简介</h3>
      <p class="detail-desc">{{ comic.description || '暂无简介' }}</p>
    </div>

    <!-- 章节列表 -->
    <div class="detail-section">
      <div class="eps-header">
        <h3 class="section-title">章节 ({{ eps.length }})</h3>
      </div>

      <div v-if="eps.length === 0" class="empty-state">暂无章节</div>

      <div v-else>
        <div class="eps-toolbar">
          <button class="eps-btn" @click="toggleSelectMode">
            {{ selectMode ? '取消选择' : '选择下载' }}
          </button>
          <template v-if="selectMode">
            <button class="eps-btn" @click="selectAll">全选</button>
            <button class="eps-btn" @click="selectNone">全不选</button>
            <button class="eps-btn eps-btn-download" @click="downloadSelected" :disabled="selectedEps.length === 0">
              📥 下载选中 ({{ selectedEps.length }})
            </button>
          </template>
        </div>

        <div class="eps-grid">
          <div
            v-for="ep in eps"
            :key="ep._id"
            class="eps-chip"
            :class="{
              'selected': selectedEps.includes(ep.order),
              'select-mode': selectMode
            }"
            @click="handleEpsClick(ep)"
          >
            <span v-if="selectMode" class="eps-check">
              {{ selectedEps.includes(ep.order) ? '☑' : '☐' }}
            </span>
            {{ ep.title || `第${ep.order}话` }}
          </div>
        </div>
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
import { getComicDetail, getComicEps, addDownload } from '@/api'
import { saveComicHistory } from '@/utils/history'
import type { Comic, EP } from '@/types'

const route = useRoute()
const router = useRouter()
const comic = ref<Comic>({} as Comic)
const eps = ref<EP[]>([])
const loading = ref(true)
const error = ref('')
const epsPage = ref(1)
const epsHasMore = ref(false)
const selectedEps = ref<number[]>([])
const selectMode = ref(false)

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
    const epsData = epsRes.data?.eps
    eps.value = Array.isArray(epsData?.docs) ? epsData.docs : (Array.isArray(epsData) ? epsData : [])
    epsHasMore.value = epsPage.value < (epsData?.pages || 1)
    saveComicHistory(comic.value)
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
    const epsData = res.data?.eps
    const newEps = Array.isArray(epsData?.docs) ? epsData.docs : (Array.isArray(epsData) ? epsData : [])
    eps.value.push(...newEps)
    epsHasMore.value = epsPage.value < (epsData?.pages || 1)
  } catch {}
}

function handleEpsClick(ep: EP) {
  if (selectMode.value) {
    // 选择模式：切换选中状态
    toggleEps(ep.order)
  } else {
    // 普通模式：打开阅读器
    router.push(`/reader/${route.params.id}/${ep.order}`)
  }
}

function toggleSelectMode() {
  selectMode.value = !selectMode.value
  if (!selectMode.value) {
    selectedEps.value = []
  }
}

function toggleEps(order: number) {
  const idx = selectedEps.value.indexOf(order)
  if (idx >= 0) {
    selectedEps.value.splice(idx, 1)
  } else {
    selectedEps.value.push(order)
  }
}

function selectAll() {
  selectedEps.value = eps.value.map(ep => ep.order)
}

function selectNone() {
  selectedEps.value = []
}

function startReading() {
  if (eps.value.length > 0) {
    router.push(`/reader/${route.params.id}/${eps.value[0].order}`)
  }
}

async function downloadComic() {
  const coverUrl = comic.value.thumb?.fileServer && comic.value.thumb?.path
    ? `${comic.value.thumb.fileServer}/static/${comic.value.thumb.path}`
    : ''
  try {
    await addDownload(route.params.id as string, comic.value.title || '', coverUrl)
    alert('已添加到下载队列！')
  } catch (e: any) {
    alert(e.message || '添加下载失败')
  }
}

async function downloadSelected() {
  if (selectedEps.value.length === 0) return
  const coverUrl = comic.value.thumb?.fileServer && comic.value.thumb?.path
    ? `${comic.value.thumb.fileServer}/static/${comic.value.thumb.path}`
    : ''
  try {
    await addDownload(route.params.id as string, comic.value.title || '', coverUrl)
    alert(`已添加 ${selectedEps.value.length} 个章节到下载队列！`)
    selectMode.value = false
    selectedEps.value = []
  } catch (e: any) {
    alert(e.message || '添加下载失败')
  }
}

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
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

.detail-actions {
  display: flex;
  gap: 8px;
}

.btn-download {
  background: #3498db;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
}

.btn-download:hover {
  background: #2980b9;
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

/* 章节网格 - 和原版一致 */
.eps-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.eps-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.eps-chip {
  padding: 8px 16px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.eps-chip:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: #fff0f5;
}

.eps-chip.selected {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.eps-chip.select-mode {
  padding: 8px 12px;
}

.eps-check {
  margin-right: 4px;
}

.eps-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.eps-btn {
  padding: 6px 14px;
  border: 1px solid var(--border);
  border-radius: 16px;
  font-size: 13px;
  cursor: pointer;
  background: var(--bg);
  color: var(--text);
  transition: all 0.2s;
}

.eps-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.eps-btn-download {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.eps-btn-download:hover {
  background: var(--primary-dark);
}

.eps-btn-download:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
