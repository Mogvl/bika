<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <h1 class="login-logo">PicACG</h1>
        <p class="login-subtitle">哔咔漫画 Web 版</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <input
            v-model="email"
            type="text"
            class="input-field"
            placeholder="用户名/邮箱"
            required
            autocomplete="username"
          />
        </div>
        <div class="form-group">
          <input
            v-model="password"
            type="password"
            class="input-field"
            placeholder="密码"
            required
            autocomplete="current-password"
          />
        </div>

        <p v-if="errorMsg" class="error-tip">{{ errorMsg }}</p>

        <button type="submit" class="btn btn-primary login-btn" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <div class="login-footer">
        <div class="login-links">
          <router-link to="/register" class="register-link">没有账号？去注册</router-link>
          <a href="#" class="forgot-link" @click.prevent="showForgot = true">忘记密码？</a>
        </div>
        <p class="disclaimer">
          使用哔咔漫画账号登录<br />
          本项目仅用于技术研究
        </p>
      </div>
    </div>

    <!-- 忘记密码弹窗 -->
    <div v-if="showForgot" class="modal-overlay" @click.self="showForgot = false">
      <div class="modal">
        <h3>重置密码</h3>
        <p class="forgot-tip">通过注册时设置的密保问题重置密码</p>
        <input v-model="forgot.email" type="text" class="input-field" placeholder="邮箱" />
        <select v-model="forgot.questionNo" class="input-field">
          <option :value="1">你最喜欢的漫画角色？</option>
          <option :value="2">你的第一个宠物叫什么？</option>
          <option :value="3">你的家乡在哪里？</option>
        </select>
        <input v-model="forgot.answer" type="text" class="input-field" placeholder="密保答案" />
        <p v-if="forgotMsg" class="error-tip" :class="{ success: forgotOk }">{{ forgotMsg }}</p>
        <div class="modal-actions">
          <button class="btn" style="background: var(--bg);" @click="showForgot = false">取消</button>
          <button class="btn btn-primary" @click="doReset" :disabled="forgotLoading">
            {{ forgotLoading ? '提交中...' : '重置密码' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { resetPassword } from '@/api'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

// 忘记密码
const showForgot = ref(false)
const forgotLoading = ref(false)
const forgotMsg = ref('')
const forgotOk = ref(false)
const forgot = reactive({
  email: '',
  questionNo: 1,
  answer: '',
})

async function handleLogin() {
  if (!email.value || !password.value) {
    errorMsg.value = '请输入用户名和密码'
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    await auth.login(email.value, password.value)
    router.push('/home')
  } catch (e: any) {
    errorMsg.value = e.message || '登录失败，请检查账号密码'
  } finally {
    loading.value = false
  }
}

async function doReset() {
  if (!forgot.email || !forgot.answer) {
    forgotMsg.value = '请填写邮箱和答案'
    forgotOk.value = false
    return
  }
  forgotLoading.value = true
  forgotMsg.value = ''
  try {
    await resetPassword(forgot.email, forgot.questionNo, forgot.answer)
    forgotMsg.value = '校验通过！请查看邮箱设置新密码'
    forgotOk.value = true
  } catch (e: any) {
    forgotMsg.value = e.message || '重置失败'
    forgotOk.value = false
  } finally {
    forgotLoading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fce4ec 0%, #f8bbd0 100%);
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 360px;
  background: white;
  border-radius: 20px;
  padding: 40px 28px;
  box-shadow: 0 8px 32px rgba(231, 76, 139, 0.15);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  font-size: 36px;
  color: var(--primary);
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  width: 100%;
}

.error-tip {
  color: #e74c3c;
  font-size: 13px;
  text-align: center;
}

.login-btn {
  width: 100%;
  padding: 14px;
  font-size: 16px;
  margin-top: 8px;
}

.login-footer {
  margin-top: 24px;
}

.login-links {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.forgot-link {
  font-size: 13px;
  color: var(--text-secondary);
}

.disclaimer {
  font-size: 12px;
  color: #999;
  line-height: 1.6;
}

.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: white; border-radius: var(--radius); padding: 24px; width: 90%; max-width: 360px; display: flex; flex-direction: column; gap: 12px; }
.modal h3 { margin-bottom: 8px; }
.forgot-tip { font-size: 12px; color: var(--text-secondary); margin-bottom: 8px; }
.error-tip { color: #e74c3c; font-size: 13px; }
.error-tip.success { color: #27ae60; }
.modal-actions { display: flex; gap: 12px; margin-top: 8px; }
.modal-actions .btn { flex: 1; padding: 10px; }
</style>
