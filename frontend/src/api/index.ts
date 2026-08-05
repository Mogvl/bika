import axios from 'axios'
import type { ApiResponse, LeaderboardTT } from '@/types'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = token
  return config
})

api.interceptors.response.use(
  (response) => {
    const data = response.data as ApiResponse
    if (data.code !== 200) return Promise.reject(new Error(data.message || '请求失败'))
    return response
  },
  (error) => {
    // 登录接口的 401 不跳转，让调用方处理错误信息
    const isLoginRequest = error.config?.url?.includes('/auth/login')
    if (error.response?.status === 401 && !isLoginRequest) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/#/login'
    }
    const msg = error.response?.data?.message || error.message || '请求失败'
    return Promise.reject(new Error(msg))
  }
)

// ==================== 认证 ====================
export async function login(email: string, password: string) {
  const res = await api.post('/auth/login', { email, password })
  return res.data
}

export async function register(data: any) {
  const res = await api.post('/auth/register', data)
  return res.data
}

export async function getProfile() {
  const res = await api.get('/auth/profile')
  return res.data
}

export async function punchIn() {
  const res = await api.post('/auth/punch-in')
  return res.data
}

export async function changePassword(old_password: string, new_password: string) {
  const res = await api.post('/auth/change-password', { old_password, new_password })
  return res.data
}

export async function forgotPassword(email: string) {
  const res = await api.post('/auth/forgot-password', { email })
  return res.data
}

// 重置密码（忘记密码后通过密保问答重置）
export async function resetPassword(email: string, questionNo: number, answer: string) {
  const res = await api.post('/auth/reset-password', { email, questionNo, answer })
  return res.data
}

// 修改头像（接收 base64 图片数据）
export async function setAvatar(avatar: string) {
  const res = await api.post('/auth/avatar', { avatar })
  return res.data
}

// 修改称号
export async function setTitle(title: string) {
  const res = await api.post('/auth/title', { title })
  return res.data
}

// 我的评论列表
export async function getMyComments(page = 1) {
  const res = await api.get('/auth/my-comments', { params: { page } })
  return res.data
}

// ==================== 分类 ====================
export async function getCategories() {
  const res = await api.get('/categories')
  return res.data
}

// ==================== 漫画 ====================
export async function getComicsByCategory(page: number, category: string, sort = 'ua') {
  const res = await api.get('/comics', { params: { page, c: category, s: sort } })
  return res.data
}

export async function searchComics(keyword: string, page = 1, category = '', sort = 'ua') {
  const res = await api.get('/comics/search', { params: { keyword, page, c: category, s: sort } })
  return res.data
}

export async function getComicDetail(id: string) {
  const res = await api.get(`/comics/${id}`)
  return res.data
}

export async function getComicEps(id: string, page = 1) {
  const res = await api.get(`/comics/${id}/eps`, { params: { page } })
  return res.data
}

export async function getComicPages(id: string, epsId: string, page = 1) {
  const res = await api.get(`/comics/${id}/eps/${epsId}/pages`, { params: { page } })
  return res.data
}

export async function getLeaderboard(tt: LeaderboardTT = 'H24') {
  const res = await api.get('/comics/leaderboard', { params: { tt } })
  return res.data
}

export async function getKnightLeaderboard() {
  const res = await api.get('/comics/knight-leaderboard')
  return res.data
}

export async function getRandomComics() {
  const res = await api.get('/comics/random')
  return res.data
}

export async function getCollections() {
  const res = await api.get('/comics/collections')
  return res.data
}

export async function getKeywords() {
  const res = await api.get('/comics/keywords')
  return res.data
}

// ==================== 评论 ====================
export async function getComments(bookId: string, page = 1) {
  const res = await api.get(`/comics/${bookId}/comments`, { params: { page } })
  return res.data
}

export async function sendComment(bookId: string, content: string) {
  const res = await api.post(`/comics/${bookId}/comments/send`, { content })
  return res.data
}

export async function likeComment(commentId: string) {
  const res = await api.post(`/comments/${commentId}/like`)
  return res.data
}

