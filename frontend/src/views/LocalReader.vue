<template>
  <div class="local-reader-page">
    <!-- 顶部工具栏 -->
    <div class="reader-toolbar" :class="{ hidden: controlsHidden }">
      <button class="toolbar-btn" @click="goBack">✕</button>
      <div class="toolbar-info">
        <span class="eps-label">{{ epsTitle }}</span>
      </div>
      <div class="toolbar-actions">
        <button class="toolbar-btn" @click="toggleMode">{{ singleMode ? '📖' : '📜' }}</button>
      </div>
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
          @error="onScrollImageError"
        />
      </div>
    </div>

    <!-- 底部章节切换 -->
    <div class="reader-bottom" :class="{ hidden: controlsHidden }">
      <button class="btn-nav" @click="prevEps">上一话</button>
      <span class="eps-label2">{{ currentIndex + 1 }} / {{ totalPages }}</span>
      <button class="btn-nav" @click="nextEps">下一话</button>
    </div>

    <div v-if="error" class="reader-error">
      <p>{{ error }}</p>
      <button class="btn btn-primary" @click="goBack">返回</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getLocalEps, getLocalImageUrl } from '@/api'

const route = useRoute()
const router = useRouter()

const comicPath = (route.query.comic as string) || ''
const epsName = (route.query.eps as string) || ''
const epsTitle = (route.query.title as string) || '本地漫画'

const pages = ref<string[]>([])
const currentIndex = ref(0)
const controlsHidden = ref(false)
const singleMode = ref(true)
const error = ref('')

let touchStartX = 0
let touchStartY = 0

const totalPages = computed(() => pages.value.length)

const currentImageUrl = computed(() => {
  if (pages.value.length === 0) return ''
  return getLocalImageUrl(pages.value[currentIndex.value])
})

const allImages = computed(() => pages.value.map(p => getLocalImageUrl(p)))

onMounted(() => {
  loadPages()
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

// 从本地库 eps 接口获取指定章节的图片列表
async function loadPages() {
  if (!comicPath || !epsName) {
    error.value = '参数错误'
    return
  }
  try {
    const res = await getLocalEps(comicPath)
    const epsList = res.data?.eps || []
    const target = epsList.find((e: any) => e.title === epsName)
    if (target && target.pages?.length) {
      // 图片路径是相对漫画目录的，需拼接漫画路径
      pages.value = target.pages.map((p: string) => comicPath + '/' + p)
      currentIndex.value = 0
    } else {
      error.value = '未找到该章节'
    }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  }
}

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

function nextPage() {
  if (currentIndex.value < pages.value.length - 1) currentIndex.value++
}

function prevPage() {
  if (currentIndex.value > 0) currentIndex.value--
}

function goBack() {
  router.push('/local')
}

function toggleControls() {
  controlsHidden.value = !controlsHidden.value
}

function toggleMode() {
  singleMode.value = !singleMode.value
  controlsHidden.value = false
}

function prevEps() { /* 本地章节切换简化：暂不支持 */ }
function nextEps() { /* 本地章节切换简化：暂不支持 */ }

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') goBack()
  else if (e.key === 'ArrowRight' || e.key === 'ArrowDown' || e.key === ' ') { e.preventDefault(); nextPage() }
  else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') prevPage()
}

function onImageError(e: Event) {
  const img = e.target as HTMLImageElement
  img.alt = '加载失败'
  img.style.background = '#f0f0f0'
  img.style.minHeight = '300px'
}

function onScrollImageError(e: Event) {
  const img = e.target as HTMLImageElement
  img.alt = '加载失败'
  img.style.background = '#f0f0f0'
}
</script>

<style scoped>
.local-reader-page { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: #1a1a1a; display: flex; flex-direction: column; z-index: 300; }
.reader-toolbar { position: fixed; top: 0; left: 0; right: 0; height: 48px; background: rgba(0,0,0,0.8); display: flex; align-items: center; justify-content: space-between; padding: 0 12px; z-index: 310; transition: opacity 0.3s; }
.reader-toolbar.hidden { opacity: 0; pointer-events: none; }
.toolbar-btn { background: none; border: none; color: white; font-size: 20px; cursor: pointer; padding: 8px; }
.toolbar-info { color: white; font-size: 14px; flex: 1; text-align: center; }
.toolbar-actions { display: flex; gap: 4px; }
.reader-content { flex: 1; overflow: hidden; position: relative; }
.single-page { display: flex; height: 100%; align-items: center; justify-content: center; }
.page-nav-area { position: absolute; top: 48px; bottom: 48px; width: 30%; z-index: 5; cursor: pointer; }
.prev-page { left: 0; }
.next-page { right: 0; }
.reader-image { max-width: 100%; max-height: 100%; object-fit: contain; user-select: none; }
.scroll-mode { height: 100%; overflow-y: auto; -webkit-overflow-scrolling: touch; }
.scroll-image { display: block; width: 100%; max-width: 800px; margin: 0 auto; }
.reader-bottom { position: fixed; bottom: 0; left: 0; right: 0; height: 48px; background: rgba(0,0,0,0.8); display: flex; align-items: center; justify-content: center; gap: 12px; padding: 0 12px; z-index: 310; transition: opacity 0.3s; }
.reader-bottom.hidden { opacity: 0; pointer-events: none; }
.btn-nav { background: rgba(255,255,255,0.15); border: none; color: white; padding: 6px 16px; border-radius: 6px; font-size: 13px; cursor: pointer; }
.eps-label2 { color: white; font-size: 13px; }
.reader-error { position: fixed; top: 0; left: 0; right: 0; bottom: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; color: white; gap: 16px; }
</style>
