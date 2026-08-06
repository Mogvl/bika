<template>
  <div class="home-page">
    <!-- 标签页 -->
    <div class="tabs">
      <div class="tab-item" :class="{ active: currentTab === 'random' }" @click="currentTab = 'random'">随机推荐</div>
      <div class="tab-item" :class="{ active: currentTab === 'collections' }" @click="currentTab = 'collections'">本子神推荐</div>
      <div class="tab-item" :class="{ active: currentTab === 'recommend' }" @click="currentTab = 'recommend'">本子魔推荐</div>
    </div>

    <!-- 随机推荐 -->
    <div v-show="currentTab === 'random'" class="tab-content">
      <div v-if="randomLoading && randomComics.length === 0" class="loading">加载中</div>
      <div v-else-if="randomError" class="error-msg">{{ randomError }}</div>
      <div v-else>
        <div class="comic-grid">
          <div
            v-for="comic in randomComics"
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
        <div class="refresh-bar">
          <button class="refresh-btn" @click="loadRandom" :disabled="randomLoading">
            {{ randomLoading ? '加载中...' : '🔄 刷新 (F5)' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 本子神推荐 -->
    <div v-show="currentTab === 'collections'" class="tab-content">
      <div v-if="collectionsLoading && collections.length === 0" class="loading">加载中</div>
      <div v-else-if="collectionsError" class="error-msg">{{ collectionsError }}</div>
      <div v-else>
        <div v-for="(collection, idx) in collections" :key="idx" class="collection-section">
          <h3 class="collection-title">{{ collection.title || `推荐 ${idx + 1}` }}</h3>
          <div class="comic-grid">
            <div
              v-for="comic in (collection.comics || [])"
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
        </div>
      </div>
    </div>

    <!-- 本子魔推荐 -->
    <div v-show="currentTab === 'recommend'" class="tab-content">
      <div v-if="recommendLoading && recommendComics.length === 0" class="loading">加载中</div>
      <div v-else-if="recommendError" class="error-msg">{{ recommendError }}</div>
      <div v-else>
        <div class="comic-grid">
          <div
            v-for="comic in recommendComics"
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
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getRandomComics, getCollections } from '@/api'
import type { Comic } from '@/types'

const router = useRouter()
const currentTab = ref('random')

// 随机推荐
const randomComics = ref<Comic[]>([])
const randomLoading = ref(false)
const randomError = ref('')

// 合集
const collections = ref<any[]>([])
const collectionsLoading = ref(false)
const collectionsError = ref('')

// 推荐
const recommendComics = ref<Comic[]>([])
const recommendLoading = ref(false)
const recommendError = ref('')

onMounted(() => {
  loadRandom()
  loadCollections()
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'F5' && currentTab.value === 'random') {
    e.preventDefault()
    loadRandom()
  }
}

async function loadRandom() {
  randomLoading.value = true
  randomError.value = ''
  try {
    const res = await getRandomComics()
    const comicsData = res.data?.comics
    randomComics.value = Array.isArray(comicsData) ? comicsData : (Array.isArray(comicsData?.docs) ? comicsData.docs : [])
  } catch (e: any) {
    randomError.value = e.message || '加载失败'
  } finally {
    randomLoading.value = false
  }
}

async function loadCollections() {
  collectionsLoading.value = true
  collectionsError.value = ''
  try {
    const res = await getCollections()
    const data = res.data?.collections
    if (Array.isArray(data)) {
      collections.value = data
      // 本子魔推荐 - 第二个合集或随机漫画
      if (data.length > 1 && data[1]?.comics) {
        recommendComics.value = data[1].comics || []
      } else {
        // 没有第二个合集，用随机漫画替代
        loadRecommend()
      }
    }
  } catch (e: any) {
    collectionsError.value = e.message || '加载失败'
  } finally {
    collectionsLoading.value = false
  }
}

async function loadRecommend() {
  recommendLoading.value = true
  recommendError.value = ''
  try {
    const res = await getRandomComics()
    const comicsData = res.data?.comics
    recommendComics.value = Array.isArray(comicsData) ? comicsData : (Array.isArray(comicsData?.docs) ? comicsData.docs : [])
  } catch (e: any) {
    recommendError.value = e.message || '加载失败'
  } finally {
    recommendLoading.value = false
  }
}

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = 'rgba(255,255,255,0.06)'
}

function goComic(id: string) {
  router.push(`/comic/${id}`)
}
</script>

<style scoped>
.home-page {
  max-width: 100%;
}

.tab-content {
  padding: 12px;
}

.collection-section {
  margin-bottom: 24px;
}

.collection-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
  padding: 0 4px;
  color: var(--primary);
}

.refresh-bar {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.refresh-btn {
  padding: 10px 32px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 15px;
  cursor: pointer;
  transition: background 0.2s;
}

.refresh-btn:hover {
  background: var(--primary-dark);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
