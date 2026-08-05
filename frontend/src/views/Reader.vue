<template>
  <div class="reader-page">
    <!-- 顶部工具栏 -->
    <div class="reader-toolbar" :class="{ hidden: controlsHidden }">
      <button class="toolbar-btn" @click="goBack">✕</button>
      <div class="toolbar-info">
        <span class="eps-label">{{ currentEpsTitle }}</span>
      </div>
      <button class="toolbar-btn" @click="toggleMode">
        {{ singleMode ? '📖' : '📜' }}
      </button>
    </div>

    <div class="reader-content" ref="readerContent" @click="toggleControls" @touchstart="onTouchStart" @touchend="onTouchEnd">
      <!-- 单页模式 -->
      <div v-if="singleMode" class="single-page" @click.stop>
        <div class="page-nav-area prev-page" @click.stop="prevPage"></div>
        <img
          v-if="currentImageUrl"
          :src="currentImageUrl"
          class="reader-image"
          alt="漫画页"
          @load="onImageLoaded"
          @error="onImageError"
        />
        <div class="page-nav-area next-page" @click.stop="nextPage"></div>
      </div>

      <!-- 滚动模式 -->
      <div v-else class="scroll-mode" ref="scrollContainer">
        <img
          v-for="(img, idx) in allImages"
          :key="idx"
          :src="img"
          class="scroll-image"
          alt="漫画页"
          loading="lazy"
          @load="onScrollImageLoaded(idx)"
          @error="onScrollImageError($event, idx)"
        />
      </div>
    </div>

    <!-- 底部章节切换 -->
    <div class="reader-bottom" :class="{ hidden: controlsHidden }">
      <button class="btn-nav" @click="prevEps">上一话</button>
      <select v-model="currentOrder" @change="switchEps" class="eps-select">
        <option v-for="ep in epsList" :key="ep.order" :value="ep.order">
          第 {{ ep.order }} 话
        </option>
      </select>
      <button class="btn-nav" @click="nextEps">下一话</button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="reader-loading">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-if="error" class="reader-error">
      <p>{{ error }}</p>
      <button class="btn btn-primary" @click="reload">重试</button>
    </div>

    <!-- 页码指示器 -->
    <div class="page-indicator" v-if="!controlsHidden">
      {{ currentIndex + 1 }} / {{ totalPages }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComicEps, getComicPages, getGameEps, getGamePages } from '@/api'
import { saveReadingProgress } from '@/utils/history'
import type { EP, Page } from '@/types'

const route = useRoute()
const router = useRouter()

const comicId = computed(() => route.params.id as string)
const isGame = computed(() => route.query.type === 'game')
const currentOrder = ref(Number(route.params.epsId) || 1)  // 章节序号
const epsList = ref<EP[]>([])
const totalEpsLoaded = ref(false)

const pages = ref<Page[]>([])
const currentIndex = ref(0)
const loading = ref(true)
const error = ref('')
const controlsHidden = ref(false)
const singleMode = ref(true)

const readerContent = ref<HTMLElement | null>(null)
const scrollContainer = ref<HTMLElement | null>(null)

const currentEpsTitle = computed(() => {
  const ep = epsList.value.find(e => e.order === currentOrder.value)
  return ep ? `第 ${ep.order} 话` : `第 ${currentOrder.value} 话`
})

const totalPages = computed(() => pages.value.length)

const currentImageUrl = computed(() => {
  if (pages.value.length === 0) return ''
  const page = pages.value[currentIndex.value]
  return getImageSrc(page)
})

const allImages = computed(() => {
  return pages.value.map(p => getImageSrc(p))
})

let touchStartX = 0
let touchStartY = 0

function onTouchStart(e: TouchEvent) {
  touchStartX = e.touches[0].clientX
  touchStartY = e.touches[0].clientY
}

function onTouchEnd(e: TouchEvent) {
  if (!singleMode.value) return
  const dx = e.changedTouches[0].clientX - touchStartX
  const dy = e.changedTouches[0].clientY - touchStartY
  if (Math.abs(dx) > Math.abs(dy) && Math.abs(dx) > 30) {
    if (dx < 0) nextPage()
    else prevPage()
  }
}

