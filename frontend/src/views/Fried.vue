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
            <img :src="getMediaUrl(media)" alt="" loading="lazy" @error="handleImgError" />
          </div>
        </div>

        <div class="post-footer">
          <span class="post-stat" @click="toggleComments(post)">💬 {{ post.totalComments }}</span>
          <span class="post-stat" :class="{ active: post.liked }" @click="togglePostLike(post)">
            ❤️ {{ post.totalLikes }}
          </span>
        </div>

        <!-- 评论展开 -->
        <div v-if="showComments[post._id]" class="post-comments">
          <div v-if="commentsLoading[post._id]" class="loading">加载中</div>
          <div v-else-if="(comments[post._id] || []).length === 0" class="sub-empty">暂无评论</div>
          <div v-else class="sub-list">
            <div v-for="c in comments[post._id] || []" :key="c._id" class="sub-comment">
              <div class="sub-body">
                <div class="sub-user">{{ c._user?.name || '匿名' }}</div>
                <div class="sub-text">{{ c.content }}</div>
                <div class="sub-actions">
                  <span class="post-stat small" :class="{ active: c.liked }" @click="toggleCommentLike(post, c)">
                    👍 {{ c.totalLikes || 0 }}
                  </span>
                </div>
              </div>
            </div>
          </div>
          <div class="sub-input">
            <input
              v-model="commentInput[post._id]"
              type="text"
              placeholder="发表评论..."
              @keydown.enter="sendComment(post)"
            />
            <button class="btn-sub" @click="sendComment(post)" :disabled="!commentInput[post._id]?.trim()">发送</button>
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
import { ref, reactive, onMounted } from 'vue'
import { getFriedPosts, getFriedComments, sendFriedComment, likeFriedComment } from '@/api'

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

const posts = ref<FriedPost[]>([])
const page = ref(1)
const hasMore = ref(true)
const loading = ref(true)

// 评论状态
const showComments = reactive<Record<string, boolean>>({})
const comments = reactive<Record<string, any[]>>({})
const commentsLoading = reactive<Record<string, boolean>>({})
const commentInput = reactive<Record<string, string>>({})

onMounted(() => loadPosts())

async function loadPosts() {
  loading.value = true
  try {
    const res = await getFriedPosts(page.value)
    const data = res.data
    const newPosts = data?.posts || []
    posts.value = page.value === 1 ? newPosts : posts.value.concat(newPosts)
    hasMore.value = posts.value.length < (data?.total || 0)
  } catch {} finally { loading.value = false }
}

async function loadMore() {
  page.value++
  await loadPosts()
}

async function toggleComments(post: FriedPost) {
  const id = post._id
  if (showComments[id]) {
    showComments[id] = false
    return
  }
  showComments[id] = true
  if (!comments[id]) {
    await loadComments(id)
  }
}

async function loadComments(postId: string) {
  commentsLoading[postId] = true
  try {
    const res = await getFriedComments(postId, 1)
    comments[postId] = res.data?.comments || []
  } catch {} finally { commentsLoading[postId] = false }
}

async function sendComment(post: FriedPost) {
  const content = (commentInput[post._id] || '').trim()
  if (!content) return
  try {
    await sendFriedComment(post._id, content)
    commentInput[post._id] = ''
    post.totalComments++
    await loadComments(post._id)
  } catch (e: any) {
    alert(e.message || '发送失败')
  }
}

async function toggleCommentLike(post: FriedPost, c: any) {
  try {
    await likeFriedComment(c._id)
    c.liked = !c.liked
    c.totalLikes = (c.totalLikes || 0) + (c.liked ? 1 : -1)
    if (c.totalLikes < 0) c.totalLikes = 0
  } catch (e: any) {
    alert(e.message || '点赞失败')
  }
}

function togglePostLike(post: FriedPost) {
  // 锅贴帖子本身点赞 API 未提供，仅做本地状态展示
  post.liked = !post.liked
  post.totalLikes += post.liked ? 1 : -1
}

function getMediaUrl(media: any): string {
  if (!media?.path) return ''
  // 外部图片地址直接使用；若是 fileServer+path 结构则走代理
  if (media.path.startsWith('http')) return media.path
  if (media.fileServer && media.path) {
    return `/api/image/proxy?fileServer=${encodeURIComponent(media.fileServer)}&path=${encodeURIComponent(media.path)}`
  }
  return `/api/image/proxy?url=${encodeURIComponent(media.path)}`
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
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.75);
  border-radius: var(--radius-sm);
  padding: 16px;
  box-shadow: var(--shadow-sm);
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
  white-space: pre-wrap;
  word-break: break-word;
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
  aspect-ratio: 1;
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

.post-stat.active {
  color: var(--primary);
}

.post-stat.small {
  font-size: 12px;
}

/* 评论 */
.post-comments {
  margin-top: 12px;
  padding: 12px;
  background: var(--bg-solid);
  border-radius: 8px;
}

.sub-empty {
  color: var(--text-secondary);
  font-size: 13px;
  padding: 8px 0;
  text-align: center;
}

.sub-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sub-comment {
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
}

.sub-comment:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.sub-user {
  font-size: 12px;
  font-weight: 500;
  color: var(--primary);
  margin-bottom: 2px;
}

.sub-text {
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.sub-actions {
  margin-top: 4px;
}

.sub-input {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.sub-input input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-pill);
  font-size: 13px;
  outline: none;
  background: rgba(255, 255, 255, 0.6);
}

.sub-input input:focus {
  border-color: var(--primary);
}

.btn-sub {
  padding: 6px 16px;
  background: linear-gradient(120deg, #b98ce8, #a475d0);
  color: white;
  border: none;
  border-radius: var(--radius-pill);
  font-size: 13px;
  cursor: pointer;
}

.btn-sub:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
