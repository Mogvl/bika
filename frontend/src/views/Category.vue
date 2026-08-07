<template>
  <div class="page-container">
    <div class="section-title">📂 分类</div>

    <div v-if="loading" class="loading">加载中</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>

    <div v-else class="category-grid">
      <!-- 固定入口 -->
      <div class="category-card" @click="goRank">
        <div class="category-icon rank-icon">🏆</div>
        <div class="category-name">排行榜</div>
      </div>
      <div class="category-card" @click="goSearch('')">
        <div class="category-icon update-icon">📅</div>
        <div class="category-name">最近更新</div>
      </div>
      <div class="category-card" @click="goRandom">
        <div class="category-icon random-icon">🎲</div>
        <div class="category-name">随机本子</div>
      </div>

      <!-- 动态分类 -->
      <div
        v-for="cat in categories"
        :key="cat.title"
        class="category-card"
        @click="goCategory(cat.title)"
      >
        <img
          :src="getCoverUrl(cat.thumb)"
          :alt="cat.title"
          class="category-img"
          loading="lazy"
          @error="handleImgError"
        />
        <div class="category-name">{{ cat.title }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories } from '@/api'
import type { Category } from '@/types'

const router = useRouter()
const categories = ref<Category[]>([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const res = await getCategories()
    const cats = res.data?.categories || []
    categories.value = cats.filter((c: Category) => c.title)
  } catch (e: any) {
    error.value = '加载分类失败'
  } finally {
    loading.value = false
  }
})

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = 'rgba(255,255,255,0.06)'
}

function goRank() {
  router.push('/leaderboard')
}

function goSearch(keyword: string) {
  router.push('/search')
}

function goRandom() {
  router.push('/home')
}

function goCategory(title: string) {
  router.push({ path: '/category/comics', query: { c: title } })
}
</script>

<style scoped>
.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 16px;
  padding: 4px;
}

.category-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 12px 8px;
  border-radius: 18px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  transition: transform 0.22s cubic-bezier(0.22, 1, 0.36, 1), box-shadow 0.22s, border-color 0.22s;
}

.category-card:hover {
  transform: translateY(-4px);
  border-color: var(--primary);
  box-shadow: var(--shadow);
}

.category-icon {
  width: 72px;
  height: 72px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30px;
  color: #fff;
  box-shadow: 0 4px 14px rgba(164, 117, 208, 0.15);
}

.rank-icon {
  background: linear-gradient(135deg, #b98ce8, #e08bb8);
}

.update-icon {
  background: linear-gradient(135deg, #6db7ff, #3a8dde);
}

.random-icon {
  background: linear-gradient(135deg, #b98ce8, #a475d0);
}

.category-img {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  object-fit: cover;
  background: var(--bg);
}

.category-name {
  font-size: 12px;
  text-align: center;
  color: var(--text);
  max-width: 88px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
