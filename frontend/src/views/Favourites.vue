<template>
  <div class="page-container">
    <div class="section-title">我的收藏</div>

    <div v-if="loading && comics.length === 0" class="loading">加载中</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>

    <div v-else-if="comics.length === 0" class="empty-state">
      <p>还没有收藏任何漫画</p>
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
import { getFavourites } from '@/api'
import type { Comic } from '@/types'

const router = useRouter()
const comics = ref<Comic[]>([])
const page = ref(1)
const hasMore = ref(false)
const loading = ref(true)
const error = ref('')

onMounted(() => loadFavourites())

async function loadFavourites() {
  loading.value = true
  error.value = ''
  try {
    const res = await getFavourites(page.value)
    const data = res.data
    comics.value = data?.comics || []
    hasMore.value = page.value < (data?.pages || 1)
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

async function loadMore() {
  page.value++
  loading.value = true
  try {
    const res = await getFavourites(page.value)
    const data = res.data
    comics.value.push(...(data?.comics || []))
    hasMore.value = page.value < (data?.pages || 1)
  } catch {} finally {
    loading.value = false
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
