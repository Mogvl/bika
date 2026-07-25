// TypeScript 类型定义

// 通用 API 响应
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分类
export interface Category {
  title: string
  thumb: string
  isActive: boolean
  link: string
}

// 漫画基本信息
export interface Comic {
  _id: string
  title: string
  author: string
  categories: string[]
  tags: string[]
  description: string
  thumb: {
    originalName: string
    path: string
    fileServer: string
  }
  epsCount: number
  pagesCount: number
  finished: boolean
  likesCount: number
  totalViews: number
  totalLikes: number
  commentsCount: number
  isFavourite: boolean
  isLiked: boolean
  updated_at: string
  created_at: string
}

// 漫画章节
export interface EP {
  _id: string
  title: string
  order: number
  updated_at: string
  pagesCount: number
}

// 漫画页面（图片）
export interface Page {
  id: number
  media: {
    originalName: string
    path: string
    fileServer: string
  }
}

// 用户信息
export interface UserProfile {
  _id: string
  email: string
  name: string
  avatar: {
    originalName: string
    path: string
    fileServer: string
  }
  level: number
  character: string
  title: string
  verified: boolean
  expired: boolean
}

// 排行榜时间范围
export type LeaderboardTT = 'H24' | 'D7' | 'D30'
