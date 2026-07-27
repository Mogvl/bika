<template>
  <div>
    <div class="search-box">
      <span class="search-icon">🔍</span>
      <input
        v-model="keyword"
        type="text"
        placeholder="搜索漫画..."
        @keydown.enter="doSearch"
        @input="onInput"
      />
      <button v-if="keyword" class="clear-btn" @click="keyword = ''; results = []">✕</button>
    </div>

    <!-- 搜索热词 -->
    <div v-if="!searched && keywords.length > 0" class="page-container">
      <div class="section-title">热门搜索</div>
      <div class="keyword-tags">
        <span v-for="kw in keywords" :key="kw" class="keyword-tag" @click="searchKeyword(kw)">
          {{ kw }}
        </span>
      </div>
    </div>

    <!-- 搜索结果 -->
    <div v-if="searched" class="page-container">
      <p v-if="results.length === 0 && !loading" class="empty-state">
        没有找到 "{{ lastKeyword }}" 的相关漫画
      </p>

      <div v-else class="comic-grid">
        <div
          v-for="comic in results"
          :key="comic._id"
          class="comic-card"
          @click="goComic(comic._id)"
        >
          <img
            :src="getCoverUrl(comic.thumb)"
            :alt="comic.title"
            class="comic-cover"
            loading="lazy"
          />
          <div class="comic-info">
            <div class="comic-title">{{ comic.title }}</div>
            <div class="comic-author">{{ comic.author }}</div>
          </div>
        </div>
      </div>

      <div v-if="loading" class="loading">搜索中</div>

      <div class="page-load-more">
        <button v-if="hasMore" class="load-more-btn" @click="loadMore">加载更多</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getKeywords, searchComics, getComicsByCategory } from '@/api'
import type { Comic } from '@/types'

const route = useRoute()
const router = useRouter()
const keyword = ref('')
const lastKeyword = ref('')
const results = ref<Comic[]>([])
const keywords = ref<string[]>([])
const page = ref(1)
const hasMore = ref(false)
const loading = ref(false)
const searched = ref(false)
let searchTimer: number | undefined

onMounted(async () => {
  // 检查是否有分类查询参数
  const category = route.query.c as string
  if (category) {
    keyword.value = category
    await doCategorySearch(category)
  } else {
    try {
      const res = await getKeywords()
      keywords.value = res.data?.keywords || []
    } catch {}
  }
})

function onInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = window.setTimeout(() => {
    if (keyword.value.trim()) doSearch()
  }, 600)
}

async function doSearch() {
  const kw = keyword.value.trim()
  if (!kw) return
  lastKeyword.value = kw
  page.value = 1
  results.value = []
  searched.value = true
  hasMore.value = false
  loading.value = true

  try {
    const res = await searchComics(kw, page.value)
    const data = res.data
    const comicsData = data?.comics
    results.value = Array.isArray(comicsData?.docs) ? comicsData.docs : (Array.isArray(comicsData) ? comicsData : [])
    hasMore.value = page.value < (comicsData?.pages || data?.pages || 1)
  } catch {} finally {
    loading.value = false
  }
}

async function doCategorySearch(category: string) {
  lastKeyword.value = category
  page.value = 1
  results.value = []
  searched.value = true
  hasMore.value = false
  loading.value = true

  try {
    const res = await getComicsByCategory(page.value, category)
    const data = res.data
    const comicsData = data?.comics
    results.value = Array.isArray(comicsData?.docs) ? comicsData.docs : (Array.isArray(comicsData) ? comicsData : [])
    hasMore.value = page.value < (comicsData?.pages || data?.pages || 1)
  } catch {} finally {
    loading.value = false
  }
}

async function loadMore() {
  page.value++
  loading.value = true
  try {
    const res = await searchComics(lastKeyword.value, page.value)
    const data = res.data
    const comicsData = data?.comics
    const newComics = Array.isArray(comicsData?.docs) ? comicsData.docs : (Array.isArray(comicsData) ? comicsData : [])
    results.value.push(...newComics)
    hasMore.value = page.value < (comicsData?.pages || data?.pages || 1)
  } catch {} finally {
    loading.value = false
  }
}

function searchKeyword(kw: string) {
  keyword.value = kw
  doSearch()
}

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function goComic(id: string) {
  router.push(`/comic/${id}`)
}
</script>

<style scoped>
.clear-btn {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: var(--text-secondary);
}

.keyword-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.keyword-tag {
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 13px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
}

.keyword-tag:hover {
  border-color: var(--primary);
  color: var(--primary);
}
</style>
