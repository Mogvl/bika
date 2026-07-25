<template>
  <div class="page-container">
    <div class="section-title">👤 个人资料</div>

    <div v-if="loading" class="loading">加载中</div>
    <div v-else-if="user" class="profile-card">
      <div class="profile-header">
        <div class="profile-avatar">
          <img v-if="user.avatar" :src="getCoverUrl(user.avatar)" alt="头像" />
          <div v-else class="avatar-placeholder">{{ user.name?.[0] || '?' }}</div>
        </div>
        <div class="profile-info">
          <h2 class="profile-name">{{ user.name }}</h2>
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
        <router-link to="/settings" class="btn" style="background: var(--bg); color: var(--text);">设置</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { punchIn as apiPunchIn } from '@/api'

const auth = useAuthStore()
const user = ref<any>(null)
const loading = ref(true)

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
</script>

<style scoped>
.profile-card { background: var(--bg-card); border-radius: var(--radius); padding: 24px; box-shadow: var(--shadow); }
.profile-header { display: flex; gap: 16px; margin-bottom: 24px; }
.profile-avatar { width: 80px; height: 80px; border-radius: 50%; overflow: hidden; flex-shrink: 0; }
.profile-avatar img { width: 100%; height: 100%; object-fit: cover; }
.avatar-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; background: var(--primary); color: white; font-size: 32px; font-weight: 700; }
.profile-name { font-size: 20px; font-weight: 600; }
.profile-level { font-size: 13px; color: var(--primary); }
.profile-title { font-size: 12px; color: var(--text-secondary); }
.profile-stats { display: flex; justify-content: space-around; padding: 16px 0; border-top: 1px solid var(--border); border-bottom: 1px solid var(--border); margin-bottom: 16px; }
.stat-item { text-align: center; }
.stat-value { display: block; font-size: 20px; font-weight: 600; color: var(--primary); }
.stat-label { font-size: 12px; color: var(--text-secondary); }
.profile-actions { display: flex; gap: 12px; }
.profile-actions .btn { flex: 1; padding: 10px; text-align: center; }
</style>
