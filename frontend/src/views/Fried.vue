<template>
  <div class="page-container">
    <div class="section-title">📱 锅贴广场</div>

    <div v-if="loading" class="loading">加载中</div>

    <div v-else-if="posts.length === 0" class="empty-state">
      <p>暂无动态</p>
    </div>

    <div v-else class="posts-list">
      <div v-for="post in posts" :key="post._id" class="post-item">
        <div class="post-header">
          <div class="post-avatar">
            <span>{{ post._user?.name?.[0] || '?' }}</span>
          </div>
          <div class="post-user">
            <div class="post-username">{{ post._user?.name || '匿名' }}</div>
            <div class="post-time">{{ formatTime(post.createdAt) }}</div>
          </div>
        </div>

        <div class="post-content">{{ post.content }}</div>

        <div v-if="post.medias?.length" class="post-media">
          <div v-for="media in post.medias" :key="media._id" class="media-item">
            <img :src="media.path" alt="" loading="lazy" @error="handleImgError" />
          </div>
        </div>

        <div class="post-footer">
          <span class="post-stat" @click="viewComments(post._id)">💬 {{ post.totalComments }}</span>
          <span class="post-stat">❤️ {{ post.totalLikes }}</span>
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
import { getDownloads } from '@/api'

interface FriedPost {
  _id: string
  content: string
  _user: {
    _id: string
    name: string
    level: number
    title: string
    character: string
    avatar: string
  }
  medias: Array<{ _id: string; path: string }>
  createdAt: string
  totalComments: number
  totalLikes: number
  liked: boolean
}

const router = useRouter()
const posts = ref<FriedPost[]>([])
const page = ref(1)
const hasMore = ref(true)
const loading = ref(true)

onMounted(() => loadPosts())

async function loadPosts() {
  loading.value = true
  try {
    const res = await getDownloads() // TODO: Replace with actual fried API
    // For now, show empty
    posts.value = []
  } catch {} finally { loading.value = false }
}

async function loadMore() {
  page.value++
  await loadPosts()
}

function viewComments(postId: string) {
  // TODO: Implement comment view
  alert('评论功能开发中')
}

function formatTime(t: string): string {
  if (!t) return ''
  const date = new Date(t)
  return date.toLocaleString('zh-CN')
}

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}
</script>

<style scoped>
.posts-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.post-item {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 16px;
  box-shadow: var(--shadow);
}

.post-header {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.post-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}

.post-username {
  font-size: 14px;
  font-weight: 500;
}

.post-time {
  font-size: 12px;
  color: var(--text-secondary);
}

.post-content {
  font-size: 15px;
  line-height: 1.6;
  margin-bottom: 12px;
}

.post-media {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  margin-bottom: 12px;
}

.media-item img {
  width: 100%;
  border-radius: 8px;
  object-fit: cover;
}

.post-footer {
  display: flex;
  gap: 24px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.post-stat {
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
}

.post-stat:hover {
  color: var(--primary);
}
</style>
