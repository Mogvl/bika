<template>
  <div class="detail-page" v-if="!loading">
    <div class="detail-header">
      <div class="detail-cover-wrap">
        <img :src="getCoverUrl(game.thumb)" :alt="game.title" class="detail-cover" />
      </div>
      <div class="detail-info">
        <h1 class="detail-title">{{ game.title }}</h1>
        <p class="detail-author">厂商: {{ game.author }}</p>
        <p class="detail-meta">章节: {{ game.epsCount }} | 喜欢: {{ game.totalLikes }}</p>
        <div class="detail-tags">
          <span v-for="cat in game.categories" :key="cat" class="tag">{{ cat }}</span>
        </div>
        <button class="btn btn-primary" @click="startPlaying" :disabled="!eps.length">
          {{ eps.length ? '开始游戏' : '暂无章节' }}
        </button>
      </div>
    </div>
    <div class="detail-section">
      <h3 class="section-title">简介</h3>
      <p class="detail-desc">{{ game.description || '暂无简介' }}</p>
    </div>
    <div class="detail-section">
      <h3 class="section-title">章节列表 ({{ eps.length }})</h3>
      <div class="eps-list">
        <div v-for="ep in eps" :key="ep.order" class="eps-item" @click="goReader(ep.order)">
          <span class="eps-order">第 {{ ep.order }} 话</span>
          <span class="eps-title">{{ ep.title || `第 ${ep.order} 话` }}</span>
        </div>
      </div>
    </div>

    <!-- 游戏评论 -->
    <div class="detail-section">
      <div class="comments-header">
        <h3 class="section-title">评论 ({{ commentsTotal }})</h3>
        <button class="comments-refresh" @click="loadComments(1)">🔄 刷新</button>
      </div>

      <div class="comment-input-wrap">
        <textarea
          v-model="newComment"
          class="comment-input"
          placeholder="发表你的评论..."
          rows="2"
        ></textarea>
        <div class="comment-send-row">
          <span v-if="commentMsg" class="comment-msg" :class="{ error: commentError }">{{ commentMsg }}</span>
          <button class="btn btn-primary btn-comment-send" @click="submitComment" :disabled="!newComment.trim()">
            发表评论
          </button>
        </div>
      </div>

      <div v-if="commentsLoading && comments.length === 0" class="loading">加载中</div>
      <div v-else-if="comments.length === 0" class="empty-state">
        <p>还没有评论</p>
      </div>

      <div v-else class="comments-list">
        <div v-for="c in comments" :key="c._id" class="comment-item">
          <div class="comment-avatar">
            <img v-if="c._user?.avatar" :src="getUserAvatar(c._user)" alt="" />
            <span v-else class="avatar-text">{{ c._user?.name?.[0] || '?' }}</span>
          </div>
          <div class="comment-body">
            <div class="comment-meta">
              <span class="comment-user">{{ c._user?.name || '匿名' }}</span>
              <span class="comment-time">{{ formatTime(c.createdAt) }}</span>
            </div>
            <div class="comment-content">{{ c.content }}</div>
            <div class="comment-actions">
              <span class="comment-action" :class="{ active: c.liked }" @click="toggleLike(c)">
                👍 {{ c.likesCount || 0 }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="commentsHasMore" class="page-load-more">
        <button class="load-more-btn" @click="loadMoreComments">加载更多评论</button>
      </div>
    </div>
  </div>
  <div v-else class="loading">加载中</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getGameDetail, getGameEps, getGameComments, sendGameComment, likeGameComment } from '@/api'

const route = useRoute()
const router = useRouter()
const game = ref<any>({})
const eps = ref<any[]>([])
const loading = ref(true)

// 评论
const comments = ref<any[]>([])
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentsHasMore = ref(false)
const commentsLoading = ref(false)
const newComment = ref('')
const commentMsg = ref('')
const commentError = ref(false)

onMounted(async () => {
  const id = route.params.id as string
  try {
    const [detailRes, epsRes] = await Promise.all([getGameDetail(id), getGameEps(id)])
    game.value = detailRes.data?.game || {}
    const epsData = epsRes.data?.eps
    eps.value = Array.isArray(epsData?.docs) ? epsData.docs : (Array.isArray(epsData) ? epsData : [])
    loadComments(1)
  } catch {} finally { loading.value = false }
})

async function loadComments(page: number) {
  commentsLoading.value = true
  try {
    const res = await getGameComments(route.params.id as string, page)
    const data = res.data
    const docs = data?.comments?.docs || []
    comments.value = page === 1 ? docs : comments.value.concat(docs)
    commentsTotal.value = data?.comments?.total || 0
    commentsPage.value = page
    commentsHasMore.value = page < (data?.comments?.pages || 1)
  } catch {} finally { commentsLoading.value = false }
}

async function loadMoreComments() {
  await loadComments(commentsPage.value + 1)
}

