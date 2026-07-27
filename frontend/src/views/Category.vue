<template>
  <div class="page-container">
    <div class="section-title">📂 分类浏览</div>

    <!-- 分类标签 -->
    <div class="category-tags">
      <span
        v-for="cat in categories"
        :key="cat.title"
        class="category-tag"
        :class="{ active: selectedCategory === cat.title }"
        @click="selectCategory(cat.title)"
      >
        {{ cat.title }}
      </span>
    </div>

    <!-- 漫画列表 -->
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

      <div class="page-load-more">
        <button v-if="hasMore" class="load-more-btn" @click="loadMore">加载更多</button>
        <span v-else class="load-more-btn" style="color: #999; cursor: default;">没有更多了</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories, getComicsByCategory } from '@/api'
import type { Comic, Category } from '@/types'

const router = useRouter()
const categories = ref<Category[]>([])
const selectedCategory = ref('')
const comics = ref<Comic[]>([])
const page = ref(1)
const hasMore = ref(true)
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  await loadCategories()
  if (categories.value.length > 0) {
    await selectCategory(categories.value[0].title)
  }
})

async function loadCategories() {
  try {
    const res = await getCategories()
    const cats = res.data?.categories || []
    categories.value = cats.filter((c: Category) => c.title)
  } catch (e: any) {
    error.value = '加载分类失败'
  }
}

async function selectCategory(cat: string) {
  if (selectedCategory.value === cat) return
  selectedCategory.value = cat
  page.value = 1
  comics.value = []
  hasMore.value = true
  await loadComics()
}

async function loadComics() {
  if (!selectedCategory.value) return
  loading.value = true
  error.value = ''
  try {
    const res = await getComicsByCategory(page.value, selectedCategory.value)
    const data = res.data
    const comicsData = data?.comics
    const newComics = Array.isArray(comicsData?.docs) ? comicsData.docs : (Array.isArray(comicsData) ? comicsData : [])
    comics.value.push(...newComics)
    hasMore.value = page.value < (comicsData?.pages || data?.pages || 1)
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  page.value++
  await loadComics()
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
.category-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px 0;
  overflow-x: auto;
}

.category-tag {
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 13px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.category-tag:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.category-tag.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}
</style>
