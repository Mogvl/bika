<template>
  <div class="search-page">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <button class="btn-back" @click="$router.back()">←</button>
      <div class="search-input-wrap">
        <input
          v-model="keyword"
          type="text"
          :placeholder="categoryName || '搜索漫画...'"
          class="search-input"
          @keydown.enter="doSearch"
        />
        <button class="search-btn" @click="doSearch">🔍</button>
      </div>
    </div>

    <!-- 分类筛选（展开/收起） -->
    <div class="filter-section">
      <button class="filter-toggle" @click="showFilters = !showFilters">
        {{ showFilters ? '收起分类 ▲' : '展开分类 ▼' }}
      </button>
      <div v-show="showFilters" class="filter-body">
        <div class="filter-actions">
          <button class="filter-action" @click="selectAll">全选</button>
          <button class="filter-action" @click="clearFilters">清除</button>
        </div>
        <div class="filter-tags">
          <span
            v-for="cat in allCategories"
            :key="cat"
            class="filter-tag"
            :class="{ active: selectedCategories.includes(cat) }"
            @click="toggleCategory(cat)"
          >
            {{ cat }}
          </span>
        </div>
      </div>
    </div>

    <!-- 排序 -->
    <div class="sort-bar">
      <div class="sort-group">
        <span class="sort-label">排序:</span>
        <select v-model="sort" @change="doSearch" class="sort-select">
          <option value="ua">最新</option>
          <option value="dd">新→旧</option>
          <option value="da">旧→新</option>
          <option value="ld">最多喜欢</option>
          <option value="vv">最多浏览</option>
        </select>
      </div>
    </div>

    <!-- 结果 -->
    <div v-if="loading && comics.length === 0" class="loading">加载中</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>
    <div v-else-if="comics.length === 0" class="empty-state">没有找到相关漫画</div>

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
        <input v-model.number="jumpPage" type="number" class="page-input" min="1" :max="totalPages" @keydown.enter="goPage" />
        <span class="page-info">/ {{ totalPages }}</span>
        <button class="page-btn" @click="goPage">跳转</button>
        <button class="page-btn" :disabled="page >= totalPages" @click="nextPage">下一页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
const jumpPage = ref(1)
const sort = ref('ua')
const loading = ref(false)
const error = ref('')
const showFilters = ref(false)
const allCategories = ref<string[]>([])
const selectedCategories = ref<string[]>([])

onMounted(async () => {
  // 加载所有分类
  try {
    const res = await getCategories()
    const cats = res.data?.categories || []
    allCategories.value = cats.map((c: any) => c.title).filter(Boolean)
  } catch {}

  // 如果指定了分类
  if (categoryName.value) {
    selectedCategories.value = [categoryName.value]
    keyword.value = categoryName.value
    await loadComics()
  }
})

async function loadComics() {
  loading.value = true
  error.value = ''
  try {
    let res
    const searchKeyword = keyword.value.trim()
    const cats = selectedCategories.value

    if (searchKeyword && cats.length > 0) {
      // 有关键词+分类：用搜索接口
      res = await searchComics(searchKeyword, page.value, cats[0], sort.value)
    } else if (searchKeyword) {
      // 纯关键词搜索
      res = await searchComics(searchKeyword, page.value, '', sort.value)
    } else if (cats.length > 0) {
      // 有分类：用分类接口（取第一个选中的分类）
      res = await getComicsByCategory(page.value, cats[0], sort.value)
    } else {
      // 无分类无关键词：默认分类
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
    jumpPage.value = page.value
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
  doSearch()
}

function selectAll() {
  selectedCategories.value = [...allCategories.value]
  doSearch()
}

function clearFilters() {
  selectedCategories.value = []
  keyword.value = ''
  doSearch()
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

function goPage() {
  const p = jumpPage.value
  if (p >= 1 && p <= totalPages.value) {
    page.value = p
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
.search-page {
  max-width: 100%;
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
}

.btn-back {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--text);
  padding: 4px 8px;
  flex-shrink: 0;
}

.search-input-wrap {
  flex: 1;
  display: flex;
  gap: 8px;
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
  flex-shrink: 0;
}

/* 分类筛选 */
.filter-section {
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
}

.filter-toggle {
  width: 100%;
  padding: 10px 16px;
  background: none;
  border: none;
  font-size: 13px;
  color: var(--primary);
  cursor: pointer;
  text-align: left;
}

.filter-body {
  padding: 0 16px 12px;
  border-top: 1px solid var(--border);
}

.filter-actions {
  display: flex;
  gap: 8px;
  padding: 8px 0;
}

.filter-action {
  padding: 4px 12px;
  border: 1px solid var(--border);
  border-radius: 12px;
  font-size: 12px;
  background: var(--bg);
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

/* 排序 */
.sort-bar {
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  padding: 8px 16px;
}

.sort-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sort-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.sort-select {
  padding: 6px 12px;
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 13px;
  background: var(--bg);
  outline: none;
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
}

/* 分页 */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px;
}

.page-btn {
  padding: 8px 16px;
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

.page-input {
  width: 50px;
  padding: 6px 8px;
  border: 1px solid var(--border);
  border-radius: 8px;
  text-align: center;
  font-size: 14px;
  outline: none;
}

.page-input:focus {
  border-color: var(--primary);
}

.page-info {
  font-size: 14px;
  color: var(--text-secondary);
}
</style>