async function submitComment() {
  const content = newComment.value.trim()
  if (!content) return
  commentMsg.value = ''
  try {
    await sendGameComment(route.params.id as string, content)
    newComment.value = ''
    commentMsg.value = '发表成功！'
    commentError.value = false
    await loadComments(1)
  } catch (e: any) {
    commentMsg.value = e.message || '发表失败'
    commentError.value = true
  }
}

async function toggleLike(c: any) {
  try {
    await likeGameComment(c._id)
    c.liked = !c.liked
    c.likesCount = (c.likesCount || 0) + (c.liked ? 1 : -1)
    if (c.likesCount < 0) c.likesCount = 0
  } catch (e: any) {
    alert(e.message || '点赞失败')
  }
}

function getUserAvatar(user: any): string {
  if (user?.avatar?.fileServer && user?.avatar?.path) {
    return `/api/image/proxy?fileServer=${encodeURIComponent(user.avatar.fileServer)}&path=${encodeURIComponent(user.avatar.path)}`
  }
  if (typeof user?.avatar === 'string' && user.avatar) {
    return `/api/image/proxy?url=${encodeURIComponent(user.avatar)}`
  }
  return ''
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

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

function startPlaying() { if (eps.value.length) goReader(eps.value[0].order) }
function goReader(order: number) { router.push({ path: `/reader/${route.params.id}/${order}`, query: { type: 'game' } }) }

function goBack() {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/games')
  }
}
</script>

<style scoped>
.detail-page { max-width: 800px; margin: 0 auto; }
.detail-header { display: flex; gap: 16px; padding: 16px; background: var(--bg-card); }
.detail-cover-wrap { flex-shrink: 0; width: 140px; }
.detail-cover { width: 100%; border-radius: 8px; aspect-ratio: 3/4; object-fit: cover; }
.detail-info { flex: 1; min-width: 0; }
.detail-title { font-size: 18px; font-weight: 600; margin-bottom: 8px; }
.detail-author { font-size: 13px; color: var(--text-secondary); margin-bottom: 4px; }
.detail-meta { font-size: 12px; color: var(--text-secondary); margin-bottom: 8px; }
.detail-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 12px; }
.tag { padding: 2px 10px; border-radius: 12px; font-size: 11px; background: rgba(255,95,168,0.14); color: var(--primary); }
.detail-section { padding: 16px; background: var(--bg-card); margin-top: 8px; }
.detail-desc { font-size: 14px; line-height: 1.7; color: var(--text-secondary); }
.eps-list { display: flex; flex-direction: column; }
.eps-item { display: flex; align-items: center; padding: 12px 8px; border-bottom: 1px solid var(--border); cursor: pointer; gap: 12px; }
.eps-item:hover { background: var(--bg); }
.eps-order { font-size: 13px; color: var(--primary); font-weight: 500; min-width: 56px; }
.eps-title { flex: 1; font-size: 14px; }

/* 评论 */
.comments-header { display: flex; justify-content: space-between; align-items: center; }
.comments-refresh { padding: 4px 12px; border: 1px solid var(--border); border-radius: 16px; font-size: 12px; background: var(--bg); cursor: pointer; color: var(--text-secondary); }
.comment-input-wrap { margin-bottom: 16px; }
.comment-input { width: 100%; padding: 12px; border: 1px solid var(--border); border-radius: 8px; font-size: 14px; outline: none; resize: vertical; box-sizing: border-box; font-family: inherit; }
.comment-input:focus { border-color: var(--primary); }
.comment-send-row { display: flex; justify-content: flex-end; align-items: center; gap: 12px; margin-top: 8px; }
.comment-msg { font-size: 13px; color: var(--success); }
.comment-msg.error { color: var(--danger); }
.btn-comment-send { padding: 8px 24px; }
.comments-list { display: flex; flex-direction: column; }
.comment-item { display: flex; gap: 10px; padding: 12px 0; border-bottom: 1px solid var(--border); }
.comment-avatar { width: 36px; height: 36px; border-radius: 50%; overflow: hidden; flex-shrink: 0; background: var(--primary); display: flex; align-items: center; justify-content: center; color: white; font-size: 14px; font-weight: 600; }
.comment-avatar img { width: 100%; height: 100%; object-fit: cover; }
.comment-body { flex: 1; min-width: 0; }
.comment-meta { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.comment-user { font-size: 13px; font-weight: 500; }
.comment-time { font-size: 11px; color: var(--text-muted); }
.comment-content { font-size: 14px; line-height: 1.6; word-break: break-word; }
.comment-actions { display: flex; gap: 16px; margin-top: 6px; }
.comment-action { font-size: 12px; color: var(--text-secondary); cursor: pointer; }
.comment-action:hover { color: var(--primary); }
.comment-action.active { color: var(--primary); }
</style>
