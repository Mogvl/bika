<template>
  <div class="page-container">
    <div class="category-header">
      <button class="btn-back" @click="$router.back()">←</button>
      <h2 class="category-title">{{ categoryName }}</h2>
    </div>

    <!-- 排序选项 -->
    <div class="sort-bar">
      <span class="sort-item" :class="{ active: sort === 'ua' }" @click="changeSort('ua')">最新</span>
      <span class="sort-item" :class="{ active: sort === 'dd' }" @click="changeSort('dd')">新→旧</span>
      <span class="sort-item" :class="{ active: sort === 'da' }" @click="changeSort('da')">旧→新</span>
      <span class="sort-item" :class="{ active: sort === 'ld' }" @click="changeSort('ld')">最多喜欢</span>
      <span class="sort-item" :class="{ active: sort === 'vv' }" @click="changeSort('vv')">最多浏览</span>
    </div>

    <div v-if="loading && comics.length === 0" class="loading">加载中</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>

    <div v-else>
      <div class="comic-grid">
        <div
          v-for="comic in comics"
          :key="comic._id"
          class="comic-card"
          @click="goComic(comic._id)"
        >
          <img :src="getCoverUrl(comic.thumb)" :alt="comic.title" class="comic-cover" loading="lazy" @error="handleImgError" />
          <div class="comic-info">
            <div class="comic-title">{{ comic.title }}</div>
            <div class="comic-author">{{ comic.author }}</div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <button class="page-btn" :disabled="page <= 1" @click="prevPage">上一页</button>
        <span class="page-info">{{ page }} / {{ totalPages }}</span>
        <button class="page-btn" :disabled="page >= totalPages" @click="nextPage">下一页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComicsByCategory } from '@/api'
import type { Comic } from '@/types'

const route = useRoute()
const router = useRouter()

const categoryName = computed(() => (route.query.c as string) || '')
const comics = ref<Comic[]>([])
const page = ref(1)
const totalPages = ref(1)
const sort = ref('ua')
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  if (categoryName.value) {
    await loadComics()
  } else {
    error.value = '未指定分类'
  }
})

async function loadComics() {
  loading.value = true
  error.value = ''
  try {
    const res = await getComicsByCategory(page.value, categoryName.value, sort.value)
    const data = res.data
    const comicsData = data?.comics
    if (Array.isArray(comicsData?.docs)) {
      comics.value = comicsData.docs
      totalPages.value = comicsData.pages || 1
    } else if (Array.isArray(comicsData)) {
      comics.value = comicsData
      totalPages.value = data?.pages || 1
    } else {
      comics.value = []
    }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function changeSort(s: string) {
  sort.value = s
  page.value = 1
  comics.value = []
  loadComics()
}

function prevPage() {
  if (page.value > 1) {
    page.value--
    loadComics()
  }
}

function nextPage() {
  if (page.value < totalPages.value) {
    page.value++
    loadComics()
  }
}

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = '#f0f0f0'
}

function goComic(id: string) {
  router.push(`/comic/${id}`)
}
</script>

<style scoped>
.category-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--bg-card);
}

.btn-back {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--text);
  padding: 4px 8px;
}

.category-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--primary);
}

.sort-bar {
  display: flex;
  gap: 0;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  overflow-x: auto;
}

.sort-item {
  padding: 10px 16px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  white-space: nowrap;
  border-bottom: 2px solid transparent;
}

.sort-item.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
  font-weight: 500;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 20px;
}

.page-btn {
  padding: 8px 20px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: var(--text-secondary);
}
</style>