export async function likeComic(bookId: string) {
  const res = await api.post(`/comics/${bookId}/like`)
  return res.data
}

// 子评论（楼中楼）
export async function getSubComments(commentId: string, page = 1) {
  const res = await api.get(`/comments/${commentId}/childrens`, { params: { page } })
  return res.data
}

export async function sendSubComment(commentId: string, content: string) {
  const res = await api.post(`/comments/${commentId}/childrens/send`, { content })
  return res.data
}

export async function reportComment(commentId: string) {
  const res = await api.post(`/comments/${commentId}/report`)
  return res.data
}

// 相关推荐
export async function getComicRecommendation(id: string) {
  const res = await api.get(`/comics/${id}/recommendation`)
  return res.data
}

// ==================== 收藏 ====================
export async function getFavourites(page = 1, sort = 'da') {
  const res = await api.get('/favourites', { params: { page, s: sort } })
  return res.data
}

export async function addFavourite(id: string) {
  const res = await api.post(`/comics/${id}/favourite`)
  return res.data
}

// ==================== 游戏 ====================
export async function getGames(page = 1) {
  const res = await api.get('/games', { params: { page } })
  return res.data
}

export async function getGameDetail(id: string) {
  const res = await api.get(`/games/${id}`)
  return res.data
}

export async function getGameEps(id: string, page = 1) {
  const res = await api.get(`/games/${id}/eps`, { params: { page } })
  return res.data
}

export async function getGamePages(id: string, epsId: string, page = 1) {
  const res = await api.get(`/games/${id}/eps/${epsId}/pages`, { params: { page } })
  return res.data
}

export async function getGameComments(id: string, page = 1) {
  const res = await api.get(`/games/${id}/comments`, { params: { page } })
  return res.data
}

export async function sendGameComment(id: string, content: string) {
  const res = await api.post(`/games/${id}/comments/send`, { content })
  return res.data
}

export async function likeGameComment(commentId: string) {
  const res = await api.post(`/game-comments/${commentId}/like`)
  return res.data
}

// ==================== 图片 ====================
export function getImageUrl(fileServer: string, path: string): string {
  return `/api/image/proxy?fileServer=${encodeURIComponent(fileServer)}&path=${encodeURIComponent(path)}`
}

// ==================== 下载 ====================
export async function getDownloads() {
  const res = await api.get('/downloads')
  return res.data
}

export async function addDownload(bookId: string, title: string, coverUrl: string) {
  const res = await api.post('/downloads/add', { bookId, title, coverUrl })
  return res.data
}

export async function getDownloadStatus(id: string) {
  const res = await api.get(`/downloads/${id}`)
  return res.data
}

export async function cancelDownload(id: string) {
  const res = await api.post(`/downloads/${id}/cancel`)
  return res.data
}

export async function removeDownload(id: string) {
  const res = await api.post(`/downloads/${id}/remove`)
  return res.data
}

// ==================== 聊天 ====================
export async function chatLogin(email: string, password: string) {
  const res = await api.post('/chat/login', { email, password })
  return res.data
}

export async function getChatRooms() {
  const res = await api.get('/chat/rooms')
  return res.data
}

export async function getChatMessages(roomId: string, page = 1) {
  const res = await api.get('/chat/messages', { params: { roomId, page } })
  return res.data
}

export async function sendChatMessage(roomId: string, message: string) {
  const res = await api.post('/chat/send', { roomId, message })
  return res.data
}

export async function getChatProfile() {
  const res = await api.get('/chat/profile')
  return res.data
}

// ==================== 锅贴（好友动态） ====================
export async function getFriedPosts(page = 1) {
  const res = await api.get('/fried/posts', { params: { page } })
  return res.data
}

export async function getFriedComments(postId: string, page = 1) {
  const res = await api.get(`/fried/posts/${postId}/comments`, { params: { page } })
  return res.data
}

export async function sendFriedComment(postId: string, content: string) {
  const res = await api.post(`/fried/posts/${postId}/comments/send`, { content })
  return res.data
}

export async function likeFriedComment(commentId: string) {
  const res = await api.post(`/fried/comments/${commentId}/like`)
  return res.data
}
