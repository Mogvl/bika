<template>
  <div class="page-container">
    <div class="section-title">💬 我的评论</div>

    <div v-if="loading" class="loading">加载中</div>

    <div v-else-if="comments.length === 0" class="empty-state">
      <p>你还没有发表过评论</p>
    </div>

    <div v-else class="comments-list">
      <div v-for="c in comments" :key="c._id" class="comment-item">
        <div class="comment-body">
          <div class="comment-meta">
            <span class="comment-comic" @click="goComic(c._comic?._id)">
              {{ c._comic?.title || '未知漫画' }}
            </span>
            <span class="comment-time">{{ formatTime(c.createdAt) }}</span>
          </div>
          <div class="comment-content">{{ c.content }}</div>
          <div class="comment-actions">
            <span class="comment-stat">👍 {{ c.likesCount || 0 }}</span>
            <span v-if="c._comic?._id" class="comment-stat link" @click="goComic(c._comic._id)">查看漫画</span>
          </div>
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
import { getMyComments } from '@/api'

const router = useRouter()
const comments = ref<any[]>([])
const page = ref(1)
const hasMore = ref(false)
const loading = ref(true)

onMounted(() => loadComments())

async function loadComments() {
  loading.value = true
  try {
    const res = await getMyComments(page.value)
    const data = res.data
    const docs = data?.comments?.docs || data?.docs || []
    comments.value = page.value === 1 ? docs : comments.value.concat(docs)
    hasMore.value = page.value < (data?.comments?.pages || data?.pages || 1)
  } catch {} finally { loading.value = false }
}

async function loadMore() {
  page.value++
  await loadComments()
}

function formatTime(t: string): string {
  if (!t) return ''
  const date = new Date(t)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / 86400000)
  if (days > 0) return `${days}天前`
  const hours = Math.floor(diff / 3600000)
  if (hours > 0) return `${hours}小时前`
  const minutes = Math.floor(diff / 60000)
  if (minutes > 0) return `${minutes}分钟前`
  return '刚刚'
}

function goComic(id: string) {
  if (id) router.push(`/comic/${id}`)
}
</script>

<style scoped>
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-item {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 16px;
  box-shadow: var(--shadow);
}

.comment-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.comment-comic {
  font-size: 13px;
  font-weight: 500;
  color: var(--primary);
  cursor: pointer;
}

.comment-time {
  font-size: 11px;
  color: #999;
}

.comment-content {
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
}

.comment-actions {
  display: flex;
  gap: 16px;
  margin-top: 8px;
}

.comment-stat {
  font-size: 12px;
  color: var(--text-secondary);
}

.comment-stat.link {
  cursor: pointer;
  color: var(--primary);
}
</style>
