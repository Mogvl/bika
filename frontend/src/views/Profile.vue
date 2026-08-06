<template>
  <div class="page-container">
    <div class="section-title">👤 个人资料</div>

    <div v-if="loading" class="loading">加载中</div>
    <div v-else-if="user" class="profile-card">
      <div class="profile-header">
        <div class="profile-avatar" @click="triggerAvatarUpload">
          <img v-if="user.avatar" :src="getCoverUrl(user.avatar)" alt="头像" />
          <div v-else class="avatar-placeholder">{{ user.name?.[0] || '?' }}</div>
          <div class="avatar-overlay">📷</div>
        </div>
        <input ref="avatarInput" type="file" accept="image/*" style="display: none;" @change="onAvatarChange" />
        <div class="profile-info">
          <h2 class="profile-name">
            {{ user.name }}
            <button class="edit-title-btn" @click="showTitleModal = true" title="修改称号">✏️</button>
          </h2>
          <p class="profile-level">Lv.{{ user.level }}</p>
          <p class="profile-title">{{ user.title || '普通用户' }}</p>
        </div>
      </div>

      <div class="profile-stats">
        <div class="stat-item">
          <span class="stat-value">{{ user.exp || 0 }}</span>
          <span class="stat-label">经验</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ user.comicsDownloaded || 0 }}</span>
          <span class="stat-label">下载</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ user.likes || 0 }}</span>
          <span class="stat-label">获赞</span>
        </div>
      </div>

      <div class="profile-actions">
        <button class="btn btn-primary" @click="punchIn">签到</button>
        <router-link to="/my-comments" class="btn btn-ghost">我的评论</router-link>
        <router-link to="/settings" class="btn btn-ghost">设置</router-link>
      </div>
    </div>

    <!-- 修改称号弹窗 -->
    <div v-if="showTitleModal" class="modal-overlay" @click.self="showTitleModal = false">
      <div class="modal">
        <h3>修改称号</h3>
        <input v-model="titleInput" type="text" class="input-field" placeholder="输入新称号" maxlength="20" />
        <p v-if="titleMsg" class="error-tip">{{ titleMsg }}</p>
        <div class="modal-actions">
          <button class="btn" style="background: var(--bg);" @click="showTitleModal = false">取消</button>
          <button class="btn btn-primary" @click="saveTitle" :disabled="!titleInput.trim()">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { punchIn as apiPunchIn, setAvatar as apiSetAvatar, setTitle as apiSetTitle } from '@/api'

const auth = useAuthStore()
const user = ref<any>(null)
const loading = ref(true)

const avatarInput = ref<HTMLInputElement | null>(null)
const showTitleModal = ref(false)
const titleInput = ref('')
const titleMsg = ref('')

onMounted(async () => {
  try {
    await auth.fetchProfile()
    user.value = auth.user
  } catch {} finally { loading.value = false }
})

function getCoverUrl(thumb: any): string {
  if (!thumb?.fileServer || !thumb?.path) return ''
  return `/api/image/proxy?fileServer=${encodeURIComponent(thumb.fileServer)}&path=${encodeURIComponent(thumb.path)}`
}

async function punchIn() {
  try {
    await apiPunchIn()
    alert('签到成功！')
  } catch (e: any) { alert(e.message || '签到失败') }
}

// ==================== 修改头像 ====================
function triggerAvatarUpload() {
  avatarInput.value?.click()
}

function onAvatarChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // 压缩到合理大小
  const reader = new FileReader()
  reader.onload = () => {
    const img = new Image()
    img.onload = async () => {
      const canvas = document.createElement('canvas')
      const maxSize = 300
      let w = img.width
      let h = img.height
      if (w > maxSize || h > maxSize) {
        const ratio = Math.min(maxSize / w, maxSize / h)
        w = Math.round(w * ratio)
        h = Math.round(h * ratio)
      }
      canvas.width = w
      canvas.height = h
      const ctx = canvas.getContext('2d')
      if (!ctx) return
      ctx.drawImage(img, 0, 0, w, h)
      // 转 jpeg base64
      const base64 = canvas.toDataURL('image/jpeg', 0.9).split(',')[1]
      try {
        await apiSetAvatar(base64)
        alert('头像修改成功！')
        await auth.fetchProfile()
        user.value = auth.user
      } catch (err: any) {
        alert(err.message || '修改头像失败')
      }
    }
    img.src = reader.result as string
  }
  reader.readAsDataURL(file)
  input.value = ''
}

// ==================== 修改称号 ====================
async function saveTitle() {
  const title = titleInput.value.trim()
  if (!title) return
  titleMsg.value = ''
  try {
    await apiSetTitle(title)
    alert('称号修改成功！')
    showTitleModal.value = false
    titleInput.value = ''
    await auth.fetchProfile()
    user.value = auth.user
  } catch (e: any) {
    titleMsg.value = e.message || '修改称号失败'
  }
}
</script>

<style scoped>
.profile-card { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 26px; box-shadow: var(--shadow); }
.profile-header { display: flex; gap: 16px; margin-bottom: 24px; }
.profile-avatar { width: 84px; height: 84px; border-radius: 50%; overflow: hidden; flex-shrink: 0; position: relative; cursor: pointer; box-shadow: 0 0 0 3px rgba(255,95,168,0.25); }
.profile-avatar img { width: 100%; height: 100%; object-fit: cover; }
.avatar-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg,#ff5fa8,#ff3d94); color: white; font-size: 32px; font-weight: 700; }
.avatar-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; color: white; font-size: 22px; opacity: 0; transition: opacity 0.2s; border-radius: 50%; }
.profile-avatar:hover .avatar-overlay { opacity: 1; }
.profile-name { font-size: 20px; font-weight: 600; }
.edit-title-btn { background: none; border: none; font-size: 14px; cursor: pointer; margin-left: 6px; opacity: 0.7; }
.edit-title-btn:hover { opacity: 1; }
.profile-level { font-size: 13px; color: var(--primary); }
.profile-title { font-size: 12px; color: var(--text-muted); }
.profile-stats { display: flex; justify-content: space-around; padding: 18px 0; border-top: 1px solid var(--border); border-bottom: 1px solid var(--border); margin-bottom: 18px; }
.stat-item { text-align: center; }
.stat-value { display: block; font-size: 22px; font-weight: 700; color: var(--primary); font-variant-numeric: tabular-nums; }
.stat-label { font-size: 12px; color: var(--text-muted); }
.profile-actions { display: flex; gap: 12px; }
.profile-actions .btn { flex: 1; padding: 10px; text-align: center; }
.btn-ghost { background: var(--bg-elev); color: var(--text); border: 1px solid var(--border); }
.btn-ghost:hover { border-color: var(--primary); color: var(--primary); }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 300; }
.modal { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 26px; width: 90%; max-width: 360px; display: flex; flex-direction: column; gap: 12px; box-shadow: var(--shadow); }
.modal h3 { margin-bottom: 8px; font-family: var(--font-display); font-size: 18px; }
.error-tip { color: var(--danger); font-size: 13px; }
.modal-actions { display: flex; gap: 12px; margin-top: 8px; }
.modal-actions .btn { flex: 1; padding: 10px; }
</style>
