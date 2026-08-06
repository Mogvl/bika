<template>
  <div class="page-container">
    <div class="section-title">⚙️ 设置</div>

    <div class="settings-group">
      <div class="setting-item">
        <span class="setting-label">修改密码</span>
        <button class="btn" style="background: var(--bg);" @click="showChangePwd = true">修改</button>
      </div>
      <div class="setting-item">
        <span class="setting-label">退出登录</span>
        <button class="btn" style="background: var(--danger); color: white;" @click="logout">退出</button>
      </div>
    </div>

    <div v-if="showChangePwd" class="modal-overlay" @click.self="showChangePwd = false">
      <div class="modal">
        <h3>修改密码</h3>
        <input v-model="oldPwd" type="password" class="input-field" placeholder="旧密码" />
        <input v-model="newPwd" type="password" class="input-field" placeholder="新密码" />
        <p v-if="pwdMsg" class="error-tip">{{ pwdMsg }}</p>
        <div class="modal-actions">
          <button class="btn" style="background: var(--bg);" @click="showChangePwd = false">取消</button>
          <button class="btn btn-primary" @click="changePwd">确认</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { changePassword } from '@/api'

const router = useRouter()
const auth = useAuthStore()
const showChangePwd = ref(false)
const oldPwd = ref('')
const newPwd = ref('')
const pwdMsg = ref('')

async function changePwd() {
  if (!oldPwd.value || !newPwd.value) { pwdMsg.value = '请填写完整'; return }
  try {
    await changePassword(oldPwd.value, newPwd.value)
    alert('密码修改成功，请重新登录')
    logout()
  } catch (e: any) { pwdMsg.value = e.message || '修改失败' }
}

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.settings-group { background: var(--bg-card); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.setting-item { display: flex; align-items: center; justify-content: space-between; padding: 16px; border-bottom: 1px solid var(--border); }
.setting-label { font-size: 15px; }
.setting-item:last-child { border-bottom: none; }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: var(--bg-card); border: 1px solid var(--border); box-shadow: var(--shadow); border-radius: var(--radius); padding: 24px; width: 90%; max-width: 360px; display: flex; flex-direction: column; gap: 12px; }
.modal h3 { margin-bottom: 8px; }
.error-tip { color: var(--danger); font-size: 13px; }
.modal-actions { display: flex; gap: 12px; margin-top: 8px; }
.modal-actions .btn { flex: 1; padding: 10px; }
</style>
