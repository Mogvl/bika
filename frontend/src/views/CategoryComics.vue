<template>
  <div class="category-comics-page">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-input-wrap">
        <input
          v-model="keyword"
          type="text"
          placeholder="搜索漫画..."
          class="search-input"
          @keydown.enter="doSearch"
        />
        <button class="search-btn" @click="doSearch">🔍</button>
      </div>
    </div>

    <!-- 分类标签过滤 -->
    <div class="filter-section">
      <div class="filter-header">
        <span class="filter-label">分类筛选</span>
        <button class="filter-toggle" @click="showFilters = !showFilters">
          {{ showFilters ? '收起' : '展开' }}
        </button>
      </div>
      <div v-show="showFilters" class="filter-tags">
        <span
          v-for="cat in allCategories"
          :key="cat"
          class="filter-tag"
          :class="{ active: selectedCategories.includes(cat) }"
          @click="toggleCategory(cat)"
        >
          {{ cat }}
        </span>
        <button v-if="selectedCategories.length > 0" class="filter-clear" @click="clearFilters">清除筛选</button>
      </div>
    </div>

    <!-- 排序选项 -->
    <div class="sort-bar">
      <div class="sort-group">
        <span class="sort-label">排序:</span>
        <span class="sort-item" :class="{ active: sort === 'ua' }" @click="changeSort('ua')">最新</span>
        <span class="sort-item" :class="{ active: sort === 'dd' }" @click="changeSort('dd')">新→旧</span>
        <span class="sort-item" :class="{ active: sort === 'da' }" @click="changeSort('da')">旧→新</span>
        <span class="sort-item" :class="{ active: sort === 'ld' }" @click="changeSort('ld')">最多喜欢</span>
        <span class="sort-item" :class="{ active: sort === 'vv' }" @click="changeSort('vv')">最多浏览</span>
      </div>
    </div>

    <!-- 搜索结果/分类漫画 -->
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
            <div class="comic-categories">
              <span v-for="cat in (comic.categories || []).slice(0, 2)" :key="cat" class="comic-cat-tag">{{ cat }}</span>
            </div>
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
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComicsByCategory, searchComics, getCategories } from '@/api'
import type { Comic } from '@/types'

const route = useRoute()
const router = useRouter()

const keyword = ref('')
const categoryName = ref((route.query.c as string) || '')
const comics = ref<Comic[]>([])
const page = ref(1)
const totalPages = ref(1)
const sort = ref('ua')
const loading = ref(false)
const error = ref('')
const showFilters = ref(false)
const allCategories = ref<string[]>([])
const selectedCategories = ref<string[]>([])

onMounted(async () => {
  // 加载所有分类用于筛选
  try {
    const res = await getCategories()
    const cats = res.data?.categories || []
    allCategories.value = cats.map((c: any) => c.title).filter(Boolean)
  } catch {}

  // 如果有指定分类，预选
  if (categoryName.value) {
    selectedCategories.value = [categoryName.value]
    await loadComics()
  }
})

async function loadComics() {
  loading.value = true
  error.value = ''
  try {
    let res
    if (keyword.value.trim()) {
      // 有搜索关键词时用搜索接口
      res = await searchComics(keyword.value, page.value, selectedCategories.value[0] || '', sort.value)
    } else if (selectedCategories.value.length === 1) {
      // 单个分类
      res = await getComicsByCategory(page.value, selectedCategories.value[0], sort.value)
    } else if (selectedCategories.value.length > 1) {
      // 多个分类用搜索接口
      res = await searchComics('', page.value, selectedCategories.value.join(','), sort.value)
    } else {
      // 无分类无关键词，获取全部
      res = await getComicsByCategory(page.value, '大家都在看', sort.value)
    }

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

function doSearch() {
  page.value = 1
  comics.value = []
  loadComics()
}

function toggleCategory(cat: string) {
  const idx = selectedCategories.value.indexOf(cat)
  if (idx >= 0) {
    selectedCategories.value.splice(idx, 1)
  } else {
    selectedCategories.value.push(cat)
  }
  page.value = 1
  comics.value = []
  loadComics()
}

function clearFilters() {
  selectedCategories.value = []
  page.value = 1
  comics.value = []
  loadComics()
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
.category-comics-page {
  max-width: 100%;
}

.search-bar {
  padding: 12px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
}

.search-input-wrap {
  display: flex;
  gap: 8px;
  max-width: 600px;
  margin: 0 auto;
}

.search-input {
  flex: 1;
  padding: 10px 16px;
  border: 1px solid var(--border);
  border-radius: 20px;
  font-size: 15px;
  outline: none;
}

.search-input:focus {
  border-color: var(--primary);
}

.search-btn {
  padding: 10px 16px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 20px;
  font-size: 16px;
  cursor: pointer;
}

/* 分类筛选 */
.filter-section {
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  padding: 8px 12px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.filter-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.filter-toggle {
  background: none;
  border: none;
  font-size: 12px;
  color: var(--primary);
  cursor: pointer;
}

.filter-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-height: 200px;
  overflow-y: auto;
}

.filter-tag {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  background: var(--bg);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
}

.filter-tag:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.filter-tag.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.filter-clear {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  background: #e74c3c;
  color: white;
  border: none;
  cursor: pointer;
}

/* 排序 */
.sort-bar {
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  padding: 8px 12px;
}

.sort-group {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow-x: auto;
}

.sort-label {
  font-size: 13px;
  color: var(--text-secondary);
  white-space: nowrap;
}

.sort-item {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  cursor: pointer;
  white-space: nowrap;
  color: var(--text-secondary);
}

.sort-item.active {
  background: var(--primary);
  color: white;
}

/* 漫画网格 */
.comic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
  padding: 12px;
}

.comic-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: var(--shadow);
  cursor: pointer;
}

.comic-cover {
  width: 100%;
  aspect-ratio: 3/4;
  object-fit: cover;
}

.comic-info {
  padding: 8px 10px;
}

.comic-title {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.comic-author {
  font-size: 11px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.comic-categories {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.comic-cat-tag {
  padding: 1px 6px;
  border-radius: 8px;
  font-size: 10px;
  background: #fce4ec;
  color: var(--primary);
}

/* 分页 */
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
