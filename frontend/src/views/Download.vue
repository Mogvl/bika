<template>
  <div class="page-container">
    <div class="section-title">📥 下载管理</div>

    <!-- 状态筛选 -->
    <div class="filter-bar">
      <button
        v-for="f in filters"
        :key="f.value"
        class="filter-btn"
        :class="{ active: currentFilter === f.value }"
        @click="currentFilter = f.value"
      >
        {{ f.label }}
      </button>
    </div>

    <div v-if="loading" class="loading">加载中</div>

    <div v-else-if="filteredTasks.length === 0" class="empty-state">
      <p>暂无下载任务</p>
      <p class="empty-tip">在漫画详情页点击"下载"按钮添加任务</p>
    </div>

    <div v-else class="download-list">
      <div v-for="task in filteredTasks" :key="task.id" class="download-item">
        <div class="download-cover">
          <img :src="getCoverUrl(task.coverUrl)" :alt="task.title" @error="handleImgError" />
        </div>
        <div class="download-info">
          <div class="download-title">{{ task.title }}</div>
          <div class="download-status">
            <span :class="'status-' + task.status">{{ getStatusText(task.status) }}</span>
            <span v-if="task.status === 'downloading'" class="download-progress">
              {{ task.downloaded }}/{{ task.totalPages || '?' }}
            </span>
            <span v-if="task.status === 'downloading' && task.speed" class="download-speed">
              ⚡ {{ task.speed }}
            </span>
          </div>
          <div v-if="task.error" class="download-error">{{ task.error }}</div>
          <div class="download-path">📁 {{ getDisplayPath(task) }}</div>
          <div class="download-time">{{ formatTime(task.updatedAt) }}</div>
        </div>
        <div class="download-actions">
          <button v-if="task.status === 'downloading'" class="btn-action btn-cancel" @click="cancel(task.id)">暂停</button>
          <button v-if="task.status === 'paused' || task.status === 'error'" class="btn-action btn-start" @click="resume(task.id)">继续</button>
          <button v-if="task.status === 'completed'" class="btn-action btn-start" @click="openLocal(task)">📚 本地</button>
          <button class="btn-action btn-remove" @click="remove(task.id, task.title)">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDownloads, cancelDownload, resumeDownload, removeDownload } from '@/api'

interface DownloadTask {
  id: string
  bookId: string
  title: string
  coverUrl: string
  savePath?: string
  totalPages: number
  downloaded: number
  speed?: string
  status: string
  error?: string
  createdAt: string
  updatedAt: string
}

const router = useRouter()
const tasks = ref<DownloadTask[]>([])
const loading = ref(true)
const currentFilter = ref('all')
let timer: number | null = null

const filters = [
  { label: '全部', value: 'all' },
  { label: '下载中', value: 'downloading' },
  { label: '已完成', value: 'completed' },
  { label: '已暂停', value: 'paused' },
]

const filteredTasks = computed(() => {
  if (currentFilter.value === 'all') return tasks.value
  return tasks.value.filter(t => t.status === currentFilter.value)
})

onMounted(() => {
  loadTasks()
  // 每2秒刷新一次状态
  timer = window.setInterval(loadTasks, 2000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

async function loadTasks() {
  try {
    const res = await getDownloads()
    tasks.value = res.data?.tasks || []
  } catch {} finally { loading.value = false }
}

async function cancel(id: string) {
  try {
    await cancelDownload(id)
    await loadTasks()
  } catch {}
}

async function resume(id: string) {
  try {
    await resumeDownload(id)
    await loadTasks()
  } catch {}
}

async function remove(id: string, title: string) {
  const deleteFile = confirm(`确定删除「${title}」的下载任务？\n\n确认则同时删除已下载文件。`)
  try {
    await removeDownload(id, deleteFile)
    await loadTasks()
  } catch {}
}

function openLocal(task: DownloadTask) {
  router.push('/local')
}

function getCoverUrl(url: string): string {
  if (!url) return ''
  if (url.includes('fileServer')) {
    return url
  }
  return `/api/image/proxy?url=${encodeURIComponent(url)}`
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = '#f0f0f0'
}

function sanitizeTitle(title: string): string {
  return title.replace(/[\/\\:*?"<>|]/g, '_').substring(0, 200)
}

function getDisplayPath(task: DownloadTask): string {
  const path = task.savePath || ('downloads/' + sanitizeTitle(task.title))
  return path
    .replace(/^\/data\/downloads/, '/volume1/bika')
    .replace(/^downloads/, '/volume1/bika')
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    waiting: '等待中',
    downloading: '下载中',
    completed: '已完成',
    error: '失败',
    paused: '已暂停',
  }
  return map[status] || status
}

function formatTime(t: string): string {
  if (!t) return ''
  const date = new Date(t)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.filter-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  overflow-x: auto;
  padding: 4px;
}

.filter-btn {
  padding: 6px 16px;
  border: 1px solid var(--border);
  border-radius: 16px;
  font-size: 13px;
  background: var(--bg-card);
  color: var(--text-secondary);
  cursor: pointer;
  white-space: nowrap;
}

.filter-btn.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.download-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.download-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.download-cover {
  width: 60px;
  height: 80px;
  flex-shrink: 0;
  border-radius: 6px;
  overflow: hidden;
}

.download-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.download-info {
  flex: 1;
  min-width: 0;
}

.download-title {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.download-status {
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  gap: 8px;
  align-items: center;
}

.status-waiting { color: #f39c12; }
.status-downloading { color: #3498db; }
.status-completed { color: #27ae60; }
.status-error { color: #e74c3c; }
.status-paused { color: #95a5a6; }

.download-progress {
  color: var(--primary);
  font-weight: 500;
}

.download-speed {
  color: #27ae60;
  font-weight: 500;
}

.download-error {
  font-size: 11px;
  color: #e74c3c;
  margin-top: 4px;
}

.download-path {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 4px;
  font-family: monospace;
  word-break: break-all;
}

.download-time {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
}

.download-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex-shrink: 0;
}

.btn-action {
  padding: 4px 12px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.btn-cancel { background: #f39c12; color: white; }
.btn-start { background: #3498db; color: white; }
.btn-remove { background: #e74c3c; color: white; }

.empty-tip {
  font-size: 13px;
  color: #999;
  margin-top: 8px;
}
</style>
