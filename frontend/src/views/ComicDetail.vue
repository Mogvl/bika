<template>
  <div class="detail-page" v-if="!loading">
    <!-- 头部信息 -->
    <div class="detail-header">
      <div class="detail-cover-wrap">
        <img :src="getCoverUrl(comic.thumb)" :alt="comic.title" class="detail-cover" />
      </div>
      <div class="detail-info">
        <h1 class="detail-title">{{ comic.title }}</h1>
        <p class="detail-author">作者: {{ comic.author }}</p>
        <p class="detail-meta">
          章节: {{ comic.epsCount }} | 喜欢: {{ comic.totalLikes || comic.likesCount }} | 浏览: {{ comic.totalViews || 0 }}
        </p>
        <div class="detail-tags">
          <span v-for="cat in comic.categories" :key="cat" class="tag">{{ cat }}</span>
        </div>
        <div class="detail-actions">
          <button class="btn btn-primary" @click="startReading" :disabled="!eps.length">
            {{ eps.length ? '开始阅读' : '暂无章节' }}
          </button>
          <button
            class="btn btn-fav"
            :class="{ 'is-fav': comic.isFavourite }"
            @click="toggleFavourite"
          >
            {{ comic.isFavourite ? '❤️ 已收藏' : '🤍 收藏' }}
          </button>
          <button class="btn btn-download" @click="downloadComic">📥 下载</button>
        </div>
        <div class="detail-actions" style="margin-top: 8px;">
          <button class="btn btn-like" :class="{ 'is-liked': comic.isLiked }" @click="toggleLike">
            {{ comic.isLiked ? '👍 已赞' : '👍 点赞' }} ({{ comic.totalLikes || comic.likesCount || 0 }})
          </button>
        </div>
      </div>
    </div>

    <!-- 简介 -->
    <div class="detail-section">
      <h3 class="section-title">简介</h3>
      <p class="detail-desc">{{ comic.description || '暂无简介' }}</p>
    </div>

    <!-- 章节列表 -->
    <div class="detail-section">
      <div class="eps-header">
        <h3 class="section-title">章节 ({{ eps.length }})</h3>
      </div>

      <div v-if="eps.length === 0" class="empty-state">暂无章节</div>

      <div v-else>
        <div class="eps-toolbar">
          <button class="eps-btn" @click="toggleSelectMode">
            {{ selectMode ? '取消选择' : '选择下载' }}
          </button>
          <template v-if="selectMode">
            <button class="eps-btn" @click="selectAll">全选</button>
            <button class="eps-btn" @click="selectNone">全不选</button>
            <button class="eps-btn eps-btn-download" @click="downloadSelected" :disabled="selectedEps.length === 0">
              📥 下载选中 ({{ selectedEps.length }})
            </button>
          </template>
        </div>

        <div class="eps-grid">
          <div
            v-for="ep in eps"
            :key="ep._id"
            class="eps-chip"
            :class="{
              'selected': selectedEps.includes(ep.order),
              'select-mode': selectMode
            }"
            @click="handleEpsClick(ep)"
          >
            <span v-if="selectMode" class="eps-check">
              {{ selectedEps.includes(ep.order) ? '☑' : '☐' }}
            </span>
            {{ ep.title || `第${ep.order}话` }}
          </div>
        </div>
      </div>

      <div v-if="epsHasMore" class="page-load-more">
        <button class="load-more-btn" @click="loadMoreEps">加载更多章节</button>
      </div>
    </div>

    <!-- 相关推荐 -->
    <div v-if="recommendations.length > 0" class="detail-section">
      <h3 class="section-title">相关推荐</h3>
      <div class="comic-grid">
        <div
          v-for="rec in recommendations"
          :key="rec._id"
          class="comic-card"
          @click="goComic(rec._id)"
        >
          <img :src="getCoverUrl(rec.thumb)" :alt="rec.title" class="comic-cover" loading="lazy" @error="handleImgError" />
          <div class="comic-info">
            <div class="comic-title">{{ rec.title }}</div>
            <div class="comic-author">{{ rec.author }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 评论区 -->
    <div class="detail-section">
      <div class="comments-header">
        <h3 class="section-title">评论 ({{ comic.commentsCount || commentsTotal || 0 }})</h3>
        <button class="comments-refresh" @click="loadComments(1)">🔄 刷新</button>
      </div>

      <!-- 发表评论 -->
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
        <p>还没有评论，快来抢沙发~</p>
      </div>

      <div v-else class="comments-list">
        <div v-for="c in comments" :key="c._id" class="comment-item">
          <div class="comment-avatar">
            <img v-if="c._user?.avatar" :src="getUserAvatar(c._user)" alt="" />
            <span v-else class="avatar-text">{{ c._user?.name?.[0] || '?' }}</span>
          </div>
          <div class="comment-body">
            <div class="comment-meta">
              <span class="comment-user">
                {{ c._user?.name || '匿名' }}
                <span v-if="c._user?.level" class="comment-level">Lv.{{ c._user.level }}</span>
              </span>
              <span class="comment-time">{{ formatTime(c.createdAt) }}</span>
            </div>
            <div class="comment-content" :class="{ 'comment-deleted': c.isDeleted }">
              {{ c.isDeleted ? '该评论已被删除' : c.content }}
            </div>
            <div class="comment-actions">
              <span
                class="comment-action"
                :class="{ active: c.liked }"
                @click="toggleCommentLike(c)"
              >
                👍 {{ c.likesCount || 0 }}
              </span>
              <span class="comment-action" @click="toggleSubComments(c)">💬 回复 ({{ c.childrenCount || subCountMap[c._id] || 0 }})</span>
              <span v-if="!c.isDeleted" class="comment-action report" @click="doReport(c)">举报</span>
            </div>

            <!-- 子评论（楼中楼） -->
            <div v-if="showSubs[c._id]" class="sub-comments">
              <div v-for="sub in subComments[c._id] || []" :key="sub._id" class="sub-comment">
                <div class="sub-avatar">
                  <img v-if="sub._user?.avatar" :src="getUserAvatar(sub._user)" alt="" />
                  <span v-else>{{ sub._user?.name?.[0] || '?' }}</span>
                </div>
                <div class="sub-body">
                  <div class="comment-meta">
                    <span class="comment-user">
                      {{ sub._user?.name || '匿名' }}
                      <span v-if="sub._user?.level" class="comment-level">Lv.{{ sub._user.level }}</span>
                    </span>
                    <span class="comment-time">{{ formatTime(sub.createdAt) }}</span>
                  </div>
                  <div class="comment-content" :class="{ 'comment-deleted': sub.isDeleted }">
                    {{ sub.isDeleted ? '该评论已被删除' : sub.content }}
                  </div>
                  <div class="comment-actions">
                    <span class="comment-action" :class="{ active: sub.liked }" @click="toggleSubLike(sub)">
                      👍 {{ sub.likesCount || 0 }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="sub-comment-input">
                <input
                  v-model="subInput[c._id]"
                  type="text"
                  placeholder="回复该评论..."
                  @keydown.enter="submitSubComment(c)"
                />
                <button class="btn-sub" @click="submitSubComment(c)" :disabled="!subInput[c._id]?.trim()">回复</button>
              </div>
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
  <div v-if="error" class="error-msg">{{ error }}</div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  getComicDetail, getComicEps, addDownload,
  getComments, sendComment, likeComment, likeComic, addFavourite,
  getSubComments, sendSubComment, reportComment, getComicRecommendation,
} from '@/api'
import { saveComicHistory } from '@/utils/history'
import type { Comic, EP } from '@/types'

const route = useRoute()
const router = useRouter()
const comic = ref<Comic>({} as Comic)
const eps = ref<EP[]>([])
const loading = ref(true)
const error = ref('')
const epsPage = ref(1)
const epsHasMore = ref(false)
const selectedEps = ref<number[]>([])
const selectMode = ref(false)

// 相关推荐
const recommendations = ref<any[]>([])

// 评论
const comments = ref<any[]>([])
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentsHasMore = ref(false)
const commentsLoading = ref(false)
const newComment = ref('')
const commentMsg = ref('')
const commentError = ref(false)
const showSubs = reactive<Record<string, boolean>>({})
const subComments = reactive<Record<string, any[]>>({})
const subCountMap = reactive<Record<string, number>>({})
const subInput = reactive<Record<string, string>>({})

onMounted(async () => {
  const id = route.params.id as string
  if (!id) {
    error.value = '参数错误'
    loading.value = false
    return
  }
  await loadDetail(id)
})

async function loadDetail(id: string) {
  try {
    const [detailRes, epsRes] = await Promise.all([
      getComicDetail(id),
      getComicEps(id, 1),
    ])
    comic.value = detailRes.data?.comic || {}
    const epsData = epsRes.data?.eps
    eps.value = Array.isArray(epsData?.docs) ? epsData.docs : (Array.isArray(epsData) ? epsData : [])
    epsHasMore.value = epsPage.value < (epsData?.pages || 1)
    saveComicHistory(comic.value)
    loadComments(1)
    loadRecommendations()
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

async function loadRecommendations() {
  try {
    const res = await getComicRecommendation(route.params.id as string)
    const data = res.data
    const recs = data?.comics
    recommendations.value = Array.isArray(recs?.docs) ? recs.docs : (Array.isArray(recs) ? recs : [])
  } catch {}
}

async function loadMoreEps() {
  epsPage.value++
  try {
    const res = await getComicEps(route.params.id as string, epsPage.value)
    const epsData = res.data?.eps
    const newEps = Array.isArray(epsData?.docs) ? epsData.docs : (Array.isArray(epsData) ? epsData : [])
    eps.value.push(...newEps)
    epsHasMore.value = epsPage.value < (epsData?.pages || 1)
  } catch {}
}

function handleEpsClick(ep: EP) {
  if (selectMode.value) {
    toggleEps(ep.order)
  } else {
    router.push(`/reader/${route.params.id}/${ep.order}`)
  }
}

function toggleSelectMode() {
  selectMode.value = !selectMode.value
  if (!selectMode.value) {
    selectedEps.value = []
  }
}

function toggleEps(order: number) {
  const idx = selectedEps.value.indexOf(order)
  if (idx >= 0) {
    selectedEps.value.splice(idx, 1)
  } else {
    selectedEps.value.push(order)
  }
}

function selectAll() {
  selectedEps.value = eps.value.map(ep => ep.order)
}

function selectNone() {
  selectedEps.value = []
}

function startReading() {
  if (eps.value.length > 0) {
    router.push(`/reader/${route.params.id}/${eps.value[0].order}`)
  }
}

async function downloadComic() {
  const coverUrl = comic.value.thumb?.fileServer && comic.value.thumb?.path
    ? `${comic.value.thumb.fileServer}/static/${comic.value.thumb.path}`
    : ''
  try {
    await addDownload(route.params.id as string, comic.value.title || '', coverUrl)
    alert('已添加到下载队列！')
  } catch (e: any) {
    alert(e.message || '添加下载失败')
  }
}

async function downloadSelected() {
  if (selectedEps.value.length === 0) return
  const coverUrl = comic.value.thumb?.fileServer && comic.value.thumb?.path
    ? `${comic.value.thumb.fileServer}/static/${comic.value.thumb.path}`
    : ''
  try {
    await addDownload(route.params.id as string, comic.value.title || '', coverUrl)
    alert(`已添加 ${selectedEps.value.length} 个章节到下载队列！`)
    selectMode.value = false
    selectedEps.value = []
  } catch (e: any) {
    alert(e.message || '添加下载失败')
  }
}

// ==================== 收藏 / 点赞 ====================
async function toggleFavourite() {
  try {
    await addFavourite(route.params.id as string)
    comic.value.isFavourite = !comic.value.isFavourite
    alert(comic.value.isFavourite ? '已加入收藏' : '已取消收藏')
  } catch (e: any) {
    alert(e.message || '操作失败')
  }
}

async function toggleLike() {
  try {
    await likeComic(route.params.id as string)
    comic.value.isLiked = !comic.value.isLiked
    const likes = comic.value.totalLikes || comic.value.likesCount || 0
    comic.value.totalLikes = comic.value.isLiked ? likes + 1 : Math.max(0, likes - 1)
  } catch (e: any) {
    alert(e.message || '点赞失败')
  }
}

// ==================== 评论 ====================
async function loadComments(page: number) {
  commentsLoading.value = true
  try {
    const res = await getComments(route.params.id as string, page)
    const data = res.data
    const docs = data?.comments?.docs || []
    commentsTotal.value = data?.comments?.total || 0
    comments.value = page === 1 ? docs : comments.value.concat(docs)
    commentsPage.value = page
    commentsHasMore.value = page < (data?.comments?.pages || 1)
  } catch (e: any) {
    commentMsg.value = e.message || '评论加载失败'
    commentError.value = true
  } finally {
    commentsLoading.value = false
  }
}

async function loadMoreComments() {
  await loadComments(commentsPage.value + 1)
}

async function submitComment() {
  const content = newComment.value.trim()
  if (!content) return
  commentMsg.value = ''
  try {
    await sendComment(route.params.id as string, content)
    newComment.value = ''
    commentMsg.value = '发表成功！'
    commentError.value = false
    await loadComments(1)
  } catch (e: any) {
    commentMsg.value = e.message || '发表失败'
    commentError.value = true
  }
}

async function toggleCommentLike(c: any) {
  try {
    await likeComment(c._id)
    c.liked = !c.liked
    c.likesCount = (c.likesCount || 0) + (c.liked ? 1 : -1)
    if (c.likesCount < 0) c.likesCount = 0
  } catch (e: any) {
    alert(e.message || '点赞失败')
  }
}

async function toggleSubLike(sub: any) {
  try {
    await likeComment(sub._id)
    sub.liked = !sub.liked
    sub.likesCount = (sub.likesCount || 0) + (sub.liked ? 1 : -1)
    if (sub.likesCount < 0) sub.likesCount = 0
  } catch (e: any) {
    alert(e.message || '点赞失败')
  }
}

async function toggleSubComments(c: any) {
  const cid = c._id
  if (showSubs[cid]) {
    showSubs[cid] = false
    return
  }
  showSubs[cid] = true
  if (!subComments[cid]) {
    await loadSubComments(cid)
  }
}

async function loadSubComments(commentId: string) {
  try {
    const res = await getSubComments(commentId, 1)
    const data = res.data
    const subs = data?.comments?.docs || data?.docs || []
    subComments[commentId] = Array.isArray(subs) ? subs : []
    subCountMap[commentId] = data?.comments?.total || data?.total || subComments[commentId].length
  } catch {}
}

async function submitSubComment(c: any) {
  const content = (subInput[c._id] || '').trim()
  if (!content) return
  try {
    await sendSubComment(c._id, content)
    subInput[c._id] = ''
    await loadSubComments(c._id)
    if (!c.childrenCount) c.childrenCount = 0
    c.childrenCount++
  } catch (e: any) {
    alert(e.message || '回复失败')
  }
}

async function doReport(c: any) {
  if (!confirm('确定举报该评论？')) return
  try {
    await reportComment(c._id)
    alert('举报成功')
  } catch (e: any) {
    alert(e.message || '举报失败')
  }
}

// ==================== 工具 ====================
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

function handleImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.background = '#f0f0f0'
}

function goComic(id: string) {
  router.push(`/comic/${id}`)
}
</script>

<style scoped>
.detail-page {
  max-width: 800px;
  margin: 0 auto;
}

.detail-header {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: var(--bg-card);
}

.detail-cover-wrap {
  flex-shrink: 0;
  width: 140px;
}

.detail-cover {
  width: 100%;
  border-radius: 8px;
  aspect-ratio: 3/4;
  object-fit: cover;
}

.detail-info {
  flex: 1;
  min-width: 0;
}

.detail-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  line-height: 1.3;
}

.detail-author {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.detail-meta {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.detail-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.tag {
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 11px;
  background: #fce4ec;
  color: var(--primary);
}

.detail-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.btn-fav {
  background: #ff6b81;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
}

.btn-fav.is-fav {
  background: #e74c3c;
}

.btn-like {
  background: var(--bg);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
}

.btn-like.is-liked {
  background: #fce4ec;
  border-color: var(--primary);
  color: var(--primary);
}

.btn-download {
  background: #3498db;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
}

.btn-download:hover {
  background: #2980b9;
}

.detail-section {
  padding: 16px;
  background: var(--bg-card);
  margin-top: 8px;
}

.detail-desc {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-secondary);
}

/* 章节网格 - 和原版一致 */
.eps-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.eps-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.eps-chip {
  padding: 8px 16px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.eps-chip:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: #fff0f5;
}

.eps-chip.selected {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.eps-chip.select-mode {
  padding: 8px 12px;
}

.eps-check {
  margin-right: 4px;
}

.eps-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.eps-btn {
  padding: 6px 14px;
  border: 1px solid var(--border);
  border-radius: 16px;
  font-size: 13px;
  cursor: pointer;
  background: var(--bg);
  color: var(--text);
  transition: all 0.2s;
}

.eps-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.eps-btn-download {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.eps-btn-download:hover {
  background: var(--primary-dark);
}

.eps-btn-download:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 评论区 */
.comments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.comments-refresh {
  padding: 4px 12px;
  border: 1px solid var(--border);
  border-radius: 16px;
  font-size: 12px;
  background: var(--bg);
  cursor: pointer;
  color: var(--text-secondary);
}

.comment-input-wrap {
  margin-bottom: 16px;
}

.comment-input {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  resize: vertical;
  box-sizing: border-box;
  font-family: inherit;
}

.comment-input:focus {
  border-color: var(--primary);
}

.comment-send-row {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
}

.comment-msg {
  font-size: 13px;
  color: #27ae60;
}

.comment-msg.error {
  color: #e74c3c;
}

.btn-comment-send {
  padding: 8px 24px;
}

.comments-list {
  display: flex;
  flex-direction: column;
}

.comment-item {
  display: flex;
  gap: 10px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 600;
}

.comment-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.comment-user {
  font-size: 13px;
  font-weight: 500;
}

.comment-level {
  font-size: 11px;
  color: var(--primary);
  background: #fce4ec;
  padding: 1px 6px;
  border-radius: 8px;
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

.comment-deleted {
  color: #999;
  font-style: italic;
}

.comment-actions {
  display: flex;
  gap: 16px;
  margin-top: 6px;
}

.comment-action {
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
}

.comment-action:hover {
  color: var(--primary);
}

.comment-action.active {
  color: var(--primary);
}

.comment-action.report {
  color: #999;
}

/* 子评论 */
.sub-comments {
  margin-top: 8px;
  padding: 8px;
  background: var(--bg);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sub-comment {
  display: flex;
  gap: 8px;
}

.sub-avatar {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: 600;
}

.sub-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.sub-body {
  flex: 1;
  min-width: 0;
}

.sub-comment-input {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.sub-comment-input input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 16px;
  font-size: 13px;
  outline: none;
}

.sub-comment-input input:focus {
  border-color: var(--primary);
}

.btn-sub {
  padding: 6px 14px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 16px;
  font-size: 13px;
  cursor: pointer;
}

.btn-sub:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
