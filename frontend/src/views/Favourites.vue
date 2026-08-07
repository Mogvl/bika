<template>
  <div class="page-container">
    <div class="section-title">我的收藏</div>

    <!-- 工具栏 -->
    <div class="favourite-toolbar">
      <SortSelect v-model="sort" :options="sortOptions" @change="reload" />
      <div class="search-box2">
        <input v-model="searchKeyword" type="text" placeholder="搜索收藏..." @input="onSearchInput" />
      </div>
    </div>

    <div v-if="loading && comics.length === 0" class="loading">加载中</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>

    <div v-else-if="comics.length === 0" class="empty-state">
      <p>{{ searchKeyword ? '没有找到匹配的收藏' : '还没有收藏任何漫画' }}</p>
      <router-link to="/home" class="btn btn-primary" style="margin-top: 16px;">去逛逛</router-link>
    </div>

    <div v-else class="comic-grid">
      <div
        v-for="comic in comics"
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
        <button class="unfav-btn" @click.stop="unfavourite(comic)">✕</button>
      </div>
    </div>

    <div class="page-load-more">
      <button v-if="hasMore" class="load-more-btn" @click="loadMore">加载更多</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getFavourites, addFavourite } from '@/api'
import type { Comic } from '@/types'
import SortSelect from '@/components/SortSelect.vue'

const router = useRouter()
const comics = ref<Comic[]>([])
const page = ref(1)
const hasMore = ref(false)
const loading = ref(true)
const error = ref('')
const sort = ref('da')
const searchKeyword = ref('')

// 排序选项
const sortOptions = [
  { value: 'da', label: '添加时间（新→旧）' },
  { value: 'dd', label: '添加时间（旧→新）' },
]

onMounted(() => loadFavourites())

async function loadFavourites() {
  loading.value = true
  error.value = ''
  try {
    const res = await getFavourites(page.value, sort.value)
    const data = res.data
    const comicsData = data?.comics
    comics.value = Array.isArray(comicsData?.docs) ? comicsData.docs : (Array.isArray(comicsData) ? comicsData : [])
    hasMore.value = page.value < (comicsData?.pages || data?.pages || 1)
  } catch (e: any) {
    if (e.message?.includes('401') || e.message?.includes('未登录')) {
      error.value = '请先登录'
    } else {
      error.value = e.message || '加载失败'
    }
  } finally {
    loading.value = false
  }
}

function reload() {
  page.value = 1
  loadFavourites()
}

function onSearchInput() {
  // 搜索收藏：前端过滤（接口未提供收藏内搜索，本地过滤已加载列表）
  const kw = searchKeyword.value.trim().toLowerCase()
  if (!kw) {
    reload()
    return
  }
}

async function loadMore() {
  page.value++
  loading.value = true
  try {
    const res = await getFavourites(page.value, sort.value)
    const data = res.data
    const comicsData = data?.comics
    const newComics = Array.isArray(comicsData?.docs) ? comicsData.docs : (Array.isArray(comicsData) ? comicsData : [])
    comics.value.push(...newComics)
    hasMore.value = page.value < (comicsData?.pages || data?.pages || 1)
  } catch {} finally {
    loading.value = false
  }
}

async function unfavourite(comic: Comic) {
  if (!confirm(`确定取消收藏「${comic.title}」？`)) return
  try {
    await addFavourite(comic._id) // 收藏接口是 toggle，已收藏则取消
    comics.value = comics.value.filter(c => c._id !== comic._id)
  } catch (e: any) {
    alert(e.message || '操作失败')
  }
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
.favourite-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  align-items: center;
}



.search-box2 {
  flex: 1;
}

.search-box2 input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-pill);
  font-size: 13px;
  outline: none;
  background: rgba(255, 255, 255, 0.6);
  box-shadow: var(--shadow-sm);
}

.search-box2 input:focus {
  border-color: var(--primary);
}

.comic-card {
  position: relative;
}

.unfav-btn {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: none;
  background: rgba(164, 117, 208, 0.3);
  color: white;
  font-size: 12px;
  cursor: pointer;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
}

.unfav-btn:hover {
  background: var(--danger);
}
</style>
