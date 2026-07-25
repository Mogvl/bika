import axios from 'axios'
import type { ApiResponse, LeaderboardTT } from '@/types'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 自动添加 token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = token
  }
  return config
})

// 响应拦截器 - 统一错误处理
api.interceptors.response.use(
  (response) => {
    const data = response.data as ApiResponse
    if (data.code !== 200) {
      return Promise.reject(new Error(data.message || '请求失败'))
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/#/login'
    }
    return Promise.reject(error)
  }
)

// ==================== 认证相关 ====================

export async function login(email: string, password: string) {
  const res = await api.post('/auth/login', { email, password })
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

// ==================== 分类相关 ====================

export async function getCategories() {
  const res = await api.get('/categories')
  return res.data
}

// ==================== 漫画相关 ====================

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

// ==================== 收藏相关 ====================

export async function getFavourites(page = 1, sort = 'da') {
  const res = await api.get('/favourites', { params: { page, s: sort } })
  return res.data
}

export async function addFavourite(id: string) {
  const res = await api.post(`/comics/${id}/favourite`)
  return res.data
}

// ==================== 图片代理 ====================

export function getImageUrl(fileServer: string, path: string): string {
  return `/api/image/proxy?fileServer=${encodeURIComponent(fileServer)}&path=${encodeURIComponent(path)}`
}

export function getImageUrlDirect(url: string): string {
  return `/api/image/proxy?url=${encodeURIComponent(url)}`
}
