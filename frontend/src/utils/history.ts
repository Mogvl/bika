// 阅读历史管理

interface HistoryItem {
  bookId: string
  title: string
  author: string
  thumb: {
    fileServer: string
    path: string
  }
  epsId: string
  epsOrder: number
  epsTitle: string
  picIndex: number
  lastRead: number
}

const STORAGE_KEY = 'read_history'
const MAX_HISTORY = 200

// 获取所有历史
export function getHistory(): HistoryItem[] {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved) {
    try { return JSON.parse(saved) } catch {}
  }
  return []
}

// 保存历史
export function saveHistory(item: HistoryItem) {
  const history = getHistory()

  // 检查是否已存在
  const existingIdx = history.findIndex(h => h.bookId === item.bookId)
  if (existingIdx >= 0) {
    // 更新已存在的记录
    history[existingIdx] = { ...history[existingIdx], ...item, lastRead: Date.now() }
  } else {
    // 添加新记录到开头
    history.unshift({ ...item, lastRead: Date.now() })
  }

  // 限制数量
  if (history.length > MAX_HISTORY) {
    history.splice(MAX_HISTORY)
  }

  localStorage.setItem(STORAGE_KEY, JSON.stringify(history))
}

// 保存漫画浏览历史（从详情页）
export function saveComicHistory(comic: any) {
  if (!comic?._id) return
  saveHistory({
    bookId: comic._id,
    title: comic.title || '',
    author: comic.author || '',
    thumb: comic.thumb || { fileServer: '', path: '' },
    epsId: '',
    epsOrder: 0,
    epsTitle: '',
    picIndex: 0,
    lastRead: Date.now(),
  })
}

// 保存阅读进度
export function saveReadingProgress(bookId: string, title: string, thumb: any, epsId: string, epsOrder: number, epsTitle: string, picIndex: number) {
  const existing = getHistory().find(h => h.bookId === bookId)
  // 传入 thumb 为空时保留历史里已有的封面信息
  const effectiveThumb = (thumb && (thumb.fileServer || thumb.path))
    ? thumb
    : (existing?.thumb || { fileServer: '', path: '' })
  saveHistory({
    bookId,
    title,
    author: existing?.author || '',
    thumb: effectiveThumb,
    epsId,
    epsOrder,
    epsTitle,
    picIndex,
    lastRead: Date.now(),
  })
}

// 删除历史
export function removeHistory(bookId: string) {
  const history = getHistory().filter(h => h.bookId !== bookId)
  localStorage.setItem(STORAGE_KEY, JSON.stringify(history))
}

// 清空历史
export function clearHistory() {
  localStorage.removeItem(STORAGE_KEY)
}