function getImageSrc(page: Page): string {
  if (!page?.media?.fileServer || !page?.media?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(page.media.fileServer)}&path=${encodeURIComponent(page.media.path)}`
}

onMounted(async () => {
  await loadEpsList()
  await loadPages()
})

watch(currentOrder, async () => {
  await loadPages()
})

async function loadEpsList() {
  try {
    const res = isGame.value
      ? await getGameEps(comicId.value, 1)
      : await getComicEps(comicId.value, 1)
    const epsData = res.data?.eps
    epsList.value = Array.isArray(epsData?.docs) ? epsData.docs : (Array.isArray(epsData) ? epsData : [])
    totalEpsLoaded.value = true
  } catch {}
}

async function loadPages() {
  loading.value = true
  error.value = ''
  currentIndex.value = 0
  pages.value = []

  try {
    const res = isGame.value
      ? await getGamePages(comicId.value, String(currentOrder.value))
      : await getComicPages(comicId.value, String(currentOrder.value))
    const data = res.data
    const pagesData = data?.pages
    pages.value = Array.isArray(pagesData?.docs) ? pagesData.docs : (Array.isArray(pagesData) ? pagesData : [])
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function nextPage() {
  if (currentIndex.value < pages.value.length - 1) {
    currentIndex.value++
    saveProgress()
  } else {
    nextEps()
  }
}

function prevPage() {
  if (currentIndex.value > 0) {
    currentIndex.value--
    saveProgress()
  }
}

function saveProgress() {
  const ep = epsList.value.find(e => e.order === currentOrder.value)
  saveReadingProgress(
    comicId.value,
    ep?.title || `第${currentOrder.value}话`,
    null,
    String(currentOrder.value),
    currentOrder.value,
    ep?.title || `第${currentOrder.value}话`,
    currentIndex.value
  )
}

function goBack() {
  router.push(`/comic/${comicId.value}`)
}

function toggleControls() {
  if (!singleMode.value) return
  controlsHidden.value = !controlsHidden.value
}

function toggleMode() {
  singleMode.value = !singleMode.value
  controlsHidden.value = false
}

function switchEps() {
  router.replace(`/reader/${comicId.value}/${currentOrder.value}`)
}

function nextEps() {
  const idx = epsList.value.findIndex(e => e.order === currentOrder.value)
  if (idx < epsList.value.length - 1) {
    currentOrder.value = epsList.value[idx + 1].order
    router.replace(`/reader/${comicId.value}/${currentOrder.value}`)
  }
}

function prevEps() {
  const idx = epsList.value.findIndex(e => e.order === currentOrder.value)
  if (idx > 0) {
    currentOrder.value = epsList.value[idx - 1].order
    router.replace(`/reader/${comicId.value}/${currentOrder.value}`)
  }
}

function reload() {
  loadPages()
}

function onImageLoaded() {}

function onImageError(e: Event) {
  const img = e.target as HTMLImageElement
  img.src = ''
  img.alt = '加载失败'
  img.style.background = '#f0f0f0'
  img.style.minHeight = '300px'
}

function onScrollImageLoaded(idx: number) {}

function onScrollImageError(e: Event, idx: number) {
  const img = e.target as HTMLImageElement
  img.alt = '加载失败'
  img.style.background = '#f0f0f0'
  img.style.minHeight = '200px'
}
</script>

<style scoped>
.reader-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #1a1a1a;
  display: flex;
  flex-direction: column;
  z-index: 300;
}

.reader-toolbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 48px;
  background: rgba(0,0,0,0.8);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  z-index: 310;
  transition: opacity 0.3s;
}

.reader-toolbar.hidden {
  opacity: 0;
  pointer-events: none;
}

.toolbar-btn {
  background: none;
  border: none;
  color: white;
  font-size: 20px;
  cursor: pointer;
  padding: 8px;
}

.toolbar-info {
  color: white;
  font-size: 14px;
}

.reader-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

/* 单页模式 */
.single-page {
  display: flex;
  height: 100%;
  align-items: center;
  justify-content: center;
}

.page-nav-area {
  position: absolute;
  top: 48px;
  bottom: 48px;
  width: 30%;
  z-index: 5;
  cursor: pointer;
}

.prev-page { left: 0; }
.next-page { right: 0; }

.reader-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  user-select: none;
  -webkit-user-drag: none;
}

/* 滚动模式 */
.scroll-mode {
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.scroll-image {
  display: block;
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

/* 底部导航 */
.reader-bottom {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 48px;
  background: rgba(0,0,0,0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 12px;
  z-index: 310;
  transition: opacity 0.3s;
}

.reader-bottom.hidden {
  opacity: 0;
  pointer-events: none;
}

.btn-nav {
  background: rgba(255,255,255,0.15);
  border: none;
  color: white;
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
}

.eps-select {
  background: rgba(255,255,255,0.15);
  border: none;
  color: white;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  max-width: 120px;
}

.eps-select option {
  background: #333;
  color: white;
}

.reader-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  gap: 16px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255,255,255,0.2);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.reader-error {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  gap: 16px;
}

.page-indicator {
  position: fixed;
  bottom: 56px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0,0,0,0.6);
  color: white;
  padding: 4px 16px;
  border-radius: 12px;
  font-size: 12px;
  z-index: 305;
}
</style>
