<template>
  <div class="page-container">
    <div class="section-title">📥 下载管理</div>

    <div v-if="loading" class="loading">加载中</div>

    <div v-else-if="tasks.length === 0" class="empty-state">
      <p>暂无下载任务</p>
      <p class="empty-tip">在漫画详情页点击"下载"按钮添加任务</p>
    </div>

    <div v-else class="download-list">
      <div v-for="task in tasks" :key="task.id" class="download-item">
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
          </div>
          <div v-if="task.error" class="download-error">{{ task.error }}</div>
          <div class="download-path">📁 {{ getDownloadPath(task.title) }}</div>
          <div class="download-time">{{ formatTime(task.updatedAt) }}</div>
        </div>
        <div class="download-actions">
          <button v-if="task.status === 'downloading'" class="btn-action btn-cancel" @click="cancel(task.id)">暂停</button>
          <button v-if="task.status === 'paused' || task.status === 'error'" class="btn-action btn-start" @click="retry(task.id)">重试</button>
          <button class="btn-action btn-remove" @click="remove(task.id)">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { getDownloads, cancelDownload, removeDownload, addDownload } from '@/api'

interface DownloadTask {
  id: string
  bookId: string
  title: string
  coverUrl: string
  totalPages: number
  downloaded: number
  status: string
  error?: string
  createdAt: string
  updatedAt: string
}

const tasks = ref<DownloadTask[]>([])
const loading = ref(true)
let timer: number | null = null

onMounted(() => {
  loadTasks()
  // 每3秒刷新一次状态
  timer = window.setInterval(loadTasks, 3000)
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

async function remove(id: string) {
  try {
    await removeDownload(id)
    await loadTasks()
  } catch {}
}

async function retry(id: string) {
  // 重新添加任务
  const task = tasks.value.find(t => t.id === id)
  if (task) {
    await addDownload(task.bookId, task.title, task.coverUrl)
    await loadTasks()
  }
}

function getCoverUrl(url: string): string {
  if (!url) return ''
  if (url.includes('fileServer')) {
    // 解析旧格式
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

function getDownloadPath(title: string): string {
  const safeName = sanitizeTitle(title)
  return `/volume1/bika/${safeName}/`
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
